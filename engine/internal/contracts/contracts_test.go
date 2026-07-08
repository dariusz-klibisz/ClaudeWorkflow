package contracts

import (
	"path/filepath"
	"testing"

	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/spec"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/store"
)

func loadSpec(t *testing.T) *spec.Spec {
	t.Helper()
	p, _ := filepath.Abs(filepath.Join("..", "..", "..", "workflow", "workflow.yaml"))
	s, err := spec.Load(p, "")
	if err != nil {
		t.Fatal(err)
	}
	return s
}

type envBuilder struct {
	events []store.Event
	n      int
}

func (b *envBuilder) add(kind string, auto bool, data map[string]any) string {
	b.n++
	id := store.NewULID()
	b.events = append(b.events, store.Event{
		ID: id, Seq: int64(b.n), Kind: kind, Auto: auto, Run: "r1", Data: data,
	})
	return id
}

func newEnv(t *testing.T, b *envBuilder, family, intent string, cfg *store.Config) *Env {
	if cfg == nil {
		cfg = &store.Config{}
	}
	return &Env{
		Spec:   loadSpec(t),
		Config: cfg,
		Run:    &store.Run{ID: "r1", Family: family, Intent: intent, Phase: "frame", Status: "active"},
		Events: b.events,
	}
}

func findingIDs(fs []Finding) map[string]Finding {
	m := map[string]Finding{}
	for _, f := range fs {
		m[f.ID] = f
	}
	return m
}

func TestEmptyFrameFindsEverything(t *testing.T) {
	env := newEnv(t, &envBuilder{}, "diff", "fix", nil)
	fs, err := EvaluatePhase(env, "frame")
	if err != nil {
		t.Fatal(err)
	}
	ids := findingIDs(fs)
	for _, want := range []string{"frame.classification", "frame.classification-approved", "frame.risk-scanned", "frame.requirements", "frame.origin"} {
		if _, ok := ids[want]; !ok {
			t.Errorf("missing finding %s", want)
		}
	}
	if f := ids["frame.classification-approved"]; !f.UserBlocked {
		t.Error("approval findings must be user-blocked")
	}
	if f := ids["frame.risk-scanned"]; f.UserBlocked {
		t.Error("record findings must not be user-blocked")
	}
}

func TestWhenIntents(t *testing.T) {
	// refactor intent: origin not required
	env := newEnv(t, &envBuilder{}, "diff", "refactor", nil)
	fs, _ := EvaluatePhase(env, "frame")
	if _, ok := findingIDs(fs)["frame.origin"]; ok {
		t.Error("frame.origin must not apply to refactor intent")
	}
}

func TestFamilyScoping(t *testing.T) {
	// assessment family: no requirements/completeness/abuse-case items at frame
	env := newEnv(t, &envBuilder{}, "assessment", "research", nil)
	fs, _ := EvaluatePhase(env, "frame")
	ids := findingIDs(fs)
	for _, absent := range []string{"frame.requirements", "frame.completeness", "frame.abuse-cases"} {
		if _, ok := ids[absent]; ok {
			t.Errorf("%s must not apply to assessment", absent)
		}
	}
}

func TestPerEachAmbiguityPerLens(t *testing.T) {
	b := &envBuilder{}
	b.add("risk", false, map[string]any{
		"signals": []any{"auth"},
		"lenses":  []any{"security", "user"},
	})
	b.add("ambiguity", false, map[string]any{"lens": "security", "disposition": "resolved"})
	env := newEnv(t, b, "diff", "new", nil)
	fs, _ := EvaluatePhase(env, "frame")
	f, ok := findingIDs(fs)["frame.ambiguity-per-lens"]
	if !ok {
		t.Fatal("user lens has no ambiguity — must be a finding")
	}
	if want := "user"; !contains(f.Detail, want) {
		t.Errorf("detail should name the missing lens: %q", f.Detail)
	}
	// satisfy the second lens
	b.add("ambiguity", false, map[string]any{"lens": "user", "none": true, "disposition": "none"})
	env = newEnv(t, b, "diff", "new", nil)
	fs, _ = EvaluatePhase(env, "frame")
	if _, ok := findingIDs(fs)["frame.ambiguity-per-lens"]; ok {
		t.Error("both lenses covered — finding must clear")
	}
}

func TestVerdictStickyAutoEvidence(t *testing.T) {
	b := &envBuilder{}
	autoID := b.add("verdict", true, map[string]any{"agent": "critic", "status": "unsafe", "criticals": 2, "majors": 0})
	// manual attempt to override the auto fail
	b.add("verdict", false, map[string]any{"agent": "critic", "status": "safe", "criticals": 0, "majors": 0})
	env := newEnv(t, b, "diff", "new", nil)
	env.Run.Phase = "design"
	fs, _ := EvaluatePhase(env, "design")
	if _, ok := findingIDs(fs)["design.critic"]; !ok {
		t.Fatal("manual verdict must not supersede auto fail (sticky evidence)")
	}
	// an explicit disposition of the auto verdict releases it
	b.add("disposition", false, map[string]any{"ref": autoID, "text": "false positive, reviewed with user"})
	env = newEnv(t, b, "diff", "new", nil)
	fs, _ = EvaluatePhase(env, "design")
	if _, ok := findingIDs(fs)["design.critic"]; ok {
		t.Error("dispositioned auto-fail must release the manual verdict")
	}
}

func TestVerdictRiskyNeedsDispositions(t *testing.T) {
	b := &envBuilder{}
	vid := b.add("verdict", true, map[string]any{"agent": "critic", "status": "risky", "criticals": 0, "majors": 1})
	env := newEnv(t, b, "diff", "new", nil)
	fs, _ := EvaluatePhase(env, "design")
	if _, ok := findingIDs(fs)["design.critic"]; !ok {
		t.Fatal("risky without dispositions must not pass")
	}
	b.add("disposition", false, map[string]any{"ref": vid, "text": "accepted: latency risk, monitored"})
	env = newEnv(t, b, "diff", "new", nil)
	fs, _ = EvaluatePhase(env, "design")
	if _, ok := findingIDs(fs)["design.critic"]; ok {
		t.Error("risky with dispositions must pass")
	}
}

func TestRedGreenPerTask(t *testing.T) {
	b := &envBuilder{}
	b.add("task", false, map[string]any{"tid": "T-1", "subject": "s", "dod": []any{"green test"}, "status": "done"})
	// green-only history must NOT satisfy test-first (the C4 fix)
	b.add("test-run", true, map[string]any{"cmd": "go test", "exit": 0, "grounded": true, "task": "T-1"})
	env := newEnv(t, b, "diff", "new", nil)
	fs, _ := EvaluatePhase(env, "build")
	if _, ok := findingIDs(fs)["build.test-first"]; !ok {
		t.Fatal("green-only history must not satisfy test-first")
	}
	// red then green satisfies
	b.add("test-run", true, map[string]any{"cmd": "go test", "exit": 1, "grounded": true, "task": "T-1"})
	b.add("test-run", true, map[string]any{"cmd": "go test", "exit": 0, "grounded": true, "task": "T-1"})
	env = newEnv(t, b, "diff", "new", nil)
	fs, _ = EvaluatePhase(env, "build")
	if f, ok := findingIDs(fs)["build.test-first"]; ok {
		t.Errorf("red→green must satisfy test-first: %s", f.Detail)
	}
	// a waiver satisfies a genuinely testless task
	b2 := &envBuilder{}
	b2.add("task", false, map[string]any{"tid": "T-2", "subject": "docs", "dod": []any{"d"}, "status": "done"})
	b2.add("waiver", false, map[string]any{"item": "T-2", "reason": "docs-only task"})
	env = newEnv(t, b2, "diff", "new", nil)
	fs, _ = EvaluatePhase(env, "build")
	if f, ok := findingIDs(fs)["build.test-first"]; ok {
		t.Errorf("waived testless task must pass: %s", f.Detail)
	}
}

func TestRedGreenOrdering(t *testing.T) {
	// isolate red-green via a task with red after green only
	b := &envBuilder{}
	b.add("task", false, map[string]any{"tid": "T-9", "subject": "s", "dod": []any{"d"}, "status": "done"})
	b.add("test-run", true, map[string]any{"cmd": "go test", "exit": 1, "grounded": true, "task": "T-9"})
	env := newEnv(t, b, "diff", "new", nil)
	ok, _, err := evalPredicate(env, spec.PredRedGreen, map[string]any{"link": "task"}, "T-9")
	if err != nil || ok {
		t.Fatal("red without green must fail")
	}
	b.add("test-run", true, map[string]any{"cmd": "go test", "exit": 0, "grounded": true, "task": "T-9"})
	env = newEnv(t, b, "diff", "new", nil)
	ok, _, _ = evalPredicate(env, spec.PredRedGreen, map[string]any{"link": "task"}, "T-9")
	if !ok {
		t.Error("red then green must pass")
	}
	// ungrounded green must not count
	b2 := &envBuilder{}
	b2.add("test-run", true, map[string]any{"cmd": "go test", "exit": 1, "grounded": true, "task": "T-9"})
	b2.add("test-run", false, map[string]any{"cmd": "go test | tail", "exit": 0, "grounded": false, "task": "T-9"})
	env = newEnv(t, b2, "diff", "new", nil)
	ok, _, _ = evalPredicate(env, spec.PredRedGreen, map[string]any{"link": "task"}, "T-9")
	if ok {
		t.Error("ungrounded green must not complete red-green")
	}
}

func TestNoOpenWithUpdateFolding(t *testing.T) {
	b := &envBuilder{}
	tid := b.add("task", false, map[string]any{"tid": "T-1", "subject": "s", "dod": []any{"d"}, "status": "open"})
	env := newEnv(t, b, "diff", "new", nil)
	fs, _ := EvaluatePhase(env, "build")
	if _, ok := findingIDs(fs)["build.tasks-closed"]; !ok {
		t.Fatal("open task must block build exit")
	}
	b.add("task", false, map[string]any{"updates": tid, "status": "done"})
	env = newEnv(t, b, "diff", "new", nil)
	fs, _ = EvaluatePhase(env, "build")
	if f, ok := findingIDs(fs)["build.tasks-closed"]; ok {
		t.Errorf("updated task still open: %s", f.Detail)
	}
}

func TestWhenConfigUX(t *testing.T) {
	env := newEnv(t, &envBuilder{}, "diff", "new", &store.Config{UX: false})
	fs, _ := EvaluatePhase(env, "design")
	if _, ok := findingIDs(fs)["design.ux-reviewed"]; ok {
		t.Error("ux item must not apply when ux=false")
	}
	env = newEnv(t, &envBuilder{}, "diff", "new", &store.Config{UX: true})
	fs, _ = EvaluatePhase(env, "design")
	if _, ok := findingIDs(fs)["design.ux-reviewed"]; !ok {
		t.Error("ux item must apply when ux=true")
	}
}

func TestWhenSignalsThreatModel(t *testing.T) {
	b := &envBuilder{}
	b.add("risk", false, map[string]any{"signals": []any{"docs"}, "lenses": []any{"user"}})
	env := newEnv(t, b, "diff", "new", nil)
	fs, _ := EvaluatePhase(env, "design")
	if _, ok := findingIDs(fs)["design.threat-model"]; ok {
		t.Error("threat model must not apply without security signals")
	}
	b2 := &envBuilder{}
	b2.add("risk", false, map[string]any{"signals": []any{"auth"}, "lenses": []any{"security"}})
	env = newEnv(t, b2, "diff", "new", nil)
	fs, _ = EvaluatePhase(env, "design")
	if _, ok := findingIDs(fs)["design.threat-model"]; !ok {
		t.Error("auth signal must require a threat model")
	}
}

func TestWaiverClearsWaivableItem(t *testing.T) {
	b := &envBuilder{}
	env := newEnv(t, b, "diff", "new", nil)
	fs, _ := EvaluatePhase(env, "context")
	if _, ok := findingIDs(fs)["context.research-grounded"]; !ok {
		t.Fatal("expected research finding")
	}
	b.add("waiver", false, map[string]any{"item": "context.research-grounded", "reason": "no external research needed"})
	env = newEnv(t, b, "diff", "new", nil)
	fs, _ = EvaluatePhase(env, "context")
	if _, ok := findingIDs(fs)["context.research-grounded"]; ok {
		t.Error("waived item must clear")
	}
	// non-waivable items cannot be waived
	b.add("waiver", false, map[string]any{"item": "context.scope-approved", "reason": "nope"})
	env = newEnv(t, b, "diff", "new", nil)
	fs, _ = EvaluatePhase(env, "context")
	if _, ok := findingIDs(fs)["context.scope-approved"]; !ok {
		t.Error("non-waivable approval must ignore waivers")
	}
}

func TestVerifyACVerdictsPerAC(t *testing.T) {
	b := &envBuilder{}
	b.add("requirement", false, map[string]any{
		"rid": "SWR-1", "level": "software", "text": "t", "status": "active",
		"acs": []any{
			map[string]any{"id": "AC-1", "text": "a", "verifiable": true},
			map[string]any{"id": "AC-2", "text": "b", "verifiable": true},
		},
	})
	b.add("ac-verdict", false, map[string]any{"ac": "AC-1", "status": "pass"})
	env := newEnv(t, b, "diff", "new", nil)
	fs, _ := EvaluatePhase(env, "verify")
	f, ok := findingIDs(fs)["verify.ac-verdicts"]
	if !ok {
		t.Fatal("AC-2 unverdicted — must be a finding")
	}
	if !contains(f.Detail, "AC-2") {
		t.Errorf("detail should name AC-2: %q", f.Detail)
	}
	// failing AC blocks separately
	b.add("ac-verdict", false, map[string]any{"ac": "AC-2", "status": "fail"})
	env = newEnv(t, b, "diff", "new", nil)
	fs, _ = EvaluatePhase(env, "verify")
	if _, ok := findingIDs(fs)["verify.no-failing-acs"]; !ok {
		t.Error("undispositioned failing AC must block")
	}
}

// The vacuous-pass fix: per-each items with min:1 fail when the element set
// is empty (an AC-less requirement used to dodge every per-AC item).
func TestPerEachMinGuardsVacuousPass(t *testing.T) {
	b := &envBuilder{}
	b.add("requirement", false, map[string]any{
		"rid": "SWR-1", "level": "software", "text": "t", "status": "active",
		"acs": []any{},
	})
	env := newEnv(t, b, "diff", "new", nil)
	fs, _ := EvaluatePhase(env, "verify")
	f, ok := findingIDs(fs)["verify.ac-verdicts"]
	if !ok {
		t.Fatal("zero ACs must fail verify.ac-verdicts (min 1), not pass vacuously")
	}
	if !contains(f.Detail, "min") {
		t.Errorf("detail should explain the min guard: %q", f.Detail)
	}
	fs, _ = EvaluatePhase(env, "plan")
	if _, ok := findingIDs(fs)["plan.verification-strategy"]; !ok {
		t.Error("zero ACs must fail plan.verification-strategy (min 1)")
	}
	// with one AC + linked records, both clear
	b2 := &envBuilder{}
	b2.add("requirement", false, map[string]any{
		"rid": "SWR-1", "level": "software", "text": "t", "status": "active",
		"acs": []any{map[string]any{"id": "AC-1", "text": "a", "verifiable": true}},
	})
	b2.add("verification-strategy", false, map[string]any{"ac": "AC-1", "method": "unit test", "command": "go test"})
	b2.add("ac-verdict", false, map[string]any{"ac": "AC-1", "status": "pass"})
	env = newEnv(t, b2, "diff", "new", nil)
	fs, _ = EvaluatePhase(env, "plan")
	if f, ok := findingIDs(fs)["plan.verification-strategy"]; ok {
		t.Errorf("one AC with a strategy must pass: %s", f.Detail)
	}
	fs, _ = EvaluatePhase(env, "verify")
	if f, ok := findingIDs(fs)["verify.ac-verdicts"]; ok {
		t.Errorf("one AC with a verdict must pass: %s", f.Detail)
	}
}

// The C10 fix: depth items demand ≥3 elements, waivable as the escape.
func TestMinContentDepthItems(t *testing.T) {
	b := &envBuilder{}
	b.add("context-map", false, map[string]any{"entries": []any{"a.go"}, "sufficiency": "ok"})
	env := newEnv(t, b, "diff", "new", nil)
	fs, _ := EvaluatePhase(env, "context")
	f, ok := findingIDs(fs)["context.map-depth"]
	if !ok {
		t.Fatal("one-entry map must fail the depth floor")
	}
	if !f.Waivable || !contains(f.Detail, "min") {
		t.Errorf("depth finding must be waivable and name the floor: %+v", f)
	}
	if _, ok := findingIDs(fs)["context.map"]; ok {
		t.Error("existence item must still pass — only depth fails")
	}
	// three entries clear it
	b.add("context-map", false, map[string]any{"entries": []any{"b.go", "b_test.go"}, "sufficiency": "ok"})
	env = newEnv(t, b, "diff", "new", nil)
	fs, _ = EvaluatePhase(env, "context")
	if _, ok := findingIDs(fs)["context.map-depth"]; ok {
		t.Error("3 entries across maps must satisfy the floor")
	}
	// completeness depth mirrors it
	b2 := &envBuilder{}
	b2.add("completeness", false, map[string]any{"items": []any{map[string]any{"case": "empty", "disposition": "ok"}}})
	env = newEnv(t, b2, "diff", "new", nil)
	fs, _ = EvaluatePhase(env, "frame")
	if _, ok := findingIDs(fs)["frame.completeness-depth"]; !ok {
		t.Error("one-case completeness must fail the depth floor")
	}
	b2.add("waiver", false, map[string]any{"item": "frame.completeness-depth", "reason": "single toggle flip, no negative space"})
	env = newEnv(t, b2, "diff", "new", nil)
	fs, _ = EvaluatePhase(env, "frame")
	if _, ok := findingIDs(fs)["frame.completeness-depth"]; ok {
		t.Error("waiver must clear the depth item")
	}
}

// Red→green pairing (the order-only hole): cross-runner pairs never satisfy
// the predicate; same-runner selector mismatches pass weakly and are
// surfaced via WeakRedGreenTasks.
func TestRedGreenPairing(t *testing.T) {
	// cross-runner: gitleaks red + pytest green must NOT pass
	b := &envBuilder{}
	b.add("task", false, map[string]any{"tid": "T-1", "subject": "s", "dod": []any{"d"}, "status": "done"})
	b.add("test-run", true, map[string]any{"cmd": "gitleaks detect", "exit": 1, "grounded": true, "task": "T-1"})
	b.add("test-run", true, map[string]any{"cmd": "pytest", "exit": 0, "grounded": true, "task": "T-1"})
	env := newEnv(t, b, "diff", "new", nil)
	ok, detail, _ := evalPredicate(env, spec.PredRedGreen, map[string]any{"link": "task"}, "T-1")
	if ok {
		t.Fatal("cross-runner red→green must not pass")
	}
	if !contains(detail, "runner") {
		t.Errorf("detail should explain the runner mismatch: %q", detail)
	}
	// red on a selector, green on the full suite: strict
	b2 := &envBuilder{}
	b2.add("task", false, map[string]any{"tid": "T-2", "subject": "s", "dod": []any{"d"}, "status": "done"})
	b2.add("test-run", true, map[string]any{"cmd": "pytest tests/test_x.py::test_new", "exit": 1, "grounded": true, "task": "T-2"})
	b2.add("test-run", true, map[string]any{"cmd": "pytest", "exit": 0, "grounded": true, "task": "T-2"})
	env = newEnv(t, b2, "diff", "new", nil)
	if ok, _, _ := evalPredicate(env, spec.PredRedGreen, map[string]any{"link": "task"}, "T-2"); !ok {
		t.Error("selector red + full-suite green must pass strictly")
	}
	if weak := WeakRedGreenTasks(env); len(weak) != 0 {
		t.Errorf("strict pair must not be reported weak: %v", weak)
	}
	// same runner, diverging selectors: weak pass, surfaced
	b3 := &envBuilder{}
	b3.add("task", false, map[string]any{"tid": "T-3", "subject": "s", "dod": []any{"d"}, "status": "done"})
	b3.add("test-run", true, map[string]any{"cmd": "pytest tests/test_x.py", "exit": 1, "grounded": true, "task": "T-3"})
	b3.add("test-run", true, map[string]any{"cmd": "pytest tests/test_y.py", "exit": 0, "grounded": true, "task": "T-3"})
	env = newEnv(t, b3, "diff", "new", nil)
	if ok, _, _ := evalPredicate(env, spec.PredRedGreen, map[string]any{"link": "task"}, "T-3"); !ok {
		t.Error("same-runner selector mismatch must weak-pass (no wedge)")
	}
	weak := WeakRedGreenTasks(env)
	if len(weak) != 1 || weak[0] != "T-3" {
		t.Errorf("weak pair must be surfaced: %v", weak)
	}
}

// A manual verdict for a gating agent must not satisfy the contract while
// verdict capture is provably alive; a disposition is the recorded escape.
func TestManualGatingVerdictRefusedWhileCaptureAlive(t *testing.T) {
	// hooks dead (no auto verdict anywhere): manual records keep working
	b := &envBuilder{}
	b.add("verdict", false, map[string]any{"agent": "critic", "status": "safe", "criticals": 0, "majors": 0})
	env := newEnv(t, b, "diff", "new", nil)
	fs, _ := EvaluatePhase(env, "design")
	if _, ok := findingIDs(fs)["design.critic"]; ok {
		t.Fatal("degraded-hooks session: manual verdict must still satisfy")
	}
	// capture alive (an auto verdict exists): the manual one is refused
	b2 := &envBuilder{}
	b2.add("verdict", true, map[string]any{"agent": "design-reviewer", "status": "clean", "criticals": 0, "majors": 0})
	manualID := b2.add("verdict", false, map[string]any{"agent": "critic", "status": "safe", "criticals": 0, "majors": 0})
	env = newEnv(t, b2, "diff", "new", nil)
	fs, _ = EvaluatePhase(env, "design")
	f, ok := findingIDs(fs)["design.critic"]
	if !ok {
		t.Fatal("manual gating verdict must be refused while capture is alive")
	}
	if !contains(f.Detail, "self-attested") {
		t.Errorf("detail should explain: %q", f.Detail)
	}
	// disposition referencing the manual verdict is the escape
	b2.add("disposition", false, map[string]any{"ref": manualID, "text": "critic ran outside wf: namespace; transcript reviewed with user"})
	env = newEnv(t, b2, "diff", "new", nil)
	fs, _ = EvaluatePhase(env, "design")
	if _, ok := findingIDs(fs)["design.critic"]; ok {
		t.Error("dispositioned manual verdict must satisfy")
	}
	// non-gating agents (researcher) are unaffected
	b3 := &envBuilder{}
	b3.add("verdict", true, map[string]any{"agent": "design-reviewer", "status": "clean", "criticals": 0, "majors": 0})
	b3.add("verdict", false, map[string]any{"agent": "researcher", "status": "n/a", "criticals": 0, "majors": 0, "reason": "no research ran"})
	env = newEnv(t, b3, "diff", "new", nil)
	fs, _ = EvaluatePhase(env, "context")
	if _, ok := findingIDs(fs)["context.research-grounded"]; ok {
		t.Error("non-gating manual verdict must remain acceptable")
	}
}

func contains(s, sub string) bool {
	return len(s) >= len(sub) && (s == sub || len(sub) == 0 || indexOf(s, sub) >= 0)
}

func indexOf(s, sub string) int {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return i
		}
	}
	return -1
}
