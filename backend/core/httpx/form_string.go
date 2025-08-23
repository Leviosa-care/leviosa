package httpx

import (
	"net/http"
	"strings"
)

func FormString(r *http.Request, key string) string {
	return strings.TrimSpace(r.FormValue(key))
}
