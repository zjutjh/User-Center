GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

GOCTL ?= goctl
GOCTL_VERSION ?= v1.10.1
STYLE ?= go_zero

API_DIR := apps/user-api
API_FILE := $(API_DIR)/user.api
RPC_DIR := apps/user-rpc
RPC_PROTO := $(RPC_DIR)/user.proto
MODULE := github.com/zjutjh/User-Center

.PHONY: build
build: ## Build the combined user center service.
	go build -o bin/usercenter .

.PHONY: build-api
build-api: ## Build the user API service.
	go build -o bin/user-api ./$(API_DIR)

.PHONY: build-rpc
build-rpc: ## Build the user RPC service.
	go build -o bin/user-rpc ./$(RPC_DIR)

.PHONY: generate
generate: generate-api generate-rpc ## Generate api and rpc code with goctl.

.PHONY: generate-api
generate-api: ## Generate HTTP service code from user.api.
	$(GOCTL) api go --api $(API_FILE) --dir $(API_DIR) --style $(STYLE)

.PHONY: generate-rpc
generate-rpc: ## Generate RPC service code from user.proto.
	$(GOCTL) rpc protoc $(RPC_PROTO) --go_out=$(RPC_DIR) --go-grpc_out=$(RPC_DIR) --zrpc_out=$(RPC_DIR) --style $(STYLE) --module $(MODULE)

.PHONY: generate-model
generate-model: ## Generate DAO model/query code.
	go run ./cmd/gen

.PHONY: configure
configure: tools ## Install the goctl toolchain.

.PHONY: tools
tools: ## Install goctl and protoc plugins used by goctl rpc protoc.
	go install github.com/zeromicro/go-zero/tools/goctl@$(GOCTL_VERSION)
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

.PHONY: fmt
fmt: ## Run go fmt against code.
	$(GOCTL) api format --dir $(API_DIR)
	go fmt ./...

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: test
test: ## Run tests.
	go test ./...

.PHONY: run
run: ## Run the combined user center service.
	go run .

.PHONY: run-api
run-api: ## Run the user API service.
	go run ./$(API_DIR)

.PHONY: run-rpc
run-rpc: ## Run the user RPC service.
	go run ./$(RPC_DIR)
