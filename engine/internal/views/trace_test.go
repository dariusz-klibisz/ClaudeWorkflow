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
	// 4 findings: force, followup, ambiguity + the vacated-signals finding
	// (frame was forced without a risk scan — the signal-conditioned
	// security items silently never applied; that must be visible)
	if len(open) != 4 {
		t.Fatalf("want 4 open findings (force, followup, ambiguity, vacated-signals), got %d\n%s", len(open), report)
	}
	if !strings.Contains(report, "force") || !strings.Contains(report, "followup") {
		t.Errorf("report must name the finding classes:\n%s", report)
	}
	if !strings.Contains(report, "without a risk scan") {
		t.Errorf("forced frame without risk record must surface the vacated security items:\n%s", report)
	}
	if !strings.Contains(report, "WAIVED") && !strings.Contains(report, "exited") {
		t.Errorf("report must show phase coverage:\n%s", report)
	}

	// idempotent: a second run creates nothing new
	_, err = Trace(c)
	if err != nil {
		t.Fatal(err)
	}
	if got := len(openFindings(t, c)); got != 4 {
		t.Fatalf("trace must be idempotent: want 4, got %d", got)
	}

	// disposition one; re-run keeps it closed and doesn't recreate
	first := open[0]
	_, err = c.Record("trace-finding", map[string]any{"updates": first.ID, "status": "resolved", "note": "converted to task"}, false, "agent")
	if err != nil {
		t.Fatal(err)
	}
	_, _ = Trace(c)
	if got := len(openFindings(t, c)); got != 3 {
		t.Fatalf("resolved finding must stay resolved after re-trace: want 3 open, got %d", got)
	}
}

// Approval drift: records added after an approval bound its refs surface as
// medium findings; re-approving binds them and creates nothing new.
func TestTraceApprovalDrift(t *testing.T) {
	c := newCtl(t)
	_, _ = c.RunStart("diff", "fix")
	rec := func(kind string, data map[string]any) {
		t.Helper()
		if _, err := c.Record(kind, data, false, "agent"); err != nil {
			t.Fatalf("%s: %v", kind, err)
		}
	}
	rec("requirement", map[string]any{"rid": "SWR-1", "level": "software", "text": "t", "status": "active",
		"acs": []any{map[string]any{"id": "AC-1", "text": "a", "verifiable": true}}})
	if _, err := c.Approve("scope", "baseline"); err != nil {
		t.Fatal(err)
	}
	// no drift yet
	if _, err := Trace(c); err != nil {
		t.Fatal(err)
	}
	for _, f := range openFindings(t, c) {
		if strings.HasPrefix(f.Data["key"].(string), "drift:") {
			t.Fatalf("no drift expected right after approval: %+v", f.Data)
		}
	}
	// a requirement appears after the approval → drift finding
	rec("requirement", map[string]any{"rid": "SWR-2", "level": "software", "text": "late addition", "status": "active",
		"acs": []any{map[string]any{"id": "AC-2", "text": "b", "verifiable": true}}})
	if _, err := Trace(c); err != nil {
		t.Fatal(err)
	}
	drift := 0
	for _, f := range openFindings(t, c) {
		if key, _ := f.Data["key"].(string); strings.HasPrefix(key, "drift:scope:req:SWR-2") {
			drift++
		}
	}
	if drift != 1 {
		t.Fatalf("want 1 drift finding for SWR-2, got %d", drift)
	}
	// re-approval binds the new baseline: no NEW drift finding on re-trace
	if _, err := c.Approve("scope", "baseline v2"); err != nil {
		t.Fatal(err)
	}
	before := len(openFindings(t, c))
	if _, err := Trace(c); err != nil {
		t.Fatal(err)
	}
	if got := len(openFindings(t, c)); got != before {
		t.Fatalf("re-approved baseline must create no new findings: %d -> %d", before, got)
	}
}

// Path-like out_of_scope entries are mechanically checked against edit
// records; prose entries are skipped (the auditor's to judge).
func TestTraceScopeBoundaryEdits(t *testing.T) {
	c := newCtl(t)
	_, _ = c.RunStart("diff", "new")
	rec := func(kind string, data map[string]any) {
		t.Helper()
		if _, err := c.Record(kind, data, false, "agent"); err != nil {
			t.Fatalf("%s: %v", kind, err)
		}
	}
	rec("scope-boundary", map[string]any{
		"in_scope":     []any{"pkg/"},
		"out_of_scope": []any{"legacy/", "*.sql", "authentication rework"},
	})
	rec("edit", map[string]any{"path": "pkg/api.go"})          // in scope
	rec("edit", map[string]any{"path": "legacy/db.go"})        // prefix hit
	rec("edit", map[string]any{"path": "migrations/init.sql"}) // glob-basename hit
	if _, err := Trace(c); err != nil {
		t.Fatal(err)
	}
	var hits []string
	for _, f := range openFindings(t, c) {
		if key, _ := f.Data["key"].(string); strings.HasPrefix(key, "scope:") {
			hits = append(hits, f.Data["text"].(string))
		}
	}
	if len(hits) != 2 {
		t.Fatalf("want 2 out-of-scope edit findings, got %d: %v", len(hits), hits)
	}
	joined := strings.Join(hits, "\n")
	for _, want := range []string{"legacy/db.go", "migrations/init.sql"} {
		if !strings.Contains(joined, want) {
			t.Errorf("missing out-of-scope hit for %s:\n%s", want, joined)
		}
	}
	if strings.Contains(joined, "pkg/api.go") {
		t.Errorf("in-scope edit must not be flagged:\n%s", joined)
	}
	// prose entry produced no finding — nothing matches "authentication rework"
}

// Edits recorded after the latest gating verdict surface as a medium
// staleness finding; a fresher verdict clears the condition.
func TestTraceVerdictStaleness(t *testing.T) {
	c := newCtl(t)
	_, _ = c.RunStart("diff", "new")
	rec := func(kind string, data map[string]any) {
		t.Helper()
		if _, err := c.Record(kind, data, false, "agent"); err != nil {
			t.Fatalf("%s: %v", kind, err)
		}
	}
	stale := func() bool {
		t.Helper()
		if _, err := Trace(c); err != nil {
			t.Fatal(err)
		}
		for _, f := range openFindings(t, c) {
			if key, _ := f.Data["key"].(string); key == "stale-verdicts" {
				return true
			}
		}
		return false
	}
	rec("edit", map[string]any{"path": "pkg/a.go"})
	rec("verdict", map[string]any{"agent": "code-quality-reviewer", "status": "clean", "criticals": 0, "majors": 0})
	if stale() {
		t.Fatal("verdict newer than all edits must not be stale")
	}
	rec("edit", map[string]any{"path": "pkg/b.go"})
	if !stale() {
		t.Fatal("edit after the latest gating verdict must surface staleness")
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
