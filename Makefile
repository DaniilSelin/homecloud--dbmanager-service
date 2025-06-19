PROTO_DIR=internal/transport/grpc/protos
PROTO_FILE=$(PROTO_DIR)/db_manager.proto
GO_OUT=$(PROTO_DIR)

.PHONY: proto build run clean

proto:
	protoc -I=$(PROTO_DIR) --go_out=$(GO_OUT) --go_opt=paths=source_relative --go-grpc_out=$(GO_OUT) --go-grpc_opt=paths=source_relative $(PROTO_FILE)

build:
	go build -o bin/dbmanager cmd/server/main.go

run: build
	./bin/dbmanager

clean:
	rm -rf bin/
	find $(GO_OUT) -name '*.pb.go' -delete 