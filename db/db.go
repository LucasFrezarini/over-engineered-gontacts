package db

import (
	"database/sql"
	"fmt"

	"github.com/LucasFrezarini/go-contacts/env"
	"github.com/google/wire"
	"go.uber.org/zap"
)

// ProvideDB opens a sql.DB connection that will be used in the whole project
func ProvideDB(logger *zap.Logger) (*sql.DB, error) {
	l := logger.Named("ProvideDB")
	e := env.GetEnvironment()

	uri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", e.MySQL.User, e.MySQL.Password, e.MySQL.Host, e.MySQL.Port, e.MySQL.Database)
	l.Info("opening connection to MySQL with URI " + uri)

	db, err := sql.Open("mysql", uri)
	if err != nil {
		return nil, fmt.Errorf("ProvideDB: error while creating sql.Conn: %w", err)
	}

	l.Info("ProvideDB: connection openned successfully")
	return db, nil
}

// DBSet is the wire.ProviderSet that represents this package
var DBSet = wire.NewSet(ProvideDB)
