package service

import (
	"testing"
	"time"

	msgmodel "github.com/aim/aim/internal/message/model"
	msgrepo "github.com/aim/aim/internal/message/repo"
	"github.com/aim/aim/internal/pkg/database"
)

func setupMessageService(t *testing.T) *MessageService {
	t.Helper()
	db, err := database.TestDB()
	if err != nil {
		t.Fatalf("TestDB() error = %v", err)
	}
	repo := msgrepo.NewMessageRepo(db)
	if err := repo.AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate() error = %v", err)
	}
	return NewMessageService(db)
}

func TestMessageService_SendAndGet(t *testing.T) {
	svc := setupMessageService(t)

	msg, err := svc.Send(1, 2, msgmodel.TargetUser, "text", "hello world")
	if err != nil {
		t.Fatalf("Send() error = %v", err)
	}
	if msg.Content != "hello world" {
		t.Errorf("Content = %q, want hello world", msg.Content)
	}
	if msg.Status != msgmodel.MsgSent {
		t.Errorf("Status = %q, want %q", msg.Status, msgmodel.MsgSent)
	}

	msgs, err := svc.GetMessages(1, 2, msgmodel.TargetUser, 0, 20)
	if err != nil {
		t.Fatalf("GetMessages() error = %v", err)
	}
	if len(msgs) != 1 {
		t.Errorf("len(msgs) = %d, want 1", len(msgs))
	}
}

func TestMessageService_SendGroup(t *testing.T) {
	svc := setupMessageService(t)

	_, err := svc.Send(1, 100, msgmodel.TargetGroup, "text", "group msg")
	if err != nil {
		t.Fatalf("Send() to group error = %v", err)
	}

	msgs, err := svc.GetMessages(1, 100, msgmodel.TargetGroup, 0, 20)
	if err != nil {
		t.Fatalf("GetMessages() error = %v", err)
	}
	if len(msgs) != 1 {
		t.Errorf("len(msgs) = %d, want 1", len(msgs))
	}
}

func TestMessageService_InvalidTarget(t *testing.T) {
	svc := setupMessageService(t)

	_, err := svc.Send(1, 2, "invalid", "text", "x")
	if err != ErrInvalidTargetType {
		t.Errorf("error = %v, want ErrInvalidTargetType", err)
	}
}

func TestMessageService_Search(t *testing.T) {
	svc := setupMessageService(t)
	svc.Send(1, 2, msgmodel.TargetUser, "text", "unique keyword test")
	svc.Send(1, 2, msgmodel.TargetUser, "text", "other message")

	results, err := svc.SearchMessages("unique", 1, 0, nil, nil, 20)
	if err != nil {
		t.Fatalf("SearchMessages() error = %v", err)
	}
	if len(results) != 1 {
		t.Errorf("len(results) = %d, want 1", len(results))
	}
}

func TestMessageService_MarkRead(t *testing.T) {
	svc := setupMessageService(t)
	svc.Send(1, 2, msgmodel.TargetUser, "text", "msg1")

	if err := svc.MarkRead(2, 1, msgmodel.TargetUser); err != nil {
		t.Fatalf("MarkRead() error = %v", err)
	}
}

func TestMessageService_GetSessions(t *testing.T) {
	svc := setupMessageService(t)
	svc.Send(1, 2, msgmodel.TargetUser, "text", "hello")

	sessions, err := svc.GetSessions(1)
	if err != nil {
		t.Fatalf("GetSessions() error = %v", err)
	}
	if len(sessions) < 1 {
		t.Errorf("len(sessions) = %d, want >= 1", len(sessions))
	}
}

func TestTruncate(t *testing.T) {
	if s := truncate("hello", 100); s != "hello" {
		t.Errorf("truncate short = %q, want hello", s)
	}
	longMsg := "你好世界测试消息内容非常长超过一百个字符" +
		"你好世界测试消息内容非常长超过一百个字符" +
		"你好世界测试消息内容非常长超过一百个字符" +
		"你好世界测试消息内容非常长超过一百个字符"
	s := truncate(longMsg, 10)
	if s == longMsg {
		t.Error("truncate should have shortened the message")
	}

	if s2 := truncate("", 5); s2 != "" {
		t.Errorf("truncate empty = %q, want empty", s2)
	}

	_ = time.Now
}
