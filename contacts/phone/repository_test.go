package phone

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
	expectedPhones := []Phone{
		Phone{ID: 1, ContactID: contactID, Number: "1122223333", Type: PhoneTypeHome},
		Phone{ID: 2, ContactID: contactID, Number: "33444445555", Type: PhoneTypeMobile},
		Phone{ID: 3, ContactID: contactID, Number: "+5511911112222", Type: PhoneTypeWork},
		Phone{ID: 4, ContactID: contactID, Number: "+5511933332222", Type: PhoneTypeFax},
	}

	rows := sqlmock.NewRows([]string{"id", "contact_id", "number", "type"})

	for _, p := range expectedPhones {
		rows.AddRow(p.ID, p.ContactID, p.Number, p.Type)
	}

	mock.ExpectQuery("SELECT (.+) FROM phone").WithArgs(contactID).WillReturnRows(rows).RowsWillBeClosed()

	repository := ProvideRepository(db, zap.NewNop())
	phones, err := repository.FindByContactID(contactID)

	if err != nil {
		t.Errorf("FindByContactID(%d) returned an error: '%v', want nil", contactID, err)
	}

	if expected, got := len(expectedPhones), len(phones); expected != got {
		t.Errorf("FindByContactID(%d) returned %d phones, want %d", contactID, got, expected)
	}

	for i, c := range phones {
		if expected := expectedPhones[i]; !reflect.DeepEqual(expected, c) {
			t.Errorf("FindByContactID(%d) phones[%d] = '%v', want '%v'", contactID, i, c, expected)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("FindByContactID(%d) unfulfilled mock expectations: %v", contactID, err)
	}
}
