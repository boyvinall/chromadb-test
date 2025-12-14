package main

import (
	"errors"
	"fmt"
	"log/slog"
	"os"

	"github.com/amikos-tech/chroma-go/pkg/embeddings"
	defaultef "github.com/amikos-tech/chroma-go/pkg/embeddings/default_ef"
	"github.com/amikos-tech/chroma-go/pkg/embeddings/gemini"
	"github.com/amikos-tech/chroma-go/pkg/embeddings/openai"
)

type embeddingConfig struct {
	embeddingType string
	openaiAPIKey  string
	openaiBaseURL string
	openaiModel   string
}

func createEmbeddingFunction(cfg embeddingConfig) (embeddings.EmbeddingFunction, error) {
	switch cfg.embeddingType {
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

	case "openai":
		if cfg.openaiAPIKey == "" {
			return nil, errors.New("--openai-api-key or OPENAI_API_KEY not set")
		}
		if cfg.openaiBaseURL == "" {
			return nil, errors.New("--openai-base-url or OPENAI_BASE_URL not set")
		}
		slog.Info("using OpenAI embedding", "base_url", cfg.openaiBaseURL, "model", cfg.openaiModel)
		ef, err := openai.NewOpenAIEmbeddingFunction(cfg.openaiAPIKey,
			openai.WithBaseURL(cfg.openaiBaseURL),
			openai.WithModel(openai.EmbeddingModel(cfg.openaiModel)),
		)
		if err != nil {
			return nil, fmt.Errorf("create OpenAI embedding: %w", err)
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
