package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/deblinux123/Qdrant-Database/day-12/services"
)

func main() {
	ctx, cancel := context.WithTimeout(
		context.Background(),
		5*time.Minute,
	)

	defer cancel()

	ollama := services.NewOllamaServiec()

	qdrant, err := services.NewQdrantService("papers")

	if err != nil {
		log.Fatal("Qdrant connection failed: %w", err)
	}

	defer qdrant.Close()

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Prompt > ")

	question, err := reader.ReadString('\n')

	if err != nil {
		log.Fatal(err)
	}

	question = strings.TrimSpace(question)

	if question == "" {
		log.Fatal("Prompt can be empty, please enter your question.")
	}

	vector, err := ollama.Embed(ctx, question)

	if err != nil {
		log.Fatal("Embedding failed:", err)
	}

	docs, err := qdrant.Search(ctx, vector, 5)

	if err != nil {
		log.Fatal("Qdrant search failed:", err)
	}

	var contextBuilder strings.Builder

	for i, doc := range docs {
		fmt.Printf("\nDocument %d | Score: %.4f\n", i+1, doc.Score)
		fmt.Printf("\n--- Document %d | Score: %.4f ---\n", i+1, doc.Score)
		fmt.Println(doc.Text)
		contextBuilder.WriteString(doc.Text)
		contextBuilder.WriteString("\n\n")
	}

	answer, err := ollama.GenerateAnswer(ctx, question, contextBuilder.String())

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("AI > ")
	fmt.Println(answer)
}
