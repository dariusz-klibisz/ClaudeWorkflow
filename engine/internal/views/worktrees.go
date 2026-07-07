package views

// Cross-worktree aggregation (08 §7, 09 Q4 first half): each worktree
// carries its own .workflow/ run stream; `wf report --worktrees` reads the
// siblings' state — lock-free by design (gates read-only, 08 §7) — and
// groups signals per tree. Agent-teams/TeammateIdle integration stays
// deferred until teams stabilize (09 Q4).

import (
	"fmt"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/runctl"
)

// TreeReport is one worktree's aggregated signals.
type TreeReport struct {
	Worktree string       `json:"worktree"` // dir basename (label)
	Path     string       `json:"path"`
	Branch   string       `json:"branch,omitempty"`
	Signals  []RunSignals `json:"signals"`
	Note     string       `json:"note,omitempty"`
}

// ReportWorktrees aggregates across every ADOPTED worktree of the repo.
// openTree builds a Ctl for a sibling root (injected so spec/contracts.d
// loading stays with the CLI); un-adopted and legacy trees are skipped.
// A non-git projectDir degrades to the current tree with a note.
func ReportWorktrees(c *runctl.Ctl, projectDir string, openTree func(root string) (*runctl.Ctl, error)) ([]TreeReport, error) {
	trees, err := gitWorktrees(projectDir)
	if err != nil || len(trees) == 0 {
		sigs, rerr := Report(c)
		if rerr != nil {
			return nil, rerr
		}
		return []TreeReport{{
			Worktree: filepath.Base(projectDir), Path: projectDir, Signals: sigs,
			Note: "not a git worktree setup — current tree only",
		}}, nil
	}
	var out []TreeReport
	for _, wt := range trees {
		tc, err := openTree(wt.path)
		if err != nil {
			continue // un-adopted or legacy tree: not ours to report
		}
		sigs, err := Report(tc)
		if err != nil {
			continue
		}
		out = append(out, TreeReport{
			Worktree: filepath.Base(wt.path), Path: wt.path, Branch: wt.branch, Signals: sigs,
		})
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("no adopted worktrees found (run /wf:init)")
	}
	return out, nil
}

type worktree struct {
	path, branch string
}

// gitWorktrees parses `git worktree list --porcelain`.
func gitWorktrees(dir string) ([]worktree, error) {
	raw, err := exec.Command("git", "-C", dir, "worktree", "list", "--porcelain").Output()
	if err != nil {
		return nil, err
	}
	var out []worktree
	var cur *worktree
	for _, line := range strings.Split(string(raw), "\n") {
		switch {
		case strings.HasPrefix(line, "worktree "):
			if cur != nil {
				out = append(out, *cur)
			}
			cur = &worktree{path: strings.TrimPrefix(line, "worktree ")}
		case strings.HasPrefix(line, "branch ") && cur != nil:
			cur.branch = strings.TrimPrefix(strings.TrimPrefix(line, "branch "), "refs/heads/")
		}
	}
	if cur != nil {
		out = append(out, *cur)
	}
	return out, nil
}

// RenderTreeReports renders the grouped text view.
func RenderTreeReports(trees []TreeReport) string {
	var b strings.Builder
	fmt.Fprintf(&b, "[wf report] %d worktree(s)\n", len(trees))
	for _, t := range trees {
		label := t.Worktree
		if t.Branch != "" {
			label += " · " + t.Branch
		}
		fmt.Fprintf(&b, "── %s (%s)\n", label, t.Path)
		if t.Note != "" {
			fmt.Fprintf(&b, "   note: %s\n", t.Note)
		}
		for _, line := range strings.Split(strings.TrimRight(RenderReport(t.Signals), "\n"), "\n") {
			fmt.Fprintf(&b, "   %s\n", line)
		}
	}
	return b.String()
}
