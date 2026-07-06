package gates

import (
	"encoding/json"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/hookio"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
)

func gitIn(t *testing.T, dir string, args ...string) {
	t.Helper()
	cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
	cmd.Env = append(os.Environ(),
		"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
	if out, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("git %v: %v %s", args, err, out)
	}
}

func commitPayload(t *testing.T, cwd, cmd string) *hookio.Input {
	t.Helper()
	raw, _ := json.Marshal(map[string]any{
		"hook_event_name": "PostToolUse", "tool_name": "Bash", "cwd": cwd,
		"tool_input": map[string]any{"command": cmd}, "tool_response": map[string]any{"exit_code": 0},
	})
	return hookInput(t, string(raw))
}

func TestCaptureCommit(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git unavailable")
	}
	c := newCtl(t)
	run, _ := c.RunStart("diff", "fix")
	dir := filepath.Dir(c.Store.Root) // the project dir
	gitIn(t, dir, "init", "-q")
	_ = os.WriteFile(filepath.Join(dir, "a.txt"), []byte("x"), 0o644)
	gitIn(t, dir, "add", ".")
	gitIn(t, dir, "commit", "-q", "-m", "untagged change")

	// non-commit commands are ignored
	_ = CaptureCommit(c, commitPayload(t, dir, "git status && ls"))
	env, _ := c.Env(run)
	if len(env.Records("commit-origin")) != 0 {
		t.Fatal("non-commit command must not record")
	}
	// echo mentioning git commit in a string must not record (head matching)
	_ = CaptureCommit(c, commitPayload(t, dir, `echo "please git commit later"`))
	env, _ = c.Env(run)
	if len(env.Records("commit-origin")) != 0 {
		t.Fatal("quoted mention must not record")
	}

	// a real commit records, flags the missing tag
	r := CaptureCommit(c, commitPayload(t, dir, `git add . && git commit -m "untagged change"`))
	env, _ = c.Env(run)
	cos := env.Records("commit-origin")
	if len(cos) != 1 {
		t.Fatalf("want 1 commit-origin, got %d", len(cos))
	}
	if tagged, _ := cos[0].Data["tagged"].(bool); tagged {
		t.Error("untagged commit must be flagged tagged=false")
	}
	if r.Stdout == "" {
		t.Error("untagged commit should surface a reminder")
	}

	// idempotent on the same sha
	_ = CaptureCommit(c, commitPayload(t, dir, `git commit --amend --no-edit`))
	// amend changes sha -> records; but repeating the SAME sha must not duplicate
	_ = CaptureCommit(c, commitPayload(t, dir, `git commit -q -m x || true`)) // no-op commit fails; HEAD unchanged
	env, _ = c.Env(run)
	n := len(env.Records("commit-origin"))
	if n > 2 {
		t.Fatalf("same sha must not duplicate: got %d records", n)
	}

	// tagged commit passes silently
	_ = os.WriteFile(filepath.Join(dir, "b.txt"), []byte("y"), 0o644)
	gitIn(t, dir, "add", ".")
	gitIn(t, dir, "commit", "-q", "-m", "fix [run:"+run.ID+"]")
	r = CaptureCommit(c, commitPayload(t, dir, `git commit -m "fix [run:`+run.ID+`]"`))
	if r.Stdout != "" {
		t.Errorf("tagged commit should be silent: %s", r.Stdout)
	}
}
