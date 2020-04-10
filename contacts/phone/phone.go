package phone

import "github.com/google/wire"

// Available phone types
const (
	PhoneTypeMobile = "mobile"
	PhoneTypeHome   = "home"
	PhoneTypeWork   = "work"
	PhoneTypeFax    = "fax"
)

// Phone represents a contact's phone entry
type Phone struct {
	ID        int    `json:"id,omitempty"`
	ContactID int    `json:"contact_id"`
	Type      string `json:"type"`
	Number    string `json:"number"`
}

// Set is a Wire set that contains all the providers for this package
var Set = wire.NewSet(RepositorySet)
