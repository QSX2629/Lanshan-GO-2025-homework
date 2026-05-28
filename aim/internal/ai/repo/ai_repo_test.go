package repo

import (
	"testing"

	aimodel "github.com/aim/aim/internal/ai/model"
	"github.com/aim/aim/internal/pkg/database"
)

func setupAIRepo(t *testing.T) *AIRepo {
	t.Helper()
	db, err := database.TestDB()
	if err != nil {
		t.Fatalf("TestDB() error = %v", err)
	}
	repo := NewAIRepo(db)
	if err := repo.AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate() error = %v", err)
	}
	return repo
}

func TestAIRepo_Bot(t *testing.T) {
	repo := setupAIRepo(t)

	b := &aimodel.Bot{Name: "MyBot", Provider: "openai", Model: "gpt-4", IsOfficial: true}
	if err := repo.CreateBot(b); err != nil {
		t.Fatalf("CreateBot() error = %v", err)
	}

	found, err := repo.FindBotByID(b.ID)
	if err != nil {
		t.Fatalf("FindBotByID() error = %v", err)
	}
	if found.Name != "MyBot" {
		t.Errorf("Name = %q, want MyBot", found.Name)
	}

	bots, err := repo.ListBots()
	if err != nil {
		t.Fatalf("ListBots() error = %v", err)
	}
	if len(bots) != 1 {
		t.Errorf("len(bots) = %d, want 1", len(bots))
	}

	b.Name = "UpdatedBot"
	if err := repo.UpdateBot(b); err != nil {
		t.Fatalf("UpdateBot() error = %v", err)
	}
}

func TestAIRepo_Session(t *testing.T) {
	repo := setupAIRepo(t)

	s := &aimodel.AISession{UserID: 1, BotID: 1, TargetID: 2, TargetType: "user", ContextTokens: 4096}
	if err := repo.CreateAISession(s); err != nil {
		t.Fatalf("CreateAISession() error = %v", err)
	}

	found, err := repo.FindAISession(s.ID)
	if err != nil {
		t.Fatalf("FindAISession() error = %v", err)
	}
	if found.ContextTokens != 4096 {
		t.Errorf("ContextTokens = %d, want 4096", found.ContextTokens)
	}

	s.ContextTokens = 8192
	if err := repo.UpdateAISession(s); err != nil {
		t.Fatalf("UpdateAISession() error = %v", err)
	}
}

func TestAIRepo_ChatRecord(t *testing.T) {
	repo := setupAIRepo(t)

	r := &aimodel.AIChatRecord{SessionID: 1, Role: "user", Content: "hello", TokenCount: 10}
	if err := repo.AddChatRecord(r); err != nil {
		t.Fatalf("AddChatRecord() error = %v", err)
	}

	history, err := repo.GetChatHistory(1, 50)
	if err != nil {
		t.Fatalf("GetChatHistory() error = %v", err)
	}
	if len(history) != 1 {
		t.Errorf("len(history) = %d, want 1", len(history))
	}
}

func TestAIRepo_Billing(t *testing.T) {
	repo := setupAIRepo(t)

	br := &aimodel.BillingRecord{UserID: 1, BotID: 1, TokensUsed: 1000, Cost: 0.03}
	if err := repo.CreateBillingRecord(br); err != nil {
		t.Fatalf("CreateBillingRecord() error = %v", err)
	}

	records, err := repo.GetBillingRecords(1, 20)
	if err != nil {
		t.Fatalf("GetBillingRecords() error = %v", err)
	}
	if len(records) != 1 {
		t.Errorf("len(records) = %d, want 1", len(records))
	}
	if records[0].Cost != 0.03 {
		t.Errorf("Cost = %f, want 0.03", records[0].Cost)
	}
}
