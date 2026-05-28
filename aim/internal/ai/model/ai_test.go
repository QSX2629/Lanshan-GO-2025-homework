package model

import (
	"testing"
	"time"
)

func TestBotFields(t *testing.T) {
	b := Bot{
		Name:         "AIM Bot",
		Provider:     "openai",
		Model:        "gpt-4",
		SystemPrompt: "You are a helpful assistant.",
		APIKey:       "sk-secret",
		IsOfficial:   true,
		PricePer1K:   0.03,
	}

	if b.Provider != "openai" {
		t.Errorf("Provider = %q, want openai", b.Provider)
	}
	if b.Model != "gpt-4" {
		t.Errorf("Model = %q, want gpt-4", b.Model)
	}
	if !b.IsOfficial {
		t.Error("IsOfficial should be true")
	}
	// APIKey should not be serialized (json:"-")
	if b.APIKey != "sk-secret" {
		t.Error("APIKey mismatch")
	}
}

func TestBotRoleConstants(t *testing.T) {
	if BotRoleUser != "user" || BotRoleAssistant != "assistant" {
		t.Error("bot role constants mismatch")
	}
}

func TestAISessionFields(t *testing.T) {
	s := AISession{
		UserID:        1,
		BotID:         1,
		TargetID:      2,
		TargetType:    "group",
		ContextTokens: 4096,
	}

	if s.UserID != 1 || s.BotID != 1 {
		t.Error("AISession fields mismatch")
	}
	if s.TargetType != "group" {
		t.Errorf("TargetType = %q, want group", s.TargetType)
	}
}

func TestAIChatRecord(t *testing.T) {
	r := AIChatRecord{
		SessionID:  1,
		Role:       BotRoleUser,
		Content:    "hello",
		TokenCount: 10,
	}

	if r.Role != BotRoleUser {
		t.Errorf("Role = %q, want %q", r.Role, BotRoleUser)
	}
}

func TestBillingRecord(t *testing.T) {
	br := BillingRecord{
		UserID:     1,
		BotID:      1,
		TokensUsed: 500,
		Cost:       0.015,
		CreatedAt:  time.Now(),
	}

	if br.TokensUsed != 500 {
		t.Errorf("TokensUsed = %d, want 500", br.TokensUsed)
	}
	if br.Cost != 0.015 {
		t.Errorf("Cost = %f, want 0.015", br.Cost)
	}
}
