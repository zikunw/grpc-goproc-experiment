PROTO_SRCS := $(shell find . -name *.proto)

deps:
	go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.28.1
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.2
	go install github.com/planetscale/vtprotobuf/cmd/protoc-gen-go-vtproto@latest

# Build pb
# Common fix `export GOBIN=~/go/bin` (can't be run from Makefile)
pb: #internal/message/rpc.proto # we have this commented out because protoc doesn't really work on Chameleon, so we need to run it locally and check the generated code into version control :vomit:
	@for PROTO in $(PROTO_SRCS); do \
		protoc \
			--go_out=. --go_opt=paths=source_relative\
			--plugin protoc-gen-go="${GOBIN}/protoc-gen-go" \
			--go-grpc_out=. --go-grpc_opt=paths=source_relative\
			--plugin protoc-gen-go-grpc="${GOBIN}/protoc-gen-go-grpc" \
			--go-vtproto_out=. --go-vtproto_opt=paths=source_relative\
			--plugin protoc-gen-go-vtproto="${GOBIN}/protoc-gen-go-vtproto" \
			--go-vtproto_opt=features=marshal+unmarshal+size+pool+clone \
			--go-vtproto_opt=pool=github.com/zikunw/grpc-goproc-experiment/message.Batch \
			--go-vtproto_opt=pool=github.com/zikunw/grpc-goproc-experiment/message.KV \
			$$PROTO; \
	done

build:
	go build -o ./bin/reciever reciever/main.go
	go build -o ./bin/sender sender/*.go