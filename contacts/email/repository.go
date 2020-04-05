package email

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/wire"
	"go.uber.org/zap"
)

// GenericRepository defines the structure of a generic email's repository
// created to facilitate the mocking in unit testing
type GenericRepository interface {
	FindByContactID(id int) ([]*Email, error)
}

// A Repository can perform all the CRUD logic of the
// contact emails
type Repository struct {
	DB     *sql.DB
	Logger *zap.Logger
}

// ProvideEmailRepository creates a new repository and return its pointer.
// Especially to be used by Wire, providing the dependencies via DI
func ProvideEmailRepository(db *sql.DB, logger *zap.Logger) *Repository {
	return &Repository{db, logger.Named("EmailRepository")}
}

// FindByContactID return all the emails registered for the contact with
// the id provided as parameter
func (r *Repository) FindByContactID(id int) ([]*Email, error) {
	raw := "SELECT id, contact_id, address FROM email WHERE contact_id = ?"

	rows, err := r.DB.Query(raw, id)
	if err != nil {
		msg := fmt.Sprintf("FindByContactID(%d): error while preparing statement: %v", id, err)
		r.Logger.Error(msg)
		return nil, errors.New(msg)
	}

	defer rows.Close()
	emails := make([]*Email, 0)

	for rows.Next() {
		var email Email

		if err := rows.Scan(&email.ID, &email.ContactID, &email.Address); err != nil {
			msg := fmt.Sprintf("FindByContactID(%d): error while scanning row: %v", id, err)
			r.Logger.Error(msg)
			return nil, errors.New(msg)
		}

		emails = append(emails, &email)
	}

	return emails, nil
}

// EmailRepositorySet is the wire set which contains all the binding necessary
// to create a new email Repository
var EmailRepositorySet = wire.NewSet(
	ProvideEmailRepository,
	wire.Bind(new(GenericRepository), new(*Repository)),
)
