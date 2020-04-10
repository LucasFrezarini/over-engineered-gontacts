package contacts

import (
	"github.com/LucasFrezarini/go-contacts/contacts/email"
	"github.com/LucasFrezarini/go-contacts/contacts/phone"
)

var contactsList = []*Contact{
	&Contact{
		ID:        1,
		FirstName: "Inosuke",
		LastName:  "Hashibira",
	},
	&Contact{
		ID:        2,
		FirstName: "Gonpachiro",
		LastName:  "Kamaboko",
	},
}

type MockedContactsRepository struct {
	id int
}

func (m *MockedContactsRepository) FindAll() ([]*Contact, error) {
	return contactsList, nil
}

func (m *MockedContactsRepository) Create(c Contact) (*Contact, error) {
	m.id++

	return &Contact{
		ID:        m.id,
		FirstName: c.FirstName,
		LastName:  c.LastName,
	}, nil
}

var emailsList = []email.Email{
	email.Email{ID: 1, ContactID: 1, Address: "inosuke@gmail.com"},
	email.Email{ID: 2, ContactID: 1, Address: "pigassault@outlook.com"},
	email.Email{ID: 3, ContactID: 2, Address: "tanjirou@gmail.com"},
}

var phonesList = []phone.Phone{
	phone.Phone{ID: 1, ContactID: 1, Number: "1122223333", Type: phone.PhoneTypeHome},
	phone.Phone{ID: 2, ContactID: 1, Number: "33444445555", Type: phone.PhoneTypeMobile},
	phone.Phone{ID: 3, ContactID: 1, Number: "+5511911112222", Type: phone.PhoneTypeWork},
	phone.Phone{ID: 4, ContactID: 1, Number: "+5511933332222", Type: phone.PhoneTypeFax},
	phone.Phone{ID: 5, ContactID: 2, Number: "11955554444", Type: phone.PhoneTypeMobile},
	phone.Phone{ID: 6, ContactID: 2, Number: "1122223333", Type: phone.PhoneTypeHome},
}

func filterEmailsByContactID(id int) []email.Email {
	filteredEmails := make([]email.Email, 0)

	for _, e := range emailsList {
		if e.ContactID == id {
			filteredEmails = append(filteredEmails, e)
		}
	}

	return filteredEmails
}

func filterPhonesByContactID(id int) []phone.Phone {
	filteredPhones := make([]phone.Phone, 0)

	for _, p := range phonesList {
		if p.ContactID == id {
			filteredPhones = append(filteredPhones, p)
		}
	}

	return filteredPhones
}

type MockedEmailRepository struct{}

func (m *MockedEmailRepository) FindByContactID(id int) ([]email.Email, error) {
	return filterEmailsByContactID(id), nil
}

type MockedPhoneRepository struct{}

func (p *MockedPhoneRepository) FindByContactID(id int) ([]phone.Phone, error) {
	return filterPhonesByContactID(id), nil
}
