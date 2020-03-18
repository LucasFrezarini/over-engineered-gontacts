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
}

type ContactsRepository struct {
	DB     *sql.DB
	logger *zap.Logger
}

func ProvideContactsRepository(db *sql.DB, logger *zap.Logger) *ContactsRepository {
	return &ContactsRepository{DB: db, logger: logger.Named("ContactsRepository")}
}

func (r *ContactsRepository) FindAll() ([]*Contact, error) {
	r.logger.Debug("Executing FindAll()")
	stmt := "SELECT id, first_name, last_name FROM contact"
	r.logger.Debug("stmt: " + stmt)
	rows, err := r.DB.Query(stmt)

	if err != nil {
		return nil, fmt.Errorf("FindAll: error while fetching contacts: %w", err)
	}

	defer rows.Close()
	contacts := make([]*Contact, 0)

	for rows.Next() {
		var contact Contact

		if err := rows.Scan(&contact.ID, &contact.FirstName, &contact.LastName); err != nil {
			return nil, fmt.Errorf("FindAll: error while scanning rows: %w", err)
		}

		contacts = append(contacts, &contact)
	}

	return contacts, nil
}

var RepositorySet = wire.NewSet(
	ProvideContactsRepository,
	wire.Bind(new(Repository), new(*ContactsRepository)),
)
