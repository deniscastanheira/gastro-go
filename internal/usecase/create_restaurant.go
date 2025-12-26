package usecase

import (
	"context"
	"errors"
	"fmt"

	"github.com/google/uuid"

	"gastro-go/internal/domain"
	"gastro-go/internal/utils"
)

// RestaurantCreator define a interface mínima necessária para criar restaurantes
// Segue Interface Segregation Principle: apenas os métodos que este use case precisa
type RestaurantCreator interface {
	SlugExists(ctx context.Context, slug string) (bool, error)
	Create(ctx context.Context, restaurant *domain.Restaurant) error
}

// CreateRestaurantUseCase implementa o caso de uso de criação de restaurante
type CreateRestaurantUseCase struct {
	repo RestaurantCreator
}

// NewCreateRestaurantUseCase cria uma nova instância do use case
func NewCreateRestaurantUseCase(repo RestaurantCreator) *CreateRestaurantUseCase {
	return &CreateRestaurantUseCase{
		repo: repo,
	}
}

// CreateRestaurantInput representa os dados de entrada para criar um restaurante
type CreateRestaurantInput struct {
	Name               string
	Slug               string // Opcional, será gerado se vazio
	Description        string
	Category           string
	DeliveryFee        int64
	MinOrderValue      int64
	PreparationTimeMin int
	SupportsPickup     bool
	SupportsDelivery   bool
	LogoURL            string
	BannerURL          string
	Address            *CreateAddressInput
}

// CreateAddressInput representa os dados de entrada para criar um endereço
type CreateAddressInput struct {
	Street     string
	Number     string
	Complement string
	City       string
	State      string
	ZipCode    string
	Lat        float64
	Lng        float64
}

// Execute executa o caso de uso de criação de restaurante
func (uc *CreateRestaurantUseCase) Execute(ctx context.Context, input CreateRestaurantInput) (*domain.Restaurant, error) {
	// Validações
	if input.Name == "" {
		return nil, fmt.Errorf("create restaurant usecase: name is required")
	}

	if input.DeliveryFee < 0 {
		return nil, fmt.Errorf("create restaurant usecase: delivery fee cannot be negative")
	}

	if input.MinOrderValue < 0 {
		return nil, fmt.Errorf("create restaurant usecase: min order value cannot be negative")
	}

	// Gerar slug se não fornecido
	slug := input.Slug
	if slug == "" {
		slug = utils.GenerateSlug(input.Name)
	}

	// Verificar se o slug já existe
	exists, err := uc.repo.SlugExists(ctx, slug)
	if err != nil {
		return nil, fmt.Errorf("create restaurant usecase: check slug uniqueness: %w", err)
	}
	if exists {
		return nil, fmt.Errorf("create restaurant usecase: slug already exists: %w", errors.New("conflict"))
	}

	// Criar restaurante
	restaurant := &domain.Restaurant{
		ID:                 uuid.New(),
		Name:               input.Name,
		Slug:               slug,
		Description:        input.Description,
		Status:             domain.StatusDraft,
		Category:           input.Category,
		Rating:             0,
		TotalReviews:       0,
		DeliveryFee:        input.DeliveryFee,
		MinOrderValue:      input.MinOrderValue,
		PreparationTimeMin: input.PreparationTimeMin,
		SupportsPickup:     input.SupportsPickup,
		SupportsDelivery:   input.SupportsDelivery,
		LogoURL:            input.LogoURL,
		BannerURL:          input.BannerURL,
	}

	// Criar endereço se fornecido
	if input.Address != nil {
		restaurant.Address = &domain.Address{
			ID:           uuid.New(),
			RestaurantID: restaurant.ID,
			Street:       input.Address.Street,
			Number:       input.Address.Number,
			Complement:   input.Address.Complement,
			City:         input.Address.City,
			State:        input.Address.State,
			ZipCode:      input.Address.ZipCode,
			Lat:          input.Address.Lat,
			Lng:          input.Address.Lng,
		}
	}

	// Salvar no banco
	if err := uc.repo.Create(ctx, restaurant); err != nil {
		return nil, fmt.Errorf("create restaurant usecase: %w", err)
	}

	return restaurant, nil
}

