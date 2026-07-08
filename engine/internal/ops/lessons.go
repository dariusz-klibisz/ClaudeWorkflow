package ops

// wf lessons — the enforcement loop for what runs teach (03 §4.7).
// Two lesson species, two delivery channels:
//   - `check:` lessons carry a contract-item fragment; on acceptance they are
//     written as ordinary `lesson.*` items into .workflow/contracts.d/
//     lessons.yaml — one representation, one evaluator, enforced next run.
//   - prose lessons (no check) regenerate .claude/rules/wf-lessons.md
//     (marker-delimited, engine-owned, committed) — unscoped rules re-inject
//     after compaction (05 §2).
// Both files are regenerated idempotently from lesson records, read from the
// committed run archives + the live log (lesson events archive with their
// run — the bounded-live-log rule, 08 §6).

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
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

// triggerClasses maps a run's health signals to the engine suggest-trigger
// classes that fire, each with its lesson text. Shared by LessonsSuggest
// (proposing) and LessonsStatus (efficacy: did an accepted lesson's trigger
// fire again in a later run?). Class names are stable — they are stamped
// into lesson records as `trigger`.
func triggerClasses(s *views.RunSignals) map[string]string {
	out := map[string]string{}
	if s.Forces > 0 {
		out["forces"] = fmt.Sprintf("This run forced %d gate(s). Review what blocked and fix the process (or propose a contract change) instead of forcing.", s.Forces)
	}
	if s.Loops >= 2 {
		out["loops"] = fmt.Sprintf("This run looped %d times. Check whether Frame/Plan missed the cause — loops past the first usually mean an upstream gap.", s.Loops)
	}
	if s.LoopsByCause["design"]+s.LoopsByCause["plan"] >= 2 {
		out["upstream-loops"] = fmt.Sprintf("Verify re-opened Design/Plan %d times (causes: design %d, plan %d). The framing or design review is letting structural gaps through — tighten those gates, don't just fix the instance.",
			s.LoopsByCause["design"]+s.LoopsByCause["plan"], s.LoopsByCause["design"], s.LoopsByCause["plan"])
	}
	for ac, n := range s.LoopsByAC {
		if n >= 2 {
			out["ac-churn"] = fmt.Sprintf("AC %s looped %d times. The criterion may be unverifiable as written, or the selected design doesn't support it — revisit the AC before the next re-entry.", ac, n)
			break
		}
	}
	if s.Verdicts >= 2 && s.AutoVerdicts == 0 {
		out["manual-verdicts"] = "All reviewer verdicts were recorded manually — the SubagentStop capture never fired. Run wf doctor --bootstrap at session start."
	}
	// test-run grounding signals are diff-family only: artifact/assessment
	// runs verify documents with manual grep-style checks by design — zero
	// auto-captures there is expected, not a lesson (the arch-design run's
	// false-positive suggestion).
	if s.Family == "diff" {
		if s.TestRuns >= 3 && s.AutoTestRuns == 0 {
			out["manual-tests"] = "No test run was auto-captured — the runner isn't recognized. Record verification-strategy commands as the real invocations, or add the wrapper to config \"runners\"."
		}
		if len(s.UngroundedACs) > 0 {
			out["ungrounded-acs"] = fmt.Sprintf("AC(s) %s passed without a grounded green test-run. Tag AC verification runs (--ac) so evidence is hook-captured.", strings.Join(s.UngroundedACs, ", "))
		}
	}
	return out
}

// LessonsSuggest proposes lessons from the active run's health signals —
// engine-spotted, user-triaged. Idempotent: an existing lesson with the same
// normalized text (any status) suppresses the suggestion. Each proposal is
// stamped with its trigger class so LessonsStatus can track recurrence.
func LessonsSuggest(c *runctl.Ctl) (string, error) {
	s, err := views.ReportRun(c, "current")
	if err != nil {
		return "", err
	}
	classes := triggerClasses(s)
	order := make([]string, 0, len(classes))
	for cl := range classes {
		order = append(order, cl)
	}
	sort.Strings(order)

	existing := map[string]bool{}
	for _, l := range lessonRecords(c) {
		text, _ := l.Data["text"].(string)
		existing[normText(text)] = true
	}
	proposed := 0
	var b strings.Builder
	for _, cl := range order {
		t := classes[cl]
		if existing[normText(t)] {
			continue
		}
		ev, err := c.Record("lesson", map[string]any{"text": t, "status": "proposed", "source": "engine", "trigger": cl}, true, "engine")
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

// LessonsStatus renders the efficacy view: per lesson — origin run, status,
// and what happened since acceptance. For check-lessons: how often the
// generated lesson.* item was waived (accepted but dodged). For
// engine-suggested lessons: whether the trigger class fired again in a
// later run (accepted but not working). Derived, never gates.
func LessonsStatus(c *runctl.Ctl) (string, error) {
	lessons := lessonRecords(c)
	if len(lessons) == 0 {
		return "[wf lessons] no lessons recorded yet", nil
	}

	// origin run per lesson + waiver runs per contract item, from raw events
	originRun := map[string]string{}
	waivedIn := map[string][]string{}
	if evs, err := c.Store.AllEvents(); err == nil {
		for _, e := range evs {
			switch e.Kind {
			case "lesson":
				if _, isUpdate := e.Data["updates"].(string); !isUpdate {
					originRun[e.ID] = e.Run
				}
			case "waiver":
				if item, _ := e.Data["item"].(string); strings.HasPrefix(item, "lesson.") {
					waivedIn[item] = append(waivedIn[item], e.Run)
				}
			}
		}
	}

	// per-run signals, run-ID ordered (IDs are date-prefixed)
	sigs, err := views.Report(c)
	if err != nil {
		return "", err
	}
	classByRun := map[string]map[string]string{}
	runOrder := make([]string, 0, len(sigs))
	for i := range sigs {
		classByRun[sigs[i].Run] = triggerClasses(&sigs[i])
		runOrder = append(runOrder, sigs[i].Run)
	}
	sort.Strings(runOrder)

	var b strings.Builder
	fmt.Fprintf(&b, "[wf lessons] %d lesson(s)\n", len(lessons))
	concerns := 0
	for _, l := range lessons {
		status, _ := l.Data["status"].(string)
		text, _ := l.Data["text"].(string)
		origin := originRun[l.ID]
		fmt.Fprintf(&b, "  %s · %s · from %s — %s\n", l.ID, status, orStr(origin, "?"), clipText(text, 90))
		if status != "accepted" {
			continue
		}
		if check, _ := l.Data["check"].(string); check != "" {
			lc := l
			if item, err := checkToItem(&lc); err == nil {
				if runs := waivedIn[item.ID]; len(runs) > 0 {
					concerns++
					fmt.Fprintf(&b, "    ⚠ enforced as %s — waived in %d run(s): %s (accepted but dodged)\n", item.ID, len(runs), strings.Join(uniq(runs), ", "))
				} else {
					fmt.Fprintf(&b, "    enforced as %s — never waived\n", item.ID)
				}
			}
			continue
		}
		if trigger, _ := l.Data["trigger"].(string); trigger != "" {
			var recurred []string
			for _, run := range runOrder {
				if origin != "" && run <= origin {
					continue
				}
				if _, hit := classByRun[run][trigger]; hit {
					recurred = append(recurred, run)
				}
			}
			if len(recurred) > 0 {
				concerns++
				fmt.Fprintf(&b, "    ⚠ trigger %q recurred in %d later run(s): %s (accepted but not working)\n", trigger, len(recurred), strings.Join(recurred, ", "))
			} else {
				fmt.Fprintf(&b, "    trigger %q has not recurred since acceptance\n", trigger)
			}
		}
	}
	if concerns > 0 {
		fmt.Fprintf(&b, "%d efficacy concern(s) — a dodged or recurring lesson is a process gap, not a formality\n", concerns)
	}
	return b.String(), nil
}

func orStr(s, def string) string {
	if s == "" {
		return def
	}
	return s
}

func clipText(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-1] + "…"
}

func uniq(xs []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, x := range xs {
		if !seen[x] {
			seen[x] = true
			out = append(out, x)
		}
	}
	return out
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

// lessonRecords folds ALL lesson events — archived runs + the live log —
// into effective records. Lesson events archive with their run (bounded
// live log); regeneration reads them back from the committed archives.
func lessonRecords(c *runctl.Ctl) []contracts.Record {
	evs, err := c.Store.AllEvents()
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
