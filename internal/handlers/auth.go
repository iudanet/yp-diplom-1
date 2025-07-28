package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/iudanet/yp-diplom-1/internal/models"
	"github.com/iudanet/yp-diplom-1/internal/pkg/auth"
)

// Register godoc
//
//	@Summary		Register a new user
//	@Description	Register a new user with login and password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.RegisterRequest	true	"Register request"
//	@Success		200		{string}	string					"Successfully registered"
//	@Failure		400		{string}	string					"Invalid request format"
//	@Failure		409		{string}	string					"Login already exists"
//	@Failure		500		{string}	string					"Internal server error"
//	@Router			/api/user/register [post]
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

// Login godoc
//
//	@Summary		Login user
//	@Description	Authenticate user with login and password
//	@Tags			auth
//	@Accept			json
//	@Produce		json
//	@Param			request	body		models.LoginRequest	true	"Login request"
//	@Success		200		{string}	string				"Successfully logged in"
//	@Failure		400		{string}	string				"Invalid request format"
//	@Failure		401		{string}	string				"Invalid credentials"
//	@Failure		500		{string}	string				"Internal server error"
//	@Router			/api/user/login [post]
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
