package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"gastro-go/internal/domain"
	"gastro-go/internal/repository"
)

// OpenRestaurantUseCase implementa o caso de uso de abrir um restaurante
type OpenRestaurantUseCase struct {
	repo repository.RestaurantRepositoryInterface
}

// NewOpenRestaurantUseCase cria uma nova instância do use case
func NewOpenRestaurantUseCase(repo repository.RestaurantRepositoryInterface) *OpenRestaurantUseCase {
	return &OpenRestaurantUseCase{
		repo: repo,
	}
}

// Execute executa o caso de uso de abrir restaurante
func (uc *OpenRestaurantUseCase) Execute(ctx context.Context, id uuid.UUID) error {
	// Buscar restaurante
	restaurant, err := uc.repo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("open restaurant usecase: %w", err)
	}

	// Validações de regras de negócio
	if restaurant.Address == nil {
		return fmt.Errorf("open restaurant usecase: restaurant must have an address to be opened")
	}

	openingHours, err := uc.repo.GetOpeningHours(ctx, id)
	if err != nil {
		return fmt.Errorf("open restaurant usecase: get opening hours: %w", err)
	}
	if len(openingHours) == 0 {
		return fmt.Errorf("open restaurant usecase: restaurant must have opening hours to be opened")
	}

	paymentMethods, err := uc.repo.GetPaymentMethods(ctx, id)
	if err != nil {
		return fmt.Errorf("open restaurant usecase: get payment methods: %w", err)
	}
	if len(paymentMethods) == 0 {
		return fmt.Errorf("open restaurant usecase: restaurant must have at least one payment method to be opened")
	}

	// Atualizar status
	if err := uc.repo.UpdateStatus(ctx, id, domain.StatusOpen); err != nil {
		return fmt.Errorf("open restaurant usecase: %w", err)
	}

	return nil
}

