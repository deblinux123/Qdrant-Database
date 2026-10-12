package ragservice

import (
	"context"
	"fmt"
	"strconv"

	"github.com/deblinux123/Qdrant-Database/day-13/config"
	"github.com/qdrant/go-client/qdrant"
)

type QdrantService struct {
	Client     *qdrant.Client
	Collection string
}

func NewClient(collection string) (*QdrantService, error) {
	cfg := config.Load()

	port, _ := strconv.Atoi(cfg.QdrantPort)

	client, err := qdrant.NewClient(&qdrant.Config{
		Host: cfg.QdrantHost,
		Port: port,
	})

	if err != nil {
		return nil, err
	}

	return &QdrantService{
		Client:     client,
		Collection: collection,
	}, nil
}

func (s QdrantService) CreateCollection(ctx context.Context) error {
	exists, err := s.Client.CollectionExists(ctx, s.Collection)

	if err != nil {
		return fmt.Errorf("failed to check collection:%w", err)
	}

	if exists {
		fmt.Println("Collection already exists.")
		return nil
	}

	err = s.Client.CreateCollection(ctx, &qdrant.CreateCollection{
		CollectionName: s.Collection,
		VectorsConfig: qdrant.NewVectorsConfig(
			&qdrant.VectorParams{
				Size:     768,
				Distance: qdrant.Distance_Cosine,
			},
		),
	})

	if err != nil {
		return fmt.Errorf("failed to create collection")
	}

	fmt.Printf("Collection %q created\n", s.Collection)

	return nil
}

func (s QdrantService) Search(ctx context.Context, vector []float32, limit uint64) ([]*qdrant.ScoredPoint, error) {
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
		return nil, fmt.Errorf(
			"failed to search points: %w",
			err,
		)
	}

	return points, nil
}

func (s QdrantService) Close() error {
	return s.Client.Close()
}
