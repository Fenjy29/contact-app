package service

import (
	"contact-app/internal/domain"
	"context"
	"fmt"
)

type ContactRepository interface {
	Create(ctx context.Context, contact domain.Contact) error
	GetByID(ctx context.Context, id int64) (domain.Contact, error)
	GetAll(ctx context.Context) ([]domain.Contact, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, inp domain.UpdateContact) error
}

type Contacts struct {
	repo ContactRepository
}

func NewContacts(repo ContactRepository) *Contacts {
	return &Contacts{
		repo: repo,
	}
}

func (c *Contacts) Create(ctx context.Context, contact domain.Contact) error {
	err := c.repo.Create(ctx, contact)
	if err != nil {
		return fmt.Errorf("failed to create contact: %w", err)
	}

	return nil
}

func (c *Contacts) GetByID(ctx context.Context, id int64) (domain.Contact, error) {
	return c.repo.GetByID(ctx, id)
}

func (c *Contacts) GetAll(ctx context.Context) ([]domain.Contact, error) {
	return c.repo.GetAll(ctx)
}

func (c *Contacts) Delete(ctx context.Context, id int64) error {
	return c.repo.Delete(ctx, id)
}

func (c *Contacts) Update(ctx context.Context, id int64, inp domain.UpdateContact) error {
	return c.repo.Update(ctx, id, inp)
}
