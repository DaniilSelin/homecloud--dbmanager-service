APP_NAME=dbmanager
BINARY=bin/$(APP_NAME)
SRC=./cmd/server/main.go
LOG=server.log

.PHONY: all build run stop restart logs proto proto-clean test clean migrate migrate-up migrate-down migrate-status migrate-reset help

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
	@echo "[PROTO BUILD]"
	./proto.sh build

proto-clean:
	@echo "[PROTO CLEAN]"
	./proto.sh clean

test:
	@echo "[TEST]"
	go test -v ./test

# Команды для работы с миграциями
migrate: migrate-up

migrate-up:
	@echo "[MIGRATE UP]"
	./migrations/migrate.sh up

migrate-down:
	@echo "[MIGRATE DOWN]"
	./migrations/migrate.sh down

migrate-status:
	@echo "[MIGRATE STATUS]"
	./migrations/migrate.sh status

migrate-reset:
	@echo "[MIGRATE RESET]"
	./migrations/migrate.sh reset

help:
	@echo "Доступные команды:"
	@echo ""
	@echo "Основные команды:"
	@echo "  build        - Собрать приложение"
	@echo "  run          - Запустить приложение"
	@echo "  stop         - Остановить приложение"
	@echo "  restart      - Перезапустить приложение"
	@echo "  logs         - Показать логи"
	@echo "  proto        - Сгенерировать protobuf файлы"
	@echo "  proto-clean  - Удалить сгенерированные protobuf файлы"
	@echo "  test         - Запустить тесты"
	@echo "  clean        - Очистить временные файлы"
	@echo ""
	@echo "Команды миграций:"
	@echo "  migrate-up     - Применить все миграции"
	@echo "  migrate-down   - Откатить последнюю миграцию"
	@echo "  migrate-status - Показать статус миграций"
	@echo "  migrate-reset  - Откатить все миграции"
	@echo "  migrate        - Алиас для migrate-up"
	@echo ""
	@echo "Для получения справки по миграциям: ./migrations/migrate.sh help"

clean:
	@echo "[CLEAN]"
	rm -f $(BINARY) $(APP_NAME).pid $(LOG) 