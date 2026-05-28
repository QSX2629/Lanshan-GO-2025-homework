package repo

import (
	"github.com/aim/aim/internal/ai/model"
	"github.com/aim/aim/internal/pkg/database"
)

// AIRepo handles AI bot and session data access.
type AIRepo struct {
	db *database.DB
}

// NewAIRepo creates a new AIRepo.
func NewAIRepo(db *database.DB) *AIRepo {
	return &AIRepo{db: db}
}

// AutoMigrate creates the AI tables.
func (r *AIRepo) AutoMigrate() error {
	return r.db.AutoMigrate(
		&model.Bot{},
		&model.AISession{},
		&model.AIChatRecord{},
		&model.BillingRecord{},
	)
}

// --- Bot ---

// CreateBot inserts a new bot configuration.
func (r *AIRepo) CreateBot(b *model.Bot) error {
	return r.db.Create(b).Error
}

// FindBotByID looks up a bot.
func (r *AIRepo) FindBotByID(id uint) (*model.Bot, error) {
	var b model.Bot
	err := r.db.First(&b, id).Error
	if err != nil {
		return nil, err
	}
	return &b, nil
}

// ListBots returns all available bots.
func (r *AIRepo) ListBots() ([]model.Bot, error) {
	var bots []model.Bot
	err := r.db.Find(&bots).Error
	return bots, err
}

// UpdateBot updates bot configuration.
func (r *AIRepo) UpdateBot(b *model.Bot) error {
	return r.db.Save(b).Error
}

// --- AI Session ---

// CreateAISession starts a new AI conversation session.
func (r *AIRepo) CreateAISession(s *model.AISession) error {
	return r.db.Create(s).Error
}

// FindAISession looks up an AI session.
func (r *AIRepo) FindAISession(sessionID uint) (*model.AISession, error) {
	var s model.AISession
	err := r.db.First(&s, sessionID).Error
	if err != nil {
		return nil, err
	}
	return &s, nil
}

// UpdateAISession updates session context.
func (r *AIRepo) UpdateAISession(s *model.AISession) error {
	return r.db.Save(s).Error
}

// --- Chat Record ---

// AddChatRecord appends a turn to the AI conversation history.
func (r *AIRepo) AddChatRecord(record *model.AIChatRecord) error {
	return r.db.Create(record).Error
}

// GetChatHistory returns recent conversation history for a session.
func (r *AIRepo) GetChatHistory(sessionID uint, limit int) ([]model.AIChatRecord, error) {
	var records []model.AIChatRecord
	err := r.db.Where("session_id = ?", sessionID).
		Order("id ASC").Limit(limit).Find(&records).Error
	return records, err
}

// --- Billing ---

// CreateBillingRecord records a billing event.
func (r *AIRepo) CreateBillingRecord(b *model.BillingRecord) error {
	return r.db.Create(b).Error
}

// GetBillingRecords returns billing history for a user.
func (r *AIRepo) GetBillingRecords(userID uint, limit int) ([]model.BillingRecord, error) {
	var records []model.BillingRecord
	query := r.db.Where("user_id = ?", userID).
		Order("created_at DESC")
	if limit > 0 {
		query = query.Limit(limit)
	}
	err := query.Find(&records).Error
	return records, err
}
