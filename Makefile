.PHONY: run build test migrate-up migrate-down migrate-create sqlc-generate help

# Variáveis
MIGRATE_CMD = migrate
MIGRATE_PATH = file://db/migrations
DB_URL ?= postgres://localhost/gastrogodb?sslmode=disable

# Executar a aplicação
run:
	go run cmd/api/main.go

# Build da aplicação
build:
	go build -o bin/api cmd/api/main.go

# Executar testes
test:
	go test ./...

# Executar testes com coverage
test-coverage:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

# Aplicar migrations
migrate-up:
	$(MIGRATE_CMD) -path db/migrations -database "$(DB_URL)" up

# Reverter migrations
migrate-down:
	$(MIGRATE_CMD) -path db/migrations -database "$(DB_URL)" down

# Criar nova migration
migrate-create:
	@read -p "Nome da migration: " name; \
	$(MIGRATE_CMD) create -ext sql -dir db/migrations -seq $$name

# Gerar código SQLC
sqlc-generate:
	sqlc generate

# Instalar dependências
deps:
	go mod download
	go mod tidy

# Formatar código
fmt:
	go fmt ./...

# Executar linter
lint:
	golangci-lint run

# Ajuda
help:
	@echo "Comandos disponíveis:"
	@echo "  make run              - Executa a aplicação"
	@echo "  make build            - Compila a aplicação"
	@echo "  make test             - Executa os testes"
	@echo "  make test-coverage    - Executa testes com coverage"
	@echo "  make migrate-up       - Aplica migrations (use DB_URL para especificar conexão)"
	@echo "  make migrate-down     - Reverte migrations"
	@echo "  make migrate-create   - Cria nova migration"
	@echo "  make sqlc-generate    - Gera código SQLC"
	@echo "  make deps             - Instala/atualiza dependências"
	@echo "  make fmt              - Formata o código"
	@echo "  make lint             - Executa linter (requer golangci-lint)"

