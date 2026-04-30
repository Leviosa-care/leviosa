package cookies

import (
	"net/http"
	"time"
)

// SetTokenCookies sets both access and refresh token cookies with appropriate security settings
func SetTokenCookies(w http.ResponseWriter, accessToken, refreshToken string, accessExpiry, refreshExpiry time.Time) {
	// Set access token cookie (available on all paths)
	http.SetCookie(w, &http.Cookie{
		Name:     AccessTokenCookieName,
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // Only sent over HTTPS
		SameSite: http.SameSiteStrictMode,
		Expires:  accessExpiry,
	})

	// Set refresh token cookie (restricted to refresh endpoint)
	http.SetCookie(w, &http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    refreshToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true, // Only sent over HTTPS
		SameSite: http.SameSiteStrictMode,
		Expires:  refreshExpiry,
	})
}
