package main

import (
	"context"
	"fmt"
	"log/slog"

	chroma "github.com/amikos-tech/chroma-go/pkg/api/v2"
	"github.com/urfave/cli/v3"
)

var sampleDocuments = []struct {
	ID       chroma.DocumentID
	Text     string
	Category string
	Type     string
}{
	{"doc1", "The quick brown fox jumps over the lazy dog", "animals", "classic"},
	{"doc2", "A fast orange cat leaps across the sleepy puppy", "animals", "variant"},
	{"doc3", "Machine learning is a subset of artificial intelligence", "technology", "ml"},
	{"doc4", "Deep learning uses neural networks with many layers", "technology", "dl"},
	{"doc5", "Go is a statically typed, compiled programming language", "technology", "programming"},
}

func addCommand(ctx context.Context, cmd *cli.Command) error {
	session, err := newSession(ctx, cmd.String("server"), cmd.String("embedding"))
	if err != nil {
		return err
	}
	defer session.Close() //nolint:errcheck

	ids := make([]chroma.DocumentID, len(sampleDocuments))
	texts := make([]string, len(sampleDocuments))
	metadatas := make([]chroma.DocumentMetadata, len(sampleDocuments))

	for i, doc := range sampleDocuments {
		ids[i] = doc.ID
		texts[i] = doc.Text
		metadatas[i], _ = chroma.NewDocumentMetadataFromMap(map[string]any{
			"category": doc.Category,
			"type":     doc.Type,
		})
	}

	slog.Info("adding documents", "count", len(sampleDocuments))
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
	slog.Info("documents added", "added", len(sampleDocuments), "total", count)

	return nil
}

func queryCommand(ctx context.Context, cmd *cli.Command) error {
	session, err := newSession(ctx, cmd.String("server"), cmd.String("embedding"))
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
