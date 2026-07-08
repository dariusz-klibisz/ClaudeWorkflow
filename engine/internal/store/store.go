// Package store owns all on-disk workflow state under .workflow/:
// the append-only event log, the derived run snapshot, project config, and
// the per-machine local dir. Guarantees: ULID event identity (merge-safe
// across branches/machines), atomic snapshot writes (temp+rename),
// single-writer lockfile, snapshot always re-derivable from the log
// (workflow-redesign/08).
package store

import (
	"bufio"
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"
)

const (
	StateSchema = 1
	DirName     = ".workflow"
)

// Event is the envelope for every record and engine transition (08 §3).
type Event struct {
	Schema int    `json:"schema"`
	ID     string `json:"id"`
	Seq    int64  `json:"seq"` // per-writer ordering hint; ties break on ID
	// Prev chains each log line to the raw bytes of the line before it
	// (sha256/16). Tamper-EVIDENCE, not cryptographic integrity: an
	// out-of-band editor must recompute every later line, which wf doctor
	// makes visible. Engine rewrites (compaction, archive) re-anchor the
	// chain legitimately.
	Prev  string         `json:"prev,omitempty"`
	TS    string         `json:"ts"`
	Run   string         `json:"run,omitempty"`
	Phase string         `json:"phase,omitempty"`
	Kind  string         `json:"kind"`
	Auto  bool           `json:"auto"`
	Actor string         `json:"actor"` // agent | engine | hook | user
	Note  string         `json:"note,omitempty"`
	Data  map[string]any `json:"data,omitempty"`
}

// Str returns a string field from Data.
func (e Event) Str(key string) string {
	v, _ := e.Data[key].(string)
	return v
}

// Run is the derived current-run snapshot (state/run.json) for O(1) gate reads.
type Run struct {
	Schema  int    `json:"schema"`
	ID      string `json:"id"`
	Family  string `json:"family"`
	Intent  string `json:"intent"`
	Phase   string `json:"phase"`
	Status  string `json:"status"` // active | parked | closed
	Started string `json:"started"`
	Parent  string `json:"parent,omitempty"`
	// counters maintained by runctl
	Loops    int            `json:"loops"`
	SlipByAC map[string]int `json:"slip_by_ac,omitempty"`
	Forces   int            `json:"forces"`
	WaivedPh []string       `json:"waived_phases,omitempty"`
	ExitedPh []string       `json:"exited_phases,omitempty"`
}

// Config is .workflow/config.json (08 §2).
type Config struct {
	Schema        int            `json:"schema"`
	PluginVersion string         `json:"plugin_version,omitempty"`
	UX            bool           `json:"ux"`
	Thresholds    map[string]any `json:"thresholds,omitempty"`
	Flags         map[string]any `json:"flags,omitempty"`
	// Runners: extra test-runner heads for Bash test capture — custom
	// wrappers (`./scripts/test.sh`) no static list or strategy-learning
	// heuristic can recognize.
	Runners []string `json:"runners,omitempty"`
}

// ConfigFlag exposes config values for `when.config` evaluation.
func (c *Config) ConfigFlag(key string) any {
	switch key {
	case "ux":
		return c.UX
	case "thresholds":
		return len(c.Thresholds) > 0
	}
	if c.Flags != nil {
		if v, ok := c.Flags[key]; ok {
			return v
		}
	}
	return nil
}

type Store struct {
	Root string // the .workflow directory
	seq  int64
	lock *os.File
}

// Open binds a store to <projectDir>/.workflow, creating the layout if init
// is true. Returns ErrNotInitialized when the dir is absent and init false.
var ErrNotInitialized = errors.New("workflow state not initialized (run /wf:init)")

func Open(projectDir string, init bool) (*Store, error) {
	root := filepath.Join(projectDir, DirName)
	if _, err := os.Stat(root); err != nil {
		if !init {
			return nil, ErrNotInitialized
		}
	}
	// Never mix state with a legacy ClaudeInit scaffold: the old generator
	// used the same .workflow/ directory with an incompatible layout,
	// identified by its manifest.json. Adoption refuses loudly; every other
	// path treats the project as not-ours (gates stay silent).
	if _, err := os.Stat(filepath.Join(root, "manifest.json")); err == nil {
		if init {
			return nil, fmt.Errorf("%s contains a legacy ClaudeInit scaffold (manifest.json); remove or rename the old tree before adopting wf", root)
		}
		return nil, ErrNotInitialized
	}
	if init {
		for _, d := range []string{"state", "log", "runs", "local", "contracts.d"} {
			if err := os.MkdirAll(filepath.Join(root, d), 0o755); err != nil {
				return nil, err
			}
		}
		gi := filepath.Join(root, ".gitignore")
		if _, err := os.Stat(gi); os.IsNotExist(err) {
			_ = os.WriteFile(gi, []byte("state/lock\nlocal/\n"), 0o644)
		}
	}
	return &Store{Root: root}, nil
}

func (s *Store) path(parts ...string) string {
	return filepath.Join(append([]string{s.Root}, parts...)...)
}

func (s *Store) EventsPath() string        { return s.path("log", "events.jsonl") }
func (s *Store) RunPath() string           { return s.path("state", "run.json") }
func (s *Store) ConfigPath() string        { return s.path("config.json") }
func (s *Store) LocalPath(n string) string { return s.path("local", n) }
func (s *Store) ContractsDir() string      { return s.path("contracts.d") }
func (s *Store) RunsDir() string           { return s.path("runs") }

// ---------------------------------------------------------------------------
// Locking (single writer; gates read lock-free)
// ---------------------------------------------------------------------------

const lockStaleAfter = 30 * time.Second

func (s *Store) Lock() error {
	lp := s.path("state", "lock")
	deadline := time.Now().Add(5 * time.Second)
	for {
		f, err := os.OpenFile(lp, os.O_CREATE|os.O_EXCL|os.O_WRONLY, 0o644)
		if err == nil {
			fmt.Fprintf(f, "%d %s", os.Getpid(), time.Now().UTC().Format(time.RFC3339))
			s.lock = f
			return nil
		}
		if st, serr := os.Stat(lp); serr == nil && time.Since(st.ModTime()) > lockStaleAfter {
			_ = os.Remove(lp) // stale lock from a crashed writer
			continue
		}
		if time.Now().After(deadline) {
			return fmt.Errorf("state locked by another wf process (remove %s if stale)", lp)
		}
		time.Sleep(50 * time.Millisecond)
	}
}

func (s *Store) Unlock() {
	if s.lock != nil {
		name := s.lock.Name()
		_ = s.lock.Close()
		_ = os.Remove(name)
		s.lock = nil
	}
}

// ---------------------------------------------------------------------------
// Event log
// ---------------------------------------------------------------------------

// Append stamps identity (ULID id, seq, ts, schema) and appends the event.
func (s *Store) Append(ev *Event) error {
	if ev.Kind == "" {
		return errors.New("event kind required")
	}
	ev.Schema = StateSchema
	if ev.ID == "" {
		ev.ID = NewULID()
	}
	if ev.TS == "" {
		ev.TS = time.Now().UTC().Format(time.RFC3339)
	}
	if s.seq == 0 {
		s.seq = s.lastSeq()
	}
	s.seq++
	ev.Seq = s.seq
	if err := os.MkdirAll(filepath.Dir(s.EventsPath()), 0o755); err != nil {
		return err
	}
	f, err := os.OpenFile(s.EventsPath(), os.O_CREATE|os.O_APPEND|os.O_RDWR, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	// Heal a torn tail (crash mid-write): ensure the log ends with a newline
	// so the previous partial line stays isolated and scan-skippable.
	if st, err := f.Stat(); err == nil && st.Size() > 0 {
		buf := make([]byte, 1)
		if _, err := f.ReadAt(buf, st.Size()-1); err == nil && buf[0] != '\n' {
			if _, err := f.Write([]byte("\n")); err != nil {
				return err
			}
		}
	}
	// Chain to the previous raw line (tamper evidence; doctor verifies).
	ev.Prev = ""
	if last, ok := lastLine(f); ok {
		ev.Prev = chainHash(last)
	}
	line, err := json.Marshal(ev)
	if err != nil {
		return err
	}
	if _, err := f.Write(append(line, '\n')); err != nil {
		return err
	}
	return f.Sync()
}

// chainHash is the line-chaining hash: sha256 over the raw line bytes
// (no trailing newline), truncated to 16 hex chars.
func chainHash(line []byte) string {
	sum := sha256.Sum256(line)
	return hex.EncodeToString(sum[:8])
}

// maxLineTail bounds the backwards read for the previous line — matches the
// scanner's max token size, so any scannable line is chainable.
const maxLineTail = 4 * 1024 * 1024

// lastLine returns the last complete line of f (which is guaranteed by the
// caller to end with '\n' when non-empty). ok=false on an empty file or when
// the final line exceeds maxLineTail.
func lastLine(f *os.File) ([]byte, bool) {
	st, err := f.Stat()
	if err != nil || st.Size() == 0 {
		return nil, false
	}
	off := int64(0)
	if st.Size() > maxLineTail {
		off = st.Size() - maxLineTail
	}
	buf := make([]byte, st.Size()-off)
	if _, err := f.ReadAt(buf, off); err != nil {
		return nil, false
	}
	buf = bytes.TrimSuffix(buf, []byte("\n"))
	if i := bytes.LastIndexByte(buf, '\n'); i >= 0 {
		return buf[i+1:], true
	}
	if off > 0 {
		return nil, false // line longer than the window — unchainable
	}
	return buf, true
}

func (s *Store) lastSeq() int64 {
	var last int64
	_ = s.scan(func(ev Event) bool {
		if ev.Seq > last {
			last = ev.Seq
		}
		return true
	})
	return last
}

// Events returns all events matching filter, in file order (append order per
// writer; merged logs are re-sorted by ID by callers that need global order).
func (s *Store) Events(filter func(Event) bool) ([]Event, error) {
	var out []Event
	err := s.scan(func(ev Event) bool {
		if filter == nil || filter(ev) {
			out = append(out, ev)
		}
		return true
	})
	return out, err
}

// RunEvents returns the events of one run.
func (s *Store) RunEvents(runID string) ([]Event, error) {
	return s.Events(func(e Event) bool { return e.Run == runID })
}

// ListArchivedRuns returns the IDs of closed runs under runs/, sorted
// (run IDs are date-prefixed, so this is chronological).
func (s *Store) ListArchivedRuns() ([]string, error) {
	entries, err := os.ReadDir(s.RunsDir())
	if err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	var ids []string
	for _, e := range entries {
		if !e.IsDir() {
			continue
		}
		if _, err := os.Stat(filepath.Join(s.RunsDir(), e.Name(), "events.jsonl")); err == nil {
			ids = append(ids, e.Name())
		}
	}
	sort.Strings(ids)
	return ids, nil
}

// AllEvents returns archived + live events in chronological order (archived
// runs are date-prefixed and sorted; the live log follows). Feeds views,
// lessons regeneration, and origin discovery — NEVER gates (07 §5): gate
// reads stay O(live log).
func (s *Store) AllEvents() ([]Event, error) {
	var out []Event
	ids, err := s.ListArchivedRuns()
	if err != nil {
		return nil, err
	}
	for _, id := range ids {
		evs, err := s.ArchivedRunEvents(id)
		if err != nil {
			continue // a damaged archive must not take down the readers
		}
		out = append(out, evs...)
	}
	live, err := s.Events(nil)
	if err != nil {
		return out, err
	}
	return append(out, live...), nil
}

// ArchivedRunEvents reads a closed run's archived event slice
// (runs/<id>/events.jsonl). Log replay is allowed here — this feeds
// doctor/report only, never gates (07 §5).
func (s *Store) ArchivedRunEvents(runID string) ([]Event, error) {
	path := filepath.Join(s.RunsDir(), runID, "events.jsonl")
	if _, err := os.Stat(path); err != nil {
		return nil, fmt.Errorf("no archived run %s: %w", runID, err)
	}
	var out []Event
	err := s.scanFile(path, func(ev Event) bool {
		out = append(out, ev)
		return true
	})
	return out, err
}

func (s *Store) scan(fn func(Event) bool) error {
	return s.scanFile(s.EventsPath(), fn)
}

func (s *Store) scanFile(path string, fn func(Event) bool) error {
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	for sc.Scan() {
		line := sc.Bytes()
		if len(line) == 0 {
			continue
		}
		var ev Event
		if err := json.Unmarshal(line, &ev); err != nil {
			continue // tolerate a torn/foreign line; doctor reports it
		}
		if !fn(ev) {
			break
		}
	}
	return sc.Err()
}

// ---------------------------------------------------------------------------
// Snapshot + config (atomic JSON files)
// ---------------------------------------------------------------------------

func writeJSONAtomic(path string, v any) error {
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		return err
	}
	raw, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	tmp := path + ".tmp"
	if err := os.WriteFile(tmp, append(raw, '\n'), 0o644); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func readJSON(path string, v any) error {
	raw, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	return json.Unmarshal(raw, v)
}

// LoadRun returns the current run snapshot, or nil when no run is active.
func (s *Store) LoadRun() (*Run, error) {
	var r Run
	if err := readJSON(s.RunPath(), &r); err != nil {
		if os.IsNotExist(err) {
			return nil, nil
		}
		return nil, err
	}
	if r.ID == "" {
		return nil, nil
	}
	return &r, nil
}

func (s *Store) SaveRun(r *Run) error {
	r.Schema = StateSchema
	return writeJSONAtomic(s.RunPath(), r)
}

func (s *Store) ClearRun() error {
	err := os.Remove(s.RunPath())
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (s *Store) LoadConfig() (*Config, error) {
	var c Config
	if err := readJSON(s.ConfigPath(), &c); err != nil {
		if os.IsNotExist(err) {
			return &Config{Schema: StateSchema}, nil
		}
		return nil, err
	}
	return &c, nil
}

func (s *Store) SaveConfig(c *Config) error {
	c.Schema = StateSchema
	return writeJSONAtomic(s.ConfigPath(), c)
}

// Local reads/writes small per-machine JSON blobs (never authoritative).
func (s *Store) LoadLocal(name string, v any) error {
	err := readJSON(s.LocalPath(name), v)
	if os.IsNotExist(err) {
		return nil
	}
	return err
}

func (s *Store) SaveLocal(name string, v any) error {
	return writeJSONAtomic(s.LocalPath(name), v)
}

// ---------------------------------------------------------------------------
// Run archival (wf run close, one transaction — A5/A6 fix)
// ---------------------------------------------------------------------------

// ArchiveRun moves the run's events into runs/<id>/events.jsonl, compacts the
// live log, and clears the snapshot. The terminal event must already be
// appended by the caller (inside the same engine command).
func (s *Store) ArchiveRun(runID string) error {
	all, err := s.Events(nil)
	if err != nil {
		return err
	}
	var mine, rest []Event
	for _, e := range all {
		// Open followups, commit-origins and lesson state stay live (08 §6).
		if e.Run == runID && !keepLive(e) {
			mine = append(mine, e)
		} else {
			rest = append(rest, e)
		}
	}
	dir := filepath.Join(s.RunsDir(), runID)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return err
	}
	if err := writeEvents(filepath.Join(dir, "events.jsonl"), mine); err != nil {
		return err
	}
	// verify counts before compacting (archive → verify → compact → clear)
	arch, err := countLines(filepath.Join(dir, "events.jsonl"))
	if err != nil || arch != len(mine) {
		return fmt.Errorf("archive verification failed (%d != %d): %v", arch, len(mine), err)
	}
	if err := writeEvents(s.EventsPath()+".tmp", rest); err != nil {
		return err
	}
	if err := os.Rename(s.EventsPath()+".tmp", s.EventsPath()); err != nil {
		return err
	}
	s.PruneLocal()
	return s.ClearRun()
}

// PruneLocal clears the per-machine run-scoped counters and mirrors at run
// close — none carries meaning across runs, and two of them (tasks-mirror,
// verdict-attempts) otherwise grow forever.
func (s *Store) PruneLocal() {
	for _, n := range []string{"tasks-mirror.json", "verdict-attempts.json", "stop-gate.json", "enforce-off.json", "challenge.json"} {
		_ = os.Remove(s.LocalPath(n))
	}
}

// keepLive: only OPEN followups survive in the live log after close. Lesson
// and commit-origin events archive with their run (bounded live log — the
// "grows forever" fix): their durable forms are the regenerated lesson
// channels and the committed runs/<id>/ archives, and their readers
// (lessons regeneration, origin discovery) fold archived events back in.
func keepLive(e Event) bool {
	if e.Kind == "followup" {
		return e.Str("status") == "open" || e.Str("status") == "next-run"
	}
	return false
}

// writeEvents rewrites a log file, re-anchoring the hash chain (engine
// rewrites — compaction and archival — are the only legitimate re-chains).
func writeEvents(path string, evs []Event) error {
	f, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)
	if err != nil {
		return err
	}
	defer f.Close()
	w := bufio.NewWriter(f)
	prev := ""
	for _, e := range evs {
		e.Prev = prev
		line, err := json.Marshal(e)
		if err != nil {
			return err
		}
		if _, err := w.Write(append(line, '\n')); err != nil {
			return err
		}
		prev = chainHash(line)
	}
	if err := w.Flush(); err != nil {
		return err
	}
	return f.Sync()
}

// ChainReport is VerifyChain's result: chain breaks and unparseable lines in
// a log file (doctor surfaces both; store.scan tolerates them silently).
type ChainReport struct {
	Breaks      []string // "line N: …" descriptions
	Unparseable int
	Lines       int
}

// VerifyChain checks the live log's line hash chain. Legacy lines (written
// before chaining) carry no prev and are tolerated — but once any line
// carries prev, every later line must, or the gap is a break (an appended
// forgery cannot simply omit the field).
func (s *Store) VerifyChain() (ChainReport, error) {
	return verifyChainFile(s.EventsPath())
}

func verifyChainFile(path string) (ChainReport, error) {
	var rep ChainReport
	f, err := os.Open(path)
	if err != nil {
		if os.IsNotExist(err) {
			return rep, nil
		}
		return rep, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), maxLineTail)
	prevHash := ""
	chained := false // seen any prev-carrying line yet
	n := 0
	for sc.Scan() {
		raw := sc.Bytes()
		if len(raw) == 0 {
			continue
		}
		n++
		var ev Event
		if err := json.Unmarshal(raw, &ev); err != nil {
			rep.Unparseable++         // reported separately by doctor
			prevHash = chainHash(raw) // later lines chain over raw bytes regardless
			continue
		}
		switch {
		case ev.Prev == "" && n == 1:
			// genesis
		case ev.Prev == "" && !chained:
			// legacy segment (pre-chaining engine)
		case ev.Prev == "" && chained:
			rep.Breaks = append(rep.Breaks, fmt.Sprintf("line %d (%s %s): chain regression — no prev after chained lines", n, ev.Kind, ev.ID))
		case ev.Prev != prevHash:
			rep.Breaks = append(rep.Breaks, fmt.Sprintf("line %d (%s %s): prev %s does not match the preceding line (%s) — the log was edited out-of-band", n, ev.Kind, ev.ID, ev.Prev, prevHash))
		}
		if ev.Prev != "" {
			chained = true
		}
		prevHash = chainHash(raw)
	}
	rep.Lines = n
	return rep, sc.Err()
}

func countLines(path string) (int, error) {
	f, err := os.Open(path)
	if err != nil {
		return 0, err
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 4*1024*1024)
	n := 0
	for sc.Scan() {
		if len(sc.Bytes()) > 0 {
			n++
		}
	}
	return n, sc.Err()
}

// DeriveRun rebuilds a Run snapshot from the event log (doctor's repair path
// and the merge-conflict recovery path).
func (s *Store) DeriveRun() (*Run, error) {
	evs, err := s.Events(nil)
	if err != nil {
		return nil, err
	}
	var r *Run
	for _, e := range evs {
		switch e.Kind {
		case "run":
			switch e.Str("action") {
			case "start", "branch", "adopt":
				r = &Run{
					Schema: StateSchema, ID: e.Run, Family: e.Str("family"),
					Intent: e.Str("intent"), Status: "active", Started: e.TS,
					Parent: e.Str("parent"), SlipByAC: map[string]int{},
				}
			case "close":
				if r != nil && r.ID == e.Run {
					r = nil
				}
			}
		case "phase":
			if r == nil || r.ID != e.Run {
				continue
			}
			switch e.Str("action") {
			case "enter":
				r.Phase = e.Str("target")
			case "exit":
				r.ExitedPh = append(r.ExitedPh, e.Phase)
			case "waive":
				r.WaivedPh = append(r.WaivedPh, e.Str("target"))
			case "loop":
				r.Loops++
				r.Phase = e.Str("target")
			case "force":
				r.Forces++
				r.ExitedPh = append(r.ExitedPh, e.Phase)
			case "park":
				r.Status = "parked"
			case "resume":
				r.Status = "active"
			}
		}
	}
	return r, nil
}
