package main

import (
	"encoding/json"
	"net/http"

	"github.com/maryelinv/go-chat-challenge/internal/auth"
	"github.com/maryelinv/go-chat-challenge/internal/db"
	"gorm.io/gorm"
)

type registerReq struct{ Username, Password string }
type loginReq struct{ Username, Password string }

func register(g *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in registerReq
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		if in.Username == "" || len(in.Password) < 6 {
			http.Error(w, "username and 6+ char password required", http.StatusBadRequest)
			return
		}
		hash, _ := auth.Hash(in.Password)
		u := db.User{Username: in.Username, PasswordHash: hash}
		if err := g.Create(&u).Error; err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		_ = auth.SetSession(w, u.ID)
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(map[string]any{"id": u.ID, "username": u.Username})
	}
}

func login(g *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var in loginReq
		if err := json.NewDecoder(r.Body).Decode(&in); err != nil {
			http.Error(w, "bad json", http.StatusBadRequest)
			return
		}
		var u db.User
		if err := g.Where("username = ?", in.Username).First(&u).Error; err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		if err := auth.Check(u.PasswordHash, in.Password); err != nil {
			http.Error(w, "invalid credentials", http.StatusUnauthorized)
			return
		}
		_ = auth.SetSession(w, u.ID)
		json.NewEncoder(w).Encode(map[string]any{"id": u.ID, "username": u.Username})
	}
}

func logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth.ClearSession(w, r)
		w.WriteHeader(http.StatusNoContent)
	}
}

func me(g *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		uid, ok := auth.CurrentUserID(r)
		if !ok {
			http.Error(w, "unauthenticated", http.StatusUnauthorized)
			return
		}
		var u db.User
		if err := g.First(&u, uid).Error; err != nil {
			http.Error(w, "not found", 404)
			return
		}
		json.NewEncoder(w).Encode(map[string]any{"id": u.ID, "username": u.Username})
	}
}
