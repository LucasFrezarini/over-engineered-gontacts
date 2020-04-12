package contacts

import (
	"reflect"
	"testing"

	"go.uber.org/zap"
)

func TestServiceFindAllContacts(t *testing.T) {
	service := ProvideContactsService(
		zap.NewNop(),
		&MockedContactsRepository{},
		&MockedEmailRepository{},
		&MockedPhoneRepository{},
	)

	contacts, err := service.FindAllContacts()

	if err != nil {
		t.Errorf("FindAllContacts() returned a non-nil error '%v', want nil", err)
	}

	if expected, got := len(contactsList), len(contacts); expected != got {
		t.Errorf("FindAllContacts() returned %d contacts, want %d", got, expected)
	}

	for i, c := range contacts {
		expectedEmails := filterEmailsByContactID(c.ID)

		if expected, got := len(expectedEmails), len(c.Emails); expected != got {
			t.Errorf("FindAllContacts() contact[%d] has %d emails, want %d", i, got, expected)
		}

		for ie, e := range c.Emails {
			if expected, got := expectedEmails[ie], e; !reflect.DeepEqual(expected, got) {
				t.Errorf("FindAllContacts() contact[%d].emails[%d] == '%v', want '%v'", i, ie, got, expected)
			}
		}

		expectedPhones := filterPhonesByContactID(c.ID)

		if expected, got := len(expectedPhones), len(c.Phones); expected != got {
			t.Errorf("FindAllContacts() contact[%d] has %d phones, want %d", i, got, expected)
		}

		for iph, p := range c.Phones {
			if expected, got := expectedPhones[iph], p; !reflect.DeepEqual(expected, got) {
				t.Errorf("FindAllContacts() contact[%d].phones[%d] == '%v', want '%v'", i, iph, got, expected)
			}
		}
	}
}

func TestServiceCreate(t *testing.T) {
	c := CreateContactData{
		FirstName: "Zenitsu",
		LastName:  "Agatsuma",
		Emails:    []string{"zenitsu01@gmail.com", "zenitsu02@gmail.com"},
	}

	service := ProvideContactsService(
		zap.NewNop(),
		&MockedContactsRepository{},
		&MockedEmailRepository{},
		&MockedPhoneRepository{},
	)

	contact, err := service.Create(c)
	if err != nil {
		t.Errorf("Create(%v) returned an non-nil error: '%v', want nil", c, err)
	}

	if contact.ID == 0 {
		t.Errorf("Create(%v) contact.ID == 0, want != 0", c)
	}

	if expected, got := c.FirstName, contact.FirstName; expected != got {
		t.Errorf("Create(%v) contact.FirstName == %s, want %s", c, got, expected)
	}

	if expected, got := c.LastName, contact.LastName; expected != got {
		t.Errorf("Create(%v) contact.LastName == %s, want %s", c, got, expected)
	}

	if expected, got := len(c.Emails), len(contact.Emails); expected != got {
		t.Errorf("Create(%v) contact has %d emails, want %d", c, got, expected)
	}
}
