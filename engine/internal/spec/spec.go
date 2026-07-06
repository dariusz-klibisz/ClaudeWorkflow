// Package spec loads and validates the declarative workflow specification
// (workflow/workflow.yaml) and merges add-only project-local extensions from
// .workflow/contracts.d/. The spec is THE single source of truth for phases,
// families, the gating roster, record kinds, and contract items
// (workflow-redesign/03 §4.0).
package spec

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	"gopkg.in/yaml.v3"
)

// SchemaVersion is the spec schema this engine reads.
const SchemaVersion = 1

type Spec struct {
	Schema    int                 `yaml:"schema"`
	Families  []string            `yaml:"families"`
	Intents   map[string][]string `yaml:"intents"`
	Verdicts  Verdicts            `yaml:"verdicts"`
	Phases    []Phase             `yaml:"phases"`
	Loops     Loops               `yaml:"loops"`
	Roster    []Agent             `yaml:"roster"`
	Records   []RecordKind        `yaml:"records"`
	Contracts []ContractItem      `yaml:"contracts"`

	pluginRoot string // dir containing workflow/, reference/, agents/ …
}

// PluginRoot is the plugin directory the spec was loaded from.
func (s *Spec) PluginRoot() string { return s.pluginRoot }

type Verdicts struct {
	Pass        []string `yaml:"pass"`
	Fail        []string `yaml:"fail"`
	Conditional []string `yaml:"conditional"`
}

type Phase struct {
	ID            string            `yaml:"id"`
	Order         int               `yaml:"order"`
	Mode          string            `yaml:"mode"` // interactive | interactive-exit | auto-advance
	Skill         string            `yaml:"skill"`
	Participation map[string]string `yaml:"participation"` // family -> required|waivable|none
}

type Loops struct {
	From         string   `yaml:"from"`
	Targets      []string `yaml:"targets"`
	MaxPerRun    int      `yaml:"max_per_run"`
	MaxSlipPerAC int      `yaml:"max_slip_per_ac"`
}

type Agent struct {
	Name     string   `yaml:"name"`
	Gating   bool     `yaml:"gating"`
	Phases   []string `yaml:"phases"`
	Model    string   `yaml:"model"`
	Tools    []string `yaml:"tools"`
	MaxTurns int      `yaml:"maxTurns"`
	Memory   string   `yaml:"memory"`
	When     *When    `yaml:"when"`
	// Corpus lists the bundled reference files/dirs this agent is routed to
	// (plugin-root-relative). Validated to exist; injected at SubagentStart.
	Corpus []string `yaml:"corpus"`
}

type RecordKind struct {
	Kind   string   `yaml:"kind"`
	Fields []string `yaml:"fields"` // trailing '?' marks optional
}

// Required returns the required field names (without '?').
func (r RecordKind) Required() []string {
	var out []string
	for _, f := range r.Fields {
		if !strings.HasSuffix(f, "?") {
			out = append(out, f)
		}
	}
	return out
}

// Known returns all field names (without '?').
func (r RecordKind) Known() []string {
	out := make([]string, 0, len(r.Fields))
	for _, f := range r.Fields {
		out = append(out, strings.TrimSuffix(f, "?"))
	}
	return out
}

// When gates whether a contract item (or roster entry) applies.
type When struct {
	Intents []string       `yaml:"intents"`
	Config  map[string]any `yaml:"config"`
	Signals []string       `yaml:"signals"`
	Records *RecordMatch   `yaml:"records"`
}

type RecordMatch struct {
	Kind   string         `yaml:"kind"`
	Filter map[string]any `yaml:"filter"`
}

// Predicate names — the closed vocabulary.
const (
	PredRecordExists = "record-exists"
	PredLinkedRecord = "linked-record"
	PredVerdictIn    = "verdict-in"
	PredApproval     = "approval"
	PredNoOpen       = "no-open"
	PredPerEach      = "per-each"
	PredAnyOf        = "any-of"
	PredRedGreen     = "red-green"
)

var predicates = map[string]bool{
	PredRecordExists: true, PredLinkedRecord: true, PredVerdictIn: true,
	PredApproval: true, PredNoOpen: true, PredPerEach: true,
	PredAnyOf: true, PredRedGreen: true,
}

type ContractItem struct {
	ID          string         `yaml:"id"`
	Phase       string         `yaml:"phase"`
	Families    []string       `yaml:"families"` // empty = all
	When        *When          `yaml:"when"`
	Predicate   string         `yaml:"predicate"`
	Params      map[string]any `yaml:"params"`
	Waivable    bool           `yaml:"waivable"`
	Remediation string         `yaml:"remediation"`
	Source      string         `yaml:"-"` // "spec" or the contracts.d file it came from
}

// SubItem extracts a nested predicate spec from params (per-each `item`,
// any-of `items` elements).
func SubItem(m map[string]any) (pred string, params map[string]any, err error) {
	p, _ := m["predicate"].(string)
	if p == "" {
		return "", nil, fmt.Errorf("nested item missing predicate")
	}
	prm := map[string]any{}
	if raw, ok := m["params"].(map[string]any); ok {
		prm = raw
	}
	return p, prm, nil
}

// ---------------------------------------------------------------------------
// Loading
// ---------------------------------------------------------------------------

// Load reads the plugin spec and merges any .workflow/contracts.d/*.yaml
// extensions (add-only). specPath is workflow/workflow.yaml inside the plugin;
// contractsDir may be "" (no project extensions).
//
// Runtime loading is TOLERANT of unknown fields: additive spec changes must
// never break an older engine within the same schema version (02 §5 — the
// lesson of the corpus-field incident). Typos are caught at dev/CI time by
// LoadStrict (gen -check, tests, doctor).
func Load(specPath, contractsDir string) (*Spec, error) {
	return load(specPath, contractsDir, false)
}

// LoadStrict rejects unknown fields — the dev/CI-side parse.
func LoadStrict(specPath, contractsDir string) (*Spec, error) {
	return load(specPath, contractsDir, true)
}

func load(specPath, contractsDir string, strict bool) (*Spec, error) {
	raw, err := os.ReadFile(specPath)
	if err != nil {
		return nil, fmt.Errorf("read spec: %w", err)
	}
	var s Spec
	dec := yaml.NewDecoder(strings.NewReader(string(raw)))
	dec.KnownFields(strict)
	if err := dec.Decode(&s); err != nil {
		return nil, fmt.Errorf("parse spec: %w", err)
	}
	if abs, err := filepath.Abs(specPath); err == nil {
		s.pluginRoot = filepath.Dir(filepath.Dir(abs))
	}
	for i := range s.Contracts {
		s.Contracts[i].Source = "spec"
	}
	if contractsDir != "" {
		if err := s.mergeContractsDir(contractsDir, strict); err != nil {
			return nil, err
		}
	}
	if err := s.Validate(); err != nil {
		return nil, err
	}
	return &s, nil
}

// extension is the shape of a contracts.d file: additional record kinds and
// contract items only.
type extension struct {
	Records   []RecordKind   `yaml:"records"`
	Contracts []ContractItem `yaml:"contracts"`
}

func (s *Spec) mergeContractsDir(dir string, strict bool) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("read contracts.d: %w", err)
	}
	names := make([]string, 0, len(entries))
	for _, e := range entries {
		if !e.IsDir() && (strings.HasSuffix(e.Name(), ".yaml") || strings.HasSuffix(e.Name(), ".yml")) {
			names = append(names, e.Name())
		}
	}
	sort.Strings(names)
	shippedIDs := map[string]bool{}
	for _, c := range s.Contracts {
		shippedIDs[c.ID] = true
	}
	shippedKinds := map[string]bool{}
	for _, r := range s.Records {
		shippedKinds[r.Kind] = true
	}
	for _, name := range names {
		raw, err := os.ReadFile(filepath.Join(dir, name))
		if err != nil {
			return fmt.Errorf("read %s: %w", name, err)
		}
		var ext extension
		dec := yaml.NewDecoder(strings.NewReader(string(raw)))
		dec.KnownFields(strict)
		if err := dec.Decode(&ext); err != nil {
			return fmt.Errorf("parse %s: %w", name, err)
		}
		for _, r := range ext.Records {
			if !strings.HasPrefix(r.Kind, "x-") {
				return fmt.Errorf("%s: project record kind %q must be namespaced x-*", name, r.Kind)
			}
			if shippedKinds[r.Kind] {
				return fmt.Errorf("%s: record kind %q already defined (add-only merge)", name, r.Kind)
			}
			shippedKinds[r.Kind] = true
			s.Records = append(s.Records, r)
		}
		for _, c := range ext.Contracts {
			if !strings.HasPrefix(c.ID, "local.") && !strings.HasPrefix(c.ID, "lesson.") {
				return fmt.Errorf("%s: project contract %q must be namespaced local.* or lesson.*", name, c.ID)
			}
			if shippedIDs[c.ID] {
				return fmt.Errorf("%s: contract %q already defined (add-only merge)", name, c.ID)
			}
			shippedIDs[c.ID] = true
			c.Source = name
			s.Contracts = append(s.Contracts, c)
		}
	}
	return nil
}

// ---------------------------------------------------------------------------
// Validation — the meta-validation that replaces v0.36's parity meta-tests:
// every reference must resolve inside this one artifact.
// ---------------------------------------------------------------------------

func (s *Spec) Validate() error {
	var errs []string
	fail := func(format string, a ...any) { errs = append(errs, fmt.Sprintf(format, a...)) }

	if s.Schema != SchemaVersion {
		fail("schema %d not supported (engine reads %d)", s.Schema, SchemaVersion)
	}
	if len(s.Families) == 0 {
		fail("no families declared")
	}
	fams := map[string]bool{}
	for _, f := range s.Families {
		fams[f] = true
	}
	for fam := range s.Intents {
		if !fams[fam] {
			fail("intents declared for unknown family %q", fam)
		}
	}

	// Phases: unique ids, contiguous order, valid modes, full participation.
	phaseIDs := map[string]bool{}
	orders := map[int]bool{}
	for _, p := range s.Phases {
		if phaseIDs[p.ID] {
			fail("duplicate phase %q", p.ID)
		}
		phaseIDs[p.ID] = true
		if orders[p.Order] {
			fail("duplicate phase order %d", p.Order)
		}
		orders[p.Order] = true
		switch p.Mode {
		case "interactive", "interactive-exit", "auto-advance":
		default:
			fail("phase %q: unknown mode %q", p.ID, p.Mode)
		}
		if p.Skill == "" {
			fail("phase %q: no skill", p.ID)
		}
		for _, f := range s.Families {
			part, ok := p.Participation[f]
			if !ok {
				fail("phase %q: no participation entry for family %q", p.ID, f)
				continue
			}
			switch part {
			case "required", "waivable", "none":
			default:
				fail("phase %q family %q: participation %q invalid", p.ID, f, part)
			}
		}
	}
	for i := 1; i <= len(s.Phases); i++ {
		if !orders[i] {
			fail("phase orders not contiguous: missing %d", i)
		}
	}

	// Loops reference real phases.
	if s.Loops.From != "" && !phaseIDs[s.Loops.From] {
		fail("loops.from %q not a phase", s.Loops.From)
	}
	for _, t := range s.Loops.Targets {
		if !phaseIDs[t] {
			fail("loop target %q not a phase", t)
		}
	}

	// Roster: unique names, real phases, resolvable corpus routes.
	agents := map[string]bool{}
	for _, a := range s.Roster {
		if agents[a.Name] {
			fail("duplicate agent %q", a.Name)
		}
		agents[a.Name] = true
		for _, p := range a.Phases {
			if !phaseIDs[p] {
				fail("agent %q: unknown phase %q", a.Name, p)
			}
		}
		// Corpus paths are checked only when the plugin actually bundles a
		// reference/ tree (spec copies in tests/dev sandboxes skip this; the
		// real repo is covered by CI's gen -check from the plugin root).
		if s.pluginRoot != "" {
			if _, err := os.Stat(filepath.Join(s.pluginRoot, "reference")); err == nil {
				for _, c := range a.Corpus {
					if _, err := os.Stat(filepath.Join(s.pluginRoot, filepath.FromSlash(c))); err != nil {
						fail("agent %q: corpus path %q not found in plugin", a.Name, c)
					}
				}
			}
		}
	}

	// Record kinds: unique.
	kinds := map[string]bool{}
	for _, r := range s.Records {
		if kinds[r.Kind] {
			fail("duplicate record kind %q", r.Kind)
		}
		kinds[r.Kind] = true
	}

	// Verdict vocabulary sanity.
	if len(s.Verdicts.Pass) == 0 || len(s.Verdicts.Fail) == 0 {
		fail("verdict vocabulary incomplete")
	}
	verdictOK := map[string]bool{}
	for _, v := range append(append(append([]string{}, s.Verdicts.Pass...), s.Verdicts.Fail...), s.Verdicts.Conditional...) {
		verdictOK[v] = true
	}

	// Contracts: unique ids, resolvable references, valid predicates.
	ids := map[string]bool{}
	for _, c := range s.Contracts {
		if ids[c.ID] {
			fail("duplicate contract id %q", c.ID)
		}
		ids[c.ID] = true
		if !phaseIDs[c.Phase] {
			fail("contract %q: unknown phase %q", c.ID, c.Phase)
		}
		for _, f := range c.Families {
			if !fams[f] {
				fail("contract %q: unknown family %q", c.ID, f)
			}
		}
		if c.Remediation == "" {
			fail("contract %q: remediation required", c.ID)
		}
		if c.When != nil && c.When.Records != nil && !kinds[c.When.Records.Kind] {
			fail("contract %q: when.records kind %q undeclared", c.ID, c.When.Records.Kind)
		}
		s.validatePredicate(c.ID, c.Predicate, c.Params, kinds, agents, verdictOK, fail)
	}

	if len(errs) > 0 {
		return fmt.Errorf("spec invalid:\n  - %s", strings.Join(errs, "\n  - "))
	}
	return nil
}

func (s *Spec) validatePredicate(id, pred string, params map[string]any, kinds, agents, verdictOK map[string]bool, fail func(string, ...any)) {
	if !predicates[pred] {
		fail("contract %q: unknown predicate %q", id, pred)
		return
	}
	kindOf := func(key string) {
		if k, ok := params[key].(string); ok {
			if !kinds[k] {
				fail("contract %q: undeclared record kind %q", id, k)
			}
		} else {
			fail("contract %q: predicate %s requires string param %q", id, pred, key)
		}
	}
	switch pred {
	case PredRecordExists, PredNoOpen:
		kindOf("kind")
		if pred == PredNoOpen {
			if _, ok := params["field"].(string); !ok {
				fail("contract %q: no-open requires field", id)
			}
			if _, ok := params["open"].([]any); !ok {
				fail("contract %q: no-open requires open list", id)
			}
		}
	case PredLinkedRecord:
		kindOf("kind")
		if _, ok := params["link"].(string); !ok {
			fail("contract %q: linked-record requires link", id)
		}
	case PredVerdictIn:
		agent, _ := params["agent"].(string)
		if !agents[agent] {
			fail("contract %q: undeclared agent %q", id, agent)
		}
		raw, ok := params["statuses"].([]any)
		if !ok || len(raw) == 0 {
			fail("contract %q: verdict-in requires statuses", id)
		}
		for _, v := range raw {
			vs, _ := v.(string)
			if !verdictOK[vs] {
				fail("contract %q: status %q not in verdict vocabulary", id, vs)
			}
		}
	case PredApproval:
		if _, ok := params["gate"].(string); !ok {
			fail("contract %q: approval requires gate", id)
		}
	case PredRedGreen:
		if _, ok := params["link"].(string); !ok {
			fail("contract %q: red-green requires link", id)
		}
	case PredPerEach:
		kindOf("kind")
		item, ok := params["item"].(map[string]any)
		if !ok {
			fail("contract %q: per-each requires item", id)
			return
		}
		p, prm, err := SubItem(item)
		if err != nil {
			fail("contract %q: %v", id, err)
			return
		}
		s.validatePredicate(id, p, prm, kinds, agents, verdictOK, fail)
	case PredAnyOf:
		items, ok := params["items"].([]any)
		if !ok || len(items) == 0 {
			fail("contract %q: any-of requires items", id)
			return
		}
		for _, it := range items {
			m, ok := it.(map[string]any)
			if !ok {
				fail("contract %q: any-of item malformed", id)
				continue
			}
			p, prm, err := SubItem(m)
			if err != nil {
				fail("contract %q: %v", id, err)
				continue
			}
			// `waiver` linked via the synthetic "item" link is allowed.
			s.validatePredicate(id, p, prm, kinds, agents, verdictOK, fail)
		}
	}
}

// ---------------------------------------------------------------------------
// Lookups
// ---------------------------------------------------------------------------

func (s *Spec) PhaseByID(id string) (Phase, bool) {
	for _, p := range s.Phases {
		if p.ID == id {
			return p, true
		}
	}
	return Phase{}, false
}

// PhasesFor returns the ordered phases a family participates in
// (participation != none).
func (s *Spec) PhasesFor(family string) []Phase {
	out := make([]Phase, 0, len(s.Phases))
	for _, p := range s.Phases {
		if p.Participation[family] != "none" {
			out = append(out, p)
		}
	}
	sort.Slice(out, func(i, j int) bool { return out[i].Order < out[j].Order })
	return out
}

// NextPhase returns the family's phase after cur, or "" when cur is last.
func (s *Spec) NextPhase(family, cur string) string {
	ph := s.PhasesFor(family)
	for i, p := range ph {
		if p.ID == cur && i+1 < len(ph) {
			return ph[i+1].ID
		}
	}
	return ""
}

// ContractsFor returns the contract items applying to (phase, family).
// `when` conditions are evaluated later, by the evaluator, against run state.
func (s *Spec) ContractsFor(phase, family string) []ContractItem {
	var out []ContractItem
	for _, c := range s.Contracts {
		if c.Phase != phase {
			continue
		}
		if len(c.Families) > 0 {
			ok := false
			for _, f := range c.Families {
				if f == family {
					ok = true
					break
				}
			}
			if !ok {
				continue
			}
		}
		out = append(out, c)
	}
	return out
}

func (s *Spec) RecordKind(kind string) (RecordKind, bool) {
	for _, r := range s.Records {
		if r.Kind == kind {
			return r, true
		}
	}
	return RecordKind{}, false
}

func (s *Spec) AgentByName(name string) (Agent, bool) {
	for _, a := range s.Roster {
		if a.Name == name {
			return a, true
		}
	}
	return Agent{}, false
}

func (s *Spec) GatingAgents() []Agent {
	var out []Agent
	for _, a := range s.Roster {
		if a.Gating {
			out = append(out, a)
		}
	}
	return out
}

func (s *Spec) ValidFamily(f string) bool {
	for _, x := range s.Families {
		if x == f {
			return true
		}
	}
	return false
}

func (s *Spec) ValidIntent(family, intent string) bool {
	for _, x := range s.Intents[family] {
		if x == intent {
			return true
		}
	}
	return false
}
