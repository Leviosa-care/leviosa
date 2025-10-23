package cookies

import "net/http"

// GetRefreshTokenFromCookies extracts refresh token from request cookies
func GetRefreshTokenFromCookies(r *http.Request) (string, error) {
	cookie, err := r.Cookie(RefreshTokenCookieName)
	if err != nil {
		return "", err
	}
	return cookie.Value, nil
}

