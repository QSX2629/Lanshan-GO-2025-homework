package service

import (
	"errors"
	"fmt"
	"strings"
	"time"

	aimodel "github.com/aim/aim/internal/ai/model"
	airepo "github.com/aim/aim/internal/ai/repo"
	"github.com/aim/aim/internal/pkg/database"
)

var (
	ErrBotNotFound    = errors.New("bot not found")
	ErrSessionExpired = errors.New("AI session expired")
)

// LLMClient defines the interface for calling LLM providers.
type LLMClient interface {
	Chat(model, systemPrompt string, history []ChatMessage, userMsg string) (string, int, error)
}

// ChatMessage is a single turn in a conversation.
type ChatMessage struct {
	Role    string
	Content string
}

// AIService handles AI assistant business logic.
type AIService struct {
	repo      *airepo.AIRepo
	llmClient LLMClient // injected, nil triggers fallback behavior
}

// NewAIService creates a new AIService.
func NewAIService(db *database.DB, llmClient LLMClient) *AIService {
	return &AIService{repo: airepo.NewAIRepo(db), llmClient: llmClient}
}

// CreateBot registers a new bot.
func (s *AIService) CreateBot(name, provider, model, systemPrompt string, isOfficial bool, pricePer1K float64) (*aimodel.Bot, error) {
	b := &aimodel.Bot{
		Name:         name,
		Provider:     provider,
		Model:        model,
		SystemPrompt: systemPrompt,
		IsOfficial:   isOfficial,
		PricePer1K:   pricePer1K,
	}
	if err := s.repo.CreateBot(b); err != nil {
		return nil, err
	}
	return b, nil
}

// ListBots returns all bots.
func (s *AIService) ListBots() ([]aimodel.Bot, error) {
	return s.repo.ListBots()
}

// GetBot returns a bot by ID.
func (s *AIService) GetBot(botID uint) (*aimodel.Bot, error) {
	b, err := s.repo.FindBotByID(botID)
	if err != nil {
		return nil, ErrBotNotFound
	}
	return b, nil
}

// Chat sends a message to the AI and returns the response.
func (s *AIService) Chat(userID, botID, targetID uint, targetType, message string) (string, error) {
	bot, err := s.GetBot(botID)
	if err != nil {
		return "", err
	}

	session, err := s.getOrCreateSession(userID, botID, targetID, targetType)
	if err != nil {
		return "", err
	}

	// Save user message.
	s.repo.AddChatRecord(&aimodel.AIChatRecord{
		SessionID:  session.ID,
		Role:       aimodel.BotRoleUser,
		Content:    message,
		TokenCount: estimateTokens(message),
		CreatedAt:  time.Now(),
	})

	history := s.buildHistory(session.ID)
	var reply string
	var tokens int

	if s.llmClient != nil {
		reply, tokens, err = s.llmClient.Chat(bot.Model, bot.SystemPrompt, history, message)
		if err != nil {
			return "", err
		}
	} else {
		// Fallback: echo with bot name when no LLM client is configured.
		reply = fmt.Sprintf("[%s] 收到消息: %s", bot.Name, message)
		tokens = estimateTokens(reply)
	}

	// Save assistant response.
	s.repo.AddChatRecord(&aimodel.AIChatRecord{
		SessionID:  session.ID,
		Role:       aimodel.BotRoleAssistant,
		Content:    reply,
		TokenCount: tokens,
		CreatedAt:  time.Now(),
	})

	// Update session context tokens.
	session.ContextTokens += estimateTokens(message) + tokens
	s.repo.UpdateAISession(session)

	// Record billing.
	s.recordBilling(userID, botID, tokens, bot.PricePer1K)

	return reply, nil
}

// Summarize generates a summary of recent chat history.
func (s *AIService) Summarize(userID, targetID uint, targetType string) (string, error) {
	history := s.buildContext(userID, targetID, targetType)
	if len(history) == 0 {
		return "没有找到对话记录", nil
	}

	prompt := "请总结以下对话内容的要点：\n" + strings.Join(history, "\n")

	if s.llmClient != nil {
		reply, _, err := s.llmClient.Chat("default", "你是一个对话总结助手。", nil, prompt)
		return reply, err
	}

	return fmt.Sprintf("对话摘要(共 %d 条消息):\n%s", len(history), strings.Join(history, "\n")), nil
}

// ExtractTodos extracts todo items from chat history.
func (s *AIService) ExtractTodos(userID, targetID uint, targetType string) (string, error) {
	history := s.buildContext(userID, targetID, targetType)
	if len(history) == 0 {
		return "没有找到对话记录", nil
	}

	prompt := "请从以下对话中提取待办事项：\n" + strings.Join(history, "\n")

	if s.llmClient != nil {
		reply, _, err := s.llmClient.Chat("default", "你是一个任务提取助手。请列出待办事项清单。", nil, prompt)
		return reply, err
	}

	return fmt.Sprintf("待办提取(共 %d 条消息):\n%s", len(history), strings.Join(history, "\n")), nil
}

// GetReplyCandidates generates reply suggestions.
func (s *AIService) GetReplyCandidates(userID, targetID uint, targetType string) ([]string, error) {
	history := s.buildContext(userID, targetID, targetType)

	prompt := "基于对话上下文，生成3个候选回复：\n" + strings.Join(history, "\n")

	if s.llmClient != nil {
		reply, _, err := s.llmClient.Chat("default", "你是一个智能回复建议助手。请生成3个简洁的候选回复，每行一个。", nil, prompt)
		if err != nil {
			return nil, err
		}
		return strings.Split(strings.TrimSpace(reply), "\n"), nil
	}

	return []string{"好的", "收到", "稍后回复"}, nil
}

// GetBilling returns billing records for a user.
func (s *AIService) GetBilling(userID uint, limit int) ([]aimodel.BillingRecord, error) {
	return s.repo.GetBillingRecords(userID, limit)
}

// --- helpers ---

func (s *AIService) getOrCreateSession(userID, botID, targetID uint, targetType string) (*aimodel.AISession, error) {
	// Look for existing active session.
	session := &aimodel.AISession{
		UserID:     userID,
		BotID:      botID,
		TargetID:   targetID,
		TargetType: targetType,
	}
	if err := s.repo.CreateAISession(session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *AIService) buildHistory(sessionID uint) []ChatMessage {
	records, _ := s.repo.GetChatHistory(sessionID, 50)
	history := make([]ChatMessage, 0, len(records))
	for _, r := range records {
		history = append(history, ChatMessage{Role: r.Role, Content: r.Content})
	}
	return history
}

func (s *AIService) buildContext(userID, targetID uint, targetType string) []string {
	// This is a simplified context builder.
	// In production, this would query the message service for recent messages.
	_ = userID
	_ = targetID
	_ = targetType
	return nil
}

func (s *AIService) recordBilling(userID, botID uint, tokensUsed int, pricePer1K float64) {
	cost := float64(tokensUsed) / 1000.0 * pricePer1K
	s.repo.CreateBillingRecord(&aimodel.BillingRecord{
		UserID:     userID,
		BotID:      botID,
		TokensUsed: tokensUsed,
		Cost:       cost,
		CreatedAt:  time.Now(),
	})
}

func estimateTokens(text string) int {
	// Rough estimate: ~4 chars per token for Chinese, ~4 for English.
	return len([]rune(text)) / 4
}
