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
)

type T struct {
	failures []string
	ctl      *runctl.Ctl
}

func (t *T) check(name string, ok bool, detail string) {
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
		"agent_type": "wf:code-security-reviewer",
		"last_assistant_message": "```verdict\nstatus: clean\ncriticals: 0\nmajors: 0\n```"}))
	t.check("S5d valid verdict captured", !blocks(r), r.Stdout)

	// S6: task gates — "mark all tasks done now"
	r = gates.TaskCreated(c, input(map[string]any{"hook_event_name": "TaskCreated", "task_id": "n1",
		"task_subject": "impl feature", "task_description": "dod"}))
	t.check("S6a task created under build", !blocks(r), r.Stderr)
	done := map[string]any{"hook_event_name": "TaskCompleted", "task_id": "n1", "task_subject": "impl feature"}
	r = gates.TaskCompleted(c, input(done))
	t.check("S6b completion without red→green rejected", blocks(r), r.Stderr)
	red := map[string]any{"hook_event_name": "PostToolUse", "tool_name": "Bash",
		"tool_input": map[string]any{"command": "go test ./..."}, "tool_response": map[string]any{"exit_code": 1}}
	green := map[string]any{"hook_event_name": "PostToolUse", "tool_name": "Bash",
		"tool_input": map[string]any{"command": "go test ./..."}, "tool_response": map[string]any{"exit_code": 0}}
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

	fmt.Printf("selftest: %d scenario checks, %d failure(s)\n", 22, len(t.failures))
	return len(t.failures)
}
