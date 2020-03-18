package contacts

import (
	"encoding/json"
	"fmt"
	"net/http"

	"go.uber.org/zap"
)

// A Controller is responsible by providing HTTP handlers to any method envolving the "contact" resource
// in the application
type Controller struct {
	repository Repository
	logger     *zap.Logger
}

// ProvideContactsController is responsible by building a ContactsController object. Designed especially for the use of
// wire, to provide the dependencies via DI
func ProvideContactsController(r Repository, logger *zap.Logger) *Controller {
	return &Controller{repository: r, logger: logger.Named("ContactsController")}
}

// FindAll searches all the contacts that exists in the database and returns it
// in a JSON response
func (c *Controller) FindAll(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content-Type", "application/json")
	c.logger.Info("Received GET /contacts/")

	contacts, err := c.repository.FindAll()

	if err != nil {
		c.logger.Error(fmt.Sprintf("GET / internal server error: %v", err))
		w.WriteHeader(500)
		return
	}

	var response = struct {
		Contacts []*Contact `json:"contacts"`
	}{contacts}

	b, err := json.Marshal(response)
	if err != nil {
		if err != nil {
			c.logger.Error(fmt.Sprintf("GET / internal server error: %v", err))
			w.WriteHeader(500)
			return
		}
	}

	w.WriteHeader(200)
	w.Write(b)
}

// Mux is responsible for building a server mux for this controller
func (c *Controller) Mux() *http.ServeMux {
	c.logger.Debug("Building the ContactsController ServeMux...")
	mux := http.NewServeMux()
	mux.HandleFunc("/", c.FindAll)

	return mux
}
