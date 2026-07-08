package runctl

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

const authoredDoc = `# Design: the thing

## Context

A realistic amount of authored prose describing the selected approach, its
constraints, and the reasons the alternatives lost. Long enough to clear the
stub heuristic by an honest margin, because this is a real document.

## Decision

The boring option. Operational familiarity beats marginal headroom.

## Consequences

Migration tooling; rejected options stay visibly rejected.
`

func writeProjectDoc(t *testing.T, c *Ctl, rel, content string) {
	t.Helper()
	abs := filepath.Join(c.ProjectDir(), filepath.FromSlash(rel))
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(abs, []byte(content), 0o644); err != nil {
		t.Fatal(err)
	}
}

// Write-time: `status: present` is a claim the disk must confirm.
func TestArtifactPresentWriteTimeValidation(t *testing.T) {
	c := newCtl(t)
	_, _ = c.RunStart("artifact", "arch-design")

	// present without a file → refused
	_, err := c.Record("artifact", map[string]any{"path": "docs/design/x.md", "status": "present"}, false, "agent")
	if err == nil || !strings.Contains(err.Error(), "does not exist on disk") {
		t.Fatalf("present without file must be refused: %v", err)
	}

	// stub content → refused
	writeProjectDoc(t, c, "docs/design/x.md", "# x\nTODO\n")
	_, err = c.Record("artifact", map[string]any{"path": "docs/design/x.md", "status": "present"}, false, "agent")
	if err == nil || !strings.Contains(err.Error(), "stub") {
		t.Fatalf("stub content must be refused: %v", err)
	}

	// authored → accepted
	writeProjectDoc(t, c, "docs/design/x.md", authoredDoc)
	if _, err = c.Record("artifact", map[string]any{"path": "docs/design/x.md", "status": "present"}, false, "agent"); err != nil {
		t.Fatalf("authored artifact must be accepted: %v", err)
	}

	// stub record stays recordable (that's what wf doc new writes)
	if _, err = c.Record("artifact", map[string]any{"path": "docs/design/y.md", "status": "stub", "template": "design"}, true, "engine"); err != nil {
		t.Fatalf("stub record must stay recordable: %v", err)
	}
}

// Write-time: flipping a stub via updates= resolves path/template from the
// original record.
func TestArtifactPresentUpdateFlip(t *testing.T) {
	c := newCtl(t)
	_, _ = c.RunStart("artifact", "arch-design")
	ev, err := c.Record("artifact", map[string]any{"path": "docs/design/z.md", "status": "stub", "template": "design"}, true, "engine")
	if err != nil {
		t.Fatal(err)
	}
	// flip before authoring → refused
	if _, err := c.Record("artifact", map[string]any{"updates": ev.ID, "status": "present"}, false, "agent"); err == nil {
		t.Fatal("flip to present without the file must be refused")
	}
	writeProjectDoc(t, c, "docs/design/z.md", authoredDoc)
	if _, err := c.Record("artifact", map[string]any{"updates": ev.ID, "status": "present"}, false, "agent"); err != nil {
		t.Fatalf("flip after authoring must pass: %v", err)
	}
}

// Transition: the next phase's entry contract blocks `wf phase exit` when
// the inputs were never produced (force-exited earlier phases).
func TestPhaseExitBlocksOnNextEntryInputs(t *testing.T) {
	c := newCtl(t)
	_, _ = c.RunStart("artifact", "arch-design")
	// force out of frame: no classification, no requirements (entry skipped)
	if _, _, err := c.PhaseExit(true, "test landing"); err != nil {
		t.Fatal(err)
	}
	r, _ := c.Store.LoadRun()
	if r.Phase != "context" {
		t.Fatalf("force must still advance: %s", r.Phase)
	}
	// satisfy context's EXIT contract without producing requirements
	rec := func(kind string, data map[string]any) {
		t.Helper()
		if _, err := c.Record(kind, data, false, "agent"); err != nil {
			t.Fatalf("%s: %v", kind, err)
		}
	}
	rec("context-map", map[string]any{"entries": []any{"docs/"}, "sufficiency": "ok"})
	rec("reclassify", map[string]any{"result": "confirmed"})
	rec("verdict", map[string]any{"agent": "researcher", "status": "n/a", "criticals": 0, "majors": 0})
	if _, err := c.Approve("scope", ""); err != nil {
		t.Fatal(err)
	}
	findings, _, err := c.PhaseExit(false, "")
	if err != nil {
		t.Fatal(err)
	}
	found := false
	for _, f := range findings {
		if f.ID == "design.entry-requirements" {
			found = true
		}
		if strings.HasPrefix(f.ID, "context.") {
			t.Fatalf("context exit items should be met, got %s", f.ID)
		}
	}
	if !found {
		t.Fatalf("design entry input must block the transition: %v", findings)
	}
	// produce the input → transition proceeds
	rec("requirement", map[string]any{"rid": "SWR-1", "level": "software", "text": "t", "status": "active",
		"acs": []any{map[string]any{"id": "AC-1", "text": "a", "verifiable": true}}})
	findings, msg, err := c.PhaseExit(false, "")
	if err != nil || len(findings) > 0 {
		t.Fatalf("transition must proceed once inputs exist: %v %v", findings, err)
	}
	if !strings.Contains(msg, "design") {
		t.Errorf("should have entered design: %s", msg)
	}
}
