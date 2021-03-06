package server

import (
	"github.com/LucasFrezarini/go-contacts/contacts"
	"github.com/LucasFrezarini/go-contacts/db"
	"github.com/LucasFrezarini/go-contacts/logger"
	"github.com/LucasFrezarini/go-contacts/server/middlewares"
	"github.com/LucasFrezarini/go-contacts/server/routes"
	"github.com/LucasFrezarini/go-contacts/server/validator"
	"github.com/google/wire"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
)

// A Server contains all the artifacts necessary to start a fresh app's server instance
type Server struct {
	router *routes.Router
	Logger *zap.Logger
	echo   *echo.Echo
}

// Start starts an go native http server. It returns the results of the http.ListenAndServer function,
// so it will ever return an error when the server goes down.
func (s *Server) Start() error {
	s.router.BuildRouter()
	s.Logger.Info("starting HTTP server on port 8080...")
	return s.echo.Start(":8080")
}

// ProvideServer provides a Server object, built for the use of wire
func ProvideServer(r *routes.Router, logger *zap.Logger, echo *echo.Echo) *Server {
	return &Server{router: r, Logger: logger.Named("Server"), echo: echo}
}

// ProvideEcho provides a brand new echo instance
func ProvideEcho(middlewares *middlewares.Container) *echo.Echo {
	e := echo.New()
	e.Validator = validator.NewCustomValidator()

	e.Use(middlewares.ZapHTTPLogger)

	return e
}

// ServerSet is the wire.ProviderSet of the server package
var ServerSet = wire.NewSet(
	ProvideEcho,
	ProvideServer,
	middlewares.ProvideMiddlewaresContainer,
	routes.ProvideRouter,
	contacts.Set,
	logger.LoggerSet,
	db.DBSet,
)
