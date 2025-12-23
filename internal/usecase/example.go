package usecase

import (
	"context"
	"fmt"

	"gastro-go/internal/domain"
)

// ExampleUseCase demonstra a estrutura de um Use Case
// Use Cases contêm a lógica de negócio
type ExampleUseCase struct {
	// Dependências: Repositories serão injetados aqui
}

// NewExampleUseCase cria uma nova instância do Use Case
func NewExampleUseCase() *ExampleUseCase {
	return &ExampleUseCase{}
}

// Execute executa a lógica de negócio
// Context sempre é o primeiro parâmetro
func (uc *ExampleUseCase) Execute(ctx context.Context, input ExampleInput) (*domain.Example, error) {
	// Lógica de negócio aqui
	// Chamadas para repositories
	// Validações
	// Transformações

	if input.Name == "" {
		return nil, fmt.Errorf("example usecase: name is required")
	}

	return &domain.Example{
		ID:   1,
		Name: input.Name,
	}, nil
}

// ExampleInput demonstra um struct de input específico para o Use Case
type ExampleInput struct {
	Name string `json:"name"`
}

