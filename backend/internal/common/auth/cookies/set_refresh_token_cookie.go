package cookies

import (
	"net/http"
	"time"
)

// SetRefreshTokenCookie sets only the refresh token cookie (for refresh operations with rotation)
func SetRefreshTokenCookie(w http.ResponseWriter, refreshToken string, expiry time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  expiry,
	})
}

