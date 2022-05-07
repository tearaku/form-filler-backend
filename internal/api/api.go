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

	http   *httpServer
	stopFn sync.Once
}

func (server *Server) Run(ctx context.Context) (err error) {
	var errorCode = make(chan error, 1)
	server.http = &httpServer{
		schoolForm: server.SchoolForm,
	}
	go func() {
		err := server.http.Run(ctx, server.HTTPAddress)
		if err != nil {
			err = fmt.Errorf("HTTP server error: %w", err)
		}
		errorCode <- err
	}()

	var errorMsg []string
	for i := 0; i < cap(errorCode); i++ {
		if err := <-errorCode; err != nil {
			errorMsg = append(errorMsg, err.Error())
			// Something about gracefully shutting down server?
			if ctx.Err() == nil {
				server.Shutdown(context.Background())
			}
		}
	}
	if len(errorMsg) > 0 {
		err = errors.New(strings.Join(errorMsg, ", "))
	}
	return err
}

// Shutting down server
func (server *Server) Shutdown(ctx context.Context) {
	server.stopFn.Do(func() {
		server.http.Shutdown(ctx)
	})
}

type httpServer struct {
	schoolForm *schoolForm.Service
	// middleware func(http.Handler) http.Handler
	http *http.Server
}

// Running HTTP server
func (server *httpServer) Run(ctx context.Context, address string) error {
	handler := NewHTTPServer(server.schoolForm)

	// TODO: Middleware here if needed

	server.http = &http.Server{
		Addr:    address,
		Handler: handler,
	}
	log.Printf("HTTP server listening at: %s\n", address)
	if err := server.http.ListenAndServe(); err != nil {
		return err
	}
	return nil
}

// Shutting down HTTP server
func (server *httpServer) Shutdown(ctx context.Context) {
	log.Println("Shutting down HTTP server")
	if server.http != nil {
		if err := server.http.Shutdown(ctx); err != nil {
			log.Println("Graceful shutdown of HTTP server: failed!")
		}
	}
}
