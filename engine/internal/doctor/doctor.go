// Package doctor verifies workflow-state health: spec validity, snapshot/log
// consistency, torn events, stale locks, and idle runs. It is the repair path
// when gates fail open (04 §7).
package doctor

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/runctl"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/spec"
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

	// idle run (E2)
	if snap != nil && snap.Started != "" {
		if t, err := time.Parse(time.RFC3339, snap.Started); err == nil && time.Since(t) > 30*24*time.Hour && snap.Status == "active" {
			f = append(f, fmt.Sprintf("run %s idle/open for >30 days — consider wf park or wf run close", snap.ID))
		}
	}

	return Report{OK: len(f) == 0, Findings: f}
}

func (r Report) String() string {
	if r.OK {
		return "wf doctor: all checks passed"
	}
	return "wf doctor: " + fmt.Sprint(len(r.Findings)) + " finding(s)\n  - " + strings.Join(r.Findings, "\n  - ")
}
