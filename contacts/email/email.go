package email

// Email represents a contact's email
type Email struct {
	ID        int    `json:"id,omitempty"`
	ContactID int    `json:"contact_id"`
	Address   string `json:"address"`
}
