package adapter

import (
	"context"
	"errors"

	"github.com/sashabaranov/go-openai"
)

type DoubaoClient struct {
	client *openai.Client
	model  string
}

func NewDoubaoClient(apiKey string) *DoubaoClient {
	config := openai.DefaultConfig(apiKey)
	config.BaseURL = "https://ark.cn-beijing.volces.com/api/v3"

	return &DoubaoClient{
		client: openai.NewClientWithConfig(config),
		model:  "doubao-seed-2-0-pro-260215",
	}
}

func (c *DoubaoClient) Chat(prompt string) (string, error) {
	resp, err := c.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: c.model,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    "user",
					Content: prompt,
				},
			},
		},
	)

	if err != nil {
		return "", err
	}

	if len(resp.Choices) == 0 {
		return "", errors.New("empty response from doubao")
	}

	return resp.Choices[0].Message.Content, nil
}
