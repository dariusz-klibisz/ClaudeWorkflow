// wf — the workflow engine. One binary behind every hook, gate, record
// command, and injection (workflow-redesign/07).
//
// Exit codes: 0 ok/allow · 2 blocked/gaps · 3 engine or contract broken.
package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/doctor"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/gates"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/hookio"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/inject"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/ops"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/runctl"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/selftest"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/spec"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/store"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/views"
)

var Version = "0.1.0-dev"

func main() {
	os.Exit(run(os.Args[1:]))
}

func run(args []string) int {
	if len(args) == 0 {
		usage()
		return 0
	}
	cmd, rest := args[0], args[1:]
	switch cmd {
	case "version", "--version":
		fmt.Println("wf", Version)
		return 0
	case "selftest":
		p, err := resolveSpecPath()
		if err != nil {
			fmt.Fprintln(os.Stderr, "wf:", err)
			return 3
		}
		if selftest.Run(p) > 0 {
			return 2
		}
		return 0
	case "help", "--help", "-h":
		usage()
		return 0
	}

	projectDir := resolveProjectDir()

	// gate/inject/capture commands read hook stdin and never hard-fail a
	// sequencing path (04 §7).
	switch cmd {
	case "gate":
		return gateCmd(projectDir, rest)
	case "inject":
		return injectCmd(projectDir, rest)
	case "capture":
		return captureCmd(projectDir, rest)
	}

	// doctor --bootstrap verifies the install before any project is adopted
	if cmd == "doctor" && len(rest) > 0 && rest[0] == "--bootstrap" {
		sp, err := loadSpecOnly()
		if err != nil {
			fmt.Fprintln(os.Stderr, "wf:", err)
			return 3
		}
		fmt.Println("wf", Version, "— spec:", len(sp.Contracts), "contract items,", len(sp.Roster), "agents")
		return bootstrapHealth()
	}

	ctl, err := openCtl(projectDir, cmd == "init")
	if err != nil {
		fmt.Fprintln(os.Stderr, "wf:", err)
		return 3
	}
	defer ctl.Store.Unlock()

	switch cmd {
	case "init":
		return initCmd(ctl)
	case "run":
		return runCmd(ctl, rest)
	case "phase":
		return phaseCmd(ctl, rest)
	case "record":
		return recordCmd(ctl, rest)
	case "approve":
		return approveCmd(ctl, rest)
	case "contract":
		return contractCmd(ctl, rest)
	case "loop":
		return loopCmd(ctl, rest)
	case "park":
		return parkCmd(ctl, rest)
	case "risk":
		return riskCmd(ctl, rest)
	case "status":
		return statusCmd(ctl, rest)
	case "doctor":
		return doctorCmd(ctl, rest)
	case "trace":
		return traceCmd(ctl)
	case "report":
		return reportCmd(ctl, rest)
	case "lessons":
		return lessonsCmd(ctl, projectDir, rest)
	case "deps":
		return depsCmd(ctl, projectDir, rest)
	case "origin":
		return originCmd(ctl, projectDir, rest)
	case "doc":
		return docCmd(ctl, projectDir, rest)
	default:
		fmt.Fprintf(os.Stderr, "wf: unknown command %q (wf help)\n", cmd)
		return 3
	}
}

// ---------------------------------------------------------------------------
// Environment resolution
// ---------------------------------------------------------------------------

func resolveProjectDir() string {
	if d := os.Getenv("CLAUDE_PROJECT_DIR"); d != "" {
		return d
	}
	// walk up from cwd to the nearest .workflow
	dir, err := os.Getwd()
	if err != nil {
		return "."
	}
	for d := dir; ; {
		if _, err := os.Stat(filepath.Join(d, store.DirName)); err == nil {
			return d
		}
		parent := filepath.Dir(d)
		if parent == d {
			return dir
		}
		d = parent
	}
}

func resolveSpecPath() (string, error) {
	if p := os.Getenv("WF_SPEC"); p != "" {
		return p, nil
	}
	if root := os.Getenv("CLAUDE_PLUGIN_ROOT"); root != "" {
		return filepath.Join(root, "workflow", "workflow.yaml"), nil
	}
	// dev fallback: relative to the executable (bin/wf → ../workflow/…)
	exe, err := os.Executable()
	if err == nil {
		for _, rel := range []string{
			filepath.Join(filepath.Dir(exe), "..", "workflow", "workflow.yaml"),
			filepath.Join(filepath.Dir(exe), "..", "..", "workflow", "workflow.yaml"),
		} {
			if _, err := os.Stat(rel); err == nil {
				return rel, nil
			}
		}
	}
	return "", fmt.Errorf("workflow spec not found (set WF_SPEC or CLAUDE_PLUGIN_ROOT)")
}

func loadSpecOnly() (*spec.Spec, error) {
	p, err := resolveSpecPath()
	if err != nil {
		return nil, err
	}
	return spec.Load(p, "")
}

func openCtl(projectDir string, initMode bool) (*runctl.Ctl, error) {
	specPath, err := resolveSpecPath()
	if err != nil {
		return nil, err
	}
	st, err := store.Open(projectDir, initMode)
	if err != nil {
		return nil, err
	}
	contractsDir := ""
	if _, err := os.Stat(st.ContractsDir()); err == nil {
		contractsDir = st.ContractsDir()
	}
	sp, err := spec.Load(specPath, contractsDir)
	if err != nil {
		return nil, err
	}
	cfg, err := st.LoadConfig()
	if err != nil {
		return nil, err
	}
	return &runctl.Ctl{Store: st, Spec: sp, Config: cfg}, nil
}

// ---------------------------------------------------------------------------
// Hook entry points
// ---------------------------------------------------------------------------

func gateCmd(projectDir string, rest []string) int {
	if len(rest) == 0 {
		fmt.Fprintln(os.Stderr, "wf gate stop|task-create|task-complete|verdict|skill|edit|bash")
		return 3
	}
	which := rest[0]
	in, err := hookio.Read(os.Stdin)
	if err != nil {
		in = &hookio.Input{}
	}
	// The catastrophic Bash net is store-free and always-on — it protects
	// un-adopted projects too and must not depend on state or spec.
	if which == "bash" {
		return gates.Bash(nil, in).Emit(os.Stdout, os.Stderr)
	}
	ctl, err := openCtl(projectDir, false)
	if err != nil {
		// An un-adopted project has no workflow to enforce: allow silently.
		if errors.Is(err, store.ErrNotInitialized) {
			return hookio.Allow().Emit(os.Stdout, os.Stderr)
		}
		// fail-safe split: sequencing gates open+loud, data gates closed
		switch which {
		case "task-create", "task-complete", "verdict":
			return hookio.Block("wf engine unavailable (fail-closed gate): "+err.Error()).Emit(os.Stdout, os.Stderr)
		default:
			return hookio.BrokenGate(err).Emit(os.Stdout, os.Stderr)
		}
	}
	var res hookio.Result
	switch which {
	case "stop":
		res = gates.Stop(ctl, in)
	case "task-create":
		res = gates.TaskCreated(ctl, in)
	case "task-complete":
		res = gates.TaskCompleted(ctl, in)
	case "verdict":
		res = gates.Verdict(ctl, in)
	case "skill":
		res = gates.Skill(ctl, in)
	case "edit":
		res = gates.Edit(ctl, in)
	case "bash":
		res = gates.Bash(ctl, in)
	default:
		res = hookio.Allow()
	}
	return res.Emit(os.Stdout, os.Stderr)
}

func injectCmd(projectDir string, rest []string) int {
	if len(rest) == 0 {
		fmt.Fprintln(os.Stderr, "wf inject session|turn|agent <name>")
		return 3
	}
	in, _ := hookio.Read(os.Stdin)
	if in == nil {
		in = &hookio.Input{}
	}
	ctl, err := openCtl(projectDir, false)
	if err != nil {
		if errors.Is(err, store.ErrNotInitialized) {
			return 0 // project not adopted: injections stay silent
		}
		return hookio.BrokenGate(err).Emit(os.Stdout, os.Stderr)
	}
	var payload string
	event := "SessionStart"
	switch rest[0] {
	case "session":
		payload, err = inject.Session(ctl)
	case "turn":
		event = "UserPromptSubmit"
		payload, err = inject.Turn(ctl)
	case "agent":
		event = "SubagentStart"
		name := in.AgentType
		if len(rest) > 1 {
			name = rest[1]
		}
		payload, err = inject.Agent(ctl, strings.TrimPrefix(name, "wf:"))
	default:
		return 3
	}
	if err != nil {
		return hookio.BrokenGate(err).Emit(os.Stdout, os.Stderr)
	}
	if strings.TrimSpace(payload) == "" {
		return 0
	}
	return hookio.AdditionalContext(event, payload).Emit(os.Stdout, os.Stderr)
}

func captureCmd(projectDir string, rest []string) int {
	if len(rest) == 0 {
		return 3
	}
	in, err := hookio.Read(os.Stdin)
	if err != nil {
		return 0 // capture never breaks the loop
	}
	ctl, err := openCtl(projectDir, false)
	if err != nil {
		return 0
	}
	switch rest[0] {
	case "bash", "test": // one Bash hook entry captures tests AND commits
		res := gates.CaptureTest(ctl, in)
		if r2 := gates.CaptureCommit(ctl, in); r2.Stdout != "" {
			res = r2
		}
		return res.Emit(os.Stdout, os.Stderr)
	case "edit":
		return gates.CaptureEdit(ctl, in).Emit(os.Stdout, os.Stderr)
	}
	return 0
}

// ---------------------------------------------------------------------------
// Agent/user commands
// ---------------------------------------------------------------------------

func initCmd(ctl *runctl.Ctl) int {
	cfg := ctl.Config
	if cfg.PluginVersion == "" {
		cfg.PluginVersion = Version
	}
	if err := ctl.Store.SaveConfig(cfg); err != nil {
		fmt.Fprintln(os.Stderr, "wf:", err)
		return 3
	}
	fmt.Println("wf: initialized .workflow/ (config, log, contracts.d)")
	return 0
}

func runCmd(ctl *runctl.Ctl, rest []string) int {
	if len(rest) == 0 {
		fmt.Fprintln(os.Stderr, "wf run start|branch|adopt|resume|close")
		return 3
	}
	fs := flag.NewFlagSet("run", flag.ContinueOnError)
	family := fs.String("family", "", "diff|artifact|assessment")
	intent := fs.String("intent", "", "intent tag")
	reason := fs.String("reason", "", "reason")
	if err := fs.Parse(rest[1:]); err != nil {
		return 3
	}
	switch rest[0] {
	case "start":
		r, err := ctl.RunStart(*family, *intent)
		if err != nil {
			fmt.Fprintln(os.Stderr, "wf:", err)
			return 2
		}
		fmt.Printf("run %s started (%s/%s) — phase: %s\n", r.ID, r.Family, orDash(r.Intent), r.Phase)
	case "branch":
		r, err := ctl.RunBranch(*family, *intent, *reason)
		if err != nil {
			fmt.Fprintln(os.Stderr, "wf:", err)
			return 2
		}
		fmt.Printf("branched: run %s (parent %s) — phase: %s\n", r.ID, r.Parent, r.Phase)
	case "adopt":
		r, err := ctl.RunAdopt()
		if err != nil {
			fmt.Fprintln(os.Stderr, "wf:", err)
			return 2
		}
		fmt.Printf("adopted run %s (%s/%s) — phase: %s, status: %s\n", r.ID, r.Family, orDash(r.Intent), r.Phase, r.Status)
	case "resume":
		r, err := ctl.Resume()
		if err != nil {
			fmt.Fprintln(os.Stderr, "wf:", err)
			return 2
		}
		fmt.Printf("resumed run %s — phase: %s\n", r.ID, r.Phase)
	case "close":
		r, err := ctl.MustRun()
		if err != nil {
			fmt.Fprintln(os.Stderr, "wf:", err)
			return 2
		}
		if err := ctl.RunClose(); err != nil {
			fmt.Fprintln(os.Stderr, "wf:", err)
			return 2
		}
		// freeze the signals snapshot into the archive (03 §4.7)
		if path, err := views.WriteRunSignals(ctl, r.ID); err != nil {
			fmt.Fprintln(os.Stderr, "wf: signals snapshot failed (run is closed):", err)
		} else {
			fmt.Println("signals snapshot:", path)
		}
		fmt.Println("run closed and archived under .workflow/runs/")
	default:
		return 3
	}
	return 0
}

func phaseCmd(ctl *runctl.Ctl, rest []string) int {
	if len(rest) == 0 {
		fmt.Fprintln(os.Stderr, "wf phase exit [--force --reason …] | waive <phase> --reason …")
		return 3
	}
	fs := flag.NewFlagSet("phase", flag.ContinueOnError)
	force := fs.Bool("force", false, "audited bypass (escalates)")
	reason := fs.String("reason", "", "reason")
	switch rest[0] {
	case "exit":
		if err := fs.Parse(rest[1:]); err != nil {
			return 3
		}
		findings, msg, err := ctl.PhaseExit(*force, *reason)
		if err != nil {
			fmt.Fprintln(os.Stderr, "wf:", err)
			if err == runctl.ErrNoRun {
				return 2
			}
			return 3
		}
		if len(findings) > 0 {
			fmt.Fprintf(os.Stderr, "phase exit blocked — %d unmet contract item(s):\n", len(findings))
			for _, f := range findings {
				d := ""
				if f.Detail != "" {
					d = " [" + f.Detail + "]"
				}
				fmt.Fprintf(os.Stderr, "  ✗ %s → %s%s\n", f.ID, f.Remediation, d)
			}
			return 2
		}
		fmt.Println(msg)
	case "waive":
		if len(rest) < 2 {
			fmt.Fprintln(os.Stderr, "wf phase waive <phase> --reason …")
			return 3
		}
		if err := fs.Parse(rest[2:]); err != nil {
			return 3
		}
		if err := ctl.PhaseWaive(rest[1], *reason); err != nil {
			fmt.Fprintln(os.Stderr, "wf:", err)
			return 2
		}
		fmt.Printf("phase %s waived (recorded, surfaced at Ship)\n", rest[1])
	default:
		return 3
	}
	return 0
}

// recordCmd: wf record <kind> [--json '{…}'] [key=value …]
func recordCmd(ctl *runctl.Ctl, rest []string) int {
	if len(rest) == 0 {
		fmt.Fprintln(os.Stderr, "wf record <kind> [--json '{…}'] [key=value …]")
		return 3
	}
	kind := rest[0]
	data := map[string]any{}
	args := rest[1:]
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--json" && i+1 < len(args) {
			if err := json.Unmarshal([]byte(args[i+1]), &data); err != nil {
				fmt.Fprintln(os.Stderr, "wf: --json:", err)
				return 2
			}
			i++
			continue
		}
		k, v, found := strings.Cut(a, "=")
		if !found {
			fmt.Fprintf(os.Stderr, "wf: expected key=value, got %q\n", a)
			return 2
		}
		k = strings.TrimPrefix(k, "--")
		var parsed any
		if err := json.Unmarshal([]byte(v), &parsed); err == nil {
			data[k] = parsed
		} else {
			data[k] = v
		}
	}
	ev, err := ctl.Record(kind, data, false, "agent")
	if err != nil {
		fmt.Fprintln(os.Stderr, "wf:", err)
		return 2
	}
	fmt.Printf("recorded %s %s\n", kind, ev.ID)
	return 0
}

func approveCmd(ctl *runctl.Ctl, rest []string) int {
	if len(rest) == 0 {
		fmt.Fprintln(os.Stderr, "wf approve <gate> [--payload …]")
		return 3
	}
	fs := flag.NewFlagSet("approve", flag.ContinueOnError)
	payload := fs.String("payload", "", "what the user approved (hashed into the record)")
	if err := fs.Parse(rest[1:]); err != nil {
		return 3
	}
	ev, err := ctl.Approve(rest[0], *payload)
	if err != nil {
		fmt.Fprintln(os.Stderr, "wf:", err)
		return 2
	}
	fmt.Printf("approval %s recorded (%s) — self-attested, reported per run\n", rest[0], ev.ID)
	return 0
}

func contractCmd(ctl *runctl.Ctl, rest []string) int {
	if len(rest) < 2 || rest[0] != "waive" {
		fmt.Fprintln(os.Stderr, "wf contract waive <item-or-element-id> --reason …")
		return 3
	}
	fs := flag.NewFlagSet("contract", flag.ContinueOnError)
	reason := fs.String("reason", "", "why this item does not apply")
	if err := fs.Parse(rest[2:]); err != nil {
		return 3
	}
	ev, err := ctl.WaiveItem(rest[1], *reason)
	if err != nil {
		fmt.Fprintln(os.Stderr, "wf:", err)
		return 2
	}
	fmt.Printf("waiver recorded %s (surfaced in trace/report)\n", ev.ID)
	return 0
}

func loopCmd(ctl *runctl.Ctl, rest []string) int {
	fs := flag.NewFlagSet("loop", flag.ContinueOnError)
	ac := fs.String("ac", "", "the failing AC")
	cause := fs.String("cause", "", "slip|design|plan")
	evidence := fs.String("evidence", "", "observed vs expected")
	if err := fs.Parse(rest); err != nil {
		return 3
	}
	target, err := ctl.Loop(*ac, *cause, *evidence)
	if err != nil {
		fmt.Fprintln(os.Stderr, "wf:", err)
		return 2
	}
	fmt.Printf("looped to %s (recorded)\n", target)
	return 0
}

func parkCmd(ctl *runctl.Ctl, rest []string) int {
	fs := flag.NewFlagSet("park", flag.ContinueOnError)
	reason := fs.String("reason", "", "why the run stops here")
	if err := fs.Parse(rest); err != nil {
		return 3
	}
	if err := ctl.Park(*reason); err != nil {
		fmt.Fprintln(os.Stderr, "wf:", err)
		return 2
	}
	fmt.Println("run parked (honest stop; resume with wf run resume)")
	return 0
}

func riskCmd(ctl *runctl.Ctl, rest []string) int {
	if len(rest) == 0 || rest[0] != "scan" {
		fmt.Fprintln(os.Stderr, "wf risk scan [--text …] [--add signal]…")
		return 3
	}
	fs := flag.NewFlagSet("risk", flag.ContinueOnError)
	text := fs.String("text", "", "task text to screen")
	var adds sliceFlag
	fs.Var(&adds, "add", "agent-judged signal (repeatable)")
	if err := fs.Parse(rest[1:]); err != nil {
		return 3
	}
	signals, lenses, err := ctl.RiskScan(*text, adds)
	if err != nil {
		fmt.Fprintln(os.Stderr, "wf:", err)
		return 2
	}
	fmt.Printf("risk signals: %s\nlenses to work: %s\n", orNone(signals), orNone(lenses))
	return 0
}

func statusCmd(ctl *runctl.Ctl, rest []string) int {
	payload, err := inject.Session(ctl)
	if err != nil {
		fmt.Fprintln(os.Stderr, "wf:", err)
		return 3
	}
	fmt.Print(payload)
	return 0
}

func doctorCmd(ctl *runctl.Ctl, rest []string) int {
	fs := flag.NewFlagSet("doctor", flag.ContinueOnError)
	bootstrap := fs.Bool("bootstrap", false, "verify install only")
	if err := fs.Parse(rest); err != nil {
		return 3
	}
	if *bootstrap {
		fmt.Println("wf", Version, "— spec:", len(ctl.Spec.Contracts), "contract items,", len(ctl.Spec.Roster), "agents")
		return bootstrapHealth()
	}
	specPath, _ := resolveSpecPath()
	rep := doctor.Run(ctl, specPath)
	fmt.Println(rep.String())
	if !rep.OK {
		return 2
	}
	return 0
}

// bootstrapHealth checks every discoverable wf plugin install for a dead
// hook engine and heals it by running that install's bootstrap script — a
// mid-session /plugin install never fires SessionStart, so without this
// every gate stays ENOENT-dead until the next session (the power-of-ten
// incident: a full run shipped with zero hook events). Exit 0 when healthy
// or healed, 2 when hooks remain dead.
func bootstrapHealth() int {
	findings, dead := doctor.HookEngineFindings("", true)
	for _, f := range findings {
		fmt.Println("  -", f)
	}
	if dead {
		return 2
	}
	return 0
}

// reportCmd — the health-signals view (08 §4): aggregate across archived +
// active runs by default; --run <id|current> for one run; --json for both.
func reportCmd(ctl *runctl.Ctl, rest []string) int {
	fs := flag.NewFlagSet("report", flag.ContinueOnError)
	asJSON := fs.Bool("json", false, "machine-readable output")
	runID := fs.String("run", "", "one run's snapshot: an archive ID or \"current\"")
	if err := fs.Parse(rest); err != nil {
		return 3
	}
	if *runID != "" {
		s, err := views.ReportRun(ctl, *runID)
		if err != nil {
			fmt.Fprintln(os.Stderr, "wf:", err)
			return 2
		}
		if *asJSON {
			raw, _ := json.MarshalIndent(s, "", "  ")
			fmt.Println(string(raw))
		} else {
			fmt.Print(views.RenderRunSignals(s))
		}
		return 0
	}
	sigs, err := views.Report(ctl)
	if err != nil {
		fmt.Fprintln(os.Stderr, "wf:", err)
		return 2
	}
	if *asJSON {
		raw, _ := json.MarshalIndent(sigs, "", "  ")
		fmt.Println(string(raw))
	} else {
		fmt.Print(views.RenderReport(sigs))
	}
	return 0
}

// lessonsCmd — the enforcement loop for what runs teach (03 §4.7):
// suggest (engine-proposed from signals), accept/reject (user disposition +
// approval + regeneration), apply (idempotent regeneration of
// contracts.d/lessons.yaml + .claude/rules/wf-lessons.md).
func lessonsCmd(ctl *runctl.Ctl, projectDir string, rest []string) int {
	if len(rest) == 0 {
		fmt.Fprintln(os.Stderr, "wf lessons suggest|accept <id>|reject <id>|apply")
		return 3
	}
	specPath, err := resolveSpecPath()
	if err != nil {
		fmt.Fprintln(os.Stderr, "wf:", err)
		return 3
	}
	var out string
	switch rest[0] {
	case "suggest":
		out, err = ops.LessonsSuggest(ctl)
	case "accept", "reject":
		if len(rest) < 2 {
			fmt.Fprintf(os.Stderr, "wf lessons %s <lesson-record-id>\n", rest[0])
			return 3
		}
		if rest[0] == "accept" {
			out, err = ops.LessonsAccept(ctl, projectDir, specPath, rest[1])
		} else {
			out, err = ops.LessonsReject(ctl, projectDir, specPath, rest[1])
		}
	case "apply":
		out, err = ops.LessonsApply(ctl, projectDir, specPath)
	default:
		fmt.Fprintln(os.Stderr, "wf lessons suggest|accept <id>|reject <id>|apply")
		return 3
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "wf:", err)
		return 2
	}
	fmt.Println(out)
	return 0
}

func traceCmd(ctl *runctl.Ctl) int {
	report, err := views.Trace(ctl)
	if err != nil {
		fmt.Fprintln(os.Stderr, "wf:", err)
		return 2
	}
	fmt.Print(report)
	return 0
}

func depsCmd(ctl *runctl.Ctl, projectDir string, rest []string) int {
	if len(rest) == 0 || rest[0] != "check" {
		fmt.Fprintln(os.Stderr, "wf deps check")
		return 3
	}
	out, err := ops.DepsCheck(ctl, projectDir)
	if err != nil {
		fmt.Fprintln(os.Stderr, "wf:", err)
		return 2
	}
	fmt.Println(out)
	if strings.HasPrefix(out, "deps: missing") {
		return 2
	}
	return 0
}

func originCmd(ctl *runctl.Ctl, projectDir string, rest []string) int {
	if len(rest) == 0 || rest[0] != "discover" {
		fmt.Fprintln(os.Stderr, "wf origin discover [--path …] [--text …]")
		return 3
	}
	fs := flag.NewFlagSet("origin", flag.ContinueOnError)
	path := fs.String("path", "", "file whose history to follow")
	text := fs.String("text", "", "code fragment to pickaxe (-S)")
	if err := fs.Parse(rest[1:]); err != nil {
		return 3
	}
	out, err := ops.OriginDiscover(ctl, projectDir, *path, *text)
	if err != nil {
		fmt.Fprintln(os.Stderr, "wf:", err)
		return 2
	}
	fmt.Println(out)
	return 0
}

func docCmd(ctl *runctl.Ctl, projectDir string, rest []string) int {
	if len(rest) < 2 || rest[0] != "new" {
		fmt.Fprintln(os.Stderr, "wf doc new adr|design|threat-model|ux|review|incident|release-notes|delivery-manifest --slug …")
		return 3
	}
	fs := flag.NewFlagSet("doc", flag.ContinueOnError)
	slug := fs.String("slug", "", "kebab-case name for the document")
	if err := fs.Parse(rest[2:]); err != nil {
		return 3
	}
	root := os.Getenv("CLAUDE_PLUGIN_ROOT")
	if root == "" {
		root = ctl.Spec.PluginRoot()
	}
	out, err := ops.DocNew(ctl, root, projectDir, rest[1], *slug)
	if err != nil {
		fmt.Fprintln(os.Stderr, "wf:", err)
		return 2
	}
	fmt.Println(out)
	return 0
}

// ---------------------------------------------------------------------------

type sliceFlag []string

func (s *sliceFlag) String() string     { return strings.Join(*s, ",") }
func (s *sliceFlag) Set(v string) error { *s = append(*s, v); return nil }

func orDash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func orNone(xs []string) string {
	if len(xs) == 0 {
		return "(none)"
	}
	return strings.Join(xs, ", ")
}

func usage() {
	fmt.Print(`wf — enforced development workflow engine

run lifecycle:   wf run start --family diff|artifact|assessment [--intent …]
                 wf run branch|adopt|resume|close
phases:          wf phase exit [--force --reason …]     (exit 0 ok / 2 gaps / 3 broken)
                 wf phase waive <phase> --reason …
records:         wf record <kind> [--json '{…}'] [key=value …]
                 wf approve <gate> [--payload …]
                 wf contract waive <id> --reason …
                 wf loop --ac AC-1 --cause slip|design|plan --evidence …
                 wf park --reason …
                 wf risk scan [--text …] [--add signal]…
grounding:       wf deps check · wf origin discover [--path …] [--text …]
                 wf doc new <type> --slug … · wf trace
lessons:         wf lessons suggest|accept <id>|reject <id>|apply
introspection:   wf status · wf report [--json] [--run <id|current>]
                 wf doctor [--bootstrap] · wf selftest · wf version
hook entries:    wf gate stop|task-create|task-complete|verdict|skill|edit|bash
                 wf inject session|turn|agent <name> · wf capture bash|edit
`)
}
