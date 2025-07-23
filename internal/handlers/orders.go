package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/iudanet/yp-diplom-1/internal/models"
)

func (s *Server) PostOrders(w http.ResponseWriter, r *http.Request) {
	// Проверяем Content-Type
	if r.Header.Get("Content-Type") != "text/plain" {
		http.Error(w, "Content-Type must be text/plain", http.StatusBadRequest)
		return
	}

	// Читаем тело запроса
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	number := string(body)
	if number == "" {
		http.Error(w, "Empty order number", http.StatusBadRequest)
		return
	}

	// Проверяем номер заказа по алгоритму Луна
	if !isValidLuhn(number) {
		http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
		return
	}

	// Проверяем авторизацию пользователя
	userID, err := s.checkAuth(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	err = s.svc.CreateOrder(r.Context(), userID, number)
	switch {
	case errors.Is(err, models.ErrOrderAlreadyUploaded):
		w.WriteHeader(http.StatusOK)
	case errors.Is(err, models.ErrOrderAlreadyUploadedByAnotherUser):
		http.Error(w, "Order already uploaded by another user", http.StatusConflict)
	case errors.Is(err, models.ErrInvalidOrderNumber):
		http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
	case err != nil:
		log.Printf("Error creating order: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusAccepted)
	}
}

func (s *Server) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID, err := s.checkAuth(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	orders, err := s.svc.GetOrders(r.Context(), userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(orders) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(orders); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
