package contracts

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func writeDoc(t *testing.T, dir, rel, content string) {
	t.Helper()
	abs := filepath.Join(dir, filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

const authored = `# ADR 1: choose the boring thing

- Status: accepted

## Context

We must decide between three storage engines with different operational
profiles and team familiarity. The workload is read-heavy with bursty writes.

## Decision

We choose the boring, well-understood option because operational familiarity
outweighs the marginal performance headroom of the alternatives.

## Consequences

Migration tooling must be written; the rejected options are recorded below.
`

func TestArtifactPresentRequiresRecordAndDisk(t *testing.T) {
	dir := t.TempDir()
	b := &envBuilder{}
	env := newEnv(t, b, "artifact", "arch-design", nil)
	env.ProjectDir = dir

	// no record at all
	ok, detail, err := EvalOne(env, "artifact-present", map[string]any{"template": "adr"}, "")
	if err != nil || ok {
		t.Fatalf("no record must fail: ok=%v err=%v", ok, err)
	}
	if !strings.Contains(detail, "no artifact record") {
		t.Errorf("detail should say no record: %s", detail)
	}

	// record exists, file missing
	b.add("artifact", true, map[string]any{"path": "docs/architecture/adr/0001-x.md", "status": "stub", "template": "adr"})
	env = newEnv(t, b, "artifact", "arch-design", nil)
	env.ProjectDir = dir
	ok, detail, _ = EvalOne(env, "artifact-present", map[string]any{"template": "adr"}, "")
	if ok {
		t.Fatal("missing file must fail")
	}
	if !strings.Contains(detail, "does not exist on disk") {
		t.Errorf("detail should say missing on disk: %s", detail)
	}

	// stub content fails
	writeDoc(t, dir, "docs/architecture/adr/0001-x.md", "# ADR\n\nTODO\n")
	ok, detail, _ = EvalOne(env, "artifact-present", map[string]any{"template": "adr"}, "")
	if ok {
		t.Fatal("stub content must fail")
	}
	if !strings.Contains(detail, "stub") {
		t.Errorf("detail should say stub: %s", detail)
	}

	// authored content passes
	writeDoc(t, dir, "docs/architecture/adr/0001-x.md", authored)
	ok, _, err = EvalOne(env, "artifact-present", map[string]any{"template": "adr"}, "")
	if err != nil || !ok {
		t.Fatalf("authored artifact must pass: ok=%v err=%v", ok, err)
	}
}

func TestArtifactPresentTemplateIdenticalIsStub(t *testing.T) {
	dir := t.TempDir()
	b := &envBuilder{}
	// the real plugin template, copied verbatim (what wf doc new does)
	sp := loadSpec(t)
	tmpl, err := os.ReadFile(filepath.Join(sp.PluginRoot(), "templates", "adr.md"))
	if err != nil {
		t.Skip("plugin templates not present in this checkout")
	}
	writeDoc(t, dir, "docs/architecture/adr/0001-y.md", string(tmpl))
	b.add("artifact", true, map[string]any{"path": "docs/architecture/adr/0001-y.md", "status": "stub", "template": "adr"})
	env := newEnv(t, b, "artifact", "arch-design", nil)
	env.ProjectDir = dir
	ok, detail, _ := EvalOne(env, "artifact-present", map[string]any{"template": "adr"}, "")
	if ok {
		t.Fatal("untouched template copy must fail")
	}
	if !strings.Contains(detail, "byte-identical") {
		t.Errorf("detail should say byte-identical: %s", detail)
	}
}

func TestArtifactPresentAbandonedIsSkipped(t *testing.T) {
	dir := t.TempDir()
	b := &envBuilder{}
	b.add("artifact", true, map[string]any{"path": "docs/x.md", "status": "missing", "template": "design"})
	env := newEnv(t, b, "artifact", "arch-design", nil)
	env.ProjectDir = dir
	ok, detail, _ := EvalOne(env, "artifact-present", map[string]any{"template": "design"}, "")
	if ok {
		t.Fatal("an abandoned-only set must fail (no present artifact)")
	}
	if !strings.Contains(detail, "abandoned") {
		t.Errorf("detail should say abandoned: %s", detail)
	}
}

func TestArtifactPresentNoFilterMatchesAny(t *testing.T) {
	dir := t.TempDir()
	writeDoc(t, dir, "docs/reviews/r.md", authored)
	b := &envBuilder{}
	b.add("artifact", false, map[string]any{"path": "docs/reviews/r.md", "status": "present"})
	env := newEnv(t, b, "assessment", "code-review", nil)
	env.ProjectDir = dir
	ok, _, err := EvalOne(env, "artifact-present", map[string]any{}, "")
	if err != nil || !ok {
		t.Fatalf("filterless artifact-present must accept any authored artifact: ok=%v err=%v", ok, err)
	}
}

func TestArtifactPresentEmptyProjectDirVacuous(t *testing.T) {
	b := &envBuilder{}
	b.add("artifact", false, map[string]any{"path": "docs/x.md", "status": "present"})
	env := newEnv(t, b, "assessment", "code-review", nil)
	ok, _, err := EvalOne(env, "artifact-present", map[string]any{}, "")
	if err != nil || !ok {
		t.Fatalf("no ProjectDir (views) must pass vacuously: ok=%v err=%v", ok, err)
	}
}

func TestEntryContractsEvaluate(t *testing.T) {
	b := &envBuilder{}
	env := newEnv(t, b, "diff", "fix", nil)
	fs, err := EvaluateEntry(env, "context")
	if err != nil {
		t.Fatal(err)
	}
	ids := findingIDs(fs)
	if _, ok := ids["context.entry-classification"]; !ok {
		t.Fatalf("context entry must demand the classification: %v", fs)
	}
	// entry items never leak into the exit evaluation
	exitFs, _ := EvaluatePhase(env, "context")
	if _, ok := findingIDs(exitFs)["context.entry-classification"]; ok {
		t.Fatal("entry item leaked into exit contract")
	}
	// satisfied input clears
	b.add("classification", false, map[string]any{"family": "diff", "intent": "fix", "restated": "r"})
	env = newEnv(t, b, "diff", "fix", nil)
	fs, _ = EvaluateEntry(env, "context")
	if _, ok := findingIDs(fs)["context.entry-classification"]; ok {
		t.Fatal("satisfied entry item must clear")
	}
}

func TestEntryContractsWaivable(t *testing.T) {
	b := &envBuilder{}
	b.add("waiver", false, map[string]any{"item": "context.entry-classification", "reason": "adopt landing"})
	env := newEnv(t, b, "diff", "fix", nil)
	fs, err := EvaluateEntry(env, "context")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := findingIDs(fs)["context.entry-classification"]; ok {
		t.Fatal("waived entry item must clear")
	}
}
