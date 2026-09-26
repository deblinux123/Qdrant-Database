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
		log.Fatalf("Connection failed: %v", err)
	}

	collectionName := "documents"

	err = client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     768,
			Distance: qdrant.Distance_Cosine,
		}),
	})

	if err != nil {
		log.Fatalf("Creating collection failed: %v", err)
	}

	fmt.Printf("Collection %s created successfully.", collectionName)

	info, err := client.GetCollectionInfo(ctx, collectionName)

	if err != nil {
		log.Fatal("Failed to get information.")
	}

	fmt.Printf("Collection info: %v\n", info)
}
