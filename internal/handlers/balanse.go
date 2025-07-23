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

	currentCents, withdrawnCents, err := s.svc.GetUserBalance(r.Context(), userID)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	// Конвертируем копейки в рубли
	currentRub := float64(currentCents) / 100
	withdrawnRub := float64(withdrawnCents) / 100

	balance := models.BalanceResponse{
		Current:   currentRub,
		Withdrawn: withdrawnRub,
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
		http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
		return
	}
	if req.Sum <= 0 {
		http.Error(w, "Sum must be positive", http.StatusBadRequest)
		return
	}
	if req.Order == "" {
		http.Error(w, "Order must not be empty", http.StatusBadRequest)
		return
	}

	if !isValidLuhn(req.Order) {
		http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
		return
	}

	// Конвертируем рубли в копейки
	sumCents := int64(req.Sum * 100)

	err = s.svc.CreateWithdrawal(r.Context(), userID, req.Order, sumCents)
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
