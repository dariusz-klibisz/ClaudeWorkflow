package ops

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/spec"
)

func specPathAbs(t *testing.T) string {
	t.Helper()
	p, _ := filepath.Abs(filepath.Join("..", "..", "..", "workflow", "workflow.yaml"))
	return p
}

const goodPack = `records:
  - { kind: x-sbom, fields: [path, status] }
contracts:
  - id: local.sbom-present
    phase: ship
    families: [diff]
    predicate: record-exists
    params: { kind: x-sbom }
    waivable: true
    remediation: "Generate the SBOM and record it: wf record x-sbom …"
`

const badPack = `contracts:
  - id: sbom-present
    phase: ship
    predicate: record-exists
    params: { kind: x-sbom }
    remediation: "…"
`

func writePack(t *testing.T, dir, name, content string) string {
	t.Helper()
	p := filepath.Join(dir, name)
	if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
	return p
}

func TestPackInstallDirectory(t *testing.T) {
	c, _ := newCtl(t)
	src := t.TempDir()
	packDir := filepath.Join(src, "regulated-baseline")
	writePack(t, packDir, "items.yaml", goodPack)
	writePack(t, packDir, "README.md", "# docs, ignored by the installer")

	out, err := PackInstall(c, specPathAbs(t), packDir)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "regulated-baseline") {
		t.Errorf("output should name the pack: %s", out)
	}
	installed := filepath.Join(c.Store.ContractsDir(), "regulated-baseline-items.yaml")
	if _, err := os.Stat(installed); err != nil {
		t.Fatalf("pack file not installed: %v", err)
	}
	// the merged spec must load and carry the pack's item
	sp, err := spec.LoadStrict(specPathAbs(t), c.Store.ContractsDir())
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, it := range sp.Contracts {
		if it.ID == "local.sbom-present" {
			found = true
		}
	}
	if !found {
		t.Fatal("pack contract item missing from merged spec")
	}
	// reinstall = collision = refused
	if _, err := PackInstall(c, specPathAbs(t), packDir); err == nil {
		t.Fatal("reinstall over existing files must be refused")
	}
}

func TestPackInstallRejectsBadNamespace(t *testing.T) {
	c, _ := newCtl(t)
	src := t.TempDir()
	p := writePack(t, src, "sloppy.yaml", badPack)
	if _, err := PackInstall(c, specPathAbs(t), p); err == nil {
		t.Fatal("un-namespaced pack must be rejected")
	}
	// nothing installed
	entries, _ := os.ReadDir(c.Store.ContractsDir())
	if len(entries) != 0 {
		t.Fatalf("rejected pack must leave contracts.d untouched: %v", entries)
	}
}
