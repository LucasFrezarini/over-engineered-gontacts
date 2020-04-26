package contacts

import (
	"github.com/LucasFrezarini/go-contacts/contacts/email"
	"github.com/LucasFrezarini/go-contacts/contacts/phone"
	"go.uber.org/zap"
)

var contactsList = []*Contact{
	{
		ID:        1,
		FirstName: "Inosuke",
		LastName:  "Hashibira",
	},
	{
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

func (m *MockedContactsRepository) DeleteByID(id int) error {
	return nil
}

var emailsList = []email.Email{
	{ID: 1, ContactID: 1, Address: "inosuke@gmail.com"},
	{ID: 2, ContactID: 1, Address: "pigassault@outlook.com"},
	{ID: 3, ContactID: 2, Address: "tanjirou@gmail.com"},
}

var phonesList = []phone.Phone{
	{ID: 1, ContactID: 1, Number: "1122223333", Type: phone.PhoneTypeHome},
	{ID: 2, ContactID: 1, Number: "33444445555", Type: phone.PhoneTypeMobile},
	{ID: 3, ContactID: 1, Number: "+5511911112222", Type: phone.PhoneTypeWork},
	{ID: 4, ContactID: 1, Number: "+5511933332222", Type: phone.PhoneTypeFax},
	{ID: 5, ContactID: 2, Number: "11955554444", Type: phone.PhoneTypeMobile},
	{ID: 6, ContactID: 2, Number: "1122223333", Type: phone.PhoneTypeHome},
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

type MockedEmailRepository struct {
	id int
}

func (m *MockedEmailRepository) FindByContactID(id int) ([]email.Email, error) {
	return filterEmailsByContactID(id), nil
}

func (m *MockedEmailRepository) Create(contactID int, emails ...string) ([]email.Email, error) {
	parsed := make([]email.Email, 0, len(emails))

	for _, e := range emails {
		m.id++
		parsed = append(parsed, email.Email{
			ID:        m.id,
			Address:   e,
			ContactID: contactID,
		})
	}

	return parsed, nil
}

type MockedPhoneRepository struct {
	id int
}

func (pr *MockedPhoneRepository) FindByContactID(id int) ([]phone.Phone, error) {
	return filterPhonesByContactID(id), nil
}

func (pr *MockedPhoneRepository) Create(contactID int, phones ...phone.CreatePhoneData) ([]phone.Phone, error) {
	parsed := make([]phone.Phone, 0, len(phones))

	for _, p := range phones {
		pr.id++
		parsed = append(parsed, phone.Phone{
			ID:        pr.id,
			ContactID: contactID,
			Type:      p.Type,
			Number:    p.Number,
		})
	}

	return parsed, nil
}

func ProvideContactMockedService() *Service {
	return ProvideContactsService(
		zap.NewNop(),
		&MockedContactsRepository{},
		&MockedEmailRepository{},
		&MockedPhoneRepository{},
	)
}
