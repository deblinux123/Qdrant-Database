package main

import (
	"context"
	"fmt"

	"github.com/qdrant/go-client/qdrant"
)

func createCollection() error {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})

	if err != nil {
		return err
	}

	err = client.CreateCollection(context.Background(), &qdrant.CreateCollection{
		CollectionName: "day04_part2",
		VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
			Size:     3,
			Distance: qdrant.Distance_Cosine,
		}),
	})

	if err != nil {
		return err
	}

	fmt.Println("Collection created.")

	return nil
}

func collectionExists() bool {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})

	if err != nil {
		return false
	}

	exists, err := client.CollectionExists(context.Background(), "day04_part2")

	if err != nil {
		return false
	}

	fmt.Println("Collection exists:", exists)

	return exists
}

func upsert() {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})

	if err != nil {
		panic(err)
	}

	response, err := client.Upsert(context.Background(), &qdrant.UpsertPoints{
		CollectionName: "day04_part2",
		Points: []*qdrant.PointStruct{
			{
				Id:      qdrant.NewIDNum(1),
				Vectors: qdrant.NewVectors(0.9, 0.1, 0.1),
				Payload: qdrant.NewValueMap(map[string]any{
					"color": "red",
				}),
			},
			{
				Id:      qdrant.NewIDNum(2),
				Vectors: qdrant.NewVectors(0.1, 0.9, 0.1),
				Payload: qdrant.NewValueMap(map[string]any{
					"color": "green",
				}),
			},
		},
	})

	if err != nil {
		panic(err)
	}

	fmt.Println("Upsert status:", response.GetStatus())
}

func main() {
	exists := collectionExists()

	if !exists {
		err := createCollection()

		if err != nil {
			panic(err)
		}
	}

	upsert()
}
