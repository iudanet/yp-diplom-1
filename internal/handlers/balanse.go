package handlers

import "net/http"

func (s *Server) Balance(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Implement balance logic here
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Balance retrieved"))
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
