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

## Desenvolvimento

### Pré-requisitos

- Go 1.21+
- PostgreSQL
- golang-migrate (para migrations)
- sqlc (para geração de código type-safe)

### Executar

```bash
go run cmd/api/main.go
```

O servidor iniciará na porta 8080 (ou na porta definida pela variável de ambiente `PORT`).

## Migrations

As migrations estão em `db/migrations/` e seguem o padrão:
- `00000N_description.up.sql` - Aplica a mudança
- `00000N_description.down.sql` - Reverte a mudança

**IMPORTANTE:** Nunca modifique uma migration já aplicada. Sempre crie uma nova migration.

