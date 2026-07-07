# wf plugin build targets
VERSION := $(shell grep '"version"' .claude-plugin/plugin.json | head -1 | sed 's/[^0-9.]*//g')
# Full version includes the git sha: the SessionStart bootstrap compares
# bin/VERSION strings, so every engine change must change the string —
# otherwise a stale binary lingers in CLAUDE_PLUGIN_DATA (the corpus-field
# incident).
GITSHA := $(shell git rev-parse --short HEAD 2>/dev/null || echo dev)
DIRTY := $(shell git diff --quiet 2>/dev/null || echo .dirty-$(shell date +%s))
# RELEASE=1: stamp the plain semver and build reproducibly (-trimpath, no
# sha/dirty suffix). Required for the fetch tier: bin/MANIFEST commits the
# checksums BEFORE the tag exists, so CI's rebuild at the tag must be
# byte-identical to the maintainer's `make dist RELEASE=1` — the sha suffix
# would break that (the manifest commit changes the sha).
ifeq ($(RELEASE),1)
VERSION_FULL := $(VERSION)
BUILDFLAGS := -trimpath
# byte-identical rebuilds need the EXACT toolchain from go.mod — GOTOOLCHAIN
# auto would silently use a newer installed Go and change the binaries.
export GOTOOLCHAIN := go$(shell sed -n 's/^go //p' engine/go.mod)
else
VERSION_FULL := $(VERSION)+$(GITSHA)$(DIRTY)
BUILDFLAGS :=
endif
PLATFORMS := linux-amd64 linux-arm64 darwin-amd64 darwin-arm64 windows-amd64 windows-arm64
LDFLAGS := -s -w -X main.Version=$(VERSION_FULL)
BASE_URL := https://github.com/dariusz-klibisz/ClaudeWorkflow/releases/download

.PHONY: test gen check build dist manifest verify-manifest clean

test:
	cd engine && go vet ./... && go test -race ./...

gen:
	cd engine && go run ./gen

check:
	cd engine && go run ./gen -check

# local dev build (current platform, into bin/ where the bootstrap finds it).
# Removes other-platform binaries first: `build` only rebuilds the current
# platform, so leftovers from an earlier `dist` would carry OLD embedded
# versions yet get freshly valid checksums (the mixed-vintage SHA256SUMS
# incident — stale 0.1.0 darwin/windows binaries shipped alongside a HEAD
# linux build). Cross-platform sets come from `dist` only.
# Note: VERSION_FULL is stamped at BUILD time — building before committing
# bakes in the previous sha + .dirty; re-run `make build` after committing.
build:
	rm -f bin/wf-* bin/SHA256SUMS
	cd engine && go build -ldflags '$(LDFLAGS)' -o ../bin/wf-$$(go env GOOS)-$$(go env GOARCH)$$(go env GOEXE) ./cmd/wf
	printf '%s' "$(VERSION_FULL)" > bin/VERSION
	cd bin && sha256sum wf-* > SHA256SUMS

# release: all platforms + checksums (RELEASE=1 for tag-bound artifacts)
dist:
	@for p in $(PLATFORMS); do \
		os=$${p%-*}; arch=$${p#*-}; ext=""; [ "$$os" = windows ] && ext=".exe"; \
		echo "building wf-$$p$$ext"; \
		cd engine && CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build $(BUILDFLAGS) -ldflags '$(LDFLAGS)' -o ../bin/wf-$$p$$ext ./cmd/wf && cd ..; \
	done
	printf '%s' "$(VERSION_FULL)" > bin/VERSION
	cd bin && sha256sum wf-* > SHA256SUMS

# manifest: the COMMITTED fetch manifest (bin/MANIFEST) — version, download
# base, per-platform sha256 of a `make dist RELEASE=1` set. The bootstrap's
# fetch tier verifies downloads against these committed sums (07 §4-B).
# Release flow: bump plugin.json → make dist RELEASE=1 manifest → commit →
# git tag v$(VERSION) → push tag (release.yml rebuilds, verifies, publishes).
manifest:
	@test "$(RELEASE)" = "1" || { echo "manifest requires RELEASE=1 (reproducible sums)"; exit 1; }
	@{ printf 'version %s\n' "$(VERSION)"; \
	   printf 'base_url %s\n' "$(BASE_URL)"; \
	   cd bin && sha256sum wf-linux-* wf-darwin-* wf-windows-*; } > bin/MANIFEST
	@echo "wrote bin/MANIFEST for v$(VERSION)"

# verify-manifest: CI release gate — a rebuild at the tag must reproduce the
# committed sums exactly (supply-chain check; also catches a stale manifest).
verify-manifest:
	@test -f bin/MANIFEST || { echo "no bin/MANIFEST"; exit 1; }
	@grep '^version $(VERSION)$$' bin/MANIFEST >/dev/null || { echo "MANIFEST version != plugin.json $(VERSION)"; exit 1; }
	@cd bin && grep ' wf-' MANIFEST | sha256sum -c -
	@echo "MANIFEST verified against built artifacts"

clean:
	rm -f bin/wf-* bin/VERSION bin/SHA256SUMS
