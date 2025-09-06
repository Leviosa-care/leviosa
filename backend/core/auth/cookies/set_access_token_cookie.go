package cookies

import (
	"net/http"
	"time"
)

// SetAccessTokenCookie sets only the access token cookie (for refresh operations)
func SetAccessTokenCookie(w http.ResponseWriter, accessToken string, expiry time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     AccessTokenCookieName,
		Value:    accessToken,
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  expiry,
	})
}