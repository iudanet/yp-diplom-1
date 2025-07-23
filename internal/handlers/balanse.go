package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"github.com/iudanet/yp-diplom-1/internal/models"
)

func (s *Server) Balance(w http.ResponseWriter, r *http.Request) {
	userID, err := s.checkAuth(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	current, withdrawn, err := s.svc.GetUserBalance(r.Context(), userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	balance := models.BalanceResponse{
		Current:   current,
		Withdrawn: withdrawn,
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(balance); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}

func (s *Server) BalanceWithdraw(w http.ResponseWriter, r *http.Request) {
	userID, err := s.checkAuth(r)
	if err != nil {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req models.BalanceWithdrawRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest) // Добавим информацию об ошибке
		return
	}

	if req.Order == "" || req.Sum <= 0 {
		http.Error(w, "Order must not be empty and sum must be positive", http.StatusBadRequest)
		return
	}

	// Проверка номера заказа по алгоритму Луна
	if !isValidLuhn(req.Order) {
		http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
		return
	}

	err = s.svc.CreateWithdrawal(r.Context(), userID, req.Order, req.Sum) // Убираем преобразование типа
	switch {
	case errors.Is(err, models.ErrInsufficientFunds):
		http.Error(w, "Insufficient funds", http.StatusPaymentRequired)
	case err != nil:
		log.Printf("Failed to create withdrawal: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	default:
		w.WriteHeader(http.StatusOK)
	}
}
