package service

import (
	"errors"
	"time"

	msgmodel "github.com/aim/aim/internal/message/model"
	msgrepo "github.com/aim/aim/internal/message/repo"
	"github.com/aim/aim/internal/pkg/database"
)

var (
	ErrInvalidTargetType = errors.New("invalid target type")
	ErrMsgNotFound       = errors.New("message not found")
)

// MessageService handles message business logic.
type MessageService struct {
	repo *msgrepo.MessageRepo
}

// NewMessageService creates a new MessageService.
func NewMessageService(db *database.DB) *MessageService {
	return &MessageService{repo: msgrepo.NewMessageRepo(db)}
}

// Send sends a message and updates the session.
func (s *MessageService) Send(fromID, toID uint, targetType, msgType, content string) (*msgmodel.Message, error) {
	if targetType != msgmodel.TargetUser && targetType != msgmodel.TargetGroup {
		return nil, ErrInvalidTargetType
	}

	msg := &msgmodel.Message{
		FromID:     fromID,
		ToID:       toID,
		TargetType: targetType,
		MsgType:    msgType,
		Content:    content,
		Status:     msgmodel.MsgSent,
		CreatedAt:  time.Now(),
	}
	if err := s.repo.Create(msg); err != nil {
		return nil, err
	}

	// Upsert session for both participants (for single chat, update both sides).
	now := time.Now()
	session := &msgmodel.Session{
		UserID:      fromID,
		TargetID:    toID,
		TargetType:  targetType,
		LastMsg:     truncate(content, 100),
		LastMsgTime: now,
		UnreadCount: 0,
	}
	s.repo.UpsertSession(session)

	// For group, increment unread for all members except sender.
	// For single chat, set unread for the receiver.
	if targetType == msgmodel.TargetUser {
		receiverSession := &msgmodel.Session{
			UserID:      toID,
			TargetID:    fromID,
			TargetType:  targetType,
			LastMsg:     truncate(content, 100),
			LastMsgTime: now,
			UnreadCount: 1,
		}
		// Get existing session to preserve unread count.
		sessions, _ := s.repo.GetSessions(toID)
		for _, s := range sessions {
			if s.TargetID == fromID && s.TargetType == targetType {
				receiverSession.UnreadCount = s.UnreadCount + 1
				break
			}
		}
		s.repo.UpsertSession(receiverSession)
	}

	return msg, nil
}

// GetMessages retrieves messages for a session.
func (s *MessageService) GetMessages(userID, targetID uint, targetType string, offset, limit int) ([]msgmodel.Message, error) {
	return s.repo.FindBySession(userID, targetID, targetType, offset, limit)
}

// SearchMessages searches messages by keyword and filters.
func (s *MessageService) SearchMessages(keyword string, fromID, toID uint, startTime, endTime *time.Time, limit int) ([]msgmodel.Message, error) {
	return s.repo.SearchMessages(keyword, fromID, toID, startTime, endTime, limit)
}

// MarkRead marks all messages in a session as read.
func (s *MessageService) MarkRead(userID, targetID uint, targetType string) error {
	return s.repo.MarkSessionRead(userID, targetID, targetType)
}

// GetSessions returns all sessions for a user.
func (s *MessageService) GetSessions(userID uint) ([]msgmodel.Session, error) {
	return s.repo.GetSessions(userID)
}

func truncate(s string, max int) string {
	runes := []rune(s)
	if len(runes) <= max {
		return s
	}
	return string(runes[:max]) + "..."
}
