package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"time"

	"teacup1592/form-filler/internal/api"
	"teacup1592/form-filler/internal/dataSrc"
	"teacup1592/form-filler/internal/database"
	"teacup1592/form-filler/internal/schoolForm"
)

var (
    httpAddr = flag.String("http", ":8080", "HTTP service address to listen for incoming requests on")
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
            []string{
                dataSrc.SCH_FORM_NAME,
                dataSrc.INSUR_FORM_NAME,
                dataSrc.MOUNT_PASS_FORM_NAME,
            },
		),
		HTTPAddress: *httpAddr,
		Converter: &api.Converter{
			Unoserver: exec.Command("unoserver", "--port", "9000"),
		},
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
