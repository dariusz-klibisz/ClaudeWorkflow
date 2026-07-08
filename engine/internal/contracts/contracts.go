// Package contracts implements the closed predicate vocabulary and the
// phase-exit evaluator (workflow-redesign/03 §4.0, 04 §1). The evaluator
// reads engine-written records only — never prose — and returns the unmet
// contract items with their remediations.
package contracts

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/cmdid"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/spec"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/store"
)

// Finding is one unmet contract item.
type Finding struct {
	ID          string `json:"id"`
	Remediation string `json:"remediation"`
	Detail      string `json:"detail,omitempty"`
	UserBlocked bool   `json:"user_blocked"` // waiting on the user, not the agent
	Waivable    bool   `json:"waivable"`
}

// Env is everything an evaluation reads.
type Env struct {
	Spec   *spec.Spec
	Config *store.Config
	Run    *store.Run
	Events []store.Event // the run's events, append order
	// ProjectDir roots the artifact-present predicate's disk checks — the
	// single deliberate filesystem window in an otherwise records-only
	// evaluator. Empty = disk checks pass vacuously (report/lessons views).
	ProjectDir string
	// derived
	records map[string][]Record // effective records by kind (updates folded)
	alias   map[string]string   // any record/update event ID -> original record ID
}

// Record is an effective record: the original event with any later
// `updates: <id>` events folded over its data.
type Record struct {
	ID    string
	Order int // position in the event stream (for red-green ordering)
	Auto  bool
	Data  map[string]any
}

// Build folds update events and indexes records by kind.
func (e *Env) build() {
	if e.records != nil {
		return
	}
	e.records = map[string][]Record{}
	e.alias = map[string]string{}
	byID := map[string]*Record{}
	kindOf := map[string]string{}
	// alias maps any event ID (original or update) to the original record —
	// chaining `updates=` onto a prior update's ID must resolve transitively
	// (the silent-no-op bug the power5 run hit).
	alias := e.alias
	for i, ev := range e.Events {
		if ev.Kind == "run" || ev.Kind == "phase" {
			continue // engine transitions, not records
		}
		if target, ok := ev.Data["updates"].(string); ok && target != "" {
			if orig, ok := alias[target]; ok {
				target = orig
			}
			if rec, ok := byID[target]; ok {
				for k, v := range ev.Data {
					if k == "updates" {
						continue
					}
					rec.Data[k] = v
				}
				rec.Auto = rec.Auto || ev.Auto
				alias[ev.ID] = target
			}
			continue
		}
		alias[ev.ID] = ev.ID
		data := map[string]any{}
		for k, v := range ev.Data {
			data[k] = v
		}
		rec := &Record{ID: ev.ID, Order: i, Auto: ev.Auto, Data: data}
		byID[ev.ID] = rec
		kindOf[ev.ID] = ev.Kind
	}
	ids := make([]string, 0, len(byID))
	for id := range byID {
		ids = append(ids, id)
	}
	sort.Slice(ids, func(a, b int) bool { return byID[ids[a]].Order < byID[ids[b]].Order })
	for _, id := range ids {
		k := kindOf[id]
		e.records[k] = append(e.records[k], *byID[id])
	}
}

// Records returns effective records of a kind, in stream order.
func (e *Env) Records(kind string) []Record {
	e.build()
	return e.records[kind]
}

// ResolveRecordID maps any record or update event ID to the original record
// ID (used to validate `updates=` targets at write time).
func (e *Env) ResolveRecordID(id string) (string, bool) {
	e.build()
	orig, ok := e.alias[id]
	return orig, ok
}

// ---------------------------------------------------------------------------
// Evaluation
// ---------------------------------------------------------------------------

// EvaluatePhase returns the unmet contract items for exiting the run's
// current phase.
func EvaluatePhase(env *Env, phase string) ([]Finding, error) {
	return evaluateItems(env, env.Spec.ContractsFor(phase, env.Run.Family))
}

// EvaluateEntry returns the unmet INPUT items for entering a phase —
// blocking at the previous phase's exit, re-checked by the skill gate on
// adopt/resume paths that never crossed the transition.
func EvaluateEntry(env *Env, phase string) ([]Finding, error) {
	return evaluateItems(env, env.Spec.EntryContractsFor(phase, env.Run.Family))
}

func evaluateItems(env *Env, items []spec.ContractItem) ([]Finding, error) {
	var out []Finding
	for _, it := range items {
		if !whenApplies(env, it.When) {
			continue
		}
		if it.Waivable && env.waived(it.ID) {
			continue
		}
		ok, detail, err := evalPredicate(env, it.Predicate, it.Params, "")
		if err != nil {
			return nil, fmt.Errorf("contract %s: %w", it.ID, err)
		}
		if !ok {
			out = append(out, Finding{
				ID:          it.ID,
				Remediation: it.Remediation,
				Detail:      detail,
				UserBlocked: it.Predicate == spec.PredApproval,
				Waivable:    it.Waivable,
			})
		}
	}
	return out, nil
}

func (e *Env) waived(itemID string) bool {
	for _, w := range e.Records("waiver") {
		if s, _ := w.Data["item"].(string); s == itemID {
			return true
		}
	}
	return false
}

func whenApplies(env *Env, w *spec.When) bool {
	if w == nil {
		return true
	}
	if len(w.Intents) > 0 {
		found := false
		for _, i := range w.Intents {
			if i == env.Run.Intent {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	for key, want := range w.Config {
		got := env.Config.ConfigFlag(key)
		if !valueEq(got, want) {
			return false
		}
	}
	if len(w.Signals) > 0 {
		if !runHasSignal(env, w.Signals) {
			return false
		}
	}
	if w.Records != nil {
		if len(matchRecords(env, w.Records.Kind, w.Records.Filter)) == 0 {
			return false
		}
	}
	return true
}

func runHasSignal(env *Env, want []string) bool {
	for _, r := range env.Records("risk") {
		sigs, _ := r.Data["signals"].([]any)
		for _, s := range sigs {
			ss := asSignal(s)
			for _, w := range want {
				if ss == w {
					return true
				}
			}
		}
	}
	return false
}

func asSignal(v any) string {
	switch t := v.(type) {
	case string:
		return t
	case map[string]any:
		s, _ := t["signal"].(string)
		return s
	}
	return ""
}

// EvalOne evaluates a single predicate expression against the env — used by
// the task gates to check one task's DoD without a full phase evaluation.
func EvalOne(env *Env, pred string, params map[string]any, ctxID string) (bool, string, error) {
	return evalPredicate(env, pred, params, ctxID)
}

// ---------------------------------------------------------------------------
// Predicates
// ---------------------------------------------------------------------------

func evalPredicate(env *Env, pred string, params map[string]any, ctxID string) (bool, string, error) {
	switch pred {
	case spec.PredRecordExists:
		return evalRecordExists(env, params)
	case spec.PredLinkedRecord:
		return evalLinkedRecord(env, params, ctxID)
	case spec.PredVerdictIn:
		return evalVerdictIn(env, params)
	case spec.PredApproval:
		return evalApproval(env, params)
	case spec.PredNoOpen:
		return evalNoOpen(env, params)
	case spec.PredPerEach:
		return evalPerEach(env, params)
	case spec.PredAnyOf:
		return evalAnyOf(env, params, ctxID)
	case spec.PredRedGreen:
		return evalRedGreen(env, params, ctxID)
	case spec.PredArtifactPresent:
		return evalArtifactPresent(env, params)
	}
	return false, "", fmt.Errorf("unknown predicate %q", pred)
}

// evalArtifactPresent: at least one artifact record matching the template/
// role filter exists, its file exists on disk, and the content is authored
// (non-stub). This is the mechanical replacement for trusting a
// self-reported `status: present` — the record claims, the disk confirms.
func evalArtifactPresent(env *Env, p map[string]any) (bool, string, error) {
	filter := map[string]any{}
	for _, key := range []string{"template", "role"} {
		if v, ok := p[key].(string); ok && v != "" {
			filter[key] = v
		}
	}
	cands := matchRecords(env, "artifact", filter)
	if len(cands) == 0 {
		return false, describeArtifactFilter(filter, "no artifact record"), nil
	}
	required := mentionIDs(env, p)
	var details []string
	for _, r := range cands {
		if s, _ := r.Data["status"].(string); s == "missing" {
			continue // explicitly abandoned
		}
		rel, _ := r.Data["path"].(string)
		ok, detail := ArtifactOnDisk(env.ProjectDir, rel, env.templatePath(r))
		if !ok {
			details = append(details, detail)
			continue
		}
		// mentions: every required record id must appear in the authored
		// file — a findings report that names none of its findings is a
		// stub with prose (env.ProjectDir=="" views skip, like the disk check)
		if missing := missingMentions(env.ProjectDir, rel, required); len(missing) > 0 {
			details = append(details, fmt.Sprintf("%s never mentions: %s", rel, strings.Join(missing, ", ")))
			continue
		}
		return true, "", nil
	}
	if len(details) == 0 {
		return false, describeArtifactFilter(filter, "every matching artifact abandoned (status missing)"), nil
	}
	return false, strings.Join(dedupe(details), "; "), nil
}

// mentionIDs collects the record ids an artifact must name (the `mentions`
// param: {kind, field}).
func mentionIDs(env *Env, p map[string]any) []string {
	m, ok := p["mentions"].(map[string]any)
	if !ok {
		return nil
	}
	kind, _ := m["kind"].(string)
	field, _ := m["field"].(string)
	if kind == "" || field == "" {
		return nil
	}
	var out []string
	for _, r := range env.Records(kind) {
		if id, _ := r.Data[field].(string); id != "" {
			out = append(out, id)
		}
	}
	return dedupe(out)
}

// missingMentions reads the artifact and returns the required ids absent
// from its content. Empty projectDir passes vacuously (views).
func missingMentions(projectDir, rel string, required []string) []string {
	if projectDir == "" || len(required) == 0 || rel == "" {
		return nil
	}
	content, err := os.ReadFile(filepath.Join(projectDir, filepath.FromSlash(rel)))
	if err != nil {
		return nil // the disk check already passed; a read race is not a mentions failure
	}
	text := string(content)
	var missing []string
	for _, id := range required {
		if !strings.Contains(text, id) {
			missing = append(missing, id)
		}
	}
	return missing
}

// templatePath resolves the plugin template an artifact record was created
// from, for the identical-to-template stub check. Empty when unknown.
func (e *Env) templatePath(r Record) string {
	tpl, _ := r.Data["template"].(string)
	if tpl == "" || e.Spec == nil || e.Spec.PluginRoot() == "" {
		return ""
	}
	return filepath.Join(e.Spec.PluginRoot(), "templates", tpl+".md")
}

func describeArtifactFilter(filter map[string]any, prefix string) string {
	var parts []string
	for k, v := range filter {
		parts = append(parts, fmt.Sprintf("%s=%v", k, v))
	}
	sort.Strings(parts)
	if len(parts) == 0 {
		return prefix
	}
	return prefix + " (" + strings.Join(parts, ", ") + ")"
}

// Stub heuristic thresholds (ported from the v0.36 artifact check, which
// they survived a year of field use in): fewer than 5 non-blank lines or
// under 200 dense characters is a skeleton, not an authored document.
const (
	stubMinLines = 5
	stubMinDense = 200
)

// ArtifactOnDisk verifies an artifact file exists and is authored. tmplPath
// (optional) enables the strongest check: byte-identical to its template =
// untouched stub. Empty projectDir passes vacuously (views without a
// filesystem root). Shared with runctl's write-time `status: present` gate.
func ArtifactOnDisk(projectDir, rel, tmplPath string) (bool, string) {
	if projectDir == "" {
		return true, ""
	}
	if rel == "" {
		return false, "artifact record has no path"
	}
	abs := filepath.Join(projectDir, filepath.FromSlash(rel))
	content, err := os.ReadFile(abs)
	if err != nil {
		return false, fmt.Sprintf("%s does not exist on disk", rel)
	}
	if tmplPath != "" {
		if tmpl, err := os.ReadFile(tmplPath); err == nil && string(tmpl) == string(content) {
			return false, fmt.Sprintf("%s is byte-identical to its template — author it", rel)
		}
	}
	lines, dense := 0, 0
	for _, ln := range strings.Split(string(content), "\n") {
		t := strings.TrimSpace(ln)
		if t == "" {
			continue
		}
		lines++
		dense += len(t)
	}
	if lines < stubMinLines || dense < stubMinDense {
		return false, fmt.Sprintf("%s looks like a stub (%d non-blank lines, %d chars)", rel, lines, dense)
	}
	return true, ""
}

func evalRecordExists(env *Env, p map[string]any) (bool, string, error) {
	kind, _ := p["kind"].(string)
	filter, _ := p["filter"].(map[string]any)
	min := intParam(p, "min", 1)
	matched := matchRecords(env, kind, filter)
	if len(matched) < min {
		return false, fmt.Sprintf("%d/%d %s record(s)", len(matched), min, kind), nil
	}
	// elements: {field, min} — a min-content floor over an array field,
	// counted across all matching records (the C10 fix: a one-entry
	// context map is a checkbox, not a map).
	if el, ok := p["elements"].(map[string]any); ok {
		field, _ := el["field"].(string)
		emin := intParam(el, "min", 1)
		total := 0
		for _, r := range matched {
			total += len(elementIDsAll(r, field))
		}
		if total < emin {
			return false, fmt.Sprintf("%s has %d %s (min %d)", kind, total, field, emin), nil
		}
	}
	return true, "", nil
}

// elementIDsAll counts every element of an array field: strings, {id} maps
// AND id-less maps (a completeness item has case/disposition, no id —
// elementIDs would drop it).
func elementIDsAll(r Record, field string) []any {
	raw, ok := r.Data[field].([]any)
	if !ok {
		return nil
	}
	return raw
}

func evalLinkedRecord(env *Env, p map[string]any, ctxID string) (bool, string, error) {
	kind, _ := p["kind"].(string)
	link, _ := p["link"].(string)
	filter, _ := p["filter"].(map[string]any)
	for _, r := range matchRecords(env, kind, filter) {
		if s, _ := r.Data[link].(string); s == ctxID && ctxID != "" {
			return true, "", nil
		}
	}
	return false, fmt.Sprintf("no %s linked to %s=%s", kind, link, ctxID), nil
}

func evalVerdictIn(env *Env, p map[string]any) (bool, string, error) {
	agent, _ := p["agent"].(string)
	scope, _ := p["scope"].(string)
	statuses := strList(p["statuses"])
	riskyOK, _ := p["risky_with_dispositions"].(bool)

	eff := effectiveVerdict(env, agent, scope)
	if eff == nil {
		if scope != "" {
			return false, fmt.Sprintf("no %s verdict (scope %s)", agent, scope), nil
		}
		return false, fmt.Sprintf("no %s verdict", agent), nil
	}
	// A manual verdict for a GATING agent is refused while hook capture is
	// provably alive (any auto verdict exists in the run): stickiness only
	// protected existing auto fails — a never-spawned reviewer could be
	// "recorded clean" by hand. Degraded-hooks sessions (no auto verdicts
	// at all) keep working; the recorded escape is a disposition
	// referencing the manual verdict (user-visible, auditor-judged).
	if !eff.Auto && env.gatingAgent(agent) && env.hasAutoVerdict() && !env.agentVerdictDispositioned(agent, scope) {
		return false, fmt.Sprintf("manual (self-attested) %s verdict while verdict capture is alive — re-spawn @wf:%s (auto-captured), or wf record disposition --ref %s --text \"why the manual record stands\"", agent, agent, eff.ID), nil
	}
	status, _ := eff.Data["status"].(string)
	for _, s := range statuses {
		if s == status {
			return true, "", nil
		}
	}
	if status == "risky" && riskyOK {
		for _, d := range env.Records("disposition") {
			if ref, _ := d.Data["ref"].(string); ref == eff.ID {
				return true, "", nil
			}
		}
		return false, fmt.Sprintf("%s verdict risky without recorded dispositions (wf record disposition --ref %s)", agent, eff.ID), nil
	}
	return false, fmt.Sprintf("%s verdict is %q (need %s)", agent, status, strings.Join(statuses, "|")), nil
}

// effectiveVerdict applies the sticky-auto-evidence rule (04 §4): a manual
// record cannot supersede an auto-captured failing verdict; only a newer auto
// capture or an explicit disposition referencing it can.
func effectiveVerdict(env *Env, agent, scope string) *Record {
	var all []Record
	for _, r := range env.Records("verdict") {
		if a, _ := r.Data["agent"].(string); a != agent {
			continue
		}
		if scope != "" {
			if s, _ := r.Data["scope"].(string); s != scope {
				continue
			}
		}
		all = append(all, r)
	}
	if len(all) == 0 {
		return nil
	}
	last := all[len(all)-1]
	if last.Auto {
		return &last
	}
	// last is manual: does an auto failing verdict exist after which no auto
	// record has run?
	var lastAuto *Record
	for i := range all {
		if all[i].Auto {
			lastAuto = &all[i]
		}
	}
	if lastAuto == nil {
		return &last
	}
	status, _ := lastAuto.Data["status"].(string)
	if !env.isFailStatus(status) {
		return &last
	}
	// sticky unless dispositioned
	if env.dispositioned(lastAuto.ID) {
		return &last
	}
	return lastAuto
}

func (e *Env) gatingAgent(name string) bool {
	if e.Spec == nil {
		return false
	}
	a, ok := e.Spec.AgentByName(name)
	return ok && a.Gating
}

func (e *Env) hasAutoVerdict() bool {
	for _, v := range e.Records("verdict") {
		if v.Auto {
			return true
		}
	}
	return false
}

// agentVerdictDispositioned: some disposition references one of this
// agent's (scope-matched) verdicts — the user-visible override for both
// sticky auto fails and manual-while-capture-alive records.
func (e *Env) agentVerdictDispositioned(agent, scope string) bool {
	for _, v := range e.Records("verdict") {
		if a, _ := v.Data["agent"].(string); a != agent {
			continue
		}
		if scope != "" {
			if s, _ := v.Data["scope"].(string); s != scope {
				continue
			}
		}
		if e.dispositioned(v.ID) {
			return true
		}
	}
	return false
}

func (e *Env) dispositioned(ref string) bool {
	for _, d := range e.Records("disposition") {
		if r, _ := d.Data["ref"].(string); r == ref {
			return true
		}
	}
	return false
}

func (e *Env) isFailStatus(s string) bool {
	for _, f := range e.Spec.Verdicts.Fail {
		if f == s {
			return true
		}
	}
	return false
}

func evalApproval(env *Env, p map[string]any) (bool, string, error) {
	gate, _ := p["gate"].(string)
	for _, r := range env.Records("approval") {
		if g, _ := r.Data["gate"].(string); g == gate {
			return true, "", nil
		}
	}
	return false, fmt.Sprintf("approval %q not recorded", gate), nil
}

func evalNoOpen(env *Env, p map[string]any) (bool, string, error) {
	kind, _ := p["kind"].(string)
	field, _ := p["field"].(string)
	open := strList(p["open"])
	filter, _ := p["filter"].(map[string]any)
	var stuck []string
	for _, r := range matchRecords(env, kind, filter) {
		v := fmt.Sprintf("%v", r.Data[field])
		for _, o := range open {
			if v == o {
				stuck = append(stuck, label(r))
			}
		}
	}
	if len(stuck) == 0 {
		return true, "", nil
	}
	return false, fmt.Sprintf("open %s: %s", kind, strings.Join(dedupe(stuck), ", ")), nil
}

func evalPerEach(env *Env, p map[string]any) (bool, string, error) {
	kind, _ := p["kind"].(string)
	each, _ := p["each"].(string)
	filter, _ := p["filter"].(map[string]any)
	item, _ := p["item"].(map[string]any)
	min := intParam(p, "min", 0)
	pred, prm, err := spec.SubItem(item)
	if err != nil {
		return false, "", err
	}
	var missing []string
	total := 0
	for _, r := range matchRecords(env, kind, filter) {
		ids := elementIDs(r, each)
		total += len(ids)
		for _, id := range ids {
			ok, _, err := evalPredicate(env, pred, prm, id)
			if err != nil {
				return false, "", err
			}
			if !ok {
				missing = append(missing, id)
			}
		}
	}
	// `min` guards the vacuous-pass edge: zero elements used to satisfy
	// per-each trivially, letting an AC-less requirement dodge every
	// downstream per-AC item.
	if total < min {
		return false, fmt.Sprintf("only %d element(s) to check (min %d)", total, min), nil
	}
	if len(missing) == 0 {
		return true, "", nil
	}
	return false, "unmet for: " + strings.Join(dedupe(missing), ", "), nil
}

// elementIDs yields the per-each element ids of one record: the record's own
// identity when `each` is empty, else the ids of the sub-array items.
func elementIDs(r Record, each string) []string {
	if each == "" {
		return []string{recordIdentity(r)}
	}
	raw, ok := r.Data[each].([]any)
	if !ok {
		return nil
	}
	var out []string
	for _, el := range raw {
		switch t := el.(type) {
		case string:
			out = append(out, t)
		case map[string]any:
			if id, _ := t["id"].(string); id != "" {
				out = append(out, id)
			}
		}
	}
	return out
}

// recordIdentity prefers a domain key (tid, rid, ac) over the event id so
// remediations speak the agent's language.
func recordIdentity(r Record) string {
	for _, k := range []string{"tid", "rid", "ac", "id"} {
		if s, _ := r.Data[k].(string); s != "" {
			return s
		}
	}
	return r.ID
}

func evalAnyOf(env *Env, p map[string]any, ctxID string) (bool, string, error) {
	items, _ := p["items"].([]any)
	var details []string
	for _, raw := range items {
		m, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		pred, prm, err := spec.SubItem(m)
		if err != nil {
			return false, "", err
		}
		ok2, d, err := evalPredicate(env, pred, prm, ctxID)
		if err != nil {
			return false, "", err
		}
		if ok2 {
			return true, "", nil
		}
		if d != "" {
			details = append(details, d)
		}
	}
	return false, strings.Join(details, "; "), nil
}

// red→green pairing classes, ordered by strength.
const (
	rgNoRed = iota
	rgNoGreen
	rgCrossRunner
	rgWeak
	rgStrict
)

// redGreenClass classifies the strongest red→green evidence pair for one
// per-each id. Order alone used to satisfy the predicate — any red before
// any green — letting a gitleaks red pair with a pytest green. Pairing now
// requires the same runner: strict when one command is a whole-token prefix
// of the other (red on a selector, green on the full suite), weak when the
// runner matches but selectors diverge (passes, surfaced by wf report).
func redGreenClass(env *Env, link, ctxID string) (cls int, lastRedCmd string) {
	type tr struct {
		order int
		cmd   string
	}
	var reds, greens []tr
	for _, r := range env.Records("test-run") {
		if s, _ := r.Data[link].(string); s != ctxID || ctxID == "" {
			continue
		}
		if g, _ := r.Data["grounded"].(bool); !g {
			continue
		}
		exit, hasExit := numValue(r.Data["exit"])
		if !hasExit {
			continue
		}
		cmd, _ := r.Data["cmd"].(string)
		if exit != 0 {
			reds = append(reds, tr{r.Order, cmd})
		} else {
			greens = append(greens, tr{r.Order, cmd})
		}
	}
	if len(reds) == 0 {
		return rgNoRed, ""
	}
	lastRedCmd = reds[len(reds)-1].cmd
	cls = rgNoGreen
	for _, g := range greens {
		for _, rd := range reds {
			if rd.order >= g.order {
				continue
			}
			if cls < rgCrossRunner {
				cls = rgCrossRunner
			}
			switch cmdid.Classify(rd.cmd, g.cmd) {
			case cmdid.PairStrict:
				return rgStrict, lastRedCmd
			case cmdid.PairWeak:
				cls = rgWeak
			}
		}
	}
	return cls, lastRedCmd
}

func evalRedGreen(env *Env, p map[string]any, ctxID string) (bool, string, error) {
	link, _ := p["link"].(string)
	cls, redCmd := redGreenClass(env, link, ctxID)
	switch cls {
	case rgStrict, rgWeak:
		return true, "", nil
	case rgNoRed:
		return false, fmt.Sprintf("no failing (red) test-run tagged %s=%s", link, ctxID), nil
	case rgCrossRunner:
		return false, fmt.Sprintf("red (%s) and the later green(s) use different runners — re-run the red's runner to green (a full-suite green counts)", redCmd), nil
	default:
		return false, fmt.Sprintf("red test-run recorded for %s but no later green", ctxID), nil
	}
}

// ApprovalRefs computes what a content-bearing approval gate binds to at
// this moment: the record identities a user saying "yes" is saying yes TO.
// Stored on the approval at wf approve time (engine-computed — the agent-
// supplied --payload proved nothing about what was actually presented);
// re-computed at trace time to detect approval drift (records added after
// the approval without re-approval).
func ApprovalRefs(env *Env, gate string) []string {
	var refs []string
	switch gate {
	case "scope":
		for _, r := range env.Records("requirement") {
			if s, _ := r.Data["status"].(string); s == "dropped" {
				continue
			}
			if rid, _ := r.Data["rid"].(string); rid != "" {
				refs = append(refs, "req:"+rid)
			}
		}
		for _, a := range env.Records("assumption") {
			if hr, _ := a.Data["high_risk"].(bool); hr {
				refs = append(refs, "assumption:"+a.ID)
			}
		}
	case "design":
		for _, os := range env.Records("option-set") {
			stage, _ := os.Data["stage"].(string)
			if sel, _ := os.Data["selected"].(string); sel != "" {
				refs = append(refs, "option:"+stage+":"+sel)
			}
		}
	case "plan":
		for _, t := range env.Records("task") {
			if tid, _ := t.Data["tid"].(string); tid != "" {
				refs = append(refs, "task:"+tid)
			}
		}
		for _, vs := range env.Records("verification-strategy") {
			if ac, _ := vs.Data["ac"].(string); ac != "" {
				refs = append(refs, "strategy:"+ac)
			}
		}
	}
	refs = dedupe(refs)
	sort.Strings(refs)
	return refs
}

// WeakRedGreenTasks lists task tids whose test-first evidence pairs only
// weakly (same runner, diverging selectors) — a wf report health signal and
// an auditor probe point, never a block.
func WeakRedGreenTasks(env *Env) []string {
	var out []string
	for _, task := range env.Records("task") {
		tid, _ := task.Data["tid"].(string)
		if tid == "" {
			continue
		}
		if cls, _ := redGreenClass(env, "task", tid); cls == rgWeak {
			out = append(out, tid)
		}
	}
	return dedupe(out)
}

// ---------------------------------------------------------------------------
// Matching helpers
// ---------------------------------------------------------------------------

func matchRecords(env *Env, kind string, filter map[string]any) []Record {
	var out []Record
	for _, r := range env.Records(kind) {
		if matchesFilter(r.Data, filter) {
			out = append(out, r)
		}
	}
	return out
}

func matchesFilter(data, filter map[string]any) bool {
	for k, want := range filter {
		got, present := data[k]
		if !present {
			return false
		}
		if list, ok := want.([]any); ok { // "in" semantics
			hit := false
			for _, w := range list {
				if valueEq(got, w) {
					hit = true
					break
				}
			}
			if !hit {
				return false
			}
			continue
		}
		if !valueEq(got, want) {
			return false
		}
	}
	return true
}

// valueEq compares JSON-decoded (float64) and YAML-decoded (int/bool/string)
// scalars.
func valueEq(got, want any) bool {
	if gn, ok := numValue(got); ok {
		if wn, ok2 := numValue(want); ok2 {
			return gn == wn
		}
	}
	if gb, ok := got.(bool); ok {
		if wb, ok2 := want.(bool); ok2 {
			return gb == wb
		}
	}
	return fmt.Sprintf("%v", got) == fmt.Sprintf("%v", want)
}

func numValue(v any) (float64, bool) {
	switch t := v.(type) {
	case float64:
		return t, true
	case int:
		return float64(t), true
	case int64:
		return float64(t), true
	}
	return 0, false
}

func intParam(p map[string]any, key string, def int) int {
	if v, ok := numValue(p[key]); ok {
		return int(v)
	}
	return def
}

func strList(v any) []string {
	raw, _ := v.([]any)
	out := make([]string, 0, len(raw))
	for _, x := range raw {
		if s, ok := x.(string); ok {
			out = append(out, s)
		}
	}
	return out
}

func label(r Record) string { return recordIdentity(r) }

func dedupe(in []string) []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range in {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}
