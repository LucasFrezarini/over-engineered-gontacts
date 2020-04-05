package contacts

import (
	"github.com/LucasFrezarini/go-contacts/contacts/email"
)

type Contact struct {
	ID        int            `json:"id,omitempty"`
	FirstName string         `json:"first_name" validate:"required"`
	LastName  string         `json:"last_name" validate:"required"`
	Emails    []*email.Email `json:"emails"`
}
