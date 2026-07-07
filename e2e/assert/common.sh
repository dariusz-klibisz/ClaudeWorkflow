# Shared assertion helpers — sourced by scenario asserts.
# Env: FIXTURE, WF, REPO, SESSION_EXIT.
# Philosophy (07 §6 layer 3): assert FROM STATE that the gate HELD —
# invariants over the ledger, not expectations about model choices.

EV="$FIXTURE/.workflow/log/events.jsonl"

fail() { echo "ASSERT FAIL: $*" >&2; exit 1; }
note() { echo "  assert: $*"; }

# all_events: live log + any archived runs, as a JSON array on stdout.
all_events() {
  {
    [ -f "$EV" ] && cat "$EV"
    for f in "$FIXTURE"/.workflow/runs/*/events.jsonl; do
      [ -f "$f" ] && cat "$f"
    done
  } 2>/dev/null | jq -s '.'
}

# count '<jq predicate over event .>'
count() {
  all_events | jq "[.[] | select($1)] | length"
}

# ledger_coherent: every line of every ledger file parses (no torn writes).
ledger_coherent() {
  for f in "$EV" "$FIXTURE"/.workflow/runs/*/events.jsonl; do
    [ -f "$f" ] || continue
    jq -c . "$f" >/dev/null 2>&1 || fail "torn/corrupt ledger: $f"
  done
  note "ledger coherent"
}

# session_not_wedged: the headless session ended by itself (any exit code
# except timeout's 124 — a wedged Stop gate would spin until the timeout).
session_not_wedged() {
  [ "$SESSION_EXIT" = 124 ] && fail "session hit the wall-clock timeout (wedged?)"
  note "session ended on its own (exit $SESSION_EXIT)"
}
