package sessionRepository

import (
	"context"
)

func (r *SessionRepository) FindSessionByTokenHash(ctx context.Context, tokenHash string) ([]byte, error) {
	sessionID, err := r.FindSessionIDByTokenHash(ctx, tokenHash)
	if err != nil {
		return nil, err
	}

	session, err := r.FindSessionByID(ctx, sessionID)
	if err != nil {
		return nil, err
	}

	return session, nil
}
