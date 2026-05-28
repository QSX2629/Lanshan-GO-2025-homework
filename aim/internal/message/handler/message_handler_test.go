package handler

import (
	"testing"

	msgmodel "github.com/aim/aim/internal/message/model"
	msgrepo "github.com/aim/aim/internal/message/repo"
	"github.com/aim/aim/internal/pkg/database"
)

func setupMessageHandler(t *testing.T) *MessageHandler {
	t.Helper()
	db, err := database.TestDB()
	if err != nil {
		t.Fatalf("TestDB() error = %v", err)
	}
	repo := msgrepo.NewMessageRepo(db)
	if err := repo.AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate() error = %v", err)
	}
	return NewMessageHandler(db)
}

func TestMessageHandler_SendAndGetMessages(t *testing.T) {
	h := setupMessageHandler(t)

	msg, err := h.Send(nil, &SendRequest{
		FromID: 1, ToID: 2, TargetType: msgmodel.TargetUser,
		MsgType: "text", Content: "hello",
	})
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if msg.Content != "hello" {
		t.Errorf("Content = %q, want hello", msg.Content)
	}

	msgs, err := h.GetMessages(nil, &GetMessagesRequest{
		UserID: 1, TargetID: 2, TargetType: msgmodel.TargetUser, Offset: 0, Limit: 20,
	})
	if err != nil {
		t.Fatalf("GetMessages() error = %v", err)
	}
	if len(msgs) != 1 {
		t.Errorf("len(msgs) = %d, want 1", len(msgs))
	}
}

func TestMessageHandler_Search(t *testing.T) {
	h := setupMessageHandler(t)
	h.Send(nil, &SendRequest{FromID: 1, ToID: 2, TargetType: msgmodel.TargetUser, MsgType: "text", Content: "find me"})

	results, err := h.SearchMessages(nil, &SearchRequest{Keyword: "find", FromID: 1, Limit: 20})
	if err != nil {
		t.Fatalf("SearchMessages() error = %v", err)
	}
	if len(results) == 0 {
		t.Error("search should return results")
	}
}

func TestMessageHandler_MarkRead(t *testing.T) {
	h := setupMessageHandler(t)
	h.Send(nil, &SendRequest{FromID: 1, ToID: 2, TargetType: msgmodel.TargetUser, MsgType: "text", Content: "msg"})

	if err := h.MarkRead(nil, 2, 1, msgmodel.TargetUser); err != nil {
		t.Fatalf("MarkRead() error = %v", err)
	}
}

func TestMessageHandler_GetSessions(t *testing.T) {
	h := setupMessageHandler(t)
	h.Send(nil, &SendRequest{FromID: 1, ToID: 2, TargetType: msgmodel.TargetUser, MsgType: "text", Content: "hi"})

	sessions, err := h.GetSessions(nil, 1)
	if err != nil {
		t.Fatalf("GetSessions() error = %v", err)
	}
	if len(sessions) < 1 {
		t.Errorf("len(sessions) = %d, want >= 1", len(sessions))
	}
}
