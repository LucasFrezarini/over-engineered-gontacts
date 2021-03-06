package contacts

import (
	"github.com/LucasFrezarini/go-contacts/contacts/email"
	"github.com/LucasFrezarini/go-contacts/contacts/phone"
	"github.com/google/wire"
)

type Contact struct {
	ID        int           `json:"id,omitempty"`
	FirstName string        `json:"first_name" validate:"required"`
	LastName  string        `json:"last_name" validate:"required"`
	Emails    []email.Email `json:"emails"`
	Phones    []phone.Phone `json:"phones"`
}

// Set is a set that contains all the Wire providers from this package
var Set = wire.NewSet(
	ControllerSet,
	ServiceSet,
	RepositorySet,
	email.Set,
	phone.Set,
)
