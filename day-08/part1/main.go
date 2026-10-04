package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/qdrant/go-client/qdrant"
)

const (
	qdrantHost     = "localhost"
	qdrantPort     = 6334
	ollamaEmbedURL = "http://localhost:11434/api/embed"
	embeddingModel = "embeddinggemma:300m"
	collectionName = "day08_vector_search"
	vectorSize     = 768
)

type EmbedRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type EmbedResponse struct {
	Embeddings [][]float32
}

func newClient() (*qdrant.Client, error) {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: qdrantHost,
		Port: qdrantPort,
	})

	if err != nil {
		return nil, fmt.Errorf(
			"qdrant connection failed: %w",
			err,
		)
	}

	return client, nil
}

func getEmbedding(ctx context.Context, text string) ([]float32, error) {
	requestBody := EmbedRequest{
		Model: embeddingModel,
		Input: text,
	}

	body, err := json.Marshal(requestBody)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		ollamaEmbedURL,
		bytes.NewReader(body),
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	response, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf(
			"ollama request failed: %w",
			err,
		)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"ollama returned status: %s",
			response.Status,
		)
	}

	var result EmbedResponse

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("no embedding returned.")
	}

	return result.Embeddings[0], nil
}

func createCollection(ctx context.Context, client *qdrant.Client) error {
	exists, err := client.CollectionExists(ctx, collectionName)

	if err != nil {
		return err
	}

	if exists {
		fmt.Printf(
			"collection %q already exists.\n",
			collectionName,
		)
		return nil
	}

	err = client.CreateCollection(
		ctx,
		&qdrant.CreateCollection{
			CollectionName: collectionName,
			VectorsConfig: qdrant.NewVectorsConfig(
				&qdrant.VectorParams{
					Size:     vectorSize,
					Distance: qdrant.Distance_Cosine,
				},
			),
		},
	)

	if err != nil {
		return fmt.Errorf(
			"failed to create %q collection.",
			collectionName,
		)
	}

	fmt.Printf(
		"collection %q created.",
		collectionName,
	)

	return nil
}

func upsertPoint(ctx context.Context, client *qdrant.Client, id uint64, text string) error {
	fmt.Println(
		"Embedding:",
		text,
	)

	vector, err := getEmbedding(ctx, text)

	if err != nil {
		return fmt.Errorf(
			"embedding faile: %w",
			err,
		)
	}

	point := &qdrant.PointStruct{
		Id:      qdrant.NewIDNum(id),
		Vectors: qdrant.NewVectors(vector...),
		Payload: qdrant.NewValueMap(map[string]any{
			"text": text,
		}),
	}

	_, err = client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collectionName,
		Points: []*qdrant.PointStruct{
			point,
		},
	})

	if err != nil {
		return fmt.Errorf(
			"upsert failed: %w",
			err,
		)
	}

	fmt.Printf(
		"%q Inserted.",
		text,
	)

	return nil
}

func search(ctx context.Context, client *qdrant.Client, query string, limit uint64) error {
	fmt.Println()
	fmt.Println("query:", query)

	queryVector, err := getEmbedding(ctx, query)

	if err != nil {
		return fmt.Errorf(
			"query embedding failed: %w",
			err,
		)
	}

	fmt.Println("query vector size:", len(queryVector))

	results, err := client.Query(ctx, &qdrant.QueryPoints{
		CollectionName: collectionName,
		Query:          qdrant.NewQuery(queryVector...),
		Limit:          qdrant.PtrOf(limit),
		WithPayload:    qdrant.NewWithPayload(true),
	})

	if err != nil {
		return fmt.Errorf(
			"vector search failed: %w",
			err,
		)
	}

	fmt.Println()
	fmt.Println("Search Results")
	fmt.Println("====================")

	for i, result := range results {
		fmt.Printf("Rank: %d\n", i+1)

		fmt.Println("ID:", result.GetId())

		fmt.Println("Score:", result.GetScore())

		fmt.Println("Payload:", result.GetPayload())

		fmt.Println("---------------------")
	}

	return nil
}

func main() {
	ctx := context.Background()

	client, err := newClient()

	if err != nil {
		panic(err)
	}

	err = createCollection(ctx, client)

	if err != nil {
		panic(err)
	}

	docs := []string{
		"Build REST API with Go",
		"Golang HTTP Server",
		"Fiber Web Framework",
		"Go Backend Tutorial",
		"Building APIs with Python",
		"Crafting RESTful APIs in Go REST (Representational State Transfer) is all about making APIs intuitive and scalable, like a well-organized library. Let’s break down the key REST principles and build a simple user management API using the Gin framework.",
		"Implementing REST APIs Using Go’s Built-in Library (net/http) Go’s standard library includes the powerful net/http package, which is all you need for basic web servers and APIs. No external dependencies required! It handles HTTP requests, responses, routing (via a simple multiplexer), and more.",
	}

	fmt.Println()
	fmt.Println("Inserting documents...")
	fmt.Println("======================")

	for i, doc := range docs {
		err := upsertPoint(
			ctx,
			client,
			uint64(i+1),
			doc,
		)

		if err != nil {
			panic(err)
		}
	}

	query := "How to write rest api in go?"

	err = search(
		ctx,
		client,
		query,
		3,
	)

	if err != nil {
		panic(err)
	}

}
