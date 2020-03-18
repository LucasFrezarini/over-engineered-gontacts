package server

import (
	"net/http"

	"github.com/LucasFrezarini/go-contacts/contacts"
	"github.com/LucasFrezarini/go-contacts/logger"
	"github.com/LucasFrezarini/go-contacts/server/routes"
	"github.com/google/wire"
	"go.uber.org/zap"
)

// A Server provides
type Server struct {
	Router *routes.Router
	Logger *zap.Logger
}

// Start starts an go native http server. It returns the results of the http.ListenAndServer function,
// so it will ever return an error when the server goes down.
func (s *Server) Start() error {
	mux := s.Router.BuildMux()
	s.Logger.Info("Starting HTTP server on port 8080...")
	return http.ListenAndServe(":8080", mux)
}

// ProvideServer provides a Server object, built for the use of wire
func ProvideServer(r *routes.Router, logger *zap.Logger) *Server {
	return &Server{Router: r, Logger: logger.Named("Server")}
}

// ServerSet is the wire.ProviderSet of the server package
var ServerSet = wire.NewSet(ProvideServer, routes.ProvideRouter, contacts.ContactSet, logger.LoggerSet)
