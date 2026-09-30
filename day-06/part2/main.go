package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/qdrant/go-client/qdrant"
)

type EmbedRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type EmbedResponse struct {
	Embeddings [][]float32 `json:"embeddings"`
}

func getEmbedding(ctx context.Context, text string) ([]float32, error) {
	requestBody := EmbedRequest{
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
			"Ollama returned status: %s",
			response.Status,
		)
	}

	var result EmbedResponse

	err = json.NewDecoder(response.Body).Decode(&result)

	if err != nil {
		return nil, err
	}

	if len(result.Embeddings) == 0 {
		return nil, fmt.Errorf("no embedding returned")
	}

	return result.Embeddings[0], nil
}

func newClient() (*qdrant.Client, error) {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})

	if err != nil {
		return nil, fmt.Errorf("Qdrant connection failed: %w", err)
	}

	return client, nil
}

func createCollection(ctx context.Context) error {
	client, err := newClient()

	if err != nil {
		return err
	}

	collection := "day06_part2"

	exists, err := client.CollectionExists(ctx, collection)

	if err != nil {
		return fmt.Errorf(
			"failed to check collection %s: %w",
			collection,
			err,
		)
	}

	if !exists {
		err := client.CreateCollection(
			ctx,
			&qdrant.CreateCollection{
				CollectionName: collection,
				VectorsConfig: qdrant.NewVectorsConfig(
					&qdrant.VectorParams{
						Size:     768,
						Distance: qdrant.Distance_Cosine,
					},
				),
			},
		)

		if err != nil {
			return fmt.Errorf(
				"failed to create collection %s: %w",
				collection,
				err,
			)
		}

		fmt.Printf("Collection %s created.\n", collection)
	} else {
		fmt.Printf(
			"Collection %s already exists.\n",
			collection,
		)
	}

	return nil
}

func upsertPoint(
	ctx context.Context,
	vector []float32,
	text string,
) error {

	client, err := newClient()

	if err != nil {
		return err
	}

	point := &qdrant.PointStruct{
		Id: qdrant.NewIDNum(1),

		Vectors: qdrant.NewVectors(
			vector...,
		),

		Payload: qdrant.NewValueMap(
			map[string]any{
				"text": text,
			},
		),
	}

	response, err := client.Upsert(
		ctx,
		&qdrant.UpsertPoints{
			CollectionName: "day06_part1",

			Points: []*qdrant.PointStruct{
				point,
			},
		},
	)

	if err != nil {
		return fmt.Errorf(
			"failed to upsert point: %w",
			err,
		)
	}

	fmt.Println(
		"Upsert status:",
		response.GetStatus(),
	)

	return nil
}

func main() {
	ctx := context.Background()

	text := "Learning golang with qdrant"

	vector, err := getEmbedding(
		ctx,
		text,
	)

	if err != nil {
		panic(err)
	}

	fmt.Println(
		"Vector Size:",
		len(vector),
	)

	fmt.Println(
		"First 5 Values:",
		vector[:5],
	)

	err = createCollection(ctx)

	if err != nil {
		fmt.Printf(
			"Error: %v\n",
			err,
		)

		return
	}

	err = upsertPoint(
		ctx,
		vector,
		text,
	)

	if err != nil {
		fmt.Printf(
			"Error: %v\n",
			err,
		)

		return
	}

	fmt.Println("Point inserted successfully.")
}
