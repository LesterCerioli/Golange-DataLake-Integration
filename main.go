package main

import (
	"log"
	"os"

	"data-lake/handlers"
	"yapp/config"

	_ "github.com/denisenkom/go-mssqldb"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Carregar variáveis de ambiente do .env
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	// Inicializar a conexão com o banco de dados
	db, err := config.InitDB()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Configuração do Fiber
	app := fiber.New()

	// Rota para pegar os pagamentos e exportar para CSV
	app.Get("/export-payments", func(c *fiber.Ctx) error {
		err := handlers.ExportPayments(db)
		if err != nil {
			return c.Status(500).SendString("Failed to export payments: " + err.Error())
		}
		return c.SendString("Payments exported successfully!")
	})

	// Iniciar o servidor
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000" // Porta padrão se não definida
	}
	log.Fatal(app.Listen(":" + port))
}
