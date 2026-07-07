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
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
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
