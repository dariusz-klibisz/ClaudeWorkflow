# wf plugin build targets
VERSION := $(shell grep '"version"' .claude-plugin/plugin.json | head -1 | sed 's/[^0-9.]*//g')
# Full version includes the git sha: the SessionStart bootstrap compares
# bin/VERSION strings, so every engine change must change the string —
# otherwise a stale binary lingers in CLAUDE_PLUGIN_DATA (the corpus-field
# incident).
GITSHA := $(shell git rev-parse --short HEAD 2>/dev/null || echo dev)
DIRTY := $(shell git diff --quiet 2>/dev/null || echo .dirty-$(shell date +%s))
VERSION_FULL := $(VERSION)+$(GITSHA)$(DIRTY)
PLATFORMS := linux-amd64 linux-arm64 darwin-amd64 darwin-arm64 windows-amd64 windows-arm64
LDFLAGS := -s -w -X main.Version=$(VERSION_FULL)

.PHONY: test gen check build dist clean

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

# release: all platforms + checksums
dist:
	@for p in $(PLATFORMS); do \
		os=$${p%-*}; arch=$${p#*-}; ext=""; [ "$$os" = windows ] && ext=".exe"; \
		echo "building wf-$$p$$ext"; \
		cd engine && CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -ldflags '$(LDFLAGS)' -o ../bin/wf-$$p$$ext ./cmd/wf && cd ..; \
	done
	printf '%s' "$(VERSION_FULL)" > bin/VERSION
	cd bin && sha256sum wf-* > SHA256SUMS

clean:
	rm -f bin/wf-* bin/VERSION bin/SHA256SUMS
