package database

import (
	"database/sql"
	"fmt"
)

type Connection struct {
	Host     string
	Port     int
	Username string
	DBName   string
	SSLMode  string
	Password string
}

func NewPostgresConnection(inf Connection) (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("host=%s port=%d username=%s dbname=%s sslmode=%s password=%s", inf.Host, inf.Port, inf.Username, inf.DBName, inf.SSLMode, inf.Password))

	if err != nil {
		return nil, err
	}

	if err := db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
