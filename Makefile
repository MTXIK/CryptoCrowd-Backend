OUTPUT := CryptoCrowd

MAIN := cmd/main.go

GOPATH := $(shell go env GOPATH)
GOBIN := $(GOPATH)/bin

.PHONY: fmt
fmt:
	@echo "Форматирование кода..."
	@go fmt ./...

.PHONY: build
build: fmt #lint отключено
	@echo "Сборка проекта..."
	@go build -o $(OUTPUT) $(MAIN)

.PHONY: migrate-up
# Применение миграций
migrate-up:
	@echo "Применение миграций..."
	@go run cmd/migrate/main.go up

.PHONY: migrate-down
# Откат миграций
migrate-down:
	@echo "Откат миграций..."
	@go run cmd/migrate/main.go down

.PHONY: migrate-status
# Статус миграций
migrate-status:
	@echo "Статус миграций..."
	@go run cmd/migrate/main.go status