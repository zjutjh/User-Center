# Global variables
include version.mk

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Git information
GIT_VERSION ?= $(shell git describe --tags --abbrev=8 --dirty) # attention: gitlab CI: git fetch should not use shallow
GIT_COMMIT_HASH ?= $(shell git rev-parse HEAD)
GIT_TREESTATE = "clean"
GIT_DIFF = $(shell git diff --quiet >/dev/null 2>&1; if [ $$? -eq 1 ]; then echo "1"; fi)
ifeq ($(GIT_DIFF), 1)
    GIT_TREESTATE = "dirty"
endif
BUILDDATE = $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')

LDFLAGS := "-X github.com/dynamia-ai/kantaloupe/pkg/version.gitVersion=$(GIT_VERSION) \
            -X github.com/dynamia-ai/kantaloupe/pkg/version.gitCommit=$(GIT_COMMIT_HASH) \
            -X github.com/dynamia-ai/kantaloupe/pkg/version.gitTreeState=$(GIT_TREESTATE) \
            -X github.com/dynamia-ai/kantaloupe/pkg/version.buildDate=$(BUILDDATE)"

# Build
.PHONY: apiserver
apiserver:
	go build -ldflags $(LDFLAGS) -o bin/usercenter cmd/apiserver/main.go

# generate code
.PHONY: genproto
genproto:
	cd ./api/ && $(MAKE) genproto

# generate code
.PHONY: install
install:
	cd ./api/ && $(MAKE) install_grpc_dep

.PHONY: genswagger
genswagger:
	cd ./api/ && $(MAKE)  genswagger

.PHONY: gen-code-gen
gen-code-gen:
	cd ./api/ && $(MAKE) gen-code-gen

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: run
run: ## Run the application.
	go run cmd/apiserver/main.go