.PHONY: run build test migrate-up migrate-down migrate-create sqlc-generate help

# Variáveis
MIGRATE_CMD = migrate
MIGRATE_PATH = file://db/migrations
# DB_URL padrão (pode ser sobrescrita via variável de ambiente ou linha de comando)
DB_URL ?= postgres://postgres:postgres@localhost:5432/gastrogo?sslmode=disable

# Executar a aplicação
run:
	@if [ -f .env ]; then export $$(cat .env | grep -v '^#' | xargs); fi; \
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
	@if [ -f .env ]; then export $$(cat .env | grep -v '^#' | xargs); fi; \
	$(MIGRATE_CMD) -path db/migrations -database "$${DATABASE_URL:-$(DB_URL)}" up

# Reverter migrations
migrate-down:
	@if [ -f .env ]; then export $$(cat .env | grep -v '^#' | xargs); fi; \
	$(MIGRATE_CMD) -path db/migrations -database "$${DATABASE_URL:-$(DB_URL)}" down

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

# Setup do banco de dados
setup-db:
	@bash scripts/setup-db.sh

# Docker Compose commands
docker-up:
	docker-compose up -d postgres

docker-down:
	docker-compose down

docker-logs:
	docker-compose logs -f postgres

docker-migrate:
	docker-compose --profile migrate up migrate

docker-reset:
	docker-compose down -v
	docker-compose up -d postgres
	@echo "Aguardando PostgreSQL iniciar..."
	@sleep 5
	@make migrate-up

# Ajuda
help:
	@echo "Comandos disponíveis:"
	@echo "  make run              - Executa a aplicação"
	@echo "  make build            - Compila a aplicação"
	@echo "  make test             - Executa os testes"
	@echo "  make test-coverage    - Executa testes com coverage"
	@echo "  make setup-db         - Setup inicial do banco de dados (cria DB e aplica migrations)"
	@echo "  make migrate-up      - Aplica migrations (use DB_URL para especificar conexão)"
	@echo "  make migrate-down    - Reverte migrations"
	@echo "  make migrate-create  - Cria nova migration"
	@echo "  make sqlc-generate   - Gera código SQLC"
	@echo ""
	@echo "Docker Compose:"
	@echo "  make docker-up        - Inicia PostgreSQL via Docker"
	@echo "  make docker-down     - Para e remove containers"
	@echo "  make docker-logs     - Mostra logs do PostgreSQL"
	@echo "  make docker-migrate   - Executa migrations via Docker"
	@echo "  make docker-reset    - Reseta banco (remove volumes e recria)"
	@echo ""
	@echo "Outros:"
	@echo "  make deps             - Instala/atualiza dependências"
	@echo "  make fmt              - Formata o código"
	@echo "  make lint             - Executa linter (requer golangci-lint)"

