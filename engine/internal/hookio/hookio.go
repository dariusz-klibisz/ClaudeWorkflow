// Package hookio implements the Claude Code hook I/O contract (verified in
// workflow-redesign/01): JSON input on stdin, exit-code protocol (0 = allow /
// JSON processed, 2 = blocking error via stderr), event-specific JSON output
// shapes, and the 10,000-character output cap. It also detects hook-vs-agent
// invocation context (the WF_ENFORCE provenance guard, 04 §7).
package hookio

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"time"
)

// Input is the union of hook input fields across events; absent fields are
// zero. Unknown fields are ignored (forward compatibility with Claude Code
// releases); Raw preserves everything.
type Input struct {
	SessionID      string `json:"session_id"`
	TranscriptPath string `json:"transcript_path"`
	CWD            string `json:"cwd"`
	PermissionMode string `json:"permission_mode"`
	HookEventName  string `json:"hook_event_name"`

	// Session / prompt
	Source string `json:"source"` // SessionStart: startup|resume|clear|compact
	Prompt string `json:"prompt"` // UserPromptSubmit

	// Tool events
	ToolName     string          `json:"tool_name"`
	ToolInput    json.RawMessage `json:"tool_input"`
	ToolResponse json.RawMessage `json:"tool_response"`
	// PostToolUseFailure only: failed tool calls never reach PostToolUse —
	// a non-zero Bash exit arrives HERE, with the code embedded in the
	// human-readable error string ("Command exited with non-zero status
	// code 1") and no tool_response at all.
	Error       string `json:"error"`
	IsInterrupt bool   `json:"is_interrupt"`

	// Stop family
	StopHookActive       bool            `json:"stop_hook_active"`
	LastAssistantMessage string          `json:"last_assistant_message"`
	BackgroundTasks      json.RawMessage `json:"background_tasks"`

	// Subagent events
	AgentID             string `json:"agent_id"`
	AgentType           string `json:"agent_type"`
	AgentTranscriptPath string `json:"agent_transcript_path"`

	// Task events
	TaskID          string `json:"task_id"`
	TaskSubject     string `json:"task_subject"`
	TaskDescription string `json:"task_description"`

	Raw map[string]any `json:"-"`
}

// Read parses hook input from r. A nil/empty stream yields an empty Input
// (agent-invoked context).
func Read(r io.Reader) (*Input, error) {
	raw, err := io.ReadAll(r)
	if err != nil {
		return nil, err
	}
	in := &Input{}
	if len(raw) == 0 {
		return in, nil
	}
	if err := json.Unmarshal(raw, in); err != nil {
		return nil, fmt.Errorf("hook input malformed: %w", err)
	}
	_ = json.Unmarshal(raw, &in.Raw)
	return in, nil
}

// FromHook reports whether this invocation carries hook input — the
// provenance test: WF_ENFORCE and other env escapes are honored only when
// true (the env belongs to Claude Code, i.e. the user), never when wf is
// called from the agent's own Bash tool.
func (in *Input) FromHook() bool {
	return in != nil && in.HookEventName != ""
}

// ToolInputField extracts a string field from tool_input.
func (in *Input) ToolInputField(key string) string {
	if len(in.ToolInput) == 0 {
		return ""
	}
	var m map[string]any
	if json.Unmarshal(in.ToolInput, &m) != nil {
		return ""
	}
	v, _ := m[key].(string)
	return v
}

// ---------------------------------------------------------------------------
// Output
// ---------------------------------------------------------------------------

// MaxOutput is the documented cap on hook output strings.
const MaxOutput = 10000

// Cap truncates s to stay under the platform's 10k-char cap with margin.
func Cap(s string) string {
	const limit = MaxOutput - 500
	if len(s) <= limit {
		return s
	}
	cut := s[:limit]
	if i := strings.LastIndexByte(cut, '\n'); i > limit/2 {
		cut = cut[:i]
	}
	return cut + "\n[wf: output truncated]"
}

// Result is what a gate returns; the CLI layer emits it and exits.
type Result struct {
	Exit   int
	Stdout string // JSON or plain text (processed on exit 0)
	Stderr string // fed to Claude on exit 2
}

func Allow() Result { return Result{} }

// AllowJSON emits a JSON decision with exit 0.
func AllowJSON(v any) Result {
	raw, err := json.Marshal(v)
	if err != nil {
		return Result{Exit: 1, Stderr: err.Error()}
	}
	return Result{Stdout: string(raw)}
}

// Block is the exit-2 blocking error: stderr is fed to Claude.
func Block(reason string) Result {
	return Result{Exit: 2, Stderr: Cap(reason)}
}

// BrokenGate is the fail-open-loud path for *sequencing* gates when the
// engine itself is unhealthy: allow, but tell the user (04 §7). Identical
// messages are rate-limited (once per hour) — degraded gates fire on every
// tool call, and repeating the same wall of text erodes the signal.
func BrokenGate(err error) Result {
	msg := fmt.Sprintf("[wf] gate degraded (fail-open): %v — run `wf doctor`", err)
	if brokenGateSeenRecently(msg) {
		return Allow()
	}
	return AllowJSON(map[string]any{"systemMessage": Cap(msg)})
}

func brokenGateSeenRecently(msg string) bool {
	h := uint64(14695981039346656037)
	for i := 0; i < len(msg); i++ {
		h ^= uint64(msg[i])
		h *= 1099511628211
	}
	marker := filepath.Join(os.TempDir(), fmt.Sprintf("wf-degraded-%x", h))
	if st, err := os.Stat(marker); err == nil && time.Since(st.ModTime()) < time.Hour {
		return true
	}
	_ = os.WriteFile(marker, nil, 0o644)
	return false
}

// StopBlock prevents Claude from stopping (Stop/SubagentStop): exit 0 JSON
// decision with the reason delivered as the next instruction.
func StopBlock(reason string) Result {
	return AllowJSON(map[string]any{"decision": "block", "reason": Cap(reason)})
}

// StopAllowMessage allows the stop but surfaces a status line to the user.
func StopAllowMessage(msg string) Result {
	return AllowJSON(map[string]any{"systemMessage": Cap(msg)})
}

// Deny blocks a PreToolUse tool call with a reason shown to Claude.
func Deny(reason string) Result {
	return AllowJSON(map[string]any{
		"hookSpecificOutput": map[string]any{
			"hookEventName":            "PreToolUse",
			"permissionDecision":       "deny",
			"permissionDecisionReason": Cap(reason),
		},
	})
}

// AdditionalContext injects context at the firing event (SessionStart,
// UserPromptSubmit, SubagentStart, …).
func AdditionalContext(event, text string) Result {
	return AllowJSON(map[string]any{
		"hookSpecificOutput": map[string]any{
			"hookEventName":     event,
			"additionalContext": Cap(text),
		},
	})
}

// Emit writes the result and returns the process exit code.
func (r Result) Emit(stdout, stderr io.Writer) int {
	if r.Stdout != "" {
		fmt.Fprintln(stdout, r.Stdout)
	}
	if r.Stderr != "" {
		fmt.Fprintln(stderr, r.Stderr)
	}
	return r.Exit
}

// EnforceDisabled reports whether the user disabled sequencing gates via
// WF_ENFORCE=0 — honored only in hook context (provenance guard).
func EnforceDisabled(in *Input) bool {
	return in.FromHook() && os.Getenv("WF_ENFORCE") == "0"
}
