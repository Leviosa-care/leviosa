package cookies

import "net/http"

// ClearTokenCookies removes both access and refresh token cookies
func ClearTokenCookies(w http.ResponseWriter) {
	// Clear access token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     AccessTokenCookieName,
		Value:    "",
		Path:     "/",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1, // Delete immediately
	})

	// Clear refresh token cookie
	http.SetCookie(w, &http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    "",
		Path:     RefreshEndpoint,
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1, // Delete immediately
	})
}
