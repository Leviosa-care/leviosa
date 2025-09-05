package auth

import "net/http"

// GetAccessTokenFromCookies extracts access token from request cookies
func GetAccessTokenFromCookies(r *http.Request) (string, error) {
	cookie, err := r.Cookie(AccessTokenCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}