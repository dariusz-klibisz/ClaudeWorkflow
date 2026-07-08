// Package doctor verifies workflow-state health: spec validity, snapshot/log
// consistency, torn events, stale locks, and idle runs. It is the repair path
// when gates fail open (04 §7).
package doctor

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/runctl"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/spec"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/store"
)

type Report struct {
	OK       bool
	Findings []string
}

func Run(c *runctl.Ctl, specPath string) Report {
	var f []string

	// strict spec re-parse: surfaces unknown fields (this engine may be older
	// than the spec — tolerated at runtime, reported here)
	if specPath != "" {
		contractsDir := ""
		if _, err := os.Stat(c.Store.ContractsDir()); err == nil {
			contractsDir = c.Store.ContractsDir()
		}
		if _, err := spec.LoadStrict(specPath, contractsDir); err != nil {
			f = append(f, fmt.Sprintf("spec has fields this engine version doesn't know (tolerated at runtime; update the plugin/engine): %v", err))
		}
	}

	// snapshot vs log consistency (merge-recovery path)
	snap, err := c.Store.LoadRun()
	if err != nil {
		f = append(f, fmt.Sprintf("run snapshot unreadable: %v (repair: wf run adopt)", err))
	}
	derived, err := c.Store.DeriveRun()
	if err != nil {
		f = append(f, fmt.Sprintf("event log unreadable: %v", err))
	}
	switch {
	case snap == nil && derived != nil:
		f = append(f, fmt.Sprintf("log shows in-flight run %s but no snapshot — a fresh clone? repair: wf run adopt", derived.ID))
	case snap != nil && derived == nil:
		f = append(f, fmt.Sprintf("snapshot names run %s but the log has no such open run — repair: wf run adopt (re-derives) or remove state/run.json", snap.ID))
	case snap != nil && derived != nil && snap.ID != derived.ID:
		f = append(f, fmt.Sprintf("snapshot run %s != log-derived run %s — repair: wf run adopt", snap.ID, derived.ID))
	}

	// stale lock
	if st, err := os.Stat(c.Store.Root + "/state/lock"); err == nil {
		if time.Since(st.ModTime()) > time.Minute {
			f = append(f, "stale lockfile state/lock (crashed writer?) — safe to remove")
		}
	}

	// ledger integrity: line hash chain + unparseable lines (store.scan
	// tolerates torn lines silently; doctor is where they get reported)
	if chain, err := c.Store.VerifyChain(); err != nil {
		f = append(f, fmt.Sprintf("event log unscannable: %v", err))
	} else {
		if chain.Unparseable > 0 {
			f = append(f, fmt.Sprintf("event log has %d unparseable line(s) — torn writes or foreign edits (events on those lines are invisible to gates)", chain.Unparseable))
		}
		const maxShown = 5
		for i, b := range chain.Breaks {
			if i >= maxShown {
				f = append(f, fmt.Sprintf("… and %d more chain finding(s)", len(chain.Breaks)-maxShown))
				break
			}
			f = append(f, "ledger chain: "+b)
		}
	}

	// hook liveness: a run past Frame with many events but zero hook-captured
	// ones means the enforcement spine is not firing (the dead-hooks
	// incident: bootstrap failed and every gate ENOENT'd silently)
	if snap != nil && snap.Status == "active" {
		if msg := HookLiveness(c, snap); msg != "" {
			f = append(f, msg)
		}
	}
	// engine reachable at the hook path? Discovered from hook-context env
	// AND installed_plugins.json, so the Bash-tool copy of wf (different
	// install, no plugin env) can still see a dead hook engine. Report
	// only — `wf doctor --bootstrap` is the heal path.
	hookFindings, _ := HookEngineFindings("", false)
	f = append(f, hookFindings...)

	// live-log growth: the log holds the active run + open followups only;
	// unusual size means unclosed runs piling up (informational threshold)
	if st, err := os.Stat(c.Store.EventsPath()); err == nil {
		const liveLogWarnBytes = 2 << 20 // 2 MiB
		if st.Size() > liveLogWarnBytes {
			f = append(f, fmt.Sprintf("live event log is %.1f MiB — it should hold only the active run + open followups; close finished runs (wf run close archives and compacts)", float64(st.Size())/(1<<20)))
		}
	}

	// idle run (E2)
	if snap != nil && snap.Started != "" {
		if t, err := time.Parse(time.RFC3339, snap.Started); err == nil && time.Since(t) > 30*24*time.Hour && snap.Status == "active" {
			f = append(f, fmt.Sprintf("run %s idle/open for >30 days — consider wf park or wf run close", snap.ID))
		}
	}

	// corpus snapshot age: reviewers cite rule IDs from these snapshots —
	// a stale bundle ages every citation (informational; the corpora
	// workflow does the authoritative drift check against the sources)
	f = append(f, corpusFindings(c.Spec.PluginRoot())...)

	return Report{OK: len(f) == 0, Findings: f}
}

// corpusFindings reports missing or stale (>180 days) bundled corpora.
func corpusFindings(pluginRoot string) []string {
	if pluginRoot == "" {
		return nil
	}
	var f []string
	for _, name := range []string{"design", "coding", "ux"} {
		ver := filepath.Join(pluginRoot, "reference", name, "VERSION")
		raw, err := os.ReadFile(ver)
		if err != nil {
			f = append(f, fmt.Sprintf("reference corpus %q missing its VERSION stamp — agents fall back to model knowledge (reinstall the plugin or run scripts/sync-corpora.sh)", name))
			continue
		}
		for _, line := range strings.Split(string(raw), "\n") {
			if d, ok := strings.CutPrefix(line, "date: "); ok {
				if t, err := time.Parse("2006-01-02", strings.TrimSpace(d)); err == nil && time.Since(t) > 180*24*time.Hour {
					f = append(f, fmt.Sprintf("reference corpus %q snapshot is %d days old — reviewer citations age with it (maintainer: scripts/sync-corpora.sh)", name, int(time.Since(t).Hours()/24)))
				}
			}
		}
	}
	return f
}

// HookLiveness returns a warning when the run's ledger shows no
// hook-captured events despite substantial activity past Frame — the
// signature of dead hooks. Empty string = healthy or not yet judgeable.
func HookLiveness(c *runctl.Ctl, snap *store.Run) string {
	past, pastPlan := false, false
	for _, ph := range snap.ExitedPh {
		if ph == "frame" {
			past = true
		}
		if ph == "plan" {
			pastPlan = true
		}
	}
	if !past {
		return ""
	}
	evs, err := c.Store.RunEvents(snap.ID)
	if err != nil {
		return ""
	}
	// Signal 1 (the power5 incident's signature): reviewer verdicts exist but
	// NONE was auto-captured — the SubagentStop gate is not firing.
	verdicts, autoVerdicts := 0, 0
	testRuns, autoTestRuns := 0, 0
	hookEvents := 0
	for _, e := range evs {
		if e.Kind == "verdict" {
			verdicts++
			if e.Auto {
				autoVerdicts++
			}
		}
		if e.Kind == "test-run" {
			testRuns++
			if e.Auto {
				autoTestRuns++
			}
		}
		if e.Actor == "hook" {
			hookEvents++
		}
	}
	if verdicts >= 2 && autoVerdicts == 0 {
		return fmt.Sprintf("run %s has %d reviewer verdicts and NONE auto-captured — the SubagentStop gate appears DEAD (bootstrap/hook wiring). Manual verdicts are honest but unenforced; run `wf doctor --bootstrap` (installs the hook engine on the spot)", snap.ID, verdicts)
	}
	// Signal 3 (the multiply-app incident's signature): tests are being
	// recorded by hand past Plan while Bash capture records none — either
	// the PostToolUse hook is dead, or the test runner isn't recognized.
	// Diff family only: artifact/assessment runs verify documents with
	// manual checks by design (the arch-design run's false positive).
	if snap.Family == "diff" && pastPlan && testRuns >= 3 && autoTestRuns == 0 {
		return fmt.Sprintf("run %s has %d test-run records and NONE auto-captured — Bash test capture is not firing. Either the hook is dead (`wf doctor --bootstrap`), or the test runner isn't recognized: make the recorded verification-strategy commands match the real invocations, or declare custom runners in .workflow/config.json (\"runners\": [\"./scripts/test.sh\"])", snap.ID, testRuns)
	}
	// Signal 2 (broader): a substantial ledger with zero hook-side events.
	if len(evs) >= 15 && hookEvents == 0 {
		return fmt.Sprintf("run %s has %d events past Frame but ZERO hook-captured ones — enforcement hooks appear DEAD (bootstrap/permissions). Gates, verdict capture, and test grounding are not firing; run `wf doctor --bootstrap` (installs the hook engine on the spot)", snap.ID, len(evs))
	}
	return ""
}

func (r Report) String() string {
	if r.OK {
		return "wf doctor: all checks passed"
	}
	return "wf doctor: " + fmt.Sprint(len(r.Findings)) + " finding(s)\n  - " + strings.Join(r.Findings, "\n  - ")
}
