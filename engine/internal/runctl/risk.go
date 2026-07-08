package runctl

import (
	"regexp"
	"sort"
	"strings"
	"sync"
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

var (
	kwOnce sync.Once
	kwRe   map[string]*regexp.Regexp
)

// keyword matching is WORD-BOUNDED: substring matching made "requirements"
// trigger the `ui` signal (live-testing bug). Multi-word keywords match as
// phrases; "rm -" style fragments keep their literal tail.
func keywordRegexps() map[string]*regexp.Regexp {
	kwOnce.Do(func() {
		isAlnum := func(b byte) bool {
			return b >= 'a' && b <= 'z' || b >= '0' && b <= '9'
		}
		kwRe = map[string]*regexp.Regexp{}
		for sig, words := range signalKeywords {
			parts := make([]string, len(words))
			for i, w := range words {
				p := regexp.QuoteMeta(w)
				if isAlnum(w[0]) {
					p = `(^|[^a-z0-9])` + p
				}
				if isAlnum(w[len(w)-1]) {
					p += `($|[^a-z0-9])`
				}
				parts[i] = p
			}
			kwRe[sig] = regexp.MustCompile(strings.Join(parts, "|"))
		}
	})
	return kwRe
}

// standingFlagSignals maps project config flags (.workflow/config.json
// `flags`) to always-on risk signals: a project that declares it handles PII
// or faces the internet gets those lenses on EVERY run, not only when the
// task text happens to mention a keyword.
var standingFlagSignals = map[string][]string{
	"pii":             {"data"},
	"internet_facing": {"network", "boundary"},
	"public_api":      {"network", "boundary"},
}

// RiskScan runs the keyword screen over text, merges agent-added signals and
// standing project flags, and records the risk event (signals[] + lens
// bindings). The "user" lens is always selected.
func (c *Ctl) RiskScan(text string, added []string) (signals, lenses []string, err error) {
	low := strings.ToLower(text)
	set := map[string]bool{}
	for sig, re := range keywordRegexps() {
		if re.MatchString(low) {
			set[sig] = true
		}
	}
	if c.Config != nil {
		for flag, sigs := range standingFlagSignals {
			if v, _ := c.Config.ConfigFlag(flag).(bool); v {
				for _, s := range sigs {
					set[s] = true
				}
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
