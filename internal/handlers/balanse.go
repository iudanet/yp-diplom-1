package handlers

import (
	"encoding/json"
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
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Implement balance withdrawal logic here
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Balance withdrawn"))
}
