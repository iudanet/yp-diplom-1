package handlers

import "net/http"

func (s *Server) Register(w http.ResponseWriter, r *http.Request) {
	// Implement authentication logic here
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Registration successful"))
}

func (s *Server) Login(w http.ResponseWriter, r *http.Request) {
	// Implement registration logic here
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Login successful"))
}
