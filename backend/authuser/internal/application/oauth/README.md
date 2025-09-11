# OAuth Setup Instructions

## Overview
This package provides OAuth authentication using the Goth library for Google and Apple providers.

## Initialization
To initialize OAuth providers in your application, call the `InitializeOAuthProviders()` function during application startup:

```go
import "github.com/Leviosa-care/authuser/internal/application/oauth"

func main() {
    // Initialize OAuth providers
    if err := oauth.InitializeOAuthProviders(); err != nil {
        log.Fatalf("Failed to initialize OAuth providers: %v", err)
    }
    
    // Continue with application setup...
}
```

## Required Environment Variables

### Google OAuth
```bash
GOOGLE_CLIENT_ID=your_google_client_id
GOOGLE_CLIENT_SECRET=your_google_client_secret
```

### Apple OAuth (Optional)
```bash
APPLE_CLIENT_ID=your_apple_client_id
APPLE_CLIENT_SECRET=your_apple_client_secret
APPLE_TEAM_ID=your_apple_team_id
APPLE_KEY_ID=your_apple_key_id
APPLE_PRIVATE_KEY=your_apple_private_key
```

### Required for Both
```bash
SESSION_SECRET=your_session_secret_key_32_chars_minimum
BASE_URL=http://localhost:8080  # or your production URL
ENV=development  # or production
```

## OAuth Flow
1. **Start OAuth**: GET `/auth/oauth/{provider}` → returns authorization URL
2. **OAuth Callback**: GET `/auth/oauth/{provider}/callback` → handles provider callback and creates/logs in user
3. **Session Management**: Uses existing session token system with HTTP-only cookies

## Integration with Existing Services
The OAuth implementation integrates with:
- **User Service**: Creates or links OAuth users
- **Session Service**: Issues access/refresh tokens
- **Encryption Service**: Encrypts sensitive user data (GDPR compliance)
- **Stripe Service**: Creates customer accounts for new OAuth users

## Testing
Add OAuth provider initialization to your test setup:

```go
func TestMain(m *testing.M) {
    // ... other setup ...
    
    // Initialize OAuth for tests
    if err := oauth.InitializeOAuthProviders(); err != nil {
        log.Fatalf("Failed to initialize OAuth providers for tests: %v", err)
    }
    
    // ... continue with test setup ...
}
```
