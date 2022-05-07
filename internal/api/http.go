package api

// Various HTTP REST endpoint handlers

import (
	"log"
	"net/http"
	"teacup1592/form-filler/internal/schoolForm"
)

func NewHTTPServer(sF *schoolForm.Service) http.Handler {
	server := &HTTPServer{
		schoolForm: sF,
		mux:        http.NewServeMux(),
	}
	server.mux.HandleFunc("/event/", server.handleGetEventInfo)
	return server.mux
}

// HTTPServer exposes inventory.Service via HTTP.
type HTTPServer struct {
	schoolForm *schoolForm.Service
	mux        *http.ServeMux
}

func (server *HTTPServer) handleGetEventInfo(w http.ResponseWriter, r *http.Request) {
	// TODO
	log.Println("A get request came in")
}
