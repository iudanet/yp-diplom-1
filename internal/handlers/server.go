package handlers

import (
	"net/http"

	"github.com/iudanet/yp-diplom-1/internal/config"
	"github.com/iudanet/yp-diplom-1/internal/service"
)

type Server struct {
	svc service.Service
	cfg *config.Config
}

func New(svc service.Service, cfg *config.Config) *Server {
	return &Server{svc: svc, cfg: cfg}
}

func (s *Server) NewMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/user/register", s.Register)
	mux.HandleFunc("POST /api/user/login", s.Login)

	// Защищенные маршруты
	mux.HandleFunc("POST /api/user/orders", s.PostOrders)
	mux.HandleFunc("GET /api/user/orders", s.GetOrders)

	mux.HandleFunc("GET /api/user/balance", s.Balance)
	mux.HandleFunc("POST /api/user/balance/withdraw", s.BalanceWithdraw)
	mux.HandleFunc("GET /api/user/withdrawals", s.Withdrawals)

	return mux
}
