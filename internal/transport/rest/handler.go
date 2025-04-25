package rest

import (
	"contact-app/internal/domain"
	"context"
)

type Contacts interface {
	Create(ctx context.Context, contact domain.Contact) error
	GetByID(ctx context.Context, id int64) (domain.Contact, error)
	GetAll(ctx context.Context) ([]domain.Contact, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, inp domain.UpdateContact) error
}

func