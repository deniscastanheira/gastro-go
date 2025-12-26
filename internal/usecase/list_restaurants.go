package usecase

import (
	"context"
	"fmt"
	"time"

	"gastro-go/internal/domain"
)

// RestaurantLister define a interface mínima necessária para listar restaurantes
// Segue Interface Segregation Principle: apenas o método que este use case precisa
type RestaurantLister interface {
	List(ctx context.Context, limit, offset int32) ([]*domain.Restaurant, error)
}

// ListRestaurantsUseCase implementa o caso de uso de listagem de restaurantes
type ListRestaurantsUseCase struct {
	repo RestaurantLister
}

// NewListRestaurantsUseCase cria uma nova instância do use case
func NewListRestaurantsUseCase(repo RestaurantLister) *ListRestaurantsUseCase {
	return &ListRestaurantsUseCase{
		repo: repo,
	}
}

// ListRestaurantsInput representa os dados de entrada para listar restaurantes
type ListRestaurantsInput struct {
	Limit  int32
	Offset int32
}

// Execute executa o caso de uso de listagem de restaurantes
func (uc *ListRestaurantsUseCase) Execute(ctx context.Context, input ListRestaurantsInput) ([]*domain.Restaurant, error) {
	if input.Limit <= 0 {
		input.Limit = 20 // Default
	}
	if input.Offset < 0 {
		input.Offset = 0
	}

	restaurants, err := uc.repo.List(ctx, input.Limit, input.Offset)
	if err != nil {
		return nil, fmt.Errorf("list restaurants usecase: %w", err)
	}

	// Calcular IsOpen para cada restaurante
	now := time.Now()
	for _, restaurant := range restaurants {
		restaurant.IsOpen = restaurant.CalculateIsOpen(now)
	}

	return restaurants, nil
}

