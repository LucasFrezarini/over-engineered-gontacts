//+build wireinject

package container

import (
	"github.com/LucasFrezarini/go-contacts/server"
	_ "github.com/go-sql-driver/mysql"
	"github.com/google/wire"
)

func InitializeServer() (*server.Server, error) {
	wire.Build(server.ServerSet)
	return &server.Server{}, nil
}
