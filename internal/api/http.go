package api

// Various HTTP REST endpoint handlers

import (
	"log"
	"net/http"
	"teacup1592/form-filler/internal/schoolForm"
)

func NewHTTPServer(sF *schoolForm.Service) http.Handler {
	s := &HTTPServer{
		schoolForm: sF,
		mux:        http.NewServeMux(),
	}
	s.mux.HandleFunc("/event/", s.handleGetEventInfo)
	return s.mux
}

// HTTPServer exposes inventory.Service via HTTP.
type HTTPServer struct {
	schoolForm *schoolForm.Service
	mux        *http.ServeMux
}

func (s *HTTPServer) handleGetEventInfo(w http.ResponseWriter, r *http.Request) {
	log.Printf("Request from: %s\n", r.URL.Host)
	switch r.Method {
	case http.MethodGet:
		eventId := r.URL.Path[len("/event/"):]
		// Fetch relevant data from database based on event ID
		eventInfo, err := s.schoolForm.GetEventInfo(r.Context(), schoolForm.GetEventInfoParams{EventID: eventId})
		if err != nil {
			log.Println(err)
			http.Error(w, "Error: fetching event data from db failed.", http.StatusInternalServerError)
		}
		err = s.schoolForm.FetchAttendances(r.Context(), schoolForm.FetchAttendancesParams{
			EventInfo: eventInfo,
		})
		if err != nil {
			log.Println(err)
			http.Error(w, "Error: fetching event data from db failed.", http.StatusInternalServerError)
		}
		// Begin modifying excel data
		// TODO: how to return the final file...? as pointer? XD
		if err = s.schoolForm.WriteExcel(eventInfo); err != nil {
			log.Println(err)
			http.Error(w, "Error: writing data to excel file failed.", http.StatusInternalServerError)
		}
		// Export to PDF via gotenburg (hosted elsewhere? Dockerfile)
	default:
		http.Error(w, "Error: http method not allowed.", http.StatusMethodNotAllowed)
	}
}
