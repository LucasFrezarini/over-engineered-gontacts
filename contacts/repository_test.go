package contacts

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"go.uber.org/zap"
)

func TestRepositoryFindAll(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error while opening a stub database connection: %v", err)
	}

	defer db.Close()

	expectedContacts := []*Contact{
		&Contact{1, "Inosuke", "Hashibira"},
		&Contact{2, "Gonpachiro", "Kamaboko"},
	}

	rows := sqlmock.NewRows([]string{"id", "first_name", "last_name"})

	for _, c := range expectedContacts {
		rows.AddRow(c.ID, c.FirstName, c.LastName)
	}

	mock.ExpectQuery("SELECT (.+) FROM contact").WillReturnRows(rows).RowsWillBeClosed()

	repository := ProvideContactsRepository(db, zap.NewNop())
	contacts, err := repository.FindAll()

	if err != nil {
		t.Errorf("FindAll() returned an error %v, want nil", err)
	}

	if expected, got := len(expectedContacts), len(contacts); got != expected {
		t.Errorf("FindAll() len(contacts) = %d, want %d", got, expected)
	}

	for i, c := range contacts {
		if expected := expectedContacts[i]; !reflect.DeepEqual(c, expected) {
			t.Errorf("FindAll contacts[%d] = %v, want %v", i, c, expected)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("unfulfilled mock expectations: %v", err)
	}
}
