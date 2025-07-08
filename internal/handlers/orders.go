package handlers

import "net/http"

func (s *Server) PostOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Implement login logic here
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Create Orders"))
}

func (s *Server) GetOrders(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	// Implement login logic here
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Orders None"))
}
