package handler

import (
	"net/http"
)

func AuthorizeUserGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, r *http.Request) {
		// ctx := context.WithValue(r.Context(), "userId", "123")
		next.ServeHTTP(rw, r.WithContext(r.Context()))
	})
}