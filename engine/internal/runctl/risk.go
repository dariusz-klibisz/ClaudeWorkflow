package runctl

import (
	"sort"
	"strings"
)

// Deterministic risk keyword screen (the v0.36 risk_scan.sh, engine-native).
// Signals bind lenses; the Frame contract then demands ≥1 ambiguity per lens.
var signalKeywords = map[string][]string{
	"auth":        {"auth", "login", "password", "token", "credential", "session", "permission", "oauth", "jwt"},
	"network":     {"http", "api", "network", "request", "socket", "endpoint", "url", "webhook"},
	"data":        {"database", "db", "migration", "schema", "persist", "storage", "file", "write", "delete"},
	"boundary":    {"parse", "input", "upload", "deserialize", "external", "untrusted", "user-provided", "import"},
	"destructive": {"drop", "truncate", "rm -", "force-push", "overwrite", "wipe", "purge"},
	"concurrency": {"concurrent", "parallel", "race", "lock", "thread", "async", "goroutine"},
	"ui":          {"ui", "screen", "form", "button", "page", "accessibility", "wcag", "frontend"},
}

var signalLenses = map[string][]string{
	"auth":        {"security"},
	"network":     {"security", "reliability"},
	"data":        {"security", "reliability"},
	"boundary":    {"security", "adversarial"},
	"destructive": {"adversarial", "operator"},
	"concurrency": {"reliability"},
	"ui":          {"usability"},
}

// RiskScan runs the keyword screen over text, merges agent-added signals, and
// records the risk event (signals[] + lens bindings). The "user" lens is
// always selected.
func (c *Ctl) RiskScan(text string, added []string) (signals, lenses []string, err error) {
	low := strings.ToLower(text)
	set := map[string]bool{}
	for sig, words := range signalKeywords {
		for _, w := range words {
			if strings.Contains(low, w) {
				set[sig] = true
				break
			}
		}
	}
	for _, a := range added {
		a = strings.TrimSpace(a)
		if a != "" {
			set[a] = true
		}
	}
	lensSet := map[string]bool{"user": true}
	for sig := range set {
		for _, l := range signalLenses[sig] {
			lensSet[l] = true
		}
	}
	for s := range set {
		signals = append(signals, s)
	}
	for l := range lensSet {
		lenses = append(lenses, l)
	}
	sort.Strings(signals)
	sort.Strings(lenses)
	sigAny := make([]any, len(signals))
	for i, s := range signals {
		sigAny[i] = s
	}
	lensAny := make([]any, len(lenses))
	for i, l := range lenses {
		lensAny[i] = l
	}
	_, err = c.Record("risk", map[string]any{"signals": sigAny, "lenses": lensAny, "text_scanned": len(text) > 0}, true, "engine")
	return signals, lenses, err
}
