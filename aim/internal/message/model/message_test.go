package model

import (
	"testing"
	"time"
)

func TestMessageFields(t *testing.T) {
	now := time.Now()
	m := Message{
		ID:         1,
		Seq:        100,
		FromID:     1,
		ToID:       2,
		TargetType: TargetUser,
		MsgType:    "text",
		Content:    "hello",
		Status:     MsgSent,
		CreatedAt:  now,
	}

	if m.TargetType != TargetUser {
		t.Errorf("TargetType = %q, want %q", m.TargetType, TargetUser)
	}
	if m.MsgType != "text" {
		t.Errorf("MsgType = %q, want text", m.MsgType)
	}
	if m.Status != MsgSent {
		t.Errorf("Status = %q, want %q", m.Status, MsgSent)
	}
}

func TestMessageStatusConstants(t *testing.T) {
	if MsgSent != "sent" || MsgDelivered != "delivered" || MsgRead != "read" {
		t.Error("message status constants mismatch")
	}
}

func TestSessionFields(t *testing.T) {
	s := Session{
		UserID:      1,
		TargetID:    2,
		TargetType:  TargetGroup,
		LastMsg:     "last message",
		UnreadCount: 5,
	}

	if s.UnreadCount != 5 {
		t.Errorf("UnreadCount = %d, want 5", s.UnreadCount)
	}
	if s.TargetType != TargetGroup {
		t.Errorf("TargetType = %q, want %q", s.TargetType, TargetGroup)
	}
}

func TestReadReceipt(t *testing.T) {
	now := time.Now()
	r := ReadReceipt{MsgID: 1, UserID: 2, ReadAt: now}
	if r.MsgID != 1 || r.UserID != 2 {
		t.Errorf("ReadReceipt fields mismatch")
	}
}
