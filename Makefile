# Global variables
include version.mk

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

# Build
.PHONY: apiserver
apiserver:
	go build  -o bin/usercenter main.go

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
	go run main.go