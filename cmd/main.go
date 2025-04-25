package main

import (
	"contact-app/pkg/database"
	"log"
	"net/http"
	"time"
)

func main() {
	db, err := database.NewPostgresConnection(database.Connection{
		Host:     "localhost",
		Port:     5432,
		Username: "postgres",
		DBName:   "postgres",
		SSLMode:  "prefer",
		Password: "1234",
	})

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler.InitRouter(),
	}

	log.Println("SERVER STARTED AT", time.Now().Format(time.RFC3339))

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
