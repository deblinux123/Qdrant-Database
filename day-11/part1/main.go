package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/google/uuid"
	"github.com/qdrant/go-client/qdrant"
)

const (
	ollamURL       = "http://localhost:11434/api/embed"
	modelName      = "embeddinggemma:300m"
	qdrantHost     = "localhost"
	qdrantPort     = 6334
	vectorSize     = uint64(768)
	collectionName = "day11_embeddings"
)

type EmbedRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type EmbedResponse struct {
	Embeddings [][]float32 `jaon:"embeddings"`
}

func getEmbedding(ctx context.Context, text string) ([]float32, error) {
	embedRequest := EmbedRequest{
		Model: modelName,
		Input: text,
	}

	body, err := json.Marshal(embedRequest)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		ollamURL,
		bytes.NewReader(body),
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	resp, err := client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama returned status: %s", resp.Status)
	}

	var result EmbedResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("no embedding returned")
	}

	return result.Embeddings[0], nil
}

func newClient() (*qdrant.Client, error) {
	return qdrant.NewClient(&qdrant.Config{
		Host: qdrantHost,
		Port: qdrantPort,
	})
}

func createCollection(ctx context.Context, client *qdrant.Client) error {
	exits, err := client.CollectionExists(ctx, collectionName)

	if err != nil {
		return fmt.Errorf("failed to check collection %q: %w", collectionName, err)
	}

	if exits {
		fmt.Printf("Collection %q already exists\n", collectionName)
		return nil
	}

	err = client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(
			&qdrant.VectorParams{
				Size:     vectorSize,
				Distance: qdrant.Distance_Cosine,
			},
		),
	})

	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	fmt.Printf("collection %q creatd.\n", collectionName)

	return nil
}

func upsertPoint(ctx context.Context, client *qdrant.Client, text string, vector []float32) error {
	if len(vector) != int(vectorSize) {
		return fmt.Errorf("expected %d dimensions, got %d", vectorSize, len(vector))
	}

	_, err := client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collectionName,
		Points: []*qdrant.PointStruct{
			{
				Id:      qdrant.NewID(uuid.NewString()),
				Vectors: qdrant.NewVectors(vector...),
				Payload: qdrant.NewValueMap(map[string]any{
					"text":  text,
					"model": modelName,
				}),
			},
		},
	})

	if err != nil {
		return fmt.Errorf("failed to upsert points: %w", err)
	}

	return nil

}

func main() {
	fmt.Println("Working with ollama new model embeddinggemma")

	ctx := context.Background()

	text := "Building REST Api with golang"

	vector, err := getEmbedding(ctx, text)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Text:", text)
	fmt.Println("Model:", modelName)
	fmt.Println("Vector dimensions:", len(vector))
	fmt.Printf("First 10 values: %v\n", vector[:10])

	client, err := newClient()

	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	if err := createCollection(ctx, client); err != nil {
		fmt.Printf("Collection error: %v\n", err)
	}

	err = upsertPoint(ctx, client, text, vector)

	if err != nil {
		fmt.Println("failed to upsert point:", err)
		return
	}

}
