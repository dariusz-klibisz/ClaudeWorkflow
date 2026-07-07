package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/spec"
)

func writeAgent(t *testing.T, dir, name, frontmatter string) string {
	t.Helper()
	p := filepath.Join(dir, name+".md")
	if err := os.WriteFile(p, []byte("---\nname: "+name+"\n"+frontmatter+"---\n\n# "+name+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

// The M4 memory subset is deliberate — roster and frontmatter must agree.
func TestCheckAgentMemory(t *testing.T) {
	dir := t.TempDir()

	// agree: memory present in both
	p := writeAgent(t, dir, "adversary", "memory: project\n")
	if err := checkAgentMemory(p, spec.Agent{Name: "adversary", Memory: "project"}); err != nil {
		t.Fatalf("matching memory must pass: %v", err)
	}
	// agree: absent in both
	p = writeAgent(t, dir, "critic", "maxTurns: 30\n")
	if err := checkAgentMemory(p, spec.Agent{Name: "critic"}); err != nil {
		t.Fatalf("both-absent must pass: %v", err)
	}
	// drift: frontmatter has it, roster doesn't
	p = writeAgent(t, dir, "auditor", "memory: project\n")
	if err := checkAgentMemory(p, spec.Agent{Name: "auditor"}); err == nil {
		t.Fatal("frontmatter-only memory must fail the check")
	}
	// drift: roster has it, frontmatter doesn't
	p = writeAgent(t, dir, "design-reviewer", "maxTurns: 30\n")
	if err := checkAgentMemory(p, spec.Agent{Name: "design-reviewer", Memory: "project"}); err == nil {
		t.Fatal("roster-only memory must fail the check")
	}
}

func TestFrontmatterField(t *testing.T) {
	content := "---\nname: x\nmemory: project\n---\n\nmemory: not-this\n"
	if got := frontmatterField(content, "memory"); got != "project" {
		t.Fatalf("got %q", got)
	}
	if got := frontmatterField("no frontmatter\nmemory: x\n", "memory"); got != "" {
		t.Fatalf("body must not count: %q", got)
	}
}
