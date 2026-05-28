package repo

import (
	"time"

	msgmodel "github.com/aim/aim/internal/message/model"
	"github.com/aim/aim/internal/pkg/database"
	"gorm.io/gorm"
)

// MessageRepo handles message data access.
type MessageRepo struct {
	db *database.DB
}

// NewMessageRepo creates a new MessageRepo.
func NewMessageRepo(db *database.DB) *MessageRepo {
	return &MessageRepo{db: db}
}

// AutoMigrate creates the message tables.
func (r *MessageRepo) AutoMigrate() error {
	return r.db.AutoMigrate(&msgmodel.Message{}, &msgmodel.Session{}, &msgmodel.ReadReceipt{})
}

// Create inserts a new message.
func (r *MessageRepo) Create(msg *msgmodel.Message) error {
	return r.db.Create(msg).Error
}

// FindBySession retrieves messages for a session with pagination.
func (r *MessageRepo) FindBySession(userID, targetID uint, targetType string, offset, limit int) ([]msgmodel.Message, error) {
	var msgs []msgmodel.Message
	err := r.db.Where(
		"((from_id = ? AND to_id = ?) OR (from_id = ? AND to_id = ?)) AND target_type = ?",
		userID, targetID, targetID, userID, targetType,
	).Order("created_at DESC").Offset(offset).Limit(limit).Find(&msgs).Error
	return msgs, err
}

// SearchMessages searches messages by keyword within a time range.
func (r *MessageRepo) SearchMessages(keyword string, fromID, toID uint, startTime, endTime *time.Time, limit int) ([]msgmodel.Message, error) {
	query := r.db.Model(&msgmodel.Message{}).Where("content LIKE ?", "%"+keyword+"%")
	if fromID > 0 {
		query = query.Where("from_id = ?", fromID)
	}
	if toID > 0 {
		query = query.Where("to_id = ?", toID)
	}
	if startTime != nil {
		query = query.Where("created_at >= ?", startTime)
	}
	if endTime != nil {
		query = query.Where("created_at <= ?", endTime)
	}
	var msgs []msgmodel.Message
	err := query.Order("created_at DESC").Limit(limit).Find(&msgs).Error
	return msgs, err
}

// UpdateStatus updates the status of a message.
func (r *MessageRepo) UpdateStatus(msgID uint, status string) error {
	return r.db.Model(&msgmodel.Message{}).Where("id = ?", msgID).Update("status", status).Error
}

// CreateReadReceipt inserts a read receipt.
func (r *MessageRepo) CreateReadReceipt(receipt *msgmodel.ReadReceipt) error {
	return r.db.Create(receipt).Error
}

// UpsertSession creates or updates a conversation session.
func (r *MessageRepo) UpsertSession(session *msgmodel.Session) error {
	var existing msgmodel.Session
	err := r.db.Where("user_id = ? AND target_id = ? AND target_type = ?",
		session.UserID, session.TargetID, session.TargetType).First(&existing).Error
	if err == gorm.ErrRecordNotFound {
		return r.db.Create(session).Error
	}
	if err != nil {
		return err
	}
	existing.LastMsg = session.LastMsg
	existing.LastMsgTime = session.LastMsgTime
	existing.UnreadCount = session.UnreadCount
	return r.db.Save(&existing).Error
}

// GetSessions returns all sessions for a user ordered by last message time.
func (r *MessageRepo) GetSessions(userID uint) ([]msgmodel.Session, error) {
	var sessions []msgmodel.Session
	err := r.db.Where("user_id = ?", userID).Order("last_msg_time DESC").Find(&sessions).Error
	return sessions, err
}

// MarkSessionRead resets the unread count for a session.
func (r *MessageRepo) MarkSessionRead(userID, targetID uint, targetType string) error {
	return r.db.Model(&msgmodel.Session{}).
		Where("user_id = ? AND target_id = ? AND target_type = ?", userID, targetID, targetType).
		Update("unread_count", 0).Error
}
