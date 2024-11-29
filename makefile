REPO = github.com/kakabei/wx-proxy-service

GIT_COMMIT := $(shell git show-branch --no-name HEAD)
GIT_BRANCH := $(shell git rev-parse --abbrev-ref HEAD)
GIT_DIRTY := $(shell test -n "`git status --porcelain`" && echo "*" || true)
BUILD_VERSION := $(shell git describe --abbrev=10 --tags --always)
BUILD_TIME := $(shell date +%FT%T%z)

LDFLAGS := "\
-X \"${REPO}/internal/version.buildGitCommit=${GIT_COMMIT} ${GIT_DIRTY}\" \
-X \"${REPO}/internal/version.buildGitBranch=${GIT_BRANCH}\" \
-X \"${REPO}/internal/version.buildVersion=${BUILD_VERSION}\" \
-X \"${REPO}/internal/version.buildTime=${BUILD_TIME}\""

export CGO_ENABLED=0

all: release

.PHONY: debug
debug:
	go vet ./...
	go build  -ldflags $(LDFLAGS)

api:
	goctl api go -api wx-proxy.api -dir . -style goZero

.PHONY: release
release:
	go vet ./...
	GOWORK=off go build -ldflags $(LDFLAGS)

.PHONY: docs
docs:
	swag init
	swagger generate markdown -f docs/swagger.yaml --output Readme.md

