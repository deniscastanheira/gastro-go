package repository

import (
	"context"
	"fmt"

	"gastro-go/internal/domain"
)

// ExampleRepository demonstra a estrutura de um Repository
// Repositories implementam interfaces do Domain usando código gerado pelo SQLC
type ExampleRepository struct {
	// Dependências: código gerado pelo SQLC será injetado aqui
	// Exemplo: db *database.Queries
}

// NewExampleRepository cria uma nova instância do Repository
func NewExampleRepository() *ExampleRepository {
	return &ExampleRepository{}
}

// GetByID busca um exemplo por ID
func (r *ExampleRepository) GetByID(ctx context.Context, id int64) (*domain.Example, error) {
	// Aqui chamaria o código gerado pelo SQLC
	// Exemplo:
	// dbExample, err := r.db.GetExampleByID(ctx, id)
	// if err != nil {
	//     return nil, fmt.Errorf("example repository: %w", err)
	// }
	// return toDomain(dbExample), nil

	return nil, fmt.Errorf("not implemented")
}

// toDomain converte uma entidade do banco para uma entidade de domínio
// func toDomain(db *database.Example) *domain.Example {
//     return &domain.Example{
//         ID:   db.ID,
//         Name: db.Name,
//     }
// }

