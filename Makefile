SHELL := /bin/bash

GO_ROOT := $(shell go env GOROOT)
GIT_SHA := $(shell git rev-parse HEAD)
GIT_SHA_SHORT := $(shell git rev-parse --short HEAD)
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
VERSION := $(shell git describe --tags)-$(GIT_SHA_SHORT)
LDFLAGS := -s -w \
	-X 'github.com/okareo-ai/okareo-cli/version.BuildDate=$(DATE)' \
	-X 'github.com/okareo-ai/okareo-cli/version.BuildVersion=$(subst v,,$(VERSION))' \
	-X 'github.com/okareo-ai/okareo-cli/version.Commit=$(GIT_SHA)'

.PHONY: build
build:
	go build -o ./dist/bin/okareo -ldflags="$(LDFLAGS)" main.go

.PHONY: test
test: build
	@TZ=UTC go test ./...

.PHONY: fmt
fmt:
	@gofumpt -w .

.PHONY: lint
lint:
	@revive -config revive.toml -formatter stylish ./...

.PHONY: install/dev
install/dev:
	go install github.com/mgechev/revive@v1.3.2
	go install github.com/securego/gosec/v2/cmd/gosec@v2.17.0
	go install honnef.co/go/tools/cmd/staticcheck@v0.4.5
	go install mvdan.cc/gofumpt@v0.5.0

.PHONY: install/goreleaser
install/goreleaser:
	go install github.com/goreleaser/goreleaser@v1.20.0

.PHONY: proto/clean
proto/clean:
	rm -rf gen/proto

.PHONY: release
release: install/goreleaser
	@goreleaser check
	@goreleaser release --snapshot --clean

.PHONY: release/publish
release/publish: install/goreleaser
	@goreleaser release

.PHONY: run
generate:
	go run main.go run
