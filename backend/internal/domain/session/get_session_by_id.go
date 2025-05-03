package sessionService

import (
	"context"
	"encoding/json"
	"errors"
	"strings"

	"github.com/hengadev/leviosa/internal/domain"
	rp "github.com/hengadev/leviosa/internal/repository"
)

func (s *service) GetSession(ctx context.Context, sessionID string) (*Session, error) {
	if strings.TrimSpace(sessionID) == "" {
		return nil, domain.NewNotFoundErr(errors.New("empty session ID"))
	}

	sessionEncoded, err := s.Repo.FindSessionByID(ctx, sessionID)
	if err != nil {
		switch {
		case errors.Is(err, rp.ErrNotFound):
			return nil, domain.NewNotFoundErr(err)
		case errors.Is(err, rp.ErrDatabase):
			return nil, domain.NewQueryFailedErr(err)
		default:
			return nil, domain.NewUnexpectTypeErr(err)
		}
	}

	var session Session
	if err = json.Unmarshal(sessionEncoded, &session); err != nil {
		return nil, domain.NewJSONUnmarshalErr(err)
	}

	return &session, nil
}
