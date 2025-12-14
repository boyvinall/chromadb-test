package main

import (
	"context"
	"log/slog"
	"os"

	"github.com/urfave/cli/v3"
)

func main() {
	slog.SetDefault(slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})))

	app := &cli.Command{
		Name:  "chromadb-test",
		Usage: "A CLI to interact with ChromaDB",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "server",
				Aliases: []string{"s"},
				Value:   "http://localhost:8001",
				Usage:   "ChromaDB server URL",
			},
			&cli.StringFlag{
				Name:    "embedding",
				Aliases: []string{"e"},
				Value:   "default",
				Usage:   "Embedding function: 'default' (ONNX), 'gemini', 'openai', or 'hash' (testing)",
			},
			&cli.StringFlag{
				Name:    "distance",
				Aliases: []string{"d"},
				Value:   "l2",
				Usage:   "Distance function: 'l2' (Euclidean), 'cosine', or 'ip' (inner product)",
			},
			&cli.StringFlag{
				Name:    "openai-api-key",
				Sources: cli.EnvVars("OPENAI_API_KEY", "ANTHROPIC_AUTH_TOKEN"),
				Usage:   "OpenAI API key (required for 'openai' embedding)",
			},
			&cli.StringFlag{
				Name:    "openai-base-url",
				Sources: cli.EnvVars("OPENAI_BASE_URL", "ANTHROPIC_BASE_URL"),
				Usage:   "OpenAI base URL (required for 'openai' embedding)",
			},
			&cli.StringFlag{
				Name:    "openai-model",
				Aliases: []string{"m"},
				Value:   "text-embedding-3-small",
				Usage:   "OpenAI embedding model",
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "add",
				Usage:  "Add documents to ChromaDB",
				Action: addCommand,
				Flags: []cli.Flag{
					&cli.StringFlag{
						Name:     "file",
						Aliases:  []string{"f"},
						Usage:    "YAML file containing documents to add",
						Required: true,
					},
				},
			},
			{
				Name:   "list",
				Usage:  "List documents in the collection",
				Action: listCommand,
				Flags: []cli.Flag{
					&cli.IntFlag{
						Name:    "limit",
						Aliases: []string{"n"},
						Value:   10,
						Usage:   "Maximum number of documents to list",
					},
					&cli.IntFlag{
						Name:  "offset",
						Value: 0,
						Usage: "Offset for pagination",
					},
				},
			},
			{
				Name:   "query",
				Usage:  "Query documents in ChromaDB",
				Action: queryCommand,
			},
			{
				Name:      "delete",
				Usage:     "Delete a document by ID",
				Action:    deleteCommand,
				ArgsUsage: "<document-id>",
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		slog.Error("command failed", "error", err)
		os.Exit(1)
	}
}
