package domain

import "errors"

var (
	ErrContactNotFound = errors.New("Contact not found")
)

type Contact struct {
	ID      int    `json: "id"`
	Name    string `json: "name"`
	Phone   string `json: "phone"`
	Email   string `json: "email"`
	Address string `json: "address"`
}

type UpdateContact struct {
	Name    *string `json: "name"`
	Phone   *string `json: "phone"`
	Email   *string `json: "email"`
	Address *string `json: "address"`
}
