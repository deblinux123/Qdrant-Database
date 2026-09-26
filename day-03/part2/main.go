package main

import (
	"context"
	"fmt"
	"log"

	"github.com/qdrant/go-client/qdrant"
)

func main() {
	ctx := context.Background()

	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})

	if err != nil {
		log.Fatal("Failed to connect.")
	}

	collections := []struct {
		Name     string
		Distance qdrant.Distance
	}{
		{"document_cosine", qdrant.Distance_Cosine},
		{"document_dot", qdrant.Distance_Dot},
		{"document_euclid", qdrant.Distance_Euclid},
	}

	for _, c := range collections {
		err := client.CreateCollection(ctx, &qdrant.CreateCollection{
			CollectionName: c.Name,
			VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
				Size:     768,
				Distance: c.Distance,
			}),
		})

		if err != nil {
			log.Fatalf("Failed to create collection: %v", err)
			continue
		}

		fmt.Printf("Collection %s created.\n", c.Name)
	}
}
