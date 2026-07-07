package doctor

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

// fakeHome builds a Claude Code plugin layout: an installed wf plugin (with
// a stub bootstrap.sh that installs a fake engine) and an EMPTY data dir —
// the exact state a mid-session /plugin install leaves behind.
func fakeHome(t *testing.T, withBootstrap bool) (home, dataDir string) {
	t.Helper()
	home = t.TempDir()
	plugins := filepath.Join(home, ".claude", "plugins")
	root := filepath.Join(plugins, "cache", "claude-workflow", "wf", "0.2.0")
	dataDir = filepath.Join(plugins, "data", "wf-claude-workflow")
	for _, d := range []string{root, dataDir} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	reg := `{"version":2,"plugins":{"wf@claude-workflow":[{"scope":"project","installPath":` +
		jsonStr(root) + `,"version":"0.2.0"}],"other@m":[{"installPath":"/nope"}]}}`
	if err := os.WriteFile(filepath.Join(plugins, "installed_plugins.json"), []byte(reg), 0o644); err != nil {
		t.Fatal(err)
	}
	if withBootstrap {
		script := "#!/bin/sh\nmkdir -p \"$CLAUDE_PLUGIN_DATA/bin\"\nprintf fake > \"$CLAUDE_PLUGIN_DATA/bin/wf\"\nchmod +x \"$CLAUDE_PLUGIN_DATA/bin/wf\"\n"
		if err := os.MkdirAll(filepath.Join(root, "scripts"), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(root, "scripts", "bootstrap.sh"), []byte(script), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	return home, dataDir
}

func jsonStr(s string) string {
	return `"` + strings.ReplaceAll(s, `\`, `\\`) + `"`
}

// The power-of-ten incident: plugin installed, data dir empty, no hook env.
// Plain findings must see the dead engine via installed_plugins.json.
func TestHookEngineDetectsDeadInstall(t *testing.T) {
	home, _ := fakeHome(t, false)
	findings, dead := HookEngineFindings(home, false)
	if len(findings) != 1 || !dead {
		t.Fatalf("expected one dead-engine finding, got dead=%v %v", dead, findings)
	}
	if !strings.Contains(findings[0], "wf-claude-workflow") {
		t.Fatalf("finding must name the data dir: %v", findings[0])
	}
}

// Heal mode must run the install's bootstrap script and revive the hooks.
func TestHookEngineHeals(t *testing.T) {
	if _, err := exec.LookPath("sh"); err != nil {
		t.Skip("no sh on PATH")
	}
	home, dataDir := fakeHome(t, true)
	findings, dead := HookEngineFindings(home, true)
	if dead {
		t.Fatalf("expected healed, got dead: %v", findings)
	}
	if len(findings) != 1 || !strings.Contains(findings[0], "installed it now") {
		t.Fatalf("expected an installed-now note: %v", findings)
	}
	if _, err := os.Stat(filepath.Join(dataDir, "bin", "wf")); err != nil {
		t.Fatalf("engine not installed: %v", err)
	}
	// second pass: healthy, no findings
	findings, dead = HookEngineFindings(home, true)
	if len(findings) != 0 || dead {
		t.Fatalf("expected healthy after heal, got dead=%v %v", dead, findings)
	}
}

// Heal without a bootstrap script must stay dead and say so.
func TestHookEngineHealFailsWithoutScript(t *testing.T) {
	if _, err := exec.LookPath("sh"); err != nil {
		t.Skip("no sh on PATH")
	}
	home, _ := fakeHome(t, false)
	findings, dead := HookEngineFindings(home, true)
	if !dead || len(findings) != 1 {
		t.Fatalf("expected one still-dead finding, got dead=%v %v", dead, findings)
	}
}

// selfUpdateFixture: a plugin root expecting version 9.9.9 with a stub
// bootstrap that "installs" a new engine, and a data dir running "1.0.0".
func selfUpdateFixture(t *testing.T, withVersionFile bool) (root, data string) {
	t.Helper()
	base := t.TempDir()
	root = filepath.Join(base, "root")
	data = filepath.Join(base, "data")
	for _, d := range []string{filepath.Join(root, "scripts"), filepath.Join(root, "bin"), filepath.Join(data, "bin")} {
		if err := os.MkdirAll(d, 0o755); err != nil {
			t.Fatal(err)
		}
	}
	script := "#!/bin/sh\nprintf NEW-ENGINE > \"$CLAUDE_PLUGIN_DATA/bin/wf\"\nchmod +x \"$CLAUDE_PLUGIN_DATA/bin/wf\"\n"
	if err := os.WriteFile(filepath.Join(root, "scripts", "bootstrap.sh"), []byte(script), 0o755); err != nil {
		t.Fatal(err)
	}
	if withVersionFile {
		if err := os.WriteFile(filepath.Join(root, "bin", "VERSION"), []byte("9.9.9"), 0o644); err != nil {
			t.Fatal(err)
		}
	} else {
		manifest := "version 9.9.9\nbase_url file:///nowhere\n"
		if err := os.WriteFile(filepath.Join(root, "bin", "MANIFEST"), []byte(manifest), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if err := os.WriteFile(filepath.Join(data, "bin", "wf"), []byte("OLD-ENGINE"), 0o755); err != nil {
		t.Fatal(err)
	}
	t.Setenv("CLAUDE_PLUGIN_ROOT", root)
	t.Setenv("CLAUDE_PLUGIN_DATA", data)
	return root, data
}

// The M5 self-update path: version skew in hook context re-runs the
// bootstrap and reports it; the rate limit stops retry loops.
func TestSelfUpdateOnSkew(t *testing.T) {
	if _, err := exec.LookPath("sh"); err != nil {
		t.Skip("no sh on PATH")
	}
	_, data := selfUpdateFixture(t, true)

	note := SelfUpdate("1.0.0")
	if !strings.Contains(note, "self-updated") {
		t.Fatalf("skew must trigger an update: %q", note)
	}
	installed, _ := os.ReadFile(filepath.Join(data, "bin", "wf"))
	if string(installed) != "NEW-ENGINE" {
		t.Fatalf("engine not replaced: %q", installed)
	}
	// rate limit: immediate second attempt for the same expected version is silent
	if note := SelfUpdate("1.0.0"); note != "" {
		t.Fatalf("rate limit must silence the retry: %q", note)
	}
}

func TestSelfUpdateNoSkewNoAction(t *testing.T) {
	_, data := selfUpdateFixture(t, true)
	if note := SelfUpdate("9.9.9"); note != "" {
		t.Fatalf("matching version must be a no-op: %q", note)
	}
	installed, _ := os.ReadFile(filepath.Join(data, "bin", "wf"))
	if string(installed) != "OLD-ENGINE" {
		t.Fatal("no-op must not touch the engine")
	}
}

func TestSelfUpdateOutsideHookContext(t *testing.T) {
	t.Setenv("CLAUDE_PLUGIN_ROOT", "")
	t.Setenv("CLAUDE_PLUGIN_DATA", "")
	if note := SelfUpdate("1.0.0"); note != "" {
		t.Fatalf("no hook env must mean no action: %q", note)
	}
}

// Fetch-tier installs carry no bin/VERSION in the root — the committed
// MANIFEST's semver is the expected version then.
func TestExpectedVersionFromManifest(t *testing.T) {
	root, _ := selfUpdateFixture(t, false)
	if got := expectedVersion(root); got != "9.9.9" {
		t.Fatalf("expectedVersion = %q, want 9.9.9", got)
	}
}

// A healthy install (engine present) yields no findings.
func TestHookEngineHealthy(t *testing.T) {
	home, dataDir := fakeHome(t, false)
	if err := os.MkdirAll(filepath.Join(dataDir, "bin"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dataDir, "bin", "wf"), []byte("x"), 0o755); err != nil {
		t.Fatal(err)
	}
	findings, dead := HookEngineFindings(home, false)
	if len(findings) != 0 || dead {
		t.Fatalf("expected healthy, got dead=%v %v", dead, findings)
	}
}
