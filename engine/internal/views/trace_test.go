package views

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/contracts"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/runctl"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/spec"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/store"
)

func newCtl(t *testing.T) *runctl.Ctl {
	t.Helper()
	s, err := store.Open(t.TempDir(), true)
	if err != nil {
		t.Fatal(err)
	}
	p, _ := filepath.Abs(filepath.Join("..", "..", "..", "workflow", "workflow.yaml"))
	sp, err := spec.Load(p, "")
	if err != nil {
		t.Fatal(err)
	}
	return &runctl.Ctl{Store: s, Spec: sp, Config: &store.Config{}}
}

func openFindings(t *testing.T, c *runctl.Ctl) []contracts.Record {
	t.Helper()
	r, _ := c.Store.LoadRun()
	env, err := c.Env(r)
	if err != nil {
		t.Fatal(err)
	}
	var out []contracts.Record
	for _, tf := range env.Records("trace-finding") {
		if s, _ := tf.Data["status"].(string); s == "open" {
			out = append(out, tf)
		}
	}
	return out
}

func TestTraceFindingsAndIdempotence(t *testing.T) {
	c := newCtl(t)
	_, _ = c.RunStart("diff", "fix")
	// a force, an open followup, an unresolved ambiguity
	_, _, err := c.PhaseExit(true, "demo force")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = c.Record("followup", map[string]any{"text": "clean up config loader", "status": "open"}, false, "agent")
	_, _ = c.Record("ambiguity", map[string]any{"lens": "user", "text": "which locale?", "disposition": "deferred"}, false, "agent")

	report, err := Trace(c)
	if err != nil {
		t.Fatal(err)
	}
	open := openFindings(t, c)
	if len(open) != 3 {
		t.Fatalf("want 3 open findings (force, followup, ambiguity), got %d\n%s", len(open), report)
	}
	if !strings.Contains(report, "force") || !strings.Contains(report, "followup") {
		t.Errorf("report must name the finding classes:\n%s", report)
	}
	if !strings.Contains(report, "WAIVED") && !strings.Contains(report, "exited") {
		t.Errorf("report must show phase coverage:\n%s", report)
	}

	// idempotent: a second run creates nothing new
	_, err = Trace(c)
	if err != nil {
		t.Fatal(err)
	}
	if got := len(openFindings(t, c)); got != 3 {
		t.Fatalf("trace must be idempotent: want 3, got %d", got)
	}

	// disposition one; re-run keeps it closed and doesn't recreate
	first := open[0]
	_, err = c.Record("trace-finding", map[string]any{"updates": first.ID, "status": "resolved", "note": "converted to task"}, false, "agent")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = Trace(c)
	if got := len(openFindings(t, c)); got != 2 {
		t.Fatalf("resolved finding must stay resolved after re-trace: want 2 open, got %d", got)
	}
}

func TestTraceMakesShipGateReal(t *testing.T) {
	c := newCtl(t)
	_, _ = c.RunStart("diff", "fix")
	_, _ = c.Record("followup", map[string]any{"text": "leak", "status": "open"}, false, "agent")
	if _, err := Trace(c); err != nil {
		t.Fatal(err)
	}
	// ship.trace-clean must now block on the open finding
	r, _ := c.Store.LoadRun()
	r.Phase = "ship"
	_ = c.Store.SaveRun(r)
	env, _ := c.Env(r)
	findings, err := contracts.EvaluatePhase(env, "ship")
	if err != nil {
		t.Fatal(err)
	}
	hit := false
	for _, f := range findings {
		if f.ID == "ship.trace-clean" {
			hit = true
		}
	}
	if !hit {
		t.Fatal("open trace-finding must block ship exit (the gate is no longer vacuous)")
	}
}
