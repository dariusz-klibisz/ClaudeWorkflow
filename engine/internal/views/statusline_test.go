package views

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestStatuslineStates(t *testing.T) {
	c := newCtl(t)

	// no run
	if got := Statusline(c); got != "wf: no run" {
		t.Fatalf("no run: %q", got)
	}

	// active run in frame: short id, phase position, unmet count
	r, err := c.RunStart("diff", "fix")
	if err != nil {
		t.Fatal(err)
	}
	got := Statusline(c)
	shortID := r.ID[strings.LastIndex(r.ID, "-")+1:]
	for _, want := range []string{"wf " + shortID, "frame (1/7)", "unmet"} {
		if !strings.Contains(got, want) {
			t.Fatalf("active line missing %q: %q", want, got)
		}
	}
	if strings.Contains(got, "\n") {
		t.Fatalf("statusline must be one line: %q", got)
	}

	// parked
	r.Status = "parked"
	_ = c.Store.SaveRun(r)
	if got := Statusline(c); !strings.Contains(got, "parked") {
		t.Fatalf("parked line: %q", got)
	}

	// all phases done
	r.Status = "active"
	r.Phase = ""
	_ = c.Store.SaveRun(r)
	if got := Statusline(c); !strings.Contains(got, "wf run close") {
		t.Fatalf("done line: %q", got)
	}
}

// Broken state must degrade to silence, never an error line.
func TestStatuslineNeverLoud(t *testing.T) {
	c := newCtl(t)
	_, _ = c.RunStart("diff", "fix")
	snap := filepath.Join(c.Store.Root, "state", "run.json")
	if err := os.WriteFile(snap, []byte("{corrupt"), 0o644); err != nil {
		t.Fatal(err)
	}
	if got := Statusline(c); got != "" {
		t.Fatalf("corrupt snapshot must yield empty output: %q", got)
	}
}
