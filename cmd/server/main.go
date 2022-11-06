package main

import (
	"fmt"
	"log"
	"os"

	"github.com/by-sabbir/company-microservice-rest/internal/company"
	"github.com/by-sabbir/company-microservice-rest/internal/db"

	transportHttp "github.com/by-sabbir/company-microservice-rest/internal/transport/http"
)

func Run() error {
	db, err := db.NewDatabase()
	if err != nil {
		fmt.Println("Failed to connect to database")
		return err
	}
	log.Println("Successfully Connected to the DB!")

	if err := db.MigrateDB(); err != nil {
		return fmt.Errorf("migrations failed because of: %w", err)
	}

	svc := company.NewService(db)

	httpHandler := transportHttp.NewHandler(svc)
	log.Println("service started at: ", httpHandler.Server.Addr)
	if err := httpHandler.Serve(); err != nil {
		log.Println("service stopped: ", err)
		os.Exit(1)
	}
	return nil
}

func main() {
	if err := Run(); err != nil {
		log.Fatalf("could not run the service: %+v\n", err)
	}
}
