package views

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/runctl"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/spec"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/store"
)

func gitc(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v\n%s", args, err, out)
	}
}

func openTree(t *testing.T) func(root string) (*runctl.Ctl, error) {
	t.Helper()
	p, _ := filepath.Abs(filepath.Join("..", "..", "..", "workflow", "workflow.yaml"))
	sp, err := spec.Load(p, "")
	if err != nil {
		t.Fatal(err)
	}
	return func(root string) (*runctl.Ctl, error) {
		s, err := store.Open(root, false)
		if err != nil {
			return nil, err
		}
		return &runctl.Ctl{Store: s, Spec: sp, Config: &store.Config{}}, nil
	}
}

// Two adopted worktrees report as two groups; an un-adopted third is skipped.
func TestReportWorktrees(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("no git")
	}
	base := t.TempDir()
	main := filepath.Join(base, "main")
	if err := os.MkdirAll(main, 0o755); err != nil {
		t.Fatal(err)
	}
	gitc(t, main, "init", "-q", "-b", "trunk")
	if err := os.WriteFile(filepath.Join(main, "f.txt"), []byte("x"), 0o644); err != nil {
		t.Fatal(err)
	}
	gitc(t, main, "add", "-A")
	gitc(t, main, "commit", "-qm", "init")
	feat := filepath.Join(base, "feat")
	bare := filepath.Join(base, "bare")
	gitc(t, main, "worktree", "add", "-q", "-b", "feature-x", feat)
	gitc(t, main, "worktree", "add", "-q", "-b", "unadopted", bare)

	open := openTree(t)
	adopt := func(root, family, intent string) *runctl.Ctl {
		t.Helper()
		s, err := store.Open(root, true)
		if err != nil {
			t.Fatal(err)
		}
		c, _ := open(root)
		c.Store = s
		if _, err := c.RunStart(family, intent); err != nil {
			t.Fatal(err)
		}
		return c
	}
	cMain := adopt(main, "diff", "fix")
	adopt(feat, "assessment", "investigate")
	// bare stays un-adopted

	trees, err := ReportWorktrees(cMain, main, open)
	if err != nil {
		t.Fatal(err)
	}
	if len(trees) != 2 {
		t.Fatalf("want 2 adopted trees, got %d: %+v", len(trees), trees)
	}
	byName := map[string]TreeReport{}
	for _, tr := range trees {
		byName[tr.Worktree] = tr
	}
	if tr, ok := byName["main"]; !ok || tr.Branch != "trunk" || len(tr.Signals) != 1 {
		t.Fatalf("main tree wrong: %+v", byName["main"])
	}
	if tr, ok := byName["feat"]; !ok || tr.Branch != "feature-x" || tr.Signals[0].Family != "assessment" {
		t.Fatalf("feat tree wrong: %+v", byName["feat"])
	}
	out := RenderTreeReports(trees)
	if !strings.Contains(out, "2 worktree(s)") || !strings.Contains(out, "feature-x") {
		t.Fatalf("render:\n%s", out)
	}
}

// A non-git project degrades to the current tree with a note.
func TestReportWorktreesNonGit(t *testing.T) {
	c := newCtl(t)
	_, _ = c.RunStart("diff", "fix")
	trees, err := ReportWorktrees(c, c.Store.Root, openTree(t))
	if err != nil {
		t.Fatal(err)
	}
	if len(trees) != 1 || trees[0].Note == "" || len(trees[0].Signals) != 1 {
		t.Fatalf("non-git fallback wrong: %+v", trees)
	}
}
