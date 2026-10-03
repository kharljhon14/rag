package app

import (
	"context"
	"log"
	"os"

	"github.com/kharljhon14/rag/internal/chat"
	"github.com/kharljhon14/rag/internal/config"
	"github.com/kharljhon14/rag/internal/llm"
	"github.com/kharljhon14/rag/internal/vector"
	"github.com/kharljhon14/rag/internal/vector/pgvector"
)

func Run(ctx context.Context, cfg config.Config) error {
	logger := log.New(os.Stderr, "[rag] ", log.LstdFlags)

	client, err := llm.New(cfg)
	if err != nil {
		return err
	}

	store, err := openStore(ctx, cfg)
	if err != nil {
		logger.Printf("vector store disabled: %v", err)
	}
	if store != nil {
		defer store.Close()
		logger.Printf("vector store ready")
	}

	return chat.RunREPL(ctx, client, chat.Options{
		SystemPromptFile: cfg.SystemPromptFile,
	})
}

func openStore(ctx context.Context, cfg config.Config) (vector.Store, error) {
	if cfg.DatabaseURL == "" {
		return nil, nil
	}

	s, err := pgvector.New(ctx, pgvector.Options{
		DSN:          cfg.DatabaseURL,
		EmbeddingDim: cfg.EmbeddingDim,
	})
	if err != nil {
		return nil, err
	}

	return s, err
}
