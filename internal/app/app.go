package app

import (
	"context"

	"github.com/kharljhon14/rag/internal/chat"
	"github.com/kharljhon14/rag/internal/config"
	"github.com/kharljhon14/rag/internal/llm"
)

func Run(ctx context.Context, cfg config.Config) error {
	client, err := llm.New(cfg)
	if err != nil {
		return err
	}

	return chat.RunREPL(ctx, client, chat.Options{
		SystemPromptFile: cfg.SystemPromptFile,
	})
}
