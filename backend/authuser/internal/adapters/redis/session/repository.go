package sessionRepository

import (
	"fmt"

	"github.com/Leviosa-care/authuser/internal/ports"
	"github.com/redis/go-redis/v9"
)

type SessionRepository struct {
	client *redis.Client
}

const (
	ResetSessionKeyPrefix = "authuser:reset_session:"
)

func New(client *redis.Client) ports.SessionRepository {
	return &SessionRepository{client: client}
}

// FormatResetSessionKey formats a reset session key for Redis storage
// This function is public to allow consistent key formatting in tests
func FormatResetSessionKey(tokenHash string) string {
	return fmt.Sprintf("%s%s", ResetSessionKeyPrefix, tokenHash)
}
