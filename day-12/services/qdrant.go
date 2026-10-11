package services

import (
	"context"
	"fmt"

	"github.com/qdrant/go-client/qdrant"
)

type qdrantService struct {
	Client     *qdrant.Client
	Collection string
}

type retriveDocument struct {
	Text  string
	Score float32
}

func NewQdrantService(collection string) (*qdrantService, error) {
	client, err := qdrant.NewClient(&qdrant.Config{
		Host: "localhost",
		Port: 6334,
	})

	if err != nil {
		return nil, err
	}

	return &qdrantService{
		Client:     client,
		Collection: collection,
	}, nil
}

func (s *qdrantService) Search(ctx context.Context, vector []float32, limit uint64) ([]retriveDocument, error) {
	points, err := s.Client.Query(
		ctx,
		&qdrant.QueryPoints{
			CollectionName: s.Collection,
			Query:          qdrant.NewQuery(vector...),
			Limit:          &limit,
			WithPayload:    qdrant.NewWithPayload(true),
		},
	)

	if err != nil {
		return nil, err
	}

	docs := make([]retriveDocument, 0, len(points))

	for _, point := range points {
		if point.Payload == nil {
			continue
		}

		value, ok := point.Payload["content"]

		if !ok {
			continue
		}

		textValue, ok := value.GetKind().(*qdrant.Value_StringValue)
		if !ok {
			continue
		}

		docs = append(docs, retriveDocument{
			Text:  textValue.StringValue,
			Score: point.Score,
		})
	}

	if len(docs) == 0 {
		return nil, fmt.Errorf(
			"no documents with a text payload found in collection %q",
			s.Collection,
		)
	}

	return docs, nil
}

func (s *qdrantService) Close() error {
	return s.Client.Close()
}
