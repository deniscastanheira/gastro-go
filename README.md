# GastroGo

Backend para sistema de delivery de comida desenvolvido em Go.

## Tech Stack

- **Language:** Go (Latest stable version)
- **Framework:** Echo v4
- **Database:** PostgreSQL
- **Drivers:** pgx/v5
- **Libraries:** sqlc (type-safe SQL), testify (testing)

## Estrutura do Projeto

```
gastro-go/
├── cmd/
│   └── api/              # Entry point da aplicação
├── internal/
│   ├── domain/           # Entidades de negócio puras
│   ├── handler/          # Controllers HTTP (Echo handlers)
│   ├── usecase/          # Lógica de negócio (um struct por ação)
│   ├── repository/       # Camada de acesso a dados
│   └── database/         # Configuração do banco e código gerado pelo SQLC
└── db/
    └── migrations/       # Migrações SQL
```

## Arquitetura

O projeto segue uma arquitetura Clean Architecture leve com Use Cases:

- **Domain:** Entidades puras de Go (sem dependências de framework ou banco)
- **Handler:** Recebe requisições HTTP, parseia DTOs, chama Use Cases e retorna respostas
- **Use Case:** Contém a lógica de negócio (um arquivo/struct por ação)
- **Repository:** Implementa interfaces do Domain usando código gerado pelo SQLC

## Regras de Negócio

- **Dinheiro:** Sempre em `int64` (centavos)
- **Tempo:** Sempre UTC

## Quick Start (Docker Compose)

A forma mais rápida de começar:

```bash
# 1. Iniciar PostgreSQL
make docker-up

# 2. Configurar variáveis de ambiente
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/gastrogo?sslmode=disable"

# 3. Executar migrations
make migrate-up

# 4. Executar a aplicação
make run
```

Pronto! A API estará rodando em `http://localhost:8080`

## Desenvolvimento

### Pré-requisitos

- Go 1.21+
- PostgreSQL 12+
- golang-migrate (para migrations)
- sqlc (para geração de código type-safe)

### Instalação das Ferramentas

#### Instalar golang-migrate

**macOS:**
```bash
brew install golang-migrate
```

**Linux:**
```bash
curl -L https://github.com/golang-migrate/migrate/releases/download/v4.16.2/migrate.linux-amd64.tar.gz | tar xvz
sudo mv migrate /usr/local/bin/migrate
```

**Windows:**
```powershell
choco install golang-migrate
```

#### Instalar sqlc

```bash
go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
```

### Configuração do Banco de Dados

#### 1. Instalar PostgreSQL

**Opção A: Docker Compose (Recomendado - Mais Fácil)**

```bash
# Iniciar PostgreSQL
make docker-up

# Ou diretamente com docker-compose
docker-compose up -d postgres
```

Isso irá:
- Criar um container PostgreSQL na porta 5432
- Criar o banco `gastrogo` automaticamente
- Configurar usuário/senha padrão (postgres/postgres)
- Persistir dados em um volume Docker

**Opção B: Instalação Local**

**macOS:**
```bash
brew install postgresql@15
brew services start postgresql@15
```

**Linux (Ubuntu/Debian):**
```bash
sudo apt update
sudo apt install postgresql postgresql-contrib
sudo systemctl start postgresql
```

**Opção C: Docker Standalone**
```bash
docker run --name gastrogo-db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=postgres \
  -e POSTGRES_DB=gastrogo \
  -p 5432:5432 \
  -d postgres:15
```

#### 2. Criar o Banco de Dados

Conecte-se ao PostgreSQL e crie o banco:

```bash
# Conectar ao PostgreSQL
psql -U postgres

# Criar o banco de dados
CREATE DATABASE gastrogo;

# Sair do psql
\q
```

Ou via linha de comando:
```bash
createdb -U postgres gastrogo
```

#### 3. Configurar Variáveis de Ambiente

Crie um arquivo `.env` na raiz do projeto (ou configure as variáveis no seu ambiente):

**Se usar Docker Compose:**
```bash
# DATABASE_URL para Docker Compose
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/gastrogo?sslmode=disable"
export PORT=8080
```

**Se usar PostgreSQL local:**
```bash
# Opção 1: Usar DATABASE_URL (recomendado)
DATABASE_URL=postgres://postgres:postgres@localhost:5432/gastrogo?sslmode=disable

# Opção 2: Usar variáveis individuais
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=gastrogo

# Porta da API (opcional, padrão: 8080)
PORT=8080
```

**Nota:** O código usa `DATABASE_URL` se disponível, caso contrário usa as variáveis individuais.

#### 4. Executar Migrations

**Opção A: Docker Compose (Mais Fácil)**

```bash
# 1. Iniciar PostgreSQL
make docker-up

# 2. Executar migrations via Docker
make docker-migrate

# Ou tudo de uma vez (reset completo)
make docker-reset
```

**Opção B: Script Automatizado**

```bash
# Executa setup completo (cria DB e aplica migrations)
make setup-db

# Ou execute o script diretamente
bash scripts/setup-db.sh
```

**Opção C: Manual**

Aplique as migrations para criar as tabelas:

```bash
# Configurar DATABASE_URL primeiro
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/gastrogo?sslmode=disable"

# Usando o Makefile
make migrate-up

# Ou diretamente com migrate
migrate -path db/migrations -database "$DATABASE_URL" up
```

Para reverter migrations:
```bash
make migrate-down
```

### Comandos Docker Compose Úteis

```bash
# Iniciar PostgreSQL
make docker-up

# Ver logs do PostgreSQL
make docker-logs

# Parar PostgreSQL
make docker-down

# Executar migrations via Docker
make docker-migrate

# Reset completo (remove volumes e recria tudo)
make docker-reset
```

### Executar a Aplicação

```bash
# Carregar variáveis de ambiente (se usar .env)
export $(cat .env | xargs)

# Executar
make run
# ou
go run cmd/api/main.go
```

O servidor iniciará na porta 8080 (ou na porta definida pela variável de ambiente `PORT`).

### Verificar se está Funcionando

```bash
# Health check
curl http://localhost:8080/health

# Deve retornar:
# {"status":"ok","service":"gastro-go"}
```

## Migrations

As migrations estão em `db/migrations/` e seguem o padrão:
- `00000N_description.up.sql` - Aplica a mudança
- `00000N_description.down.sql` - Reverte a mudança

**IMPORTANTE:** Nunca modifique uma migration já aplicada. Sempre crie uma nova migration.

### Comandos Úteis

```bash
# Aplicar todas as migrations
make migrate-up

# Reverter última migration
make migrate-down

# Criar nova migration
make migrate-create

# Ver status das migrations
migrate -path db/migrations -database "$DB_URL" version
```

## Estrutura do Banco de Dados

Após executar as migrations, você terá as seguintes tabelas:

- `restaurants` - Dados principais dos restaurantes
- `restaurant_addresses` - Endereços dos restaurantes
- `restaurant_opening_hours` - Horários de funcionamento
- `restaurant_payment_methods` - Métodos de pagamento aceitos

Todas as tabelas têm índices apropriados e constraints de integridade referencial.

