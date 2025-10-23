# Authentication Test Utilities - Usage Guide

This guide demonstrates how to use the authentication test utilities in `core/testutils/auth.go` for integration testing across all Leviosa microservices.

## Overview

The auth test utilities provide:
- ✅ Role-based user setup (6 roles: Visitor, Standard, Premium, Guest, Partner, Administrator)
- ✅ ENCX-encrypted session management
- ✅ HTTP request helpers for middleware testing
- ✅ Support for expired/invalid session scenarios
- ✅ Granular cleanup utilities
- ✅ Cross-service compatibility

## Quick Start

```go
// Setup auth context for your tests
authCtx := &testutils.AuthTestContext{
    Pool:   testPool,    // Your test PostgreSQL pool
    Redis:  testClient,  // Your test Redis client
    Crypto: crypto,      // Your ENCX crypto service
}

// Create a user with specific role
accessToken := testutils.SetupAdminUser(t, ctx, authCtx)

// Create authenticated HTTP request
req := testutils.CreateAuthenticatedRequest("GET", "/api/settings", accessToken)

// Test your endpoint
resp := httptest.NewRecorder()
yourHandler.ServeHTTP(resp, req)

// Verify response
assert.Equal(t, http.StatusOK, resp.Code)

// Cleanup after test
testutils.ClearAuthData(t, ctx, authCtx)
```

## Role-Based User Setup

### Basic Role Functions

```go
// Individual role setup
visitorToken := testutils.SetupVisitorUser(t, ctx, authCtx)
standardToken := testutils.SetupStandardUser(t, ctx, authCtx)
premiumToken := testutils.SetupPremiumUser(t, ctx, authCtx)
guestToken := testutils.SetupGuestUser(t, ctx, authCtx)
partnerToken := testutils.SetupPartnerUser(t, ctx, authCtx)
adminToken := testutils.SetupAdminUser(t, ctx, authCtx)

// Generic role setup
token := testutils.SetupUserWithRole(t, ctx, identity.Premium, authCtx)
```

### Multiple Users Setup

```go
// Setup users for multiple roles at once
roles := []identity.Role{
    identity.Standard,
    identity.Premium,
    identity.Administrator,
}
tokens := testutils.SetupMultipleUsers(t, ctx, roles, authCtx)

// Access tokens by role
standardToken := tokens[identity.Standard]
premiumToken := tokens[identity.Premium]
adminToken := tokens[identity.Administrator]
```

### Custom User Data

```go
// Create user with custom data
token := testutils.SetupUserWithCustomData(
    t, ctx,
    identity.Premium,                    // Role
    "john.doe@example.com",             // Email
    "John",                             // First name
    "Doe",                              // Last name
    "+33612345678",                     // Phone (optional)
    authCtx,
)
```

## Session Testing Scenarios

### Active Sessions (Normal Flow)

```go
// Standard active user session
token := testutils.SetupPremiumUser(t, ctx, authCtx)

// Request will pass through auth middleware successfully
req := testutils.CreateAuthenticatedRequest("GET", "/api/catalog/products", token)
```

### Pending Sessions (Registration Flow)

```go
// User with pending session (email verification required)
token := testutils.SetupPendingUserWithRole(t, ctx, identity.Standard, authCtx)

// Request should be handled appropriately by your middleware
req := testutils.CreateAuthenticatedRequest("GET", "/api/user/profile", token)
```

### Expired Sessions (Timeout Testing)

```go
// User with expired session
token := testutils.SetupExpiredUserWithRole(t, ctx, identity.Standard, authCtx)

// Request should be rejected by auth middleware
req := testutils.CreateAuthenticatedRequest("GET", "/api/protected", token)

// Test that your middleware returns 401 Unauthorized
resp := httptest.NewRecorder()
yourAuthMiddleware(yourHandler).ServeHTTP(resp, req)
assert.Equal(t, http.StatusUnauthorized, resp.Code)
```

## HTTP Request Helpers

### Basic Authenticated Requests

```go
// Simple GET request with auth
req := testutils.CreateAuthenticatedRequest("GET", "/api/settings", token)

// POST request with JSON body
body := strings.NewReader(`{"name": "Test Setting"}`)
req := testutils.CreateAuthenticatedRequestWithBody("POST", "/api/settings", token, body)
```

### Custom Headers and Cookies

```go
// Create Authorization header only
headers := testutils.CreateAuthHeader(token)
req := httptest.NewRequest("GET", "/api/data", nil)
req.Header = headers

// Create auth cookie only
cookie := testutils.CreateAuthCookie(token)
req.AddCookie(cookie)

// Create refresh cookie for testing refresh endpoints
refreshCookie := testutils.CreateRefreshCookie(refreshToken)
req.AddCookie(refreshCookie)
```

## Role-Based Authorization Testing

### Test Access Control

```go
func TestAdminOnlyEndpoint(t *testing.T) {
    ctx := context.Background()

    // Test with admin user (should succeed)
    adminToken := testutils.SetupAdminUser(t, ctx, authCtx)
    req := testutils.CreateAuthenticatedRequest("GET", "/api/admin/users", adminToken)
    resp := httptest.NewRecorder()
    adminOnlyHandler.ServeHTTP(resp, req)
    assert.Equal(t, http.StatusOK, resp.Code)

    // Clean up for next test
    testutils.ClearSessionsOnly(t, ctx, authCtx)

    // Test with standard user (should fail)
    standardToken := testutils.SetupStandardUser(t, ctx, authCtx)
    req = testutils.CreateAuthenticatedRequest("GET", "/api/admin/users", standardToken)
    resp = httptest.NewRecorder()
    adminOnlyHandler.ServeHTTP(resp, req)
    assert.Equal(t, http.StatusForbidden, resp.Code)
}
```

### Test Role Hierarchy

```go
func TestRoleHierarchy(t *testing.T) {
    ctx := context.Background()

    // Test minimum role requirements
    testCases := []struct {
        name         string
        role         identity.Role
        setupFunc    func(*testing.T, context.Context, *testutils.AuthTestContext) string
        expectStatus int
    }{
        {"Visitor access", identity.Visitor, testutils.SetupVisitorUser, http.StatusOK},
        {"Standard access", identity.Standard, testutils.SetupStandardUser, http.StatusOK},
        {"Premium access", identity.Premium, testutils.SetupPremiumUser, http.StatusOK},
        {"Guest access", identity.Guest, testutils.SetupGuestUser, http.StatusForbidden},
    }

    for _, tc := range testCases {
        t.Run(tc.name, func(t *testing.T) {
            token := tc.setupFunc(t, ctx, authCtx)
            req := testutils.CreateAuthenticatedRequest("GET", "/api/premium-content", token)
            resp := httptest.NewRecorder()
            premiumHandler.ServeHTTP(resp, req)
            assert.Equal(t, tc.expectStatus, resp.Code)

            testutils.ClearSessionsOnly(t, ctx, authCtx)
        })
    }
}
```

## Test Data Verification

### Verify User Creation

```go
// Check that user was created correctly
token := testutils.SetupPremiumUser(t, ctx, authCtx)
userCount := testutils.CountAuthUsers(t, ctx, authCtx.Pool)
assert.Equal(t, 1, userCount)

// Verify specific user exists
exists := testutils.UserExists(t, ctx, "premium@leviosa.care", authCtx.Pool, authCtx.Crypto)
assert.True(t, exists)
```

### Verify Session Creation

```go
// Check that session was created
token := testutils.SetupAdminUser(t, ctx, authCtx)
sessionCount := testutils.CountActiveSessions(t, ctx, authCtx.Redis)
assert.Equal(t, 1, sessionCount)
```

## Cleanup Strategies

### Complete Cleanup

```go
// Clean up everything (users + sessions)
defer testutils.ClearAuthData(t, ctx, authCtx)
```

### Granular Cleanup

```go
// Keep users, clear sessions only (useful for multiple session tests)
testutils.ClearSessionsOnly(t, ctx, authCtx)

// Keep sessions, clear users only (useful for session cleanup tests)
testutils.ClearUsersOnly(t, ctx, authCtx)
```

### Test Isolation

```go
func TestMyFeature(t *testing.T) {
    ctx := context.Background()

    // Clean up before test
    testutils.ClearAuthData(t, ctx, authCtx)

    // Your test logic here
    token := testutils.SetupStandardUser(t, ctx, authCtx)
    // ... test implementation

    // Clean up after test
    testutils.ClearAuthData(t, ctx, authCtx)
}
```

## Integration with Different Services

### Catalog Service Example

```go
func TestCatalogProductAccess(t *testing.T) {
    authCtx := &testutils.AuthTestContext{
        Pool:   catalogTestPool,
        Redis:  catalogTestRedis,
        Crypto: catalogCrypto,
    }

    // Test public access (no auth required)
    req := httptest.NewRequest("GET", "/api/catalog/products", nil)
    resp := httptest.NewRecorder()
    catalogHandler.ServeHTTP(resp, req)
    assert.Equal(t, http.StatusOK, resp.Code)

    // Test premium content access
    premiumToken := testutils.SetupPremiumUser(t, ctx, authCtx)
    req = testutils.CreateAuthenticatedRequest("GET", "/api/catalog/premium-products", premiumToken)
    resp = httptest.NewRecorder()
    catalogHandler.ServeHTTP(resp, req)
    assert.Equal(t, http.StatusOK, resp.Code)

    // Test standard user blocked from premium content
    standardToken := testutils.SetupStandardUser(t, ctx, authCtx)
    req = testutils.CreateAuthenticatedRequest("GET", "/api/catalog/premium-products", standardToken)
    resp = httptest.NewRecorder()
    catalogHandler.ServeHTTP(resp, req)
    assert.Equal(t, http.StatusForbidden, resp.Code)
}
```

### Settings Service Example

```go
func TestSettingsManagement(t *testing.T) {
    authCtx := &testutils.AuthTestContext{
        Pool:   settingsTestPool,
        Redis:  settingsTestRedis,
        Crypto: settingsCrypto,
    }

    // Test admin can update company settings
    adminToken := testutils.SetupAdminUser(t, ctx, authCtx)
    body := strings.NewReader(`{"company_name": "New Company Name"}`)
    req := testutils.CreateAuthenticatedRequestWithBody("PUT", "/api/settings/company", adminToken, body)
    resp := httptest.NewRecorder()
    settingsHandler.ServeHTTP(resp, req)
    assert.Equal(t, http.StatusOK, resp.Code)

    // Test standard user cannot update company settings
    standardToken := testutils.SetupStandardUser(t, ctx, authCtx)
    req = testutils.CreateAuthenticatedRequestWithBody("PUT", "/api/settings/company", standardToken, body)
    resp = httptest.NewRecorder()
    settingsHandler.ServeHTTP(resp, req)
    assert.Equal(t, http.StatusForbidden, resp.Code)
}
```

## Best Practices

1. **Always clean up** after tests to ensure test isolation
2. **Use granular cleanup** when testing multiple scenarios with same users
3. **Test both success and failure** scenarios for each role
4. **Verify database state** in addition to HTTP responses
5. **Use descriptive test names** that indicate the role and expected outcome
6. **Test edge cases** like expired sessions, invalid tokens, etc.
7. **Keep test data consistent** across different service tests

## Error Handling

The utilities will automatically fail tests with clear error messages if:
- User creation fails (database errors, encryption errors)
- Session creation fails (Redis errors, token generation errors)
- HTTP request creation fails (invalid inputs)

All error messages include context about what failed and why, making debugging easier.