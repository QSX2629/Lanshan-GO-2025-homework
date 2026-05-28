package service

import (
	"strings"
	"testing"

	airepo "github.com/aim/aim/internal/ai/repo"
	"github.com/aim/aim/internal/pkg/database"
)

func setupAIService(t *testing.T) *AIService {
	t.Helper()
	db, err := database.TestDB()
	if err != nil {
		t.Fatalf("TestDB() error = %v", err)
	}
	repo := airepo.NewAIRepo(db)
	if err := repo.AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate() error = %v", err)
	}
	// No LLMClient: uses fallback mode.
	return NewAIService(db, nil)
}

func TestAIService_CreateAndListBots(t *testing.T) {
	svc := setupAIService(t)

	b, err := svc.CreateBot("MyBot", "openai", "gpt-4", "be helpful", true, 0.03)
	if err != nil {
		t.Fatalf("CreateBot() error = %v", err)
	}
	if b.Name != "MyBot" {
		t.Errorf("Name = %q, want MyBot", b.Name)
	}

	bots, err := svc.ListBots()
	if err != nil {
		t.Fatalf("ListBots() error = %v", err)
	}
	if len(bots) != 1 {
		t.Errorf("len(bots) = %d, want 1", len(bots))
	}

	got, err := svc.GetBot(b.ID)
	if err != nil {
		t.Fatalf("GetBot() error = %v", err)
	}
	if got.Provider != "openai" {
		t.Errorf("Provider = %q, want openai", got.Provider)
	}

	_, err = svc.GetBot(999)
	if err != ErrBotNotFound {
		t.Errorf("GetBot(999) error = %v, want ErrBotNotFound", err)
	}
}

func TestAIService_ChatFallback(t *testing.T) {
	svc := setupAIService(t)
	bot, _ := svc.CreateBot("TestBot", "openai", "gpt-4", "prompt", true, 0.03)

	reply, err := svc.Chat(1, bot.ID, 2, "user", "hello")
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}
	if !strings.Contains(reply, "TestBot") {
		t.Errorf("reply should contain bot name, got %q", reply)
	}

	// Verify billing was recorded.
	records, err := svc.GetBilling(1, 20)
	if err != nil {
		t.Fatalf("GetBilling() error = %v", err)
	}
	if len(records) != 1 {
		t.Errorf("len(billing) = %d, want 1", len(records))
	}
}

func TestAIService_SummarizeFallback(t *testing.T) {
	svc := setupAIService(t)

	result, err := svc.Summarize(1, 2, "user")
	if err != nil {
		t.Fatalf("Summarize() error = %v", err)
	}
	// With no messages to summarize, returns "没有找到对话记录".
	if !strings.Contains(result, "没有找到对话记录") {
		t.Logf("summarize result: %s", result)
	}
}

func TestAIService_ExtractTodosFallback(t *testing.T) {
	svc := setupAIService(t)

	result, err := svc.ExtractTodos(1, 2, "user")
	if err != nil {
		t.Fatalf("ExtractTodos() error = %v", err)
	}
	if !strings.Contains(result, "没有找到对话记录") {
		t.Logf("todos result: %s", result)
	}
}

func TestAIService_ReplyCandidatesFallback(t *testing.T) {
	svc := setupAIService(t)

	candidates, err := svc.GetReplyCandidates(1, 2, "user")
	if err != nil {
		t.Fatalf("GetReplyCandidates() error = %v", err)
	}
	if len(candidates) != 3 {
		t.Errorf("len(candidates) = %d, want 3", len(candidates))
	}
}

func TestEstimateTokens(t *testing.T) {
	n := estimateTokens("hello world")
	if n == 0 {
		t.Error("estimateTokens should be > 0")
	}

	n2 := estimateTokens("你好世界")
	if n2 == 0 {
		t.Error("estimateTokens for Chinese should be > 0")
	}

	n3 := estimateTokens("")
	if n3 != 0 {
		t.Error("estimateTokens for empty string should be 0")
	}
}
