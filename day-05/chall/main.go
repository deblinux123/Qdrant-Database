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
		panic(err)
	}

	fmt.Println("Qdrant connected successfully.")

	ctx := context.Background()

	collection := "test1"

	exists, err := client.CollectionExists(ctx, collection)

	if err != nil {
		panic(err)
	}

	if !exists {
		err = client.CreateCollection(
			ctx,
			&qdrant.CreateCollection{
				CollectionName: collection,
				VectorsConfig: qdrant.NewVectorsConfig(
					&qdrant.VectorParams{
						Size:     3,
						Distance: qdrant.Distance_Cosine,
					},
				),
			},
		)

		if err != nil {
			panic(err)
		}

		fmt.Println("Collection:", collection)
	} else {
		fmt.Println("Collection already exists.")
	}

	fmt.Println("Ready for points")

	response, err := client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collection,
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
