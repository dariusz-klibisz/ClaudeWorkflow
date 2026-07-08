// Package ops implements the agent-facing engine operations that ground
// records in reality: dependency checking, defect-origin discovery, and
// engine-mediated document creation (03 §4, 06 §6).
package ops

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/runctl"
)

// ---------------------------------------------------------------------------
// wf deps check
// ---------------------------------------------------------------------------

var manifestFiles = []string{
	"go.mod", "package.json", "pyproject.toml", "requirements.txt",
	"Cargo.toml", "pom.xml", "build.gradle", "Gemfile", "composer.json",
	"*.csproj",
}

// DepsCheck verifies (best-effort) that the project's declared dependency
// world and the plan's verification tooling exist, then records the deps
// verdict: present | missing | n/a. Missing blocks the Plan gate.
func DepsCheck(c *runctl.Ctl, projectDir string) (string, error) {
	r, err := c.MustRun()
	if err != nil {
		return "", err
	}
	var manifests []string
	for _, pat := range manifestFiles {
		hits, _ := filepath.Glob(filepath.Join(projectDir, pat))
		for _, h := range hits {
			manifests = append(manifests, filepath.Base(h))
		}
	}

	// every verification-strategy command's tool must resolve on PATH
	env, err := c.Env(r)
	if err != nil {
		return "", err
	}
	var missing []string
	checked := 0
	for _, vs := range env.Records("verification-strategy") {
		cmd, _ := vs.Data["command"].(string)
		if cmd == "" {
			continue
		}
		checked++
		tool := strings.Fields(cmd)[0]
		if _, err := exec.LookPath(tool); err != nil {
			missing = append(missing, fmt.Sprintf("%s (verification for %v)", tool, vs.Data["ac"]))
		}
	}

	verdict := "present"
	detail := fmt.Sprintf("manifests: %s; %d verification tool(s) resolved", orNone(manifests), checked)
	switch {
	case len(missing) > 0:
		verdict = "missing"
		detail = "missing tools: " + strings.Join(missing, ", ")
	case len(manifests) == 0 && checked == 0:
		verdict = "n/a"
		detail = "no dependency manifest and no verification commands to check"
	}
	if _, err := c.Record("deps", map[string]any{"verdict": verdict, "detail": detail}, true, "engine"); err != nil {
		return "", err
	}
	return fmt.Sprintf("deps: %s — %s", verdict, detail), nil
}

// ---------------------------------------------------------------------------
// wf origin discover
// ---------------------------------------------------------------------------

// OriginDiscover attributes a defect to its origin commit via git (pickaxe
// on --text and/or history of --path). Records the origin on success;
// returns an error with guidance when inconclusive (the agent falls back to
// a manual record with reduced confidence).
func OriginDiscover(c *runctl.Ctl, projectDir, path, text string) (string, error) {
	if _, err := c.MustRun(); err != nil {
		return "", err
	}
	if path == "" && text == "" {
		return "", fmt.Errorf("origin discover needs --path and/or --text")
	}
	args := []string{"-C", projectDir, "log", "-n", "3", "--format=%h %ad %an: %s", "--date=short"}
	if text != "" {
		args = append(args, "-S", text)
	}
	if path != "" {
		args = append(args, "--follow", "--", path)
	}
	out, err := exec.Command("git", args...).Output()
	lines := strings.TrimSpace(string(out))
	if err != nil || lines == "" {
		return "", fmt.Errorf("git found no candidate commits (path=%q text=%q) — record manually with reduced confidence: wf record origin attribution=\"…\"", path, text)
	}
	first := strings.SplitN(lines, "\n", 2)[0]
	commit := strings.Fields(first)[0]
	attribution := fmt.Sprintf("candidate origin (git pickaxe/history): %s", first)
	// the ledger beats git heuristics: a commit-origin record (this run or an
	// archived one) names the run that produced the commit — survives
	// squash/rebase where git history rewrites lose the trail
	if run := commitOriginRun(c, commit); run != "" {
		attribution += fmt.Sprintf(" — originating run %s (ledger commit-origin)", run)
	}
	if _, err := c.Record("origin", map[string]any{"commit": commit, "attribution": attribution}, true, "engine"); err != nil {
		return "", err
	}
	return "origin recorded: " + attribution + "\ncandidates:\n" + lines, nil
}

// commitOriginRun searches live + archived commit-origin records for a
// commit sha (prefix-tolerant both ways: git output is abbreviated).
func commitOriginRun(c *runctl.Ctl, sha string) string {
	if sha == "" {
		return ""
	}
	evs, err := c.Store.AllEvents()
	if err != nil {
		return ""
	}
	for _, e := range evs {
		if e.Kind != "commit-origin" {
			continue
		}
		full, _ := e.Data["commit"].(string)
		if full == "" {
			continue
		}
		if strings.HasPrefix(full, sha) || strings.HasPrefix(sha, full) {
			if run, _ := e.Data["run"].(string); run != "" {
				return run
			}
			return e.Run
		}
	}
	return ""
}

// ---------------------------------------------------------------------------
// wf doc new
// ---------------------------------------------------------------------------

type docSpec struct {
	dest string // pattern with <slug>; <nnnn> = next ADR number
	role string
}

var docTypes = map[string]docSpec{
	"adr":               {dest: "docs/architecture/adr/<nnnn>-<slug>.md"},
	"design":            {dest: "docs/design/<slug>.md"},
	"threat-model":      {dest: "docs/design/threat-model-<slug>.md"},
	"abuse-cases":       {dest: "docs/design/abuse-cases-<slug>.md"},
	"attack-tree":       {dest: "docs/design/attack-tree-<slug>.md"},
	"ux":                {dest: "docs/design/ux-<slug>.md"},
	"review":            {dest: "docs/reviews/<slug>.md", role: "deliverable-report"},
	"red-team-report":   {dest: "docs/reviews/red-team-<slug>.md"},
	"test-plan":         {dest: "docs/test/<slug>.md"},
	"runbook":           {dest: "docs/operations/runbook-<slug>.md"},
	"incident":          {dest: "docs/incidents/<slug>.md"},
	"release-notes":     {dest: "docs/releases/<slug>.md"},
	"delivery-manifest": {dest: "docs/releases/delivery-manifest-<slug>.md", role: "delivery-manifest"},
	"retro":             {dest: "docs/retros/<slug>.md"},
	"research-findings": {
		dest: "docs/research/<slug>.md", role: "deliverable-report"},
	"investigation-findings": {
		dest: "docs/investigations/<slug>.md", role: "deliverable-report"},
}

var slugRe = regexp.MustCompile(`[^a-z0-9]+`)

// DocNew copies a plugin template into its docs/ destination and records the
// artifact as a stub — the engine-mediated copy makes "authored the template
// in place" impossible (06 §6). The agent flips status to `present` after
// authoring: wf record artifact updates=<id> status=present.
func DocNew(c *runctl.Ctl, pluginRoot, projectDir, docType, slug string) (string, error) {
	if _, err := c.MustRun(); err != nil {
		return "", err
	}
	ds, ok := docTypes[docType]
	if !ok {
		keys := make([]string, 0, len(docTypes))
		for k := range docTypes {
			keys = append(keys, k)
		}
		return "", fmt.Errorf("unknown doc type %q (one of: %s)", docType, strings.Join(keys, ", "))
	}
	slug = strings.Trim(slugRe.ReplaceAllString(strings.ToLower(slug), "-"), "-")
	if slug == "" {
		return "", fmt.Errorf("doc new needs --slug")
	}
	tmpl := filepath.Join(pluginRoot, "templates", docType+".md")
	content, err := os.ReadFile(tmpl)
	if err != nil {
		return "", fmt.Errorf("template %s missing in plugin: %w", docType, err)
	}
	rel := strings.ReplaceAll(ds.dest, "<slug>", slug)
	if strings.Contains(rel, "<nnnn>") {
		n := 1
		if hits, _ := filepath.Glob(filepath.Join(projectDir, "docs", "architecture", "adr", "*.md")); hits != nil {
			n = len(hits) + 1
		}
		rel = strings.ReplaceAll(rel, "<nnnn>", fmt.Sprintf("%04d", n))
	}
	abs := filepath.Join(projectDir, filepath.FromSlash(rel))
	if _, err := os.Stat(abs); err == nil {
		return "", fmt.Errorf("%s already exists — pick another slug", rel)
	}
	if err := os.MkdirAll(filepath.Dir(abs), 0o755); err != nil {
		return "", err
	}
	if err := os.WriteFile(abs, content, 0o644); err != nil {
		return "", err
	}
	data := map[string]any{"path": rel, "status": "stub", "template": docType}
	if ds.role != "" {
		data["role"] = ds.role
	}
	ev, err := c.Record("artifact", data, true, "engine")
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("created %s (record %s) — author it, then: wf record artifact updates=%s status=present", rel, ev.ID, ev.ID), nil
}

func orNone(xs []string) string {
	if len(xs) == 0 {
		return "(none)"
	}
	return strings.Join(xs, ", ")
}
