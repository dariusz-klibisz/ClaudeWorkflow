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

// Pack docs (.md) travel to .workflow/packs/<pack>/ — the injected
// checklist path for pack-referencing reviewers.
func TestPackInstallCopiesDocs(t *testing.T) {
	c, _ := newCtl(t)
	src := t.TempDir()
	packDir := filepath.Join(src, "my-std")
	if err := os.MkdirAll(packDir, 0o755); err != nil {
		t.Fatal(err)
	}
	writePack(t, packDir, "contracts.yaml", goodPack)
	if err := os.WriteFile(filepath.Join(packDir, "my-std.md"), []byte("# checklist\n- item"), 0o644); err != nil {
		t.Fatal(err)
	}

	out, err := PackInstall(c, specPathAbs(t), packDir)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(out, "1 doc(s)") {
		t.Errorf("install output missing doc note: %s", out)
	}
	doc := filepath.Join(c.Store.Root, "packs", "my-std", "my-std.md")
	if _, err := os.Stat(doc); err != nil {
		t.Fatalf("pack doc not copied to %s: %v", doc, err)
	}
}

// The shipped regulated packs must install cleanly against the shipped spec
// (the same check the exemplar smoke ran, kept green in CI).
func TestShippedPacksInstall(t *testing.T) {
	root, _ := filepath.Abs(filepath.Join("..", "..", ".."))
	packs := []string{
		filepath.Join(root, "packs", "sbom"),
		filepath.Join(root, "packs", "regulated", "iso-26262"),
		filepath.Join(root, "packs", "regulated", "iec-62304"),
		filepath.Join(root, "packs", "regulated", "do-178c"),
		filepath.Join(root, "packs", "regulated", "iec-61508"),
		filepath.Join(root, "packs", "regulated", "en-50128"),
		filepath.Join(root, "packs", "regulated", "nist-800-53"),
	}
	c, _ := newCtl(t)
	for _, p := range packs {
		if _, err := PackInstall(c, specPathAbs(t), p); err != nil {
			t.Errorf("shipped pack %s must install: %v", filepath.Base(p), err)
		}
	}
	// all six standards in force, distinct scopes
	sp, err := spec.Load(specPathAbs(t), c.Store.ContractsDir())
	if err != nil {
		t.Fatal(err)
	}
	if stds := sp.ComplianceStandards(); len(stds) != 6 {
		t.Errorf("want 6 standards in force, got %v", stds)
	}
}
