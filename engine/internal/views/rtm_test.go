package views

import (
	"strings"
	"testing"
)

func TestRTMRendersChain(t *testing.T) {
	c := newCtl(t)
	r, err := c.RunStart("diff", "new")
	if err != nil {
		t.Fatal(err)
	}
	id := r.ID

	raw(t, c, id, "frame", "requirement", false, "agent", map[string]any{
		"rid": "SWR-1", "level": "software", "status": "active",
		"text": "parse uploads safely",
		"acs": []any{
			map[string]any{"id": "AC-1", "text": "rejects oversized files", "verifiable": true},
			map[string]any{"id": "AC-2", "text": "accepts valid files", "verifiable": true},
		},
	})
	raw(t, c, id, "plan", "verification-strategy", false, "agent", map[string]any{
		"ac": "AC-1", "method": "test", "command": "go test ./upload/",
	})
	raw(t, c, id, "plan", "task", false, "agent", map[string]any{
		"tid": "T-1", "subject": "size guard", "dod": "d", "status": "done",
		"ac_links": []any{"AC-1"},
	})
	// evidence: grounded red then green for AC-1; nothing for AC-2
	raw(t, c, id, "build", "test-run", true, "hook", map[string]any{"cmd": "go test ./upload/", "exit": 1.0, "grounded": true, "ac": "AC-1"})
	raw(t, c, id, "build", "test-run", true, "hook", map[string]any{"cmd": "go test ./upload/", "exit": 0.0, "grounded": true, "ac": "AC-1"})
	raw(t, c, id, "verify", "ac-verdict", false, "agent", map[string]any{"ac": "AC-1", "status": "pass"})
	raw(t, c, id, "verify", "loop", false, "agent", map[string]any{"ac": "AC-2", "cause": "slip", "evidence": "e", "target": "build"})
	// orphan evidence: tagged to an AC no requirement declares
	raw(t, c, id, "build", "test-run", true, "hook", map[string]any{"cmd": "go test ./ghost/", "exit": 0.0, "grounded": true, "ac": "AC-9"})

	out, err := RTM(c)
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"# RTM — run " + id,
		"## SWR-1 — parse uploads safely",
		"level: software",
		"| AC-1 | rejects oversized files | test: `go test ./upload/` | 1 grounded green, 1 red | pass | T-1 | — |",
		"| AC-2 | accepts valid files | — | none | — | — | 1 (slip) |",
		"Unlinked ACs",
		"AC-9",
		"requirements: 1 · ACs: 2 · with grounded green evidence: 1 · verdict pass: 1",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("RTM missing %q:\n%s", want, out)
		}
	}
}

func TestRTMEmptyRun(t *testing.T) {
	c := newCtl(t)
	if _, err := c.RunStart("diff", "fix"); err != nil {
		t.Fatal(err)
	}
	out, err := RTM(c)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "no requirement records") {
		t.Errorf("empty run must say so:\n%s", out)
	}
}
