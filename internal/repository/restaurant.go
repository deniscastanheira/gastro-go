package repository

import (
	"context"
	"fmt"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"

	"gastro-go/internal/database"
	"gastro-go/internal/domain"
)

// RestaurantRepository implementa operações de acesso a dados para restaurantes
type RestaurantRepository struct {
	queries *database.Queries
}

// NewRestaurantRepository cria uma nova instância do repository
func NewRestaurantRepository(queries *database.Queries) *RestaurantRepository {
	return &RestaurantRepository{
		queries: queries,
	}
}

// Create cria um novo restaurante
func (r *RestaurantRepository) Create(ctx context.Context, restaurant *domain.Restaurant) error {
	// Converter para modelo do banco
	params := database.CreateRestaurantParams{
		Name:               restaurant.Name,
		Slug:               restaurant.Slug,
		Status:             restaurant.Status,
		Rating:             int32(restaurant.Rating),
		TotalReviews:       int32(restaurant.TotalReviews),
		DeliveryFee:        restaurant.DeliveryFee,
		MinOrderValue:      restaurant.MinOrderValue,
		PreparationTimeMin: int32(restaurant.PreparationTimeMin),
		SupportsPickup:     restaurant.SupportsPickup,
		SupportsDelivery:   restaurant.SupportsDelivery,
	}

	if restaurant.Description != "" {
		params.Description = pgtype.Text{String: restaurant.Description, Valid: true}
	}
	if restaurant.Category != "" {
		params.Category = pgtype.Text{String: restaurant.Category, Valid: true}
	}
	if restaurant.LogoURL != "" {
		params.LogoUrl = pgtype.Text{String: restaurant.LogoURL, Valid: true}
	}
	if restaurant.BannerURL != "" {
		params.BannerUrl = pgtype.Text{String: restaurant.BannerURL, Valid: true}
	}

	dbRestaurant, err := r.queries.CreateRestaurant(ctx, params)
	if err != nil {
		return fmt.Errorf("restaurant repository: create restaurant: %w", err)
	}

	// Atualizar o restaurante com os dados retornados
	restaurant.ID = dbRestaurant.ID
	restaurant.CreatedAt = dbRestaurant.CreatedAt.Time
	restaurant.UpdatedAt = dbRestaurant.UpdatedAt.Time

	// Criar endereço se fornecido
	if restaurant.Address != nil {
		addrParams := database.CreateRestaurantAddressParams{
			RestaurantID: restaurant.ID,
			Street:       restaurant.Address.Street,
			Number:       restaurant.Address.Number,
			City:         restaurant.Address.City,
			State:        restaurant.Address.State,
			ZipCode:      restaurant.Address.ZipCode,
			Lat:          restaurant.Address.Lat,
			Lng:          restaurant.Address.Lng,
		}
		if restaurant.Address.Complement != "" {
			addrParams.Complement = pgtype.Text{String: restaurant.Address.Complement, Valid: true}
		}

		dbAddress, err := r.queries.CreateRestaurantAddress(ctx, addrParams)
		if err != nil {
			return fmt.Errorf("restaurant repository: create address: %w", err)
		}

		restaurant.Address.ID = dbAddress.ID
	}

	return nil
}

// GetByID busca um restaurante por ID
func (r *RestaurantRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Restaurant, error) {
	dbRestaurant, err := r.queries.GetRestaurantByID(ctx, id)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("restaurant repository: restaurant not found: %w", err)
		}
		return nil, fmt.Errorf("restaurant repository: get by id: %w", err)
	}

	// Carregar relacionamentos
	address, _ := r.queries.GetRestaurantAddress(ctx, dbRestaurant.ID)
	openingHours, _ := r.queries.GetOpeningHoursByRestaurant(ctx, dbRestaurant.ID)
	paymentMethods, _ := r.queries.GetPaymentMethodsByRestaurant(ctx, dbRestaurant.ID)

	return r.toDomain(&dbRestaurant, &address, openingHours, paymentMethods)
}

// GetBySlug busca um restaurante por slug
func (r *RestaurantRepository) GetBySlug(ctx context.Context, slug string) (*domain.Restaurant, error) {
	dbRestaurant, err := r.queries.GetRestaurantBySlug(ctx, slug)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, fmt.Errorf("restaurant repository: restaurant not found: %w", err)
		}
		return nil, fmt.Errorf("restaurant repository: get by slug: %w", err)
	}

	// Carregar relacionamentos
	address, _ := r.queries.GetRestaurantAddress(ctx, dbRestaurant.ID)
	openingHours, _ := r.queries.GetOpeningHoursByRestaurant(ctx, dbRestaurant.ID)
	paymentMethods, _ := r.queries.GetPaymentMethodsByRestaurant(ctx, dbRestaurant.ID)

	return r.toDomain(&dbRestaurant, &address, openingHours, paymentMethods)
}

// SlugExists verifica se um slug já existe
func (r *RestaurantRepository) SlugExists(ctx context.Context, slug string) (bool, error) {
	_, err := r.queries.GetRestaurantBySlug(ctx, slug)
	if err != nil {
		if err == pgx.ErrNoRows {
			return false, nil
		}
		return false, fmt.Errorf("restaurant repository: check slug exists: %w", err)
	}
	return true, nil
}

// List lista restaurantes com paginação
func (r *RestaurantRepository) List(ctx context.Context, limit, offset int32) ([]*domain.Restaurant, error) {
	dbRestaurants, err := r.queries.ListRestaurants(ctx, database.ListRestaurantsParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, fmt.Errorf("restaurant repository: list: %w", err)
	}

	restaurants := make([]*domain.Restaurant, 0, len(dbRestaurants))
	for _, dbRestaurant := range dbRestaurants {
		// Carregar relacionamentos para cada restaurante
		address, _ := r.queries.GetRestaurantAddress(ctx, dbRestaurant.ID)
		openingHours, _ := r.queries.GetOpeningHoursByRestaurant(ctx, dbRestaurant.ID)
		paymentMethods, _ := r.queries.GetPaymentMethodsByRestaurant(ctx, dbRestaurant.ID)

		restaurant, err := r.toDomain(&dbRestaurant, &address, openingHours, paymentMethods)
		if err != nil {
			return nil, err
		}
		restaurants = append(restaurants, restaurant)
	}

	return restaurants, nil
}

// UpdateStatus atualiza o status de um restaurante
func (r *RestaurantRepository) UpdateStatus(ctx context.Context, id uuid.UUID, status string) error {
	_, err := r.queries.UpdateRestaurantStatus(ctx, database.UpdateRestaurantStatusParams{
		ID:     id,
		Status: status,
	})
	if err != nil {
		return fmt.Errorf("restaurant repository: update status: %w", err)
	}
	return nil
}

// CreateAddress cria ou atualiza o endereço de um restaurante
func (r *RestaurantRepository) CreateAddress(ctx context.Context, address *domain.Address) error {
	params := database.CreateRestaurantAddressParams{
		RestaurantID: address.RestaurantID,
		Street:       address.Street,
		Number:       address.Number,
		City:         address.City,
		State:        address.State,
		ZipCode:      address.ZipCode,
		Lat:          address.Lat,
		Lng:          address.Lng,
	}
	if address.Complement != "" {
		params.Complement = pgtype.Text{String: address.Complement, Valid: true}
	}

	dbAddress, err := r.queries.CreateRestaurantAddress(ctx, params)
	if err != nil {
		return fmt.Errorf("restaurant repository: create address: %w", err)
	}

	address.ID = dbAddress.ID
	return nil
}

// UpdateAddress atualiza o endereço de um restaurante
func (r *RestaurantRepository) UpdateAddress(ctx context.Context, address *domain.Address) error {
	params := database.UpdateRestaurantAddressParams{
		RestaurantID: address.RestaurantID,
		Street:       address.Street,
		Number:       address.Number,
		City:         address.City,
		State:        address.State,
		ZipCode:      address.ZipCode,
		Lat:          address.Lat,
		Lng:          address.Lng,
	}
	if address.Complement != "" {
		params.Complement = pgtype.Text{String: address.Complement, Valid: true}
	}

	dbAddress, err := r.queries.UpdateRestaurantAddress(ctx, params)
	if err != nil {
		return fmt.Errorf("restaurant repository: update address: %w", err)
	}

	address.ID = dbAddress.ID
	return nil
}

// GetAddress busca o endereço de um restaurante
func (r *RestaurantRepository) GetAddress(ctx context.Context, restaurantID uuid.UUID) (*domain.Address, error) {
	dbAddress, err := r.queries.GetRestaurantAddress(ctx, restaurantID)
	if err != nil {
		if err == pgx.ErrNoRows {
			return nil, nil
		}
		return nil, fmt.Errorf("restaurant repository: get address: %w", err)
	}

	return &domain.Address{
		ID:           dbAddress.ID,
		RestaurantID: dbAddress.RestaurantID,
		Street:       dbAddress.Street,
		Number:       dbAddress.Number,
		Complement:   dbAddress.Complement.String,
		City:         dbAddress.City,
		State:        dbAddress.State,
		ZipCode:      dbAddress.ZipCode,
		Lat:          dbAddress.Lat,
		Lng:          dbAddress.Lng,
	}, nil
}

// CreateOpeningHour cria um horário de funcionamento
func (r *RestaurantRepository) CreateOpeningHour(ctx context.Context, hour *domain.OpeningHour) error {
	dbHour, err := r.queries.CreateOpeningHour(ctx, database.CreateOpeningHourParams{
		RestaurantID: hour.RestaurantID,
		Weekday:      int32(hour.Weekday),
		OpensAt:      int32(hour.OpensAt),
		ClosesAt:     int32(hour.ClosesAt),
	})
	if err != nil {
		return fmt.Errorf("restaurant repository: create opening hour: %w", err)
	}

	hour.ID = dbHour.ID
	return nil
}

// DeleteOpeningHoursByRestaurant remove todos os horários de um restaurante
func (r *RestaurantRepository) DeleteOpeningHoursByRestaurant(ctx context.Context, restaurantID uuid.UUID) error {
	err := r.queries.DeleteOpeningHoursByRestaurant(ctx, restaurantID)
	if err != nil {
		return fmt.Errorf("restaurant repository: delete opening hours: %w", err)
	}
	return nil
}

// GetOpeningHours busca os horários de funcionamento de um restaurante
func (r *RestaurantRepository) GetOpeningHours(ctx context.Context, restaurantID uuid.UUID) ([]*domain.OpeningHour, error) {
	dbHours, err := r.queries.GetOpeningHoursByRestaurant(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("restaurant repository: get opening hours: %w", err)
	}

	hours := make([]*domain.OpeningHour, 0, len(dbHours))
	for _, dbHour := range dbHours {
		hours = append(hours, &domain.OpeningHour{
			ID:           dbHour.ID,
			RestaurantID: dbHour.RestaurantID,
			Weekday:      int(dbHour.Weekday),
			OpensAt:      int(dbHour.OpensAt),
			ClosesAt:     int(dbHour.ClosesAt),
		})
	}

	return hours, nil
}

// CreatePaymentMethod cria um método de pagamento
func (r *RestaurantRepository) CreatePaymentMethod(ctx context.Context, method *domain.PaymentMethod) error {
	dbMethod, err := r.queries.CreatePaymentMethod(ctx, database.CreatePaymentMethodParams{
		RestaurantID: method.RestaurantID,
		Method:       method.Method,
	})
	if err != nil {
		return fmt.Errorf("restaurant repository: create payment method: %w", err)
	}

	method.ID = dbMethod.ID
	return nil
}

// DeletePaymentMethodsByRestaurant remove todos os métodos de pagamento de um restaurante
func (r *RestaurantRepository) DeletePaymentMethodsByRestaurant(ctx context.Context, restaurantID uuid.UUID) error {
	err := r.queries.DeletePaymentMethodsByRestaurant(ctx, restaurantID)
	if err != nil {
		return fmt.Errorf("restaurant repository: delete payment methods: %w", err)
	}
	return nil
}

// GetPaymentMethods busca os métodos de pagamento de um restaurante
func (r *RestaurantRepository) GetPaymentMethods(ctx context.Context, restaurantID uuid.UUID) ([]*domain.PaymentMethod, error) {
	dbMethods, err := r.queries.GetPaymentMethodsByRestaurant(ctx, restaurantID)
	if err != nil {
		return nil, fmt.Errorf("restaurant repository: get payment methods: %w", err)
	}

	methods := make([]*domain.PaymentMethod, 0, len(dbMethods))
	for _, dbMethod := range dbMethods {
		methods = append(methods, &domain.PaymentMethod{
			ID:           dbMethod.ID,
			RestaurantID: dbMethod.RestaurantID,
			Method:       dbMethod.Method,
		})
	}

	return methods, nil
}

// toDomain converte modelos do banco para entidades de domínio
func (r *RestaurantRepository) toDomain(
	dbRestaurant *database.Restaurant,
	dbAddress *database.RestaurantAddress,
	dbHours []database.RestaurantOpeningHour,
	dbMethods []database.RestaurantPaymentMethod,
) (*domain.Restaurant, error) {
	restaurant := &domain.Restaurant{
		ID:                 dbRestaurant.ID,
		Name:               dbRestaurant.Name,
		Slug:               dbRestaurant.Slug,
		Status:             dbRestaurant.Status,
		Rating:             int(dbRestaurant.Rating),
		TotalReviews:       int(dbRestaurant.TotalReviews),
		DeliveryFee:        dbRestaurant.DeliveryFee,
		MinOrderValue:      dbRestaurant.MinOrderValue,
		PreparationTimeMin: int(dbRestaurant.PreparationTimeMin),
		SupportsPickup:     dbRestaurant.SupportsPickup,
		SupportsDelivery:   dbRestaurant.SupportsDelivery,
		CreatedAt:          dbRestaurant.CreatedAt.Time,
		UpdatedAt:          dbRestaurant.UpdatedAt.Time,
	}

	if dbRestaurant.Description.Valid {
		restaurant.Description = dbRestaurant.Description.String
	}
	if dbRestaurant.Category.Valid {
		restaurant.Category = dbRestaurant.Category.String
	}
	if dbRestaurant.LogoUrl.Valid {
		restaurant.LogoURL = dbRestaurant.LogoUrl.String
	}
	if dbRestaurant.BannerUrl.Valid {
		restaurant.BannerURL = dbRestaurant.BannerUrl.String
	}

	// Converter endereço
	if dbAddress != nil {
		restaurant.Address = &domain.Address{
			ID:           dbAddress.ID,
			RestaurantID: dbAddress.RestaurantID,
			Street:       dbAddress.Street,
			Number:       dbAddress.Number,
			City:         dbAddress.City,
			State:        dbAddress.State,
			ZipCode:      dbAddress.ZipCode,
			Lat:          dbAddress.Lat,
			Lng:          dbAddress.Lng,
		}
		if dbAddress.Complement.Valid {
			restaurant.Address.Complement = dbAddress.Complement.String
		}
	}

	// Converter horários
	restaurant.OpeningHours = make([]domain.OpeningHour, 0, len(dbHours))
	for _, dbHour := range dbHours {
		restaurant.OpeningHours = append(restaurant.OpeningHours, domain.OpeningHour{
			ID:           dbHour.ID,
			RestaurantID: dbHour.RestaurantID,
			Weekday:      int(dbHour.Weekday),
			OpensAt:      int(dbHour.OpensAt),
			ClosesAt:     int(dbHour.ClosesAt),
		})
	}

	// Converter métodos de pagamento
	restaurant.PaymentMethods = make([]domain.PaymentMethod, 0, len(dbMethods))
	for _, dbMethod := range dbMethods {
		restaurant.PaymentMethods = append(restaurant.PaymentMethods, domain.PaymentMethod{
			ID:           dbMethod.ID,
			RestaurantID: dbMethod.RestaurantID,
			Method:       dbMethod.Method,
		})
	}

	return restaurant, nil
}

// WithTx retorna um repository com transação
func (r *RestaurantRepository) WithTx(tx pgx.Tx) *RestaurantRepository {
	return &RestaurantRepository{
		queries: r.queries.WithTx(tx),
	}
}
