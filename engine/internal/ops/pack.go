// Contract packs (09 Q1): a pack is a directory (or single file) of
// contracts.d-shaped YAML — record kinds namespaced x-*, contract items
// namespaced local.* — plus any documentation (.md, copied to
// .workflow/packs/<pack>/ so pack-referencing agents like the
// compliance-reviewer can read their checklists project-side). `wf pack
// install` validates the merged result strictly against a temp mirror of
// contracts.d (the same zero-risk pattern as lesson acceptance), then
// copies the files in. The merge stays add-only: a pack can never weaken
// or replace shipped items.
package ops

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/runctl"
	"github.com/dariusz-klibisz/ClaudeWorkflow/engine/internal/spec"
)

var packNameRe = regexp.MustCompile(`[^a-z0-9]+`)

// PackInstall installs the pack at src (file or directory) into the
// project's contracts.d.
func PackInstall(c *runctl.Ctl, specPath, src string) (string, error) {
	st, err := os.Stat(src)
	if err != nil {
		return "", fmt.Errorf("pack source: %w", err)
	}
	pack := strings.Trim(packNameRe.ReplaceAllString(
		strings.ToLower(strings.TrimSuffix(filepath.Base(src), filepath.Ext(src))), "-"), "-")
	if pack == "" {
		return "", fmt.Errorf("cannot derive a pack name from %q", src)
	}

	// collect the pack's yaml files (+ .md docs when a directory)
	var files, docs []string
	if st.IsDir() {
		entries, err := os.ReadDir(src)
		if err != nil {
			return "", err
		}
		for _, e := range entries {
			if e.IsDir() {
				continue
			}
			switch {
			case strings.HasSuffix(e.Name(), ".yaml") || strings.HasSuffix(e.Name(), ".yml"):
				files = append(files, filepath.Join(src, e.Name()))
			case strings.HasSuffix(e.Name(), ".md"):
				docs = append(docs, filepath.Join(src, e.Name()))
			}
		}
	} else {
		if !strings.HasSuffix(src, ".yaml") && !strings.HasSuffix(src, ".yml") {
			return "", fmt.Errorf("a pack file must be .yaml/.yml")
		}
		files = []string{src}
	}
	if len(files) == 0 {
		return "", fmt.Errorf("no yaml files in pack %s", src)
	}
	sort.Strings(files)

	// destination names: <pack>-<original>.yaml — collision = refuse
	dests := make(map[string]string, len(files)) // dest name -> src path
	for _, f := range files {
		name := filepath.Base(f)
		base := strings.TrimSuffix(name, filepath.Ext(name))
		if base != pack && !strings.HasPrefix(base, pack+"-") {
			name = pack + "-" + name
		}
		if _, err := os.Stat(filepath.Join(c.Store.ContractsDir(), name)); err == nil {
			return "", fmt.Errorf("contracts.d/%s already exists — uninstall/rename first (add-only merge)", name)
		}
		dests[name] = f
	}

	// strict validation against a temp mirror of contracts.d + the pack
	tmp, err := os.MkdirTemp("", "wf-pack-*")
	if err != nil {
		return "", err
	}
	defer os.RemoveAll(tmp)
	if entries, err := os.ReadDir(c.Store.ContractsDir()); err == nil {
		for _, e := range entries {
			name := e.Name()
			if e.IsDir() || (!strings.HasSuffix(name, ".yaml") && !strings.HasSuffix(name, ".yml")) {
				continue
			}
			raw, err := os.ReadFile(filepath.Join(c.Store.ContractsDir(), name))
			if err != nil {
				return "", err
			}
			if err := os.WriteFile(filepath.Join(tmp, name), raw, 0o644); err != nil {
				return "", err
			}
		}
	}
	for name, srcPath := range dests {
		raw, err := os.ReadFile(srcPath)
		if err != nil {
			return "", err
		}
		if err := os.WriteFile(filepath.Join(tmp, name), raw, 0o644); err != nil {
			return "", err
		}
	}
	if _, err := spec.LoadStrict(specPath, tmp); err != nil {
		return "", fmt.Errorf("pack %s rejected (nothing installed): %w", pack, err)
	}

	// install
	names := make([]string, 0, len(dests))
	for name, srcPath := range dests {
		raw, err := os.ReadFile(srcPath)
		if err != nil {
			return "", err
		}
		if err := os.WriteFile(filepath.Join(c.Store.ContractsDir(), name), raw, 0o644); err != nil {
			return "", err
		}
		names = append(names, name)
	}
	sort.Strings(names)

	// docs travel with the pack: .workflow/packs/<pack>/ (committed) — the
	// compliance-reviewer's injected checklist paths point here
	docNote := ""
	if len(docs) > 0 {
		packDir := filepath.Join(c.Store.Root, "packs", pack)
		if err := os.MkdirAll(packDir, 0o755); err != nil {
			return "", err
		}
		for _, d := range docs {
			raw, err := os.ReadFile(d)
			if err != nil {
				return "", err
			}
			if err := os.WriteFile(filepath.Join(packDir, filepath.Base(d)), raw, 0o644); err != nil {
				return "", err
			}
		}
		docNote = fmt.Sprintf(" · %d doc(s) → .workflow/packs/%s/", len(docs), pack)
	}
	return fmt.Sprintf("pack %s installed: contracts.d/{%s}%s — commit .workflow/contracts.d (and packs/) so the pack travels with the repo (remove the files to uninstall)",
		pack, strings.Join(names, ", "), docNote), nil
}
