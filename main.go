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
				Category: "core",
				Name:     "server",
				Aliases:  []string{"s"},
				Value:    "http://localhost:8001",
				Usage:    "ChromaDB server URL",
			},
			&cli.StringFlag{
				Category: "core",
				Name:     "embedding",
				Aliases:  []string{"e"},
				Value:    "default",
				Usage:    "Embedding function: 'default' (ONNX), 'gemini', 'openai', or 'hash'",
			},
			&cli.StringFlag{
				Category: "core",
				Name:     "distance",
				Aliases:  []string{"d"},
				Value:    "l2",
				Usage:    "Distance function: 'l2' (Euclidean), 'cosine', or 'ip' (inner product)",
			},
			&cli.StringFlag{
				Category: "openai",
				Name:     "openai-api-key",
				Sources:  cli.EnvVars("OPENAI_API_KEY", "ANTHROPIC_AUTH_TOKEN"),
				Usage:    "OpenAI API key (e.g. `sk-...`)",
			},
			&cli.StringFlag{
				Category: "openai",
				Name:     "openai-base-url",
				Sources:  cli.EnvVars("OPENAI_BASE_URL", "ANTHROPIC_BASE_URL"),
				Usage:    "OpenAI base URL",
				Value:    "https://api.openai.com/v1",
			},
			&cli.StringFlag{
				Category: "openai",
				Name:     "openai-model",
				Aliases:  []string{"m"},
				Value:    "text-embedding-3-small",
				Usage:    "OpenAI embedding model",
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
					&cli.BoolFlag{
						Name:    "upsert",
						Aliases: []string{"u"},
						Usage:   "Update existing documents instead of failing on duplicates",
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
				Flags: []cli.Flag{
					&cli.StringSliceFlag{
						Name:    "query",
						Aliases: []string{"q"},
						Usage:   "Query text to search for (can be repeated)",
					},
					&cli.IntFlag{
						Name:    "results",
						Aliases: []string{"n"},
						Value:   3,
						Usage:   "Number of results to return",
					},
				},
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
