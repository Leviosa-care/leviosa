package throttlerService

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
	// TODO: add this later with the right error message and checking
	// "github.com/hengadev/leviosa/internal/domain"
	// rp "github.com/hengadev/leviosa/internal/repository"
)

func (s *service) RegisterAttempt(ctx context.Context, email string) error {
	// TODO: make better error handling for this function
	now := time.Now()
	infoEncoded, err := s.repo.IsLocked(ctx, email)
	if err != nil {
		return fmt.Errorf("locked for the current email")
	}
	var info Info
	if err := json.Unmarshal(infoEncoded, &info); err != nil {

	}
	if time.Now().After(info.LockedUntil) {
		return fmt.Errorf("user locked, too many request: %w", err)
	}
	if err := s.repo.MakeAttempt(ctx, email, now); err != nil {
		return fmt.Errorf("user made a login attempt: %w", err)
	}
	return nil
}
