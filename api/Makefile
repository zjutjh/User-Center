API_PKG    := ./api
.PHONY: clean
clean:
	hack/clean.sh

.PHONY: genproto
genproto:
	hack/genproto.sh

.PHONY: genswagger
genswagger:
	hack/gen-openapi.sh

.PHONY: install_grpc_dep
install_grpc_dep:
	go env -w GOPROXY=https://goproxy.cn,direct && \
	go install \
    	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway \
    	github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2 \
    	google.golang.org/protobuf/cmd/protoc-gen-go \
    	google.golang.org/grpc/cmd/protoc-gen-go-grpc && \
	go install github.com/golang/protobuf/protoc-gen-go@v1.5.4 && \
	go install github.com/grpc-ecosystem/protoc-gen-grpc-gateway-ts@v1.1.2
