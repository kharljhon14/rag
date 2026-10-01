package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	BaseURL          string
	APIKey           string
	Model            string
	SystemPromptFile string
}

func Load() Config {
	_ = godotenv.Load()

	cfg := Config{
		BaseURL:          os.Getenv("DEEPSEEK_BASE_URL"),
		APIKey:           os.Getenv("DEEKSEEK_API_KEY"),
		Model:            os.Getenv("DEEPSEEK_MODEL"),
		SystemPromptFile: os.Getenv("SYSTEM_PROMPT_FILE"),
	}

	if cfg.BaseURL == "" {
		cfg.BaseURL = "https://api.deepseek.com/"
	}

	if cfg.Model == "" {
		cfg.Model = "deepseek-flash"
	}

	return cfg
}
