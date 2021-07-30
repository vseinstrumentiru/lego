# Dependency versions
PROTOC_VERSION = 3.14.0
BUF_VERSION = 0.35.1
PROTOC_GEN_KIT_VERSION = 0.2.0

bin/protoc: bin/protoc-${PROTOC_VERSION}
	@ln -sf protoc-${PROTOC_VERSION}/bin/protoc bin/protoc
bin/protoc-${PROTOC_VERSION}:
	@mkdir -p bin/protoc-${PROTOC_VERSION}
ifeq (${OS}, darwin)
	curl -L https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-osx-x86_64.zip > bin/protoc.zip
endif
ifeq (${OS}, linux)
	curl -L https://github.com/protocolbuffers/protobuf/releases/download/v${PROTOC_VERSION}/protoc-${PROTOC_VERSION}-linux-x86_64.zip > bin/protoc.zip
endif
	unzip bin/protoc.zip -d bin/protoc-${PROTOC_VERSION}
	rm bin/protoc.zip

#bin/protoc-gen-go: bin/protoc-gen-go-${PROTOC_GEN_GO_VERSION}
#	@ln -sf protoc-gen-go-${PROTOC_GEN_GO_VERSION} bin/protoc-gen-go
bin/protoc-gen-go: gotools.mod
	@mkdir -p bin
	go build -modfile gotools.mod -o bin/protoc-gen-go github.com/golang/protobuf/protoc-gen-go

bin/protoc-gen-go-grpc: gotools.mod
	@mkdir -p bin
	go build -modfile gotools.mod -o bin/protoc-gen-go-grpc google.golang.org/grpc/cmd/protoc-gen-go-grpc

bin/protoc-gen-kit: bin/protoc-gen-kit-${PROTOC_GEN_KIT_VERSION}
	@ln -sf protoc-gen-kit-${PROTOC_GEN_KIT_VERSION} bin/protoc-gen-kit
bin/protoc-gen-kit-${PROTOC_GEN_KIT_VERSION}:
	@mkdir -p bin
	curl -L https://github.com/sagikazarmark/protoc-gen-kit/releases/download/v${PROTOC_GEN_KIT_VERSION}/protoc-gen-kit_${OS}_amd64.tar.gz | tar -zOxf - protoc-gen-kit > ./bin/protoc-gen-kit-${PROTOC_GEN_KIT_VERSION} && chmod +x ./bin/protoc-gen-kit-${PROTOC_GEN_KIT_VERSION}

bin/buf: bin/buf-${BUF_VERSION}
	@ln -sf buf-${BUF_VERSION} bin/buf
bin/buf-${BUF_VERSION}:
	@mkdir -p bin
	curl -L https://github.com/bufbuild/buf/releases/download/v${BUF_VERSION}/buf-${OS}-x86_64 -o ./bin/buf-${BUF_VERSION} && chmod +x ./bin/buf-${BUF_VERSION}

.PHONY: buf
buf: bin/buf ## Generate client and server stubs from the protobuf definition
	bin/buf build -o /dev/null
	bin/buf lint

.PHONY: proto
proto: bin/protoc bin/protoc-gen-go bin/protoc-gen-go-grpc bin/protoc-gen-kit buf ## Generate client and server stubs from the protobuf definition
	bin/buf build -o - | protoc --descriptor_set_in=/dev/stdin --go_out=paths=source_relative:api --go-grpc_out=paths=source_relative:api --kit_out=paths=source_relative:api $(shell bin/buf build -o - | bin/buf ls-files --input - | grep -v google)
