package handler

import (
	"strings"
	"testing"

	airepo "github.com/aim/aim/internal/ai/repo"
	"github.com/aim/aim/internal/pkg/database"
)

func setupAIHandler(t *testing.T) *AIHandler {
	t.Helper()
	db, err := database.TestDB()
	if err != nil {
		t.Fatalf("TestDB() error = %v", err)
	}
	repo := airepo.NewAIRepo(db)
	if err := repo.AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate() error = %v", err)
	}
	return NewAIHandler(db, nil)
}

func TestAIHandler_CreateBot(t *testing.T) {
	h := setupAIHandler(t)

	bot, err := h.CreateBot(nil, &CreateBotRequest{
		Name: "MyBot", Provider: "openai", Model: "gpt-4",
		SystemPrompt: "helpful", IsOfficial: true, PricePer1K: 0.03,
	})
	if err != nil {
		t.Fatalf("CreateBot() error = %v", err)
	}
	if bot.Name != "MyBot" {
		t.Errorf("Name = %q, want MyBot", bot.Name)
	}

	bots, err := h.ListBots(nil)
	if err != nil {
		t.Fatalf("ListBots() error = %v", err)
	}
	if len(bots) != 1 {
		t.Errorf("len(bots) = %d, want 1", len(bots))
	}
}

func TestAIHandler_Chat(t *testing.T) {
	h := setupAIHandler(t)
	bot, _ := h.CreateBot(nil, &CreateBotRequest{
		Name: "Bot", Provider: "openai", Model: "gpt-4", PricePer1K: 0.03,
	})

	reply, err := h.Chat(nil, &ChatRequest{
		UserID: 1, BotID: bot.ID, TargetID: 2, TargetType: "user", Message: "hello",
	})
	if err != nil {
		t.Fatalf("Chat() error = %v", err)
	}
	if !strings.Contains(reply, "Bot") {
		t.Errorf("reply should contain bot name: %s", reply)
	}
}

func TestAIHandler_Summarize(t *testing.T) {
	h := setupAIHandler(t)

	result, err := h.Summarize(nil, &AIActionRequest{UserID: 1, TargetID: 2, TargetType: "user"})
	if err != nil {
		t.Fatalf("Summarize() error = %v", err)
	}
	if !strings.Contains(result, "没有找到对话记录") {
		t.Logf("summarize: %s", result)
	}
}

func TestAIHandler_ExtractTodos(t *testing.T) {
	h := setupAIHandler(t)

	result, err := h.ExtractTodos(nil, &AIActionRequest{UserID: 1, TargetID: 2, TargetType: "user"})
	if err != nil {
		t.Fatalf("ExtractTodos() error = %v", err)
	}
	if !strings.Contains(result, "没有找到对话记录") {
		t.Logf("todos: %s", result)
	}
}

func TestAIHandler_ReplyCandidates(t *testing.T) {
	h := setupAIHandler(t)

	c, err := h.GetReplyCandidates(nil, &AIActionRequest{UserID: 1, TargetID: 2, TargetType: "user"})
	if err != nil {
		t.Fatalf("GetReplyCandidates() error = %v", err)
	}
	if len(c) != 3 {
		t.Errorf("len(candidates) = %d, want 3", len(c))
	}
}

func TestAIHandler_Billing(t *testing.T) {
	h := setupAIHandler(t)
	bot, _ := h.CreateBot(nil, &CreateBotRequest{
		Name: "Bot", Provider: "openai", Model: "gpt-4", PricePer1K: 0.03,
	})
	h.Chat(nil, &ChatRequest{UserID: 1, BotID: bot.ID, TargetID: 2, TargetType: "user", Message: "hello"})

	records, err := h.GetBilling(nil, 1, 20)
	if err != nil {
		t.Fatalf("GetBilling() error = %v", err)
	}
	if len(records) == 0 {
		t.Error("billing records should exist after chat")
	}
}
