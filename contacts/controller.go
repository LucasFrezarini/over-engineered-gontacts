package contacts

import (
	"fmt"
	"net/http"

	"github.com/LucasFrezarini/go-contacts/db"
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// A Controller is responsible by providing HTTP handlers to any method envolving the "contact" resource
// in the application
type Controller struct {
	repository Repository // TODO: move the create method to service as well
	service    *Service
	logger     *zap.Logger
	echo       *echo.Echo
}

// ProvideContactsController is responsible by building a ContactsController object. Designed especially for the use of
// wire, to provide the dependencies via DI
func ProvideContactsController(s *Service, r Repository, logger *zap.Logger, echo *echo.Echo) *Controller {
	return &Controller{service: s, repository: r, logger: logger.Named("ContactsController"), echo: echo}
}

// FindAll searches all the contacts that exists in the database and returns it
// in a JSON response
func (ct *Controller) FindAll(c echo.Context) error {
	contacts, err := ct.service.FindAllContacts()

	if err != nil {
		ct.logger.Error(fmt.Sprintf("GET / internal server error: %v", err))
		c.NoContent(http.StatusInternalServerError)
		return err
	}

	var response = struct {
		Contacts []*Contact `json:"contacts"`
	}{contacts}

	return c.JSON(http.StatusOK, response)
}

// Create creates a new contact with the info provided in the body
func (ct *Controller) Create(c echo.Context) (err error) {
	body := new(Contact)

	if err = c.Bind(body); err != nil {
		return
	}

	if err = c.Validate(body); err != nil {
		c.JSON(400, map[string]interface{}{
			"message": err.Error(),
		})

		return
	}

	created, err := ct.repository.Create(*body)
	if err != nil {
		return
	}

	return c.JSON(http.StatusCreated, created)
}

// EchoGroup is responsible for building an echo group with all routes for this controller
func (ct *Controller) EchoGroup() *echo.Group {
	ct.logger.Debug("Building the ContactsController routing group...")

	gp := ct.echo.Group("/contacts")
	gp.GET("/", ct.FindAll)
	gp.POST("/", ct.Create)

	return gp
}

// ControllerSet is a wire set which contains all the bindings needed for building all
// the resources present in this package
var ControllerSet = wire.NewSet(ProvideContactsController, ServiceSet, db.DBSet)
