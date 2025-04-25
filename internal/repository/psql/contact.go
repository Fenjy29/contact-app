package psql

import (
	"contact-app/internal/domain"
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Contacts struct {
	db *sql.DB
}

func NewContacts(db *sql.DB) *Contacts {
	return &Contacts{db}
}

func (c *Contacts) Create(ctx context.Context, contact domain.Contact) error {
	_, err := c.db.Exec("INSERT INTO contacts (name, phone, email, address) VALUES ($1, $2, $3, $4)", contact.Name, contact.Phone, contact.Email, contact.Address)

	return err
}

func (c *Contacts) GetByID(ctx context.Context, id int64) (domain.Contact, error) {
	var contact domain.Contact
	err := c.db.QueryRow("SELECT (name, phone, email, address) FROM contacts 	WHERE id = $1", id).Scan(&contact.ID, &contact.Name, &contact.Phone, &contact.Email, &contact.Address)
	if err == sql.ErrNoRows {
		return contact, domain.ErrContactNotFound
	}

	return contact, err
}

func (c *Contacts) GetAll(ctx context.Context) ([]domain.Contact, error) {
	contacts := make([]domain.Contact, 0)

	rows, err := c.db.QueryContext(ctx, "SELECT id, name, phone, email, address FROM contacts")
	if err != nil {
		return nil, fmt.Errorf("failed to query contacts: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var contact domain.Contact
		// Теперь сканируем все поля, включая ID
		if err := rows.Scan(
			&contact.ID,
			&contact.Name,
			&contact.Phone,
			&contact.Email,
			&contact.Address,
		); err != nil {
			return nil, fmt.Errorf("failed to scan contact: %w", err)
		}
		contacts = append(contacts, contact)
	}

	// Проверяем ошибки после итерации
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("rows iteration error: %w", err)
	}

	return contacts, nil
}

func (c *Contacts) Delete(ctx context.Context, id int64) error {
	_, err := c.db.Exec("DELETE FROM contacts WHERE id = $1", id)

	return err
}

func (c *Contacts) Update(ctx context.Context, id int64, inp domain.UpdateContact) error {
	setValues := make([]string, 0)
	args := make([]interface{}, 0)
	argId := 1

	if inp.Name != nil {
		setValues = append(setValues, fmt.Sprintf("name=$%d", argId))
		args = append(args, *inp.Name)
		argId++
	}

	if inp.Email != nil {
		setValues = append(setValues, fmt.Sprintf("email=$%d", argId))
		args = append(args, *inp.Email)
		argId++
	}

	if inp.Phone != nil {
		setValues = append(setValues, fmt.Sprintf("phone=$%d", argId))
		args = append(args, *inp.Phone)
		argId++
	}

	if inp.Address != nil {
		setValues = append(setValues, fmt.Sprintf("address=$%d", argId))
		args = append(args, *inp.Address)
		argId++
	}

	setQuery := strings.Join(setValues, ", ")

	query := fmt.Sprintf("UPDATE contacts SET %s WHERE id = %d", setQuery, argId)
	args = append(args, id)

	_, err := c.db.Exec(query, args...)

	return err
}
