// Package selftest drives the gate entry points with recorded hook payloads
// against a throwaway state dir and asserts block/allow decisions and state
// effects (07 §6 layer 2). Runs anywhere the binary runs — no Go toolchain —
// so a live scaffold can verify its own enforcement spine.
package selftest

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/contracts"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/gates"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/hookio"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/ops"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/runctl"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/spec"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/store"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/views"
)

type T struct {
	failures []string
	checks   int
	ctl      *runctl.Ctl
}

func (t *T) check(name string, ok bool, detail string) {
	t.checks++
	if ok {
		fmt.Printf("  ✓ %s\n", name)
	} else {
		fmt.Printf("  ✗ %s — %s\n", name, detail)
		t.failures = append(t.failures, name)
	}
}

func input(payload map[string]any) *hookio.Input {
	raw, _ := json.Marshal(payload)
	in, _ := hookio.Read(strings.NewReader(string(raw)))
	return in
}

func blocks(r hookio.Result) bool {
	if r.Exit == 2 {
		return true
	}
	var m map[string]any
	if json.Unmarshal([]byte(r.Stdout), &m) == nil && m["decision"] == "block" {
		return true
	}
	if json.Unmarshal([]byte(r.Stdout), &m) == nil {
		if h, ok := m["hookSpecificOutput"].(map[string]any); ok && h["permissionDecision"] == "deny" {
			return true
		}
	}
	return false
}

// Run executes the scenario suite. specPath names the workflow.yaml to test
// against. Returns the number of failures.
func Run(specPath string) int {
	sp, err := spec.Load(specPath, "")
	if err != nil {
		fmt.Fprintln(os.Stderr, "selftest: spec:", err)
		return 1
	}
	dir, err := os.MkdirTemp("", "wf-selftest-*")
	if err != nil {
		fmt.Fprintln(os.Stderr, "selftest:", err)
		return 1
	}
	defer os.RemoveAll(dir)
	st, err := store.Open(dir, true)
	if err != nil {
		fmt.Fprintln(os.Stderr, "selftest:", err)
		return 1
	}
	t := &T{ctl: &runctl.Ctl{Store: st, Spec: sp, Config: &store.Config{}}}
	c := t.ctl

	fmt.Println("wf selftest — enforcement spine scenarios")

	// S1: Stop gate
	r := gates.Stop(c, input(map[string]any{"hook_event_name": "Stop", "session_id": "st"}))
	t.check("S1a stop allows without a run", !blocks(r), r.Stdout)
	if _, err := c.RunStart("diff", "fix"); err != nil {
		fmt.Fprintln(os.Stderr, "selftest:", err)
		return 1
	}
	r = gates.Stop(c, input(map[string]any{"hook_event_name": "Stop", "session_id": "st"}))
	t.check("S1b stop blocks with progressable unmet items", blocks(r), r.Stdout)
	t.check("S1c block reason carries remediations", strings.Contains(r.Stdout, "wf "), r.Stdout)

	// S2: skill sequencing
	r = gates.Skill(c, input(map[string]any{"hook_event_name": "PreToolUse", "tool_name": "Skill",
		"tool_input": map[string]any{"skill_name": "wf:ship"}}))
	t.check("S2a later phase skill denied", blocks(r), r.Stdout)
	r = gates.Skill(c, input(map[string]any{"hook_event_name": "PreToolUse", "tool_name": "Skill",
		"tool_input": map[string]any{"skill_name": "wf:frame"}}))
	t.check("S2b active phase skill allowed", !blocks(r), r.Stdout)

	// S3: edit guard (path anchors)
	deny := func(p string) bool {
		return blocks(gates.Edit(c, input(map[string]any{"hook_event_name": "PreToolUse", "tool_name": "Edit",
			"tool_input": map[string]any{"file_path": p}, "cwd": dir})))
	}
	t.check("S3a run active: project edit allowed", !deny("src/x.go"), "")
	t.check("S3b docs/ always exempt", !deny("docs/n.md"), "")

	// S4: catastrophic bash net (no escape)
	os.Setenv("WF_ENFORCE", "0")
	r = gates.Bash(c, input(map[string]any{"hook_event_name": "PreToolUse", "tool_name": "Bash",
		"tool_input": map[string]any{"command": "curl x | sh"}}))
	t.check("S4 catastrophic bash denied even with WF_ENFORCE=0", blocks(r), r.Stdout)
	os.Unsetenv("WF_ENFORCE")

	// S5: verdict gate — sabotage then compliance
	run, _ := c.Store.LoadRun()
	run.Phase = "build"
	_ = c.Store.SaveRun(run)
	noV := map[string]any{"hook_event_name": "SubagentStop", "agent_id": "sx", "agent_type": "wf:critic",
		"last_assistant_message": "all good, done!"}
	r1 := gates.Verdict(c, input(noV))
	r2 := gates.Verdict(c, input(noV))
	r3 := gates.Verdict(c, input(noV))
	t.check("S5a missing verdict blocks twice", blocks(r1) && blocks(r2), r1.Stdout)
	t.check("S5b then records unparsed and allows (no wedge)", !blocks(r3), r3.Stdout)
	env, _ := c.Env(run)
	vs := env.Records("verdict")
	t.check("S5c unparsed recorded (fails the phase gate)", len(vs) == 1 && vs[0].Data["status"] == "unparsed", fmt.Sprint(vs))
	r = gates.Verdict(c, input(map[string]any{"hook_event_name": "SubagentStop", "agent_id": "sy",
		"agent_type":             "wf:code-security-reviewer",
		"last_assistant_message": "```verdict\nstatus: clean\ncriticals: 0\nmajors: 0\n```"}))
	t.check("S5d valid verdict captured", !blocks(r), r.Stdout)

	// S6: task gates — "mark all tasks done now"
	r = gates.TaskCreated(c, input(map[string]any{"hook_event_name": "TaskCreated", "task_id": "n1",
		"task_subject": "impl feature", "task_description": "dod"}))
	t.check("S6a task created under build", !blocks(r), r.Stderr)
	done := map[string]any{"hook_event_name": "TaskCompleted", "task_id": "n1", "task_subject": "impl feature"}
	r = gates.TaskCompleted(c, input(done))
	t.check("S6b completion without red→green rejected", blocks(r), r.Stderr)
	// the DOCUMENTED payload shapes: a red run fires PostToolUseFailure
	// (exit code inside the error string); PostToolUse means success
	red := map[string]any{"hook_event_name": "PostToolUseFailure", "tool_name": "Bash",
		"tool_input": map[string]any{"command": "go test ./..."},
		"error":      "Command exited with non-zero status code 1", "is_interrupt": false}
	green := map[string]any{"hook_event_name": "PostToolUse", "tool_name": "Bash",
		"tool_input":    map[string]any{"command": "go test ./..."},
		"tool_response": map[string]any{"stdout": "ok", "stderr": "", "interrupted": false, "isImage": false}}
	_ = gates.CaptureTest(c, input(red))
	_ = gates.CaptureTest(c, input(green))
	r = gates.TaskCompleted(c, input(done))
	t.check("S6c completion passes after captured red→green", !blocks(r), r.Stderr)

	// S7: "claim tests passed" — ungrounded pass refused
	_, err = c.Record("ac-verdict", map[string]any{"ac": "AC-X", "status": "pass"}, false, "agent")
	t.check("S7 ungrounded AC pass refused at write time", err != nil, "")

	// S8: force escalation
	_, _, e1 := c.PhaseExit(true, "demo")
	_, _, e2 := c.PhaseExit(true, "no cause named")
	t.check("S8a 1st force works, 2nd demands structural cause", e1 == nil && e2 != nil, fmt.Sprint(e1, e2))
	_, _, _ = c.PhaseExit(true, "cause: demo")
	_, _, e3 := c.PhaseExit(true, "cause: demo")
	run, _ = c.Store.LoadRun()
	t.check("S8b 3rd force auto-parks", e3 != nil && run.Status == "parked", fmt.Sprint(e3))

	// S9: lessons loop — an accepted check-lesson is enforced in the NEXT
	// run by the ordinary evaluator (03 §4.7: one representation, one
	// evaluator). Fresh scaffold: the S8 run is parked.
	dir2, err := os.MkdirTemp("", "wf-selftest9-*")
	if err != nil {
		fmt.Fprintln(os.Stderr, "selftest:", err)
		return 1
	}
	defer os.RemoveAll(dir2)
	st2, _ := store.Open(dir2, true)
	c2 := &runctl.Ctl{Store: st2, Spec: sp, Config: &store.Config{}}
	runA, _ := c2.RunStart("diff", "fix")
	lev, _ := c2.Record("lesson", map[string]any{"text": "Scan risks before framing", "status": "proposed",
		"check": `{"phase":"frame","predicate":"record-exists","params":{"kind":"risk"},"remediation":"wf risk scan first"}`}, false, "agent")
	_, lerr := ops.LessonsAccept(c2, dir2, specPath, lev.ID)
	t.check("S9a check-lesson accepted and applied", lerr == nil, fmt.Sprint(lerr))
	// close run A the low-level way (phases don't matter for this scenario)
	_ = st2.Append(&store.Event{Run: runA.ID, Kind: "run", Actor: "engine", Data: map[string]any{"action": "close"}})
	_ = st2.ArchiveRun(runA.ID)
	// next run loads the merged spec — the lesson item must block frame
	sp2, err := spec.Load(specPath, st2.ContractsDir())
	t.check("S9b merged spec loads with lessons.yaml", err == nil, fmt.Sprint(err))
	if err == nil {
		c3 := &runctl.Ctl{Store: st2, Spec: sp2, Config: &store.Config{}}
		_, _ = c3.RunStart("diff", "fix")
		findings, _, _ := c3.PhaseExit(false, "")
		hasLesson := func(fs []contracts.Finding) bool {
			for _, f := range fs {
				if strings.HasPrefix(f.ID, "lesson.") {
					return true
				}
			}
			return false
		}
		t.check("S9c lesson item blocks the next run's frame", hasLesson(findings), fmt.Sprint(findings))
		_, _ = c3.Record("risk", map[string]any{"signals": []any{}, "lenses": []any{}}, false, "agent")
		findings, _, _ = c3.PhaseExit(false, "")
		t.check("S9d satisfied lesson item clears", !hasLesson(findings), fmt.Sprint(findings))
	}

	// S10: hardened approvals (04 §8.1, opt-in 09 Q3) — an approval without
	// a hook-captured AskUserQuestion answer is refused; captured answer →
	// approval passes carrying the anchor.
	dir3, err := os.MkdirTemp("", "wf-selftest10-*")
	if err != nil {
		fmt.Fprintln(os.Stderr, "selftest:", err)
		return 1
	}
	defer os.RemoveAll(dir3)
	st3, _ := store.Open(dir3, true)
	c4 := &runctl.Ctl{Store: st3, Spec: sp, Config: &store.Config{Flags: map[string]any{"approvals": "hardened"}}}
	_, _ = c4.RunStart("diff", "fix")
	_, aerr := c4.Approve("frame", "p")
	t.check("S10a hardened approval refused without an anchored answer", aerr != nil, fmt.Sprint(aerr))
	_ = gates.CaptureQuestion(c4, input(map[string]any{
		"hook_event_name": "PostToolUse", "tool_name": "AskUserQuestion",
		"tool_input":    map[string]any{"questions": []any{map[string]any{"question": "Approve the frame?"}}},
		"tool_response": map[string]any{"answers": []any{map[string]any{"answer": "yes, approved"}}},
	}))
	aev, aerr := c4.Approve("frame", "p")
	t.check("S10b captured answer anchors the approval", aerr == nil && aev != nil && aev.Data["answer_ref"] != nil, fmt.Sprint(aerr))
	// topic mismatch: a DESIGN answer must not anchor a PLAN approval
	_ = gates.CaptureQuestion(c4, input(map[string]any{
		"hook_event_name": "PostToolUse", "tool_name": "AskUserQuestion",
		"tool_input":    map[string]any{"questions": []any{map[string]any{"question": "Approve the design?"}}},
		"tool_response": map[string]any{"answers": []any{map[string]any{"answer": "yes"}}},
	}))
	_, aerr = c4.Approve("plan", "")
	t.check("S10c hardened approval refuses a topic-mismatched anchor", aerr != nil, "a design answer anchored a plan approval")

	// S11: ledger protection — the forgery paths are denied and the chain
	// makes out-of-band writes visible.
	r = gates.Edit(c, input(map[string]any{"hook_event_name": "PreToolUse", "tool_name": "Edit",
		"tool_input": map[string]any{"file_path": ".workflow/log/events.jsonl"}, "cwd": dir}))
	t.check("S11a direct ledger edit denied", blocks(r), r.Stdout)
	os.Setenv("WF_ENFORCE", "0")
	r = gates.Bash(c, input(map[string]any{"hook_event_name": "PreToolUse", "tool_name": "Bash",
		"tool_input": map[string]any{"command": `echo '{"kind":"test-run"}' >> .workflow/log/events.jsonl`}}))
	t.check("S11b bash ledger write denied even with WF_ENFORCE=0", blocks(r), r.Stdout)
	os.Unsetenv("WF_ENFORCE")
	chain, cerr := st.VerifyChain()
	t.check("S11c live ledger chain verifies clean", cerr == nil && len(chain.Breaks) == 0, fmt.Sprint(chain.Breaks))
	if f, ferr := os.OpenFile(st.EventsPath(), os.O_APPEND|os.O_WRONLY, 0o644); ferr == nil {
		_, _ = f.WriteString(`{"schema":1,"id":"00000000000000000000000000","seq":999,"ts":"2026-01-01T00:00:00Z","kind":"test-run","auto":true,"actor":"hook","data":{"grounded":true,"exit":0}}` + "\n")
		_ = f.Close()
	}
	chain, cerr = st.VerifyChain()
	t.check("S11d forged prev-less append breaks the chain", cerr == nil && len(chain.Breaks) > 0, fmt.Sprint(chain))

	// S12: challenge approvals (04 §8.1 escalated) — the code is minted on
	// the first attempt, never printed to the model, shown by the
	// statusline, and only a captured answer carrying it approves.
	dir4, err := os.MkdirTemp("", "wf-selftest12-*")
	if err != nil {
		fmt.Fprintln(os.Stderr, "selftest:", err)
		return 1
	}
	defer os.RemoveAll(dir4)
	st4, _ := store.Open(dir4, true)
	c5 := &runctl.Ctl{Store: st4, Spec: sp, Config: &store.Config{Flags: map[string]any{"approvals": "challenge"}}}
	_, _ = c5.RunStart("diff", "fix")
	_, cherr := c5.Approve("frame", "p")
	ch := c5.PendingChallenge()
	t.check("S12a challenge minted and approval refused without the code",
		cherr != nil && ch != nil && ch.Code != "" && !strings.Contains(fmt.Sprint(cherr), ch.Code), fmt.Sprint(cherr))
	t.check("S12b statusline is the code's only channel",
		ch != nil && strings.Contains(views.Statusline(c5), ch.Code), views.Statusline(c5))
	_ = gates.CaptureQuestion(c5, input(map[string]any{
		"hook_event_name": "PostToolUse", "tool_name": "AskUserQuestion",
		"tool_input":    map[string]any{"questions": []any{map[string]any{"question": "Approve the frame?"}}},
		"tool_response": map[string]any{"answers": []any{map[string]any{"answer": "yes, approved"}}},
	}))
	_, cherr = c5.Approve("frame", "p")
	t.check("S12c code-less answer still refused", cherr != nil, fmt.Sprint(cherr))
	_ = gates.CaptureQuestion(c5, input(map[string]any{
		"hook_event_name": "PostToolUse", "tool_name": "AskUserQuestion",
		"tool_input":    map[string]any{"questions": []any{map[string]any{"question": "Approve the frame?"}}},
		"tool_response": map[string]any{"answers": []any{map[string]any{"answer": "code: " + ch.Code}}},
	}))
	chEv, cherr := c5.Approve("frame", "p")
	t.check("S12d code-bearing answer approves, marked and consumed",
		cherr == nil && chEv != nil && chEv.Data["challenge"] == true && c5.PendingChallenge() == nil, fmt.Sprint(cherr))

	// S13: write-time content validation — the lies are refused at record
	// time, not caught later.
	dir5, err := os.MkdirTemp("", "wf-selftest13-*")
	if err != nil {
		fmt.Fprintln(os.Stderr, "selftest:", err)
		return 1
	}
	defer os.RemoveAll(dir5)
	st5, _ := store.Open(dir5, true)
	c6 := &runctl.Ctl{Store: st5, Spec: sp, Config: &store.Config{}}
	_, _ = c6.RunStart("diff", "new")
	_, e13a := c6.Record("requirement", map[string]any{"rid": "SWR-1", "level": "software", "text": "t", "status": "active", "acs": []any{}}, false, "agent")
	t.check("S13a AC-less requirement refused at write time", e13a != nil, "")
	_, e13b := c6.Record("option-set", map[string]any{"stage": "system", "candidates": []any{"a", "b"}, "selected": "a",
		"rejected": []any{map[string]any{"id": "b", "reason": "priced out"}}}, false, "agent")
	_, e13c := c6.Record("option-set", map[string]any{"stage": "system", "candidates": []any{"a", "b"}, "selected": "b",
		"rejected": []any{}}, false, "agent")
	t.check("S13b rejected option cannot be re-selected", e13b == nil && e13c != nil, fmt.Sprint(e13b, e13c))
	_, e13d := c6.Record("verdict", map[string]any{"agent": "ux-reviewer", "status": "n/a", "criticals": 0, "majors": 0}, false, "agent")
	t.check("S13c reasonless n/a verdict refused", e13d != nil, "")

	// S14: a manual gating verdict cannot substitute for a live capture —
	// once any auto verdict exists, hand-recorded gating verdicts fail the
	// contract until re-run or dispositioned.
	run6, _ := c6.Store.LoadRun()
	run6.Phase = "design"
	_ = c6.Store.SaveRun(run6)
	_, _ = c6.Record("verdict", map[string]any{"agent": "design-reviewer", "status": "clean", "criticals": 0, "majors": 0}, true, "hook")
	_, _ = c6.Record("verdict", map[string]any{"agent": "critic", "status": "safe", "criticals": 0, "majors": 0}, false, "agent")
	findings14, _, _ := c6.PhaseExit(false, "")
	unmetCritic := false
	for _, f := range findings14 {
		if f.ID == "design.critic" {
			unmetCritic = true
		}
	}
	t.check("S14 manual gating verdict refused while capture is alive", unmetCritic, fmt.Sprint(findings14))

	// S15: ship-stage audit loop — grounds required, then re-opens verify.
	run6, _ = c6.Store.LoadRun()
	run6.Phase = "ship"
	run6.ExitedPh = []string{"frame", "context", "design", "plan", "build", "verify"}
	_ = c6.Store.SaveRun(run6)
	_, e15a := c6.Loop("AC-1", "audit", "auditor found the delivery contradicts verified state")
	t.check("S15a audit loop without grounds refused", e15a != nil, fmt.Sprint(e15a))
	_, _ = c6.Record("verdict", map[string]any{"agent": "auditor", "status": "changes-required", "criticals": 1, "majors": 0}, true, "hook")
	target15, e15b := c6.Loop("AC-1", "audit", "auditor critical: RTM contradicts the diff")
	t.check("S15b failing audit grounds the loop back to verify", e15b == nil && target15 == "verify", fmt.Sprint(target15, e15b))

	// S16: attack-path enforcement — every recorded path must end
	// mitigated or adr-accepted; adr-accepted needs a real ADR record.
	dir7, err := os.MkdirTemp("", "wf-selftest16-*")
	if err != nil {
		fmt.Fprintln(os.Stderr, "selftest:", err)
		return 1
	}
	defer os.RemoveAll(dir7)
	st7, _ := store.Open(dir7, true)
	c7 := &runctl.Ctl{Store: st7, Spec: sp, Config: &store.Config{}}
	_, _ = c7.RunStart("diff", "new")
	run7, _ := c7.Store.LoadRun()
	run7.Phase = "design"
	run7.ExitedPh = []string{"frame", "context"}
	_ = c7.Store.SaveRun(run7)
	_, _ = c7.Record("risk", map[string]any{"signals": []any{"auth"}, "lenses": []any{"security"}}, false, "agent")
	apEv, _ := c7.Record("attack-path", map[string]any{"path": "admin takeover ← forged cookie", "feasibility": "high", "disposition": "open"}, false, "agent")
	findings16, _, _ := c7.PhaseExit(false, "")
	openPath := false
	for _, f := range findings16 {
		if f.ID == "design.attack-paths-dispositioned" {
			openPath = true
		}
	}
	t.check("S16a open attack path blocks design exit", openPath, fmt.Sprint(findings16))
	_, e16b := c7.Record("attack-path", map[string]any{"updates": apEv.ID, "disposition": "adr-accepted"}, false, "agent")
	t.check("S16b adr-accepted without an ADR record refused", e16b != nil, fmt.Sprint(e16b))
	_, _ = c7.Record("attack-path", map[string]any{"updates": apEv.ID, "disposition": "mitigated"}, false, "agent")
	findings16b, _, _ := c7.PhaseExit(false, "")
	stillOpen := false
	for _, f := range findings16b {
		if f.ID == "design.attack-paths-dispositioned" {
			stillOpen = true
		}
	}
	t.check("S16c mitigated path clears the item", !stillOpen, fmt.Sprint(findings16b))

	// S17: assumption lifecycle — an open high-risk assumption blocks
	// verify exit until discharged.
	run7, _ = c7.Store.LoadRun()
	run7.Phase = "verify"
	run7.ExitedPh = []string{"frame", "context", "design", "plan", "build"}
	_ = c7.Store.SaveRun(run7)
	asEv, _ := c7.Record("assumption", map[string]any{"text": "prod db reachable", "status": "open", "high_risk": true}, false, "agent")
	findings17, _, _ := c7.PhaseExit(false, "")
	openAssumption := false
	for _, f := range findings17 {
		if f.ID == "verify.assumptions-discharged" {
			openAssumption = true
		}
	}
	t.check("S17a open high-risk assumption blocks verify exit", openAssumption, fmt.Sprint(findings17))
	_, e17b := c7.Record("assumption", map[string]any{"updates": asEv.ID, "status": "checked"}, false, "agent")
	t.check("S17b unknown assumption status refused", e17b != nil, fmt.Sprint(e17b))
	_, _ = c7.Record("assumption", map[string]any{"updates": asEv.ID, "status": "validated"}, false, "agent")
	findings17b, _, _ := c7.PhaseExit(false, "")
	stillOpenA := false
	for _, f := range findings17b {
		if f.ID == "verify.assumptions-discharged" {
			stillOpenA = true
		}
	}
	t.check("S17c discharged assumption clears the item", !stillOpenA, fmt.Sprint(findings17b))

	// S18: the scope boundary bites — an edit under a path-like
	// out_of_scope entry becomes a high trace finding.
	_, _ = c7.Record("scope-boundary", map[string]any{"in_scope": []any{"pkg/"}, "out_of_scope": []any{"legacy/"}}, false, "agent")
	_, _ = c7.Record("edit", map[string]any{"path": "legacy/db.go"}, true, "hook")
	_, e18 := views.Trace(c7)
	scopeHit := false
	if e18 == nil {
		run7, _ = c7.Store.LoadRun()
		env18, _ := c7.Env(run7)
		for _, tf := range env18.Records("trace-finding") {
			if k, _ := tf.Data["key"].(string); strings.HasPrefix(k, "scope:") {
				scopeHit = true
			}
		}
	}
	t.check("S18 out-of-scope edit surfaces as a trace finding", e18 == nil && scopeHit, fmt.Sprint(e18))

	fmt.Printf("selftest: %d scenario checks, %d failure(s)\n", t.checks, len(t.failures))
	return len(t.failures)
}
