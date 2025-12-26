package usecase

import (
	"context"
	"fmt"
	"time"

	"gastro-go/internal/domain"
	"gastro-go/internal/repository"
)

// GetRestaurantBySlugUseCase implementa o caso de uso de buscar restaurante por slug
type GetRestaurantBySlugUseCase struct {
	repo repository.RestaurantRepositoryInterface
}

// NewGetRestaurantBySlugUseCase cria uma nova instância do use case
func NewGetRestaurantBySlugUseCase(repo repository.RestaurantRepositoryInterface) *GetRestaurantBySlugUseCase {
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

