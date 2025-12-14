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
				Usage:   "Embedding function: 'default' (ONNX), 'gemini' (requires GEMINI_API_KEY), or 'hash' (testing)",
			},
		},
		Commands: []*cli.Command{
			{
				Name:   "add",
				Usage:  "Add documents to ChromaDB",
				Action: addCommand,
			},
			{
				Name:   "query",
				Usage:  "Query documents in ChromaDB",
				Action: queryCommand,
			},
		},
	}

	if err := app.Run(context.Background(), os.Args); err != nil {
		slog.Error("command failed", "error", err)
		os.Exit(1)
	}
}
