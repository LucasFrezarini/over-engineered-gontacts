package contacts

import (
	"errors"
	"fmt"

	"github.com/LucasFrezarini/go-contacts/contacts/email"
	"github.com/LucasFrezarini/go-contacts/contacts/phone"
	"github.com/google/wire"
	"go.uber.org/zap"
)

// A Service contains all the business logic related to the contact resource in the application
type Service struct {
	Logger             *zap.Logger
	ContactsRepository Repository
	EmailRepository    email.GenericRepository
	PhoneRepository    phone.GenericRepository
}

// ProvideContactsService creates a new Service with the provided dependencies.
// Created especially for the use of Wire, who will inject the dependencies via DI
func ProvideContactsService(logger *zap.Logger, cr Repository, er email.GenericRepository, pr phone.GenericRepository) *Service {
	return &Service{logger.Named("ContactsService"), cr, er, pr}
}

// FindAllContacts fetches all the contacts registered in the application, as well as its emails and phones
func (s *Service) FindAllContacts() ([]*Contact, error) {
	contacts, err := s.ContactsRepository.FindAll()
	if err != nil {
		msg := fmt.Sprintf("FindAllContacts() error while trying to fetch contacts: %v", err)
		s.Logger.Error(msg)
		return nil, errors.New(msg)
	}

	for _, c := range contacts {
		emails, err := s.EmailRepository.FindByContactID(c.ID)
		if err != nil {
			msg := fmt.Sprintf("FindAllContacts() error while trying to fetch contact's emails: %v", err)
			s.Logger.Error(msg)
			return nil, errors.New(msg)
		}

		c.Emails = emails

		phones, err := s.PhoneRepository.FindByContactID(c.ID)
		if err != nil {
			msg := fmt.Sprintf("FindAllContacts() error while trying to fetch contact's phones: %v", err)
			s.Logger.Error(msg)
			return nil, errors.New(msg)
		}

		c.Phones = phones
	}

	return contacts, nil
}

// CreateContactData is the structure of a contact data that will be created
type CreateContactData struct {
	FirstName string
	LastName  string
	Emails    []string
	Phones    []map[string]string
}

// Create creates a new contact with the data provided as parameter.
// If the contact is created successfully, it will return a formated Contact object
func (s *Service) Create(c CreateContactData) (*Contact, error) {
	contact, err := s.ContactsRepository.Create(Contact{
		FirstName: c.FirstName,
		LastName:  c.LastName,
	})

	if err != nil {
		msg := fmt.Sprintf("error while creating a new contact: %v", err)
		s.Logger.Error(msg)
		return nil, errors.New(msg)
	}

	if len(c.Emails) != 0 {
		emails, err := s.EmailRepository.Create(contact.ID, c.Emails...)
		if err != nil {
			msg := fmt.Sprintf("error while inserting contact's emails: %v", err)
			s.Logger.Error(msg)
			return nil, errors.New(msg)
		}

		contact.Emails = emails
	}

	return contact, nil
}

// ServiceSet is a wire set which contains all the bindings needed for creating a new service
var ServiceSet = wire.NewSet(ProvideContactsService)
