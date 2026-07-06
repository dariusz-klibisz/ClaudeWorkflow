# wf plugin build targets
VERSION := $(shell grep '"version"' .claude-plugin/plugin.json | head -1 | sed 's/[^0-9.]*//g')
PLATFORMS := linux-amd64 linux-arm64 darwin-amd64 darwin-arm64 windows-amd64 windows-arm64
LDFLAGS := -s -w -X main.Version=$(VERSION)

.PHONY: test gen check build dist clean

test:
	cd engine && go vet ./... && go test -race ./...

gen:
	cd engine && go run ./gen

check:
	cd engine && go run ./gen -check

# local dev build (current platform, into bin/ where the bootstrap finds it)
build:
	cd engine && go build -ldflags '$(LDFLAGS)' -o ../bin/wf-$$(go env GOOS)-$$(go env GOARCH)$$(go env GOEXE) ./cmd/wf
	printf '%s' "$(VERSION)" > bin/VERSION

# release: all platforms + checksums
dist:
	@for p in $(PLATFORMS); do \
		os=$${p%-*}; arch=$${p#*-}; ext=""; [ "$$os" = windows ] && ext=".exe"; \
		echo "building wf-$$p$$ext"; \
		cd engine && CGO_ENABLED=0 GOOS=$$os GOARCH=$$arch go build -ldflags '$(LDFLAGS)' -o ../bin/wf-$$p$$ext ./cmd/wf && cd ..; \
	done
	printf '%s' "$(VERSION)" > bin/VERSION
	cd bin && sha256sum wf-* > SHA256SUMS

clean:
	rm -f bin/wf-* bin/VERSION bin/SHA256SUMS
