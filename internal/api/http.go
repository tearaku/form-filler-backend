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
	switch r.Method {
	case http.MethodGet:
		log.Println("Handling GET request...")
		eventId := r.URL.Path[len("/event/"):]
		// Fetch relevant data from database based on event ID
		eventInfo, err := server.schoolForm.GetEventInfo(r.Context(), schoolForm.GetEventInfoParams{EventID: eventId})
		if err != nil {
			log.Println(err)
			http.Error(w, "Error: fetching event data from db failed.", http.StatusInternalServerError)
		}
		log.Println(eventInfo.Attendants)
		// Create temp file for the excel
		// Modify the excel file
		// Export to PDF via gotenburg (hosted elsewhere? Dockerfile)
	default:
		http.Error(w, "Error: http method not allowed.", http.StatusMethodNotAllowed)
	}
}
