package runctl

// Challenge approvals (`approvals: "challenge"`, 04 §8.1 escalated): the
// engine mints a single-use code into per-machine local state and shows it
// ONLY through the user's statusline — never on engine stdout, so the model
// cannot see it before the user types it into the AskUserQuestion exchange.
// Honest bounds, stated: the code is on disk under .workflow/local/ (tool
// gates deny model reads of that path, MCP/externals are out of scope), and
// after the user answers, the model sees the consumed code — which by then
// anchors nothing new (single-use).

import (
	"crypto/rand"
	"fmt"
	"os"
	"time"
)

// Challenge is the pending statusline-delivered approval code.
type Challenge struct {
	Gate    string `json:"gate"`
	Code    string `json:"code"`
	Created string `json:"created"`
}

const challengeFile = "challenge.json"

// ensureChallenge returns the pending code for the gate, minting one when
// none exists (or the pending one belongs to a different gate). The code is
// never part of the returned error path — callers must not print it.
func (c *Ctl) ensureChallenge(gate string) (string, error) {
	var ch Challenge
	if err := c.Store.LoadLocal(challengeFile, &ch); err == nil && ch.Gate == gate && ch.Code != "" {
		return ch.Code, nil
	}
	code, err := mintCode()
	if err != nil {
		return "", fmt.Errorf("challenge code generation failed: %w", err)
	}
	ch = Challenge{Gate: gate, Code: code, Created: time.Now().UTC().Format(time.RFC3339)}
	if err := c.Store.SaveLocal(challengeFile, &ch); err != nil {
		return "", err
	}
	return code, nil
}

// clearChallenge consumes the pending code (single-use).
func (c *Ctl) clearChallenge() {
	_ = os.Remove(c.Store.LocalPath(challengeFile))
}

// PendingChallenge is the statusline's read path. nil = nothing pending.
func (c *Ctl) PendingChallenge() *Challenge {
	var ch Challenge
	if err := c.Store.LoadLocal(challengeFile, &ch); err != nil || ch.Code == "" {
		return nil
	}
	return &ch
}

// mintCode returns 6 chars from an unambiguous uppercase alphabet
// (no 0/O/1/I) — typed by a human from a statusline, so legibility counts.
func mintCode() (string, error) {
	const alphabet = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	raw := make([]byte, 6)
	if _, err := rand.Read(raw); err != nil {
		return "", err
	}
	out := make([]byte, len(raw))
	for i, b := range raw {
		out[i] = alphabet[int(b)%len(alphabet)]
	}
	return string(out), nil
}
