package bot

import (
	"testing"

	aimodel "github.com/aim/aim/internal/ai/model"
)

func TestEngine_Parse(t *testing.T) {
	e := &Engine{}

	bots := []aimodel.Bot{
		{ID: 1, Name: "AIMBot"},
		{ID: 2, Name: "Helper"},
	}

	tests := []struct {
		text        string
		wantBot     bool
		wantBotID   uint
		wantContent string
	}{
		{"@AIMBot 你好", true, 1, "你好"},
		{"@AIMBot 今天天气怎么样", true, 1, "今天天气怎么样"},
		{"@Helper help me", true, 2, "help me"},
		{"你好 @AIMBot", false, 0, "你好 @AIMBot"}, // Not at start.
		{"普通消息", false, 0, "普通消息"},
		{"", false, 0, ""},
	}

	for _, tt := range tests {
		m := e.Parse(tt.text, bots)
		if m.IsBot != tt.wantBot {
			t.Errorf("Parse(%q).IsBot = %v, want %v", tt.text, m.IsBot, tt.wantBot)
		}
		if m.IsBot && m.BotID != tt.wantBotID {
			t.Errorf("Parse(%q).BotID = %d, want %d", tt.text, m.BotID, tt.wantBotID)
		}
		if tt.wantBot && m.Content != tt.wantContent {
			t.Errorf("Parse(%q).Content = %q, want %q", tt.text, m.Content, tt.wantContent)
		}
	}
}

func TestEngine_Parse_EmptyBots(t *testing.T) {
	e := &Engine{}
	m := e.Parse("@Bot hello", nil)
	if m.IsBot {
		t.Error("no bots registered, should not match")
	}
}

func TestEngine_DetectAndProcess_NotBot(t *testing.T) {
	e := &Engine{}
	_, handled, err := e.DetectAndProcess(1, 2, "user", "hello", nil)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if handled {
		t.Error("should not handle non-bot message")
	}
}
