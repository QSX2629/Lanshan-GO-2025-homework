package billing

import (
	"testing"

	airepo "github.com/aim/aim/internal/ai/repo"
	"github.com/aim/aim/internal/pkg/database"
)

func setupBillingManager(t *testing.T) *Manager {
	t.Helper()
	db, err := database.TestDB()
	if err != nil {
		t.Fatalf("TestDB() error = %v", err)
	}
	repo := airepo.NewAIRepo(db)
	if err := repo.AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate() error = %v", err)
	}
	return NewManager(db)
}

func TestManager_RecordAndHistory(t *testing.T) {
	m := setupBillingManager(t)

	if err := m.RecordUsage(1, 1, 1000, 0.03); err != nil {
		t.Fatalf("RecordUsage() error = %v", err)
	}

	history, err := m.GetHistory(1, 20)
	if err != nil {
		t.Fatalf("GetHistory() error = %v", err)
	}
	if len(history) != 1 {
		t.Errorf("len(history) = %d, want 1", len(history))
	}

	cost, err := m.GetTotalCost(1)
	if err != nil {
		t.Fatalf("GetTotalCost() error = %v", err)
	}
	if cost != 0.03 {
		t.Errorf("total cost = %f, want 0.03", cost)
	}

	totalTokens, err := m.GetTotalTokens(1)
	if err != nil {
		t.Fatalf("GetTotalTokens() error = %v", err)
	}
	if totalTokens != 1000 {
		t.Errorf("total tokens = %d, want 1000", totalTokens)
	}
}

func TestManager_MultipleRecords(t *testing.T) {
	m := setupBillingManager(t)
	m.RecordUsage(1, 1, 500, 0.02)
	m.RecordUsage(1, 1, 500, 0.02)

	cost, _ := m.GetTotalCost(1)
	if cost != 0.02 { // 500*0.02/1000 + 500*0.02/1000 = 0.01 + 0.01
		t.Errorf("total cost = %f, want 0.02", cost)
	}
}

func TestCalculateCost(t *testing.T) {
	c := CalculateCost(2000, 0.03)
	if c != 0.06 {
		t.Errorf("CalculateCost = %f, want 0.06", c)
	}

	c2 := CalculateCost(0, 0.03)
	if c2 != 0 {
		t.Errorf("CalculateCost(0) = %f, want 0", c2)
	}
}
