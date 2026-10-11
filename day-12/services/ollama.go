package services

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type OllamaService struct {
	BaseURL        string
	EmbeddingModel string
	ChatModel      string
	Client         *http.Client
}

type embedRequest struct {
	Model string `json:"model"`
	Input string `json:"input"`
}

type embedResponse struct {
	Embeddings [][]float32 `json:"embeddings"`
}

type chatMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type chatRequest struct {
	Model    string        `json:"model"`
	Messages []chatMessage `json:"messages"`
	Stream   bool          `json:"stream"`
}

type chatResponse struct {
	Message chatMessage `json:"message"`
}

func NewOllamaServiec() *OllamaService {
	return &OllamaService{
		BaseURL:        "http://localhost:11434",
		EmbeddingModel: "embeddinggemma:300m",
		ChatModel:      "gemma3:4b",
		Client:         &http.Client{},
	}
}

func (s *OllamaService) Embed(ctx context.Context, text string) ([]float32, error) {

	request := embedRequest{
		Model: s.EmbeddingModel,
		Input: text,
	}

	body, err := json.Marshal(request)

	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.BaseURL+"/api/embed",
		bytes.NewReader(body),
	)

	if err != nil {
		return nil, err
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := s.Client.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		message, _ := io.ReadAll(resp.Body)

		return nil, fmt.Errorf(
			"ollama embed error: %s:%s", resp.Status, string(message),
		)
	}

	var result embedResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	if len(result.Embeddings) == 0 {
		return nil, fmt.Errorf(
			"ollama returned no embeddings",
		)
	}

	return result.Embeddings[0], nil
}

func (s *OllamaService) GenerateAnswer(ctx context.Context, question string, contextText string) (string, error) {
	prompt := fmt.Sprintf(
		`به زبان فارسی و فقط بر اساس متن مرجع زیر پاسخ بده.
اگر اطلاعات کافی در متن وجود ندارد، صریحاً بگو پاسخ در اسناد بازیابی‌شده پیدا نشد.

متن مرجع:
%s

سؤال کاربر:
%s`,
		contextText,
		question,
	)

	body, err := json.Marshal(chatRequest{
		Model: s.ChatModel,
		Messages: []chatMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
		Stream: false,
	})

	if err != nil {
		return "", err
	}

	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		s.BaseURL+"/api/chat",
		bytes.NewReader(body),
	)

	if err != nil {
		return "", err
	}

	req.Header.Set(
		"Content-Type",
		"application/json",
	)

	resp, err := s.Client.Do(req)

	if err != nil {
		return "", err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		message, _ := io.ReadAll(resp.Body)

		return "", fmt.Errorf(
			"ollama chat error: %s:%s", resp.Status, string(message),
		)
	}

	var result chatResponse

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Message.Content, nil
}
