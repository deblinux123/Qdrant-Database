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
	qdrantHost            = "localhost"
	qdrantPort            = 6334
	ollamEmbedURL         = "http://localhost:11434/api/embed"
	embeddingModel        = "embeddinggemma:300m"
	collectionName        = "day07_payload"
	vectorSize     uint64 = 768
)

type EmbedRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type EmbedResponse struct {
	Embeddings [][]float32 `json:"embeddings"`
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

func getEmbeddings(ctx context.Context, text string) ([]float32, error) {
	requestBody := EmbedRequest{
		Model: embeddingModel,
		Input: text,
	}

	body, err := json.Marshal(requestBody)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to marshal embedding request: %w",
			err,
		)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		ollamEmbedURL,
		bytes.NewReader(body),
	)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to create request: %w",
			err,
		)
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
			"Ollama returned status: %w",
			response.Status,
		)
	}

	var result EmbedResponse

	err = json.NewDecoder(response.Body).Decode(&result)

	if err != nil {
		return nil, fmt.Errorf(
			"failed to decode ollama response: %w",
			err,
		)
	}

	if len(result.Embeddings) == 0 {
		return nil, fmt.Errorf(
			"no embeddings returned: %w",
			err,
		)
	}

	return result.Embeddings[0], nil
}

func createCollection(ctx context.Context, client *qdrant.Client) error {
	exists, err := client.CollectionExists(ctx, collectionName)

	if err != nil {
		return fmt.Errorf(
			"failed to check collection: %w",
			err,
		)
	}

	if exists {
		fmt.Printf(
			"Collection %s already exists.\n",
			collectionName,
		)

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
		return fmt.Errorf(
			"failed to create collection: %w",
			err,
		)
	}

	fmt.Printf(
		"collection %s created.\n",
		collectionName,
	)

	return nil
}

func upsertPoint(ctx context.Context, client *qdrant.Client, vector []float32) error {
	point := &qdrant.PointStruct{
		Id:      qdrant.NewIDNum(1),
		Vectors: qdrant.NewVectors(vector...),
		Payload: qdrant.NewValueMap(map[string]any{
			"title":    "golang course",
			"author":   "John",
			"category": "programming",
			"language": "en",
			"year":     2026,
		}),
	}

	response, err := client.Upsert(ctx, &qdrant.UpsertPoints{
		CollectionName: collectionName,
		Points: []*qdrant.PointStruct{
			point,
		},
	})

	if err != nil {
		return fmt.Errorf(
			"failed to upsert point: %w",
			err,
		)
	}

	fmt.Printf(
		"Upser status: %v\n",
		response.GetStatus(),
	)

	return nil
}
func main() {
	ctx := context.Background()

	text := "Learning golang programming with qdrant."

	fmt.Print("creating qdrant client...")

	client, err := newClient()

	if err != nil {
		panic(err)
	}

	fmt.Println("creating collection...")

	err = createCollection(ctx, client)

	if err != nil {
		panic(err)
	}

	fmt.Println("getting embedding from ollama...")

	vector, err := getEmbeddings(ctx, text)

	if err != nil {
		panic(err)
	}

	fmt.Println("Vector size:", len(vector))

	fmt.Println("first 5 values:", vector[:5])

	err = upsertPoint(ctx, client, vector)

	if err != nil {
		panic(err)
	}

	fmt.Println("Point inserted successfully.")
}
