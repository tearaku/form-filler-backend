package api

// Setting up, starting & shutting down of server(s)

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strings"
	"sync"
	"teacup1592/form-filler/internal/schoolForm"
)

// Server for the API.
type Server struct {
	HTTPAddress string
	// GRPCAddress string

	SchoolForm *schoolForm.Service

	// grpc   *grpcServer
	http *httpServer

	stopFn sync.Once
}

func (s *Server) Run(ctx context.Context) (err error) {
	var ec = make(chan error, 1)
	s.http = &httpServer{
		schoolForm: s.SchoolForm,
	}
	go func() {
		err := s.http.Run(ctx, s.HTTPAddress)
		if err != nil {
			err = fmt.Errorf("HTTP server error: %w", err)
		}
		ec <- err
	}()

	var eMsg []string
	for i := 0; i < cap(ec); i++ {
		if err := <-ec; err != nil {
			eMsg = append(eMsg, err.Error())
			// Something about gracefully shutting down server?
			if ctx.Err() == nil {
				s.Shutdown(context.Background())
			}
		}
	}
	if len(eMsg) > 0 {
		err = errors.New(strings.Join(eMsg, ", "))
	}
	return err
}

// Shutting down server
func (s *Server) Shutdown(ctx context.Context) {
	s.stopFn.Do(func() {
		s.http.Shutdown(ctx)
	})
}

type httpServer struct {
	schoolForm *schoolForm.Service
	// middleware func(http.Handler) http.Handler
	http *http.Server
}

// Running HTTP server
func (s *httpServer) Run(ctx context.Context, address string) error {
	handler := NewHTTPServer(s.schoolForm)

	// TODO: Middleware here if needed

	s.http = &http.Server{
		Addr:    address,
		Handler: handler,
	}
	log.Printf("HTTP server listening at: %s\n", address)
	if err := s.http.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

// Shutting down HTTP server
func (s *httpServer) Shutdown(ctx context.Context) {
	log.Println("Shutting down HTTP server")
	if s.http != nil {
		if err := s.http.Shutdown(ctx); err != nil {
			log.Println("Graceful shutdown of HTTP server: failed!")
		}
	}
}
