package main

import (
	"context"
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"teacup1592/form-filler/internal/api"
	datasrc "teacup1592/form-filler/internal/dataSrc"
	"teacup1592/form-filler/internal/database"
	"teacup1592/form-filler/internal/schoolForm"
)

var (
	httpAddr = flag.String("http", "localhost:8080", "HTTP service address to listen for incoming requests on")
)

func main() {
	// Local settings setup
	datasrc.LocalEnvSetup()

	// // Creating copy of source.xlsx
	// var rawSrcFile *os.File = datasrc.SourceLocal()
	// defer os.Remove(rawSrcFile.Name())
	// // Begin modifying file
	// formFile, err := excelize.OpenReader(rawSrcFile)
	// if err != nil {
	// 	log.Fatal("Excelize failed to open source file. ", err)
	// }
	// log.Println("Currently active sheet index: ", formFile.GetActiveSheetIndex())

	connPool, err := database.NewDBPool(context.Background(), os.Getenv("DATABASE_URL"))
	if err != nil {
		log.Fatal(err)
	}
	defer connPool.Close()

	server := &api.Server{
		SchoolForm: schoolForm.NewService(&database.DB{
			DbPool: connPool,
		}),
		HTTPAddress: *httpAddr,
	}
	errChann := make(chan error, 1)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	go func() {
		errChann <- server.Run(context.Background())
	}()

	select {
	case err = <-errChann:
	case <-ctx.Done():
		haltCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(haltCtx)
		stop()
		err = <-errChann
	}
	if err != nil {
		log.Fatal(err)
	}
}
