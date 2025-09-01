package auth

import (
	"net/http"
	"time"
)

const (
	AccessTokenCookieName  = "leviosa_access_token"
	RefreshTokenCookieName = "leviosa_refresh_token"
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
		Path:     "/auth/refresh", // Restrict to refresh endpoint only
		HttpOnly: true,
		Secure:   true, // Only sent over HTTPS
		SameSite: http.SameSiteStrictMode,
		Expires:  refreshExpiry,
	})
}

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
		Path:     "/auth/refresh",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		MaxAge:   -1, // Delete immediately
	})
}

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

// SetRefreshTokenCookie sets only the refresh token cookie (for refresh operations with rotation)
func SetRefreshTokenCookie(w http.ResponseWriter, refreshToken string, expiry time.Time) {
	http.SetCookie(w, &http.Cookie{
		Name:     RefreshTokenCookieName,
		Value:    refreshToken,
		Path:     "/auth/refresh",
		HttpOnly: true,
		Secure:   true,
		SameSite: http.SameSiteStrictMode,
		Expires:  expiry,
	})
}

// GetAccessTokenFromCookies extracts access token from request cookies
func GetAccessTokenFromCookies(r *http.Request) (string, error) {
	cookie, err := r.Cookie(AccessTokenCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

// GetRefreshTokenFromCookies extracts refresh token from request cookies
func GetRefreshTokenFromCookies(r *http.Request) (string, error) {
	cookie, err := r.Cookie(RefreshTokenCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}
