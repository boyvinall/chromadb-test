package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"os"

	chroma "github.com/amikos-tech/chroma-go/pkg/api/v2"
	"github.com/urfave/cli/v3"
	"gopkg.in/yaml.v3"
)

type document struct {
	ID       string            `yaml:"id"`
	Text     string            `yaml:"text"`
	Metadata map[string]string `yaml:"metadata"`
}

func loadDocuments(filename string) ([]document, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, fmt.Errorf("open file: %w", err)
	}
	defer f.Close()

	var documents []document
	decoder := yaml.NewDecoder(f)
	for {
		var doc document
		if err := decoder.Decode(&doc); err != nil {
			if errors.Is(err, io.EOF) {
				break
			}
			return nil, fmt.Errorf("parse yaml: %w", err)
		}
		documents = append(documents, doc)
	}

	return documents, nil
}

func addCommand(ctx context.Context, cmd *cli.Command) error {
	filename := cmd.String("file")
	if filename == "" {
		return fmt.Errorf("--file flag is required")
	}

	documents, err := loadDocuments(filename)
	if err != nil {
		return fmt.Errorf("load documents: %w", err)
	}

	if len(documents) == 0 {
		return fmt.Errorf("no documents found in %s", filename)
	}

	session, err := newSession(ctx, cmd.String("server"), embeddingConfig{
		embeddingType: cmd.String("embedding"),
		openaiAPIKey:  cmd.String("openai-api-key"),
		openaiBaseURL: cmd.String("openai-base-url"),
		openaiModel:   cmd.String("openai-model"),
	}, cmd.String("distance"))
	if err != nil {
		return err
	}
	defer session.Close() //nolint:errcheck

	ids := make([]chroma.DocumentID, len(documents))
	texts := make([]string, len(documents))
	metadatas := make([]chroma.DocumentMetadata, len(documents))

	for i, doc := range documents {
		ids[i] = chroma.DocumentID(doc.ID)
		texts[i] = doc.Text
		meta := make(map[string]any)
		for k, v := range doc.Metadata {
			meta[k] = v
		}
		metadatas[i], _ = chroma.NewDocumentMetadataFromMap(meta)
	}

	slog.Info("adding documents", "count", len(documents), "file", filename)
	if err := session.collection.Add(ctx,
		chroma.WithIDs(ids...),
		chroma.WithTexts(texts...),
		chroma.WithMetadatas(metadatas...),
	); err != nil {
		return fmt.Errorf("add documents: %w", err)
	}

	count, err := session.collection.Count(ctx)
	if err != nil {
		return fmt.Errorf("count documents: %w", err)
	}
	slog.Info("documents added", "added", len(documents), "total", count)

	return nil
}

func listCommand(ctx context.Context, cmd *cli.Command) error {
	session, err := newSession(ctx, cmd.String("server"), embeddingConfig{
		embeddingType: cmd.String("embedding"),
		openaiAPIKey:  cmd.String("openai-api-key"),
		openaiBaseURL: cmd.String("openai-base-url"),
		openaiModel:   cmd.String("openai-model"),
	}, cmd.String("distance"))
	if err != nil {
		return err
	}
	defer session.Close() //nolint:errcheck

	limit := int(cmd.Int("limit"))
	offset := int(cmd.Int("offset"))

	results, err := session.collection.Get(ctx,
		chroma.WithLimitGet(limit),
		chroma.WithOffsetGet(offset),
		chroma.WithIncludeGet(chroma.IncludeDocuments, chroma.IncludeMetadatas),
	)
	if err != nil {
		return fmt.Errorf("get documents: %w", err)
	}

	ids := results.GetIDs()
	docs := results.GetDocuments()
	metadatas := results.GetMetadatas()

	if len(ids) == 0 {
		fmt.Println("No documents found.")
		return nil
	}

	fmt.Printf("=== Documents (showing %d, offset %d) ===\n", len(ids), offset)
	for i, id := range ids {
		fmt.Printf("\n[%d] ID: %s\n", i+1, id)
		if i < len(docs) && docs[i] != nil {
			text := docs[i].ContentString()
			if len(text) > 100 {
				text = text[:100] + "..."
			}
			fmt.Printf("    Text: %s\n", text)
		}
		if i < len(metadatas) && metadatas[i] != nil {
			if meta, ok := metadatas[i].(*chroma.DocumentMetadataImpl); ok {
				for _, key := range meta.Keys() {
					if val, ok := meta.GetString(key); ok {
						fmt.Printf("    %s: %v\n", key, val)
					}
				}
			}
		}
	}

	count, err := session.collection.Count(ctx)
	if err != nil {
		return fmt.Errorf("count documents: %w", err)
	}
	fmt.Printf("\nTotal documents in collection: %d\n", count)

	return nil
}

func queryCommand(ctx context.Context, cmd *cli.Command) error {
	session, err := newSession(ctx, cmd.String("server"), embeddingConfig{
		embeddingType: cmd.String("embedding"),
		openaiAPIKey:  cmd.String("openai-api-key"),
		openaiBaseURL: cmd.String("openai-base-url"),
		openaiModel:   cmd.String("openai-model"),
	}, cmd.String("distance"))
	if err != nil {
		return err
	}
	defer session.Close() //nolint:errcheck

	queryText := "neural networks and AI"
	slog.Info("querying collection", "query", queryText)

	results, err := session.collection.Query(ctx,
		chroma.WithQueryTexts(queryText),
		chroma.WithNResults(3),
	)
	if err != nil {
		return fmt.Errorf("query: %w", err)
	}

	idGroups := results.GetIDGroups()
	docGroups := results.GetDocumentsGroups()
	distanceGroups := results.GetDistancesGroups()

	if len(idGroups) == 0 || len(idGroups[0]) == 0 {
		slog.Info("no results found")
		return nil
	}

	fmt.Println("\n=== Query Results ===")
	for i, id := range idGroups[0] {
		doc := ""
		if i < len(docGroups[0]) {
			doc = docGroups[0][i].ContentString()
		}
		distance := float32(0)
		if i < len(distanceGroups[0]) {
			distance = float32(distanceGroups[0][i])
		}
		fmt.Printf("  [%d] ID: %s, Distance: %.4f\n      Document: %s\n", i+1, id, distance, doc)
	}

	return nil
}

func deleteCommand(ctx context.Context, cmd *cli.Command) error {
	args := cmd.Args()
	if args.Len() == 0 {
		return fmt.Errorf("document ID required")
	}
	docID := args.First()

	session, err := newSession(ctx, cmd.String("server"), embeddingConfig{
		embeddingType: cmd.String("embedding"),
		openaiAPIKey:  cmd.String("openai-api-key"),
		openaiBaseURL: cmd.String("openai-base-url"),
		openaiModel:   cmd.String("openai-model"),
	}, cmd.String("distance"))
	if err != nil {
		return err
	}
	defer session.Close() //nolint:errcheck

	slog.Info("deleting document", "id", docID)
	if err := session.collection.Delete(ctx,
		chroma.WithIDsDelete(chroma.DocumentID(docID)),
	); err != nil {
		return fmt.Errorf("delete document: %w", err)
	}

	slog.Info("document deleted", "id", docID)
	return nil
}
