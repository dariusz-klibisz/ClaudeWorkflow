package store

import (
	"encoding/json"
	"os"
	"strings"
	"testing"
)

func appendN(t *testing.T, s *Store, n int, run string) {
	t.Helper()
	for i := 0; i < n; i++ {
		if err := s.Append(&Event{Kind: "risk", Run: run, Actor: "engine",
			Data: map[string]any{"i": i, "signals": []any{}, "lenses": []any{}}}); err != nil {
			t.Fatal(err)
		}
	}
}

func TestChainAppendAndVerify(t *testing.T) {
	s := newStore(t)
	appendN(t, s, 5, "r1")
	rep, err := s.VerifyChain()
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Breaks) != 0 || rep.Unparseable != 0 {
		t.Fatalf("fresh log must verify clean: %+v", rep)
	}
	if rep.Lines != 5 {
		t.Fatalf("want 5 lines, got %d", rep.Lines)
	}
	// every line after genesis must carry prev
	evs, _ := s.Events(nil)
	if evs[0].Prev != "" {
		t.Error("genesis event must have no prev")
	}
	for _, e := range evs[1:] {
		if e.Prev == "" {
			t.Fatalf("event %s missing prev", e.ID)
		}
	}
}

func TestChainDetectsInPlaceEdit(t *testing.T) {
	s := newStore(t)
	appendN(t, s, 4, "r1")
	raw, _ := os.ReadFile(s.EventsPath())
	lines := strings.Split(strings.TrimSuffix(string(raw), "\n"), "\n")
	// forge line 2: flip its payload but keep it valid JSON
	lines[1] = strings.Replace(lines[1], `"i":1`, `"i":99`, 1)
	_ = os.WriteFile(s.EventsPath(), []byte(strings.Join(lines, "\n")+"\n"), 0o644)
	rep, err := s.VerifyChain()
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Breaks) == 0 {
		t.Fatal("an edited line must break the chain")
	}
	if !strings.Contains(rep.Breaks[0], "edited out-of-band") {
		t.Errorf("break should name the cause: %s", rep.Breaks[0])
	}
}

func TestChainDetectsAppendedForgery(t *testing.T) {
	s := newStore(t)
	appendN(t, s, 3, "r1")
	// forged event appended without prev (the lazy forgery)
	forged, _ := json.Marshal(Event{Schema: 1, ID: NewULID(), Seq: 99, Kind: "test-run",
		Run: "r1", Auto: true, Actor: "hook", Data: map[string]any{"grounded": true, "exit": 0}})
	f, _ := os.OpenFile(s.EventsPath(), os.O_APPEND|os.O_WRONLY, 0o644)
	f.Write(append(forged, '\n'))
	f.Close()
	rep, _ := s.VerifyChain()
	found := false
	for _, b := range rep.Breaks {
		if strings.Contains(b, "chain regression") {
			found = true
		}
	}
	if !found {
		t.Fatalf("appended prev-less forgery must be a chain regression: %+v", rep.Breaks)
	}
}

func TestChainToleratesLegacyLog(t *testing.T) {
	s := newStore(t)
	// legacy lines: hand-written, no prev
	f, _ := os.OpenFile(s.EventsPath(), os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	for i := 0; i < 3; i++ {
		raw, _ := json.Marshal(Event{Schema: 1, ID: NewULID(), Seq: int64(i + 1), Kind: "risk", Run: "r0"})
		f.Write(append(raw, '\n'))
	}
	f.Close()
	// modern appends chain from there
	appendN(t, s, 2, "r0")
	rep, err := s.VerifyChain()
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Breaks) != 0 {
		t.Fatalf("legacy prefix must be tolerated: %+v", rep.Breaks)
	}
}

func TestChainRecomputedByArchive(t *testing.T) {
	s := newStore(t)
	_ = s.Append(&Event{Kind: "run", Run: "r1", Actor: "engine", Data: map[string]any{"action": "start", "family": "diff"}})
	appendN(t, s, 3, "r1")
	// a second run's events interleaved after (these stay live)
	_ = s.Append(&Event{Kind: "run", Run: "r2", Actor: "engine", Data: map[string]any{"action": "start", "family": "diff"}})
	_ = s.Append(&Event{Kind: "run", Run: "r1", Actor: "engine", Data: map[string]any{"action": "close"}})
	if err := s.ArchiveRun("r1"); err != nil {
		t.Fatal(err)
	}
	// live log re-chained
	rep, err := s.VerifyChain()
	if err != nil {
		t.Fatal(err)
	}
	if len(rep.Breaks) != 0 {
		t.Fatalf("compacted live log must verify clean: %+v", rep.Breaks)
	}
	// archived slice has its own clean chain
	arep, err := verifyChainFile(s.RunsDir() + "/r1/events.jsonl")
	if err != nil {
		t.Fatal(err)
	}
	if arep.Lines == 0 || len(arep.Breaks) != 0 {
		t.Fatalf("archived log must verify clean: %+v", arep)
	}
}

func TestChainSurvivesTornTailHealing(t *testing.T) {
	s := newStore(t)
	appendN(t, s, 2, "r1")
	f, _ := os.OpenFile(s.EventsPath(), os.O_APPEND|os.O_WRONLY, 0o644)
	f.WriteString(`{"schema":1,"kind":"torn`) // crash mid-write, no newline
	f.Close()
	appendN(t, s, 1, "r1") // heals the tail, chains over the torn raw line
	rep, err := s.VerifyChain()
	if err != nil {
		t.Fatal(err)
	}
	if rep.Unparseable != 1 {
		t.Fatalf("torn line must be counted unparseable: %+v", rep)
	}
	if len(rep.Breaks) != 0 {
		t.Fatalf("healed append must still chain correctly: %+v", rep.Breaks)
	}
}
