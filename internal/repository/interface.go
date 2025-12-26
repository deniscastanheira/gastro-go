package repository

import (
	"context"

	"github.com/google/uuid"

	"gastro-go/internal/domain"
)

// RestaurantRepositoryInterface define a interface do repository de restaurantes
type RestaurantRepositoryInterface interface {
	Create(ctx context.Context, restaurant *domain.Restaurant) error
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Restaurant, error)
	GetBySlug(ctx context.Context, slug string) (*domain.Restaurant, error)
	SlugExists(ctx context.Context, slug string) (bool, error)
	List(ctx context.Context, limit, offset int32) ([]*domain.Restaurant, error)
	UpdateStatus(ctx context.Context, id uuid.UUID, status string) error
	CreateAddress(ctx context.Context, address *domain.Address) error
	UpdateAddress(ctx context.Context, address *domain.Address) error
	GetAddress(ctx context.Context, restaurantID uuid.UUID) (*domain.Address, error)
	CreateOpeningHour(ctx context.Context, hour *domain.OpeningHour) error
	DeleteOpeningHoursByRestaurant(ctx context.Context, restaurantID uuid.UUID) error
	GetOpeningHours(ctx context.Context, restaurantID uuid.UUID) ([]*domain.OpeningHour, error)
	CreatePaymentMethod(ctx context.Context, method *domain.PaymentMethod) error
	DeletePaymentMethodsByRestaurant(ctx context.Context, restaurantID uuid.UUID) error
	GetPaymentMethods(ctx context.Context, restaurantID uuid.UUID) ([]*domain.PaymentMethod, error)
}

