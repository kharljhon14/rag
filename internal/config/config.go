package config

import (
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type Config struct {
	BaseURL          string
	APIKey           string
	Model            string
	SystemPromptFile string
	DatabaseURL      string
	EmbeddingDim     int
}

func Load() Config {
	_ = godotenv.Load()

	cfg := Config{
		BaseURL:          os.Getenv("DEEPSEEK_BASE_URL"),
		APIKey:           os.Getenv("DEEPSEEK_API_KEY"),
		Model:            os.Getenv("DEEPSEEK_MODEL"),
		SystemPromptFile: os.Getenv("SYSTEM_PROMPT_FILE"),
		DatabaseURL:      os.Getenv("DATABASE_URL"),
		EmbeddingDim:     atoiOr(os.Getenv("EMBEDDING_DIM"), 0),
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.deepseek.com/"
	}

	if cfg.Model == "" {
		cfg.Model = "deepseek-flash"
	}

	if cfg.EmbeddingDim == 0 {
		cfg.EmbeddingDim = 768
	}

	return cfg
}

func atoiOr(s string, fallback int) int {
	if s == "" {
		return fallback
	}

	n, err := strconv.Atoi(s)
	if err != nil {
		return fallback
	}

	return n
}
