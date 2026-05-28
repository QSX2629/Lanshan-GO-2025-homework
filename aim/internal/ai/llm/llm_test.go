package llm

import (
	"strings"
	"testing"

	"github.com/aim/aim/internal/ai/service"
)

func TestChatClient_OpenAI(t *testing.T) {
	client := NewChatClient(Config{Name: ProviderOpenAI, Endpoint: "https://api.openai.com/v1"})

	reply, tokens, err := client.Chat("gpt-4", "be helpful", nil, "hello")
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}
	if !strings.Contains(reply, "OpenAI") {
		t.Errorf("reply should contain provider name: %s", reply)
	}
	if tokens == 0 {
		t.Error("tokens should be > 0")
	}
}

func TestChatClient_Anthropic(t *testing.T) {
	client := NewChatClient(Config{Name: ProviderAnthropic, Endpoint: "https://api.anthropic.com/v1"})

	reply, tokens, err := client.Chat("claude-3", "be helpful", nil, "hello")
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}
	if !strings.Contains(reply, "Anthropic") {
		t.Errorf("reply should contain provider name: %s", reply)
	}
	if tokens == 0 {
		t.Error("tokens should be > 0")
	}
}

func TestChatClient_Unknown(t *testing.T) {
	client := NewChatClient(Config{Name: "unknown"})

	reply, _, err := client.Chat("model", "", nil, "hello")
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}
	if reply == "" {
		t.Error("reply should not be empty")
	}
}

func TestChatClient_WithHistory(t *testing.T) {
	client := NewChatClient(Config{Name: ProviderOpenAI})
	history := []service.ChatMessage{
		{Role: "user", Content: "hi"},
		{Role: "assistant", Content: "hello"},
	}
	reply, _, err := client.Chat("gpt-4", "system", history, "how are you")
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}
	if reply == "" {
		t.Error("reply should not be empty")
	}
}

func TestMultiProvider(t *testing.T) {
	mp := NewMultiProvider([]Config{
		{Name: ProviderOpenAI, Endpoint: "https://api.openai.com/v1"},
		{Name: ProviderAnthropic, Endpoint: "https://api.anthropic.com/v1"},
	})

	c, ok := mp.Get(ProviderOpenAI)
	if !ok || c == nil {
		t.Error("should get OpenAI client")
	}

	_, ok = mp.Get("nonexistent")
	if ok {
		t.Error("should not find nonexistent provider")
	}
}

func TestChatClientImplementsInterface(t *testing.T) {
	// Compile-time check: ChatClient should implement LLMClient.
	var _ service.LLMClient = NewChatClient(Config{Name: "test"})
}
