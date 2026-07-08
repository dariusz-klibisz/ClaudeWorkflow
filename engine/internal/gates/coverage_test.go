package gates

import (
	"encoding/json"
	"testing"

	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/store"
)

func TestCoverageFromResponse(t *testing.T) {
	cases := []struct {
		name, stdout string
		want         float64
		ok           bool
	}{
		{"go", "ok  \tpkg\t0.5s\tcoverage: 85.2% of statements\n", 85.2, true},
		{"pytest-cov", "----------\nTOTAL                 250     25    90%\n", 90, true},
		{"jest", "All files          |   87.5 |    80 |   85 |   87.5 |\n", 87.5, true},
		{"tarpaulin", "|| Tested/Total Lines:\n72.34% coverage, 34/47 lines covered\n", 72.34, true},
		{"lcov", "  lines......: 91.3% (2000 of 2190 lines)\n", 91.3, true},
		{"none", "all tests passed\n", 0, false},
		{"absurd", "coverage: 200.0% of statements\n", 0, false},
	}
	for _, tc := range cases {
		raw, _ := json.Marshal(map[string]any{"stdout": tc.stdout, "stderr": ""})
		got, ok := coverageFromResponse(raw)
		if ok != tc.ok || (ok && got != tc.want) {
			t.Errorf("%s: got %v/%v want %v/%v", tc.name, got, ok, tc.want, tc.ok)
		}
	}
}

func greenRunPayload(t *testing.T, cmd, stdout string) string {
	t.Helper()
	raw, _ := json.Marshal(map[string]any{
		"session_id": "s1", "hook_event_name": "PostToolUse", "tool_name": "Bash",
		"tool_input":    map[string]any{"command": cmd},
		"tool_response": map[string]any{"stdout": stdout, "stderr": "", "interrupted": false},
	})
	return string(raw)
}

func TestCaptureTestRecordsGroundedCoverage(t *testing.T) {
	c := newCtl(t)
	c.Config = &store.Config{Thresholds: map[string]any{"coverage": 80}}
	run, _ := c.RunStart("diff", "fix")
	run.Phase = "build"
	_ = c.Store.SaveRun(run)

	// below the floor
	_ = CaptureTest(c, hookInput(t, greenRunPayload(t, "go test ./... -cover", "coverage: 61.0% of statements\n")))
	env, _ := c.Env(run)
	ms := env.Records("metric")
	if len(ms) != 1 {
		t.Fatalf("want 1 coverage metric, got %d", len(ms))
	}
	if g, _ := ms[0].Data["grounded"].(bool); !g {
		t.Error("scraped coverage must be grounded")
	}
	if bt, _ := ms[0].Data["below_threshold"].(bool); !bt {
		t.Errorf("61 < 80 must set below_threshold: %v", ms[0].Data)
	}

	// re-measure above the floor UPDATES the same record (no stale block)
	_ = CaptureTest(c, hookInput(t, greenRunPayload(t, "go test ./... -cover", "coverage: 85.0% of statements\n")))
	env, _ = c.Env(run)
	ms = env.Records("metric")
	if len(ms) != 1 {
		t.Fatalf("re-measure must fold into one effective record, got %d", len(ms))
	}
	if bt, _ := ms[0].Data["below_threshold"].(bool); bt {
		t.Errorf("85 >= 80 must clear below_threshold: %v", ms[0].Data)
	}
}

func TestManualMetricSelfAttested(t *testing.T) {
	c := newCtl(t)
	c.Config = &store.Config{Thresholds: map[string]any{"coverage": 80}}
	run, _ := c.RunStart("diff", "fix")
	ev, err := c.Record("metric", map[string]any{"name": "coverage", "value": 90.0}, false, "agent")
	if err != nil {
		t.Fatal(err)
	}
	if g, _ := ev.Data["grounded"].(bool); g {
		t.Error("manual metric must default to grounded=false")
	}
	if bt, ok := ev.Data["below_threshold"].(bool); !ok || bt {
		t.Errorf("below_threshold must be engine-computed false: %v", ev.Data)
	}
	_ = run
}

func TestCaptureTestLeadingCdAllowed(t *testing.T) {
	c := newCtl(t)
	run, _ := c.RunStart("diff", "fix")
	run.Phase = "build"
	_ = c.Store.SaveRun(run)
	_ = CaptureTest(c, hookInput(t, greenRunPayload(t, "cd services/api && pytest tests/", "5 passed\n")))
	env, _ := c.Env(run)
	trs := env.Records("test-run")
	if len(trs) != 1 {
		t.Fatalf("cd-prefixed runner must be captured, got %d records", len(trs))
	}
	if g, _ := trs[0].Data["grounded"].(bool); !g {
		t.Errorf("single leading cd must stay grounded: %v", trs[0].Data)
	}
	// deeper chains are not even recognized (head stays cd-prefixed) —
	// never a grounded record either way
	_ = CaptureTest(c, hookInput(t, greenRunPayload(t, "cd a && cd b && pytest", "ok\n")))
	env, _ = c.Env(run)
	trs = env.Records("test-run")
	if len(trs) != 1 {
		t.Fatalf("double chain must not be captured: got %d records", len(trs))
	}
	// runner followed by a further chain: recognized but ungrounded
	_ = CaptureTest(c, hookInput(t, greenRunPayload(t, "cd a && pytest && echo done", "ok\n")))
	env, _ = c.Env(run)
	trs = env.Records("test-run")
	if len(trs) != 2 {
		t.Fatalf("runner+chain must be captured ungrounded: got %d records", len(trs))
	}
	if g, _ := trs[1].Data["grounded"].(bool); g {
		t.Error("trailing chain must stay ungrounded")
	}
}
