APP_NAME=dbmanager
BINARY=bin/$(APP_NAME)
SRC=./cmd/server/main.go
LOG=server.log

.PHONY: all build run stop restart logs proto test clean

all: build

build:
	@echo "[BUILD]"
	go build -o $(BINARY) $(SRC)

run:
	@echo "[RUN]"
	@nohup $(BINARY) > $(LOG) 2>&1 & echo $$! > $(APP_NAME).pid
	@sleep 1
	@echo "Started $(APP_NAME) with PID $$(cat $(APP_NAME).pid)"

stop:
	@echo "[STOP]"
	@if [ -f $(APP_NAME).pid ]; then \
		kill -9 $$(cat $(APP_NAME).pid) && rm -f $(APP_NAME).pid && echo "Stopped $(APP_NAME)"; \
	else \
		echo "No PID file found"; \
	fi

restart: stop build run

logs:
	@echo "[LOGS]"
	tail -f $(LOG)

proto:
	@echo "[PROTO]"
	cd internal/transport/grpc/protos && \
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative db_manager.proto

test:
	@echo "[TEST]"
	go test -v ./test

clean:
	@echo "[CLEAN]"
	rm -f $(BINARY) $(APP_NAME).pid $(LOG) 