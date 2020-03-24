package contacts

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// A Controller is responsible by providing HTTP handlers to any method envolving the "contact" resource
// in the application
type Controller struct {
	repository Repository
	logger     *zap.Logger
	echo       *echo.Echo
}

// ProvideContactsController is responsible by building a ContactsController object. Designed especially for the use of
// wire, to provide the dependencies via DI
func ProvideContactsController(r Repository, logger *zap.Logger, echo *echo.Echo) *Controller {
	return &Controller{repository: r, logger: logger.Named("ContactsController"), echo: echo}
}

// FindAll searches all the contacts that exists in the database and returns it
// in a JSON response
func (ct *Controller) FindAll(c echo.Context) error {
	contacts, err := ct.repository.FindAll()

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
func (ct *Controller) Create(c echo.Context) error {
	body := new(Contact)
	if err := c.Bind(body); err != nil {
		return err
	}

	created, err := ct.repository.Create(*body)
	if err != nil {
		return err
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
