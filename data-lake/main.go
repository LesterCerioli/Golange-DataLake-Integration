package main

import (
	"data-lake/services"
	"log"

	"github.com/gofiber/fiber/v2"
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"time"
	"github.com/Azure/azure-sdk-for-go/sdk/storage/azdatalake"
	_ "github.com/jackc/pgx/v4/stdlib" 

	
)

const (
	accountName = "your-account-name"
	accountKey = "your-account-key"
	fileSystem    = "your-file-system-name" 
	azureFilePath = "exports/data.csv" 
)

const (
	postgresDSN = "postgres://user:password@localhost:5432/payments_db"
)

func main() {
	app := fiber.New()

	app.Post("/loign", services.Login)
	app.Get("/secure-data", services.SecureData)

	log.Fatal(app.Listen(":3000"))
}
