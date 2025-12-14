package main

import (
	"context"
	"fmt"
	"log/slog"

	chroma "github.com/amikos-tech/chroma-go/pkg/api/v2"
	"github.com/amikos-tech/chroma-go/pkg/embeddings"
)

const collectionName = "example_collection"

type chromaSession struct {
	client     chroma.Client
	collection chroma.Collection
}

func (s *chromaSession) Close() error {
	if s.client != nil {
		return s.client.Close()
	}
	return nil
}

func parseDistanceMetric(distanceType string) (embeddings.DistanceMetric, error) {
	switch distanceType {
	case "l2":
		return embeddings.L2, nil
	case "cosine":
		return embeddings.COSINE, nil
	case "ip":
		return embeddings.IP, nil
	default:
		return "", fmt.Errorf("unknown distance function: %s (valid: l2, cosine, ip)", distanceType)
	}
}

func newSession(ctx context.Context, serverURL string, embeddingCfg embeddingConfig, distanceType string) (*chromaSession, error) {
	slog.Info("connecting to ChromaDB", "url", serverURL)

	client, err := chroma.NewHTTPClient(
		chroma.WithBaseURL(serverURL),
		chroma.WithDefaultDatabaseAndTenant(),
	)
	if err != nil {
		return nil, fmt.Errorf("create client: %w", err)
	}

	if err := client.Heartbeat(ctx); err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("heartbeat: %w", err)
	}
	slog.Info("connected to ChromaDB")

	embeddingFunc, err := createEmbeddingFunction(embeddingCfg)
	if err != nil {
		_ = client.Close()
		return nil, err
	}

	distanceMetric, err := parseDistanceMetric(distanceType)
	if err != nil {
		_ = client.Close()
		return nil, err
	}

	slog.Info("creating collection",
		"name", collectionName,
		"distance", distanceType,
		"embedding", embeddingCfg.embeddingType,
	)
	coll, err := client.GetOrCreateCollection(ctx, collectionName,
		chroma.WithEmbeddingFunctionCreate(embeddingFunc),
		chroma.WithCollectionMetadataCreate(chroma.NewMetadata(
			chroma.NewStringAttribute("description", "An example collection for testing"),
		)),
		chroma.WithHNSWSpaceCreate(distanceMetric),
	)
	if err != nil {
		_ = client.Close()
		return nil, fmt.Errorf("get/create collection: %w", err)
	}

	slog.Info("collection ready")
	return &chromaSession{
		client:     client,
		collection: coll,
	}, nil
}
