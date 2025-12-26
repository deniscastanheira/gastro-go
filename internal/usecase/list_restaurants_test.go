package usecase

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"gastro-go/internal/domain"
)

func TestListRestaurantsUseCase_Execute_Success(t *testing.T) {
	// Input
	ctx := context.Background()
	input := ListRestaurantsInput{
		Limit:  10,
		Offset: 0,
	}

	// Mock data
	restaurant1ID := uuid.New()
	restaurant1 := &domain.Restaurant{
		ID:                 restaurant1ID,
		Name:               "Pizza do João",
		Slug:               "pizza-do-joao",
		Status:             domain.StatusOpen,
		DeliveryFee:        500,
		MinOrderValue:      2000,
		PreparationTimeMin: 30,
		SupportsPickup:     true,
		SupportsDelivery:   true,
		OpeningHours: []domain.OpeningHour{
			{
				ID:           uuid.New(),
				RestaurantID: restaurant1ID,
				Weekday:      1, // Segunda-feira
				OpensAt:      480,  // 08:00
				ClosesAt:     1200, // 20:00
			},
		},
	}

	restaurant2 := &domain.Restaurant{
		ID:                 uuid.New(),
		Name:               "Burgers King",
		Slug:               "burgers-king",
		Status:             domain.StatusDraft,
		DeliveryFee:        300,
		MinOrderValue:      1500,
		PreparationTimeMin: 20,
		SupportsPickup:     true,
		SupportsDelivery:   true,
	}

	// Mock
	mockRepo := new(MockRestaurantRepository)
	mockRepo.On("List", ctx, int32(10), int32(0)).Return([]*domain.Restaurant{restaurant1, restaurant2}, nil)

	// Execute
	uc := NewListRestaurantsUseCase(mockRepo)
	restaurants, err := uc.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, restaurants)
	assert.Len(t, restaurants, 2)
	assert.Equal(t, "Pizza do João", restaurants[0].Name)
	assert.Equal(t, "Burgers King", restaurants[1].Name)
	// IsOpen deve ser calculado
	assert.NotNil(t, restaurants[0].IsOpen)
	assert.NotNil(t, restaurants[1].IsOpen)
	mockRepo.AssertExpectations(t)
}

func TestListRestaurantsUseCase_Execute_DefaultValues(t *testing.T) {
	// Input
	ctx := context.Background()
	input := ListRestaurantsInput{
		Limit:  0,  // Deve usar default
		Offset: -1, // Deve ser ajustado para 0
	}

	// Mock
	mockRepo := new(MockRestaurantRepository)
	mockRepo.On("List", ctx, int32(20), int32(0)).Return([]*domain.Restaurant{}, nil)

	// Execute
	uc := NewListRestaurantsUseCase(mockRepo)
	restaurants, err := uc.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, restaurants)
	mockRepo.AssertExpectations(t)
}

