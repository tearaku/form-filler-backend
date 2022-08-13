package api

// Various HTTP REST endpoint handlers

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

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
	if r.Method == http.MethodOptions {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		w.Header().Set("Access-Control-Max-Age", "3600")
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Access-Control-Allow-Origin", os.Getenv("FRONTEND_URL"))
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
	log.Printf("Request from: %s\n", r.URL.Host)
	switch r.Method {
	case http.MethodGet:
		eventId := r.URL.Path[len("/event/"):]
		// Ping gotenberg service to wake it up
		convErr := make(chan error, 1)
		go func() {
			if err := schoolForm.GotenbergPing(); err != nil {
				convErr <- err
				return
			}
			convErr <- nil
		}()
		// Fetch relevant data from database based on event ID
		eventInfo, err := s.schoolForm.GetEventInfo(r.Context(), schoolForm.GetEventInfoParams{EventID: eventId})
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Error: fetching event data from db failed.", http.StatusInternalServerError)
			return
		}
		err = s.schoolForm.FetchAttendances(r.Context(), schoolForm.FetchAttendancesParams{
			EventInfo: eventInfo,
		})
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Error: fetching event data from db failed.", http.StatusInternalServerError)
			return
		}
		clubLeader, err := s.schoolForm.GetMemberByDept(r.Context(), schoolForm.GetMemberByDeptParams{
			Description: "社長",
		})
		if err != nil {
			log.Println(err.Error())
			http.Error(w, "Error: fetching club leader from db failed.", http.StatusInternalServerError)
			return
		}

		// Begin modifying excel data, and adding to zip file
		zA, err := schoolForm.NewArchiver("event_" + eventId + "_*.zip")
		if err != nil {
			http.Error(w, "Error: failed to create zip file.", http.StatusInternalServerError)
			return
		}
		defer zA.Cleanup()
		ec := make(chan error, 3)
		go func() {
			// Check if gotenberg service is healthy
			if err := <-convErr; err != nil {
				ec <- err
				return
			}
			// Being conversion
			if err = s.schoolForm.WriteSchForm(eventInfo, clubLeader, zA); err != nil {
				log.Printf("writeSchForm err: %v\n", err.Error())
				ec <- errors.New("error: writing data to school form failed")
				return
			}
			ec <- nil
		}()
		go func() {
			if err = s.schoolForm.WriteInsuranceForm(eventInfo, zA); err != nil {
				log.Printf("writeInsForm err: %v\n", err.Error())
				ec <- errors.New("error: writing data to insurance form failed")
				return
			}
			ec <- nil
		}()
		go func() {
			if err = s.schoolForm.WriteMountPass(eventInfo, zA); err != nil {
				log.Printf("writeMountPass err: %v\n", err.Error())
				ec <- errors.New("error: writing data to mountain pass form failed")
				return
			}
			ec <- nil
		}()
		var eMsg strings.Builder
		for i := 0; i < 3; i++ {
			if err := <-ec; err != nil {
				eMsg.WriteString(fmt.Sprintf("err %d: %s\n", i+1, err.Error()))
			}
		}
		if eMsg.Len() != 0 {
			http.Error(w, eMsg.String(), http.StatusInternalServerError)
			return
		}
		zA.ZipW.Close()
		http.ServeFile(w, r, zA.TempF.Name())
	default:
		http.Error(w, "Error: http method not allowed.", http.StatusMethodNotAllowed)
		return
	}
}
