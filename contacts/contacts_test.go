package contacts

import "github.com/LucasFrezarini/go-contacts/contacts/email"

var contactsList = []*Contact{
	&Contact{1, "Inosuke", "Hashibira", nil},
	&Contact{2, "Gonpachiro", "Kamaboko", nil},
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

var emailsList = []*email.Email{
	&email.Email{ID: 1, ContactID: 1, Address: "inosuke@gmail.com"},
	&email.Email{ID: 2, ContactID: 1, Address: "pigassault@outlook.com"},
	&email.Email{ID: 3, ContactID: 2, Address: "tanjirou@gmail.com"},
}

func filterEmailsByContactID(id int) []*email.Email {
	filteredEmails := make([]*email.Email, 0)

	for _, e := range emailsList {
		if e.ContactID == id {
			filteredEmails = append(filteredEmails, e)
		}
	}

	return filteredEmails
}

type MockedEmailsRepository struct{}

func (m *MockedEmailsRepository) FindByContactID(id int) ([]*email.Email, error) {
	return filterEmailsByContactID(id), nil
}
