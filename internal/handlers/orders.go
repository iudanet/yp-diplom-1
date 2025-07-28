package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/iudanet/yp-diplom-1/internal/models"
	"github.com/iudanet/yp-diplom-1/internal/pkg/luhn"
)

// PostOrders godoc
//
//	@Summary		Upload order number
//	@Description	Upload order number for accrual calculation
//	@Tags			orders
//	@Accept			text/plain
//	@Produce		json
//	@Security		BearerAuth
//	@Param			order	body		string	true	"Order number"
//	@Success		200		{string}	string	"Order already uploaded"
//	@Success		202		{string}	string	"Order accepted for processing"
//	@Failure		400		{string}	string	"Invalid request format"
//	@Failure		401		{string}	string	"Unauthorized"
//	@Failure		409		{string}	string	"Order already uploaded by another user"
//	@Failure		422		{string}	string	"Invalid order number format"
//	@Failure		500		{string}	string	"Internal server error"
//	@Router			/api/user/orders [post]
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
	if !luhn.IsValid(number) {
		http.Error(w, "Invalid order number", http.StatusUnprocessableEntity)
		return
	}

	// Получаем userID из контекста
	userID := r.Context().Value(userIDKey).(int64)

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

// GetOrders godoc
//
//	@Summary		Get user's orders
//	@Description	Get list of orders uploaded by user
//	@Tags			orders
//	@Accept			json
//	@Produce		json
//	@Security		BearerAuth
//	@Success		200	{array}		models.OrderUser	"List of user's orders"
//	@Success		204	{string}	string				"No orders found"
//	@Failure		401	{string}	string				"Unauthorized"
//	@Failure		500	{string}	string				"Internal server error"
//	@Router			/api/user/orders [get]
func (s *Server) GetOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value(userIDKey).(int64)

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
		log.Printf("Encode error: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
	}
}
