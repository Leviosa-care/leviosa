package auth

import (
	"context"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/Leviosa-care/leviosa/backend/internal/common/contracts/services"
	"github.com/Leviosa-care/leviosa/backend/internal/common/ctxutil"
	"github.com/hengadev/encx"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// ---------------------------------------------------------------------------
// Mock SecretReader
// ---------------------------------------------------------------------------

// mockSecretReader is a testify mock implementing SecretReader.
type mockSecretReader struct {
	mock.Mock
}

func (m *mockSecretReader) Read(ctx context.Context, path string) (*SecretData, error) {
	args := m.Called(ctx, path)
	if v := args.Get(0); v == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*SecretData), args.Error(1)
}

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

// newTestMiddlewareCrypto creates a crypto service suitable for hashing keys in tests.
func newTestMiddlewareCrypto(t *testing.T) encx.CryptoService {
	t.Helper()
	c, err := NewTestCrypto(t)
	require.NoError(t, err)
	return c
}

// newRequestWithLogger creates an httptest.Request whose context carries a logger,
// so the middleware doesn't 500 on ctxutil.GetLoggerFromContext.
func newRequestWithLogger(method, target string) *http.Request {
	req := httptest.NewRequest(method, target, nil)
	ctx := context.WithValue(req.Context(), ctxutil.LoggerKey, slog.Default())
	return req.WithContext(ctx)
}

// computeKeyHash is a test helper that hashes a plaintext key using the provided crypto service.
func computeKeyHash(t *testing.T, crypto encx.CryptoService, key string) string {
	t.Helper()
	keyBytes, err := encx.SerializeValue(key)
	require.NoError(t, err)
	return crypto.HashBasic(context.Background(), keyBytes)
}

// buildMiddleware creates a SessionAuthMiddleware wired with a mock secret reader
// and an optional cache TTL override.
func buildMiddleware(t *testing.T, secretReader SecretReader, cacheTTL time.Duration) *SessionAuthMiddleware {
	t.Helper()
	crypto := newTestMiddlewareCrypto(t)
	m := NewSessionAuthMiddleware(nil, crypto, nil,
		WithSecretReader(secretReader),
		WithServiceKeyCacheTTL(cacheTTL),
	)
	return m.(*SessionAuthMiddleware)
}

// ---------------------------------------------------------------------------
// Header validation tests
// ---------------------------------------------------------------------------

func TestRequireServiceAuth_MissingServiceNameHeader(t *testing.T) {
	sr := new(mockSecretReader)
	mw := buildMiddleware(t, sr, time.Minute)

	nextCalled := false
	next := func(w http.ResponseWriter, r *http.Request) { nextCalled = true }

	req := newRequestWithLogger("GET", "/internal/test")
	req.Header.Set(services.ServiceKeyHeader, "some-key")
	rr := httptest.NewRecorder()

	mw.RequireServiceAuth(next)(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.False(t, nextCalled)
	sr.AssertNotCalled(t, "Read", mock.Anything, mock.Anything)
}

func TestRequireServiceAuth_MissingServiceKeyHeader(t *testing.T) {
	sr := new(mockSecretReader)
	mw := buildMiddleware(t, sr, time.Minute)

	nextCalled := false
	next := func(w http.ResponseWriter, r *http.Request) { nextCalled = true }

	req := newRequestWithLogger("GET", "/internal/test")
	req.Header.Set(services.ServiceNameHeader, services.Catalog)
	rr := httptest.NewRecorder()

	mw.RequireServiceAuth(next)(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.False(t, nextCalled)
	sr.AssertNotCalled(t, "Read", mock.Anything, mock.Anything)
}

func TestRequireServiceAuth_InvalidServiceName(t *testing.T) {
	sr := new(mockSecretReader)
	mw := buildMiddleware(t, sr, time.Minute)

	nextCalled := false
	next := func(w http.ResponseWriter, r *http.Request) { nextCalled = true }

	req := newRequestWithLogger("GET", "/internal/test")
	req.Header.Set(services.ServiceNameHeader, "bogus-service")
	req.Header.Set(services.ServiceKeyHeader, "some-key")
	rr := httptest.NewRecorder()

	mw.RequireServiceAuth(next)(rr, req)

	assert.Equal(t, http.StatusBadRequest, rr.Code)
	assert.False(t, nextCalled)
	sr.AssertNotCalled(t, "Read", mock.Anything, mock.Anything)
}

// ---------------------------------------------------------------------------
// Key validation tests
// ---------------------------------------------------------------------------

func TestRequireServiceAuth_ValidKey_PassesThrough(t *testing.T) {
	crypto := newTestMiddlewareCrypto(t)
	testKey := "test-api-key-abc123"
	expectedHash := computeKeyHash(t, crypto, testKey)

	sr := new(mockSecretReader)
	vaultPath := services.ServiceAPIKeyPath(services.Catalog)
	sr.On("Read", mock.Anything, vaultPath).Return(&SecretData{
		Data: map[string]interface{}{
			"data": map[string]interface{}{
				"key_hash":     expectedHash,
				"service_name": services.Catalog,
			},
		},
	}, nil)

	mw := &SessionAuthMiddleware{
		crypto:          crypto,
		secretReader:    sr,
		serviceKeyCache: make(map[string]*cachedKeyEntry),
		cacheTTL:        time.Minute,
	}

	var capturedServiceInfo *ServiceInfo
	nextCalled := false
	next := func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		si, err := GetServiceInfoFromContext(r.Context())
		if err == nil {
			capturedServiceInfo = si
		}
		w.WriteHeader(http.StatusOK)
	}

	req := newRequestWithLogger("GET", "/internal/test")
	req.Header.Set(services.ServiceNameHeader, services.Catalog)
	req.Header.Set(services.ServiceKeyHeader, testKey)
	rr := httptest.NewRecorder()

	mw.RequireServiceAuth(next)(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.True(t, nextCalled)
	require.NotNil(t, capturedServiceInfo)
	assert.Equal(t, services.Catalog, capturedServiceInfo.Name)
	sr.AssertCalled(t, "Read", mock.Anything, vaultPath)
}

func TestRequireServiceAuth_InvalidKey_Rejected(t *testing.T) {
	crypto := newTestMiddlewareCrypto(t)

	// Store a hash for a *different* key.
	storedHash := computeKeyHash(t, crypto, "wrong-key")

	sr := new(mockSecretReader)
	vaultPath := services.ServiceAPIKeyPath(services.Catalog)
	sr.On("Read", mock.Anything, vaultPath).Return(&SecretData{
		Data: map[string]interface{}{
			"data": map[string]interface{}{
				"key_hash":     storedHash,
				"service_name": services.Catalog,
			},
		},
	}, nil)

	mw := &SessionAuthMiddleware{
		crypto:          crypto,
		secretReader:    sr,
		serviceKeyCache: make(map[string]*cachedKeyEntry),
		cacheTTL:        time.Minute,
	}

	nextCalled := false
	next := func(w http.ResponseWriter, r *http.Request) { nextCalled = true }

	req := newRequestWithLogger("GET", "/internal/test")
	req.Header.Set(services.ServiceNameHeader, services.Catalog)
	req.Header.Set(services.ServiceKeyHeader, "attacker-key")
	rr := httptest.NewRecorder()

	mw.RequireServiceAuth(next)(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.False(t, nextCalled)
}

func TestRequireServiceAuth_VaultReadError_Rejected(t *testing.T) {
	crypto := newTestMiddlewareCrypto(t)

	sr := new(mockSecretReader)
	vaultPath := services.ServiceAPIKeyPath(services.Catalog)
	sr.On("Read", mock.Anything, vaultPath).Return(nil, assert.AnError)

	mw := &SessionAuthMiddleware{
		crypto:          crypto,
		secretReader:    sr,
		serviceKeyCache: make(map[string]*cachedKeyEntry),
		cacheTTL:        time.Minute,
	}

	nextCalled := false
	next := func(w http.ResponseWriter, r *http.Request) { nextCalled = true }

	req := newRequestWithLogger("GET", "/internal/test")
	req.Header.Set(services.ServiceNameHeader, services.Catalog)
	req.Header.Set(services.ServiceKeyHeader, "some-key")
	rr := httptest.NewRecorder()

	mw.RequireServiceAuth(next)(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.False(t, nextCalled)
}

func TestRequireServiceAuth_SecretNotFound_Rejected(t *testing.T) {
	crypto := newTestMiddlewareCrypto(t)

	sr := new(mockSecretReader)
	vaultPath := services.ServiceAPIKeyPath(services.Catalog)
	sr.On("Read", mock.Anything, vaultPath).Return(nil, nil)

	mw := &SessionAuthMiddleware{
		crypto:          crypto,
		secretReader:    sr,
		serviceKeyCache: make(map[string]*cachedKeyEntry),
		cacheTTL:        time.Minute,
	}

	nextCalled := false
	next := func(w http.ResponseWriter, r *http.Request) { nextCalled = true }

	req := newRequestWithLogger("GET", "/internal/test")
	req.Header.Set(services.ServiceNameHeader, services.Catalog)
	req.Header.Set(services.ServiceKeyHeader, "some-key")
	rr := httptest.NewRecorder()

	mw.RequireServiceAuth(next)(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.False(t, nextCalled)
}

// ---------------------------------------------------------------------------
// Cache tests
// ---------------------------------------------------------------------------

func TestRequireServiceAuth_CacheHit_SkipsVaultRead(t *testing.T) {
	crypto := newTestMiddlewareCrypto(t)
	testKey := "cached-test-key"
	expectedHash := computeKeyHash(t, crypto, testKey)

	// Pre-populate the cache so Vault is never called.
	mw := &SessionAuthMiddleware{
		crypto:       crypto,
		secretReader: new(mockSecretReader), // no expectations → panic if called
		serviceKeyCache: map[string]*cachedKeyEntry{
			services.Catalog: {
				keyHash:   expectedHash,
				fetchedAt: time.Now(),
			},
		},
		cacheTTL: time.Minute,
	}

	nextCalled := false
	next := func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	}

	req := newRequestWithLogger("GET", "/internal/test")
	req.Header.Set(services.ServiceNameHeader, services.Catalog)
	req.Header.Set(services.ServiceKeyHeader, testKey)
	rr := httptest.NewRecorder()

	mw.RequireServiceAuth(next)(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.True(t, nextCalled)
}

func TestRequireServiceAuth_CacheExpired_VaultReadCalled(t *testing.T) {
	crypto := newTestMiddlewareCrypto(t)
	testKey := "expired-test-key"
	expectedHash := computeKeyHash(t, crypto, testKey)

	// Pre-populate the cache with an expired entry.
	sr := new(mockSecretReader)
	mw := &SessionAuthMiddleware{
		crypto:       crypto,
		secretReader: sr,
		serviceKeyCache: map[string]*cachedKeyEntry{
			services.Catalog: {
				keyHash:   expectedHash,
				fetchedAt: time.Now().Add(-2 * time.Minute), // expired
			},
		},
		cacheTTL: time.Minute,
	}

	vaultPath := services.ServiceAPIKeyPath(services.Catalog)
	sr.On("Read", mock.Anything, vaultPath).Return(&SecretData{
		Data: map[string]interface{}{
			"data": map[string]interface{}{
				"key_hash":     expectedHash,
				"service_name": services.Catalog,
			},
		},
	}, nil)

	nextCalled := false
	next := func(w http.ResponseWriter, r *http.Request) {
		nextCalled = true
		w.WriteHeader(http.StatusOK)
	}

	req := newRequestWithLogger("GET", "/internal/test")
	req.Header.Set(services.ServiceNameHeader, services.Catalog)
	req.Header.Set(services.ServiceKeyHeader, testKey)
	rr := httptest.NewRecorder()

	mw.RequireServiceAuth(next)(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.True(t, nextCalled)
	sr.AssertCalled(t, "Read", mock.Anything, vaultPath)
}

func TestRequireServiceAuth_CachePopulatedAfterFirstRequest(t *testing.T) {
	crypto := newTestMiddlewareCrypto(t)
	testKey := "populating-key"
	expectedHash := computeKeyHash(t, crypto, testKey)

	sr := new(mockSecretReader)
	vaultPath := services.ServiceAPIKeyPath(services.Catalog)
	sr.On("Read", mock.Anything, vaultPath).Return(&SecretData{
		Data: map[string]interface{}{
			"data": map[string]interface{}{
				"key_hash":     expectedHash,
				"service_name": services.Catalog,
			},
		},
	}, nil).Once() // only one Vault call expected

	mw := &SessionAuthMiddleware{
		crypto:          crypto,
		secretReader:    sr,
		serviceKeyCache: make(map[string]*cachedKeyEntry),
		cacheTTL:        time.Minute,
	}

	next := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}

	// First request → should hit Vault.
	req1 := newRequestWithLogger("GET", "/internal/test")
	req1.Header.Set(services.ServiceNameHeader, services.Catalog)
	req1.Header.Set(services.ServiceKeyHeader, testKey)
	rr1 := httptest.NewRecorder()
	mw.RequireServiceAuth(next)(rr1, req1)
	assert.Equal(t, http.StatusOK, rr1.Code)

	// Second request → should be served from cache (no additional Vault call).
	req2 := newRequestWithLogger("GET", "/internal/test")
	req2.Header.Set(services.ServiceNameHeader, services.Catalog)
	req2.Header.Set(services.ServiceKeyHeader, testKey)
	rr2 := httptest.NewRecorder()
	mw.RequireServiceAuth(next)(rr2, req2)
	assert.Equal(t, http.StatusOK, rr2.Code)

	sr.AssertNumberOfCalls(t, "Read", 1)
}

func TestRequireServiceAuth_CachedWrongKey_StillRejected(t *testing.T) {
	crypto := newTestMiddlewareCrypto(t)

	// Cache holds the hash for "correct-key".
	correctHash := computeKeyHash(t, crypto, "correct-key")

	mw := &SessionAuthMiddleware{
		crypto:       crypto,
		secretReader: new(mockSecretReader), // shouldn't be called
		serviceKeyCache: map[string]*cachedKeyEntry{
			services.Catalog: {
				keyHash:   correctHash,
				fetchedAt: time.Now(),
			},
		},
		cacheTTL: time.Minute,
	}

	nextCalled := false
	next := func(w http.ResponseWriter, r *http.Request) { nextCalled = true }

	// But the request sends a different key.
	req := newRequestWithLogger("GET", "/internal/test")
	req.Header.Set(services.ServiceNameHeader, services.Catalog)
	req.Header.Set(services.ServiceKeyHeader, "wrong-key")
	rr := httptest.NewRecorder()

	mw.RequireServiceAuth(next)(rr, req)

	assert.Equal(t, http.StatusUnauthorized, rr.Code)
	assert.False(t, nextCalled)
}

// ---------------------------------------------------------------------------
// Cache helper unit tests
// ---------------------------------------------------------------------------

func TestGetCachedKeyHash_EmptyCache(t *testing.T) {
	mw := &SessionAuthMiddleware{
		serviceKeyCache: make(map[string]*cachedKeyEntry),
		cacheTTL:        time.Minute,
	}

	hash, ok := mw.getCachedKeyHash(services.Catalog)
	assert.False(t, ok)
	assert.Empty(t, hash)
}

func TestGetCachedKeyHash_ExpiredEntry(t *testing.T) {
	mw := &SessionAuthMiddleware{
		serviceKeyCache: map[string]*cachedKeyEntry{
			services.Catalog: {
				keyHash:   "old-hash",
				fetchedAt: time.Now().Add(-2 * time.Minute),
			},
		},
		cacheTTL: time.Minute,
	}

	hash, ok := mw.getCachedKeyHash(services.Catalog)
	assert.False(t, ok)
	assert.Empty(t, hash)
}

func TestSetCachedKeyHash(t *testing.T) {
	mw := &SessionAuthMiddleware{
		serviceKeyCache: make(map[string]*cachedKeyEntry),
		cacheTTL:        time.Minute,
	}

	mw.setCachedKeyHash(services.Catalog, "new-hash")

	hash, ok := mw.getCachedKeyHash(services.Catalog)
	assert.True(t, ok)
	assert.Equal(t, "new-hash", hash)
}
