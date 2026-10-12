package main

import (
	"day-13/config"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. using system envirnoment")
	}

	cfg := config.Load()

	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "day-13",
		})
	})

	log.Println("Server starting on: 8080")

	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
