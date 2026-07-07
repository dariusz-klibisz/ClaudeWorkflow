// Package scripts tests the repo's shell bootstrap (scripts/bootstrap.sh)
// end-to-end — the fetch tier especially, using file:// URLs so no network
// is touched. The sh script is the first thing a fresh install runs; it
// deserves real tests, not hope.
package scripts

import (
	"crypto/sha256"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

func needTools(t *testing.T) {
	t.Helper()
	for _, tool := range []string{"sh", "curl", "sha256sum"} {
		if _, err := exec.LookPath(tool); err != nil {
			t.Skipf("no %s on PATH", tool)
		}
	}
}

// fixture builds a fake plugin root with a MANIFEST pointing at a local
// file:// "release", plus the real bootstrap.sh.
func fixture(t *testing.T, binaryContent string, tamper bool) (root, data string) {
	t.Helper()
	base := t.TempDir()
	root = filepath.Join(base, "root")
	data = filepath.Join(base, "data")
	name := fmt.Sprintf("wf-%s-%s", runtime.GOOS, runtime.GOARCH)

	// the "release" asset
	rel := filepath.Join(base, "releases", "v9.9.9")
	if err := os.MkdirAll(rel, 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(rel, name), []byte(binaryContent), 0o755); err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256([]byte(binaryContent))
	if tamper {
		sum = sha256.Sum256([]byte("something else entirely"))
	}

	// plugin root: scripts/bootstrap.sh (the real one) + bin/MANIFEST
	realScript, err := os.ReadFile(filepath.Join("..", "..", "..", "scripts", "bootstrap.sh"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.MkdirAll(filepath.Join(root, "scripts"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "scripts", "bootstrap.sh"), realScript, 0o755); err != nil {
		t.Fatal(err)
	}
	manifest := fmt.Sprintf("version 9.9.9\nbase_url file://%s\n%x  %s\n",
		filepath.ToSlash(filepath.Join(base, "releases")), sum, name)
	if err := os.MkdirAll(filepath.Join(root, "bin"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(root, "bin", "MANIFEST"), []byte(manifest), 0o644); err != nil {
		t.Fatal(err)
	}
	return root, data
}

func runBootstrap(t *testing.T, root, data string) string {
	t.Helper()
	cmd := exec.Command("sh", filepath.Join(root, "scripts", "bootstrap.sh"))
	cmd.Env = append(os.Environ(),
		"CLAUDE_PLUGIN_ROOT="+root,
		"CLAUDE_PLUGIN_DATA="+data,
		"PATH="+os.Getenv("PATH"),
	)
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("bootstrap exited non-zero (must always exit 0): %v\n%s", err, out)
	}
	return string(out)
}

func TestFetchTierInstallsVerified(t *testing.T) {
	needTools(t)
	root, data := fixture(t, "FAKE-ENGINE-v9.9.9", false)

	out := runBootstrap(t, root, data)
	if !strings.Contains(out, "fetched") || !strings.Contains(out, "checksum verified") {
		t.Fatalf("fetch tier did not run:\n%s", out)
	}
	installed, err := os.ReadFile(filepath.Join(data, "bin", "wf"))
	if err != nil || string(installed) != "FAKE-ENGINE-v9.9.9" {
		t.Fatalf("engine not installed: %v %q", err, installed)
	}
	// fetched binary persisted into plugin bin/ → bundled tier next time
	name := fmt.Sprintf("wf-%s-%s", runtime.GOOS, runtime.GOARCH)
	if _, err := os.Stat(filepath.Join(root, "bin", name)); err != nil {
		t.Fatalf("fetched binary not persisted to plugin bin/: %v", err)
	}
	// sha256 stamp written → second run is a silent no-op
	stamp, _ := os.ReadFile(filepath.Join(data, "bin", "VERSION"))
	if !strings.HasPrefix(string(stamp), "sha256:") {
		t.Fatalf("stamp: %q", stamp)
	}
	out = runBootstrap(t, root, data)
	if strings.Contains(out, "installed") || strings.Contains(out, "fetched") {
		t.Fatalf("second run must be a no-op:\n%s", out)
	}
}

func TestFetchTierRefusesChecksumMismatch(t *testing.T) {
	needTools(t)
	root, data := fixture(t, "TAMPERED-ENGINE", true)

	out := runBootstrap(t, root, data)
	if !strings.Contains(out, "checksum mismatch") || !strings.Contains(out, "refusing") {
		t.Fatalf("tampered download must be refused:\n%s", out)
	}
	if _, err := os.Stat(filepath.Join(data, "bin", "wf")); !os.IsNotExist(err) {
		t.Fatal("nothing may be installed on checksum mismatch")
	}
	// and nothing persisted into plugin bin/ either
	name := fmt.Sprintf("wf-%s-%s", runtime.GOOS, runtime.GOARCH)
	if _, err := os.Stat(filepath.Join(root, "bin", name)); !os.IsNotExist(err) {
		t.Fatal("tampered binary must not be persisted")
	}
}

func TestFetchTierNoManifestEntryFallsThrough(t *testing.T) {
	needTools(t)
	root, data := fixture(t, "x", false)
	// strip the platform line — only version/base_url remain
	mp := filepath.Join(root, "bin", "MANIFEST")
	raw, _ := os.ReadFile(mp)
	var kept []string
	for _, l := range strings.Split(string(raw), "\n") {
		if strings.HasPrefix(l, "version ") || strings.HasPrefix(l, "base_url ") {
			kept = append(kept, l)
		}
	}
	if err := os.WriteFile(mp, []byte(strings.Join(kept, "\n")+"\n"), 0o644); err != nil {
		t.Fatal(err)
	}

	out := runBootstrap(t, root, data)
	if !strings.Contains(out, "fail open") {
		t.Fatalf("no entry + no Go fallback must fail open loud:\n%s", out)
	}
	if _, err := os.Stat(filepath.Join(data, "bin", "wf")); !os.IsNotExist(err) {
		t.Fatal("nothing may be installed")
	}
}

func TestBundledTierStillPreferred(t *testing.T) {
	needTools(t)
	root, data := fixture(t, "FETCHED", false)
	// a bundled binary present → fetch must not run at all
	name := fmt.Sprintf("wf-%s-%s", runtime.GOOS, runtime.GOARCH)
	if err := os.WriteFile(filepath.Join(root, "bin", name), []byte("BUNDLED"), 0o755); err != nil {
		t.Fatal(err)
	}
	out := runBootstrap(t, root, data)
	if strings.Contains(out, "fetched") {
		t.Fatalf("bundled binary present — fetch must not run:\n%s", out)
	}
	installed, _ := os.ReadFile(filepath.Join(data, "bin", "wf"))
	if string(installed) != "BUNDLED" {
		t.Fatalf("bundled binary must win: %q", installed)
	}
}
