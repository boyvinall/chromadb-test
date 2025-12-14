package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/amikos-tech/chroma-go/pkg/embeddings"
	defaultef "github.com/amikos-tech/chroma-go/pkg/embeddings/default_ef"
	"github.com/amikos-tech/chroma-go/pkg/embeddings/gemini"
)

func createEmbeddingFunction(embeddingType string) (embeddings.EmbeddingFunction, error) {
	switch embeddingType {
	case "hash":
		slog.Info("using hash embedding (testing only)")
		return embeddings.NewConsistentHashEmbeddingFunction(), nil

	case "gemini":
		apiKey := os.Getenv("GEMINI_API_KEY")
		if apiKey == "" {
			return nil, errors.New("GEMINI_API_KEY not set (see https://aistudio.google.com/app/api-keys)")
		}
		slog.Info("using Gemini embedding")
		ef, err := gemini.NewGeminiEmbeddingFunction(gemini.WithAPIKey(apiKey))
		if err != nil {
			return nil, fmt.Errorf("create Gemini embedding: %w", err)
		}
		return ef, nil

	default:
		slog.Info("using default ONNX embedding", "model", "all-MiniLM-L6-v2")
		ef, _, err := defaultef.NewDefaultEmbeddingFunction()
		if err != nil {
			return nil, fmt.Errorf("create default embedding: %w", err)
		}
		return ef, nil
	}
}

