package middleware

import "net/http"

// here I just have the auth thing to check if the user has the right role for this action

type Handler = func(w http.ResponseWriter, r *http.Request)

// NOTE: here is how I am going to do the middleware, I think
func Somefunc(next http.Handler) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		next.ServeHTTP(w, r)
	}
}
