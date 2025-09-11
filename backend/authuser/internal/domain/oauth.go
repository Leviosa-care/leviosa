package domain

// OAuthStartRequest represents the request to start OAuth flow
type OAuthStartRequest struct {
	Provider string `json:"provider" validate:"required,oneof=google apple"`
}

// OAuthStartResponse contains the authorization URL to redirect to
type OAuthStartResponse struct {
	AuthorizationURL string `json:"authorization_url"`
	State            string `json:"state"`
}

// OAuthCallbackRequest represents the OAuth callback with authorization code
type OAuthCallbackRequest struct {
	Provider string `json:"provider" validate:"required,oneof=google apple"`
	Code     string `json:"code" validate:"required"`
	State    string `json:"state" validate:"required"`
}

// OAuthCallbackResponse contains the session information after successful OAuth
type OAuthCallbackResponse struct {
	AccessToken        string `json:"access_token"`
	RefreshToken       string `json:"refresh_token"`
	AccessTokenExpiry  int64  `json:"access_token_expiry"`
	RefreshTokenExpiry int64  `json:"refresh_token_expiry"`
	IsNewUser          bool   `json:"is_new_user"`
}