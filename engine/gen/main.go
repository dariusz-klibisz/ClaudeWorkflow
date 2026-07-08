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
	"regexp"
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
	// strict: gen is the dev/CI-side parse that catches spec typos
	sp, err := spec.LoadStrict(filepath.Join(root, "workflow", "workflow.yaml"), "")
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
	if *check {
		if err := checkProseParity(root, sp); err != nil {
			fatal(err)
		}
	}
	if *check && drift {
		fatal(fmt.Errorf("generated files drift from workflow.yaml — run `go generate ./...` (go run ./gen)"))
	}
}

// ---------------------------------------------------------------------------
// prose↔contract parity (the v0.36 gate_contract_map lesson, wf-native):
// skill/agent prose may only reference contract ids that exist, and every
// phase skill must mention each of its exit items' key noun — a workflow.yaml
// item nobody's procedure mentions (or a skill teaching a removed item) is
// CI-visible drift, not a silent gap.
// ---------------------------------------------------------------------------

var contractIDRe = regexp.MustCompile(`\b(frame|context|design|plan|build|verify|ship)\.[a-z][a-z0-9-]*\b`)

// parityAllowlist: exit items whose procedure lives outside the phase skill
// (engine-automated, or taught by a shared skill). Every entry needs a
// reason — this list is the honest record of what is NOT skill-taught.
var parityAllowlist = map[string]string{
	"build.tasks-closed":        "task lifecycle is native TaskCreate/TaskUpdate, gated by TaskCompleted — not a skill instruction",
	"ship.tasks-resolved":       "followup conversion is part of the close checklist in the dev skill",
	"verify.artifacts-authored": "generic stub sweep — the authoring instruction lives with each doc-creating item",
	"ship.artifacts-authored":   "generic stub sweep — see verify.artifacts-authored",
	"build.artifacts-authored":  "generic stub sweep — see verify.artifacts-authored",
	"design.artifacts-authored": "generic stub sweep — see verify.artifacts-authored",
}

func checkProseParity(root string, sp *spec.Spec) error {
	ids := map[string]bool{}
	for _, c := range sp.Contracts {
		ids[c.ID] = true
	}

	// 1. every contract-id-shaped token in prose must exist in the spec
	proseFiles, _ := filepath.Glob(filepath.Join(root, "skills", "*", "SKILL.md"))
	agentFiles, _ := filepath.Glob(filepath.Join(root, "agents", "*.md"))
	for _, p := range append(proseFiles, agentFiles...) {
		raw, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		for _, tok := range contractIDRe.FindAllString(string(raw), -1) {
			if !ids[tok] {
				return fmt.Errorf("%s references contract id %q which does not exist in workflow.yaml", p, tok)
			}
		}
	}

	// 2. each phase skill mentions every exit item's key noun
	for _, ph := range sp.Phases {
		p := filepath.Join(root, "skills", ph.Skill, "SKILL.md")
		raw, err := os.ReadFile(p)
		if err != nil {
			return fmt.Errorf("phase %s: skill file: %w", ph.ID, err)
		}
		content := strings.ToLower(string(raw))
		for _, it := range sp.Contracts {
			if it.Phase != ph.ID || it.Stage == "entry" {
				continue
			}
			if _, ok := parityAllowlist[it.ID]; ok {
				continue
			}
			mentions := itemMentions(it)
			if len(mentions) == 0 {
				continue
			}
			hit := false
			for _, m := range mentions {
				if strings.Contains(content, strings.ToLower(m)) {
					hit = true
					break
				}
			}
			if !hit {
				return fmt.Errorf("skill %s never mentions contract item %s (expected one of: %s) — teach it or allowlist it with a reason",
					p, it.ID, strings.Join(mentions, ", "))
			}
		}
	}
	return nil
}

// itemMentions derives the noun(s) a skill must utter to count as teaching
// the item: the record kind, the gating agent, or the approve command.
func itemMentions(it spec.ContractItem) []string {
	return predicateMentions(it.Predicate, it.Params)
}

func predicateMentions(pred string, params map[string]any) []string {
	switch pred {
	case spec.PredRecordExists, spec.PredNoOpen, spec.PredLinkedRecord, spec.PredPerEach:
		if pred == spec.PredPerEach {
			// the per-each ITEM is what the agent produces
			if item, ok := params["item"].(map[string]any); ok {
				if p, prm, err := spec.SubItem(item); err == nil {
					return predicateMentions(p, prm)
				}
			}
		}
		if k, _ := params["kind"].(string); k != "" {
			return []string{k}
		}
	case spec.PredVerdictIn:
		if a, _ := params["agent"].(string); a != "" {
			return []string{a}
		}
	case spec.PredApproval:
		if g, _ := params["gate"].(string); g != "" {
			return []string{"approve " + g}
		}
	case spec.PredRedGreen:
		return []string{"red", "test-first"}
	case spec.PredArtifactPresent:
		var out []string
		if t, _ := params["template"].(string); t != "" {
			out = append(out, t)
		}
		if r, _ := params["role"].(string); r != "" {
			out = append(out, r)
		}
		if len(out) == 0 {
			out = append(out, "artifact")
		}
		return out
	case spec.PredAnyOf:
		var out []string
		if items, ok := params["items"].([]any); ok {
			for _, raw := range items {
				if m, ok := raw.(map[string]any); ok {
					if p, prm, err := spec.SubItem(m); err == nil {
						out = append(out, predicateMentions(p, prm)...)
					}
				}
			}
		}
		return out
	}
	return nil
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
	// PowerShell. Native Windows is covered without hook tricks (M5):
	// the FIRST install is one manual step (bootstrap.ps1 or `wf doctor
	// --bootstrap` — see skills/init), and every later update is
	// engine-mediated — `wf inject session` (entry 2, exec form, runs on
	// all platforms once a binary exists) detects version skew against the
	// plugin root and re-runs the bootstrap itself (doctor.SelfUpdate).
	// Invoked via `sh <script>` so it works even when the cache copy lost
	// the executable bit (the dead-hooks incident: git had mode 100644).
	bootstrapSh := map[string]any{
		"type": "command", "timeout": 60,
		"command": `sh "${CLAUDE_PLUGIN_ROOT}/scripts/bootstrap.sh"`,
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
				group("Bash", exec(10, "capture", "bash")),
				group("Edit|Write", exec(10, "capture", "edit")),
				// approval anchoring (04 §8.1): the user's answers become
				// hook-captured user-answer records that wf approve links
				group("AskUserQuestion", exec(10, "capture", "question")),
			},
			// failed tool calls never fire PostToolUse — a RED test run
			// (non-zero exit) arrives here, with the exit code embedded in
			// the error string. Without this entry no red evidence is ever
			// auto-captured (the four-TestRepo-runs incident).
			"PostToolUseFailure": []any{
				group("Bash", exec(10, "capture", "bash")),
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
			"$schema":    "https://json-schema.org/draft/2020-12/schema",
			"$id":        "https://github.com/dariusz-klibisz/ClaudeWorkflow/workflow/schemas/" + rk.Kind + ".json",
			"title":      "wf record: " + rk.Kind,
			"type":       "object",
			"required":   rk.Required(),
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
			// hand-authored content is never overwritten — but the roster
			// and the frontmatter must agree on memory (the M4 `memory:
			// project` subset is deliberate: design-reviewer,
			// code-quality-reviewer, adversary — 06 §"Frontmatter").
			if check {
				if err := checkAgentMemory(p, a); err != nil {
					return err
				}
				if err := checkAgentModel(p, a); err != nil {
					return err
				}
			}
			continue
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

// checkAgentMemory verifies the roster's `memory:` field matches the agent
// file's frontmatter — drift here silently changes which agents accumulate
// cross-run recall.
func checkAgentMemory(path string, a spec.Agent) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	got := frontmatterField(string(raw), "memory")
	if got != a.Memory {
		return fmt.Errorf("agent %q: frontmatter memory %q != roster memory %q (align agents/%s.md with workflow.yaml)", a.Name, got, a.Memory, a.Name)
	}
	return nil
}

// checkAgentModel verifies the roster's `model:` field matches the agent
// file's frontmatter — drift here silently changes which model reviews
// (the reviewer-vs-author asymmetry is roster-declared, WS-E).
func checkAgentModel(path string, a spec.Agent) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	got := frontmatterField(string(raw), "model")
	if got != a.Model {
		return fmt.Errorf("agent %q: frontmatter model %q != roster model %q (align agents/%s.md with workflow.yaml)", a.Name, got, a.Model, a.Name)
	}
	return nil
}

// frontmatterField extracts a scalar field from the leading --- block.
func frontmatterField(content, field string) string {
	lines := strings.Split(content, "\n")
	if len(lines) == 0 || strings.TrimSpace(lines[0]) != "---" {
		return ""
	}
	for _, line := range lines[1:] {
		if strings.TrimSpace(line) == "---" {
			break
		}
		if v, ok := strings.CutPrefix(strings.TrimSpace(line), field+":"); ok {
			return strings.TrimSpace(v)
		}
	}
	return ""
}

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
