package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/iudanet/yp-diplom-1/internal/models"
)

func (s *Server) Withdrawals(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int64)

	withdrawals, err := s.svc.GetWithdrawals(r.Context(), userID)
	if err != nil {
		log.Printf("Failed to get withdrawals: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	if len(withdrawals) == 0 {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	// Конвертируем копейки в рубли
	response := make([]models.WithdrawalResponse, len(withdrawals))
	for i, w := range withdrawals {
		response[i] = models.WithdrawalResponse{
			Order:       w.Order,
			Sum:         float64(w.Sum) / 100,
			ProcessedAt: w.ProcessedAt,
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
