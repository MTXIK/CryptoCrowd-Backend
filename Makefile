OUTPUT := CryptoCrowd

MAIN := cmd/main.go

GOPATH := $(shell go env GOPATH)
GOBIN := $(GOPATH)/bin

.PHONY: fmt
fmt:
	@echo "Форматирование кода..."
	@go fmt ./...

.PHONY: build
build:
	@echo "Сборка проекта..."
	@go build -o $(OUTPUT) $(MAIN)

# Применение миграций
.PHONY: migrate-up
migrate-up:
	@echo "Применение миграций..."
	@go run cmd/migrate/main.go up

# Откат миграций
.PHONY: migrate-down
migrate-down:
	@echo "Откат миграций..."
	@go run cmd/migrate/main.go down

# Статус миграций
.PHONY: migrate-status
migrate-status:
	@echo "Статус миграций..."
	@go run cmd/migrate/main.go status

# Запуск базы данных
.PHONY: compose-up
compose-up:
	@echo "Запуск PostgreSQL..."
	@docker compose up -d

# Остановка базы данных
.PHONY: compose-down
compose-down:
	@echo "Остановка PostgreSQL..."
	@docker compose down

# Удаление томов базы данных
.PHONY: compose-down-vol
compose-down-vol:
	@echo "Остановка PostgreSQL..."
	@docker compose down -v

# Цель запуска: сначала сборка, затем запуск исполняемого файла
.PHONY: run
run: build
	@echo "Запуск проекта..."
	@./$(OUTPUT)

# Полный запуск (БД + миграции + приложение)
.PHONY: start
start: compose-up
	@echo "Ожидание готовности БД..."
	@sleep 2
	@$(MAKE) migrate-up
	@$(MAKE) run
