// Package llm provides LLM provider adapters for different vendors.
package llm

import (
	"fmt"

	"github.com/aim/aim/internal/ai/service"
)

// Provider names.
const (
	ProviderOpenAI    = "openai"
	ProviderAnthropic = "anthropic"
)

// Config holds the configuration for an LLM provider.
type Config struct {
	Name     string
	Endpoint string
	APIKey   string
}

// ChatClient implements service.LLMClient for a specific provider.
type ChatClient struct {
	cfg Config
}

// NewChatClient creates a new ChatClient.
func NewChatClient(cfg Config) *ChatClient {
	return &ChatClient{cfg: cfg}
}

// Chat calls the LLM provider's chat endpoint.
// This is a stub implementation — in production, this would make actual HTTP calls.
func (c *ChatClient) Chat(model, systemPrompt string, history []service.ChatMessage, userMsg string) (string, int, error) {
	switch c.cfg.Name {
	case ProviderOpenAI:
		return c.chatOpenAI(model, systemPrompt, history, userMsg)
	case ProviderAnthropic:
		return c.chatAnthropic(model, systemPrompt, history, userMsg)
	default:
		return fmt.Sprintf("[%s/%s] 收到: %s", c.cfg.Name, model, userMsg), len([]rune(userMsg)) / 4, nil
	}
}

func (c *ChatClient) chatOpenAI(model, systemPrompt string, history []service.ChatMessage, userMsg string) (string, int, error) {
	// Stub: In production, call OpenAI API.
	reply := fmt.Sprintf("[OpenAI/%s] 回复: %s", model, userMsg)
	return reply, len([]rune(reply)) / 4, nil
}

func (c *ChatClient) chatAnthropic(model, systemPrompt string, history []service.ChatMessage, userMsg string) (string, int, error) {
	// Stub: In production, call Anthropic API.
	reply := fmt.Sprintf("[Anthropic/%s] 回复: %s", model, userMsg)
	return reply, len([]rune(reply)) / 4, nil
}

// Ensure ChatClient satisfies LLMClient.
var _ service.LLMClient = (*ChatClient)(nil)

// MultiProvider manages multiple LLM providers.
type MultiProvider struct {
	clients map[string]*ChatClient
}

// NewMultiProvider creates a MultiProvider from configs.
func NewMultiProvider(cfgs []Config) *MultiProvider {
	mp := &MultiProvider{clients: make(map[string]*ChatClient)}
	for _, cfg := range cfgs {
		mp.clients[cfg.Name] = NewChatClient(cfg)
	}
	return mp
}

// Get returns a provider by name.
func (mp *MultiProvider) Get(name string) (*ChatClient, bool) {
	c, ok := mp.clients[name]
	return c, ok
}
