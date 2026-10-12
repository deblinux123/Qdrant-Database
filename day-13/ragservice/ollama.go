package ragservice

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"

	"github.com/deblinux123/Qdrant-Database/day-13/config"
)

type OllamaService struct {
	URL   string
	Model string
}

type EmbedRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type EmbedRespnse struct {
	Embeddings [][]float32 `json:"embeddings"`
}

func NewOllamaService() *OllamaService {
	cfg := config.Load()

	return &OllamaService{
		URL:   cfg.OllamaURL,
		Model: cfg.OllamaModel,
	}
}

func (s *OllamaService) GetEmbedding(ctx context.Context, text string) ([]float32, error) {
	cfg := config.Load()

	request := EmbedRequest{
		Model: cfg.OllamaModel,
		Input: text,
	}

	bocy, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		cfg.OllamaURL+"/api/embed",
		bytes.NewReader(bocy),
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		message, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf(
			"ollama embed error: %s:%s",
			resp.Status,
			string(message),
		)
	}

	var result EmbedRespnse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Embeddings) == 0 {
		return nil, fmt.Errorf("No embedding returned")
	}

	return result.Embeddings[0], nil
}
