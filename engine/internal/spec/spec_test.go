package spec

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// repoSpec loads the real shipped workflow.yaml.
func repoSpec(t *testing.T, contractsDir string) *Spec {
	t.Helper()
	s, err := Load(shippedSpecPath(t), contractsDir)
	if err != nil {
		t.Fatalf("load shipped spec: %v", err)
	}
	return s
}

func shippedSpecPath(t *testing.T) string {
	t.Helper()
	// engine/internal/spec -> repo root
	p, err := filepath.Abs(filepath.Join("..", "..", "..", "workflow", "workflow.yaml"))
	if err != nil {
		t.Fatal(err)
	}
	return p
}

func TestShippedSpecValid(t *testing.T) {
	s := repoSpec(t, "")
	if len(s.Phases) != 7 {
		t.Errorf("want 7 phases, got %d", len(s.Phases))
	}
	if len(s.Families) != 3 {
		t.Errorf("want 3 families, got %d", len(s.Families))
	}
	if got := len(s.GatingAgents()); got != 11 {
		t.Errorf("want 11 gating agents, got %d", got)
	}
	if len(s.Roster) != 15 {
		t.Errorf("want 15 agents, got %d", len(s.Roster))
	}
}

func TestPhaseSequenceByFamily(t *testing.T) {
	s := repoSpec(t, "")
	// assessment skips design
	for _, p := range s.PhasesFor("assessment") {
		if p.ID == "design" {
			t.Error("assessment must not include design")
		}
	}
	if got := s.NextPhase("assessment", "context"); got != "plan" {
		t.Errorf("assessment after context: want plan, got %q", got)
	}
	if got := s.NextPhase("diff", "context"); got != "design" {
		t.Errorf("diff after context: want design, got %q", got)
	}
	if got := s.NextPhase("diff", "ship"); got != "" {
		t.Errorf("after ship: want end, got %q", got)
	}
}

func TestContractsForFamilyFiltering(t *testing.T) {
	s := repoSpec(t, "")
	// deps is diff-only (the B2 fix)
	for _, c := range s.ContractsFor("plan", "artifact") {
		if c.ID == "plan.deps" {
			t.Error("plan.deps must not apply to artifact family")
		}
	}
	found := false
	for _, c := range s.ContractsFor("plan", "diff") {
		if c.ID == "plan.deps" {
			found = true
		}
	}
	if !found {
		t.Error("plan.deps missing for diff family")
	}
}

func writeSpecDir(t *testing.T, contractsD map[string]string) (specPath, cdir string) {
	t.Helper()
	dir := t.TempDir()
	raw, err := os.ReadFile(shippedSpecPath(t))
	if err != nil {
		t.Fatal(err)
	}
	specPath = filepath.Join(dir, "workflow.yaml")
	if err := os.WriteFile(specPath, raw, 0o644); err != nil {
		t.Fatal(err)
	}
	cdir = filepath.Join(dir, "contracts.d")
	if err := os.MkdirAll(cdir, 0o755); err != nil {
		t.Fatal(err)
	}
	for name, content := range contractsD {
		if err := os.WriteFile(filepath.Join(cdir, name), []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return specPath, cdir
}

func TestContractsDAddOnly(t *testing.T) {
	good := `
contracts:
  - id: local.changelog
    phase: ship
    predicate: record-exists
    params: { kind: artifact, filter: { role: changelog } }
    remediation: "add a changelog entry"
`
	specPath, cdir := writeSpecDir(t, map[string]string{"10-local.yaml": good})
	s, err := Load(specPath, cdir)
	if err != nil {
		t.Fatalf("valid extension rejected: %v", err)
	}
	found := false
	for _, c := range s.ContractsFor("ship", "diff") {
		if c.ID == "local.changelog" {
			found = true
			if c.Source != "10-local.yaml" {
				t.Errorf("source not tracked: %q", c.Source)
			}
		}
	}
	if !found {
		t.Error("extension contract not merged")
	}
}

func TestContractsDRejectsOverride(t *testing.T) {
	override := `
contracts:
  - id: local.x
    phase: ship
    predicate: record-exists
    params: { kind: artifact }
    remediation: "r"
  - id: ship.audited
    phase: ship
    predicate: record-exists
    params: { kind: artifact }
    remediation: "weakened!"
`
	specPath, cdir := writeSpecDir(t, map[string]string{"a.yaml": override})
	_, err := Load(specPath, cdir)
	if err == nil {
		t.Fatal("shipped-id override must be rejected")
	}
	// namespacing enforced
	bad := `
contracts:
  - id: mystuff.x
    phase: ship
    predicate: record-exists
    params: { kind: artifact }
    remediation: "r"
`
	specPath, cdir = writeSpecDir(t, map[string]string{"b.yaml": bad})
	if _, err := Load(specPath, cdir); err == nil {
		t.Fatal("non-namespaced project contract must be rejected")
	}
}

func TestContractsDCustomRecordKinds(t *testing.T) {
	ext := `
records:
  - { kind: x-signoff, fields: [approver, "role?"] }
contracts:
  - id: local.signoff
    phase: ship
    predicate: record-exists
    params: { kind: x-signoff }
    remediation: "record the signoff"
`
	specPath, cdir := writeSpecDir(t, map[string]string{"c.yaml": ext})
	s, err := Load(specPath, cdir)
	if err != nil {
		t.Fatalf("custom kind rejected: %v", err)
	}
	rk, ok := s.RecordKind("x-signoff")
	if !ok {
		t.Fatal("x-signoff not merged")
	}
	if req := rk.Required(); len(req) != 1 || req[0] != "approver" {
		t.Errorf("required fields wrong: %v", req)
	}
	// non-x- kind rejected
	bad := `
records:
  - { kind: signoff, fields: [approver] }
`
	specPath, cdir = writeSpecDir(t, map[string]string{"d.yaml": bad})
	if _, err := Load(specPath, cdir); err == nil {
		t.Fatal("non-x- record kind must be rejected")
	}
}

func TestValidationCatchesDanglingRefs(t *testing.T) {
	raw, err := os.ReadFile(shippedSpecPath(t))
	if err != nil {
		t.Fatal(err)
	}
	// Introduce a contract referencing an undeclared agent.
	broken := string(raw) + `
  - id: frame.phantom
    phase: frame
    predicate: verdict-in
    params: { agent: adversarial-reviewer, statuses: [clean] }
    remediation: "the A1 phantom"
`
	dir := t.TempDir()
	p := filepath.Join(dir, "workflow.yaml")
	if err := os.WriteFile(p, []byte(broken), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err = Load(p, "")
	if err == nil {
		t.Fatal("phantom agent must fail validation (the A1 class)")
	}
	if !strings.Contains(err.Error(), "adversarial-reviewer") {
		t.Errorf("error should name the phantom: %v", err)
	}
}
