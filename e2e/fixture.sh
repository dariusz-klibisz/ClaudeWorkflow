#!/bin/sh
# Builds the throwaway E2E fixture repo: git repo, wf-adopted (.workflow/,
# settings, CLAUDE.md block), pointing at the checkout as the plugin.
# Usage: fixture.sh <dir> <repo-checkout> <wf-binary>
set -eu

FIXTURE="$1"; REPO="$2"; WF="$3"

cd "$FIXTURE"
git init -q
git config user.email e2e@wf.test
git config user.name "wf e2e"

mkdir -p .claude
cat > .claude/settings.json <<'EOF'
{
  "permissions": { "allow": ["Bash(wf *)"] }
}
EOF

# adopt: wf init + the CLAUDE.md block (what /wf:init would materialize)
CLAUDE_PLUGIN_ROOT="$REPO" "$WF" init >/dev/null

cat > CLAUDE.md <<'EOF'
<!-- wf:begin -->
## Workflow (wf)
- Work happens inside runs: start every task with /wf:dev.
- .workflow/ on disk is the source of truth; after compaction or resume,
  trust the injected [wf] status block over memory.
- Record facts only via wf commands — never edit .workflow/state|log by hand.
- Deliverable documents go under docs/.
- Audited escapes: /wf:park (honest stop), /wf:force (bypass, escalates).
<!-- wf:end -->
EOF

printf '# e2e fixture\n' > README.md
git add -A
git commit -qm "e2e fixture"
