package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"gastro-go/internal/domain"
)

// RestaurantCloser define a interface mínima necessária para fechar restaurantes
// Segue Interface Segregation Principle: apenas os métodos que este use case precisa
type RestaurantCloser interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Restaurant, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
}

// CloseRestaurantUseCase implementa o caso de uso de fechar um restaurante
type CloseRestaurantUseCase struct {
	repo RestaurantCloser
}

// NewCloseRestaurantUseCase cria uma nova instância do use case
func NewCloseRestaurantUseCase(repo RestaurantCloser) *CloseRestaurantUseCase {
	return &CloseRestaurantUseCase{
		repo: repo,
	}
}

// Execute executa o caso de uso de fechar restaurante
func (uc *CloseRestaurantUseCase) Execute(ctx context.Context, id uuid.UUID) error {
	// Verificar se o restaurante existe
	_, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("close restaurant usecase: %w", err)
	}

	// Atualizar status
	if err := uc.repo.UpdateStatus(ctx, id, domain.StatusClosed); err != nil {
		return fmt.Errorf("close restaurant usecase: %w", err)
	}

	return nil
}

