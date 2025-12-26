package handler

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"

	"gastro-go/internal/usecase"
)

// RestaurantHandler gerencia os endpoints HTTP relacionados a restaurantes
type RestaurantHandler struct {
	createUseCase           *usecase.CreateRestaurantUseCase
	listUseCase             *usecase.ListRestaurantsUseCase
	getBySlugUseCase        *usecase.GetRestaurantBySlugUseCase
	openUseCase             *usecase.OpenRestaurantUseCase
	closeUseCase            *usecase.CloseRestaurantUseCase
	updateOpeningHoursUseCase *usecase.UpdateOpeningHoursUseCase
	updatePaymentMethodsUseCase *usecase.UpdatePaymentMethodsUseCase
}

// NewRestaurantHandler cria uma nova instância do handler
func NewRestaurantHandler(
	createUseCase *usecase.CreateRestaurantUseCase,
	listUseCase *usecase.ListRestaurantsUseCase,
	getBySlugUseCase *usecase.GetRestaurantBySlugUseCase,
	openUseCase *usecase.OpenRestaurantUseCase,
	closeUseCase *usecase.CloseRestaurantUseCase,
	updateOpeningHoursUseCase *usecase.UpdateOpeningHoursUseCase,
	updatePaymentMethodsUseCase *usecase.UpdatePaymentMethodsUseCase,
) *RestaurantHandler {
	return &RestaurantHandler{
		createUseCase:              createUseCase,
		listUseCase:                listUseCase,
		getBySlugUseCase:          getBySlugUseCase,
		openUseCase:                openUseCase,
		closeUseCase:               closeUseCase,
		updateOpeningHoursUseCase:  updateOpeningHoursUseCase,
		updatePaymentMethodsUseCase: updatePaymentMethodsUseCase,
	}
}

// CreateRestaurantRequest representa o payload de criação de restaurante
type CreateRestaurantRequest struct {
	Name               string                  `json:"name"`
	Slug               string                  `json:"slug,omitempty"`
	Description        string                  `json:"description,omitempty"`
	Category           string                  `json:"category,omitempty"`
	DeliveryFee        int64                   `json:"delivery_fee"`
	MinOrderValue      int64                   `json:"min_order_value"`
	PreparationTimeMin int                     `json:"preparation_time_min"`
	SupportsPickup     bool                    `json:"supports_pickup"`
	SupportsDelivery   bool                    `json:"supports_delivery"`
	LogoURL            string                  `json:"logo_url,omitempty"`
	BannerURL          string                  `json:"banner_url,omitempty"`
	Address            *CreateAddressRequest    `json:"address,omitempty"`
}

// CreateAddressRequest representa o payload de criação de endereço
type CreateAddressRequest struct {
	Street     string  `json:"street"`
	Number     string  `json:"number"`
	Complement string  `json:"complement,omitempty"`
	City       string  `json:"city"`
	State      string  `json:"state"`
	ZipCode    string  `json:"zip_code"`
	Lat        float64 `json:"lat"`
	Lng        float64 `json:"lng"`
}

// UpdateOpeningHoursRequest representa o payload de atualização de horários
type UpdateOpeningHoursRequest struct {
	Hours []OpeningHourRequest `json:"hours"`
}

// OpeningHourRequest representa um horário de funcionamento
type OpeningHourRequest struct {
	Weekday  int `json:"weekday"`
	OpensAt  int `json:"opens_at"`
	ClosesAt int `json:"closes_at"`
}

// UpdatePaymentMethodsRequest representa o payload de atualização de métodos de pagamento
type UpdatePaymentMethodsRequest struct {
	Methods []string `json:"methods"`
}

// CreateRestaurant cria um novo restaurante
// POST /restaurants
func (h *RestaurantHandler) CreateRestaurant(c echo.Context) error {
	var req CreateRestaurantRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	var addressInput *usecase.CreateAddressInput
	if req.Address != nil {
		addressInput = &usecase.CreateAddressInput{
			Street:     req.Address.Street,
			Number:     req.Address.Number,
			Complement: req.Address.Complement,
			City:       req.Address.City,
			State:      req.Address.State,
			ZipCode:    req.Address.ZipCode,
			Lat:        req.Address.Lat,
			Lng:        req.Address.Lng,
		}
	}

	input := usecase.CreateRestaurantInput{
		Name:               req.Name,
		Slug:               req.Slug,
		Description:        req.Description,
		Category:           req.Category,
		DeliveryFee:        req.DeliveryFee,
		MinOrderValue:      req.MinOrderValue,
		PreparationTimeMin: req.PreparationTimeMin,
		SupportsPickup:     req.SupportsPickup,
		SupportsDelivery:   req.SupportsDelivery,
		LogoURL:            req.LogoURL,
		BannerURL:          req.BannerURL,
		Address:            addressInput,
	}

	restaurant, err := h.createUseCase.Execute(c.Request().Context(), input)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusCreated, restaurant)
}

// ListRestaurants lista restaurantes com paginação
// GET /restaurants
func (h *RestaurantHandler) ListRestaurants(c echo.Context) error {
	limitStr := c.QueryParam("limit")
	offsetStr := c.QueryParam("offset")

	limit := int32(20) // Default
	offset := int32(0)

	if limitStr != "" {
		l, err := strconv.ParseInt(limitStr, 10, 32)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid limit parameter",
			})
		}
		limit = int32(l)
	}

	if offsetStr != "" {
		o, err := strconv.ParseInt(offsetStr, 10, 32)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{
				"error": "invalid offset parameter",
			})
		}
		offset = int32(o)
	}

	input := usecase.ListRestaurantsInput{
		Limit:  limit,
		Offset: offset,
	}

	restaurants, err := h.listUseCase.Execute(c.Request().Context(), input)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, restaurants)
}

// GetRestaurantBySlug busca um restaurante por slug
// GET /restaurants/{slug}
func (h *RestaurantHandler) GetRestaurantBySlug(c echo.Context) error {
	slug := c.Param("slug")
	if slug == "" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "slug is required",
		})
	}

	restaurant, err := h.getBySlugUseCase.Execute(c.Request().Context(), slug)
	if err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, restaurant)
}

// OpenRestaurant abre um restaurante
// PATCH /restaurants/{id}/open
func (h *RestaurantHandler) OpenRestaurant(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid restaurant id",
		})
	}

	if err := h.openUseCase.Execute(c.Request().Context(), id); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "restaurant opened successfully",
	})
}

// CloseRestaurant fecha um restaurante
// PATCH /restaurants/{id}/close
func (h *RestaurantHandler) CloseRestaurant(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid restaurant id",
		})
	}

	if err := h.closeUseCase.Execute(c.Request().Context(), id); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "restaurant closed successfully",
	})
}

// UpdateOpeningHours atualiza os horários de funcionamento
// PUT /restaurants/{id}/hours
func (h *RestaurantHandler) UpdateOpeningHours(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid restaurant id",
		})
	}

	var req UpdateOpeningHoursRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	hours := make([]usecase.OpeningHourInput, 0, len(req.Hours))
	for _, h := range req.Hours {
		hours = append(hours, usecase.OpeningHourInput{
			Weekday:  h.Weekday,
			OpensAt:  h.OpensAt,
			ClosesAt: h.ClosesAt,
		})
	}

	input := usecase.UpdateOpeningHoursInput{
		RestaurantID: id,
		Hours:        hours,
	}

	if err := h.updateOpeningHoursUseCase.Execute(c.Request().Context(), input); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "opening hours updated successfully",
	})
}

// UpdatePaymentMethods atualiza os métodos de pagamento
// PUT /restaurants/{id}/payments
func (h *RestaurantHandler) UpdatePaymentMethods(c echo.Context) error {
	idStr := c.Param("id")
	id, err := uuid.Parse(idStr)
	if err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid restaurant id",
		})
	}

	var req UpdatePaymentMethodsRequest
	if err := c.Bind(&req); err != nil {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": "invalid request body",
		})
	}

	input := usecase.UpdatePaymentMethodsInput{
		RestaurantID: id,
		Methods:      req.Methods,
	}

	if err := h.updatePaymentMethodsUseCase.Execute(c.Request().Context(), input); err != nil {
		return h.handleError(c, err)
	}

	return c.JSON(http.StatusOK, map[string]string{
		"message": "payment methods updated successfully",
	})
}

// handleError trata erros e retorna a resposta HTTP apropriada
func (h *RestaurantHandler) handleError(c echo.Context, err error) error {
	errMsg := err.Error()

	// Verificar tipo de erro
	if errors.Is(err, errors.New("conflict")) || errMsg == "create restaurant usecase: slug already exists: conflict" {
		return c.JSON(http.StatusConflict, map[string]string{
			"error": "slug already exists",
		})
	}

	if errMsg == "restaurant repository: restaurant not found: no rows in result set" ||
		errMsg == "get restaurant by slug usecase: restaurant repository: restaurant not found: no rows in result set" ||
		errMsg == "open restaurant usecase: restaurant repository: restaurant not found: no rows in result set" ||
		errMsg == "close restaurant usecase: restaurant repository: restaurant not found: no rows in result set" {
		return c.JSON(http.StatusNotFound, map[string]string{
			"error": "restaurant not found",
		})
	}

	// Erros de validação (400)
	if errMsg == "create restaurant usecase: name is required" ||
		errMsg == "create restaurant usecase: delivery fee cannot be negative" ||
		errMsg == "create restaurant usecase: min order value cannot be negative" ||
		errMsg == "open restaurant usecase: restaurant must have an address to be opened" ||
		errMsg == "open restaurant usecase: restaurant must have opening hours to be opened" ||
		errMsg == "open restaurant usecase: restaurant must have at least one payment method to be opened" {
		return c.JSON(http.StatusBadRequest, map[string]string{
			"error": errMsg,
		})
	}

	// Erro genérico (500)
	return c.JSON(http.StatusInternalServerError, map[string]string{
		"error": "internal server error",
	})
}

