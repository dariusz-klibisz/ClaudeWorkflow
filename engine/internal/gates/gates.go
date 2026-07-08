// Package gates implements the four-gate enforcement spine (04): the Stop
// gate, the task gates (TaskCreated/TaskCompleted), the SubagentStop verdict
// gate, and the PreToolUse tool gates (skill sequencing, stray-edit guard,
// catastrophic Bash net) plus PostToolUse capture.
//
// Fail-safe split (04 §7): gates that *sequence* (stop, skill, edit) fail
// open with a loud systemMessage when the engine is unhealthy; gates that
// *record or protect data* (task-complete, verdict) fail closed.
package gates

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/cmdid"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/contracts"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/hookio"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/runctl"
)

// ---------------------------------------------------------------------------
// Gate 1 — Stop (04 §2)
// ---------------------------------------------------------------------------

type stopState struct {
	SessionID string `json:"session_id"`
	Hash      string `json:"hash"`
	Count     int    `json:"count"`
}

const stopSelfCap = 3 // self-imposed, under the platform's 8

func Stop(c *runctl.Ctl, in *hookio.Input) hookio.Result {
	if hookio.EnforceDisabled(in) {
		recordEnforceOff(c, in, "stop")
		return hookio.StopAllowMessage("[wf] WF_ENFORCE=0 — Stop gate downgraded to a warning (recorded)")
	}
	r, err := c.Store.LoadRun()
	if err != nil {
		return hookio.BrokenGate(err)
	}
	if r == nil || r.Status != "active" {
		return hookio.Allow()
	}
	if r.Phase == "" {
		return stopBlockCounted(c, in, []string{"run.close"},
			"All phases are complete but the run is open. Close it: wf run close")
	}
	env, err := c.Env(r)
	if err != nil {
		return hookio.BrokenGate(err)
	}
	findings, err := contracts.EvaluatePhase(env, r.Phase)
	if err != nil {
		return hookio.BrokenGate(err)
	}
	var agentItems, userItems []contracts.Finding
	for _, f := range findings {
		if f.UserBlocked {
			userItems = append(userItems, f)
		} else {
			agentItems = append(agentItems, f)
		}
	}
	if len(findings) == 0 {
		return stopBlockCounted(c, in, []string{"phase.exit"},
			fmt.Sprintf("Phase %s contract is met. Exit it before stopping: wf phase exit", r.Phase))
	}
	if len(agentItems) == 0 {
		return hookio.StopAllowMessage("[wf] waiting on the user: " + userItems[0].Remediation)
	}
	// a gating reviewer still running in the background is progress in flight
	if reviewerInFlight(c, in) {
		return hookio.StopAllowMessage("[wf] gating reviewer still running in the background — verdict will be captured at SubagentStop")
	}
	ids := make([]string, 0, len(agentItems))
	var b strings.Builder
	fmt.Fprintf(&b, "The %s phase has %d unmet contract item(s) you can progress without the user:\n", r.Phase, len(agentItems))
	for i, f := range agentItems {
		ids = append(ids, f.ID)
		if i >= 5 {
			fmt.Fprintf(&b, "… and %d more (wf status)\n", len(agentItems)-5)
			break
		}
		line := f.Remediation
		if f.Detail != "" {
			line += " [" + f.Detail + "]"
		}
		fmt.Fprintf(&b, "%d. %s → %s\n", i+1, f.ID, line)
	}
	b.WriteString("Honest stops: /wf:park (reason recorded) if genuinely blocked.")
	return stopBlockCounted(c, in, ids, b.String())
}

// stopBlockCounted blocks, but allows after stopSelfCap identical blocks in
// one session (prevents burn loops on a genuinely stuck item).
func stopBlockCounted(c *runctl.Ctl, in *hookio.Input, ids []string, reason string) hookio.Result {
	sort.Strings(ids)
	hash := strings.Join(ids, ",")
	var st stopState
	_ = c.Store.LoadLocal("stop-gate.json", &st)
	if st.SessionID == in.SessionID && st.Hash == hash {
		st.Count++
	} else {
		st = stopState{SessionID: in.SessionID, Hash: hash, Count: 1}
	}
	_ = c.Store.SaveLocal("stop-gate.json", &st)
	if in.StopHookActive && st.Count > stopSelfCap {
		// The release valve is itself an escape — record it (once per
		// identical finding set) so wf report surfaces it, never silence.
		if st.Count == stopSelfCap+1 {
			_, _ = c.Record("escape", map[string]any{
				"action": "stop-cap",
				"reason": "stop-gate self-cap released after " + fmt.Sprint(stopSelfCap) + " identical blocks: " + hash,
			}, true, "hook")
		}
		return hookio.StopAllowMessage("[wf] the same items blocked " + fmt.Sprint(stopSelfCap) + "× — allowing the stop (recorded). If stuck, /wf:park records the honest state")
	}
	return hookio.StopBlock(reason)
}

// enforceOffState rate-limits enforce-off escape events to one per gate per
// session — degraded sequencing gates fire on every tool call.
type enforceOffState struct {
	SessionID string          `json:"session_id"`
	Gates     map[string]bool `json:"gates"`
}

// recordEnforceOff appends the escape event the WF_ENFORCE=0 downgrade paths
// promised ("recorded") but never wrote. The provenance guard already proved
// the env came from the user; the hook is the observer, hence actor "hook".
func recordEnforceOff(c *runctl.Ctl, in *hookio.Input, gate string) {
	var st enforceOffState
	_ = c.Store.LoadLocal("enforce-off.json", &st)
	if st.SessionID != in.SessionID || st.Gates == nil {
		st = enforceOffState{SessionID: in.SessionID, Gates: map[string]bool{}}
	}
	if st.Gates[gate] {
		return
	}
	if r, err := c.Store.LoadRun(); err != nil || r == nil || r.Status != "active" {
		return // no run to attach the event to — nothing being escaped
	}
	if _, err := c.Record("escape", map[string]any{
		"action": "enforce-off", "reason": "WF_ENFORCE=0 set by the user", "gate": gate,
	}, true, "hook"); err != nil {
		return
	}
	st.Gates[gate] = true
	_ = c.Store.SaveLocal("enforce-off.json", &st)
}

func reviewerInFlight(c *runctl.Ctl, in *hookio.Input) bool {
	if len(in.BackgroundTasks) == 0 {
		return false
	}
	raw := string(in.BackgroundTasks)
	for _, a := range c.Spec.GatingAgents() {
		if strings.Contains(raw, a.Name) {
			return true
		}
	}
	return false
}

// ---------------------------------------------------------------------------
// Gate 2 — task gates (04 §3)
// ---------------------------------------------------------------------------

type taskMirror struct {
	// native task id -> wf task record (event) id
	Map map[string]string `json:"map"`
}

// TaskCreated enforces task shape and mirrors native tasks into wf task
// records (the no-leak funnel, 03 §7). Fail-closed on engine errors.
func TaskCreated(c *runctl.Ctl, in *hookio.Input) hookio.Result {
	r, err := c.Store.LoadRun()
	if err != nil {
		return hookio.Block("wf engine unhealthy (task gate is fail-closed): " + err.Error())
	}
	if r == nil || r.Status != "active" {
		return hookio.Block("No active wf run — tasks are created inside a run (start with /wf:dev). Native task rolled back.")
	}
	if r.Phase != "plan" && r.Phase != "build" {
		return hookio.Block(fmt.Sprintf("Tasks are created under Plan or Build (current phase: %s). Follow the phase procedure.", r.Phase))
	}
	env, err := c.Env(r)
	if err != nil {
		return hookio.Block("wf engine unhealthy: " + err.Error())
	}
	// mirror: link the native task to an existing wf task record via the
	// matching ladder (tid token → exact subject → containment), else create.
	subject := strings.TrimSpace(in.TaskSubject)
	recID := matchTaskRecord(env, subject)
	if recID == "" {
		tid := fmt.Sprintf("T-%d", len(env.Records("task"))+1)
		dod := in.TaskDescription
		if strings.TrimSpace(dod) == "" {
			dod = subject
		}
		ev, err := c.Record("task", map[string]any{
			"tid": tid, "subject": subject, "dod": []any{dod}, "status": "open",
		}, true, "hook")
		if err != nil {
			return hookio.Block("task record refused: " + err.Error())
		}
		recID = ev.ID
	}
	var m taskMirror
	_ = c.Store.LoadLocal("tasks-mirror.json", &m)
	if m.Map == nil {
		m.Map = map[string]string{}
	}
	m.Map[in.TaskID] = recID
	_ = c.Store.SaveLocal("tasks-mirror.json", &m)
	return hookio.Allow()
}

// TaskCompleted verifies the task's DoD records exist before the native task
// may close (the anti-"checked it off anyway" gate). Fail-closed.
func TaskCompleted(c *runctl.Ctl, in *hookio.Input) hookio.Result {
	r, err := c.Store.LoadRun()
	if err != nil {
		return hookio.Block("wf engine unhealthy (task gate is fail-closed): " + err.Error())
	}
	if r == nil {
		return hookio.Allow() // tasks outside runs are not ours
	}
	env, err := c.Env(r)
	if err != nil {
		return hookio.Block("wf engine unhealthy: " + err.Error())
	}
	rec := findTaskRecord(c, env, in)
	if rec == nil {
		return hookio.Allow() // not a mirrored wf task
	}
	tid, _ := rec.Data["tid"].(string)
	if r.Family == "diff" {
		ok, detail, err := contracts.EvalOne(env, "any-of", map[string]any{
			"items": []any{
				map[string]any{"predicate": "red-green", "params": map[string]any{"link": "task"}},
				map[string]any{"predicate": "linked-record", "params": map[string]any{"kind": "waiver", "link": "item"}},
			},
		}, tid)
		if err != nil {
			return hookio.Block("wf engine unhealthy: " + err.Error())
		}
		if !ok {
			return hookio.Block(fmt.Sprintf(
				"Task %s (%s) is not complete: no red→green test-run pair tagged task=%s. Run the failing test first, make it pass (auto-captured), or for a genuinely testless task: wf contract waive %s --reason \"…\". [%s]",
				tid, in.TaskSubject, tid, tid, detail))
		}
	}
	// mark the wf record done
	_, _ = c.Record("task", map[string]any{"updates": rec.ID, "status": "done"}, true, "hook")
	return hookio.Allow()
}

func findTaskRecord(c *runctl.Ctl, env *contracts.Env, in *hookio.Input) *contracts.Record {
	var m taskMirror
	_ = c.Store.LoadLocal("tasks-mirror.json", &m)
	if id, ok := m.Map[in.TaskID]; ok {
		for _, tr := range env.Records("task") {
			if tr.ID == id {
				return &tr
			}
		}
	}
	if id := matchTaskRecord(env, strings.TrimSpace(in.TaskSubject)); id != "" {
		for _, tr := range env.Records("task") {
			if tr.ID == id {
				return &tr
			}
		}
	}
	return nil
}

// tidToken matches a leading "T-<n>" (with optional ":" separator) in a
// native task subject — the linking convention the plan skill prescribes.
var tidToken = regexp.MustCompile(`^\s*(T-\d+)\b:?\s*`)

// matchTaskRecord implements the matching ladder that keeps native tasks and
// wf task records one-to-one (the T-3/T-4 duplication bug from live testing):
//  1. a leading "T-<n>:" token in the native subject → the record with that tid
//  2. case-insensitive exact subject match
//  3. normalized containment (either subject contains the other)
//
// Returns the record ID, or "" when nothing matches (caller creates one).
func matchTaskRecord(env *contracts.Env, subject string) string {
	tasks := env.Records("task")
	if m := tidToken.FindStringSubmatch(subject); m != nil {
		for _, tr := range tasks {
			if tid, _ := tr.Data["tid"].(string); strings.EqualFold(tid, m[1]) {
				return tr.ID
			}
		}
	}
	stripped := strings.TrimSpace(tidToken.ReplaceAllString(subject, ""))
	for _, tr := range tasks {
		if s, _ := tr.Data["subject"].(string); strings.EqualFold(strings.TrimSpace(s), subject) ||
			strings.EqualFold(strings.TrimSpace(s), stripped) {
			return tr.ID
		}
	}
	norm := func(s string) string { return strings.ToLower(strings.Join(strings.Fields(s), " ")) }
	ns := norm(stripped)
	if len(ns) >= 8 { // containment on very short strings is noise
		for _, tr := range tasks {
			s, _ := tr.Data["subject"].(string)
			nr := norm(s)
			if nr != "" && (strings.Contains(nr, ns) || strings.Contains(ns, nr)) {
				return tr.ID
			}
		}
	}
	return ""
}

// ---------------------------------------------------------------------------
// Gate 3 — SubagentStop verdict gate (04 §4)
// ---------------------------------------------------------------------------

var verdictBlockRe = regexp.MustCompile("(?s)```verdict\\s*\n(.*?)```")

type parsedVerdict struct {
	status            string
	criticals, majors int
	scope             string
	reason            string
	ok                bool
}

func parseVerdict(text string) parsedVerdict {
	m := verdictBlockRe.FindStringSubmatch(text)
	if m == nil {
		return parsedVerdict{}
	}
	v := parsedVerdict{criticals: -1, majors: -1}
	for _, line := range strings.Split(m[1], "\n") {
		k, val, found := strings.Cut(line, ":")
		if !found {
			continue
		}
		k = strings.TrimSpace(strings.ToLower(k))
		val = strings.TrimSpace(val)
		switch k {
		case "status":
			v.status = strings.ToLower(val)
		case "criticals":
			fmt.Sscanf(val, "%d", &v.criticals)
		case "majors":
			fmt.Sscanf(val, "%d", &v.majors)
		case "scope":
			v.scope = val
		case "reason":
			v.reason = val
		}
	}
	// n/a self-attests inapplicability — a reasonless n/a is not a verdict
	// (the agent prose always demanded the reason; now the parse does too)
	v.ok = v.status != "" && v.criticals >= 0 && v.majors >= 0 &&
		(v.status != "n/a" || v.reason != "")
	return v
}

type verdictAttempts struct {
	Attempts map[string]int `json:"attempts"` // agent_id -> blocks so far
}

const verdictBlockFormat = "The verdict is machine-parsed and required. End the final message with exactly:\n```verdict\nstatus: <clean|changes-required|safe|risky|unsafe|n/a>\ncriticals: <int>\nmajors: <int>\nscope: <assigned mode/lens, if any>\nreason: <required for n/a — one line: why this review does not apply>\n```"

// Verdict anchors verdict capture at SubagentStop: it blocks the reviewer
// until a parseable block is emitted (2 attempts), then records `unparsed`
// (fails the phase gate) and allows — no wedge (04 §4).
func Verdict(c *runctl.Ctl, in *hookio.Input) hookio.Result {
	agentName := strings.TrimPrefix(in.AgentType, "wf:")
	ag, ok := c.Spec.AgentByName(agentName)
	if !ok || !ag.Gating {
		return hookio.Allow()
	}
	text := in.LastAssistantMessage
	pv := parseVerdict(text)
	if !pv.ok && in.AgentTranscriptPath != "" {
		pv = parseVerdict(tailFile(in.AgentTranscriptPath, 16*1024))
	}
	if !pv.ok {
		var va verdictAttempts
		_ = c.Store.LoadLocal("verdict-attempts.json", &va)
		if va.Attempts == nil {
			va.Attempts = map[string]int{}
		}
		va.Attempts[in.AgentID]++
		n := va.Attempts[in.AgentID]
		_ = c.Store.SaveLocal("verdict-attempts.json", &va)
		if n <= 2 {
			return hookio.StopBlock(verdictBlockFormat)
		}
		_, err := c.Record("verdict", map[string]any{
			"agent": agentName, "status": "unparsed", "criticals": 0, "majors": 0,
			"scope": defaultScope(c, agentName),
		}, true, "hook")
		if err != nil {
			return hookio.Block("verdict recording failed (fail-closed): " + err.Error())
		}
		return hookio.StopAllowMessage("[wf] " + agentName + " produced no parseable verdict after 2 retries — recorded `unparsed` (fails the phase gate)")
	}
	scope := pv.scope
	if scope == "" {
		scope = defaultScope(c, agentName)
	}
	data := map[string]any{
		"agent": agentName, "status": pv.status,
		"criticals": pv.criticals, "majors": pv.majors,
	}
	if scope != "" {
		data["scope"] = scope
	}
	if pv.reason != "" {
		data["reason"] = pv.reason
	}
	ev, err := c.Record("verdict", data, true, "hook")
	if err != nil {
		return hookio.Block("verdict recording failed (fail-closed): " + err.Error())
	}
	status, _ := ev.Data["status"].(string)
	msg := fmt.Sprintf("[wf] %s verdict recorded: %s (criticals=%d majors=%d)", agentName, status, pv.criticals, pv.majors)
	if d, _ := ev.Data["downgraded"].(bool); d {
		msg += " — auto-downgraded: clean/safe cannot carry findings"
	}
	return hookio.StopAllowMessage(msg)
}

// defaultScope infers the reviewer's mode from the phase when the verdict
// block omits it (adversary's modes are phase-bound).
func defaultScope(c *runctl.Ctl, agent string) string {
	r, err := c.Store.LoadRun()
	if err != nil || r == nil {
		return ""
	}
	if agent == "adversary" {
		switch r.Phase {
		case "frame":
			return "abuse-case"
		case "design":
			return "attack-tree"
		default:
			return "red-team"
		}
	}
	if agent == "lens-reviewer" && r.Phase == "frame" {
		return "security" // the gated lens; others recorded via explicit scope
	}
	if agent == "compliance-reviewer" {
		// exactly one standard in force ⇒ unambiguous; several ⇒ the
		// reviewer must self-declare (the briefing says so)
		if stds := c.Spec.ComplianceStandards(); len(stds) == 1 {
			return stds[0]
		}
	}
	return ""
}

func tailFile(path string, n int64) string {
	f, err := os.Open(path)
	if err != nil {
		return ""
	}
	defer f.Close()
	st, err := f.Stat()
	if err != nil {
		return ""
	}
	off := int64(0)
	if st.Size() > n {
		off = st.Size() - n
	}
	buf := make([]byte, st.Size()-off)
	if _, err := f.ReadAt(buf, off); err != nil {
		return ""
	}
	// transcripts are JSONL with escaped newlines — unescape enough to match
	s := strings.ReplaceAll(string(buf), `\n`, "\n")
	s = strings.ReplaceAll(s, "\\u0060", "`")
	return s
}

// ---------------------------------------------------------------------------
// Gate 4 — PreToolUse tool gates (04 §5)
// ---------------------------------------------------------------------------

// Skill denies invoking a phase skill that is not the active phase (or a
// legal loop target). Sequencing gate: fails open + loud.
func Skill(c *runctl.Ctl, in *hookio.Input) hookio.Result {
	if hookio.EnforceDisabled(in) {
		recordEnforceOff(c, in, "skill")
		return hookio.Allow()
	}
	name := skillName(in)
	name = strings.TrimPrefix(name, "wf:")
	// which phase (if any) does this skill drive?
	var target string
	for _, p := range c.Spec.Phases {
		if p.Skill == name {
			target = p.ID
			break
		}
	}
	if target == "" {
		return hookio.Allow() // not a phase skill (dev/init/status/park/force/other plugins)
	}
	r, err := c.Store.LoadRun()
	if err != nil {
		return hookio.BrokenGate(err)
	}
	if r == nil || r.Status != "active" {
		return hookio.Deny("No active run — phase procedures run inside a run. Start with /wf:dev.")
	}
	if target == r.Phase {
		// entry re-check: adopt/resume/force paths can land in a phase whose
		// inputs were never produced — the transition gate never saw them
		if env, err := c.Env(r); err == nil {
			if entry, err := contracts.EvaluateEntry(env, target); err == nil && len(entry) > 0 {
				return hookio.Deny(fmt.Sprintf("Phase %s is missing its inputs: %s → %s (deliberate skip: wf contract waive %s --reason …)",
					target, entry[0].ID, entry[0].Remediation, entry[0].ID))
			}
		}
		return hookio.Allow()
	}
	for _, t := range c.Spec.Loops.Targets {
		if target == t && r.Phase == c.Spec.Loops.From {
			return hookio.Allow() // loop-back procedure is legal from verify
		}
	}
	return hookio.Deny(fmt.Sprintf("The active phase is %s (its procedure: /wf:%s). Phase order is enforced — /wf:%s is not available now.",
		r.Phase, skillOf(c, r.Phase), name))
}

func skillName(in *hookio.Input) string {
	for _, k := range []string{"skill_name", "name", "skill", "command"} {
		if v := in.ToolInputField(k); v != "" {
			return v
		}
	}
	return ""
}

func skillOf(c *runctl.Ctl, phase string) string {
	if p, ok := c.Spec.PhaseByID(phase); ok {
		return p.Skill
	}
	return "dev"
}

// exemptEditPrefixes are path anchors (never basenames — the C7 fix).
var exemptEditPrefixes = []string{".workflow/", "docs/", ".claude/", "CLAUDE.md", "AGENTS.md"}

// protectedStatePaths are the engine-written ledger, snapshot, archives and
// project config: the evidence chain. Agents never write these through
// tools — the engine writes them from its own process (wf record/…), so a
// tool-level deny costs nothing legitimate and closes the forgery path
// (an agent minting `auto:true grounded` events directly). Unconditional:
// active run or not, WF_ENFORCE or not.
var protectedStatePaths = []string{
	".workflow/log/", ".workflow/state/", ".workflow/runs/", ".workflow/config.json",
}

func protectedStatePath(rel string) bool {
	for _, p := range protectedStatePaths {
		if strings.HasPrefix(rel, p) || rel == strings.TrimSuffix(p, "/") {
			return true
		}
	}
	return false
}

// Edit is the stray-edit guard: project files change only inside an active
// run. Sequencing gate: fails open + loud — EXCEPT the protected-state deny,
// which is data protection and has no escape hatch.
func Edit(c *runctl.Ctl, in *hookio.Input) hookio.Result {
	path := in.ToolInputField("file_path")
	rel := relToProject(path, in.CWD)
	if protectedStatePath(rel) {
		return hookio.Deny("wf state under " + rel + " is engine-written (no override): use wf record/approve/… — config.json changes belong to the user")
	}
	if hookio.EnforceDisabled(in) {
		recordEnforceOff(c, in, "edit")
		return hookio.Allow()
	}
	for _, p := range exemptEditPrefixes {
		if strings.HasPrefix(rel, p) {
			return hookio.Allow()
		}
	}
	r, err := c.Store.LoadRun()
	if err != nil {
		return hookio.BrokenGate(err)
	}
	if r == nil || r.Status != "active" {
		return hookio.Deny("No active run — project files change inside a run (start or resume with /wf:dev). Bookkeeping under .workflow/ and docs/ is exempt.")
	}
	return hookio.Allow()
}

func relToProject(path, cwd string) string {
	root := os.Getenv("CLAUDE_PROJECT_DIR")
	if root == "" {
		root = cwd
	}
	p := strings.ReplaceAll(path, "\\", "/")
	root = strings.ReplaceAll(root, "\\", "/")
	p = strings.TrimPrefix(p, strings.TrimSuffix(root, "/")+"/")
	return p
}

// catastrophic Bash patterns — the always-on net with NO escape hatch.
var catastrophic = []*regexp.Regexp{
	regexp.MustCompile(`\brm\s+(-[a-zA-Z]*r[a-zA-Z]*f|-[a-zA-Z]*f[a-zA-Z]*r)[a-zA-Z]*\s+(/|~|\$HOME)(\s|$)`),
	regexp.MustCompile(`\bgit\s+push\s+[^|;&]*--force(-with-lease)?\s+[^|;&]*\b(main|master)\b`),
	regexp.MustCompile(`\bgit\s+push\s+[^|;&]*\b(origin\s+)?\+?(main|master)\b[^|;&]*--force`),
	regexp.MustCompile(`\b(curl|wget)\b[^|;&]*\|\s*(sudo\s+)?(ba|z|da)?sh\b`),
	regexp.MustCompile(`\bmkfs\b|\bdd\s+[^|;&]*of=/dev/`),
	regexp.MustCompile(`>\s*/etc/`),
	regexp.MustCompile(`\bchmod\s+(-[a-zA-Z]*R[a-zA-Z]*\s+)?777\s+/(\s|$)`),
}

// protectedTarget matches the protected .workflow state paths inside a
// command string (relative or absolute).
var protectedTarget = regexp.MustCompile(`(^|[/\s"'=])\.workflow/(log|state|runs)(/|\s|"|'|$)|(^|[/\s"'=])\.workflow/config\.json`)

// stateWriters are Bash constructs able to write/destroy files. A segment
// that both starts with (or contains, for redirects) one of these and
// mentions a protected path is denied. Reading (cat/grep/less) stays free —
// the ledger is meant to be read.
var (
	redirectIntoState = regexp.MustCompile(`>>?\s*['"]?[^|;&<>\s'"]*\.workflow/(log/|state/|runs/|config\.json)`)
	writerHeads       = []string{"tee", "sed", "cp", "mv", "rm", "install", "rsync", "truncate", "dd", "ln", "touch", "chmod", "chown"}
)

// wfStateTamper detects Bash writes into the protected .workflow state —
// the forgery path the Edit gate cannot see. Heuristic on purpose: the
// hash-chained ledger (doctor) is the backstop for what slips past.
func wfStateTamper(cmd string) bool {
	if redirectIntoState.MatchString(cmd) {
		return true
	}
	for _, seg := range segmentSplit.Split(cmd, -1) {
		for _, p := range strings.Split(seg, "|") {
			head := commandHead(strings.TrimSpace(p))
			for _, w := range writerHeads {
				if (head == w || strings.HasPrefix(head, w+" ")) && protectedTarget.MatchString(p) {
					return true
				}
			}
		}
	}
	return false
}

// localStateTarget matches .workflow/local — the engine-private per-machine
// state (counters, mirrors, pending approval challenge codes). Nothing
// legitimate reads or writes it through tools; challenge codes make model
// reads actively harmful (04 §8.1: the code must reach the model only
// through the user's typed answer).
var localStateTarget = regexp.MustCompile(`(^|[/\s"'=])\.workflow/local(/|\s|"|'|$)`)

// statuslineInvocation matches agent-run `wf statusline` — the statusline
// payload is the user's channel and may carry a pending challenge code.
var statuslineInvocation = regexp.MustCompile(`\bwf(\.exe|-[a-z0-9-]+)?["']?\s+statusline\b`)

// Bash is the catastrophic-command net (deny; duplicated as permission rules
// where expressible). Store-free and always-on: it also protects the wf
// ledger from direct writes even before adoption state can be read.
func Bash(_ *runctl.Ctl, in *hookio.Input) hookio.Result {
	cmd := in.ToolInputField("command")
	for _, re := range catastrophic {
		if re.MatchString(cmd) {
			return hookio.Deny("Catastrophic command blocked by wf (no override): " + re.String())
		}
	}
	if wfStateTamper(cmd) {
		return hookio.Deny("Direct writes into .workflow/{log,state,runs,config.json} are blocked (no override): the ledger is engine-written — use wf record/approve/…; config.json changes belong to the user")
	}
	if localStateTarget.MatchString(cmd) {
		return hookio.Deny(".workflow/local is engine-private per-machine state (no override): it holds nothing an agent needs — and may hold a pending approval challenge code")
	}
	if statuslineInvocation.MatchString(cmd) {
		return hookio.Deny("wf statusline is the user's display channel (it may carry an approval challenge code) — agents read run state via wf status")
	}
	return hookio.Allow()
}

// ReadTool guards the Read/Grep/Glob tools against .workflow/local — the
// only project path whose CONTENT is model-sensitive (pending challenge
// codes). Store-free data protection, no escape hatch; everything else in
// .workflow stays readable (the ledger is meant to be read).
func ReadTool(in *hookio.Input) hookio.Result {
	for _, field := range []string{"file_path", "path", "pattern", "include"} {
		v := in.ToolInputField(field)
		if v == "" {
			continue
		}
		if localStateTarget.MatchString(v) || strings.Contains(strings.ReplaceAll(v, "\\", "/"), ".workflow/local") {
			return hookio.Deny(".workflow/local is engine-private per-machine state (no override): it holds nothing an agent needs — and may hold a pending approval challenge code")
		}
	}
	return hookio.Allow()
}

// ---------------------------------------------------------------------------
// PostToolUse capture (04 §5)
// ---------------------------------------------------------------------------

// test runners recognized in the command HEAD only (the G1 fix). This list
// is the zero-config fast path — NOT the whole story: runners are also
// learned from the run's own verification-strategy commands (any language)
// and from config `runners` (custom wrappers). The power-of-ten incident:
// `python3 -m unittest` was missing here and every test-run went uncaptured.
var runnerHeads = []string{
	"go test", "gotestsum", "pytest", "python -m pytest", "python3 -m pytest",
	"python -m unittest", "python3 -m unittest", "py -m unittest",
	"npm test", "npm run test", "pnpm test", "pnpm run test", "yarn test",
	"yarn run test", "npx jest", "jest", "npx vitest", "vitest", "mocha",
	"npx mocha", "npx playwright test", "playwright test", "cypress run",
	"npx cypress run", "ng test", "deno test", "bun test",
	"cargo test", "dotnet test", "mvn test", "gradle test", "./gradlew test",
	"make test", "ctest", "rspec", "bundle exec rspec", "phpunit",
	"vendor/bin/phpunit", "composer test", "mix test", "sbt test",
	"stack test", "cabal test", "swift test", "dart test", "flutter test",
	"tox", "nox", "zig build test", "gitleaks",
}

// matchHead: rh matches at a token boundary — or a `:` (npm-style script
// variants: `npm run test:unit`). Bare prefixes are NOT enough ("toxiproxy"
// must not match "tox").
func matchHead(head, rh string) bool {
	return head == rh || strings.HasPrefix(head, rh+" ") || strings.HasPrefix(head, rh+":")
}

var filterPipe = regexp.MustCompile(`\|\s*(grep|head|tail|awk|sed|rg)\b`)

// CaptureTest turns recognized runner invocations into grounded test-run
// records. Rules (G1): match the head only, skip wf self-calls, treat
// filter-pipes and missing exit codes as ungrounded. Recognition, in order:
// static runnerHeads, config `runners`, then runners learned from the run's
// recorded verification-strategy commands — the last makes capture
// language-agnostic, since Plan already declared this run's test commands.
func CaptureTest(c *runctl.Ctl, in *hookio.Input) hookio.Result {
	cmd := strings.TrimSpace(in.ToolInputField("command"))
	if cmd == "" || strings.HasPrefix(cmd, "wf ") || strings.Contains(cmd, "/wf ") {
		return hookio.Allow()
	}
	r, err := c.Store.LoadRun()
	if err != nil || r == nil || r.Status != "active" {
		return hookio.Allow()
	}
	env, envErr := c.Env(r)

	// `cd <dir> && <runner>` is the single chained form common enough to
	// honor (it drove users toward the config-runners escape): strip ONE
	// leading cd for recognition and grounding. Honest bound: a failing cd
	// reports its own exit — a false red at worst, never a false green.
	effective := cmdid.Effective(cmd)
	head := commandHead(effective)
	var category, ac string
	matched := false
	for _, rh := range runnerHeads {
		if matchHead(head, rh) {
			matched = true
			if strings.HasPrefix(rh, "gitleaks") {
				category = "secret-scan"
			}
			break
		}
	}
	if !matched && c.Config != nil {
		for _, rh := range c.Config.Runners {
			if rh != "" && matchHead(head, rh) {
				matched = true
				break
			}
		}
	}
	if !matched && envErr == nil {
		matched, ac = strategyMatch(env, head)
	}
	if !matched {
		return hookio.Allow()
	}
	if interrupted(in) {
		return hookio.Allow() // a ctrl-C'd run is not evidence, red or green
	}
	exit, hasExit := commandExit(in)
	grounded := hasExit && !filterPipe.MatchString(effective) && !chained(effective)
	data := map[string]any{"cmd": cmd, "grounded": grounded}
	if hasExit {
		data["exit"] = exit
	} else {
		data["exit"] = nil
	}
	if category != "" {
		data["category"] = category
	}
	if ac != "" {
		data["ac"] = ac // the exact per-AC verification command was run
	}
	// bind to the single in-progress task (and its ACs) when unambiguous
	if envErr == nil {
		if tid, acs := activeTask(env); tid != "" {
			data["task"] = tid
			if _, tagged := data["ac"]; !tagged && len(acs) == 1 {
				data["ac"] = acs[0]
			}
		}
	}
	if _, err := c.Record("test-run", data, true, "hook"); err != nil {
		return hookio.Allow() // capture must never break the loop
	}
	// grounded green runner output is the one honest source for coverage —
	// scrape it (ClaudeInit's proven extraction, engine-native) so the
	// verify.quality-floor contract runs on measured numbers, not claims
	if grounded && hasExit && exit == 0 && envErr == nil {
		if cov, ok := coverageFromResponse(in.ToolResponse); ok {
			recordCoverage(c, env, cov)
		}
	}
	return hookio.Allow()
}

// coveragePatterns match the coverage summary lines of the common runners.
// First match wins; the value is the percentage.
var coveragePatterns = []*regexp.Regexp{
	regexp.MustCompile(`coverage: (\d+(?:\.\d+)?)% of statements`),   // go test -cover
	regexp.MustCompile(`(?m)^TOTAL\s+.*?(\d+(?:\.\d+)?)%`),           // coverage.py / pytest-cov
	regexp.MustCompile(`(?m)^All files\s*\|\s*(\d+(?:\.\d+)?)\s*\|`), // jest/istanbul text table
	regexp.MustCompile(`(\d+(?:\.\d+)?)% coverage`),                  // cargo-tarpaulin
	regexp.MustCompile(`lines\.+:\s*(\d+(?:\.\d+)?)%`),               // lcov
}

func coverageFromResponse(raw json.RawMessage) (float64, bool) {
	if len(raw) == 0 {
		return 0, false
	}
	var m struct {
		Stdout string `json:"stdout"`
		Stderr string `json:"stderr"`
	}
	if json.Unmarshal(raw, &m) != nil {
		return 0, false
	}
	out := m.Stdout + "\n" + m.Stderr
	for _, re := range coveragePatterns {
		if hit := re.FindStringSubmatch(out); hit != nil {
			if v, err := strconv.ParseFloat(hit[1], 64); err == nil && v >= 0 && v <= 100 {
				return v, true
			}
		}
	}
	return 0, false
}

// recordCoverage records the grounded coverage metric, UPDATING the run's
// existing coverage record so a re-measure supersedes (no stale
// below_threshold=true blocking after the fix). below_threshold is computed
// by runctl's write-time validation against config.thresholds.
func recordCoverage(c *runctl.Ctl, env *contracts.Env, cov float64) {
	data := map[string]any{"name": "coverage", "value": cov, "grounded": true}
	for _, m := range env.Records("metric") {
		if n, _ := m.Data["name"].(string); n == "coverage" {
			data["updates"] = m.ID
			break
		}
	}
	_, _ = c.Record("metric", data, true, "hook")
}

// strategyMatch recognizes test invocations from the run's own
// verification-strategy records. Two tiers:
//  1. the command IS a recorded per-AC verification command (whole-token
//     prefix either way, ≥2 shared tokens — flag variations tolerated)
//     → matched, tagged with that AC;
//  2. the command shares a strategy's learned runner head (same runner,
//     different selector — where Build's red/green runs live) → matched.
func strategyMatch(env *contracts.Env, head string) (bool, string) {
	cmdTok := strings.Fields(head)
	strategies := env.Records("verification-strategy")
	for _, s := range strategies {
		sc, _ := s.Data["command"].(string)
		if sc == "" {
			continue
		}
		if tokenPrefix(cmdTok, strings.Fields(commandHead(sc))) {
			acv, _ := s.Data["ac"].(string)
			return true, acv
		}
	}
	for _, s := range strategies {
		sc, _ := s.Data["command"].(string)
		lh := strings.Fields(learnedHead(sc))
		if len(lh) == 0 || len(cmdTok) < len(lh) {
			continue
		}
		match := true
		for i := range lh {
			if cmdTok[i] != lh[i] {
				match = false
				break
			}
		}
		if match {
			return true, ""
		}
	}
	return false, ""
}

// tokenPrefix, learnedHead and commandHead moved to the shared cmdid
// package (contracts needs the same command-identity reasoning for
// red→green pairing and must not import gates). Thin aliases keep the
// call sites readable.
func tokenPrefix(a, b []string) bool { return cmdid.TokenPrefix(a, b) }
func learnedHead(cmd string) string  { return cmdid.LearnedHead(cmd) }
func commandHead(cmd string) string  { return cmdid.Head(cmd) }

// commandExit resolves the exit code from the DOCUMENTED hook event
// semantics (the four-TestRepo-runs discovery): the Bash tool_response
// carries no exit-code field of any name, and a non-zero exit never fires
// PostToolUse at all —
//   - PostToolUse means the command "completed successfully" ⇒ exit 0;
//   - PostToolUseFailure carries the code inside the error string
//     ("Command exited with non-zero status code 1").
//
// An explicit response field still wins if a future release adds one.
func commandExit(in *hookio.Input) (int, bool) {
	if exit, ok := responseExit(in.ToolResponse); ok {
		return exit, true
	}
	switch in.HookEventName {
	case "PostToolUse":
		return 0, true
	case "PostToolUseFailure":
		if m := failureExitRe.FindStringSubmatch(in.Error); m != nil {
			if n, err := strconv.Atoi(m[1]); err == nil {
				return n, true
			}
		}
		// a failure without a parseable code (timeout, spawn error…) ran
		// but proves nothing — recorded ungrounded
		return 0, false
	}
	return 0, false
}

var failureExitRe = regexp.MustCompile(`non-zero (?:status|exit) code (\d+)`)

// interrupted: a user-interrupted command is not evidence in either
// direction. PostToolUse carries `interrupted` in the response; the failure
// event carries top-level `is_interrupt`.
func interrupted(in *hookio.Input) bool {
	if in.IsInterrupt {
		return true
	}
	var m struct {
		Interrupted bool `json:"interrupted"`
	}
	if len(in.ToolResponse) > 0 && json.Unmarshal(in.ToolResponse, &m) == nil {
		return m.Interrupted
	}
	return false
}

// chained: `&&`/`||`/`;`/newline chains report the LAST command's exit, not
// the runner's — evidence quality guard (run 4's "compound piped commands"
// confusion, codified). Heuristic on purpose; false positives only cost a
// grounded flag, never a record.
func chained(cmd string) bool {
	return strings.Contains(cmd, "&&") || strings.Contains(cmd, "||") ||
		strings.Contains(cmd, ";") || strings.Contains(cmd, "\n")
}

// responseExit checks for an explicit exit-code field in the tool response
// (none is documented today; kept for forward compatibility).
func responseExit(raw json.RawMessage) (int, bool) {
	if len(raw) == 0 {
		return 0, false
	}
	var m map[string]any
	if json.Unmarshal(raw, &m) != nil {
		return 0, false
	}
	for _, k := range []string{"exit_code", "exitCode", "code", "returnCodeInterpretation"} {
		if v, ok := m[k]; ok {
			if f, ok := v.(float64); ok {
				return int(f), true
			}
		}
	}
	return 0, false
}

func activeTask(env *contracts.Env) (string, []string) {
	// prefer the single in_progress task; fall back to a single open one
	// (the common one-task flow where the agent skipped the status update)
	for _, statuses := range [][]string{{"in_progress"}, {"in_progress", "open"}} {
		var tid string
		var acs []string
		count := 0
		for _, tr := range env.Records("task") {
			s, _ := tr.Data["status"].(string)
			match := false
			for _, want := range statuses {
				if s == want {
					match = true
				}
			}
			if !match {
				continue
			}
			count++
			tid, _ = tr.Data["tid"].(string)
			acs = nil
			if raw, ok := tr.Data["ac_links"].([]any); ok {
				for _, a := range raw {
					if s, ok := a.(string); ok {
						acs = append(acs, s)
					}
				}
			}
		}
		if count == 1 {
			return tid, acs
		}
	}
	return "", nil
}

// CaptureCommit records durable commit→run attribution (08 §6) when the
// Bash call contains a `git commit` subcommand that succeeded. Recording
// only — never blocks, never guesses on failure.
func CaptureCommit(c *runctl.Ctl, in *hookio.Input) hookio.Result {
	cmd := in.ToolInputField("command")
	if !hasGitCommitSegment(cmd) {
		return hookio.Allow()
	}
	r, err := c.Store.LoadRun()
	if err != nil || r == nil || r.Status != "active" {
		return hookio.Allow()
	}
	dir := in.CWD
	if dir == "" {
		dir = "."
	}
	sha, err := gitOut(dir, "rev-parse", "HEAD")
	if err != nil || sha == "" {
		return hookio.Allow()
	}
	// idempotence: skip if this sha is already recorded
	if env, err := c.Env(r); err == nil {
		for _, co := range env.Records("commit-origin") {
			if s, _ := co.Data["commit"].(string); s == sha {
				return hookio.Allow()
			}
		}
	}
	msg, _ := gitOut(dir, "log", "-1", "--format=%s")
	tagged := strings.Contains(msg, "[run:"+r.ID+"]")
	_, _ = c.Record("commit-origin", map[string]any{
		"commit": sha, "run": r.ID, "tagged": tagged, "subject": msg,
	}, true, "hook")
	if !tagged {
		return hookio.AllowJSON(map[string]any{
			"systemMessage": "[wf] commit " + sha[:min(8, len(sha))] + " recorded, but its message lacks the [run:" + r.ID + "] tag — include it in future commits",
		})
	}
	return hookio.Allow()
}

var segmentSplit = regexp.MustCompile(`&&|\|\||;`)

func hasGitCommitSegment(cmd string) bool {
	for _, seg := range segmentSplit.Split(cmd, -1) {
		head := commandHead(strings.TrimSpace(seg))
		if strings.HasPrefix(head, "git commit") {
			return true
		}
	}
	return false
}

func gitOut(dir string, args ...string) (string, error) {
	out, err := exec.Command("git", append([]string{"-C", dir}, args...)...).Output()
	return strings.TrimSpace(string(out)), err
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// CaptureEdit appends the edit→task binding ledger (never blocks).
// CaptureQuestion turns an AskUserQuestion exchange into a hook-captured
// user-answer record — the anchoring evidence `wf approve` links via
// answer_ref (04 §8.1: still not proof a human typed it, noted as such;
// it is one layer harder to fabricate than a bare wf approve). Defensive
// by design: the tool's payload shape is not documented upstream, so this
// extracts what it can and records only when BOTH sides yielded text.
// Never blocks.
func CaptureQuestion(c *runctl.Ctl, in *hookio.Input) hookio.Result {
	r, err := c.Store.LoadRun()
	if err != nil || r == nil || r.Status != "active" {
		return hookio.Allow()
	}
	question := questionText(in.ToolInput)
	answer := answerText(in.ToolResponse)
	if question == "" || answer == "" {
		return hookio.Allow()
	}
	data := map[string]any{"question": clip(question, 300), "answer": clip(answer, 300)}
	if topic := questionTopic(question); topic != "" {
		data["topic"] = topic
	}
	_, _ = c.Record("user-answer", data, true, "hook")
	return hookio.Allow()
}

// topicKeywords maps approval gates to the words an approval question about
// that gate carries (the skills phrase them so). Used by wf approve to
// anchor to the RIGHT answer: before topics, "the newest answer after the
// last approval" let any unrelated question anchor a hardened approval.
var topicKeywords = map[string][]string{
	"frame":     {"classif", "family", "intent", "framing", "frame"},
	"scope":     {"scope", "requirement", "assumption"},
	"design":    {"design", "option", "architecture"},
	"plan":      {"plan", "task breakdown", "verification strategy"},
	"deferral":  {"defer"},
	"deviation": {"deviation", "depart"},
	"lesson":    {"lesson"},
}

// questionTopic infers the single approval gate a question is about;
// ambiguous questions (words of several gates) get no topic — a wrong
// anchor is worse than none.
func questionTopic(q string) string {
	low := strings.ToLower(q)
	found := ""
	for gate, words := range topicKeywords {
		for _, w := range words {
			if strings.Contains(low, w) {
				if found != "" && found != gate {
					return ""
				}
				found = gate
				break
			}
		}
	}
	return found
}

// questionText joins the question strings from AskUserQuestion tool_input
// ({"questions":[{"question":…},…]} — with a generic-walk fallback).
func questionText(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	var in struct {
		Questions []struct {
			Question string `json:"question"`
		} `json:"questions"`
	}
	if json.Unmarshal(raw, &in) == nil && len(in.Questions) > 0 {
		var qs []string
		for _, q := range in.Questions {
			if q.Question != "" {
				qs = append(qs, q.Question)
			}
		}
		if len(qs) > 0 {
			return strings.Join(qs, " | ")
		}
	}
	return strings.Join(collectStrings(raw, map[string]bool{"question": true}), " | ")
}

// answerText extracts the CHOSEN answer(s) from tool_response, tiered so
// the echoed option labels can't flood the evidence (the circle-area run:
// "No, needs changes | Yes, that's correct" — every label, not the pick):
//  1. the documented `answers` object ({question: chosen label}) or
//     array-of-{question,answer} variants;
//  2. explicit answer-ish keys;
//  3. bare `label` values — last resort, may include unchosen options.
func answerText(raw json.RawMessage) string {
	if len(raw) == 0 {
		return ""
	}
	if got := answersValues(raw); len(got) > 0 {
		return strings.Join(got, " | ")
	}
	if got := collectStrings(raw, map[string]bool{"answer": true, "selected": true, "choice": true}); len(got) > 0 {
		return strings.Join(got, " | ")
	}
	if got := collectStrings(raw, map[string]bool{"label": true}); len(got) > 0 {
		return strings.Join(got, " | ")
	}
	return ""
}

// answersValues finds `answers` nodes anywhere in the payload and collects
// the chosen values.
func answersValues(raw json.RawMessage) []string {
	var v any
	if json.Unmarshal(raw, &v) != nil {
		return nil
	}
	var out []string
	var walk func(node any)
	walk = func(node any) {
		switch n := node.(type) {
		case map[string]any:
			for k, child := range n {
				if k == "answers" {
					answerNode(child, &out)
					continue
				}
				walk(child)
			}
		case []any:
			for _, c := range n {
				walk(c)
			}
		}
	}
	walk(v)
	sort.Strings(out)
	return out
}

// answerNode flattens an `answers` value: {question: label} maps (take the
// values), arrays of {question, answer} objects, or plain strings.
func answerNode(node any, out *[]string) {
	switch n := node.(type) {
	case map[string]any:
		keys := make([]string, 0, len(n))
		for k := range n {
			keys = append(keys, k)
		}
		sort.Strings(keys)
		for _, k := range keys {
			if s, ok := n[k].(string); ok {
				if s != "" {
					*out = append(*out, s)
				}
				continue
			}
			answerNode(n[k], out)
		}
	case []any:
		for _, c := range n {
			if m, ok := c.(map[string]any); ok {
				if a, ok := m["answer"].(string); ok && a != "" {
					*out = append(*out, a)
					continue
				}
			}
			answerNode(c, out)
		}
	case string:
		if n != "" {
			*out = append(*out, n)
		}
	}
}

// collectStrings walks arbitrary JSON and gathers string values found under
// the wanted keys (arrays of strings included).
func collectStrings(raw json.RawMessage, want map[string]bool) []string {
	var v any
	if json.Unmarshal(raw, &v) != nil {
		return nil
	}
	var out []string
	var walk func(node any, wanted bool)
	walk = func(node any, wanted bool) {
		switch n := node.(type) {
		case map[string]any:
			for k, child := range n {
				walk(child, want[k])
			}
		case []any:
			for _, child := range n {
				walk(child, wanted)
			}
		case string:
			if wanted && n != "" {
				out = append(out, n)
			}
		}
	}
	walk(v, false)
	sort.Strings(out)
	return out
}

func clip(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "…"
}

func CaptureEdit(c *runctl.Ctl, in *hookio.Input) hookio.Result {
	r, err := c.Store.LoadRun()
	if err != nil || r == nil || r.Status != "active" {
		return hookio.Allow()
	}
	path := relToProject(in.ToolInputField("file_path"), in.CWD)
	if path == "" || strings.HasPrefix(path, ".workflow/") {
		return hookio.Allow()
	}
	data := map[string]any{"path": path}
	if env, err := c.Env(r); err == nil {
		if tid, _ := activeTask(env); tid != "" {
			data["task"] = tid
		}
	}
	_, _ = c.Record("edit", data, true, "hook")
	return hookio.Allow()
}
