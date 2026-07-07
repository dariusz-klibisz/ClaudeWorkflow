package views

// wf statusline — the one-line statusLine payload (09 Q8): run, phase
// position, unmet count, who is blocked, dead-hooks marker. Cheap by
// construction (snapshot + one phase evaluation, the inject read path) and
// never loud: statuslines re-run constantly, so errors degrade to "".

import (
	"fmt"
	"strings"

	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/contracts"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/doctor"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/runctl"
)

// Statusline renders the payload. "" = print nothing (broken state is the
// doctor's job, not the statusline's).
func Statusline(c *runctl.Ctl) string {
	r, err := c.Store.LoadRun()
	if err != nil {
		return ""
	}
	if r == nil {
		return "wf: no run"
	}
	short := r.ID
	if i := strings.LastIndex(short, "-"); i >= 0 {
		short = short[i+1:]
	}
	if r.Status == "parked" {
		return fmt.Sprintf("wf %s · parked — /wf:dev to resume", short)
	}
	if r.Phase == "" {
		return fmt.Sprintf("wf %s · all phases done — wf run close", short)
	}

	line := fmt.Sprintf("wf %s · %s", short, r.Phase)
	phases := c.Spec.PhasesFor(r.Family)
	for i, p := range phases {
		if p.ID == r.Phase {
			line += fmt.Sprintf(" (%d/%d)", i+1, len(phases))
			break
		}
	}

	env, err := c.Env(r)
	if err != nil {
		return line
	}
	findings, err := contracts.EvaluatePhase(env, r.Phase)
	if err != nil {
		return line
	}
	if len(findings) == 0 {
		line += " · contract met"
	} else {
		userBlocked := 0
		for _, f := range findings {
			if f.UserBlocked {
				userBlocked++
			}
		}
		line += fmt.Sprintf(" · %d unmet", len(findings))
		if userBlocked == len(findings) {
			line += " · waiting: user"
		}
	}
	if doctor.HookLiveness(c, r) != "" {
		line += " · ⚠ hooks dead"
	}
	return line
}
