package main

import (
	"data-lake/services"
	"log"

	"github.com/gofiber/fiber/v2"
	
)

func main() {
	app := fiber.New()

	app.Post("/loign", services.Login)
	app.Get("/secure-data", services.SecureData)

	log.Fatal(app.Listen(":3000"))
}
