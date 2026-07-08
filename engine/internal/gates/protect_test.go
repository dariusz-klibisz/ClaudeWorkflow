package gates

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/hookio"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/runctl"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/store"
)

func editPayload(t *testing.T, path string) string {
	t.Helper()
	raw, _ := json.Marshal(map[string]any{
		"session_id": "s1", "hook_event_name": "PreToolUse", "tool_name": "Edit",
		"tool_input": map[string]any{"file_path": path},
	})
	return string(raw)
}

func bashPayload(t *testing.T, cmd string) string {
	t.Helper()
	raw, _ := json.Marshal(map[string]any{
		"session_id": "s1", "hook_event_name": "PreToolUse", "tool_name": "Bash",
		"tool_input": map[string]any{"command": cmd},
	})
	return string(raw)
}

// ---------------------------------------------------------------------------
// Protected-state Edit gate (the ledger-forgery fix)
// ---------------------------------------------------------------------------

func TestEditDeniesProtectedState(t *testing.T) {
	c := newCtl(t)
	for _, p := range []string{
		".workflow/log/events.jsonl",
		".workflow/state/run.json",
		".workflow/runs/20260101-abc/events.jsonl",
		".workflow/config.json",
	} {
		if r := Edit(c, hookInput(t, editPayload(t, p))); !isDeny(r) {
			t.Errorf("edit of %s must be denied: %+v", p, r)
		}
	}
}

func TestEditProtectedStateIgnoresEnforceOff(t *testing.T) {
	c := newCtl(t)
	t.Setenv("WF_ENFORCE", "0")
	if r := Edit(c, hookInput(t, editPayload(t, ".workflow/log/events.jsonl"))); !isDeny(r) {
		t.Fatal("protected-state deny is data protection: WF_ENFORCE must not lift it")
	}
}

func TestEditAllowsUnprotectedWorkflowPaths(t *testing.T) {
	c := newCtl(t)
	for _, p := range []string{
		".workflow/contracts.d/local-extra.yaml",
		".workflow/.gitignore",
		"docs/design/x.md",
	} {
		if r := Edit(c, hookInput(t, editPayload(t, p))); isDeny(r) {
			t.Errorf("edit of %s must stay allowed: %+v", p, r)
		}
	}
}

// ---------------------------------------------------------------------------
// Bash state-tamper net
// ---------------------------------------------------------------------------

func TestBashDeniesStateTamper(t *testing.T) {
	c := newCtl(t)
	for _, cmd := range []string{
		`echo '{"kind":"test-run","auto":true}' >> .workflow/log/events.jsonl`,
		`printf '{}' > .workflow/state/run.json`,
		`echo x > .workflow/runs/20260101-abc/events.jsonl`,
		`cat forged.json >> "` + `.workflow/log/events.jsonl"`,
		`echo y | tee -a .workflow/log/events.jsonl`,
		`sed -i 's/hardened//' .workflow/config.json`,
		`cp /tmp/forged.jsonl .workflow/log/events.jsonl`,
		`mv /tmp/run.json .workflow/state/run.json`,
		`rm .workflow/log/events.jsonl`,
		`git status && echo x >> .workflow/log/events.jsonl`,
		`truncate -s 0 .workflow/log/events.jsonl`,
	} {
		if r := Bash(c, hookInput(t, bashPayload(t, cmd))); !isDeny(r) {
			t.Errorf("tamper must be denied: %s", cmd)
		}
	}
}

func TestBashTamperIgnoresEnforceOff(t *testing.T) {
	c := newCtl(t)
	t.Setenv("WF_ENFORCE", "0")
	cmd := `echo x >> .workflow/log/events.jsonl`
	if r := Bash(c, hookInput(t, bashPayload(t, cmd))); !isDeny(r) {
		t.Fatal("state-tamper net must ignore WF_ENFORCE")
	}
}

func TestBashAllowsLegitimateStateAccess(t *testing.T) {
	c := newCtl(t)
	for _, cmd := range []string{
		`cat .workflow/log/events.jsonl`,
		`grep test-run .workflow/log/events.jsonl | head -5`,
		`cat .workflow/log/events.jsonl > /tmp/backup.jsonl`,
		`git add .workflow/log/events.jsonl && git commit -m "x [run:r1]"`,
		`wf record test-run --cmd pytest --exit 0`,
		`ls .workflow/runs/`,
		`echo hello > /tmp/out.txt`,
		`tee /tmp/notes.txt < input`,
		`rm -f build/artifacts.zip`,
		`jq . .workflow/state/run.json`,
	} {
		if r := Bash(c, hookInput(t, bashPayload(t, cmd))); isDeny(r) {
			t.Errorf("legitimate command must stay allowed: %s (%+v)", cmd, r)
		}
	}
}

// ---------------------------------------------------------------------------
// Escape events (WS-B: no silent escapes)
// ---------------------------------------------------------------------------

func escapesByAction(t *testing.T, c *runctl.Ctl, action string) int {
	t.Helper()
	r, err := c.Store.LoadRun()
	if err != nil || r == nil {
		// run may be mid-mutation; read all events instead
		evs, _ := c.Store.Events(nil)
		n := 0
		for _, e := range evs {
			if e.Kind == "escape" && e.Str("action") == action {
				n++
			}
		}
		return n
	}
	evs, _ := c.Store.RunEvents(r.ID)
	n := 0
	for _, e := range evs {
		if e.Kind == "escape" && e.Str("action") == action {
			n++
		}
	}
	return n
}

func TestStopCapReleaseRecordsEscape(t *testing.T) {
	c := newCtl(t)
	_, _ = c.RunStart("diff", "fix")
	active := `{"session_id":"s1","hook_event_name":"Stop","stop_hook_active":true}`
	var last hookio.Result
	for i := 0; i < 6; i++ {
		last = Stop(c, hookInput(t, active))
	}
	if isBlockDecision(last) {
		t.Fatal("self-cap must have released")
	}
	if n := escapesByAction(t, c, "stop-cap"); n != 1 {
		t.Fatalf("exactly one stop-cap escape event must be recorded, got %d", n)
	}
	if !strings.Contains(last.Stdout, "recorded") {
		t.Errorf("release message should say it was recorded: %s", last.Stdout)
	}
}

func TestEnforceOffRecordsEscapeOncePerGate(t *testing.T) {
	c := newCtl(t)
	_, _ = c.RunStart("diff", "fix")
	t.Setenv("WF_ENFORCE", "0")
	stop := `{"session_id":"s1","hook_event_name":"Stop"}`
	_ = Stop(c, hookInput(t, stop))
	_ = Stop(c, hookInput(t, stop)) // second downgrade: rate-limited
	_ = Edit(c, hookInput(t, editPayload(t, "src/app.go")))
	if n := escapesByAction(t, c, "enforce-off"); n != 2 {
		t.Fatalf("want 2 enforce-off escapes (stop + edit, once each), got %d", n)
	}
}

// ---------------------------------------------------------------------------
// Topic-anchored approvals (WS-C)
// ---------------------------------------------------------------------------

func questionPayload(t *testing.T, question, answer string) string {
	raw, _ := json.Marshal(map[string]any{
		"hook_event_name": "PostToolUse", "tool_name": "AskUserQuestion",
		"tool_input":    map[string]any{"questions": []any{map[string]any{"question": question}}},
		"tool_response": map[string]any{"answers": []any{map[string]any{"answer": answer}}},
	})
	return string(raw)
}

func TestQuestionTopicInference(t *testing.T) {
	cases := map[string]string{
		"Do you approve the design option B?":       "design",
		"Approve the plan (task breakdown)?":        "plan",
		"Is the classification diff/fix correct?":   "frame",
		"OK to defer AC-3?":                         "deferral",
		"Shall we proceed?":                         "", // no gate words
		"Approve the design and the plan together?": "", // ambiguous
	}
	for q, want := range cases {
		if got := questionTopic(q); got != want {
			t.Errorf("%q: topic %q, want %q", q, got, want)
		}
	}
}

func TestApproveAnchorsToMatchingTopic(t *testing.T) {
	c := newCtl(t)
	c.Config = &store.Config{Flags: map[string]any{"approvals": "hardened"}}
	_, _ = c.RunStart("diff", "fix")

	// only a DESIGN answer exists — a hardened PLAN approval must refuse
	_ = CaptureQuestion(c, hookInput(t, questionPayload(t, "Do you approve the design?", "yes")))
	if _, err := c.Approve("plan", ""); err == nil {
		t.Fatal("hardened approval must refuse a topic-mismatched anchor")
	}
	// the design approval anchors to it
	ev, err := c.Approve("design", "")
	if err != nil || ev.Data["answer_ref"] == nil {
		t.Fatalf("design approval must anchor: %v %v", err, ev)
	}
	// a topicless answer anchors anything
	_ = CaptureQuestion(c, hookInput(t, questionPayload(t, "Shall we proceed?", "go ahead")))
	ev, err = c.Approve("plan", "")
	if err != nil || ev.Data["answer_ref"] == nil {
		t.Fatalf("topicless answer must still anchor: %v %v", err, ev)
	}
}
