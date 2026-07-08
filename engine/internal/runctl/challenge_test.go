package runctl

import (
	"strings"
	"testing"
)

// Challenge approvals: the code is minted on the first attempt (statusline-
// delivered, never printed), and only an answer containing it approves.
func TestApproveChallengeFlow(t *testing.T) {
	c := newCtl(t)
	c.Config.Flags = map[string]any{"approvals": "challenge"}
	_, _ = c.RunStart("diff", "fix")

	// first attempt: refused, challenge minted — and the error must not
	// leak the code
	_, err := c.Approve("frame", "p")
	if err == nil {
		t.Fatal("challenge mode must refuse without the code")
	}
	ch := c.PendingChallenge()
	if ch == nil || ch.Gate != "frame" || len(ch.Code) != 6 {
		t.Fatalf("challenge not minted: %+v", ch)
	}
	if strings.Contains(err.Error(), ch.Code) {
		t.Fatal("the refusal message must never carry the code")
	}

	// a captured answer WITHOUT the code still refuses (hardened-plus)
	_, _ = c.Record("user-answer", map[string]any{"question": "approve frame?", "answer": "yes go ahead"}, true, "hook")
	if _, err := c.Approve("frame", "p"); err == nil {
		t.Fatal("an answer without the code must not approve")
	}
	// the pending code survives failed attempts (statusline stays stable)
	if ch2 := c.PendingChallenge(); ch2 == nil || ch2.Code != ch.Code {
		t.Fatal("failed attempts must not rotate the code")
	}

	// the user types the code (case-insensitively) → approved + consumed
	_, _ = c.Record("user-answer", map[string]any{"question": "approve frame?", "answer": "ok: " + strings.ToLower(ch.Code)}, true, "hook")
	ev, err := c.Approve("frame", "p")
	if err != nil {
		t.Fatalf("code-bearing answer must approve: %v", err)
	}
	if chal, _ := ev.Data["challenge"].(bool); !chal {
		t.Fatal("approval must be marked challenge-verified")
	}
	if ref, _ := ev.Data["answer_ref"].(string); ref == "" {
		t.Fatal("challenge approval must carry the anchor")
	}
	if c.PendingChallenge() != nil {
		t.Fatal("the code is single-use — consumed on success")
	}

	// consumed code cannot approve again: next approval mints a NEW code
	_, err = c.Approve("frame", "p2")
	if err == nil {
		t.Fatal("a consumed code must not approve a second time")
	}
	if ch3 := c.PendingChallenge(); ch3 == nil || ch3.Code == ch.Code {
		t.Fatal("re-approval must mint a fresh code")
	}
}

// An answer about a DIFFERENT gate never verifies the challenge, even when
// it contains the code (the anchoring-race fix applies here too).
func TestApproveChallengeTopicMismatch(t *testing.T) {
	c := newCtl(t)
	c.Config.Flags = map[string]any{"approvals": "challenge"}
	_, _ = c.RunStart("diff", "fix")

	_, _ = c.Approve("plan", "p") // mints the plan code
	ch := c.PendingChallenge()
	if ch == nil {
		t.Fatal("challenge not minted")
	}
	// design-topic answer carrying the plan code: must not verify plan
	_, _ = c.Record("user-answer", map[string]any{"question": "approve the design?", "answer": ch.Code, "topic": "design"}, true, "hook")
	if _, err := c.Approve("plan", "p"); err == nil {
		t.Fatal("an answer about another gate must not verify the challenge")
	}
}

// Switching gates re-mints: a pending frame code cannot leak into a scope
// approval.
func TestApproveChallengeGateSwitch(t *testing.T) {
	c := newCtl(t)
	c.Config.Flags = map[string]any{"approvals": "challenge"}
	_, _ = c.RunStart("diff", "fix")

	_, _ = c.Approve("frame", "p")
	frameCode := c.PendingChallenge().Code
	_, _ = c.Approve("scope", "p")
	ch := c.PendingChallenge()
	if ch.Gate != "scope" || ch.Code == frameCode {
		t.Fatalf("gate switch must mint a fresh code: %+v", ch)
	}
}
