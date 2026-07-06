// gen emits the spec-derived views (07 §4): hooks/hooks.json (matchers from
// the gating roster), workflow/schemas/*.json (record-kind schemas), and
// skeleton agent files for roster entries that lack one (never overwrites).
// CI fails when generated files drift from the spec.
//
// Usage: go run ./gen [-check] (from engine/; repo root = ../)
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/spec"
)

const engineCmd = "${CLAUDE_PLUGIN_DATA}/bin/wf"

func main() {
	check := flag.Bool("check", false, "verify generated files are current (CI drift check)")
	flag.Parse()

	root, err := filepath.Abs("..")
	if err != nil {
		fatal(err)
	}
	sp, err := spec.Load(filepath.Join(root, "workflow", "workflow.yaml"), "")
	if err != nil {
		fatal(err)
	}

	drift := false
	drift = writeOrCheck(*check, filepath.Join(root, "hooks", "hooks.json"), hooksJSON(sp)) || drift
	for kind, schema := range recordSchemas(sp) {
		p := filepath.Join(root, "workflow", "schemas", kind+".json")
		drift = writeOrCheck(*check, p, schema) || drift
	}
	if err := agentSkeletons(*check, root, sp); err != nil {
		fatal(err)
	}
	if *check && drift {
		fatal(fmt.Errorf("generated files drift from workflow.yaml — run `go generate ./...` (go run ./gen)"))
	}
}

// ---------------------------------------------------------------------------
// hooks.json
// ---------------------------------------------------------------------------

func hooksJSON(sp *spec.Spec) []byte {
	gating := make([]string, 0)
	for _, a := range sp.GatingAgents() {
		gating = append(gating, a.Name)
	}
	agentMatcher := "^wf:(" + strings.Join(gating, "|") + ")$"

	exec := func(timeout int, args ...string) map[string]any {
		return map[string]any{
			"type": "command", "command": engineCmd, "args": args, "timeout": timeout,
		}
	}
	group := func(matcher string, handlers ...map[string]any) map[string]any {
		g := map[string]any{"hooks": handlers}
		if matcher != "" {
			g["matcher"] = matcher
		}
		return g
	}

	// Bootstrap ships as a sh script only: hooks cannot be platform-scoped,
	// and a `shell:"powershell"` entry errors visibly on hosts without
	// PowerShell. Native-Windows bootstrap is deferred to M5 (run
	// scripts/bootstrap.ps1 once by hand until then — see skills/init).
	bootstrapSh := map[string]any{
		"type": "command", "timeout": 60,
		"command": `"${CLAUDE_PLUGIN_ROOT}/scripts/bootstrap.sh"`,
	}

	doc := map[string]any{
		"description": "wf enforcement spine: Stop/task/verdict gates, tool gates, capture, and context injection (all generated from workflow/workflow.yaml)",
		"hooks": map[string]any{
			"SessionStart": []any{
				group("", bootstrapSh, exec(15, "inject", "session")),
			},
			"UserPromptSubmit": []any{
				group("", exec(10, "inject", "turn")),
			},
			"Stop": []any{
				group("", exec(15, "gate", "stop")),
			},
			"SubagentStart": []any{
				group(agentMatcher, exec(10, "inject", "agent")),
			},
			"SubagentStop": []any{
				group(agentMatcher, exec(15, "gate", "verdict")),
			},
			"TaskCreated": []any{
				group("", exec(10, "gate", "task-create")),
			},
			"TaskCompleted": []any{
				group("", exec(15, "gate", "task-complete")),
			},
			"PreToolUse": []any{
				group("Skill", exec(10, "gate", "skill")),
				group("Edit|Write", exec(10, "gate", "edit")),
				group("Bash", exec(10, "gate", "bash")),
			},
			"PostToolUse": []any{
				group("Bash", exec(10, "capture", "test")),
				group("Edit|Write", exec(10, "capture", "edit")),
			},
		},
	}
	raw, _ := json.MarshalIndent(doc, "", "  ")
	return append(raw, '\n')
}

// ---------------------------------------------------------------------------
// record schemas (external validation; the engine validates natively)
// ---------------------------------------------------------------------------

func recordSchemas(sp *spec.Spec) map[string][]byte {
	out := map[string][]byte{}
	for _, rk := range sp.Records {
		props := map[string]any{}
		for _, f := range rk.Known() {
			props[f] = map[string]any{}
		}
		schema := map[string]any{
			"$schema":  "https://json-schema.org/draft/2020-12/schema",
			"$id":      "https://github.com/dariusz-klibisz/ClaudeWorkflow/workflow/schemas/" + rk.Kind + ".json",
			"title":    "wf record: " + rk.Kind,
			"type":     "object",
			"required": rk.Required(),
			"properties": props,
		}
		raw, _ := json.MarshalIndent(schema, "", "  ")
		out[rk.Kind] = append(raw, '\n')
	}
	return out
}

// ---------------------------------------------------------------------------
// agent skeletons (create-if-missing; content is hand-authored afterwards)
// ---------------------------------------------------------------------------

const verdictSection = `
## Verdict (machine-parsed — required)

End the final message with exactly this fenced block (nothing after it):

` + "```" + `verdict
status: <clean|changes-required|safe|risky|unsafe|n/a>
criticals: <int>
majors: <int>
scope: <assigned mode/lens, when one was given>
` + "```" + `

Rules: clean/safe require criticals=0 and majors=0. risky requires each
concern listed above the block for disposition. n/a requires one line of
reason. The SubagentStop gate blocks completion until this block parses.
`

func agentSkeletons(check bool, root string, sp *spec.Spec) error {
	dir := filepath.Join(root, "agents")
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	for _, a := range sp.Roster {
		p := filepath.Join(dir, a.Name+".md")
		if _, err := os.Stat(p); err == nil {
			continue // hand-authored content is never overwritten
		}
		if check {
			return fmt.Errorf("agent file missing for roster entry %q (run go generate)", a.Name)
		}
		var b strings.Builder
		b.WriteString("---\n")
		fmt.Fprintf(&b, "name: %s\n", a.Name)
		fmt.Fprintf(&b, "description: wf %s (phases: %s). Spawned by the wf workflow with scope injected at start.\n", a.Name, strings.Join(a.Phases, ", "))
		if a.Model != "" {
			fmt.Fprintf(&b, "model: %s\n", a.Model)
		}
		if len(a.Tools) > 0 {
			fmt.Fprintf(&b, "tools: %s\n", strings.Join(a.Tools, ", "))
		}
		if a.MaxTurns > 0 {
			fmt.Fprintf(&b, "maxTurns: %d\n", a.MaxTurns)
		}
		if a.Memory != "" {
			fmt.Fprintf(&b, "memory: %s\n", a.Memory)
		}
		b.WriteString("---\n\n")
		fmt.Fprintf(&b, "# %s\n\nTODO(M2): full mandate. Follow the scope injected at SubagentStart.\n", a.Name)
		if a.Gating {
			b.WriteString(verdictSection)
		}
		if err := os.WriteFile(p, []byte(b.String()), 0o644); err != nil {
			return err
		}
		fmt.Println("created", p)
	}
	return nil
}

// ---------------------------------------------------------------------------

func writeOrCheck(check bool, path string, content []byte) (drift bool) {
	cur, err := os.ReadFile(path)
	if err == nil && string(cur) == string(content) {
		return false
	}
	if check {
		fmt.Fprintln(os.Stderr, "drift:", path)
		return true
	}
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		fatal(err)
	}
	if err := os.WriteFile(path, content, 0o644); err != nil {
		fatal(err)
	}
	fmt.Println("wrote", path)
	return false
}

func fatal(err error) {
	fmt.Fprintln(os.Stderr, "gen:", err)
	os.Exit(1)
}
