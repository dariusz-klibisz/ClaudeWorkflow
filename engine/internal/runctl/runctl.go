// Package runctl owns the run state machine (04 §1): run start/branch/adopt/
// close, phase transitions, loops with engine-enforced caps, park/force with
// escalation, and write-time-validated record/approve commands. The engine is
// the only writer of phase transitions; skills instruct, hooks enforce.
package runctl

import (
	"errors"
	"fmt"
	"path/filepath"
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
	return &contracts.Env{Spec: c.Spec, Config: c.Config, Run: r, Events: evs,
		ProjectDir: c.ProjectDir()}, nil
}

// ProjectDir is the directory .workflow/ lives in.
func (c *Ctl) ProjectDir() string {
	if c.Store == nil || c.Store.Root == "" {
		return ""
	}
	return filepath.Dir(c.Store.Root)
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
		// entry contract of the NEXT phase: its inputs must exist before
		// the transition (waivers are the recorded escape; --force skips)
		if next := c.nextUnwaived(r); next != "" {
			entry, err := contracts.EvaluateEntry(env, next)
			if err != nil {
				return nil, "", err
			}
			if len(entry) > 0 {
				return entry, "", nil
			}
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
// A ship-stage loop (spec loops.ship) re-opens Verify on audit-cause
// discoveries — dispositions were the only remedy for a defect the auditor
// found after Verify exited.
func (c *Ctl) Loop(ac, cause, evidence string) (string, error) {
	r, err := c.MustRun()
	if err != nil {
		return "", err
	}
	lp := c.Spec.Loops
	var target string
	if lp.Ship != nil && r.Phase == "ship" && contains(lp.Ship.Causes, cause) {
		if err := c.shipLoopGrounds(r); err != nil {
			return "", err
		}
		target = lp.Ship.Target
	} else {
		if r.Phase != lp.From {
			return "", fmt.Errorf("loops start from %s (current phase %s)", lp.From, r.Phase)
		}
		switch cause {
		case "slip":
			target = "build"
		case "design":
			target = "design"
		case "plan":
			target = "plan"
		default:
			return "", fmt.Errorf("cause must be slip|design|plan (or audit, from ship)")
		}
		if !contains(lp.Targets, target) {
			return "", fmt.Errorf("target %s not a legal loop target", target)
		}
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

// shipLoopGrounds: an audit loop must have something at Ship to loop ON —
// a failing auditor verdict or an open trace-finding. Without grounds the
// loop is a phase-order escape, not a remediation.
func (c *Ctl) shipLoopGrounds(r *store.Run) error {
	env, err := c.Env(r)
	if err != nil {
		return err
	}
	for _, tf := range env.Records("trace-finding") {
		if s, _ := tf.Data["status"].(string); s == "open" {
			return nil
		}
	}
	last := ""
	for _, v := range env.Records("verdict") {
		if a, _ := v.Data["agent"].(string); a == "auditor" {
			last, _ = v.Data["status"].(string)
		}
	}
	if c.Spec != nil {
		for _, f := range c.Spec.Verdicts.Fail {
			if last == f {
				return nil // the LATEST audit failed — grounds to loop
			}
		}
	}
	return errors.New("nothing at ship to loop on: no open trace-finding and no failing auditor verdict (wf trace first)")
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
	case "requirement":
		// AC-less requirements vacuously satisfied every downstream per-AC
		// item (strategy, ac-verdict, red-green) — the sharpest hole in the
		// "built the right thing" chain. Refuse them at write time. Updates
		// that don't touch acs (status flips at Context baseline) pass.
		if acs, present := data["acs"]; present {
			if err := validateACs(acs); err != nil {
				rid, _ := data["rid"].(string)
				return fmt.Errorf("requirement %s: %v (--ac \"AC-1: …\" or acs=[{id,text}])", rid, err)
			}
		}
	case "context-map":
		// judgment records used to be existence-checked only — an empty map
		// was a checkbox. Structural floor here; the depth floor (≥N
		// entries) is the waivable context.map-depth contract item.
		if entries, present := data["entries"]; present {
			if raw, ok := entries.([]any); !ok || len(raw) == 0 {
				return errors.New("context-map requires at least one entry")
			}
		}
		if s, present := data["sufficiency"]; present {
			if str, _ := s.(string); strings.TrimSpace(str) == "" {
				return errors.New("context-map requires a non-empty sufficiency judgment")
			}
		}
	case "completeness":
		if items, present := data["items"]; present {
			raw, ok := items.([]any)
			if !ok || len(raw) == 0 {
				return errors.New("completeness requires at least one case")
			}
			for i, el := range raw {
				m, ok := el.(map[string]any)
				if !ok {
					return fmt.Errorf("completeness item %d must be {case, disposition}", i+1)
				}
				cs, _ := m["case"].(string)
				disp, _ := m["disposition"].(string)
				if strings.TrimSpace(cs) == "" || strings.TrimSpace(disp) == "" {
					return fmt.Errorf("completeness item %d needs both case and disposition", i+1)
				}
			}
		}
	case "option-set":
		// 03 §4.3 promised "the engine cross-checks rejected IDs — never
		// re-proposed"; until now only design-reviewer prose enforced it.
		if err := c.validateOptionSet(r, data); err != nil {
			return err
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
		// n/a self-attests inapplicability: the reason is part of the
		// verdict (agent prose always demanded it; the SubagentStop gate
		// blocks reasonless auto n/a before it ever reaches here)
		if status == "n/a" {
			if s, _ := data["reason"].(string); strings.TrimSpace(s) == "" {
				return errors.New("n/a verdict requires reason=\"why this review does not apply\"")
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
	case "metric":
		if _, ok := data["grounded"]; !ok {
			data["grounded"] = auto // manual metrics are self-attested (reported as such)
		}
		// below_threshold is engine-computed, never self-declared: the value
		// may be a claim, the comparison is mechanical
		if name, _ := data["name"].(string); name != "" && c.Config != nil {
			if th := numOr(c.Config.Thresholds[name], -1); th >= 0 {
				if v := numOr(data["value"], -1); v >= 0 {
					data["below_threshold"] = v < th
				}
			}
		}
	case "finding":
		if fid, _ := data["fid"].(string); strings.TrimSpace(fid) == "" {
			return errors.New("finding requires a non-empty fid (it must appear verbatim in the report)")
		}
		if sev, present := data["severity"]; present {
			switch sev {
			case "critical", "major", "minor", "info":
			default:
				return fmt.Errorf("finding severity must be critical|major|minor|info, got %q", sev)
			}
		}
		if txt, present := data["text"]; present {
			if s, _ := txt.(string); strings.TrimSpace(s) == "" {
				return errors.New("finding requires non-empty text")
			}
		}
	case "waiver":
		if s, _ := data["reason"].(string); strings.TrimSpace(s) == "" {
			return errors.New("waiver requires --reason")
		}
	case "artifact":
		// `status: present` is a claim the disk must confirm (write-time,
		// the stronger layer): the file exists and is authored — not the
		// untouched template, not a skeleton.
		if s, _ := data["status"].(string); s == "present" {
			rel, tmpl := c.artifactPathAndTemplate(r, data)
			if rel == "" {
				return errors.New("artifact status=present requires a path (on the record or its original)")
			}
			tmplPath := ""
			if tmpl != "" && c.Spec.PluginRoot() != "" {
				tmplPath = filepath.Join(c.Spec.PluginRoot(), "templates", tmpl+".md")
			}
			if ok, detail := contracts.ArtifactOnDisk(c.ProjectDir(), rel, tmplPath); !ok {
				return fmt.Errorf("artifact status=present refused: %s", detail)
			}
		}
	}
	return nil
}

// artifactPathAndTemplate resolves path+template for an artifact write,
// following `updates=` to the original record when the update omits them.
func (c *Ctl) artifactPathAndTemplate(r *store.Run, data map[string]any) (string, string) {
	rel, _ := data["path"].(string)
	tmpl, _ := data["template"].(string)
	if rel != "" && tmpl != "" {
		return rel, tmpl
	}
	target, _ := data["updates"].(string)
	if target == "" {
		return rel, tmpl
	}
	env, err := c.Env(r)
	if err != nil {
		return rel, tmpl
	}
	for _, a := range env.Records("artifact") {
		if a.ID == target {
			if rel == "" {
				rel, _ = a.Data["path"].(string)
			}
			if tmpl == "" {
				tmpl, _ = a.Data["template"].(string)
			}
			break
		}
	}
	return rel, tmpl
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
// layer harder to fabricate. Two opt-in friction levels (09 Q3):
//   - `approvals: hardened` — an approval WITHOUT such an answer is refused.
//   - `approvals: challenge` — the anchoring answer must additionally
//     contain a single-use code the engine shows ONLY in the user's
//     statusline (never on stdout): the model cannot manufacture its own
//     anchor because it never sees the code before the user types it.
func (c *Ctl) Approve(gate, payload string) (*store.Event, error) {
	r, err := c.MustRun()
	if err != nil {
		return nil, err
	}
	data := map[string]any{"gate": gate}
	if payload != "" {
		data["payload_hash"] = fmt.Sprintf("%x", hash(payload))
	}
	mode := ""
	if c.Config != nil {
		mode, _ = c.Config.ConfigFlag("approvals").(string)
	}
	if env, err := c.Env(r); err == nil {
		// bind the approval to WHAT is being approved: the engine computes
		// the record identities in scope right now (requirement rids,
		// selected options, task tids…) — the agent-supplied payload was a
		// free string proving nothing. Trace re-computes and reports drift.
		if refs := contracts.ApprovalRefs(env, gate); len(refs) > 0 {
			anyRefs := make([]any, len(refs)) // JSON round-trip shape
			for i, ref := range refs {
				anyRefs[i] = ref
			}
			data["approved_refs"] = anyRefs
			data["refs_hash"] = fmt.Sprintf("%x", hash(strings.Join(refs, "\n")))
		}
		lastApproval := -1
		for _, a := range env.Records("approval") {
			if g, _ := a.Data["gate"].(string); g == gate && a.Order > lastApproval {
				lastApproval = a.Order
			}
		}
		// topic-aware anchoring: prefer the newest answer ABOUT this gate;
		// tolerate topicless answers; never link an answer about a
		// DIFFERENT gate (the anchoring race: any unrelated AskUserQuestion
		// used to anchor a hardened approval)
		type cand struct{ id, answer string }
		var topicCands, plainCands []cand // stream order: newest last
		for _, ua := range env.Records("user-answer") {
			if ua.Order <= lastApproval {
				continue
			}
			answer, _ := ua.Data["answer"].(string)
			switch topic, _ := ua.Data["topic"].(string); topic {
			case gate:
				topicCands = append(topicCands, cand{ua.ID, answer})
			case "":
				plainCands = append(plainCands, cand{ua.ID, answer})
			}
		}
		if mode == "challenge" {
			// single-use code, statusline-delivered: verify it appears in a
			// captured answer, newest-first, topic-matched preferred
			code, cerr := c.ensureChallenge(gate)
			if cerr != nil {
				return nil, cerr
			}
			ref := ""
			for _, cs := range [][]cand{topicCands, plainCands} {
				for i := len(cs) - 1; i >= 0; i-- {
					if strings.Contains(strings.ToUpper(cs[i].answer), code) {
						ref = cs[i].id
						break
					}
				}
				if ref != "" {
					break
				}
			}
			if ref == "" {
				return nil, fmt.Errorf("approvals require the user's challenge code for this project: ask via AskUserQuestion (naming the %s gate) for the code shown in the user's wf statusline, then re-run wf approve %s", gate, gate)
			}
			c.clearChallenge() // single-use: consumed on success
			data["answer_ref"] = ref
			data["challenge"] = true
			return c.append(r, "approval", false, "user", data)
		}
		ref := ""
		if len(topicCands) > 0 {
			ref = topicCands[len(topicCands)-1].id
		} else if len(plainCands) > 0 {
			ref = plainCands[len(plainCands)-1].id
		}
		switch {
		case ref != "":
			data["answer_ref"] = ref
		case mode == "hardened":
			return nil, fmt.Errorf("approvals are hardened for this project: pose the %s question via AskUserQuestion first (the hook records the answer; phrase the question so it names the %s gate), then re-run wf approve %s", gate, gate, gate)
		}
	}
	return c.append(r, "approval", false, "user", data)
}

// WaiveItem records a contract-item (or per-each element) waiver.
func (c *Ctl) WaiveItem(item, reason string) (*store.Event, error) {
	return c.Record("waiver", map[string]any{"item": item, "reason": reason}, false, "user")
}

// validateOptionSet enforces genuine option evaluation at write time:
// a known stage, ≥2 candidates, a non-empty selection — and the selection
// must not be an option a prior option-set of the same stage REJECTED,
// unless a disposition references that prior record (the recorded escape
// when a rejection genuinely no longer holds).
func (c *Ctl) validateOptionSet(r *store.Run, data map[string]any) error {
	if stage, present := data["stage"]; present {
		s, _ := stage.(string)
		switch s {
		case "system", "software", "ux":
		default:
			return fmt.Errorf("option-set stage must be system|software|ux, got %q", s)
		}
	}
	if cands, present := data["candidates"]; present {
		raw, ok := cands.([]any)
		if !ok || len(raw) < 2 {
			return errors.New("option-set requires ≥2 genuine candidates (a single option is a decision, not an evaluation)")
		}
	}
	selected, hasSel := data["selected"].(string)
	if hasSel && strings.TrimSpace(selected) == "" {
		return errors.New("option-set requires a non-empty selected option")
	}
	if !hasSel {
		return nil // no selection in this write — nothing to cross-check
	}
	stage, _ := data["stage"].(string)
	env, err := c.Env(r)
	if err != nil {
		return err
	}
	dispositioned := map[string]bool{}
	for _, d := range env.Records("disposition") {
		if ref, _ := d.Data["ref"].(string); ref != "" {
			dispositioned[ref] = true
		}
	}
	for _, prior := range env.Records("option-set") {
		if s, _ := prior.Data["stage"].(string); stage != "" && s != stage {
			continue
		}
		if dispositioned[prior.ID] {
			continue
		}
		for _, rej := range rejectedIDs(prior.Data["rejected"]) {
			if strings.EqualFold(rej, selected) {
				return fmt.Errorf("option %q was rejected in option-set %s — select another, or record a disposition referencing it first (wf record disposition --ref %s --text \"why the rejection no longer holds\")",
					selected, prior.ID, prior.ID)
			}
		}
	}
	return nil
}

// rejectedIDs flattens an option-set's rejected list: strings or {id,…}.
func rejectedIDs(v any) []string {
	raw, _ := v.([]any)
	var out []string
	for _, el := range raw {
		switch t := el.(type) {
		case string:
			if t != "" {
				out = append(out, t)
			}
		case map[string]any:
			if id, _ := t["id"].(string); id != "" {
				out = append(out, id)
			}
		}
	}
	return out
}

// validateACs checks a requirement's acs payload: a non-empty array whose
// elements are non-empty strings ("AC-1: …") or maps with a non-empty id —
// the same shapes contracts.elementIDs iterates per-each over.
func validateACs(v any) error {
	raw, ok := v.([]any)
	if !ok || len(raw) == 0 {
		return errors.New("requires at least one AC")
	}
	for i, el := range raw {
		switch t := el.(type) {
		case string:
			if strings.TrimSpace(t) == "" {
				return fmt.Errorf("AC %d is empty", i+1)
			}
		case map[string]any:
			if id, _ := t["id"].(string); strings.TrimSpace(id) == "" {
				return fmt.Errorf("AC %d has no id", i+1)
			}
		default:
			return fmt.Errorf("AC %d must be a string or {id,text} object", i+1)
		}
	}
	return nil
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
