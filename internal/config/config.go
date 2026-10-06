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
	EmbeddingBaseUrl string
	EmbeddingAPIKey  string
	EmbeddingModel   string
	IngestDir        string
	ProcessDir       string
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
		EmbeddingBaseUrl: os.Getenv("EMBEDDING_BASE_URL"),
		EmbeddingAPIKey:  os.Getenv("EMBEDDING_API_KEY"),
		EmbeddingModel:   os.Getenv("EMBEDDING_MODEL"),
		IngestDir:        os.Getenv("INGEST_DIR"),
		ProcessDir:       os.Getenv("PROCESS_DIR"),
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

	if cfg.EmbeddingBaseUrl == "" {
		cfg.EmbeddingBaseUrl = cfg.BaseURL
	}

	if cfg.EmbeddingModel == "" {
		cfg.EmbeddingModel = "nomic-embed-text"
	}

	if cfg.IngestDir == "" {
		cfg.IngestDir = "./documents"
	}

	if cfg.ProcessDir == "" {
		cfg.ProcessDir = "./documents/processed"
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
