GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)

.PHONY: build
build:
	go build  -o bin/usercenter main.go

.PHONY: generate
generate: ## Run buf generate including code and swagger generate.
	cd api && buf generate

.PHONY: configure
configure: ## Run configure
	cd api && buf dep update

.PHONY: fmt
fmt: ## Run go fmt against code.
	go fmt ./...
	clang-format -i $$(find . -name '*.proto') --style="{BasedOnStyle: Google, IndentWidth: 4, ColumnLimit: 0, AlignConsecutiveAssignments: true}"

.PHONY: vet
vet: ## Run go vet against code.
	go vet ./...

.PHONY: run
run: ## Run the application.
	go run main.go