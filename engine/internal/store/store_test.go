package store

import (
	"os"
	"path/filepath"
	"sort"
	"testing"
)

func newStore(t *testing.T) *Store {
	t.Helper()
	s, err := Open(t.TempDir(), true)
	if err != nil {
		t.Fatal(err)
	}
	return s
}

func TestULIDMonotonicAndUnique(t *testing.T) {
	seen := map[string]bool{}
	var ids []string
	for i := 0; i < 5000; i++ {
		id := NewULID()
		if len(id) != 26 {
			t.Fatalf("bad length %d: %q", len(id), id)
		}
		if seen[id] {
			t.Fatalf("duplicate ULID %q", id)
		}
		seen[id] = true
		ids = append(ids, id)
	}
	if !sort.StringsAreSorted(ids) {
		t.Error("same-process ULIDs must be monotonic")
	}
}

func TestAppendAndScan(t *testing.T) {
	s := newStore(t)
	for i, kind := range []string{"run", "classification", "risk"} {
		ev := &Event{Kind: kind, Run: "r1", Actor: "engine", Data: map[string]any{"i": i}}
		if err := s.Append(ev); err != nil {
			t.Fatal(err)
		}
		if ev.ID == "" || ev.Seq == 0 || ev.TS == "" {
			t.Fatalf("identity not stamped: %+v", ev)
		}
	}
	evs, err := s.RunEvents("r1")
	if err != nil {
		t.Fatal(err)
	}
	if len(evs) != 3 {
		t.Fatalf("want 3 events, got %d", len(evs))
	}
	if evs[0].Seq >= evs[1].Seq || evs[1].Seq >= evs[2].Seq {
		t.Error("seq not increasing")
	}
}

func TestSeqResumesAfterReopen(t *testing.T) {
	dir := t.TempDir()
	s, _ := Open(dir, true)
	_ = s.Append(&Event{Kind: "run", Run: "r1"})
	_ = s.Append(&Event{Kind: "risk", Run: "r1"})
	s2, _ := Open(dir, false)
	ev := &Event{Kind: "task", Run: "r1"}
	_ = s2.Append(ev)
	if ev.Seq != 3 {
		t.Errorf("seq must resume from log: want 3, got %d", ev.Seq)
	}
}

func TestScanToleratesTornLine(t *testing.T) {
	s := newStore(t)
	_ = s.Append(&Event{Kind: "run", Run: "r1"})
	f, _ := os.OpenFile(s.EventsPath(), os.O_APPEND|os.O_WRONLY, 0o644)
	f.WriteString(`{"schema":1,"kind":"torn`) // crash mid-write
	f.Close()
	s2, _ := Open(filepath.Dir(s.Root), false)
	_ = s2.Append(&Event{Kind: "risk", Run: "r1"})
	evs, err := s2.RunEvents("r1")
	if err != nil {
		t.Fatal(err)
	}
	if len(evs) != 2 {
		t.Fatalf("torn line must be skipped, valid ones kept: got %d", len(evs))
	}
}

func TestRunSnapshotRoundTrip(t *testing.T) {
	s := newStore(t)
	if r, _ := s.LoadRun(); r != nil {
		t.Fatal("no run expected")
	}
	r := &Run{ID: "r1", Family: "diff", Intent: "fix", Phase: "frame", Status: "active"}
	if err := s.SaveRun(r); err != nil {
		t.Fatal(err)
	}
	got, err := s.LoadRun()
	if err != nil || got == nil {
		t.Fatalf("load: %v", err)
	}
	if got.ID != "r1" || got.Phase != "frame" || got.Family != "diff" {
		t.Errorf("round trip mismatch: %+v", got)
	}
	_ = s.ClearRun()
	if r, _ := s.LoadRun(); r != nil {
		t.Error("cleared run still present")
	}
}

func TestLockExclusive(t *testing.T) {
	s := newStore(t)
	if err := s.Lock(); err != nil {
		t.Fatal(err)
	}
	s2 := &Store{Root: s.Root}
	done := make(chan error, 1)
	go func() {
		err := s2.Lock()
		if err == nil {
			s2.Unlock()
		}
		done <- err
	}()
	s.Unlock()
	if err := <-done; err != nil {
		t.Fatalf("second locker should acquire after release: %v", err)
	}
}

func TestDeriveRunFromLog(t *testing.T) {
	s := newStore(t)
	app := func(kind string, run string, data map[string]any) {
		if err := s.Append(&Event{Kind: kind, Run: run, Actor: "engine", Data: data}); err != nil {
			t.Fatal(err)
		}
	}
	app("run", "r1", map[string]any{"action": "start", "family": "diff", "intent": "fix"})
	app("phase", "r1", map[string]any{"action": "enter", "target": "frame"})
	app("phase", "r1", map[string]any{"action": "exit"})
	app("phase", "r1", map[string]any{"action": "enter", "target": "context"})
	r, err := s.DeriveRun()
	if err != nil {
		t.Fatal(err)
	}
	if r == nil || r.ID != "r1" || r.Phase != "context" || r.Family != "diff" {
		t.Fatalf("derive mismatch: %+v", r)
	}
	app("run", "r1", map[string]any{"action": "close"})
	r, _ = s.DeriveRun()
	if r != nil {
		t.Error("closed run must derive to nil")
	}
}

func TestArchiveRunKeepsDurableEvents(t *testing.T) {
	s := newStore(t)
	app := func(ev *Event) {
		if err := s.Append(ev); err != nil {
			t.Fatal(err)
		}
	}
	app(&Event{Kind: "run", Run: "r1", Data: map[string]any{"action": "start"}})
	app(&Event{Kind: "risk", Run: "r1"})
	app(&Event{Kind: "followup", Run: "r1", Data: map[string]any{"status": "next-run", "text": "carry me"}})
	app(&Event{Kind: "commit-origin", Run: "r1", Data: map[string]any{"commit": "abc"}})
	app(&Event{Kind: "run", Run: "r2", Data: map[string]any{"action": "start"}})
	if err := s.ArchiveRun("r1"); err != nil {
		t.Fatal(err)
	}
	live, _ := s.Events(nil)
	kinds := map[string]int{}
	for _, e := range live {
		kinds[e.Kind]++
	}
	if kinds["risk"] != 0 {
		t.Error("archived event still live")
	}
	if kinds["followup"] != 1 {
		t.Errorf("open followups must stay live: %v", kinds)
	}
	// bounded live log: commit-origin archives with its run — readers get it
	// back via AllEvents (archives + live)
	if kinds["commit-origin"] != 0 {
		t.Errorf("commit-origin must archive with its run: %v", kinds)
	}
	if kinds["run"] != 1 {
		t.Errorf("other runs' events must stay: %v", kinds)
	}
	arch, err := os.ReadFile(filepath.Join(s.RunsDir(), "r1", "events.jsonl"))
	if err != nil {
		t.Fatal(err)
	}
	if len(arch) == 0 {
		t.Error("archive empty")
	}
	all, err := s.AllEvents()
	if err != nil {
		t.Fatal(err)
	}
	allKinds := map[string]int{}
	for _, e := range all {
		allKinds[e.Kind]++
	}
	if allKinds["commit-origin"] != 1 || allKinds["risk"] != 1 {
		t.Errorf("AllEvents must fold archived events back in: %v", allKinds)
	}
}

func TestPruneLocalOnArchive(t *testing.T) {
	s := newStore(t)
	_ = s.SaveLocal("tasks-mirror.json", map[string]any{"map": map[string]string{"1": "x"}})
	_ = s.SaveLocal("verdict-attempts.json", map[string]any{"attempts": map[string]int{"a": 3}})
	_ = s.Append(&Event{Kind: "run", Run: "r1", Data: map[string]any{"action": "start"}})
	_ = s.Append(&Event{Kind: "run", Run: "r1", Data: map[string]any{"action": "close"}})
	if err := s.ArchiveRun("r1"); err != nil {
		t.Fatal(err)
	}
	for _, n := range []string{"tasks-mirror.json", "verdict-attempts.json"} {
		if _, err := os.Stat(s.LocalPath(n)); !os.IsNotExist(err) {
			t.Errorf("local/%s must be pruned at close", n)
		}
	}
}
