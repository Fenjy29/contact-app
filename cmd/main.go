package main

import (
	"contact-app/internal/handlers/rest"
	"contact-app/internal/repository/psql"
	"contact-app/internal/service"
	"contact-app/pkg/database"
	_ "github.com/lib/pq"
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
		SSLMode:  "disable",
		Password: "1234",
	})

	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	booksRepo := psql.NewContacts(db)
	booksService := service.NewContacts(booksRepo)
	handler := rest.NewHandler(booksService)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: handler.InitRouter(),
	}

	log.Println("SERVER STARTED AT", time.Now().Format(time.RFC3339))

	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
