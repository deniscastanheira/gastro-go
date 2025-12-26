package usecase

import (
	"context"
	"fmt"
	"time"

	"gastro-go/internal/domain"
)

// RestaurantGetterBySlug define a interface mínima necessária para buscar restaurante por slug
// Segue Interface Segregation Principle: apenas o método que este use case precisa
type RestaurantGetterBySlug interface {
	GetBySlug(ctx context.Context, slug string) (*domain.Restaurant, error)
}

// GetRestaurantBySlugUseCase implementa o caso de uso de buscar restaurante por slug
type GetRestaurantBySlugUseCase struct {
	repo RestaurantGetterBySlug
}

// NewGetRestaurantBySlugUseCase cria uma nova instância do use case
func NewGetRestaurantBySlugUseCase(repo RestaurantGetterBySlug) *GetRestaurantBySlugUseCase {
	return &GetRestaurantBySlugUseCase{
		repo: repo,
	}
}

// Execute executa o caso de uso de buscar restaurante por slug
func (uc *GetRestaurantBySlugUseCase) Execute(ctx context.Context, slug string) (*domain.Restaurant, error) {
	restaurant, err := uc.repo.GetBySlug(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("get restaurant by slug usecase: %w", err)
	}

	// Calcular IsOpen
	now := time.Now()
	restaurant.IsOpen = restaurant.CalculateIsOpen(now)

	return restaurant, nil
}

