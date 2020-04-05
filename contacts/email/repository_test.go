package email

import (
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"go.uber.org/zap"
)

func TestFindByContactID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error while opening a stub database connection: %v", err)
	}

	defer db.Close()

	contactID := 2

	expectedEmails := []*Email{
		&Email{ID: 1, ContactID: contactID, Address: "inosuke@gmail.com"},
		&Email{ID: 2, ContactID: contactID, Address: "zenitsu@yahoo.com"},
	}

	rows := sqlmock.NewRows([]string{"id", "contact_id", "address"})

	for _, e := range expectedEmails {
		rows.AddRow(e.ID, e.ContactID, e.Address)
	}

	mock.ExpectQuery("SELECT (.+) FROM email").WillReturnRows(rows).RowsWillBeClosed()

	repository := ProvideEmailRepository(db, zap.NewNop())
	emails, err := repository.FindByContactID(contactID)

	if err != nil {
		t.Errorf("FindByEmailID(%d) returned an error: '%v', want nil", contactID, err)
	}

	if expected, got := len(expectedEmails), len(emails); expected != got {
		t.Errorf("FindByEmailID(%d) returned %d emails, want %d", contactID, got, expected)
	}

	for i, c := range emails {
		if expected := expectedEmails[i]; !reflect.DeepEqual(expected, c) {
			t.Errorf("FindByEmail(%d) contacts[%d] = %v, want %v", contactID, i, c, expected)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("FindByEmail(%d) unfulfilled mock expectations: %v", contactID, err)
	}
}
