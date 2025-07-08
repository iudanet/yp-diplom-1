package handlers

import "net/http"

func (s *Server) Withdrawals(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Implement withdrawals logic here
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Withdrawals processed"))
}
