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
	FindByContactID(id int) ([]Email, error)
	Create(contactID int, emails ...string) ([]Email, error)
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
func (r *Repository) FindByContactID(id int) ([]Email, error) {
	raw := "SELECT id, contact_id, address FROM email WHERE contact_id = ?"

	rows, err := r.DB.Query(raw, id)
	if err != nil {
		msg := fmt.Sprintf("FindByContactID(%d): error while preparing statement: %v", id, err)
		r.Logger.Error(msg)
		return nil, errors.New(msg)
	}

	defer rows.Close()
	emails := make([]Email, 0)

	for rows.Next() {
		var email Email

		if err := rows.Scan(&email.ID, &email.ContactID, &email.Address); err != nil {
			msg := fmt.Sprintf("FindByContactID(%d): error while scanning row: %v", id, err)
			r.Logger.Error(msg)
			return nil, errors.New(msg)
		}

		emails = append(emails, email)
	}

	return emails, nil
}

// Create creates one or more emails for the contactID provided
func (r *Repository) Create(contactID int, emails ...string) ([]Email, error) {
	insertedEmails := make([]Email, 0, len(emails))

	for _, address := range emails {
		email, err := r.createSingleEmail(contactID, address)

		if err != nil {
			return nil, fmt.Errorf("error while inserting email into the database: %w", err)
		}

		insertedEmails = append(insertedEmails, email)
	}

	return insertedEmails, nil
}

func (r *Repository) createSingleEmail(contactID int, address string) (Email, error) {
	raw := "INSERT INTO email (contact_id, address) VALUES (?, ?)"

	stmt, err := r.DB.Prepare(raw)
	if err != nil {
		return Email{}, fmt.Errorf("createSingleEmail: error while preparing statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(contactID, address)
	if err != nil {
		return Email{}, fmt.Errorf("createSingleEmail: error while executing insert query: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return Email{}, fmt.Errorf("createSingleEmail: error while fetching the last inserted ID: %w", err)
	}

	return Email{
		ID:        int(id),
		ContactID: contactID,
		Address:   address,
	}, nil
}

// RepositorySet is the wire set which contains all the binding necessary
// to create a new email Repository
var RepositorySet = wire.NewSet(
	ProvideEmailRepository,
	wire.Bind(new(GenericRepository), new(*Repository)),
)
