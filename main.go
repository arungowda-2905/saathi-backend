package main

import (
	"log"

	"github.com/gofiber/fiber/v2"

	"saathi-backend/config"
	"saathi-backend/routes"
	"saathi-backend/test_data"
)

func main() {

	err := config.ConnectDB()
	if err != nil {
		log.Fatal(err)
	}

	defer config.DisconnectDB()

	err = test_data.SeedTutorials()
	if err != nil {
		log.Fatal(err)
	}

	app := fiber.New()
	routes.SetupRoutes(app)
	log.Fatal(app.Listen(":3000"))
}
