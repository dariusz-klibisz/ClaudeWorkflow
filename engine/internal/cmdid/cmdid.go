// Package cmdid reasons about shell-command identity for test-evidence
// capture and red→green pairing. Shared by the capture hooks (gates) and
// the contract evaluator (contracts) — the latter must not import gates.
package cmdid

import (
	"regexp"
	"strings"
)

// Head strips leading env assignments and returns the command's start.
func Head(cmd string) string {
	fields := strings.Fields(cmd)
	i := 0
	for i < len(fields) && strings.Contains(fields[i], "=") && !strings.HasPrefix(fields[i], "=") {
		i++
	}
	return strings.Join(fields[i:], " ")
}

// leadingCd matches exactly one `cd <dir> &&` prefix (quoted dirs included).
var leadingCd = regexp.MustCompile(`^cd\s+("[^"]*"|'[^']*'|[^\s&|;]+)\s*&&\s*(.+)$`)

// Effective strips a single leading `cd <dir> &&` — the one chained form
// common enough to honor for recognition and pairing.
func Effective(cmd string) string {
	if m := leadingCd.FindStringSubmatch(strings.TrimSpace(cmd)); m != nil {
		return strings.TrimSpace(m[2])
	}
	return strings.TrimSpace(cmd)
}

// GenericInterpreters are launchers whose bare name proves nothing about
// testing — a strategy of `python3 test_app.py` must learn that exact
// two-token invocation, never bare `python3` (or every script run would
// count as red/green evidence).
var GenericInterpreters = map[string]bool{
	"python": true, "python3": true, "python2": true, "py": true,
	"node": true, "deno": true, "bun": true, "ruby": true, "perl": true,
	"php": true, "sh": true, "bash": true, "zsh": true, "dash": true,
	"pwsh": true, "powershell": true, "java": true, "dotnet": true,
	"go": true, "cargo": true, "npx": true, "uv": true, "uvx": true,
	"make": true, "npm": true, "pnpm": true, "yarn": true,
}

// TokenPrefix: the shorter command is a whole-token prefix of the longer,
// sharing at least 2 tokens — so `python3 test_app.py` matches strategy
// `python3 test_app.py -v`, but a bare interpreter never matches anything.
func TokenPrefix(a, b []string) bool {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	if n < 2 {
		return false
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

// LearnedHead extracts the runner-invocation head from a verification
// command: tokens up to the first selector (path / dotted test id / pytest
// `::`) or flag, keeping interpreter module flags (`python -m unittest`).
// A single-token head that is a generic interpreter is discarded.
func LearnedHead(cmd string) string {
	fields := strings.Fields(Head(cmd))
	var head []string
	for i, tok := range fields {
		if strings.ContainsAny(tok, "/\\.") || strings.Contains(tok, "::") {
			break
		}
		if strings.HasPrefix(tok, "-") {
			if tok == "-m" && i == 1 {
				head = append(head, tok)
				continue
			}
			break
		}
		head = append(head, tok)
		if len(head) == 4 {
			break
		}
	}
	if len(head) == 0 || (len(head) == 1 && GenericInterpreters[head[0]]) {
		return ""
	}
	// a dangling module flag proves nothing either ("python3 -m")
	if head[len(head)-1] == "-m" {
		return ""
	}
	return strings.Join(head, " ")
}

// Pair classifies how strongly a red and a green test invocation talk about
// the same code.
type Pair int

const (
	// PairNone: different runners — the green proves nothing about the red.
	PairNone Pair = iota
	// PairWeak: same runner, diverging selectors — plausible but not
	// certain the green exercises what the red exercised.
	PairWeak
	// PairStrict: same command, or one is a whole-token prefix of the
	// other (red on a selector, green on the full suite — or vice versa).
	PairStrict
)

// Classify judges a red→green command pair. Both commands are normalized
// (env assignments, one leading `cd &&`) before comparison. Both sides were
// already recognized as test runs by capture, so a shared launcher token is
// enough for a weak pair — the runner-mismatch class this exists to kill is
// cross-tool (gitleaks red "pairing" a pytest green).
func Classify(redCmd, greenCmd string) Pair {
	r := strings.Fields(Head(Effective(redCmd)))
	g := strings.Fields(Head(Effective(greenCmd)))
	if len(r) == 0 || len(g) == 0 {
		return PairNone
	}
	if prefixEither(r, g) {
		return PairStrict
	}
	if r[0] == g[0] {
		return PairWeak
	}
	return PairNone
}

// prefixEither: one token list is a whole-token prefix of the other
// (≥1 shared token — both sides were already recognized as test runs).
func prefixEither(a, b []string) bool {
	n := len(a)
	if len(b) < n {
		n = len(b)
	}
	if n < 1 {
		return false
	}
	for i := 0; i < n; i++ {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}
