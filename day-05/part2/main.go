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

	exists, err := client.CollectionExists(context.Background(), "day4_part2")

	if err != nil {
		panic(err)
	}

	fmt.Println("Collection Exists:", exists)
}
