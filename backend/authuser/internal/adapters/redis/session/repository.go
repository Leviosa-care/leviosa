package sessionRepository

import (
	"fmt"

	"github.com/Leviosa-care/authuser/internal/ports"
	"github.com/redis/go-redis/v9"
)

const (
	SessionKeyPrefix      = "authuser:session:"
	TokenKeyPrefix        = "authuser:token:"
	AccessTokenKeyPrefix  = "authuser:access:"
	RefreshTokenKeyPrefix = "authuser:refresh:"
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

func FormatAccessTokenKey(accessTokenHash string) string {
	return fmt.Sprintf("%s%s", AccessTokenKeyPrefix, accessTokenHash)
}

func FormatRefreshTokenKey(refreshTokenHash string) string {
	return fmt.Sprintf("%s%s", RefreshTokenKeyPrefix, refreshTokenHash)
}
