package repo

import (
	"testing"
	"time"

	msgmodel "github.com/aim/aim/internal/message/model"
	"github.com/aim/aim/internal/pkg/database"
)

func setupMessageRepo(t *testing.T) *MessageRepo {
	t.Helper()
	db, err := database.TestDB()
	if err != nil {
		t.Fatalf("TestDB() error = %v", err)
	}
	repo := NewMessageRepo(db)
	if err := repo.AutoMigrate(); err != nil {
		t.Fatalf("AutoMigrate() error = %v", err)
	}
	return repo
}

func TestMessageRepo_CreateAndFind(t *testing.T) {
	repo := setupMessageRepo(t)

	msg := &msgmodel.Message{
		Seq:        1,
		FromID:     1,
		ToID:       2,
		TargetType: msgmodel.TargetUser,
		MsgType:    "text",
		Content:    "hello",
	}
	if err := repo.Create(msg); err != nil {
		t.Fatalf("Create() error = %v", err)
	}

	msgs, err := repo.FindBySession(1, 2, msgmodel.TargetUser, 0, 20)
	if err != nil {
		t.Fatalf("FindBySession() error = %v", err)
	}
	if len(msgs) != 1 {
		t.Fatalf("len(msgs) = %d, want 1", len(msgs))
	}
	if msgs[0].Content != "hello" {
		t.Errorf("Content = %q, want hello", msgs[0].Content)
	}
}

func TestMessageRepo_Search(t *testing.T) {
	repo := setupMessageRepo(t)

	repo.Create(&msgmodel.Message{Seq: 1, FromID: 1, ToID: 2, TargetType: msgmodel.TargetUser, Content: "search target"})
	repo.Create(&msgmodel.Message{Seq: 2, FromID: 1, ToID: 2, TargetType: msgmodel.TargetUser, Content: "other"})

	msgs, err := repo.SearchMessages("search", 1, 0, nil, nil, 20)
	if err != nil {
		t.Fatalf("SearchMessages() error = %v", err)
	}
	if len(msgs) != 1 {
		t.Errorf("len(msgs) = %d, want 1", len(msgs))
	}
}

func TestMessageRepo_UpdateStatus(t *testing.T) {
	repo := setupMessageRepo(t)

	msg := &msgmodel.Message{Seq: 1, FromID: 1, ToID: 2, TargetType: msgmodel.TargetUser, Content: "test"}
	repo.Create(msg)

	if err := repo.UpdateStatus(msg.ID, msgmodel.MsgRead); err != nil {
		t.Fatalf("UpdateStatus() error = %v", err)
	}
}

func TestMessageRepo_Session(t *testing.T) {
	repo := setupMessageRepo(t)

	s := &msgmodel.Session{
		UserID:      1,
		TargetID:    2,
		TargetType:  msgmodel.TargetUser,
		LastMsg:     "hello",
		LastMsgTime: time.Now(),
		UnreadCount: 3,
	}
	if err := repo.UpsertSession(s); err != nil {
		t.Fatalf("UpsertSession() error = %v", err)
	}

	sessions, err := repo.GetSessions(1)
	if err != nil {
		t.Fatalf("GetSessions() error = %v", err)
	}
	if len(sessions) != 1 {
		t.Fatalf("len(sessions) = %d, want 1", len(sessions))
	}
	if sessions[0].UnreadCount != 3 {
		t.Errorf("UnreadCount = %d, want 3", sessions[0].UnreadCount)
	}

	// Mark as read.
	if err := repo.MarkSessionRead(1, 2, msgmodel.TargetUser); err != nil {
		t.Fatalf("MarkSessionRead() error = %v", err)
	}

	sessions, _ = repo.GetSessions(1)
	if sessions[0].UnreadCount != 0 {
		t.Errorf("UnreadCount = %d, want 0 (after mark read)", sessions[0].UnreadCount)
	}
}

func TestMessageRepo_ReadReceipt(t *testing.T) {
	repo := setupMessageRepo(t)

	r := &msgmodel.ReadReceipt{MsgID: 1, UserID: 2, ReadAt: time.Now()}
	if err := repo.CreateReadReceipt(r); err != nil {
		t.Fatalf("CreateReadReceipt() error = %v", err)
	}
}
