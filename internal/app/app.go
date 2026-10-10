package app

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/kharljhon14/rag/internal/chat"
	"github.com/kharljhon14/rag/internal/config"
	"github.com/kharljhon14/rag/internal/ingest"
	"github.com/kharljhon14/rag/internal/llm"
	"github.com/kharljhon14/rag/internal/rag"
	"github.com/kharljhon14/rag/internal/vector"
	"github.com/kharljhon14/rag/internal/vector/pgvector"
	"github.com/kharljhon14/rag/internal/web"
)

func Run(parentCtx context.Context, cfg config.Config) error {
	logger := log.New(os.Stderr, "[rag] ", log.LstdFlags)

	ctx, cancel := context.WithCancel(parentCtx)
	defer cancel()

	client := llm.New(cfg)

	embedder := llm.NewEmbedder(cfg)

	store, err := openStore(ctx, cfg)
	if err != nil {
		logger.Printf("vector store disabled: %v", err)
	}
	if store != nil {
		defer store.Close()
		logger.Printf("vector store ready")
	}

	var wg sync.WaitGroup
	if store != nil {
		wg.Go(func() {
			opts := ingest.Options{
				SourceDir:    cfg.IngestDir,
				ProcessedDir: cfg.ProcessDir,
			}

			if err := ingest.Watch(ctx, opts, embedder, store, logger); err != nil && ctx.Err() == nil {
				logger.Printf("watcher stopped: %v", err)
			}
		})
		logger.Printf("watching %s for new documents", cfg.IngestDir)
	}

	// Get Retriver and Rewriter
	var retriever *rag.Retriver
	if store != nil {
		retriever = rag.New(embedder, store, rag.Options{
			TopK:     5,
			Rewriter: rag.NewRewriter(client),
		})
	}

	if cfg.HTTPAddr != "" {
		srv, err := web.New(client, embedder, retriever, web.Options{
			Addr:             cfg.HTTPAddr,
			SystemPromptFile: cfg.SystemPromptFile,
			Store:            store,
			ProcessedDir:     cfg.ProcessDir,
			ImagesDir:        cfg.ImageDir,
		})
		if err != nil {
			logger.Printf("web server disabled: %v", err)
		} else {
			wg.Go(func() {
				if err := srv.Run(ctx, cfg.HTTPAddr); err != nil && ctx.Err() == nil {
					logger.Printf("web server stopped: %v", err)
				}
			})
			logger.Printf("web chat at http://localhost%s/chat", cfg.HTTPAddr)
		}
	}

	replErr := chat.RunREPL(ctx, client, retriever, chat.Options{
		SystemPromptFile: cfg.SystemPromptFile,
	})
	cancel()
	wg.Wait()

	return replErr
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
