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
		&Contact{ID: 1, FirstName: "Inosuke", LastName: "Hashibira"},
		&Contact{ID: 2, FirstName: "Gonpachiro", LastName: "Kamaboko"},
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

func TestRepositoryCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error while opening a stub database connection: %v", err)
	}

	defer db.Close()

	data := Contact{
		FirstName: "Zenitsu",
		LastName:  "Agatsuma",
	}

	mock.ExpectPrepare("INSERT INTO contact").ExpectExec().WithArgs(data.FirstName, data.LastName).WillReturnResult(sqlmock.NewResult(1, 1))

	repository := ProvideContactsRepository(db, zap.NewNop())
	contact, err := repository.Create(data)
	if err != nil {
		t.Errorf("repository.Create(%T): returned an error while creating a new contact: %v", data, err)
	}

	if expected := 1; contact.ID != expected {
		t.Errorf("repository.Create(%T): contact.ID == %d, want %d", data, contact.ID, expected)
	}

	if expected := data.FirstName; contact.FirstName != expected {
		t.Errorf("repository.Create(%T): contact.FirstName == %s, want %s", data, contact.FirstName, expected)
	}

	if expected := data.LastName; contact.LastName != expected {
		t.Errorf("repository.Create(%T): contact.LastName == %s, want %s", data, contact.LastName, expected)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("repository.Create(%T): unfulfilled mock expectations: %v", data, err)
	}

}
