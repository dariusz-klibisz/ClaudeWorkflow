package runctl

import (
	"path/filepath"
	"strings"
	"testing"

	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/spec"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/store"
)

func newCtl(t *testing.T) *Ctl {
	t.Helper()
	s, err := store.Open(t.TempDir(), true)
	if err != nil {
		t.Fatal(err)
	}
	p, _ := filepath.Abs(filepath.Join("..", "..", "..", "workflow", "workflow.yaml"))
	sp, err := spec.Load(p, "")
	if err != nil {
		t.Fatal(err)
	}
	return &Ctl{Store: s, Spec: sp, Config: &store.Config{}}
}

func TestRunStartValidation(t *testing.T) {
	c := newCtl(t)
	if _, err := c.RunStart("nonsense", ""); err == nil {
		t.Fatal("unknown family must be rejected")
	}
	if _, err := c.RunStart("diff", "doc-create"); err == nil {
		t.Fatal("wrong-family intent must be rejected")
	}
	r, err := c.RunStart("diff", "fix")
	if err != nil {
		t.Fatal(err)
	}
	if r.Phase != "frame" {
		t.Errorf("first phase must be frame, got %s", r.Phase)
	}
	if _, err := c.RunStart("diff", "fix"); err == nil {
		t.Fatal("second concurrent run must be rejected")
	}
}

func TestPhaseExitBlocksOnFindings(t *testing.T) {
	c := newCtl(t)
	_, _ = c.RunStart("diff", "fix")
	findings, _, err := c.PhaseExit(false, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) == 0 {
		t.Fatal("empty frame must not exit")
	}
	r, _ := c.Store.LoadRun()
	if r.Phase != "frame" {
		t.Error("phase must not advance on findings")
	}
}

// satisfyFrame records everything the diff/fix frame contract needs.
func satisfyFrame(t *testing.T, c *Ctl) {
	t.Helper()
	rec := func(kind string, data map[string]any) {
		t.Helper()
		if _, err := c.Record(kind, data, false, "agent"); err != nil {
			t.Fatalf("record %s: %v", kind, err)
		}
	}
	rec("classification", map[string]any{"family": "diff", "intent": "fix", "restated": "fix the empty-file crash"})
	rec("risk", map[string]any{"signals": []any{"data"}, "lenses": []any{"security", "adversarial"}})
	rec("ambiguity", map[string]any{"lens": "security", "text": "input validated?", "disposition": "resolved"})
	rec("ambiguity", map[string]any{"lens": "adversarial", "none": true, "disposition": "none"})
	rec("requirement", map[string]any{"rid": "SWR-1", "level": "software", "text": "handle empty files", "status": "active",
		"acs": []any{map[string]any{"id": "AC-1", "text": "empty file yields error msg", "verifiable": true}}})
	rec("completeness", map[string]any{"items": []any{map[string]any{"case": "empty", "disposition": "covered"}}})
	rec("verdict", map[string]any{"agent": "adversary", "scope": "abuse-case", "status": "clean", "criticals": 0, "majors": 0})
	rec("verdict", map[string]any{"agent": "lens-reviewer", "scope": "security", "status": "clean", "criticals": 0, "majors": 0})
	rec("origin", map[string]any{"attribution": "introduced in abc123"})
	if _, err := c.Approve("frame", "diff/fix: fix empty-file crash"); err != nil {
		t.Fatal(err)
	}
}

func TestFullFrameExits(t *testing.T) {
	c := newCtl(t)
	_, _ = c.RunStart("diff", "fix")
	satisfyFrame(t, c)
	findings, msg, err := c.PhaseExit(false, "")
	if err != nil {
		t.Fatal(err)
	}
	if len(findings) > 0 {
		t.Fatalf("frame should exit, findings: %+v", findings)
	}
	if !strings.Contains(msg, "context") {
		t.Errorf("should enter context: %q", msg)
	}
}

func TestForceEscalation(t *testing.T) {
	c := newCtl(t)
	_, _ = c.RunStart("diff", "fix")
	// 1st force: reason enough
	if _, _, err := c.PhaseExit(true, "demo bypass"); err != nil {
		t.Fatalf("1st force: %v", err)
	}
	// 2nd force: must name a structural cause
	if _, _, err := c.PhaseExit(true, "again"); err == nil {
		t.Fatal("2nd force without cause: must refuse")
	}
	if _, _, err := c.PhaseExit(true, "cause: gate misfires on assessment-shaped work"); err != nil {
		t.Fatalf("2nd force with cause: %v", err)
	}
	// 3rd force: auto-park
	_, _, err := c.PhaseExit(true, "cause: whatever")
	if err == nil || !strings.Contains(err.Error(), "auto-parked") {
		t.Fatalf("3rd force must auto-park: %v", err)
	}
	r, _ := c.Store.LoadRun()
	if r.Status != "parked" {
		t.Errorf("run must be parked, got %s", r.Status)
	}
}

func TestParkResumeAlwaysAvailable(t *testing.T) {
	c := newCtl(t)
	_, _ = c.RunStart("diff", "fix")
	if err := c.Park(""); err == nil {
		t.Fatal("park requires reason")
	}
	if err := c.Park("blocked on credentials"); err != nil {
		t.Fatal(err)
	}
	r, _ := c.Store.LoadRun()
	if r.Status != "parked" {
		t.Fatal("not parked")
	}
	if _, _, err := c.PhaseExit(false, ""); err == nil {
		t.Fatal("parked run must not exit phases")
	}
	if _, err := c.Resume(); err != nil {
		t.Fatal(err)
	}
	r, _ = c.Store.LoadRun()
	if r.Status != "active" {
		t.Fatal("not resumed")
	}
}

func TestLoopCaps(t *testing.T) {
	c := newCtl(t)
	r, _ := c.RunStart("diff", "fix")
	r.Phase = "verify"
	_ = c.Store.SaveRun(r)
	if _, err := c.Loop("AC-1", "slip", ""); err == nil {
		t.Fatal("loop without evidence must refuse")
	}
	if _, err := c.Loop("AC-1", "slip", "expected err msg, got panic"); err != nil {
		t.Fatal(err)
	}
	r, _ = c.Store.LoadRun()
	if r.Phase != "build" || r.Loops != 1 || r.SlipByAC["AC-1"] != 1 {
		t.Fatalf("loop state wrong: %+v", r)
	}
	// 2nd slip on same AC ok, 3rd refused
	r.Phase = "verify"
	_ = c.Store.SaveRun(r)
	if _, err := c.Loop("AC-1", "slip", "still failing"); err != nil {
		t.Fatal(err)
	}
	r, _ = c.Store.LoadRun()
	r.Phase = "verify"
	_ = c.Store.SaveRun(r)
	if _, err := c.Loop("AC-1", "slip", "third time"); err == nil {
		t.Fatal("3rd slip-loop on one AC must force a structural cause")
	}
	if _, err := c.Loop("AC-1", "design", "the design lacks a state"); err != nil {
		t.Fatalf("design-cause loop after slip cap: %v", err)
	}
}

func TestACVerdictGrounding(t *testing.T) {
	c := newCtl(t)
	_, _ = c.RunStart("diff", "fix")
	if _, err := c.Record("ac-verdict", map[string]any{"ac": "AC-1", "status": "pass"}, false, "agent"); err == nil {
		t.Fatal("ungrounded AC pass must be refused at write time")
	}
	// grounded green enables it
	if _, err := c.Record("test-run", map[string]any{"cmd": "go test", "exit": 0, "ac": "AC-1", "grounded": true}, true, "hook"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Record("ac-verdict", map[string]any{"ac": "AC-1", "status": "pass"}, false, "agent"); err != nil {
		t.Fatalf("grounded pass refused: %v", err)
	}
	// deferred needs approval
	if _, err := c.Record("ac-verdict", map[string]any{"ac": "AC-2", "status": "deferred"}, false, "agent"); err == nil {
		t.Fatal("deferred without approval must be refused")
	}
	_, _ = c.Approve("deferral", "AC-2 deferred to next run")
	if _, err := c.Record("ac-verdict", map[string]any{"ac": "AC-2", "status": "deferred"}, false, "agent"); err != nil {
		t.Fatalf("approved deferral refused: %v", err)
	}
}

func TestVerdictWriteValidation(t *testing.T) {
	c := newCtl(t)
	_, _ = c.RunStart("diff", "fix")
	// manual contradictory verdict refused
	if _, err := c.Record("verdict", map[string]any{"agent": "critic", "status": "clean", "criticals": 2, "majors": 0}, false, "agent"); err == nil {
		t.Fatal("manual clean-with-criticals must be refused")
	}
	// auto contradictory verdict downgraded
	ev, err := c.Record("verdict", map[string]any{"agent": "critic", "status": "clean", "criticals": 2, "majors": 0}, true, "hook")
	if err != nil {
		t.Fatal(err)
	}
	if ev.Data["status"] != "changes-required" {
		t.Errorf("auto verdict must downgrade, got %v", ev.Data["status"])
	}
	// unknown status refused
	if _, err := c.Record("verdict", map[string]any{"agent": "critic", "status": "lgtm", "criticals": 0, "majors": 0}, false, "agent"); err == nil {
		t.Fatal("unknown verdict status must be refused")
	}
}

func TestPhaseWaiveDesignForDiff(t *testing.T) {
	c := newCtl(t)
	_, _ = c.RunStart("diff", "fix")
	if err := c.PhaseWaive("design", "one-line fix, no design decisions"); err != nil {
		t.Fatal(err)
	}
	satisfyFrame(t, c)
	_, msg, err := c.PhaseExit(false, "")
	if err != nil {
		t.Fatal(err)
	}
	_ = msg
	// exit context (satisfy quickly)
	rec := func(kind string, data map[string]any) {
		if _, err := c.Record(kind, data, false, "agent"); err != nil {
			t.Fatalf("%s: %v", kind, err)
		}
	}
	rec("context-map", map[string]any{"entries": []any{"pkg/reader.go"}, "sufficiency": "single file change"})
	rec("reclassify", map[string]any{"result": "confirmed"})
	rec("waiver", map[string]any{"item": "context.research-grounded", "reason": "no research"})
	_, _ = c.Approve("scope", "scope ok")
	_, msg, err = c.PhaseExit(false, "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(msg, "plan") {
		t.Errorf("waived design must be skipped: %q", msg)
	}
	// artifact family cannot waive design
	c2 := newCtl(t)
	_, _ = c2.RunStart("artifact", "arch-design")
	if err := c2.PhaseWaive("design", "nah"); err == nil {
		t.Fatal("artifact design is required, not waivable")
	}
}

// Regression for the power5 run: chaining updates= onto a prior update's
// event ID must resolve to the original record — and unknown targets must be
// rejected at write time, never silently no-op.
func TestUpdatesChainingAndValidation(t *testing.T) {
	c := newCtl(t)
	r, _ := c.RunStart("diff", "fix")
	orig, err := c.Record("task", map[string]any{"tid": "T-1", "subject": "s", "dod": []any{"d"}, "status": "open"}, false, "agent")
	if err != nil {
		t.Fatal(err)
	}
	up1, err := c.Record("task", map[string]any{"updates": orig.ID, "status": "in_progress"}, false, "agent")
	if err != nil {
		t.Fatal(err)
	}
	// the live bug: targeting the UPDATE event's ID instead of the original
	if _, err := c.Record("task", map[string]any{"updates": up1.ID, "status": "done"}, false, "agent"); err != nil {
		t.Fatalf("chained update must resolve transitively: %v", err)
	}
	env, _ := c.Env(r)
	tasks := env.Records("task")
	if len(tasks) != 1 {
		t.Fatalf("want 1 effective task, got %d", len(tasks))
	}
	if s, _ := tasks[0].Data["status"].(string); s != "done" {
		t.Fatalf("chained update lost: status=%s", s)
	}
	// unknown target: hard error, not a silent no-op
	if _, err := c.Record("task", map[string]any{"updates": "01NOTAREALID0000000000000", "status": "done"}, false, "agent"); err == nil {
		t.Fatal("update targeting an unknown record must be rejected")
	}
}

// Regression for the power5 run: verdicts recorded under scoped or unknown
// agent names silently failed the contracts.
func TestVerdictAgentNameValidation(t *testing.T) {
	c := newCtl(t)
	r, _ := c.RunStart("diff", "fix")
	// scoped name normalizes
	ev, err := c.Record("verdict", map[string]any{"agent": "wf:adversary", "status": "clean", "criticals": 0, "majors": 0, "scope": "abuse-case"}, false, "agent")
	if err != nil {
		t.Fatal(err)
	}
	if ev.Data["agent"] != "adversary" {
		t.Errorf("scoped name must normalize: %v", ev.Data["agent"])
	}
	// unknown name rejected with the roster in the message
	_, err = c.Record("verdict", map[string]any{"agent": "abuse-case-analyst", "status": "clean", "criticals": 0, "majors": 0}, false, "agent")
	if err == nil || !strings.Contains(err.Error(), "roster") {
		t.Fatalf("unknown agent must be rejected naming the roster: %v", err)
	}
	_ = r
}

func TestBranchCarriesLineage(t *testing.T) {
	c := newCtl(t)
	parent, _ := c.RunStart("diff", "fix")
	child, err := c.RunBranch("diff", "refactor", "reclassify: scope grew")
	if err != nil {
		t.Fatal(err)
	}
	if child.Parent != parent.ID {
		t.Error("child must carry parent lineage")
	}
	if child.Intent != "refactor" {
		t.Error("child intent not applied")
	}
}

func TestAdoptFromLog(t *testing.T) {
	c := newCtl(t)
	r, _ := c.RunStart("diff", "fix")
	// simulate a fresh clone: snapshot gone, log present
	_ = c.Store.ClearRun()
	got, err := c.RunAdopt()
	if err != nil {
		t.Fatal(err)
	}
	if got.ID != r.ID || got.Phase != "frame" || got.Family != "diff" {
		t.Fatalf("adopt mismatch: %+v", got)
	}
}

func TestRunCloseRequiresShipExit(t *testing.T) {
	c := newCtl(t)
	_, _ = c.RunStart("diff", "fix")
	if err := c.RunClose(); err == nil {
		t.Fatal("close before ship exit must refuse")
	}
}
