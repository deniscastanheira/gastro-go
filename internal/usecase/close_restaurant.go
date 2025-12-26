package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"gastro-go/internal/domain"
	"gastro-go/internal/repository"
)

// CloseRestaurantUseCase implementa o caso de uso de fechar um restaurante
type CloseRestaurantUseCase struct {
	repo repository.RestaurantRepositoryInterface
}

// NewCloseRestaurantUseCase cria uma nova instância do use case
func NewCloseRestaurantUseCase(repo repository.RestaurantRepositoryInterface) *CloseRestaurantUseCase {
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

