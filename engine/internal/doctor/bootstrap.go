package doctor

// Hook-engine reachability: every generated hook invokes
// ${CLAUDE_PLUGIN_DATA}/bin/wf, which only exists after the SessionStart
// bootstrap ran. A plugin installed or reloaded MID-session never fires
// SessionStart, so the data dir stays empty and every gate ENOENTs for the
// rest of the session (the power-of-ten incident: a full run shipped with
// zero hook events). These checks discover wf installs from the hook env
// and Claude Code's installed_plugins.json, report the dead state, and —
// in heal mode — run the plugin's own bootstrap script immediately, which
// revives the hooks without waiting for a restart (hooks resolve the path
// at fire time).

import (
	"crypto/sha256"
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

// hookInstall pairs a plugin install dir (root) with its data dir — the
// two env vars the bootstrap script needs.
type hookInstall struct {
	root string // may be empty when only the data path is known
	data string
}

// HookEngineFindings reports, for every discoverable wf plugin install,
// a missing hook engine at <data>/bin/wf. With heal=true it runs the
// install's scripts/bootstrap.sh (or .ps1 on native Windows) with the
// plugin env set, then re-checks. No findings = healthy or no installs
// discoverable; dead=true means at least one install still lacks the
// engine after any heal attempt. home is the user home dir
// ("" = os.UserHomeDir).
func HookEngineFindings(home string, heal bool) (findings []string, dead bool) {
	if home == "" {
		home, _ = os.UserHomeDir()
	}
	for _, in := range discoverInstalls(home) {
		enginePath := filepath.Join(in.data, "bin", "wf")
		if _, err := os.Stat(enginePath); err == nil {
			continue
		}
		msg := fmt.Sprintf("hook engine missing at %s — the SessionStart bootstrap has not run (mid-session plugin installs never fire it); every gate is dead", enginePath)
		healed := false
		if heal && in.root != "" {
			if out, err := runBootstrap(in.root, in.data); err != nil {
				msg += fmt.Sprintf("; bootstrap attempt failed: %v (%s)", err, strings.TrimSpace(out))
			} else if _, err := os.Stat(enginePath); err == nil {
				msg = fmt.Sprintf("hook engine was missing at %s — installed it now via %s/scripts/bootstrap.sh; hooks fire from the next event on (restart the session if hook errors persist)", enginePath, in.root)
				healed = true
			} else {
				msg += fmt.Sprintf("; bootstrap ran but installed nothing (%s)", strings.TrimSpace(out))
			}
		} else if heal {
			msg += "; plugin root unknown — reinstall the plugin or restart the session so SessionStart bootstraps it"
		}
		if !healed {
			dead = true
		}
		findings = append(findings, msg)
	}
	return findings, dead
}

// discoverInstalls returns wf installs from (a) the hook-context env vars,
// and (b) ~/.claude/plugins/installed_plugins.json — the Bash-tool copy of
// wf runs OUTSIDE hook context, where the env vars are unset, so (b) is
// what lets a plain `wf doctor --bootstrap` see the dead install at all.
func discoverInstalls(home string) []hookInstall {
	var installs []hookInstall
	seen := map[string]bool{}
	add := func(in hookInstall) {
		if in.data == "" || seen[in.data] {
			return
		}
		seen[in.data] = true
		installs = append(installs, in)
	}

	if data := os.Getenv("CLAUDE_PLUGIN_DATA"); data != "" {
		add(hookInstall{root: os.Getenv("CLAUDE_PLUGIN_ROOT"), data: data})
	}
	if home == "" {
		return installs
	}

	plugins := filepath.Join(home, ".claude", "plugins")
	raw, err := os.ReadFile(filepath.Join(plugins, "installed_plugins.json"))
	if err != nil {
		return installs
	}
	var reg struct {
		Plugins map[string][]struct {
			InstallPath string `json:"installPath"`
		} `json:"plugins"`
	}
	if err := json.Unmarshal(raw, &reg); err != nil {
		return installs
	}
	for key, entries := range reg.Plugins {
		// key is "<name>@<marketplace>"; the data dir is "<name>-<marketplace>"
		name, marketplace, ok := strings.Cut(key, "@")
		if !ok || name != "wf" {
			continue
		}
		root := ""
		for _, e := range entries {
			if e.InstallPath != "" {
				root = e.InstallPath
				break
			}
		}
		add(hookInstall{root: root, data: filepath.Join(plugins, "data", name+"-"+marketplace)})
	}
	return installs
}

// SelfUpdate keeps the hook engine current on EVERY platform without
// platform-scoped hook wiring (which Claude Code does not offer): all hooks
// run ${CLAUDE_PLUGIN_DATA}/bin/wf, so once installed the engine re-runs
// the bootstrap itself when the plugin root expects a different version —
// plugin updated mid-flight, or native Windows where the sh SessionStart
// entry cannot run. Called from `wf inject session` (hook context only; the
// env guard makes it a no-op everywhere else). Returns a note for the
// session block, "" when nothing was needed. Rate-limited to one attempt
// per expected version per hour so a broken bootstrap can't loop.
func SelfUpdate(engineVersion string) string {
	root := os.Getenv("CLAUDE_PLUGIN_ROOT")
	data := os.Getenv("CLAUDE_PLUGIN_DATA")
	if root == "" || data == "" || engineVersion == "" {
		return ""
	}
	expected := expectedVersion(root)
	if expected == "" || expected == engineVersion {
		return ""
	}
	stamp := filepath.Join(data, "bin", ".selfupdate")
	var last struct {
		Expected string `json:"expected"`
		TS       int64  `json:"ts"`
	}
	if raw, err := os.ReadFile(stamp); err == nil {
		if json.Unmarshal(raw, &last) == nil && last.Expected == expected && time.Since(time.Unix(last.TS, 0)) < time.Hour {
			return "" // already attempted recently; doctor --bootstrap is the manual path
		}
	}
	last.Expected, last.TS = expected, time.Now().Unix()
	if raw, err := json.Marshal(last); err == nil {
		_ = os.MkdirAll(filepath.Dir(stamp), 0o755)
		_ = os.WriteFile(stamp, raw, 0o644)
	}

	before := fileDigest(filepath.Join(data, "bin", "wf"))
	out, err := runBootstrap(root, data)
	after := fileDigest(filepath.Join(data, "bin", "wf"))
	switch {
	case err == nil && after != "" && after != before:
		return fmt.Sprintf("engine self-updated (%s → plugin's %s) — hooks run the new engine from the next event on", engineVersion, expected)
	case err != nil:
		return fmt.Sprintf("engine version skew (running %s, plugin expects %s) and the bootstrap failed: %v — run `wf doctor --bootstrap`", engineVersion, expected, err)
	default:
		return fmt.Sprintf("engine version skew (running %s, plugin expects %s) — bootstrap ran but installed nothing (%s); run `wf doctor --bootstrap`", engineVersion, expected, strings.TrimSpace(out))
	}
}

// expectedVersion is what the plugin root wants running: bin/VERSION when
// it ships (dev/release builds), else the committed MANIFEST's semver.
func expectedVersion(root string) string {
	if raw, err := os.ReadFile(filepath.Join(root, "bin", "VERSION")); err == nil {
		return strings.TrimSpace(string(raw))
	}
	raw, err := os.ReadFile(filepath.Join(root, "bin", "MANIFEST"))
	if err != nil {
		return ""
	}
	for _, line := range strings.Split(string(raw), "\n") {
		if v, ok := strings.CutPrefix(line, "version "); ok {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

// fileDigest is a change detector, not a security check ("" = unreadable).
func fileDigest(path string) string {
	raw, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(raw)
	return fmt.Sprintf("%x", sum[:8])
}

// runBootstrap executes the plugin's bootstrap script with the plugin env
// set — the same thing SessionStart would have done.
func runBootstrap(root, data string) (string, error) {
	sh := filepath.Join(root, "scripts", "bootstrap.sh")
	var cmd *exec.Cmd
	if _, err := exec.LookPath("sh"); err == nil {
		if _, err := os.Stat(sh); err != nil {
			return "", fmt.Errorf("no bootstrap script at %s", sh)
		}
		cmd = exec.Command("sh", sh)
	} else if runtime.GOOS == "windows" {
		ps1 := filepath.Join(root, "scripts", "bootstrap.ps1")
		if _, err := os.Stat(ps1); err != nil {
			return "", fmt.Errorf("no bootstrap script at %s", ps1)
		}
		cmd = exec.Command("powershell", "-NoProfile", "-ExecutionPolicy", "Bypass", "-File", ps1)
	} else {
		return "", fmt.Errorf("no sh on PATH to run %s", sh)
	}
	cmd.Env = append(os.Environ(),
		"CLAUDE_PLUGIN_ROOT="+root,
		"CLAUDE_PLUGIN_DATA="+data,
	)
	out, err := cmd.CombinedOutput()
	return string(out), err
}
