package gates

import (
	"testing"
)

// Regression for the live T-3/T-4 duplication: the agent records T-1/T-2,
// then mirrors native tasks with "T-1: <subject>" prefixes — the gate must
// link them, never create duplicates.
func TestTaskCreatedMatchingLadder(t *testing.T) {
	c := newCtl(t)
	run, _ := c.RunStart("diff", "new")
	run.Phase = "plan"
	_ = c.Store.SaveRun(run)
	rec := func(tid, subject string) {
		t.Helper()
		if _, err := c.Record("task", map[string]any{
			"tid": tid, "subject": subject, "status": "open",
			"dod": []any{"d"}, "ac_links": []any{"AC-1"},
		}, false, "agent"); err != nil {
			t.Fatal(err)
		}
	}
	rec("T-1", "Write failing test for prime output")
	rec("T-2", "Implement primes.py to satisfy AC-1/AC-2")

	// exactly the live-session native subjects
	for id, subject := range map[string]string{
		"1": "T-1: Write failing test for prime output",
		"2": "T-2: Implement primes.py to satisfy AC-1/AC-2",
	} {
		if r := TaskCreated(c, hookInput(t, taskPayload(id, subject, "dod"))); r.Exit != 0 {
			t.Fatalf("task %s rejected: %s", id, r.Stderr)
		}
	}
	env, _ := c.Env(run)
	if n := len(env.Records("task")); n != 2 {
		t.Fatalf("duplication bug: want 2 task records, got %d", n)
	}

	// exact match without prefix still links
	if r := TaskCreated(c, hookInput(t, taskPayload("3", "write failing test for PRIME output", "d"))); r.Exit != 0 {
		t.Fatal(r.Stderr)
	}
	env, _ = c.Env(run)
	if n := len(env.Records("task")); n != 2 {
		t.Fatalf("case-insensitive exact match failed: %d records", n)
	}

	// containment: a slightly expanded native subject links to the record
	if r := TaskCreated(c, hookInput(t, taskPayload("4", "Implement primes.py to satisfy AC-1/AC-2 (comma-separated output)", "d"))); r.Exit != 0 {
		t.Fatal(r.Stderr)
	}
	env, _ = c.Env(run)
	if n := len(env.Records("task")); n != 2 {
		t.Fatalf("containment match failed: %d records", n)
	}

	// a genuinely new subject still creates a record
	if r := TaskCreated(c, hookInput(t, taskPayload("5", "totally unrelated cleanup chore", "d"))); r.Exit != 0 {
		t.Fatal(r.Stderr)
	}
	env, _ = c.Env(run)
	if n := len(env.Records("task")); n != 3 {
		t.Fatalf("new task must still be creatable: %d records", n)
	}

	// and completion via the tid-token path closes the RIGHT record
	_, _ = c.Record("test-run", map[string]any{"cmd": "pytest", "exit": 1, "grounded": true, "task": "T-1"}, true, "hook")
	_, _ = c.Record("test-run", map[string]any{"cmd": "pytest", "exit": 0, "grounded": true, "task": "T-1"}, true, "hook")
	done := `{"hook_event_name":"TaskCompleted","task_id":"1","task_subject":"T-1: Write failing test for prime output"}`
	if r := TaskCompleted(c, hookInput(t, done)); r.Exit != 0 {
		t.Fatalf("completion via mirror failed: %s", r.Stderr)
	}
	env, _ = c.Env(run)
	for _, tr := range env.Records("task") {
		if tid, _ := tr.Data["tid"].(string); tid == "T-1" {
			if s, _ := tr.Data["status"].(string); s != "done" {
				t.Errorf("T-1 must be done, got %s", s)
			}
		}
	}
}
