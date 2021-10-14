package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ONSdigital/dp-find-insights-poc-api/api"
)

type Server struct {
}

func New() *Server {
	return &Server{}
}

func (svr *Server) GetHello(w http.ResponseWriter, r *http.Request) {
	response := &api.HelloResponse{
		Message: "Hello, world!",
	}
	sendResult(w, http.StatusOK, response)
}

func sendResult(w http.ResponseWriter, code int, result interface{}) {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(result)
	if err != nil {
		log.Printf("SendResult: %s", err)
	}
}

func sendError(w http.ResponseWriter, code int, message string) {
	e := api.Error{
		Error: message,
	}
	w.WriteHeader(code)
	err := json.NewEncoder(w).Encode(e)
	if err != nil {
		log.Printf("SendResult: %s", err)
	}
}
