package usecase

import (
	"context"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"

	"gastro-go/internal/domain"
)

// MockRestaurantRepository é um mock do RestaurantRepository para testes
type MockRestaurantRepository struct {
	mock.Mock
}

func (m *MockRestaurantRepository) Create(ctx context.Context, restaurant *domain.Restaurant) error {
	args := m.Called(ctx, restaurant)
	return args.Error(0)
}

func (m *MockRestaurantRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Restaurant, error) {
	args := m.Called(ctx, id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Restaurant), args.Error(1)
}

func (m *MockRestaurantRepository) GetBySlug(ctx context.Context, slug string) (*domain.Restaurant, error) {
	args := m.Called(ctx, slug)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Restaurant), args.Error(1)
}

func (m *MockRestaurantRepository) SlugExists(ctx context.Context, slug string) (bool, error) {
	args := m.Called(ctx, slug)
	return args.Bool(0), args.Error(1)
}

func (m *MockRestaurantRepository) List(ctx context.Context, limit, offset int32) ([]*domain.Restaurant, error) {
	args := m.Called(ctx, limit, offset)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.Restaurant), args.Error(1)
}

func (m *MockRestaurantRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	args := m.Called(ctx, id, status)
	return args.Error(0)
}

func (m *MockRestaurantRepository) CreateAddress(ctx context.Context, address *domain.Address) error {
	args := m.Called(ctx, address)
	return args.Error(0)
}

func (m *MockRestaurantRepository) UpdateAddress(ctx context.Context, address *domain.Address) error {
	args := m.Called(ctx, address)
	return args.Error(0)
}

func (m *MockRestaurantRepository) GetAddress(ctx context.Context, restaurantID uuid.UUID) (*domain.Address, error) {
	args := m.Called(ctx, restaurantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*domain.Address), args.Error(1)
}

func (m *MockRestaurantRepository) CreateOpeningHour(ctx context.Context, hour *domain.OpeningHour) error {
	args := m.Called(ctx, hour)
	return args.Error(0)
}

func (m *MockRestaurantRepository) DeleteOpeningHoursByRestaurant(ctx context.Context, restaurantID uuid.UUID) error {
	args := m.Called(ctx, restaurantID)
	return args.Error(0)
}

func (m *MockRestaurantRepository) GetOpeningHours(ctx context.Context, restaurantID uuid.UUID) ([]*domain.OpeningHour, error) {
	args := m.Called(ctx, restaurantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.OpeningHour), args.Error(1)
}

func (m *MockRestaurantRepository) CreatePaymentMethod(ctx context.Context, method *domain.PaymentMethod) error {
	args := m.Called(ctx, method)
	return args.Error(0)
}

func (m *MockRestaurantRepository) DeletePaymentMethodsByRestaurant(ctx context.Context, restaurantID uuid.UUID) error {
	args := m.Called(ctx, restaurantID)
	return args.Error(0)
}

func (m *MockRestaurantRepository) GetPaymentMethods(ctx context.Context, restaurantID uuid.UUID) ([]*domain.PaymentMethod, error) {
	args := m.Called(ctx, restaurantID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).([]*domain.PaymentMethod), args.Error(1)
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
	mockRepo := new(MockRestaurantRepository)
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
	mockRepo := new(MockRestaurantRepository)
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
	mockRepo := new(MockRestaurantRepository)

	// Execute
	uc := NewCreateRestaurantUseCase(mockRepo)
	restaurant, err := uc.Execute(ctx, input)

	// Assert
	assert.Error(t, err)
	assert.Nil(t, restaurant)
	assert.Contains(t, err.Error(), "name is required")
	mockRepo.AssertNotCalled(t, "Create")
}

