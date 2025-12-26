package usecase

import (
	"context"
	"fmt"

	"github.com/google/uuid"

	"gastro-go/internal/domain"
	"gastro-go/internal/repository"
)

// UpdateOpeningHoursUseCase implementa o caso de uso de atualizar horários de funcionamento
type UpdateOpeningHoursUseCase struct {
	repo repository.RestaurantRepositoryInterface
}

// NewUpdateOpeningHoursUseCase cria uma nova instância do use case
func NewUpdateOpeningHoursUseCase(repo repository.RestaurantRepositoryInterface) *UpdateOpeningHoursUseCase {
	return &UpdateOpeningHoursUseCase{
		repo: repo,
	}
}

// OpeningHourInput representa um horário de funcionamento
type OpeningHourInput struct {
	Weekday  int // 0=Domingo, 1=Segunda ... 6=Sábado
	OpensAt  int // Minutos a partir da meia-noite (0-1439)
	ClosesAt int // Minutos a partir da meia-noite (0-1439)
}

// UpdateOpeningHoursInput representa os dados de entrada para atualizar horários
type UpdateOpeningHoursInput struct {
	RestaurantID uuid.UUID
	Hours        []OpeningHourInput
}

// Execute executa o caso de uso de atualizar horários
func (uc *UpdateOpeningHoursUseCase) Execute(ctx context.Context, input UpdateOpeningHoursInput) error {
	// Verificar se o restaurante existe
	_, err := uc.repo.GetByID(ctx, input.RestaurantID)
	if err != nil {
		return fmt.Errorf("update opening hours usecase: %w", err)
	}

	// Validar horários
	for _, hour := range input.Hours {
		if hour.Weekday < 0 || hour.Weekday > 6 {
			return fmt.Errorf("update opening hours usecase: weekday must be between 0 and 6")
		}
		if hour.OpensAt < 0 || hour.OpensAt >= 1440 {
			return fmt.Errorf("update opening hours usecase: opens_at must be between 0 and 1439")
		}
		if hour.ClosesAt < 0 || hour.ClosesAt >= 1440 {
			return fmt.Errorf("update opening hours usecase: closes_at must be between 0 and 1439")
		}
	}

	// Verificar colisões de horários no mesmo dia
	if err := uc.validateCollisions(input.Hours); err != nil {
		return fmt.Errorf("update opening hours usecase: %w", err)
	}

	// Deletar horários existentes
	if err := uc.repo.DeleteOpeningHoursByRestaurant(ctx, input.RestaurantID); err != nil {
		return fmt.Errorf("update opening hours usecase: delete existing hours: %w", err)
	}

	// Criar novos horários
	for _, hourInput := range input.Hours {
		hour := &domain.OpeningHour{
			ID:           uuid.New(),
			RestaurantID: input.RestaurantID,
			Weekday:      hourInput.Weekday,
			OpensAt:      hourInput.OpensAt,
			ClosesAt:     hourInput.ClosesAt,
		}

		if err := uc.repo.CreateOpeningHour(ctx, hour); err != nil {
			return fmt.Errorf("update opening hours usecase: create hour: %w", err)
		}
	}

	return nil
}

// validateCollisions verifica se há colisões de horários no mesmo dia
func (uc *UpdateOpeningHoursUseCase) validateCollisions(hours []OpeningHourInput) error {
	// Agrupar por dia da semana
	hoursByWeekday := make(map[int][]OpeningHourInput)
	for _, hour := range hours {
		hoursByWeekday[hour.Weekday] = append(hoursByWeekday[hour.Weekday], hour)
	}

	// Verificar colisões em cada dia
	for weekday, dayHours := range hoursByWeekday {
		if len(dayHours) <= 1 {
			continue
		}

		// Verificar todas as combinações de horários no mesmo dia
		for i := 0; i < len(dayHours); i++ {
			for j := i + 1; j < len(dayHours); j++ {
				if uc.hoursOverlap(dayHours[i], dayHours[j]) {
					return fmt.Errorf("opening hours overlap on weekday %d", weekday)
				}
			}
		}
	}

	return nil
}

// hoursOverlap verifica se dois horários se sobrepõem
func (uc *UpdateOpeningHoursUseCase) hoursOverlap(h1, h2 OpeningHourInput) bool {
	// Caso 1: Horário normal (opens_at < closes_at)
	// Caso 2: Horário que cruza meia-noite (closes_at < opens_at)

	// Normalizar ambos os horários para comparação
	// Se closes_at < opens_at, significa que fecha no dia seguinte
	// Vamos tratar como dois intervalos separados ou um intervalo que cruza

	// Simplificação: verificar se há sobreposição considerando ambos os casos
	// Para horários normais
	if h1.OpensAt < h1.ClosesAt && h2.OpensAt < h2.ClosesAt {
		// Ambos são horários normais
		return !(h1.ClosesAt <= h2.OpensAt || h2.ClosesAt <= h1.OpensAt)
	}

	// Se um cruza meia-noite, precisamos verificar de forma diferente
	// Para simplificar, vamos considerar que se closes_at < opens_at,
	// o horário vai de opens_at até 1439 e depois de 0 até closes_at
	// Isso significa que sempre há sobreposição potencial

	// Verificação mais simples: se ambos cruzam meia-noite, sempre há sobreposição
	if h1.ClosesAt < h1.OpensAt && h2.ClosesAt < h2.OpensAt {
		return true
	}

	// Se apenas um cruza meia-noite
	if h1.ClosesAt < h1.OpensAt {
		// h1 cruza meia-noite: de h1.OpensAt até 1439 e de 0 até h1.ClosesAt
		// h2 é normal: de h2.OpensAt até h2.ClosesAt
		return h2.OpensAt >= h1.OpensAt || h2.ClosesAt <= h1.ClosesAt || h2.OpensAt < h1.ClosesAt || h2.ClosesAt > h1.OpensAt
	}

	if h2.ClosesAt < h2.OpensAt {
		// h2 cruza meia-noite
		return h1.OpensAt >= h2.OpensAt || h1.ClosesAt <= h2.ClosesAt || h1.OpensAt < h2.ClosesAt || h1.ClosesAt > h2.OpensAt
	}

	return false
}

