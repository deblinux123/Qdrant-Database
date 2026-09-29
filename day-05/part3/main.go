package main

import (
	"context"
	"fmt"

	"github.com/qdrant/go-client/qdrant"
)

func main() {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})

	if err != nil {
		fmt.Println("Error:", err)
	}

	ctx := context.Background()

	collection := "day05_part3"

	exists, err := client.CollectionExists(ctx, collection)

	if err != nil {
		fmt.Println("Error:", err)
	}

	if !exists {
		err = client.CreateCollection(ctx, &qdrant.CreateCollection{
			CollectionName: collection,
			VectorsConfig: qdrant.NewVectorsConfig(
				&qdrant.VectorParams{
					Size:     3,
					Distance: qdrant.Distance_Cosine,
				},
			),
		})

		if err != nil {
			fmt.Println("Error:", err)
		}

		fmt.Println("Collection created.")
	}

	fmt.Println("Ready")
}
