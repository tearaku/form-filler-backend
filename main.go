package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"teacup1592/form-filler/internal/api"
	"teacup1592/form-filler/internal/dataSrc"
	"teacup1592/form-filler/internal/database"
	"teacup1592/form-filler/internal/schoolForm"
)

var (
	httpAddr = flag.String("http", "localhost:8080", "HTTP service address to listen for incoming requests on")
)

func main() {
	// TODO: Local settings setup
	dataSrc.LocalEnvSetup()

	connPool, err := database.NewDBPool(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer connPool.Close()

	s := &api.Server{
		SchoolForm: schoolForm.NewService(
			&database.DB{DbPool: connPool},
			schoolForm.FormFiller{},
		),
		HTTPAddress: *httpAddr,
	}
	ec := make(chan error, 1)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	go func() {
		ec <- s.Run(context.Background())
	}()

	select {
	case err = <-ec:
	case <-ctx.Done():
		haltCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		s.Shutdown(haltCtx)
		stop()
		err = <-ec
	}
	if err != nil {
		log.Fatal(err)
	}
}
