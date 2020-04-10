package phone

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/wire"
	"go.uber.org/zap"
)

// GenericRepository defines the structure of this package's repository
// Defined especially to allow mocking in unit testing
type GenericRepository interface {
	FindByContactID(id int) ([]Phone, error)
}

// Repository contains all the persistence related methods for the phone entity
type Repository struct {
	DB     *sql.DB
	Logger *zap.Logger
}

// ProvideRepository creates a new Repository with the dependencies provided.
// Created especially for the use of Wire, which will inject the dependencies via DI
func ProvideRepository(db *sql.DB, logger *zap.Logger) *Repository {
	return &Repository{db, logger.Named("PhoneRepository")}
}

// FindByContactID returns all the phones registered for the provided contact id
func (r *Repository) FindByContactID(id int) ([]Phone, error) {
	raw := "SELECT id, contact_id, number, type FROM phone WHERE contact_id = ?"
	rows, err := r.DB.Query(raw, id)

	if err != nil {
		msg := fmt.Sprintf("FindByContactID(%d): error while executing query: %v", id, err)
		r.Logger.Error(msg)
		return nil, errors.New(msg)
	}

	defer rows.Close()
	phones := make([]Phone, 0)

	for rows.Next() {
		var phone Phone

		if err := rows.Scan(&phone.ID, &phone.ContactID, &phone.Number, &phone.Type); err != nil {
			msg := fmt.Sprintf("FindByContactID(%d): error while scanning rows: %v", id, err)
			r.Logger.Error(msg)
			return nil, errors.New(msg)
		}

		phones = append(phones, phone)
	}

	return phones, nil
}

// RepositorySet is the wire set that contains all the provides for this repository
var RepositorySet = wire.NewSet(
	ProvideRepository,
	wire.Bind(new(GenericRepository), new(*Repository)),
)
