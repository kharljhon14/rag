package llm

import (
	"context"
	"errors"
	"fmt"
	"io"
	"strings"

	deepseek "github.com/cohesion-org/deepseek-go"
	"github.com/kharljhon14/rag/internal/config"
)

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type Client struct {
	cfg config.Config
	sdk *deepseek.Client
}

func New(cfg config.Config) (*Client, error) {
	opts := []deepseek.Option{}

	if cfg.BaseURL != "" {
		opts = append(opts, deepseek.WithBaseURL(cfg.BaseURL))
	}

	sdk, err := deepseek.NewClientWithOptions(cfg.APIKey, opts...)
	if err != nil {
		return nil, fmt.Errorf("deepseek new client options: %w", err)
	}

	return &Client{cfg: cfg, sdk: sdk}, nil
}

func (c *Client) ChatStream(ctx context.Context, messages []Message, onDelta func(string)) (Message, error) {
	stream, err := c.sdk.CreateChatCompletionStream(ctx, &deepseek.StreamChatCompletionRequest{
		Model:    c.cfg.Model,
		Messages: toSDKMessages(messages),
		Stream:   true,
	})
	if err != nil {
		return Message{}, fmt.Errorf("create chat completion options: %w", err)
	}
	defer stream.Close()

	var content strings.Builder
	role := "assistant"

	for {
		response, err := stream.Recv()
		if errors.Is(err, io.EOF) {
			fmt.Println("\n stream finished")
			break
		}

		if err != nil {
			return Message{}, fmt.Errorf("\nstream error: %v\n", err)
		}

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

	return Message{Role: role, Content: content.String()}, nil
}

func toSDKMessages(messages []Message) []deepseek.ChatCompletionMessage {
	out := make([]deepseek.ChatCompletionMessage, 0, len(messages))

	for _, m := range messages {
		switch m.Role {
		case "system":
			out = append(out, deepseek.ChatCompletionMessage{Role: deepseek.ChatMessageRoleSystem, Content: m.Content})
		case "assistant":
			out = append(out, deepseek.ChatCompletionMessage{Role: deepseek.ChatMessageRoleAssistant, Content: m.Content})
		default:
			out = append(out, deepseek.ChatCompletionMessage{Role: deepseek.ChatMessageRoleUser, Content: m.Content})
		}
	}

	return out
}
