package api

// Various HTTP REST endpoint handlers

import (
	"archive/zip"
	"log"
	"net/http"
	"os"
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
		clubLeader, err := s.schoolForm.GetMemberByDept(r.Context(), schoolForm.GetMemberByDeptParams{
			Description: "社長",
		})
		if err != nil {
			log.Println(err)
			http.Error(w, "Error: fetching club leader from db failed.", http.StatusInternalServerError)
		}
		// Begin modifying excel data, and adding to zip file
		zF, err := os.CreateTemp("", "event_"+eventId+"_*.zip")
		zFName := zF.Name()
		if err != nil {
			http.Error(w, "Error: failed to create zip file.", http.StatusInternalServerError)
		}
		defer os.Remove(zF.Name())
		zipWriter := zip.NewWriter(zF)
		if err = s.schoolForm.WriteSchForm(eventInfo, clubLeader, zipWriter); err != nil {
			log.Println(err)
			http.Error(w, "Error: writing data to school form failed.", http.StatusInternalServerError)
		}
		if err = s.schoolForm.WriteInsuranceForm(eventInfo, zipWriter); err != nil {
			log.Println(err)
			http.Error(w, "Error: writing data to insurance form failed.", http.StatusInternalServerError)
		}
		if err = s.schoolForm.WriteMountPass(eventInfo, zipWriter); err != nil {
			log.Println(err)
			http.Error(w, "Error: writing data to insurance form failed.", http.StatusInternalServerError)
		}
		// TODO: Close the temp file, zipWriter & serve the file
		zipWriter.Close()
		zF.Close()
		w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		http.ServeFile(w, r, zFName)
		// TODO: Export to PDF via gotenburg (hosted elsewhere? Dockerfile)
		// -> need to convert to ods first before converting to PDF @@
	default:
		http.Error(w, "Error: http method not allowed.", http.StatusMethodNotAllowed)
	}
}
