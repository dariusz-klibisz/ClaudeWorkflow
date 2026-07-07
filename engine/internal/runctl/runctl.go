// Package runctl owns the run state machine (04 §1): run start/branch/adopt/
// close, phase transitions, loops with engine-enforced caps, park/force with
// escalation, and write-time-validated record/approve commands. The engine is
// the only writer of phase transitions; skills instruct, hooks enforce.
package runctl

import (
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/contracts"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/spec"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/store"
)

var (
	ErrNoRun     = errors.New("no active run (wf run start, or /wf:dev)")
	ErrRunActive = errors.New("a run is already active (close, park, or branch it first)")
)

type Ctl struct {
	Store  *store.Store
	Spec   *spec.Spec
	Config *store.Config
}

func (c *Ctl) MustRun() (*store.Run, error) {
	r, err := c.Store.LoadRun()
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, ErrNoRun
	}
	return r, nil
}

// Env builds a contract-evaluation environment for the current run.
func (c *Ctl) Env(r *store.Run) (*contracts.Env, error) {
	evs, err := c.Store.RunEvents(r.ID)
	if err != nil {
		return nil, err
	}
	return &contracts.Env{Spec: c.Spec, Config: c.Config, Run: r, Events: evs}, nil
}

func (c *Ctl) append(r *store.Run, kind string, auto bool, actor string, data map[string]any) (*store.Event, error) {
	ev := &store.Event{Kind: kind, Auto: auto, Actor: actor, Data: data}
	if r != nil {
		ev.Run = r.ID
		ev.Phase = r.Phase
	}
	if err := c.Store.Append(ev); err != nil {
		return nil, err
	}
	return ev, nil
}

// ---------------------------------------------------------------------------
// Run lifecycle
// ---------------------------------------------------------------------------

func mintRunID() string {
	return time.Now().UTC().Format("20060102") + "-" + strings.ToLower(store.NewULID()[18:])
}

func (c *Ctl) RunStart(family, intent string) (*store.Run, error) {
	if cur, _ := c.Store.LoadRun(); cur != nil && cur.Status == "active" {
		return nil, ErrRunActive
	}
	if !c.Spec.ValidFamily(family) {
		return nil, fmt.Errorf("unknown family %q (one of %s)", family, strings.Join(c.Spec.Families, "|"))
	}
	if intent != "" && !c.Spec.ValidIntent(family, intent) {
		return nil, fmt.Errorf("intent %q not valid for family %s (one of %s)", intent, family, strings.Join(c.Spec.Intents[family], "|"))
	}
	r := &store.Run{
		ID: mintRunID(), Family: family, Intent: intent, Status: "active",
		Started: time.Now().UTC().Format(time.RFC3339), SlipByAC: map[string]int{},
	}
	first := c.Spec.PhasesFor(family)[0]
	r.Phase = first.ID
	if _, err := c.append(r, "run", false, "engine", map[string]any{"action": "start", "family": family, "intent": intent}); err != nil {
		return nil, err
	}
	if _, err := c.append(r, "phase", false, "engine", map[string]any{"action": "enter", "target": first.ID}); err != nil {
		return nil, err
	}
	return r, c.Store.SaveRun(r)
}

// RunBranch parks the parent and starts a child run carrying lineage.
func (c *Ctl) RunBranch(family, intent, reason string) (*store.Run, error) {
	parent, err := c.MustRun()
	if err != nil {
		return nil, err
	}
	if _, err := c.append(parent, "phase", false, "engine", map[string]any{"action": "park", "reason": "branched: " + reason}); err != nil {
		return nil, err
	}
	parent.Status = "parked"
	if err := c.Store.SaveRun(parent); err != nil {
		return nil, err
	}
	if family == "" {
		family = parent.Family
	}
	if intent == "" {
		intent = parent.Intent
	}
	child, err := c.RunStart(family, intent)
	if err != nil {
		return nil, err
	}
	child.Parent = parent.ID
	if _, err := c.append(child, "run", false, "engine", map[string]any{"action": "branch", "family": family, "intent": intent, "parent": parent.ID, "reason": reason}); err != nil {
		return nil, err
	}
	return child, c.Store.SaveRun(child)
}

// RunAdopt re-derives the snapshot from the committed log (fresh clone /
// second machine — the G2 fix).
func (c *Ctl) RunAdopt() (*store.Run, error) {
	r, err := c.Store.DeriveRun()
	if err != nil {
		return nil, err
	}
	if r == nil {
		return nil, errors.New("no in-flight run found in the event log")
	}
	if _, err := c.append(r, "run", false, "engine", map[string]any{"action": "adopt"}); err != nil {
		return nil, err
	}
	return r, c.Store.SaveRun(r)
}

// RunClose is the single atomic close transaction (A5/A6/G5 fix): terminal
// event → archive → verify → compact → clear.
func (c *Ctl) RunClose() error {
	r, err := c.MustRun()
	if err != nil {
		return err
	}
	last := c.Spec.PhasesFor(r.Family)
	lastID := last[len(last)-1].ID
	if !contains(r.ExitedPh, lastID) {
		return fmt.Errorf("cannot close: phase %s not exited (wf phase exit)", lastID)
	}
	if _, err := c.append(r, "run", false, "engine", map[string]any{"action": "close"}); err != nil {
		return err
	}
	return c.Store.ArchiveRun(r.ID)
}

// ---------------------------------------------------------------------------
// Phase transitions
// ---------------------------------------------------------------------------

// PhaseExit evaluates the current phase contract. On success it advances (or
// reports the run closeable); on gaps it returns findings (CLI → exit 2).
func (c *Ctl) PhaseExit(force bool, reason string) ([]contracts.Finding, string, error) {
	r, err := c.MustRun()
	if err != nil {
		return nil, "", err
	}
	if r.Status != "active" {
		return nil, "", fmt.Errorf("run is %s (wf run resume)", r.Status)
	}
	if !force {
		env, err := c.Env(r)
		if err != nil {
			return nil, "", err
		}
		findings, err := contracts.EvaluatePhase(env, r.Phase)
		if err != nil {
			return nil, "", err // broken contract → exit 3 at CLI
		}
		if len(findings) > 0 {
			return findings, "", nil
		}
		if _, err := c.append(r, "phase", false, "engine", map[string]any{"action": "exit"}); err != nil {
			return nil, "", err
		}
	} else {
		if err := c.forceEscalation(r, reason); err != nil {
			return nil, "", err
		}
		if _, err := c.append(r, "phase", false, "engine", map[string]any{"action": "force", "reason": reason}); err != nil {
			return nil, "", err
		}
		r.Forces++
	}
	r.ExitedPh = append(r.ExitedPh, r.Phase)
	next := c.nextUnwaived(r)
	if next == "" {
		r.Phase = ""
		if err := c.Store.SaveRun(r); err != nil {
			return nil, "", err
		}
		return nil, "all phases complete — wf run close", nil
	}
	r.Phase = next
	if _, err := c.append(r, "phase", false, "engine", map[string]any{"action": "enter", "target": next}); err != nil {
		return nil, "", err
	}
	return nil, "entered " + next, c.Store.SaveRun(r)
}

func (c *Ctl) nextUnwaived(r *store.Run) string {
	next := c.Spec.NextPhase(r.Family, r.Phase)
	for next != "" && contains(r.WaivedPh, next) {
		next = c.Spec.NextPhase(r.Family, next)
	}
	return next
}

// forceEscalation implements the G4 fix: 1st force needs a reason; the 2nd a
// structural cause; the 3rd auto-parks with a repair checklist.
func (c *Ctl) forceEscalation(r *store.Run, reason string) error {
	if strings.TrimSpace(reason) == "" {
		return errors.New("force requires --reason")
	}
	switch {
	case r.Forces >= 2:
		_, _ = c.append(r, "escape", false, "user", map[string]any{"action": "park", "reason": "3rd force in run — auto-parked", "level": 3})
		_, _ = c.append(r, "phase", false, "engine", map[string]any{"action": "park", "reason": "force escalation"})
		r.Status = "parked"
		_ = c.Store.SaveRun(r)
		return errors.New("3rd force in this run: auto-parked. Repair checklist: (1) wf status — read the unmet items; (2) fix or waive each with a reason; (3) wf run resume. Forcing is no longer available for this run")
	case r.Forces == 1:
		if !strings.Contains(reason, "cause:") {
			return errors.New("2nd force in this run requires naming the structural cause: --reason \"cause: <what makes this gate wrong here>\"")
		}
	}
	_, err := c.append(r, "escape", false, "user", map[string]any{"action": "force", "reason": reason, "level": r.Forces + 1})
	return err
}

// PhaseWaive records skipping a waivable phase (family-contract-checked).
func (c *Ctl) PhaseWaive(phase, reason string) error {
	r, err := c.MustRun()
	if err != nil {
		return err
	}
	p, ok := c.Spec.PhaseByID(phase)
	if !ok {
		return fmt.Errorf("unknown phase %q", phase)
	}
	if p.Participation[r.Family] != "waivable" {
		return fmt.Errorf("phase %s is not waivable for family %s", phase, r.Family)
	}
	if strings.TrimSpace(reason) == "" {
		return errors.New("waiving a phase requires --reason")
	}
	if contains(r.ExitedPh, phase) {
		return fmt.Errorf("phase %s already exited", phase)
	}
	if _, err := c.append(r, "phase", false, "user", map[string]any{"action": "waive", "target": phase, "reason": reason}); err != nil {
		return err
	}
	r.WaivedPh = append(r.WaivedPh, phase)
	if r.Phase == phase {
		r.Phase = c.nextUnwaived(r)
		if _, err := c.append(r, "phase", false, "engine", map[string]any{"action": "enter", "target": r.Phase}); err != nil {
			return err
		}
	}
	return c.Store.SaveRun(r)
}

// Loop re-opens a target phase from verify with engine-enforced caps (03 §6).
func (c *Ctl) Loop(ac, cause, evidence string) (string, error) {
	r, err := c.MustRun()
	if err != nil {
		return "", err
	}
	lp := c.Spec.Loops
	if r.Phase != lp.From {
		return "", fmt.Errorf("loops start from %s (current phase %s)", lp.From, r.Phase)
	}
	var target string
	switch cause {
	case "slip":
		target = "build"
	case "design":
		target = "design"
	case "plan":
		target = "plan"
	default:
		return "", fmt.Errorf("cause must be slip|design|plan")
	}
	if !contains(lp.Targets, target) {
		return "", fmt.Errorf("target %s not a legal loop target", target)
	}
	if strings.TrimSpace(evidence) == "" {
		return "", errors.New("loop requires --evidence (discriminating: observed vs expected)")
	}
	if r.Loops >= lp.MaxPerRun {
		_, _ = c.append(r, "phase", false, "engine", map[string]any{"action": "park", "reason": "loop cap reached"})
		r.Status = "parked"
		_ = c.Store.SaveRun(r)
		return "", fmt.Errorf("loop cap (%d/run) reached: run auto-parked. Record the root cause and branch or resume deliberately", lp.MaxPerRun)
	}
	if cause == "slip" && r.SlipByAC[ac] >= lp.MaxSlipPerAC {
		return "", fmt.Errorf("%d slip-loops already on %s: the defect is structural — use --cause design or --cause plan", lp.MaxSlipPerAC, ac)
	}
	if _, err := c.append(r, "loop", false, "agent", map[string]any{"ac": ac, "cause": cause, "evidence": evidence, "target": target}); err != nil {
		return "", err
	}
	if _, err := c.append(r, "phase", false, "engine", map[string]any{"action": "loop", "target": target, "reason": cause}); err != nil {
		return "", err
	}
	r.Loops++
	if cause == "slip" {
		if r.SlipByAC == nil {
			r.SlipByAC = map[string]int{}
		}
		r.SlipByAC[ac]++
	}
	r.ExitedPh = remove(r.ExitedPh, target)
	r.Phase = target
	return target, c.Store.SaveRun(r)
}

// Park is the always-available honest stop; it clears every sequencing gate.
func (c *Ctl) Park(reason string) error {
	r, err := c.MustRun()
	if err != nil {
		return err
	}
	if strings.TrimSpace(reason) == "" {
		return errors.New("park requires --reason")
	}
	if _, err := c.append(r, "escape", false, "user", map[string]any{"action": "park", "reason": reason}); err != nil {
		return err
	}
	if _, err := c.append(r, "phase", false, "engine", map[string]any{"action": "park", "reason": reason}); err != nil {
		return err
	}
	r.Status = "parked"
	return c.Store.SaveRun(r)
}

func (c *Ctl) Resume() (*store.Run, error) {
	r, err := c.MustRun()
	if err != nil {
		return nil, err
	}
	if r.Status != "parked" {
		return nil, fmt.Errorf("run is %s, not parked", r.Status)
	}
	if _, err := c.append(r, "phase", false, "engine", map[string]any{"action": "resume"}); err != nil {
		return nil, err
	}
	r.Status = "active"
	return r, c.Store.SaveRun(r)
}

// ---------------------------------------------------------------------------
// Records and approvals (write-time validation — the stronger layer)
// ---------------------------------------------------------------------------

// Record validates and appends a record of a spec-declared kind.
func (c *Ctl) Record(kind string, data map[string]any, auto bool, actor string) (*store.Event, error) {
	r, err := c.MustRun()
	if err != nil {
		return nil, err
	}
	rk, ok := c.Spec.RecordKind(kind)
	if !ok {
		return nil, fmt.Errorf("unknown record kind %q", kind)
	}
	if target, isUpdate := data["updates"].(string); isUpdate {
		// an update must reference a known record (original or one of its
		// update events — resolved transitively). Silent no-ops broke the
		// power5 run's task bindings.
		env, err := c.Env(r)
		if err != nil {
			return nil, err
		}
		orig, ok := env.ResolveRecordID(target)
		if !ok {
			return nil, fmt.Errorf("updates=%s does not reference any record in this run — target the record's creation ID (wf status / the ID printed when it was recorded)", target)
		}
		data["updates"] = orig // normalize chains to the original at write time
	} else {
		for _, f := range rk.Required() {
			if _, ok := data[f]; !ok {
				return nil, fmt.Errorf("%s requires field %q (fields: %s)", kind, f, strings.Join(rk.Fields, ", "))
			}
		}
	}
	if err := c.validateRecord(r, kind, data, auto); err != nil {
		return nil, err
	}
	if actor == "" {
		actor = "agent"
	}
	return c.append(r, kind, auto, actor, data)
}

func (c *Ctl) validateRecord(r *store.Run, kind string, data map[string]any, auto bool) error {
	switch kind {
	case "classification":
		fam, _ := data["family"].(string)
		intent, _ := data["intent"].(string)
		if !c.Spec.ValidFamily(fam) {
			return fmt.Errorf("unknown family %q", fam)
		}
		if !c.Spec.ValidIntent(fam, intent) {
			return fmt.Errorf("intent %q invalid for family %s", intent, fam)
		}
		// keep the snapshot in sync when Frame corrects the provisional call
		if fam != r.Family || intent != r.Intent {
			r.Family, r.Intent = fam, intent
			if err := c.Store.SaveRun(r); err != nil {
				return err
			}
		}
	case "verdict":
		// agent must be a roster name (scoped "wf:" prefixes are normalized;
		// unknown names silently failed contracts in the power5 run)
		agent, _ := data["agent"].(string)
		agent = strings.TrimPrefix(agent, "wf:")
		if _, ok := c.Spec.AgentByName(agent); !ok {
			names := make([]string, 0, len(c.Spec.Roster))
			for _, a := range c.Spec.Roster {
				names = append(names, a.Name)
			}
			return fmt.Errorf("verdict agent %q is not in the roster (%s)", agent, strings.Join(names, ", "))
		}
		data["agent"] = agent
		status, _ := data["status"].(string)
		if !c.verdictKnown(status) {
			return fmt.Errorf("verdict status %q not in vocabulary", status)
		}
		crit := numOr(data["criticals"], -1)
		maj := numOr(data["majors"], -1)
		if crit < 0 || maj < 0 {
			return errors.New("verdict requires integer criticals and majors")
		}
		if (status == "clean" || status == "safe") && crit+maj > 0 {
			if auto {
				data["status"] = "changes-required" // auto-downgrade (04 §4)
				data["downgraded"] = true
			} else {
				return fmt.Errorf("%s verdict with criticals=%v majors=%v is contradictory", status, crit, maj)
			}
		}
	case "ac-verdict":
		status, _ := data["status"].(string)
		ac, _ := data["ac"].(string)
		switch status {
		case "pass":
			if r.Family == "diff" && !c.hasGroundedGreen(r, ac) {
				return fmt.Errorf("AC %s pass refused: no grounded green test-run tagged --ac %s (run the test, or wf capture/record test)", ac, ac)
			}
		case "fail", "deferred":
			if status == "deferred" && !c.hasApproval(r, "deferral") {
				return fmt.Errorf("deferring %s requires user approval first: wf approve deferral", ac)
			}
		default:
			return fmt.Errorf("ac-verdict status must be pass|fail|deferred")
		}
	case "test-run":
		if _, ok := data["grounded"]; !ok {
			data["grounded"] = auto // manual records are ungrounded unless captured
		}
	case "waiver":
		if s, _ := data["reason"].(string); strings.TrimSpace(s) == "" {
			return errors.New("waiver requires --reason")
		}
	}
	return nil
}

func (c *Ctl) verdictKnown(s string) bool {
	for _, v := range append(append(append([]string{}, c.Spec.Verdicts.Pass...), c.Spec.Verdicts.Fail...), c.Spec.Verdicts.Conditional...) {
		if v == s {
			return true
		}
	}
	return false
}

func (c *Ctl) hasGroundedGreen(r *store.Run, ac string) bool {
	env, err := c.Env(r)
	if err != nil {
		return false
	}
	for _, tr := range env.Records("test-run") {
		if a, _ := tr.Data["ac"].(string); a != ac {
			continue
		}
		if g, _ := tr.Data["grounded"].(bool); !g {
			continue
		}
		if exit := numOr(tr.Data["exit"], -1); exit == 0 {
			return true
		}
	}
	return false
}

func (c *Ctl) hasApproval(r *store.Run, gate string) bool {
	env, err := c.Env(r)
	if err != nil {
		return false
	}
	for _, a := range env.Records("approval") {
		if g, _ := a.Data["gate"].(string); g == gate {
			return true
		}
	}
	return false
}

// Approve records a user approval — always auto:false, always self-attested
// and reported (honest bounds, 04 §8). Anchoring (04 §8.1): when a
// hook-captured user-answer exists after this gate's previous approval, it
// is linked via answer_ref — still not proof a human typed it, but one
// layer harder to fabricate. With config `approvals: hardened`, an approval
// WITHOUT such an answer is refused (opt-in friction, 09 Q3).
func (c *Ctl) Approve(gate, payload string) (*store.Event, error) {
	r, err := c.MustRun()
	if err != nil {
		return nil, err
	}
	data := map[string]any{"gate": gate}
	if payload != "" {
		data["payload_hash"] = fmt.Sprintf("%x", hash(payload))
	}
	if env, err := c.Env(r); err == nil {
		lastApproval := -1
		for _, a := range env.Records("approval") {
			if g, _ := a.Data["gate"].(string); g == gate && a.Order > lastApproval {
				lastApproval = a.Order
			}
		}
		ref := ""
		for _, ua := range env.Records("user-answer") {
			if ua.Order > lastApproval {
				ref = ua.ID // stream order: newest wins
			}
		}
		switch {
		case ref != "":
			data["answer_ref"] = ref
		case c.Config != nil && c.Config.ConfigFlag("approvals") == "hardened":
			return nil, fmt.Errorf("approvals are hardened for this project: pose the %s question via AskUserQuestion first (the hook records the answer), then re-run wf approve %s", gate, gate)
		}
	}
	return c.append(r, "approval", false, "user", data)
}

// WaiveItem records a contract-item (or per-each element) waiver.
func (c *Ctl) WaiveItem(item, reason string) (*store.Event, error) {
	return c.Record("waiver", map[string]any{"item": item, "reason": reason}, false, "user")
}

// ---------------------------------------------------------------------------

func contains(xs []string, s string) bool {
	for _, x := range xs {
		if x == s {
			return true
		}
	}
	return false
}

func remove(xs []string, s string) []string {
	out := xs[:0]
	for _, x := range xs {
		if x != s {
			out = append(out, x)
		}
	}
	return out
}

func numOr(v any, def float64) float64 {
	switch t := v.(type) {
	case float64:
		return t
	case int:
		return float64(t)
	}
	return def
}

func hash(s string) uint64 {
	var h uint64 = 14695981039346656037
	for i := 0; i < len(s); i++ {
		h ^= uint64(s[i])
		h *= 1099511628211
	}
	return h
}
