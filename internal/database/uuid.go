package database

import "github.com/google/uuid"

// UUID é um alias para o tipo UUID do pacote google/uuid
// Este arquivo garante que a dependência seja mantida no go.mod
// mesmo antes do SQLC gerar código que use UUIDs
type UUID = uuid.UUID

