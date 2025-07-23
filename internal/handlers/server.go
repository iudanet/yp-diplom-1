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

	// Создаем подмаршрут для защищенных endpoint'ов
	protected := http.NewServeMux()
	protected.HandleFunc("POST /api/user/orders", s.PostOrders)
	protected.HandleFunc("GET /api/user/orders", s.GetOrders)
	protected.HandleFunc("GET /api/user/balance", s.Balance)
	protected.HandleFunc("POST /api/user/balance/withdraw", s.BalanceWithdraw)
	protected.HandleFunc("GET /api/user/withdrawals", s.Withdrawals)

	// Применяем middleware к защищенным маршрутам
	mux.Handle("/", s.authMiddleware(protected))

	return mux
}
