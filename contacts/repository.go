package contacts

import (
	"database/sql"
	"fmt"

	"github.com/google/wire"
	"go.uber.org/zap"
)

// Repository defines the structure of a generic contact repository
// this interface was created to facilitate the mocking in the unit tests
type Repository interface {
	FindAll() ([]*Contact, error)
	Create(c Contact) (*Contact, error)
	DeleteByID(id int) error
}

type ContactsRepository struct {
	DB     *sql.DB
	Logger *zap.Logger
}

func ProvideContactsRepository(db *sql.DB, logger *zap.Logger) *ContactsRepository {
	return &ContactsRepository{DB: db, Logger: logger.Named("ContactsRepository")}
}

func (r *ContactsRepository) FindAll() ([]*Contact, error) {
	stmt := `SELECT id, first_name, last_name FROM contact`
	rows, err := r.DB.Query(stmt)

	if err != nil {
		return nil, fmt.Errorf("FindAll(): error while fetching contacts: %w", err)
	}

	defer rows.Close()
	contacts := make([]*Contact, 0)

	for rows.Next() {
		var contact Contact

		if err := rows.Scan(&contact.ID, &contact.FirstName, &contact.LastName); err != nil {
			return nil, fmt.Errorf("FindAll(): error while scanning rows: %w", err)
		}

		contacts = append(contacts, &contact)
	}

	return contacts, nil
}

func (r *ContactsRepository) Create(c Contact) (*Contact, error) {
	raw := "INSERT INTO contact (first_name, last_name) VALUES (?, ?)"

	stmt, err := r.DB.Prepare(raw)
	if err != nil {
		return nil, fmt.Errorf("create: error while preparing statement: %w", err)
	}
	defer stmt.Close()

	result, err := stmt.Exec(c.FirstName, c.LastName)
	if err != nil {
		return nil, fmt.Errorf("create: error while executing insert query: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("create: error while fetching the last inserted ID: %w", err)
	}

	c.ID = int(id)
	return &c, nil
}

func (r *ContactsRepository) DeleteByID(id int) error {
	raw := "DELETE FROM contact WHERE id = ?"

	stmt, err := r.DB.Prepare(raw)
	if err != nil {
		return fmt.Errorf("deleteByID: error while preparing statement: %w", err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(id)
	if err != nil {
		return fmt.Errorf("deleteByID: error while executing the delete query: %w", err)
	}

	return nil
}

var RepositorySet = wire.NewSet(
	ProvideContactsRepository,
	wire.Bind(new(Repository), new(*ContactsRepository)),
)
