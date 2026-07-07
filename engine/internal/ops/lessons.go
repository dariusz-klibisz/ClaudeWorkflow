package ops

// wf lessons — the enforcement loop for what runs teach (03 §4.7).
// Two lesson species, two delivery channels:
//   - `check:` lessons carry a contract-item fragment; on acceptance they are
//     written as ordinary `lesson.*` items into .workflow/contracts.d/
//     lessons.yaml — one representation, one evaluator, enforced next run.
//   - prose lessons (no check) regenerate .claude/rules/wf-lessons.md
//     (marker-delimited, engine-owned, committed) — unscoped rules re-inject
//     after compaction (05 §2).
// Both files are regenerated idempotently from lesson records, which survive
// run close in the live log (store keepLive, 08 §6).

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/contracts"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/runctl"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/spec"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/views"
)

const (
	lessonsYAML     = "lessons.yaml"
	lessonsRulePath = ".claude/rules/wf-lessons.md"
	lessonsBegin    = "<!-- wf:lessons:begin -->"
	lessonsEnd      = "<!-- wf:lessons:end -->"
)

// LessonsSuggest proposes lessons from the active run's health signals —
// engine-spotted, user-triaged. Idempotent: an existing lesson with the same
// normalized text (any status) suppresses the suggestion.
func LessonsSuggest(c *runctl.Ctl) (string, error) {
	s, err := views.ReportRun(c, "current")
	if err != nil {
		return "", err
	}
	var texts []string
	if s.Forces > 0 {
		texts = append(texts, fmt.Sprintf("This run forced %d gate(s). Review what blocked and fix the process (or propose a contract change) instead of forcing.", s.Forces))
	}
	if s.Loops >= 2 {
		texts = append(texts, fmt.Sprintf("This run looped %d times. Check whether Frame/Plan missed the cause — loops past the first usually mean an upstream gap.", s.Loops))
	}
	if s.Verdicts >= 2 && s.AutoVerdicts == 0 {
		texts = append(texts, "All reviewer verdicts were recorded manually — the SubagentStop capture never fired. Run wf doctor --bootstrap at session start.")
	}
	// test-run grounding signals are diff-family only: artifact/assessment
	// runs verify documents with manual grep-style checks by design — zero
	// auto-captures there is expected, not a lesson (the arch-design run's
	// false-positive suggestion).
	if s.Family == "diff" {
		if s.TestRuns >= 3 && s.AutoTestRuns == 0 {
			texts = append(texts, "No test run was auto-captured — the runner isn't recognized. Record verification-strategy commands as the real invocations, or add the wrapper to config \"runners\".")
		}
		if len(s.UngroundedACs) > 0 {
			texts = append(texts, fmt.Sprintf("AC(s) %s passed without a grounded green test-run. Tag AC verification runs (--ac) so evidence is hook-captured.", strings.Join(s.UngroundedACs, ", ")))
		}
	}

	existing := map[string]bool{}
	for _, l := range lessonRecords(c) {
		text, _ := l.Data["text"].(string)
		existing[normText(text)] = true
	}
	proposed := 0
	var b strings.Builder
	for _, t := range texts {
		if existing[normText(t)] {
			continue
		}
		ev, err := c.Record("lesson", map[string]any{"text": t, "status": "proposed", "source": "engine"}, true, "engine")
		if err != nil {
			return "", err
		}
		proposed++
		fmt.Fprintf(&b, "  %s — %s\n", ev.ID, t)
	}
	if proposed == 0 {
		return "wf lessons suggest: nothing to propose (signals clean or already covered)", nil
	}
	return fmt.Sprintf("proposed %d lesson(s) — triage each with wf lessons accept|reject <id>:\n%s", proposed, b.String()), nil
}

// LessonsAccept validates the lesson's check (if any) against the full
// merged spec BEFORE recording acceptance — a broken check must fail loudly
// here, not at the next run's load. On success: approval event + status
// flip + regeneration of both delivery channels.
func LessonsAccept(c *runctl.Ctl, projectDir, specPath, id string) (string, error) {
	rec, err := findLesson(c, id)
	if err != nil {
		return "", err
	}
	if check, _ := rec.Data["check"].(string); check != "" {
		item, err := checkToItem(rec)
		if err != nil {
			return "", fmt.Errorf("check invalid — fix the lesson's check field before accepting: %w", err)
		}
		// validate the WOULD-BE lessons.yaml (current accepted set + this one)
		items := acceptedCheckItems(c)
		items = append(items, *item)
		if err := validateLessonItems(c, specPath, items); err != nil {
			return "", fmt.Errorf("check does not validate against the spec: %w", err)
		}
	}
	if _, err := c.Approve("lesson", id); err != nil {
		return "", err
	}
	if _, err := c.Record("lesson", map[string]any{"updates": rec.ID, "status": "accepted"}, false, "user"); err != nil {
		return "", err
	}
	applied, err := LessonsApply(c, projectDir, specPath)
	if err != nil {
		return "", err
	}
	return "lesson " + rec.ID + " accepted\n" + applied, nil
}

// LessonsReject flips the lesson to rejected (with the approval event) and
// regenerates — a previously accepted-then-rejected lesson's artifacts go.
func LessonsReject(c *runctl.Ctl, projectDir, specPath, id string) (string, error) {
	rec, err := findLesson(c, id)
	if err != nil {
		return "", err
	}
	if _, err := c.Approve("lesson", id); err != nil {
		return "", err
	}
	if _, err := c.Record("lesson", map[string]any{"updates": rec.ID, "status": "rejected"}, false, "user"); err != nil {
		return "", err
	}
	applied, err := LessonsApply(c, projectDir, specPath)
	if err != nil {
		return "", err
	}
	return "lesson " + rec.ID + " rejected\n" + applied, nil
}

// LessonsApply regenerates both delivery channels from ALL accepted lessons
// (across runs — lesson records stay live). Validation before write: the
// merged spec must load strictly with the new lessons.yaml, or nothing
// changes on disk.
func LessonsApply(c *runctl.Ctl, projectDir, specPath string) (string, error) {
	items := acceptedCheckItems(c)
	if err := validateLessonItems(c, specPath, items); err != nil {
		return "", fmt.Errorf("refusing to write lessons.yaml: %w", err)
	}

	// channel 1: contracts.d/lessons.yaml
	target := filepath.Join(c.Store.ContractsDir(), lessonsYAML)
	if len(items) == 0 {
		if err := os.Remove(target); err != nil && !os.IsNotExist(err) {
			return "", err
		}
	} else {
		var b strings.Builder
		b.WriteString("# GENERATED by `wf lessons apply` from accepted check-lessons — do not edit.\n")
		b.WriteString("# A lesson-check IS a contract item (03 §4.0): one representation, one evaluator.\n")
		raw, err := yaml.Marshal(map[string]any{"contracts": items})
		if err != nil {
			return "", err
		}
		b.Write(raw)
		if err := writeAtomic(target, []byte(b.String())); err != nil {
			return "", err
		}
	}

	// channel 2: .claude/rules/wf-lessons.md (prose lessons)
	var prose []string
	for _, l := range lessonRecords(c) {
		if status, _ := l.Data["status"].(string); status != "accepted" {
			continue
		}
		if check, _ := l.Data["check"].(string); check != "" {
			continue
		}
		text, _ := l.Data["text"].(string)
		prose = append(prose, fmt.Sprintf("- %s <!-- %s -->", text, l.ID))
	}
	rulePath := filepath.Join(projectDir, filepath.FromSlash(lessonsRulePath))
	if err := writeProseRules(rulePath, prose); err != nil {
		return "", err
	}

	return fmt.Sprintf("applied: %d check-lesson contract item(s) → %s · %d prose lesson(s) → %s",
		len(items), filepath.ToSlash(filepath.Join(".workflow", "contracts.d", lessonsYAML)),
		len(prose), lessonsRulePath), nil
}

// ---------------------------------------------------------------------------
// internals
// ---------------------------------------------------------------------------

// lessonRecords folds ALL lesson events in the live log (any run — lessons
// survive close) into effective records.
func lessonRecords(c *runctl.Ctl) []contracts.Record {
	evs, err := c.Store.Events(nil)
	if err != nil {
		return nil
	}
	env := &contracts.Env{Events: evs}
	return env.Records("lesson")
}

func findLesson(c *runctl.Ctl, id string) (*contracts.Record, error) {
	for _, l := range lessonRecords(c) {
		if l.ID == id {
			return &l, nil
		}
	}
	return nil, fmt.Errorf("no lesson record %s (wf report --run current shows lesson counts; IDs are printed at proposal time)", id)
}

// checkToItem parses a lesson's check field — a contract-item fragment in
// YAML or JSON (JSON is valid YAML) — and namespaces its ID.
func checkToItem(rec *contracts.Record) (*spec.ContractItem, error) {
	check, _ := rec.Data["check"].(string)
	var item spec.ContractItem
	if err := yaml.Unmarshal([]byte(check), &item); err != nil {
		return nil, fmt.Errorf("check must be a contract-item fragment (phase, predicate, params, remediation): %w", err)
	}
	if item.ID == "" {
		text, _ := rec.Data["text"].(string)
		item.ID = "lesson." + slugID(text, rec.ID)
	} else if !strings.HasPrefix(item.ID, "lesson.") {
		item.ID = "lesson." + item.ID
	}
	if item.Remediation == "" {
		text, _ := rec.Data["text"].(string)
		item.Remediation = text
	}
	return &item, nil
}

// acceptedCheckItems builds the contract items for every accepted lesson
// carrying a check. Unparseable checks are skipped here (accept-time
// validation prevents them; historical bad data must not brick apply).
func acceptedCheckItems(c *runctl.Ctl) []spec.ContractItem {
	var items []spec.ContractItem
	seen := map[string]bool{}
	for _, l := range lessonRecords(c) {
		if status, _ := l.Data["status"].(string); status != "accepted" {
			continue
		}
		if check, _ := l.Data["check"].(string); check == "" {
			continue
		}
		lc := l
		item, err := checkToItem(&lc)
		if err != nil || seen[item.ID] {
			continue
		}
		seen[item.ID] = true
		items = append(items, *item)
	}
	return items
}

// validateLessonItems round-trips the would-be lessons.yaml through the
// strict spec loader against a mirror of contracts.d — full validation
// (namespacing, predicates, phases, kinds), zero risk to the real dir.
func validateLessonItems(c *runctl.Ctl, specPath string, items []spec.ContractItem) error {
	tmp, err := os.MkdirTemp("", "wf-lessons-*")
	if err != nil {
		return err
	}
	defer os.RemoveAll(tmp)
	// mirror existing extension files (minus any current lessons.yaml)
	if entries, err := os.ReadDir(c.Store.ContractsDir()); err == nil {
		for _, e := range entries {
			name := e.Name()
			if e.IsDir() || name == lessonsYAML || (!strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml")) {
				continue
			}
			raw, err := os.ReadFile(filepath.Join(c.Store.ContractsDir(), name))
			if err != nil {
				return err
			}
			if err := os.WriteFile(filepath.Join(tmp, name), raw, 0o644); err != nil {
				return err
			}
		}
	}
	if len(items) > 0 {
		raw, err := yaml.Marshal(map[string]any{"contracts": items})
		if err != nil {
			return err
		}
		if err := os.WriteFile(filepath.Join(tmp, lessonsYAML), raw, 0o644); err != nil {
			return err
		}
	}
	_, err = spec.LoadStrict(specPath, tmp)
	return err
}

// writeProseRules regenerates the marker-delimited block, preserving any
// user content outside the markers. No prose lessons + nothing else in the
// file ⇒ the file is removed.
func writeProseRules(path string, prose []string) error {
	block := lessonsBegin + "\n# Workflow lessons (wf)\nAccepted prose lessons from past runs — engine-generated, never hand-edit.\n\n" +
		strings.Join(prose, "\n") + "\n" + lessonsEnd
	existing, err := os.ReadFile(path)
	switch {
	case err == nil:
		content := string(existing)
		var rebuilt string
		if i, j := strings.Index(content, lessonsBegin), strings.Index(content, lessonsEnd); i >= 0 && j > i {
			rebuilt = content[:i] + block + content[j+len(lessonsEnd):]
		} else {
			rebuilt = strings.TrimRight(content, "\n") + "\n\n" + block + "\n"
		}
		if len(prose) == 0 {
			outside := strings.TrimSpace(strings.ReplaceAll(rebuilt, block, ""))
			if outside == "" {
				return os.Remove(path)
			}
		}
		return writeAtomic(path, []byte(rebuilt))
	case os.IsNotExist(err):
		if len(prose) == 0 {
			return nil
		}
		return writeAtomic(path, []byte(block+"\n"))
	default:
		return err
	}
}

func writeAtomic(path string, data []byte) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, data, 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

var nonWord = regexp.MustCompile(`[^a-z0-9]+`)

// normText normalizes lesson text for duplicate suppression.
func normText(s string) string {
	return strings.Trim(nonWord.ReplaceAllString(strings.ToLower(s), "-"), "-")
}

// slugID derives a stable lesson.* suffix from the text (first ~40 chars),
// falling back to the record's ULID for uniqueness.
func slugID(text, ulid string) string {
	s := strings.Trim(nonWord.ReplaceAllString(strings.ToLower(text), "-"), "-")
	if len(s) > 40 {
		s = strings.Trim(s[:40], "-")
	}
	if s == "" {
		return strings.ToLower(ulid)
	}
	return s
}
