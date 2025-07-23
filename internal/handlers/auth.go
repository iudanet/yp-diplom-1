package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/iudanet/yp-diplom-1/internal/models"
	"github.com/iudanet/yp-diplom-1/internal/pkg/auth"
)

func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	var req models.RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if req.Login == "" || req.Password == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err := s.svc.Register(r.Context(), req.Login, req.Password)
	if err != nil {
		if errors.Is(err, models.ErrUserAlreadyExists) {
			w.WriteHeader(http.StatusConflict)
			return
		}

		log.Printf("Failed to register user: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	user, err := s.svc.Login(r.Context(), req.Login, req.Password)
	if err != nil {
		log.Printf("Failed to login after registration: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	token, err := auth.GenerateToken(user.ID, s.cfg.SecretKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Устанавливаем токен в заголовок Authorization
	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	var req models.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	user, err := s.svc.Login(r.Context(), req.Login, req.Password)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	token, err := auth.GenerateToken(user.ID, s.cfg.SecretKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Устанавливаем токен в заголовок Authorization
	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}
