package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"gastro-go/internal/domain"
)

// PaymentMethodsUpdater define a interface mínima necessária para atualizar métodos de pagamento
// Segue Interface Segregation Principle: apenas os métodos que este use case precisa
type PaymentMethodsUpdater interface {
	GetByID(ctx context.Context, id uuid.UUID) (*domain.Restaurant, error)
	DeletePaymentMethodsByRestaurant(ctx context.Context, restaurantID uuid.UUID) error
	CreatePaymentMethod(ctx context.Context, method *domain.PaymentMethod) error
}

// UpdatePaymentMethodsUseCase implementa o caso de uso de atualizar métodos de pagamento
type UpdatePaymentMethodsUseCase struct {
	repo PaymentMethodsUpdater
}

// NewUpdatePaymentMethodsUseCase cria uma nova instância do use case
func NewUpdatePaymentMethodsUseCase(repo PaymentMethodsUpdater) *UpdatePaymentMethodsUseCase {
	return &UpdatePaymentMethodsUseCase{
		repo: repo,
	}
}

// UpdatePaymentMethodsInput representa os dados de entrada para atualizar métodos de pagamento
type UpdatePaymentMethodsInput struct {
	RestaurantID uuid.UUID
	Methods      []string // "PIX", "CREDIT_CARD", "DEBIT_CARD"
}

// Execute executa o caso de uso de atualizar métodos de pagamento
func (uc *UpdatePaymentMethodsUseCase) Execute(ctx context.Context, input UpdatePaymentMethodsInput) error {
	// Verificar se o restaurante existe
	_, err := uc.repo.GetByID(ctx, input.RestaurantID)
	if err != nil {
		return fmt.Errorf("update payment methods usecase: %w", err)
	}

	// Validar métodos
	validMethods := map[string]bool{
		domain.PaymentMethodPIX:        true,
		domain.PaymentMethodCreditCard: true,
		domain.PaymentMethodDebitCard:  true,
	}

	for _, method := range input.Methods {
		if !validMethods[method] {
			return fmt.Errorf("update payment methods usecase: invalid payment method: %s", method)
		}
	}

	// Deletar métodos existentes
	if err := uc.repo.DeletePaymentMethodsByRestaurant(ctx, input.RestaurantID); err != nil {
		return fmt.Errorf("update payment methods usecase: delete existing methods: %w", err)
	}

	// Criar novos métodos
	for _, methodStr := range input.Methods {
		method := &domain.PaymentMethod{
			ID:           uuid.New(),
			RestaurantID: input.RestaurantID,
			Method:       methodStr,
		}

		if err := uc.repo.CreatePaymentMethod(ctx, method); err != nil {
			return fmt.Errorf("update payment methods usecase: create method: %w", err)
		}
	}

	return nil
}

