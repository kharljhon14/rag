package web

import (
	"embed"
	"errors"
	"fmt"
	"html/template"
	"io/fs"
	"os"
	"strings"

	"github.com/kharljhon14/rag/internal/llm"
	"github.com/kharljhon14/rag/internal/rag"
	"github.com/kharljhon14/rag/internal/vector"
)

//go:embed templates/*.go.html
var templateFS embed.FS

type Options struct {
	Addr             string
	SystemPromptFile string
	Title            string
	Store            vector.Store
	ProcessedDir     string
	ImagesDir        string
}

type Server struct {
	client       *llm.Client
	embedded     *llm.Client
	retriever    *rag.Retriver
	store        vector.Store
	processedDir string
	imagesDir    string
	tpl          *template.Template
	system       string
	title        string
}

func New(client, embedder *llm.Client, retriver *rag.Retriver, opts Options) (*Server, error) {
	tpl, err := template.ParseFS(templateFS, "templates/*.go.html")
	if err != nil {
		return nil, fmt.Errorf("parse templates: %w", err)
	}

	title := opts.Title
	if title == "" {
		title = "RAG Chat"
	}

	return &Server{
		client:       client,
		embedded:     embedder,
		retriever:    retriver,
		store:        opts.Store,
		processedDir: opts.ProcessedDir,
		imagesDir:    opts.ImagesDir,
		tpl:          tpl,
		system:       readSystemPropmt(opts.SystemPromptFile),
		title:        title,
	}, nil
}

func readSystemPropmt(path string) string {
	if path == "" {
		return ""
	}

	data, err := os.ReadFile(path)
	if errors.Is(err, fs.ErrNotExist) {
		return ""
	}

	return strings.TrimSpace(string(data))
}
