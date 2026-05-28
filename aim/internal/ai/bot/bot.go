// Package bot provides the bot engine for processing @Bot messages.
package bot

import (
	"strings"

	"github.com/aim/aim/internal/ai/model"
	"github.com/aim/aim/internal/ai/service"
)

// Engine handles @Bot message detection and dispatch.
type Engine struct {
	aiSvc *service.AIService
}

// NewEngine creates a new bot Engine.
func NewEngine(aiSvc *service.AIService) *Engine {
	return &Engine{aiSvc: aiSvc}
}

// Mention represents a parsed @Bot mention.
type Mention struct {
	IsBot    bool
	BotID    uint
	BotName  string
	Content  string // message content without the @mention
	FullText string
}

// Parse detects if a message contains an @Bot mention.
func (e *Engine) Parse(text string, bots []model.Bot) *Mention {
	for _, bot := range bots {
		mention := "@" + bot.Name
		if strings.HasPrefix(strings.TrimSpace(text), mention) {
			content := strings.TrimSpace(strings.TrimPrefix(strings.TrimSpace(text), mention))
			return &Mention{
				IsBot:    true,
				BotID:    bot.ID,
				BotName:  bot.Name,
				Content:  content,
				FullText: text,
			}
		}
	}
	return &Mention{IsBot: false, FullText: text, Content: text}
}

// Process handles a detected @Bot mention and returns the bot's response.
func (e *Engine) Process(userID, botID, targetID uint, targetType, message string) (string, error) {
	return e.aiSvc.Chat(userID, botID, targetID, targetType, message)
}

// DetectAndProcess parses the message and, if it's a bot mention, processes it.
func (e *Engine) DetectAndProcess(userID, targetID uint, targetType, message string, bots []model.Bot) (string, bool, error) {
	m := e.Parse(message, bots)
	if !m.IsBot {
		return "", false, nil
	}
	reply, err := e.Process(userID, m.BotID, targetID, targetType, m.Content)
	return reply, true, err
}
