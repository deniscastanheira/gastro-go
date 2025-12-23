package domain

// Example é um exemplo de entidade de domínio
// As entidades de domínio são puras Go structs sem dependências de framework ou banco
type Example struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"` // ISO 8601 format, UTC
}

