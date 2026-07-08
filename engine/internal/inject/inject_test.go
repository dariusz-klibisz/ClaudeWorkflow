package inject

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

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

// The designer (author-side, corpus-routed) must receive a briefing: scope,
// stage assignment, and corpus paths — with no verdict-block contract.
func TestAgentBriefsNonGatingDesigner(t *testing.T) {
	c := newCtl(t)
	if _, err := c.RunStart("diff", "new"); err != nil {
		t.Fatal(err)
	}
	// designer only acts in design; the briefing itself is phase-agnostic
	out, err := Agent(c, "designer")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "work scope") {
		t.Errorf("designer briefing missing work scope line:\n%s", out)
	}
	if !strings.Contains(out, "assigned design stage for this spawn: system") {
		t.Errorf("first designer spawn must be assigned the system stage:\n%s", out)
	}
	if !strings.Contains(out, "reference/design") {
		t.Errorf("designer briefing must route the design corpus:\n%s", out)
	}
	if strings.Contains(out, "```verdict") {
		t.Errorf("author-side briefing must not carry the verdict contract:\n%s", out)
	}
}

// Stage derivation is deterministic from recorded option-sets: system first,
// then software, then explicit loop re-entry guidance.
func TestDesignerStageProgression(t *testing.T) {
	c := newCtl(t)
	if _, err := c.RunStart("diff", "new"); err != nil {
		t.Fatal(err)
	}
	rec := func(stage string) {
		t.Helper()
		_, err := c.Record("option-set", map[string]any{
			"stage": stage, "candidates": []any{"a", "b"},
			"selected": "a", "rejected": []any{map[string]any{"id": "b", "reason": "r"}},
		}, false, "agent")
		if err != nil {
			t.Fatal(err)
		}
	}

	out, _ := Agent(c, "designer")
	if !strings.Contains(out, "stage for this spawn: system") {
		t.Fatalf("expected system stage first:\n%s", out)
	}
	rec("system")
	out, _ = Agent(c, "designer")
	if !strings.Contains(out, "stage for this spawn: software") {
		t.Fatalf("after system option-set, expected software stage:\n%s", out)
	}
	rec("software")
	out, _ = Agent(c, "designer")
	if !strings.Contains(out, "loop re-entry") || !strings.Contains(out, "rejected option IDs") {
		t.Fatalf("with both stages recorded, expected loop re-entry guidance:\n%s", out)
	}
}

// ux-designer is always the ux stage.
func TestUXDesignerStage(t *testing.T) {
	c := newCtl(t)
	if _, err := c.RunStart("diff", "new"); err != nil {
		t.Fatal(err)
	}
	out, err := Agent(c, "ux-designer")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "assigned design stage for this spawn: ux") {
		t.Errorf("ux-designer must be assigned the ux stage:\n%s", out)
	}
}

// The implementer briefing scopes the spawn to the single active task:
// tid, DoD, AC text + verification command, design refs, and the recorded
// out-of-scope boundary — with no verdict-block contract.
func TestImplementerBriefing(t *testing.T) {
	c := newCtl(t)
	if _, err := c.RunStart("diff", "new"); err != nil {
		t.Fatal(err)
	}
	rec := func(kind string, data map[string]any) {
		t.Helper()
		if _, err := c.Record(kind, data, false, "agent"); err != nil {
			t.Fatalf("%s: %v", kind, err)
		}
	}

	// no task in flight → the briefing says so instead of guessing
	out, err := Agent(c, "implementer")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "no single active task") {
		t.Errorf("taskless spawn must be told to return:\n%s", out)
	}

	rec("requirement", map[string]any{"rid": "SWR-1", "level": "software", "text": "t", "status": "active",
		"acs": []any{map[string]any{"id": "AC-1", "text": "empty file yields error", "verifiable": true}}})
	rec("verification-strategy", map[string]any{"ac": "AC-1", "method": "unit", "command": "go test ./pkg/ -run TestEmpty"})
	rec("option-set", map[string]any{"stage": "system", "candidates": []any{"a", "b"},
		"selected": "a", "rejected": []any{map[string]any{"id": "b", "reason": "r"}}})
	rec("scope-boundary", map[string]any{"in_scope": []any{"pkg/"}, "out_of_scope": []any{"legacy/"}})
	rec("task", map[string]any{"tid": "T-1", "subject": "handle empty files", "status": "in_progress",
		"dod": []any{"red then green"}, "ac_links": []any{"AC-1"}})

	out, err = Agent(c, "implementer")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{
		"work scope",
		"assigned task for this spawn: T-1 — handle empty files",
		"red then green",                // DoD
		"AC-1: empty file yields error", // AC text
		"go test ./pkg/ -run TestEmpty", // exact verification invocation
		"approved design selections",    // conformance anchor
		"system:a",                      // the selected option
		"out of scope (do not touch",    // boundary
		"legacy/",                       //
		"reference/coding",              // corpus routing
	} {
		if !strings.Contains(out, want) {
			t.Errorf("implementer briefing missing %q:\n%s", want, out)
		}
	}
	if strings.Contains(out, "```verdict") {
		t.Errorf("author-side briefing must not carry the verdict contract:\n%s", out)
	}
}

// On a designer loop re-entry, the briefing surfaces the finding CONTENT of
// the latest failing design-stage verdicts — the rework must see WHAT
// failed, not just the counts.
func TestDesignerLoopReentryCarriesPriorFindings(t *testing.T) {
	c := newCtl(t)
	if _, err := c.RunStart("diff", "new"); err != nil {
		t.Fatal(err)
	}
	rec := func(kind string, data map[string]any) {
		t.Helper()
		if _, err := c.Record(kind, data, false, "agent"); err != nil {
			t.Fatalf("%s: %v", kind, err)
		}
	}
	for _, stage := range []string{"system", "software"} {
		rec("option-set", map[string]any{"stage": stage, "candidates": []any{"a", "b"},
			"selected": "a", "rejected": []any{map[string]any{"id": "b", "reason": "r"}}})
	}
	rec("verdict", map[string]any{"agent": "design-reviewer", "status": "changes-required",
		"criticals": 1, "majors": 0,
		"findings": []any{"[critical] pkg/api.go: unbounded queue — 01 §backpressure"}})

	out, err := Agent(c, "designer")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "loop re-entry") {
		t.Fatalf("expected loop re-entry stage:\n%s", out)
	}
	if !strings.Contains(out, "unbounded queue") {
		t.Errorf("re-entry briefing must carry the prior failing findings:\n%s", out)
	}

	// a later clean verdict supersedes the failure — no stale findings
	rec("verdict", map[string]any{"agent": "design-reviewer", "status": "clean", "criticals": 0, "majors": 0})
	out, _ = Agent(c, "designer")
	if strings.Contains(out, "unbounded queue") {
		t.Errorf("superseded failure must not haunt the briefing:\n%s", out)
	}
}

// Gating reviewers keep the full contract: review scope, mode when derived,
// corpus, and the verdict block.
func TestAgentBriefsGatingReviewer(t *testing.T) {
	c := newCtl(t)
	if _, err := c.RunStart("diff", "new"); err != nil {
		t.Fatal(err)
	}
	out, err := Agent(c, "adversary")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"review scope", "abuse-case", "reference/coding/04-security.md", "```verdict"} {
		if !strings.Contains(out, want) {
			t.Errorf("adversary briefing missing %q:\n%s", want, out)
		}
	}
	// briefing paths are /-separated on every platform (vacuous on unix,
	// the real regression guard on Windows CI where filepath.Join used to
	// leak backslashes into the corpus routes)
	if strings.Contains(out, `\`) {
		t.Errorf("briefing must not carry backslash paths:\n%s", out)
	}
}

// Agents that are neither gating nor corpus-routed stay silent.
func TestAgentSkipsUnroutedAuthors(t *testing.T) {
	c := newCtl(t)
	if _, err := c.RunStart("diff", "new"); err != nil {
		t.Fatal(err)
	}
	out, err := Agent(c, "researcher")
	if err != nil {
		t.Fatal(err)
	}
	if out != "" {
		t.Errorf("researcher (no corpus, non-gating) must get no briefing, got:\n%s", out)
	}
	out, _ = Agent(c, "no-such-agent")
	if out != "" {
		t.Errorf("unknown agent must get no briefing, got:\n%s", out)
	}
}

// compliance-reviewer briefing: idle without packs; with pack items, names
// the standards in force and the project-side checklist paths.
func TestComplianceBriefing(t *testing.T) {
	c := newCtl(t)
	if _, err := c.RunStart("diff", "new"); err != nil {
		t.Fatal(err)
	}
	out, err := Agent(c, "compliance-reviewer")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "no regulated standards are in force") {
		t.Errorf("packless briefing must say the reviewer is idle:\n%s", out)
	}

	// simulate an installed pack: contracts.d item + pack doc
	itemYAML := "contracts:\n  - id: local.iso-26262.design-review\n    phase: design\n    predicate: verdict-in\n    params: { agent: compliance-reviewer, scope: iso-26262, statuses: [clean, n/a] }\n    remediation: \"spawn it\"\n"
	if err := os.WriteFile(filepath.Join(c.Store.ContractsDir(), "iso.yaml"), []byte(itemYAML), 0o644); err != nil {
		t.Fatal(err)
	}
	docDir := filepath.Join(c.Store.Root, "packs", "iso-26262")
	if err := os.MkdirAll(docDir, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(docDir, "iso-26262.md"), []byte("# checklist"), 0o644); err != nil {
		t.Fatal(err)
	}
	p, _ := filepath.Abs(filepath.Join("..", "..", "..", "workflow", "workflow.yaml"))
	sp, err := spec.Load(p, c.Store.ContractsDir())
	if err != nil {
		t.Fatal(err)
	}
	c.Spec = sp

	out, err = Agent(c, "compliance-reviewer")
	if err != nil {
		t.Fatal(err)
	}
	for _, want := range []string{"standards in force for this project: iso-26262", "iso-26262.md", "```verdict"} {
		if !strings.Contains(out, want) {
			t.Errorf("compliance briefing missing %q:\n%s", want, out)
		}
	}
}
