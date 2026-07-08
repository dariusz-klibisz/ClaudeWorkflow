package runctl

import (
	"path/filepath"
	"testing"

	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/spec"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/store"
)

func riskCtl(t *testing.T) *Ctl {
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

func TestRiskScanWordBoundaries(t *testing.T) {
	c := riskCtl(t)
	_, _ = c.RunStart("diff", "new")

	// the live bug: "requires/requirements" contains "ui" as a substring
	signals, lenses, err := c.RiskScan("create an application that requires computing prime numbers per the requirements", nil)
	if err != nil {
		t.Fatal(err)
	}
	for _, s := range signals {
		if s == "ui" {
			t.Fatalf("'requirements' must not trigger the ui signal: %v", signals)
		}
	}
	for _, l := range lenses {
		if l == "usability" {
			t.Fatalf("no usability lens without a real ui signal: %v", lenses)
		}
	}
	if len(lenses) == 0 || lenses[len(lenses)-1] != "user" && lenses[0] != "user" {
		t.Errorf("the user lens is always selected: %v", lenses)
	}
}

func TestRiskScanRealSignals(t *testing.T) {
	c := riskCtl(t)
	_, _ = c.RunStart("diff", "new")
	cases := []struct {
		text string
		want string
	}{
		{"add a login form with password reset", "auth"},
		{"build the settings UI page", "ui"},
		{"protect against rm -rf in cleanup scripts", "destructive"},
		{"parse user-provided upload files", "boundary"},
		{"fix a race in the concurrent worker pool", "concurrency"},
	}
	for _, tc := range cases {
		signals, _, err := c.RiskScan(tc.text, nil)
		if err != nil {
			t.Fatal(err)
		}
		found := false
		for _, s := range signals {
			if s == tc.want {
				found = true
			}
		}
		if !found {
			t.Errorf("%q: want signal %s, got %v", tc.text, tc.want, signals)
		}
	}
}

func TestRiskScanStandingFlags(t *testing.T) {
	c := riskCtl(t)
	c.Config = &store.Config{Flags: map[string]any{"pii": true, "internet_facing": true}}
	_, _ = c.RunStart("diff", "new")
	signals, lenses, err := c.RiskScan("compute prime numbers", nil)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]bool{"data": false, "network": false, "boundary": false}
	for _, s := range signals {
		if _, ok := want[s]; ok {
			want[s] = true
		}
	}
	for k, seen := range want {
		if !seen {
			t.Errorf("standing flags must arm signal %q: %v", k, signals)
		}
	}
	hasLens := func(l string) bool {
		for _, x := range lenses {
			if x == l {
				return true
			}
		}
		return false
	}
	if !hasLens("security") || !hasLens("adversarial") {
		t.Errorf("standing signals must bind their lenses: %v", lenses)
	}
}

func TestRiskScanFlagsOffAddNothing(t *testing.T) {
	c := riskCtl(t)
	_, _ = c.RunStart("diff", "new")
	signals, _, err := c.RiskScan("compute prime numbers", nil)
	if err != nil {
		t.Fatal(err)
	}
	if len(signals) != 0 {
		t.Errorf("no flags, benign text: want no signals, got %v", signals)
	}
}
