package env

import (
	"os"
	"sync"
)

// MySQLEnvironment defines the structure of the MySQL env variables that are used across this project
type MySQLEnvironment struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
}

// Environment is a struct that defines all environment variables that are used across this project
type Environment struct {
	MySQL MySQLEnvironment
}

var environment Environment
var once sync.Once

// GetEnvironment returns a Environment singleton
func GetEnvironment() Environment {
	once.Do(func() {
		mySQLEnv := MySQLEnvironment{
			Host:     os.Getenv("MYSQL_HOST"),
			Port:     os.Getenv("MYSQL_PORT"),
			User:     os.Getenv("MYSQL_USER"),
			Password: os.Getenv("MYSQL_PASSWORD"),
			Database: os.Getenv("MYSQL_DATABASE"),
		}

		environment = Environment{MySQL: mySQLEnv}
	})

	return environment
}
