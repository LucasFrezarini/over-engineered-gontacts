package phone

import (
	"database/sql"
	"errors"
	"fmt"

	"github.com/google/wire"
	"go.uber.org/zap"
)

// CreatePhoneData defines the fields that need to be provided in order to create a new
// phone record in the database
type CreatePhoneData struct {
	Number string `json:"number"`
	Type   string `json:"type"`
}

// GenericRepository defines the structure of this package's repository
// Defined especially to allow mocking in unit testing
type GenericRepository interface {
	FindByContactID(id int) ([]Phone, error)
	Create(contactID int, phones ...CreatePhoneData) ([]Phone, error)
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

// Create creates new phones registred for the provided ContactID
func (r *Repository) Create(contactID int, phones ...CreatePhoneData) ([]Phone, error) {
	insertedPhones := make([]Phone, 0, len(phones))

	for _, data := range phones {
		phone, err := r.createSinglePhone(contactID, data)
		if err != nil {
			msg := fmt.Sprintf("Create: error while creating phone: %v", err)
			r.Logger.Error(msg)
			return nil, errors.New(msg)
		}

		insertedPhones = append(insertedPhones, phone)
	}

	return insertedPhones, nil
}

func (r *Repository) createSinglePhone(contactID int, phone CreatePhoneData) (Phone, error) {
	raw := "INSERT INTO phone (contact_id, type, number) VALUES (?, ?, ?)"
	stmt, err := r.DB.Prepare(raw)

	if err != nil {
		return Phone{}, fmt.Errorf("createSinglePhone: error while preparing statement: %w", err)
	}

	defer stmt.Close()

	result, err := stmt.Exec(contactID, phone.Type, phone.Number)
	if err != nil {
		return Phone{}, fmt.Errorf("createSinglePhone: error while executing statement: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Phone{}, fmt.Errorf("createSinglePhone: error while retrieving the last inserted ID: %w", err)
	}

	return Phone{
		ID:        int(id),
		ContactID: contactID,
		Type:      phone.Type,
		Number:    phone.Number,
	}, nil

}

// RepositorySet is the wire set that contains all the provides for this repository
var RepositorySet = wire.NewSet(
	ProvideRepository,
	wire.Bind(new(GenericRepository), new(*Repository)),
)
