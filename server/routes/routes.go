package routes

import (
	"github.com/LucasFrezarini/go-contacts/contacts"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// A Router provides functions to build the server routing, based on Go *http.ServeMux
type Router struct {
	contactsController *contacts.Controller
	logger             *zap.Logger
	echo               *echo.Echo
}

// BuildRouter initialize all routing groups of the server
func (r *Router) BuildRouter() {
	r.contactsController.EchoGroup()
}

// ProvideRouter is responsible by building the Router object. Designed especially for the use of
// wire, to provide the dependencies via DI
func ProvideRouter(cc *contacts.Controller, logger *zap.Logger, echo *echo.Echo) *Router {
	return &Router{contactsController: cc, logger: logger, echo: echo}
}
