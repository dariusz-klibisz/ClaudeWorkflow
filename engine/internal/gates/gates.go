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
	"regexp"
	"sort"
	"strings"

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
		return hookio.StopAllowMessage("[wf] the same items blocked " + fmt.Sprint(stopSelfCap) + "× — allowing the stop. If stuck, /wf:park records the honest state")
	}
	return hookio.StopBlock(reason)
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
	// mirror: reuse an existing wf task record with the same subject, else create
	subject := strings.TrimSpace(in.TaskSubject)
	var recID string
	for _, tr := range env.Records("task") {
		if s, _ := tr.Data["subject"].(string); strings.EqualFold(s, subject) {
			recID = tr.ID
			break
		}
	}
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
	for _, tr := range env.Records("task") {
		if s, _ := tr.Data["subject"].(string); strings.EqualFold(s, strings.TrimSpace(in.TaskSubject)) {
			return &tr
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Gate 3 — SubagentStop verdict gate (04 §4)
// ---------------------------------------------------------------------------

var verdictBlockRe = regexp.MustCompile("(?s)```verdict\\s*\n(.*?)```")

type parsedVerdict struct {
	status            string
	criticals, majors int
	scope             string
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
		}
	}
	v.ok = v.status != "" && v.criticals >= 0 && v.majors >= 0
	return v
}

type verdictAttempts struct {
	Attempts map[string]int `json:"attempts"` // agent_id -> blocks so far
}

const verdictBlockFormat = "The verdict is machine-parsed and required. End the final message with exactly:\n```verdict\nstatus: <clean|changes-required|safe|risky|unsafe|n/a>\ncriticals: <int>\nmajors: <int>\nscope: <assigned mode/lens, if any>\n```"

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

// Edit is the stray-edit guard: project files change only inside an active
// run. Sequencing gate: fails open + loud.
func Edit(c *runctl.Ctl, in *hookio.Input) hookio.Result {
	if hookio.EnforceDisabled(in) {
		return hookio.Allow()
	}
	path := in.ToolInputField("file_path")
	rel := relToProject(path, in.CWD)
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

// Bash is the catastrophic-command net (deny; duplicated as permission rules
// where expressible).
func Bash(_ *runctl.Ctl, in *hookio.Input) hookio.Result {
	cmd := in.ToolInputField("command")
	for _, re := range catastrophic {
		if re.MatchString(cmd) {
			return hookio.Deny("Catastrophic command blocked by wf (no override): " + re.String())
		}
	}
	return hookio.Allow()
}

// ---------------------------------------------------------------------------
// PostToolUse capture (04 §5)
// ---------------------------------------------------------------------------

// test runners recognized in the command HEAD only (the G1 fix).
var runnerHeads = []string{
	"go test", "gotestsum", "pytest", "python -m pytest", "python3 -m pytest",
	"npm test", "npm run test", "pnpm test", "yarn test", "npx jest", "jest",
	"npx vitest", "vitest", "cargo test", "dotnet test", "mvn test", "gradle test",
	"./gradlew test", "make test", "ctest", "rspec", "phpunit", "gitleaks",
}

var filterPipe = regexp.MustCompile(`\|\s*(grep|head|tail|awk|sed|rg)\b`)

// CaptureTest turns recognized runner invocations into grounded test-run
// records. Rules (G1): match the head only, skip wf self-calls, treat
// filter-pipes and missing exit codes as ungrounded.
func CaptureTest(c *runctl.Ctl, in *hookio.Input) hookio.Result {
	cmd := strings.TrimSpace(in.ToolInputField("command"))
	if cmd == "" || strings.HasPrefix(cmd, "wf ") || strings.Contains(cmd, "/wf ") {
		return hookio.Allow()
	}
	head := commandHead(cmd)
	var category string
	matched := false
	for _, rh := range runnerHeads {
		if strings.HasPrefix(head, rh) {
			matched = true
			if strings.HasPrefix(rh, "gitleaks") {
				category = "secret-scan"
			}
			break
		}
	}
	if !matched {
		return hookio.Allow()
	}
	r, err := c.Store.LoadRun()
	if err != nil || r == nil || r.Status != "active" {
		return hookio.Allow()
	}
	exit, hasExit := responseExit(in.ToolResponse)
	grounded := hasExit && !filterPipe.MatchString(cmd)
	data := map[string]any{"cmd": cmd, "grounded": grounded}
	if hasExit {
		data["exit"] = exit
	} else {
		data["exit"] = nil
	}
	if category != "" {
		data["category"] = category
	}
	// bind to the single in-progress task (and its ACs) when unambiguous
	if env, err := c.Env(r); err == nil {
		if tid, acs := activeTask(env); tid != "" {
			data["task"] = tid
			if len(acs) == 1 {
				data["ac"] = acs[0]
			}
		}
	}
	if _, err := c.Record("test-run", data, true, "hook"); err != nil {
		return hookio.Allow() // capture must never break the loop
	}
	return hookio.Allow()
}

// commandHead strips leading env assignments and returns the command's start.
func commandHead(cmd string) string {
	fields := strings.Fields(cmd)
	i := 0
	for i < len(fields) && strings.Contains(fields[i], "=") && !strings.HasPrefix(fields[i], "=") {
		i++
	}
	return strings.Join(fields[i:], " ")
}

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
	// Bash tool responses without an explicit code: interrupted=false and no
	// error usually means success, but we refuse to guess — ungrounded.
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

// CaptureEdit appends the edit→task binding ledger (never blocks).
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
