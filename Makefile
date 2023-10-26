build:
	go build -o user_center

build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 GIN_MODE=release go build -o user_center

.PHONY: build build-linux