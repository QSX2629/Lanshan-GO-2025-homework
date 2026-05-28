package model

import "time"

// BotRole constants.
const (
	BotRoleUser      = "user"
	BotRoleAssistant = "assistant"
)

// Bot represents an AI bot configuration.
type Bot struct {
	ID           uint      `gorm:"primaryKey" json:"id"`
	Name         string    `gorm:"size:128;not null" json:"name"`
	Provider     string    `gorm:"size:64;not null" json:"provider"`
	Model        string    `gorm:"size:64;not null" json:"model"`
	SystemPrompt string    `gorm:"type:text" json:"system_prompt"`
	APIKey       string    `gorm:"size:512" json:"-"` // user-provided, encrypted
	IsOfficial   bool      `gorm:"not null;default:false" json:"is_official"`
	PricePer1K   float64   `json:"price_per_1k"` // cost per 1000 tokens
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// AISession holds the context window for an AI conversation.
type AISession struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	UserID        uint      `gorm:"index;not null" json:"user_id"`
	BotID         uint      `gorm:"not null" json:"bot_id"`
	TargetID      uint      `gorm:"not null" json:"target_id"` // user or group ID
	TargetType    string    `gorm:"size:16;not null" json:"target_type"`
	ContextTokens int       `gorm:"not null;default:0" json:"context_tokens"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

// AIChatRecord stores a single turn of an AI conversation.
type AIChatRecord struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	SessionID  uint      `gorm:"index;not null" json:"session_id"`
	Role       string    `gorm:"size:16;not null" json:"role"` // "user" or "assistant"
	Content    string    `gorm:"type:text" json:"content"`
	TokenCount int       `gorm:"not null;default:0" json:"token_count"`
	CreatedAt  time.Time `json:"created_at"`
}

// BillingRecord tracks token usage and cost per request.
type BillingRecord struct {
	ID         uint      `gorm:"primaryKey" json:"id"`
	UserID     uint      `gorm:"index;not null" json:"user_id"`
	BotID      uint      `gorm:"not null" json:"bot_id"`
	TokensUsed int       `gorm:"not null" json:"tokens_used"`
	Cost       float64   `gorm:"not null" json:"cost"`
	CreatedAt  time.Time `json:"created_at"`
}
