package registerService

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/hengadev/leviosa/internal/domain"
	"github.com/hengadev/leviosa/internal/domain/event/models"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (s *service) CreateRegistration(ctx context.Context, userID, spotStr string, event *models.Event) error {
	day, month, year := parseEventTimeForRegistration(event)
	// NOTE: why is the spot str a string and not an int ?
	spot, err := strconv.Atoi(spotStr)
	if err != nil {
		return fmt.Errorf("convert spot from string to int: %w", err)
	}
	// offsetDuration := time.Duration(int(event.SessionDuration) * (spot - 1))
	offsetDuration := time.Duration(int(20) * (spot - 1))
	registrationBeginAt := event.BeginAt.Add(offsetDuration)
	_ = registrationBeginAt
	registration := NewRegistration(
		userID,
		"",
		"event",
		time.Now(),
		time.Now().Add(24*time.Hour),
		nil,
		nil,
	)
	err = s.Repo.AddRegistration(ctx, registration, day, year, month)
	if err != nil {
		switch {
		case errors.Is(err, rp.ErrContext):
			return err
		case errors.Is(err, rp.ErrDatabase):
			return domain.NewQueryFailedErr(err)
		case errors.Is(err, rp.ErrNotCreated):
			return domain.NewNotCreatedErr(err)
		default:
			return fmt.Errorf("add registration: %w", err)
		}
	}
	return nil
}

func parseEventTimeForRegistration(e *models.Event) (int, string, int) {
	month := strings.ToLower(time.Month(e.Month).String())
	return e.Day, month, e.Year
}
