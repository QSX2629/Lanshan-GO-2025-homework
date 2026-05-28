// Package billing provides usage and cost tracking for AI services.
package billing

import (
	"time"

	aimodel "github.com/aim/aim/internal/ai/model"
	airepo "github.com/aim/aim/internal/ai/repo"
	"github.com/aim/aim/internal/pkg/database"
)

// Manager handles billing calculations and record keeping.
type Manager struct {
	repo *airepo.AIRepo
}

// NewManager creates a new billing Manager.
func NewManager(db *database.DB) *Manager {
	return &Manager{repo: airepo.NewAIRepo(db)}
}

// RecordUsage logs a billing event for token usage.
func (m *Manager) RecordUsage(userID, botID uint, tokensUsed int, pricePer1K float64) error {
	cost := float64(tokensUsed) / 1000.0 * pricePer1K
	return m.repo.CreateBillingRecord(&aimodel.BillingRecord{
		UserID:     userID,
		BotID:      botID,
		TokensUsed: tokensUsed,
		Cost:       cost,
		CreatedAt:  time.Now(),
	})
}

// GetHistory returns billing history for a user.
func (m *Manager) GetHistory(userID uint, limit int) ([]aimodel.BillingRecord, error) {
	return m.repo.GetBillingRecords(userID, limit)
}

// GetTotalCost calculates the total cost for a user.
func (m *Manager) GetTotalCost(userID uint) (float64, error) {
	records, err := m.repo.GetBillingRecords(userID, 0)
	if err != nil {
		return 0, err
	}
	var total float64
	for _, r := range records {
		total += r.Cost
	}
	return total, nil
}

// GetTotalTokens returns the total tokens used by a user.
func (m *Manager) GetTotalTokens(userID uint) (int, error) {
	records, err := m.repo.GetBillingRecords(userID, 0)
	if err != nil {
		return 0, err
	}
	var total int
	for _, r := range records {
		total += r.TokensUsed
	}
	return total, nil
}

// CalculateCost computes the cost for a given token count.
func CalculateCost(tokensUsed int, pricePer1K float64) float64 {
	return float64(tokensUsed) / 1000.0 * pricePer1K
}
