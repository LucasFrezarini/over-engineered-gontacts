package contacts

import (
	"github.com/LucasFrezarini/go-contacts/db"
	"github.com/google/wire"
)

type Contact struct {
	ID        int    `json:"id,omitempty"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
}

var ContactSet = wire.NewSet(ProvideContactsController, RepositorySet, db.DBSet)
