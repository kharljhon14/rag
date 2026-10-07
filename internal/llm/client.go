package llm

import (
	"context"
	"fmt"
	"strings"

	"github.com/kharljhon14/rag/internal/config"
	"github.com/openai/openai-go/v3"
	"github.com/openai/openai-go/v3/option"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Client struct {
	cfg config.Config
	sdk openai.Client
}

func New(cfg config.Config) *Client {
	return newClient(cfg, cfg.BaseURL, cfg.APIKey)
}

func NewEmbedder(cfg config.Config) *Client {
	return newClient(cfg, cfg.EmbeddingBaseUrl, cfg.EmbeddingAPIKey)
}

func newClient(cfg config.Config, baseUrl, apiKey string) *Client {
	opts := []option.RequestOption{option.WithAPIKey(cfg.APIKey)}

	if baseUrl != "" {
		opts = append(opts, option.WithBaseURL(baseUrl))
	}

	if apiKey != "" {
		opts = append(opts, option.WithAPIKey(apiKey))
	}

	sdk := openai.NewClient(opts...)

	return &Client{cfg: cfg, sdk: sdk}
}

func (c *Client) ChatStream(ctx context.Context, messages []Message, onDelta func(string)) (Message, error) {
	stream := c.sdk.Chat.Completions.NewStreaming(ctx, openai.ChatCompletionNewParams{
		Model:    c.cfg.Model,
		Messages: toSDKMessages(messages),
	})
	defer stream.Close()

	var content strings.Builder
	role := "assistant"

	for stream.Next() {
		response := stream.Current()

		for _, choice := range response.Choices {
			delta := choice.Delta
			if delta.Role != "" {
				role = delta.Role
			}

			if delta.Content != "" {
				content.WriteString(delta.Content)
				if onDelta != nil {
					onDelta(delta.Content)
				}
			}
		}
	}

	if err := stream.Err(); err != nil {
		return Message{}, fmt.Errorf("\nstream error: %v\n", err)
	}

	fmt.Println("\n stream finished")

	return Message{Role: role, Content: content.String()}, nil
}

func toSDKMessages(messages []Message) []openai.ChatCompletionMessageParamUnion {
	out := make([]openai.ChatCompletionMessageParamUnion, 0, len(messages))

	for _, m := range messages {
		switch m.Role {
		case "system":
			out = append(out, openai.SystemMessage(m.Content))
		case "assistant":
			out = append(out, openai.AssistantMessage(m.Content))
		default:
			out = append(out, openai.UserMessage(m.Content))
		}
	}

	return out
}
