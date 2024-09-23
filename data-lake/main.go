package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"data-lake/services"
)

func main() {
	app := fiber.New()

	// Endpoint to export data from PostgreSQL and upload it to Azure Data Lake
	app.Get("/export", func(c *fiber.Ctx) error {
		err := services.ExportAndUpload()
		if err != nil {
			log.Printf("Failed to export and upload data: %v", err)
			return c.Status(fiber.StatusInternalServerError).SendString("Data export failed")
		}
		return c.SendString("Data exported and uploaded successfully")
	})

	log.Fatal(app.Listen(":3000"))
}
