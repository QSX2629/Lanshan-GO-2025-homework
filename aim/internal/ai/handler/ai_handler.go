package handler

import (
	"context"

	aimodel "github.com/aim/aim/internal/ai/model"
	aiservice "github.com/aim/aim/internal/ai/service"
	"github.com/aim/aim/internal/pkg/database"
)

// AIHandler provides handlers for AI operations (gRPC + OpenAPI).
type AIHandler struct {
	svc *aiservice.AIService
}

// NewAIHandler creates a new AIHandler.
func NewAIHandler(db *database.DB, llmClient aiservice.LLMClient) *AIHandler {
	return &AIHandler{svc: aiservice.NewAIService(db, llmClient)}
}

// CreateBotRequest is the input for creating a bot.
type CreateBotRequest struct {
	Name         string  `json:"name"`
	Provider     string  `json:"provider"`
	Model        string  `json:"model"`
	SystemPrompt string  `json:"system_prompt"`
	IsOfficial   bool    `json:"is_official"`
	PricePer1K   float64 `json:"price_per_1k"`
}

// ChatRequest is the input for an AI chat turn.
type ChatRequest struct {
	UserID     uint   `json:"user_id"`
	BotID      uint   `json:"bot_id"`
	TargetID   uint   `json:"target_id"`
	TargetType string `json:"target_type"`
	Message    string `json:"message"`
}

// AIActionRequest is used for summarize/todos/candidates.
type AIActionRequest struct {
	UserID     uint   `json:"user_id"`
	TargetID   uint   `json:"target_id"`
	TargetType string `json:"target_type"`
}

// CreateBot registers a new bot.
func (h *AIHandler) CreateBot(_ context.Context, req *CreateBotRequest) (*aimodel.Bot, error) {
	return h.svc.CreateBot(req.Name, req.Provider, req.Model, req.SystemPrompt, req.IsOfficial, req.PricePer1K)
}

// ListBots returns all bots.
func (h *AIHandler) ListBots(_ context.Context) ([]aimodel.Bot, error) {
	return h.svc.ListBots()
}

// Chat sends a message to an AI bot.
func (h *AIHandler) Chat(_ context.Context, req *ChatRequest) (string, error) {
	return h.svc.Chat(req.UserID, req.BotID, req.TargetID, req.TargetType, req.Message)
}

// Summarize generates a conversation summary.
func (h *AIHandler) Summarize(_ context.Context, req *AIActionRequest) (string, error) {
	return h.svc.Summarize(req.UserID, req.TargetID, req.TargetType)
}

// ExtractTodos extracts todo items.
func (h *AIHandler) ExtractTodos(_ context.Context, req *AIActionRequest) (string, error) {
	return h.svc.ExtractTodos(req.UserID, req.TargetID, req.TargetType)
}

// GetReplyCandidates generates reply suggestions.
func (h *AIHandler) GetReplyCandidates(_ context.Context, req *AIActionRequest) ([]string, error) {
	return h.svc.GetReplyCandidates(req.UserID, req.TargetID, req.TargetType)
}

// GetBilling returns billing records.
func (h *AIHandler) GetBilling(_ context.Context, userID uint, limit int) ([]aimodel.BillingRecord, error) {
	return h.svc.GetBilling(userID, limit)
}
