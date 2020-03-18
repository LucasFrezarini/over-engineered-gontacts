package routes

import (
	"net/http"

	"github.com/LucasFrezarini/go-contacts/contacts"
	"go.uber.org/zap"
)

// A Router provides functions to build the server routing, based on Go *http.ServeMux
type Router struct {
	contactsController *contacts.Controller
	logger             *zap.Logger
}

// BuildMux creates a new *http.ServerMux with the routing of all controller from the application
func (r *Router) BuildMux() *http.ServeMux {
	r.logger.Debug("Building server mux...")
	mux := http.NewServeMux()
	mux.Handle("/contacts", r.contactsController.Mux())

	return mux
}

// ProvideRouter is responsible by building the Router object. Designed especially for the use of
// wire, to provide the dependencies via DI
func ProvideRouter(cc *contacts.Controller, logger *zap.Logger) *Router {
	return &Router{contactsController: cc, logger: logger}
}
