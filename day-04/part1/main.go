package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

const qdrantURL = "http://localhost:6333"

type Point struct {
	ID      int               `json:"id"`
	Vector  []float32         `json:"vector"`
	Payload map[string]string `json:"payload"`
}

func createCollection() error {
	url := qdrantURL + "/collections/day04_points"

	body := map[string]any{
		"vectors": map[string]any{
			"size":     3,
			"distance": "Cosine",
		},
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	req, err := http.NewRequest(
		http.MethodPut,
		url,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	fmt.Println("Collection Created:", resp.Status)

	return nil
}

func insertPoints() error {
	points := []Point{
		{
			ID:     1,
			Vector: []float32{0.10, 0.20, 0.30},
			Payload: map[string]string{
				"title":    "Golang Basics",
				"category": "programming",
			},
		},
		{
			ID:     2,
			Vector: []float32{0.11, 0.21, 0.31},
			Payload: map[string]string{
				"title":    "Go Variables",
				"category": "programming",
			},
		},
		{
			ID:     3,
			Vector: []float32{0.12, 0.22, 0.32},
			Payload: map[string]string{
				"title":    "Go Functions",
				"category": "programming",
			},
		},
		{
			ID:     4,
			Vector: []float32{0.13, 0.23, 0.33},
			Payload: map[string]string{
				"title":    "Go Structs",
				"category": "programming",
			},
		},
		{
			ID:     5,
			Vector: []float32{0.14, 0.24, 0.34},
			Payload: map[string]string{
				"title":    "Qdrant Introduction",
				"category": "database",
			},
		},
		{
			ID:     6,
			Vector: []float32{0.15, 0.25, 0.35},
			Payload: map[string]string{
				"title":    "Qdrant Collections",
				"category": "database",
			},
		},
		{
			ID:     7,
			Vector: []float32{0.16, 0.26, 0.36},
			Payload: map[string]string{
				"title":    "Qdrant Points",
				"category": "database",
			},
		},
		{
			ID:     8,
			Vector: []float32{0.17, 0.27, 0.37},
			Payload: map[string]string{
				"title":    "Vector Search",
				"category": "ai",
			},
		},
		{
			ID:     9,
			Vector: []float32{0.18, 0.28, 0.38},
			Payload: map[string]string{
				"title":    "Embeddings",
				"category": "ai",
			},
		},
		{
			ID:     10,
			Vector: []float32{0.19, 0.29, 0.39},
			Payload: map[string]string{
				"title":    "RAG with Qdrant",
				"category": "ai",
			},
		},
	}

	body := map[string]any{
		"points": points,
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		return err
	}

	url := qdrantURL + "/collections/day04_points/points?wait=true"

	req, err := http.NewRequest(
		http.MethodPut,
		url,
		bytes.NewBuffer(jsonData),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}

	resp, err := client.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	fmt.Println("Insert Points:", resp.Status)

	return nil
}

func main() {
	err := createCollection()
	if err != nil {
		panic(err)
	}

	err = insertPoints()
	if err != nil {
		panic(err)
	}
}
