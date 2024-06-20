package server

import (
	"encoding/json"
	"log"
	"net/http"
	"strings"

	ethparser "eth-parser/eth_parser"
)

type Server struct {
	parser ethparser.Parser
	mux    *http.ServeMux
}

func NewServer(parser ethparser.Parser) *Server {
	s := &Server{
		parser: parser,
		mux:    http.NewServeMux(),
	}
	return s
}

func (s *Server) Start() error {
	s.mux.HandleFunc("/currentBlock", s.handleCurrentBlock)
	s.mux.HandleFunc("/subscribe", s.handleSubscribe)
	s.mux.HandleFunc("/transactions", s.handleTransactions)
	log.Println("Server started on :8080")
	return http.ListenAndServe(":8080", s.mux)
}

func (s *Server) handleCurrentBlock(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	currentBlock := s.parser.GetCurrentBlock()
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]int{"currentBlock": currentBlock})
}

func (s *Server) handleSubscribe(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	address := r.URL.Query().Get("address")
	if !ethparser.ValidateEthAddr(address) {
		http.Error(w, "Invalid address", http.StatusBadRequest)
		return
	}
	success := s.parser.Subscribe(strings.ToLower(address))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]bool{"success": success})
}

func (s *Server) handleTransactions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		return
	}
	address := r.URL.Query().Get("address")
	if !ethparser.ValidateEthAddr(address) {
		http.Error(w, "Invalid address", http.StatusBadRequest)
		return
	}
	transactions := s.parser.GetTransactions(strings.ToLower(address))
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string][]ethparser.Transaction{"transactions": transactions})
}
