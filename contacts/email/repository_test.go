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

	expectedEmails := []Email{
		Email{ID: 1, ContactID: contactID, Address: "inosuke@gmail.com"},
		Email{ID: 2, ContactID: contactID, Address: "zenitsu@yahoo.com"},
	}

	rows := sqlmock.NewRows([]string{"id", "contact_id", "address"})

	for _, e := range expectedEmails {
		rows.AddRow(e.ID, e.ContactID, e.Address)
	}

	mock.ExpectQuery("SELECT (.+) FROM email").WillReturnRows(rows).RowsWillBeClosed()

	repository := ProvideEmailRepository(db, zap.NewNop())
	emails, err := repository.FindByContactID(contactID)

	if err != nil {
		t.Errorf("FindByContactID(%d) returned an error: '%v', want nil", contactID, err)
	}

	if expected, got := len(expectedEmails), len(emails); expected != got {
		t.Errorf("FindByContactID(%d) returned %d emails, want %d", contactID, got, expected)
	}

	for i, c := range emails {
		if expected := expectedEmails[i]; !reflect.DeepEqual(expected, c) {
			t.Errorf("FindByContactID(%d) emails[%d] = %v, want %v", contactID, i, c, expected)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("FindByContactID(%d) unfulfilled mock expectations: %v", contactID, err)
	}
}

func TestEmailCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error while opening a stub database connection: %v", err)
	}

	defer db.Close()

	contactID := 2

	emails := []string{"zenitsu01@gmail.com", "zenitsu02@yahoo.com"}

	for i, e := range emails {
		mock.ExpectPrepare("INSERT INTO email").ExpectExec().WithArgs(contactID, e).WillReturnResult(sqlmock.NewResult(int64(i+1), 1))
	}

	repository := ProvideEmailRepository(db, zap.NewNop())
	insertedEmails, err := repository.Create(contactID, emails...)

	if err != nil {
		t.Errorf("Create(%d, %v) returned a non-nil error '%v', want nil", contactID, emails, err)
	}

	if expected, got := len(emails), len(insertedEmails); expected != got {
		t.Errorf("Create(%d, %v) returned %d emails, want %d", contactID, emails, got, expected)
	}

	for i, e := range insertedEmails {
		if e.ID == 0 {
			t.Errorf("Create(%d, %v) email[%d].ID == 0, want != 0", contactID, emails, i)
		}
	}

	if err = mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Create(%d, %v): unfulfilled mock expectations: %v", contactID, emails, err)
	}
}
