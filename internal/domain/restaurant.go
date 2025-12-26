package domain

import (
	"time"

	"github.com/google/uuid"
)

// Restaurant é o Aggregate Root do domínio de restaurantes
type Restaurant struct {
	ID                 uuid.UUID
	Name               string
	Slug               string // "pizza-do-joao" (Unique)
	Description        string
	Status             string // "DRAFT", "OPEN", "CLOSED", "SUSPENDED"
	Category           string // "Pizza", "Burgers"
	Rating             int    // 0, 1, 2, 3, 4 ou 5
	TotalReviews       int    // Default 0
	IsOpen             bool   // Campo computado
	DeliveryFee        int64  // unidades monetárias (centavos)
	MinOrderValue      int64  // unidades monetárias (centavos)
	PreparationTimeMin int    // em minutos
	SupportsPickup     bool
	SupportsDelivery   bool
	LogoURL            string // Não obrigatório
	BannerURL          string // Não obrigatório
	CreatedAt          time.Time
	UpdatedAt          time.Time

	// Relacionamentos (Carregados com o Aggregate)
	Address        *Address
	OpeningHours   []OpeningHour
	PaymentMethods []PaymentMethod
}

// Address representa o endereço de um restaurante
type Address struct {
	ID           uuid.UUID
	RestaurantID uuid.UUID
	Street       string
	Number       string
	Complement   string
	City         string
	State        string // char(2)
	ZipCode      string
	Lat          float64
	Lng          float64
}

// OpeningHour representa um horário de funcionamento de um restaurante
type OpeningHour struct {
	ID           uuid.UUID
	RestaurantID uuid.UUID
	Weekday      int // 0=Domingo, 1=Segunda ... 6=Sábado
	OpensAt      int // Minutos a partir da meia-noite (ex: 480 = 08:00)
	ClosesAt     int // Minutos a partir da meia-noite (ex: 120 = 02:00 do dia seguinte)
}

// PaymentMethod representa um método de pagamento aceito pelo restaurante
type PaymentMethod struct {
	ID           uuid.UUID
	RestaurantID uuid.UUID
	Method       string // "PIX", "CREDIT_CARD", "DEBIT_CARD"
}

// Constantes para Status
const (
	StatusDraft     = "DRAFT"
	StatusOpen      = "OPEN"
	StatusClosed    = "CLOSED"
	StatusSuspended = "SUSPENDED"
)

// Constantes para métodos de pagamento
const (
	PaymentMethodPIX        = "PIX"
	PaymentMethodCreditCard = "CREDIT_CARD"
	PaymentMethodDebitCard  = "DEBIT_CARD"
)

// CalculateIsOpen calcula se o restaurante está aberto no momento atual
// Retorna true apenas se: Status == OPEN E horário atual está dentro de um intervalo válido
func (r *Restaurant) CalculateIsOpen(now time.Time) bool {
	if r.Status != StatusOpen {
		return false
	}

	if len(r.OpeningHours) == 0 {
		return false
	}

	currentWeekday := int(now.Weekday())
	currentMinutes := now.Hour()*60 + now.Minute()

	for _, hour := range r.OpeningHours {
		if hour.Weekday == currentWeekday {
			// Se closes_at < opens_at, significa que fecha no dia seguinte
			if hour.ClosesAt < hour.OpensAt {
				// Caso especial: horário que cruza a meia-noite
				// Exemplo: abre 22:00 (1320min) e fecha 02:00 (120min)
				if currentMinutes >= hour.OpensAt || currentMinutes < hour.ClosesAt {
					return true
				}
			} else {
				// Caso normal: horário no mesmo dia
				if currentMinutes >= hour.OpensAt && currentMinutes < hour.ClosesAt {
					return true
				}
			}
		}
	}

	return false
}

