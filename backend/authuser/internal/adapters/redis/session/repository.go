package sessionRepository

import (
	"fmt"

	"github.com/Leviosa-care/authuser/internal/ports"
	"github.com/redis/go-redis/v9"
)

// NOTE: what do I store :
// session ID -> JSON
// token hash -> session ID

const (
	SessionKeyPrefix = "authuser:session:"
	TokenKeyPrefix   = "authuser:token:"
)

type SessionRepository struct {
	client *redis.Client
}

func New(client *redis.Client) ports.SessionRepository {
	return &SessionRepository{
		client: client,
	}
}

func FormatSessionKey(sessionID string) string {
	return fmt.Sprintf("%s%s", SessionKeyPrefix, sessionID)
}

func FormatTokenKey(tokenHash string) string {
	return fmt.Sprintf("%s%s", TokenKeyPrefix, tokenHash)
}
