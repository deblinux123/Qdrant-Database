package main

import (
	"context"
	"log"

	"github.com/deblinux123/Qdrant-Database/day-13/config"
	"github.com/deblinux123/Qdrant-Database/day-13/ragservice"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

type SearchRequest struct {
	Query string `json:"query"`
	Limit uint64 `json:"limit"`
}

func main() {

	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found. using system envirnoment")
	}

	cfg := config.Load()
	ctx := context.Background()

	qdrantClient, err := ragservice.NewClient("papers")

	if err != nil {
		log.Fatal(err)
	}

	defer qdrantClient.Close()

	err = qdrantClient.CreateCollection(ctx)

	if err != nil {
		log.Fatal(err)
	}

	ollama := ragservice.NewOllamaService()

	app := fiber.New()

	app.Get("/health", func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{
			"status":  "ok",
			"service": "day-13",
		})
	})

	app.Post("/search", func(c *fiber.Ctx) error {
		var req SearchRequest

		if err := c.BodyParser(&req); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "invalid request body",
			})
		}

		if req.Query == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "query is required.",
			})
		}

		if req.Limit == 0 {
			req.Limit = 3
		}

		ctx := c.UserContext()

		vector, err := ollama.GetEmbedding(ctx, req.Query)

		if err != nil {
			log.Printf("Embedding error: %v", err)

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to generate embeddng.",
			})
		}

		points, err := qdrantClient.Search(ctx, vector, req.Limit)

		if err != nil {
			log.Printf("Qdrant search failed:%v", err)

			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
				"error": "failed to search any points",
			})
		}

		return c.JSON(fiber.Map{
			"query":  req.Query,
			"count":  len(points),
			"result": points,
		})
	})

	log.Println("Server starting on: 8080")

	if err := app.Listen(":" + cfg.Port); err != nil {
		log.Fatal(err)
	}
}
