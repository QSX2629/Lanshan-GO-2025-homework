package handler

import (
	"context"
	"time"

	"github.com/aim/aim/internal/message/model"
	"github.com/aim/aim/internal/message/service"
	"github.com/aim/aim/internal/pkg/database"
)

// MessageHandler provides the gRPC-compatible handler for message operations.
type MessageHandler struct {
	svc *service.MessageService
}

// NewMessageHandler creates a new MessageHandler.
func NewMessageHandler(db *database.DB) *MessageHandler {
	return &MessageHandler{svc: service.NewMessageService(db)}
}

// SendRequest is the input for sending a message.
type SendRequest struct {
	FromID     uint   `json:"from_id"`
	ToID       uint   `json:"to_id"`
	TargetType string `json:"target_type"`
	MsgType    string `json:"msg_type"`
	Content    string `json:"content"`
}

// GetMessagesRequest is the paginated query input.
type GetMessagesRequest struct {
	UserID     uint   `json:"user_id"`
	TargetID   uint   `json:"target_id"`
	TargetType string `json:"target_type"`
	Offset     int    `json:"offset"`
	Limit      int    `json:"limit"`
}

// SearchRequest is the search query input.
type SearchRequest struct {
	Keyword   string     `json:"keyword"`
	FromID    uint       `json:"from_id"`
	ToID      uint       `json:"to_id"`
	StartTime *time.Time `json:"start_time,omitempty"`
	EndTime   *time.Time `json:"end_time,omitempty"`
	Limit     int        `json:"limit"`
}

// SessionResponse is the public session view.
type SessionResponse struct {
	UserID      uint      `json:"user_id"`
	TargetID    uint      `json:"target_id"`
	TargetType  string    `json:"target_type"`
	LastMsg     string    `json:"last_msg"`
	LastMsgTime time.Time `json:"last_msg_time"`
	UnreadCount int       `json:"unread_count"`
}

// Send handles sending a message.
func (h *MessageHandler) Send(_ context.Context, req *SendRequest) (*model.Message, error) {
	return h.svc.Send(req.FromID, req.ToID, req.TargetType, req.MsgType, req.Content)
}

// GetMessages retrieves messages for a session.
func (h *MessageHandler) GetMessages(_ context.Context, req *GetMessagesRequest) ([]model.Message, error) {
	return h.svc.GetMessages(req.UserID, req.TargetID, req.TargetType, req.Offset, req.Limit)
}

// SearchMessages searches messages.
func (h *MessageHandler) SearchMessages(_ context.Context, req *SearchRequest) ([]model.Message, error) {
	return h.svc.SearchMessages(req.Keyword, req.FromID, req.ToID, req.StartTime, req.EndTime, req.Limit)
}

// MarkRead marks a session as read.
func (h *MessageHandler) MarkRead(_ context.Context, userID, targetID uint, targetType string) error {
	return h.svc.MarkRead(userID, targetID, targetType)
}

// GetSessions returns all sessions for a user.
func (h *MessageHandler) GetSessions(_ context.Context, userID uint) ([]SessionResponse, error) {
	sessions, err := h.svc.GetSessions(userID)
	if err != nil {
		return nil, err
	}
	out := make([]SessionResponse, 0, len(sessions))
	for _, s := range sessions {
		out = append(out, SessionResponse{
			UserID:      s.UserID,
			TargetID:    s.TargetID,
			TargetType:  s.TargetType,
			LastMsg:     s.LastMsg,
			LastMsgTime: s.LastMsgTime,
			UnreadCount: s.UnreadCount,
		})
	}
	return out, nil
}
