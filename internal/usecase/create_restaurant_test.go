package usecase

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gastro-go/internal/domain"
)

// MockRestaurantCreator é um mock específico para RestaurantCreator
// Implementa apenas os métodos necessários para CreateRestaurantUseCase
// Segue Interface Segregation Principle: mock focado e simples
type MockRestaurantCreator struct {
	mock.Mock
}

func (m *MockRestaurantCreator) SlugExists(ctx context.Context, slug string) (bool, error) {
	args := m.Called(ctx, slug)
	return args.Bool(0), args.Error(1)
}

func (m *MockRestaurantCreator) Create(ctx context.Context, restaurant *domain.Restaurant) error {
	args := m.Called(ctx, restaurant)
	return args.Error(0)
}

func TestCreateRestaurantUseCase_Execute_Success(t *testing.T) {
	// Input
	ctx := context.Background()
	input := CreateRestaurantInput{
		Name:               "Pizza do João",
		Description:        "Melhor pizza da cidade",
		Category:           "Pizza",
		DeliveryFee:        500,  // R$ 5,00
		MinOrderValue:      2000, // R$ 20,00
		PreparationTimeMin: 30,
		SupportsPickup:     true,
		SupportsDelivery:   true,
	}

	// Mock
	mockRepo := new(MockRestaurantCreator)
	mockRepo.On("SlugExists", ctx, mock.AnythingOfType("string")).Return(false, nil)
	mockRepo.On("Create", ctx, mock.AnythingOfType("*domain.Restaurant")).Return(nil)

	// Execute
	uc := NewCreateRestaurantUseCase(mockRepo)
	restaurant, err := uc.Execute(ctx, input)

	// Assert
	assert.NoError(t, err)
	assert.NotNil(t, restaurant)
	assert.Equal(t, "Pizza do João", restaurant.Name)
	assert.Equal(t, domain.StatusDraft, restaurant.Status)
	assert.Equal(t, int64(500), restaurant.DeliveryFee)
	assert.Equal(t, int64(2000), restaurant.MinOrderValue)
	assert.NotEmpty(t, restaurant.Slug)
	mockRepo.AssertExpectations(t)
}

func TestCreateRestaurantUseCase_Execute_SlugConflict(t *testing.T) {
	// Input
	ctx := context.Background()
	input := CreateRestaurantInput{
		Name: "Pizza do João",
	}

	// Mock
	mockRepo := new(MockRestaurantCreator)
	mockRepo.On("SlugExists", ctx, mock.AnythingOfType("string")).Return(true, nil)

	// Execute
	uc := NewCreateRestaurantUseCase(mockRepo)
	restaurant, err := uc.Execute(ctx, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, restaurant)
	assert.Contains(t, err.Error(), "slug already exists")
	mockRepo.AssertExpectations(t)
}

func TestCreateRestaurantUseCase_Execute_ValidationError(t *testing.T) {
	// Input
	ctx := context.Background()
	input := CreateRestaurantInput{
		Name:          "", // Nome vazio
		DeliveryFee:   -100,
		MinOrderValue: -50,
	}

	// Mock
	mockRepo := new(MockRestaurantCreator)

	// Execute
	uc := NewCreateRestaurantUseCase(mockRepo)
	restaurant, err := uc.Execute(ctx, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, restaurant)
	assert.Contains(t, err.Error(), "name is required")
	mockRepo.AssertNotCalled(t, "Create")
}

