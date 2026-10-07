package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/qdrant/go-client/qdrant"
)

type Document struct {
	ID       uint64
	Text     string
	Language string
	Category string
	Year     int
}

type EmbedRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type EmbedResponse struct {
	Embeddings [][]float32 `json:"embeddings"`
}

const (
	embedModel     = "embeddinggemma:300m"
	ollamaEmbedURL = "http://localhost:11434/api/embed"
	// collectionName = "day10_hnsw"
	collectionName = "papers"

	qdrantHost = "localhost"
	qdrantPort = 6334
	embedSize  = 768
	distance   = qdrant.Distance_Cosine
)

func main() {
	ctx := context.Background()

	documents := []Document{
		{
			ID:       1,
			Text:     "Build rest api with go",
			Language: "en",
			Category: "programming",
			Year:     2026,
		},
		{
			ID:       2,
			Text:     "ساخت REST API با زبان Go",
			Language: "fa",
			Category: "programming",
			Year:     2026,
		},
		{
			ID:       3,
			Text:     "Introduction to Artificial Intelligence",
			Language: "en",
			Category: "AI",
			Year:     2024,
		},
		{
			ID:       4,
			Text:     "مقدمه‌ای بر هوش مصنوعی",
			Language: "fa",
			Category: "AI",
			Year:     2025,
		},
		{
			ID:       5,
			Text:     "Building AI Agents with Python",
			Language: "en",
			Category: "AI",
			Year:     2026,
		},
	}

	client, err := newClient()

	if err != nil {
		log.Fatal(err)
	}

	defer client.Close()

	err = createCollection(ctx, client)

	if err != nil {
		log.Fatal(err)
	}

	for _, doc := range documents {
		vector, err := getEmbedding(ctx, doc.Text)

		if err != nil {
			log.Fatal(err)
		}

		err = upsertDocument(ctx, client, doc, vector)

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("Inserted:", doc.ID, doc.Text)
	}

	// query := "How to build an API with Go"
	query := "باشگاه پژوهشگران مهندسی توسعه"

	queryVector, err := getEmbedding(ctx, query)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("\nSearch tesults:")

	err = searchDocument(ctx, client, queryVector)

	if err != nil {
		log.Fatal(err)
	}
}

func newClient() (*qdrant.Client, error) {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: qdrantHost,
		Port: qdrantPort,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to connect qdrant: %w", err)
	}

	return client, nil
}

func createCollection(ctx context.Context, client *qdrant.Client) error {
	exists, err := client.CollectionExists(ctx, collectionName)

	if err != nil {
		return fmt.Errorf("failed to check collection: %w", err)
	}

	if exists {
		fmt.Printf("collection %q already exists.\n", collectionName)
		return nil
	}

	err = client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: collectionName,
		VectorsConfig: qdrant.NewVectorsConfig(
			&qdrant.VectorParams{
				Size:     embedSize,
				Distance: distance,
			},
		),
	})

	if err != nil {
		return fmt.Errorf("failed to create collection: %w", err)
	}

	fmt.Println("collection created:", collectionName)

	return nil
}

func getEmbedding(ctx context.Context, text string) ([]float32, error) {
	requestBody := EmbedRequest{
		Model: embedModel,
		Input: text,
	}

	body, err := json.Marshal(requestBody)

	if err != nil {
		return nil, fmt.Errorf("failed to encode request: %w", err)
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		ollamaEmbedURL,
		bytes.NewReader(body),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set(
		"ontent-Type",
		"application/json",
	)

	response, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	defer response.Body.Close()

	if response.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ollama returned status: %s", response.Status)
	}

	var result EmbedResponse

	if err := json.NewDecoder(response.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode embedding response: %w", err)
	}

	if len(result.Embeddings) == 0 {
		return nil, fmt.Errorf("no embedding returned.")
	}

	return result.Embeddings[0], nil
}

func upsertDocument(ctx context.Context, client *qdrant.Client, doc Document, vector []float32) error {

	payload := qdrant.NewValueMap(map[string]any{
		"text":     doc.Text,
		"language": doc.Language,
		"category": doc.Category,
		"year":     doc.Year,
	})

	point := &qdrant.PointStruct{
		Id:      qdrant.NewIDNum(doc.ID),
		Vectors: qdrant.NewVectors(vector...),
		Payload: payload,
	}

	_, err := client.Upsert(
		ctx,
		&qdrant.UpsertPoints{
			CollectionName: collectionName,
			Points: []*qdrant.PointStruct{
				point,
			},
		},
	)

	if err != nil {
		return fmt.Errorf(
			"failed to upsert documents: %d:%w",
			doc.ID,
			err,
		)
	}

	return nil
}

func searchDocument(ctx context.Context, client *qdrant.Client, queryVector []float32) error {
	start := time.Now()

	results, err := client.Query(
		ctx, &qdrant.QueryPoints{
			CollectionName: collectionName,
			Query:          qdrant.NewQuery(queryVector...),
			Params: &qdrant.SearchParams{
				HnswEf: qdrant.PtrOf(uint64(100)),
				Exact:  qdrant.PtrOf(false),
			},
			Limit:       qdrant.PtrOf(uint64(3)),
			WithPayload: qdrant.NewWithPayload(true),
		},
	)

	if err != nil {
		return fmt.Errorf("failed to search documents: %w", err)
	}

	latency := time.Since(start)

	for i, result := range results {
		fmt.Printf("Rank: %d\n", i+1)
		fmt.Println("ID:", result.GetId())
		fmt.Println("Score:", result.GetScore())
		fmt.Println("Payload:", result.GetPayload())
		fmt.Println("--------------------")
	}

	fmt.Println(strings.Repeat("=", 30))
	fmt.Println("\n Search latency:", latency)
	fmt.Println(strings.Repeat("=", 30))

	return nil
}
