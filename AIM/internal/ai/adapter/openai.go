package adapter

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
)

type OpenAIClient struct {
	APIKey string
	Model  string
}

func NewOpenAIClient(apiKey, model string) *OpenAIClient {
	return &OpenAIClient{APIKey: apiKey, Model: model}
}

func (c *OpenAIClient) Chat(prompt string) (string, error) {
	url := "https://api.openai.com/v1/chat/completions"
	reqBody := map[string]any{
		"model": c.Model,
		"messages": []map[string]string{
			{"role": "user", "content": prompt},
		},
	}
	jsonBody, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest("POST", url, bytes.NewBuffer(jsonBody))
	req.Header.Set("Authorization", "Bearer "+c.APIKey)
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}
	choices := result["choices"].([]any)
	if len(choices) == 0 {
		return "", fmt.Errorf("no response")
	}
	return choices[0].(map[string]any)["message"].(map[string]any)["content"].(string), nil
}
