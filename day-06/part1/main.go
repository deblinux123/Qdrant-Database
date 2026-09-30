package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/qdrant/go-client/qdrant"
)

type EmbeddRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type EmbeddResponse struct {
	Embeddings [][]float32 `json:"embeddings"`
}

func getEmbedding(ctx context.Context, text string) ([]float32, error) {
	requestBody := EmbeddRequest{
		Model: "embeddinggemma:300m",
		Input: text,
	}

	body, err := json.Marshal(requestBody)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		"http://localhost:11434/api/embed",
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
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf(
			"ollama returned status: %s",
			response.Status,
		)
	}

	var result EmbeddResponse

	err = json.NewDecoder(response.Body).Decode(&result)

	if err != nil {
		return nil, err
	}

	if len(result.Embeddings) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}

	return result.Embeddings[0], nil
}

func createCollection(ctx context.Context) error {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})

	if err != nil {
		return fmt.Errorf("Qdrant Connection failed: %v", err)
	}

	collection := "day06_part1"

	exists, err := client.CollectionExists(ctx, collection)

	if err != nil {
		return fmt.Errorf("Failed to check collection %s: %w", collection, err)
	}

	if !exists {
		err = client.CreateCollection(ctx, &qdrant.CreateCollection{
			CollectionName: collection,
			VectorsConfig: qdrant.NewVectorsConfig(&qdrant.VectorParams{
				Size:     768,
				Distance: qdrant.Distance_Cosine,
			}),
		})

		if err != nil {
			return fmt.Errorf("Error: %v", err)
		}

		fmt.Printf("Collectoin %s created.\n", collection)
	} else {
		fmt.Printf("Collection %s already exists.\n", collection)
	}

	return nil
}

func Upsert(ctx context.Context, vector []float32, text string) error {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})

	if err != nil {
		return fmt.Errorf("Qdrant connection failed: %w", err)
	}

	point := &qdrant.PointStruct{
		Id:      qdrant.NewIDNum(1),
		Vectors: qdrant.NewVectors(vector...),
		Payload: qdrant.NewValueMap(map[string]any{
			"text": text,
		}),
	}

	response, err := client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: "day06_part1",
		Points: []*qdrant.PointStruct{
			point,
		},
	})

	if err != nil {
		return fmt.Errorf("Failed to upsert point: %w", err)
	}

	fmt.Println("Upsert status:", response.GetStatus())

	return nil
}

func main() {
	ctx := context.Background()

	vector, err := getEmbedding(ctx, "Learning golang with qdrant")

	if err != nil {
		panic(err)
	}

	fmt.Println("Vector Size:", len(vector))
	fmt.Println("First 5 Values:", vector[:5])
	err = createCollection(ctx)

	if err != nil {
		fmt.Printf("Error: %v", err)
	}

	err = Upsert(ctx, vector, "This is the best course to learn qdrant with golang.")

	if err != nil {
		fmt.Printf("Error: %v\n", err)
	}
}
