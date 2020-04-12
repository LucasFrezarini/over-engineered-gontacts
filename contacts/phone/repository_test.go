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

func TestPhoneCreate(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("unexpected error while opening a stub database connection: %v", err)
	}

	defer db.Close()

	contactID := 2
	phonesData := []CreatePhoneData{
		CreatePhoneData{ContactID: contactID, Number: "1122223333", Type: PhoneTypeHome},
		CreatePhoneData{ContactID: contactID, Number: "33444445555", Type: PhoneTypeMobile},
		CreatePhoneData{ContactID: contactID, Number: "+5511911112222", Type: PhoneTypeWork},
		CreatePhoneData{ContactID: contactID, Number: "+5511933332222", Type: PhoneTypeFax},
	}

	for i, phoneData := range phonesData {
		mock.ExpectPrepare("INSERT INTO phone").ExpectExec().WithArgs(contactID, phoneData.Type, phoneData.Number).WillReturnResult(sqlmock.NewResult(int64(i+1), 1))
	}

	repository := ProvideRepository(db, zap.NewNop())
	insertedPhones, err := repository.Create(contactID, phonesData...)
	if err != nil {
		t.Errorf("Create(%d, %v) returned a non-nil error '%v', want nil", contactID, phonesData, err)
	}

	if expected, got := len(phonesData), len(insertedPhones); expected != got {
		t.Errorf("Create(%d, %v) returned %d phones, want %d", contactID, phonesData, got, expected)
	}

	for i, p := range insertedPhones {
		if p.ID == 0 {
			t.Errorf("Create(%d, %v) phone[%d].ID == 0, want != 0", contactID, phonesData, i)
		}

		if expected, got := contactID, p.ContactID; expected != got {
			t.Errorf("Create(%d, %v) phone[%d].ContactID == %d, want %d", contactID, phonesData, i, got, expected)
		}
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Create(%d, %v) mock expectations weren't met: %v", contactID, phonesData, err)
	}
}
