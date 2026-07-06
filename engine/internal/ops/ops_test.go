package ops

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

func newCtl(t *testing.T) (*runctl.Ctl, string) {
	t.Helper()
	dir := t.TempDir()
	s, err := store.Open(dir, true)
	if err != nil {
		t.Fatal(err)
	}
	p, _ := filepath.Abs(filepath.Join("..", "..", "..", "workflow", "workflow.yaml"))
	sp, err := spec.Load(p, "")
	if err != nil {
		t.Fatal(err)
	}
	return &runctl.Ctl{Store: s, Spec: sp, Config: &store.Config{}}, dir
}

func rec(t *testing.T, c *runctl.Ctl, kind string, data map[string]any) {
	t.Helper()
	if _, err := c.Record(kind, data, false, "agent"); err != nil {
		t.Fatalf("record %s: %v", kind, err)
	}
}

func TestDepsCheckVerdicts(t *testing.T) {
	c, dir := newCtl(t)
	_, _ = c.RunStart("diff", "fix")

	// no manifests, no strategies -> n/a
	out, err := DepsCheck(c, dir)
	if err != nil || !strings.Contains(out, "n/a") {
		t.Fatalf("empty project: want n/a, got %q (%v)", out, err)
	}

	// resolvable tool -> present
	rec(t, c, "verification-strategy", map[string]any{"ac": "AC-1", "method": "unit", "command": "git version"})
	out, _ = DepsCheck(c, dir)
	if !strings.Contains(out, "present") {
		t.Fatalf("resolvable tool: want present, got %q", out)
	}

	// unresolvable tool -> missing
	rec(t, c, "verification-strategy", map[string]any{"ac": "AC-2", "method": "unit", "command": "definitely-not-a-tool-xyz --run"})
	out, _ = DepsCheck(c, dir)
	if !strings.Contains(out, "missing") || !strings.Contains(out, "definitely-not-a-tool-xyz") {
		t.Fatalf("missing tool must be named: %q", out)
	}
	// the missing verdict is recorded (blocks the plan gate)
	r, _ := c.Store.LoadRun()
	env, _ := c.Env(r)
	deps := env.Records("deps")
	last := deps[len(deps)-1]
	if last.Data["verdict"] != "missing" {
		t.Fatalf("deps record: %v", last.Data)
	}
}

func TestOriginDiscover(t *testing.T) {
	if _, err := exec.LookPath("git"); err != nil {
		t.Skip("git unavailable")
	}
	c, dir := newCtl(t)
	_, _ = c.RunStart("diff", "fix")
	git := func(args ...string) {
		t.Helper()
		cmd := exec.Command("git", append([]string{"-C", dir}, args...)...)
		cmd.Env = append(os.Environ(),
			"GIT_AUTHOR_NAME=t", "GIT_AUTHOR_EMAIL=t@t", "GIT_COMMITTER_NAME=t", "GIT_COMMITTER_EMAIL=t@t")
		if out, err := cmd.CombinedOutput(); err != nil {
			t.Fatalf("git %v: %v %s", args, err, out)
		}
	}
	git("init", "-q")
	_ = os.WriteFile(filepath.Join(dir, "reader.go"), []byte("func Read() { crashHere() }\n"), 0o644)
	git("add", ".")
	git("commit", "-q", "-m", "introduce reader")

	out, err := OriginDiscover(c, dir, "reader.go", "crashHere")
	if err != nil {
		t.Fatalf("discover: %v", err)
	}
	if !strings.Contains(out, "introduce reader") {
		t.Errorf("attribution should carry the commit subject: %q", out)
	}
	r, _ := c.Store.LoadRun()
	env, _ := c.Env(r)
	if len(env.Records("origin")) != 1 {
		t.Fatal("origin record missing")
	}

	// inconclusive -> error with manual-fallback guidance, nothing recorded
	_, err = OriginDiscover(c, dir, "", "no-such-fragment-anywhere")
	if err == nil || !strings.Contains(err.Error(), "record manually") {
		t.Fatalf("inconclusive search must guide the fallback: %v", err)
	}
}

func TestDocNew(t *testing.T) {
	c, dir := newCtl(t)
	_, _ = c.RunStart("diff", "new")
	pluginRoot, _ := filepath.Abs(filepath.Join("..", "..", ".."))

	out, err := DocNew(c, pluginRoot, dir, "adr", "Use Event Log!")
	if err != nil {
		t.Fatal(err)
	}
	want := filepath.Join(dir, "docs", "architecture", "adr", "0001-use-event-log.md")
	if _, err := os.Stat(want); err != nil {
		t.Fatalf("adr not created at %s: %v\n%s", want, err, out)
	}
	// artifact recorded as stub
	r, _ := c.Store.LoadRun()
	env, _ := c.Env(r)
	arts := env.Records("artifact")
	if len(arts) != 1 || arts[0].Data["status"] != "stub" || arts[0].Data["template"] != "adr" {
		t.Fatalf("artifact record wrong: %+v", arts)
	}
	// numbering advances; duplicate slug at same number is a new file
	_, err = DocNew(c, pluginRoot, dir, "adr", "second decision")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(dir, "docs", "architecture", "adr", "0002-second-decision.md")); err != nil {
		t.Fatal("adr numbering must advance")
	}
	// unknown type refused
	if _, err := DocNew(c, pluginRoot, dir, "poem", "x"); err == nil {
		t.Fatal("unknown doc type must be refused")
	}
	// delivery-manifest carries its role
	if _, err := DocNew(c, pluginRoot, dir, "delivery-manifest", "v1"); err != nil {
		t.Fatal(err)
	}
	env, _ = c.Env(r)
	found := false
	for _, a := range env.Records("artifact") {
		if a.Data["role"] == "delivery-manifest" {
			found = true
		}
	}
	if !found {
		t.Fatal("delivery-manifest role missing")
	}
}
