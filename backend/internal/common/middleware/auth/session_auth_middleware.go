package auth

import (
	"sync"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/auth/session"

	"github.com/hashicorp/vault/api"
	"github.com/hengadev/encx"
)

// DefaultServiceKeyCacheTTL is the default time-to-live for cached service key hashes.
const DefaultServiceKeyCacheTTL = 5 * time.Minute

// SessionAuthMiddleware implements AuthMiddleware using session repository
type SessionAuthMiddleware struct {
	sessionRepo session.SessionRepository
	crypto      encx.CryptoService
	vaultClient *api.Client

	// secretReader abstracts Vault reads so tests can inject mocks.
	secretReader SecretReader

	// serviceKeyCache stores key_hash values fetched from Vault per service.
	serviceKeyCache map[string]*cachedKeyEntry
	cacheTTL       time.Duration
	cacheMu        sync.RWMutex
}

// cachedKeyEntry is a single cached key-hash entry with its fetch timestamp.
type cachedKeyEntry struct {
	keyHash   string
	fetchedAt time.Time
}

// SessionAuthMiddlewareOption configures a SessionAuthMiddleware after creation.
type SessionAuthMiddlewareOption func(*SessionAuthMiddleware)

// WithServiceKeyCacheTTL sets the TTL for cached service key hashes.
func WithServiceKeyCacheTTL(ttl time.Duration) SessionAuthMiddlewareOption {
	return func(m *SessionAuthMiddleware) {
		m.cacheTTL = ttl
	}
}

// WithSecretReader sets a custom SecretReader (useful for testing).
func WithSecretReader(reader SecretReader) SessionAuthMiddlewareOption {
	return func(m *SessionAuthMiddleware) {
		m.secretReader = reader
	}
}

// NewSessionAuthMiddleware creates a new session-based auth middleware.
func NewSessionAuthMiddleware(sessionRepo session.SessionRepository, crypto encx.CryptoService, vaultClient *api.Client, opts ...SessionAuthMiddlewareOption) AuthMiddleware {
	m := &SessionAuthMiddleware{
		sessionRepo:     sessionRepo,
		crypto:          crypto,
		vaultClient:     vaultClient,
		serviceKeyCache: make(map[string]*cachedKeyEntry),
		cacheTTL:        DefaultServiceKeyCacheTTL,
	}

	// Default: use the real Vault client as the secret reader.
	if vaultClient != nil {
		m.secretReader = NewVaultSecretReader(vaultClient)
	}

	for _, opt := range opts {
		opt(m)
	}

	return m
}
