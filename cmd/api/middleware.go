package main

import (
	"net/http"

	"github.com/maryelinv/go-chat-challenge/internal/auth"
)

func requireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, ok := auth.CurrentUserID(r); !ok {
			http.Error(w, "unauthenticated", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
