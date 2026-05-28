// Package protocol defines the WebSocket message frame format for AIM.
package protocol

import "encoding/json"

// MsgType defines the type of a message.
type MsgType string

const (
	TypeText  MsgType = "text"
	TypeImage MsgType = "image"
	TypeFile  MsgType = "file"
	TypeVoice MsgType = "voice"
	TypeEvent MsgType = "event"
)

// CmdType defines command/event types carried over the wire.
type CmdType string

const (
	CmdChatMsg      CmdType = "chat.msg"      // Chat message (single/group)
	CmdChatAck      CmdType = "chat.ack"      // Read receipt
	CmdTyping       CmdType = "chat.typing"   // Typing indicator
	CmdOnlineChange CmdType = "user.online"   // Online status change
	CmdHeartbeat    CmdType = "heartbeat"     // Keep-alive ping/pong
	CmdError        CmdType = "error"         // Error frame
	CmdAIBotMsg     CmdType = "ai.bot.msg"    // AI bot message
	CmdAIBotStream  CmdType = "ai.bot.stream" // AI streaming chunk
)

// Frame is the wire-format message exchanged over WebSocket.
type Frame struct {
	Cmd     CmdType         `json:"cmd"`                // Command/event type
	Seq     int64           `json:"seq,omitempty"`      // Sequence number for ack
	From    string          `json:"from,omitempty"`     // Sender user ID
	To      string          `json:"to,omitempty"`       // Target user/group ID
	MsgType MsgType         `json:"msg_type,omitempty"` // Message content type
	Content string          `json:"content,omitempty"`  // Payload (text or JSON)
	Extra   json.RawMessage `json:"extra,omitempty"`    // Optional metadata
}

// ChatContent carries structured chat payload inside Frame.Content.
type ChatContent struct {
	Text     string `json:"text,omitempty"`
	FileURL  string `json:"file_url,omitempty"`
	FileName string `json:"file_name,omitempty"`
	FileSize int64  `json:"file_size,omitempty"`
	Duration int    `json:"duration,omitempty"` // Voice duration (seconds)
}

// TypingPayload is the payload for CmdTyping frames.
type TypingPayload struct {
	UserID     string `json:"user_id"`
	TargetID   string `json:"target_id"`   // User or group ID
	TargetType string `json:"target_type"` // "user" or "group"
	Typing     bool   `json:"typing"`
}

// AckPayload is the payload for CmdChatAck frames.
type AckPayload struct {
	MsgIDs []string `json:"msg_ids"`
}

// Encode serializes a Frame to JSON bytes.
func Encode(f Frame) ([]byte, error) {
	return json.Marshal(f)
}

// Decode deserializes a Frame from JSON bytes.
func Decode(data []byte) (*Frame, error) {
	f := &Frame{}
	if err := json.Unmarshal(data, f); err != nil {
		return nil, err
	}
	return f, nil
}

// ValidMsgType reports whether t is a known message type.
func ValidMsgType(t MsgType) bool {
	switch t {
	case TypeText, TypeImage, TypeFile, TypeVoice, TypeEvent:
		return true
	default:
		return false
	}
}
