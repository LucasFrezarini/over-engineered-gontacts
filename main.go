package main

import (
	"log"

	"github.com/LucasFrezarini/go-contacts/container"
	_ "github.com/joho/godotenv/autoload"
)

func main() {
	app, err := container.InitializeServer()
	if err != nil {
		log.Fatal(err)
	}

	defer app.Logger.Sync()

	log.Fatal(app.Start())

	app.Logger.Info("Closing application...")
}
