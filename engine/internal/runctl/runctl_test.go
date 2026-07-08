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
	orig, err := c.Record("origin", map[string]any{"attribution": "introduced in abc123"}, false, "agent")
	if err != nil {
		t.Fatalf("record origin: %v", err)
	}
	// fix intent: the regression requirement traces to the origin
	// (frame.fix-regression)
	rec("requirement", map[string]any{"rid": "SWR-1", "level": "software", "text": "handle empty files", "status": "active",
		"origin": orig.ID,
		"acs": []any{map[string]any{"id": "AC-1", "text": "empty file yields error msg", "verifiable": true}}})
	rec("completeness", map[string]any{"items": []any{
		map[string]any{"case": "empty", "disposition": "covered:AC-1"},
		map[string]any{"case": "error", "disposition": "covered:AC-1"},
		map[string]any{"case": "max", "disposition": "accepted-risk: bounded input"},
	}})
	rec("verdict", map[string]any{"agent": "adversary", "scope": "abuse-case", "status": "clean", "criticals": 0, "majors": 0})
	rec("verdict", map[string]any{"agent": "lens-reviewer", "scope": "security", "status": "clean", "criticals": 0, "majors": 0})
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

// Ship-stage discoveries loop back to Verify (cause=audit) instead of
// dispositions-or-nothing — but only with grounds (open finding / failing
// audit) and never as a phase-order escape.
func TestShipAuditLoop(t *testing.T) {
	c := newCtl(t)
	r, _ := c.RunStart("diff", "fix")
	r.Phase = "ship"
	r.ExitedPh = []string{"frame", "context", "design", "plan", "build", "verify"}
	_ = c.Store.SaveRun(r)
	// no grounds: refused
	if _, err := c.Loop("AC-1", "audit", "auditor found the RTM contradicts the diff"); err == nil {
		t.Fatal("audit loop without grounds must be refused")
	}
	// a clean audit is not grounds
	if _, err := c.Record("verdict", map[string]any{"agent": "auditor", "status": "clean", "criticals": 0, "majors": 0}, true, "hook"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Loop("AC-1", "audit", "e"); err == nil {
		t.Fatal("clean latest audit must not ground a loop")
	}
	// a failing LATEST audit is grounds
	if _, err := c.Record("verdict", map[string]any{"agent": "auditor", "status": "changes-required", "criticals": 1, "majors": 0}, true, "hook"); err != nil {
		t.Fatal(err)
	}
	target, err := c.Loop("AC-1", "audit", "auditor critical: delivery artifact contradicts verified state")
	if err != nil {
		t.Fatal(err)
	}
	if target != "verify" {
		t.Fatalf("audit loop must target verify, got %s", target)
	}
	r, _ = c.Store.LoadRun()
	if r.Phase != "verify" || r.Loops != 1 || contains(r.ExitedPh, "verify") {
		t.Fatalf("loop state wrong: %+v", r)
	}
	// slip/design/plan causes still refuse from ship
	r.Phase = "ship"
	_ = c.Store.SaveRun(r)
	if _, err := c.Loop("AC-1", "slip", "e"); err == nil {
		t.Fatal("slip cause must not loop from ship")
	}
	// audit cause refuses from verify
	r, _ = c.Store.LoadRun()
	r.Phase = "verify"
	_ = c.Store.SaveRun(r)
	if _, err := c.Loop("AC-1", "audit", "e"); err == nil {
		t.Fatal("audit cause must not loop from verify")
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

// The AC-less-requirement hole: a requirement without ACs vacuously passed
// every downstream per-AC item. Refused at write time now.
func TestRequirementRequiresACs(t *testing.T) {
	c := newCtl(t)
	_, _ = c.RunStart("diff", "fix")
	base := map[string]any{"rid": "SWR-1", "level": "software", "text": "t", "status": "active"}
	with := func(acs any) map[string]any {
		d := map[string]any{}
		for k, v := range base {
			d[k] = v
		}
		d["acs"] = acs
		return d
	}
	if _, err := c.Record("requirement", with([]any{}), false, "agent"); err == nil {
		t.Fatal("empty acs must be refused")
	}
	if _, err := c.Record("requirement", with([]any{map[string]any{"text": "no id"}}), false, "agent"); err == nil {
		t.Fatal("AC without id must be refused")
	}
	if _, err := c.Record("requirement", with([]any{"  "}), false, "agent"); err == nil {
		t.Fatal("blank string AC must be refused")
	}
	// string ACs and {id,text} objects are both fine
	ev, err := c.Record("requirement", with([]any{"AC-1: empty file yields error"}), false, "agent")
	if err != nil {
		t.Fatalf("string AC refused: %v", err)
	}
	// an update that doesn't touch acs (Context baseline status flip) passes
	if _, err := c.Record("requirement", map[string]any{"updates": ev.ID, "status": "dropped"}, false, "agent"); err != nil {
		t.Fatalf("acs-less update refused: %v", err)
	}
	// an update that empties acs is refused
	if _, err := c.Record("requirement", map[string]any{"updates": ev.ID, "acs": []any{}}, false, "agent"); err == nil {
		t.Fatal("update emptying acs must be refused")
	}
}

// Judgment records get structural floors at write time (the C10 class).
func TestJudgmentRecordWriteValidation(t *testing.T) {
	c := newCtl(t)
	_, _ = c.RunStart("diff", "fix")
	if _, err := c.Record("context-map", map[string]any{"entries": []any{}, "sufficiency": "ok"}, false, "agent"); err == nil {
		t.Fatal("empty context-map entries must be refused")
	}
	if _, err := c.Record("context-map", map[string]any{"entries": []any{"a.go"}, "sufficiency": "  "}, false, "agent"); err == nil {
		t.Fatal("blank sufficiency must be refused")
	}
	if _, err := c.Record("context-map", map[string]any{"entries": []any{"a.go"}, "sufficiency": "ok"}, false, "agent"); err != nil {
		t.Fatalf("valid context-map refused: %v", err)
	}
	if _, err := c.Record("completeness", map[string]any{"items": []any{}}, false, "agent"); err == nil {
		t.Fatal("empty completeness must be refused")
	}
	if _, err := c.Record("completeness", map[string]any{"items": []any{map[string]any{"case": "empty"}}}, false, "agent"); err == nil {
		t.Fatal("dispositionless completeness item must be refused")
	}
	if _, err := c.Record("completeness", map[string]any{"items": []any{map[string]any{"case": "empty", "disposition": "covered:AC-1"}}}, false, "agent"); err != nil {
		t.Fatalf("valid completeness refused: %v", err)
	}
	// disposition vocabulary: free prose let dispositioned cases drop
	// silently between Frame and Build
	if _, err := c.Record("completeness", map[string]any{"items": []any{map[string]any{"case": "empty", "disposition": "will handle later"}}}, false, "agent"); err == nil {
		t.Fatal("free-prose completeness disposition must be refused")
	}
	if _, err := c.Record("completeness", map[string]any{"items": []any{map[string]any{"case": "empty", "disposition": "covered:"}}}, false, "agent"); err == nil {
		t.Fatal("covered: without an AC id must be refused")
	}
	for _, ok := range []string{"out-of-scope: separate service", "accepted-risk: bounded input"} {
		if _, err := c.Record("completeness", map[string]any{"items": []any{map[string]any{"case": "x", "disposition": ok}}}, false, "agent"); err != nil {
			t.Fatalf("disposition %q refused: %v", ok, err)
		}
	}
}

// attack-path records: one per adversary path, enum-validated, and
// adr-accepted must be backed by an ADR artifact record.
func TestAttackPathWriteValidation(t *testing.T) {
	c := newCtl(t)
	_, _ = c.RunStart("diff", "new")
	if _, err := c.Record("attack-path", map[string]any{"path": "  ", "feasibility": "high", "disposition": "open"}, false, "agent"); err == nil {
		t.Fatal("blank path must be refused")
	}
	if _, err := c.Record("attack-path", map[string]any{"path": "p", "feasibility": "certain", "disposition": "open"}, false, "agent"); err == nil {
		t.Fatal("unknown feasibility must be refused")
	}
	if _, err := c.Record("attack-path", map[string]any{"path": "p", "feasibility": "high", "disposition": "ignored"}, false, "agent"); err == nil {
		t.Fatal("unknown disposition must be refused")
	}
	// adr-accepted without any ADR record is a hollow claim
	if _, err := c.Record("attack-path", map[string]any{"path": "p", "feasibility": "high", "disposition": "adr-accepted"}, false, "agent"); err == nil {
		t.Fatal("adr-accepted without an ADR artifact must be refused")
	}
	ev, err := c.Record("attack-path", map[string]any{"path": "admin takeover ← forged cookie", "feasibility": "high", "disposition": "open"}, false, "agent")
	if err != nil {
		t.Fatalf("valid attack-path refused: %v", err)
	}
	if _, err := c.Record("artifact", map[string]any{"path": "docs/architecture/adr/001.md", "template": "adr", "status": "stub"}, false, "agent"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Record("attack-path", map[string]any{"updates": ev.ID, "disposition": "adr-accepted"}, false, "agent"); err != nil {
		t.Fatalf("adr-accepted with an ADR record refused: %v", err)
	}
}

// assumption status lifecycle enum.
func TestAssumptionStatusValidation(t *testing.T) {
	c := newCtl(t)
	_, _ = c.RunStart("diff", "new")
	if _, err := c.Record("assumption", map[string]any{"text": "x", "status": "maybe"}, false, "agent"); err == nil {
		t.Fatal("unknown assumption status must be refused")
	}
	ev, err := c.Record("assumption", map[string]any{"text": "prod db reachable", "status": "open", "high_risk": true}, false, "agent")
	if err != nil {
		t.Fatalf("valid assumption refused: %v", err)
	}
	if _, err := c.Record("assumption", map[string]any{"updates": ev.ID, "status": "validated"}, false, "agent"); err != nil {
		t.Fatalf("discharge update refused: %v", err)
	}
}

// 03 §4.3: a rejected option may never be re-proposed — now engine-checked.
func TestOptionSetRejectedCrossCheck(t *testing.T) {
	c := newCtl(t)
	_, _ = c.RunStart("diff", "new")
	set := func(stage, selected string, rejected ...string) map[string]any {
		var rej []any
		for _, r := range rejected {
			rej = append(rej, map[string]any{"id": r, "reason": "priced out"})
		}
		return map[string]any{"stage": stage, "candidates": []any{"a", "b", "c"},
			"selected": selected, "rejected": rej}
	}
	if _, err := c.Record("option-set", map[string]any{"stage": "system", "candidates": []any{"only-one"},
		"selected": "only-one", "rejected": []any{}}, false, "agent"); err == nil {
		t.Fatal("single-candidate option-set must be refused")
	}
	if _, err := c.Record("option-set", set("bogus-stage", "a", "b"), false, "agent"); err == nil {
		t.Fatal("unknown stage must be refused")
	}
	first, err := c.Record("option-set", set("system", "a", "b", "c"), false, "agent")
	if err != nil {
		t.Fatal(err)
	}
	// loop re-entry re-proposing a rejected option: refused
	if _, err := c.Record("option-set", set("system", "b"), false, "agent"); err == nil {
		t.Fatal("re-selecting a rejected option must be refused")
	}
	// a different stage may reuse the id space
	if _, err := c.Record("option-set", set("software", "b", "a"), false, "agent"); err != nil {
		t.Fatalf("other-stage selection wrongly blocked: %v", err)
	}
	// a disposition referencing the rejecting set is the recorded escape
	if _, err := c.Record("disposition", map[string]any{"ref": first.ID, "text": "constraint lifted: managed service now approved"}, false, "agent"); err != nil {
		t.Fatal(err)
	}
	if _, err := c.Record("option-set", set("system", "b", "a"), false, "agent"); err != nil {
		t.Fatalf("dispositioned rejection must release the option: %v", err)
	}
}

func TestFindingWriteValidation(t *testing.T) {
	c := newCtl(t)
	_, _ = c.RunStart("assessment", "code-review")
	if _, err := c.Record("finding", map[string]any{"fid": " ", "severity": "major", "text": "t"}, false, "agent"); err == nil {
		t.Fatal("blank fid must be refused")
	}
	if _, err := c.Record("finding", map[string]any{"fid": "F-1", "severity": "catastrophic", "text": "t"}, false, "agent"); err == nil {
		t.Fatal("unknown severity must be refused")
	}
	if _, err := c.Record("finding", map[string]any{"fid": "F-1", "severity": "major", "text": ""}, false, "agent"); err == nil {
		t.Fatal("empty text must be refused")
	}
	if _, err := c.Record("finding", map[string]any{"fid": "F-1", "severity": "major", "text": "auth bypass", "evidence": "poc"}, false, "agent"); err != nil {
		t.Fatalf("valid finding refused: %v", err)
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
	// n/a requires a reason (manual and auto alike)
	if _, err := c.Record("verdict", map[string]any{"agent": "ux-reviewer", "status": "n/a", "criticals": 0, "majors": 0}, false, "agent"); err == nil {
		t.Fatal("reasonless n/a must be refused")
	}
	if _, err := c.Record("verdict", map[string]any{"agent": "ux-reviewer", "status": "n/a", "criticals": 0, "majors": 0, "reason": "no UI in this diff"}, false, "agent"); err != nil {
		t.Fatalf("reasoned n/a refused: %v", err)
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
	rec("context-map", map[string]any{"entries": []any{"pkg/reader.go", "pkg/reader_test.go", "cmd/main.go (caller)"}, "sufficiency": "single file change"})
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

// Approval anchoring (04 §8.1) + the opt-in hardened dial (09 Q3).
func TestApproveAnchoring(t *testing.T) {
	c := newCtl(t)
	_, _ = c.RunStart("diff", "fix")

	// plain approval: self-attested, no anchor
	ev, err := c.Approve("frame", "p")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := ev.Data["answer_ref"]; ok {
		t.Fatal("no user-answer captured — no anchor expected")
	}

	// a hook-captured answer AFTER the last frame approval anchors the next one
	ua, err := c.Record("user-answer", map[string]any{"question": "approve scope?", "answer": "yes"}, true, "hook")
	if err != nil {
		t.Fatal(err)
	}
	ev, err = c.Approve("scope", "p")
	if err != nil {
		t.Fatal(err)
	}
	if ref, _ := ev.Data["answer_ref"].(string); ref != ua.ID {
		t.Fatalf("answer_ref = %v, want %s", ev.Data["answer_ref"], ua.ID)
	}

	// the same answer must NOT anchor a re-approval of a gate it precedes:
	// scope's last approval is now NEWER than the answer
	ev, err = c.Approve("scope", "p2")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := ev.Data["answer_ref"]; ok {
		t.Fatal("stale answer must not anchor a later re-approval")
	}
}

// Approvals bind engine-computed refs — what the "yes" was a yes to.
func TestApproveBindsRefs(t *testing.T) {
	c := newCtl(t)
	_, _ = c.RunStart("diff", "fix")
	rec := func(kind string, data map[string]any) {
		t.Helper()
		if _, err := c.Record(kind, data, false, "agent"); err != nil {
			t.Fatalf("%s: %v", kind, err)
		}
	}
	rec("requirement", map[string]any{"rid": "SWR-1", "level": "software", "text": "t", "status": "active",
		"acs": []any{map[string]any{"id": "AC-1", "text": "a", "verifiable": true}}})
	rec("assumption", map[string]any{"text": "prod db reachable", "status": "open", "high_risk": true})
	rec("assumption", map[string]any{"text": "low-risk detail", "status": "open", "high_risk": false})
	ev, err := c.Approve("scope", "presented")
	if err != nil {
		t.Fatal(err)
	}
	refs, _ := ev.Data["approved_refs"].([]any)
	if len(refs) != 2 {
		t.Fatalf("scope approval must bind the requirement + high-risk assumption, got %v", refs)
	}
	if h, _ := ev.Data["refs_hash"].(string); h == "" {
		t.Error("refs_hash missing")
	}
	// frame approvals bind nothing (no content baseline) — no refs field
	ev, err = c.Approve("frame", "p")
	if err != nil {
		t.Fatal(err)
	}
	if _, ok := ev.Data["approved_refs"]; ok {
		t.Error("frame approval must not carry refs")
	}
}

func TestApproveHardenedRefusesUnanchored(t *testing.T) {
	c := newCtl(t)
	c.Config.Flags = map[string]any{"approvals": "hardened"}
	_, _ = c.RunStart("diff", "fix")

	if _, err := c.Approve("frame", "p"); err == nil {
		t.Fatal("hardened mode must refuse an unanchored approval")
	}
	_, _ = c.Record("user-answer", map[string]any{"question": "approve frame?", "answer": "yes"}, true, "hook")
	ev, err := c.Approve("frame", "p")
	if err != nil {
		t.Fatalf("anchored approval must pass under hardened mode: %v", err)
	}
	if ref, _ := ev.Data["answer_ref"].(string); ref == "" {
		t.Fatal("hardened approval must carry the anchor")
	}
}
