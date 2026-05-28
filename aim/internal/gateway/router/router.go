package router

import (
	"github.com/aim/aim/internal/gateway/session"
	"github.com/aim/aim/internal/gateway/ws"
	"github.com/aim/aim/internal/pkg/protocol"
)

// Router dispatches incoming WebSocket frames to the appropriate handler.
type Router struct {
	session *session.Manager
}

// NewRouter creates a new Router.
func NewRouter(sm *session.Manager) *Router {
	return &Router{session: sm}
}

// Route processes an incoming raw message and routes it accordingly.
func (r *Router) Route(raw *ws.RawMessage) {
	frame, err := protocol.Decode(raw.Data)
	if err != nil {
		r.sendError(raw.ConnID, "decode error")
		return
	}

	switch frame.Cmd {
	case protocol.CmdChatMsg:
		r.handleChatMsg(raw, frame)
	case protocol.CmdChatAck:
		r.handleAck(raw, frame)
	case protocol.CmdTyping:
		r.handleTyping(raw, frame)
	case protocol.CmdHeartbeat:
		r.handleHeartbeat(raw.ConnID)
	case protocol.CmdAIBotMsg:
		r.handleAIBotMsg(raw, frame)
	default:
		r.sendError(raw.ConnID, "unknown command")
	}
}

func (r *Router) handleChatMsg(raw *ws.RawMessage, frame *protocol.Frame) {
	targetID := frame.To
	if targetID == "" {
		r.sendError(raw.ConnID, "missing target")
		return
	}

	out, _ := protocol.Encode(*frame)
	r.session.SendToUser(raw.UserID, out)

	if frame.MsgType == protocol.TypeText {
		_ = targetID
	}
}

func (r *Router) handleAck(raw *ws.RawMessage, frame *protocol.Frame) {
	out, _ := protocol.Encode(*frame)
	r.session.SendToUser(raw.UserID, out)
}

func (r *Router) handleTyping(raw *ws.RawMessage, frame *protocol.Frame) {
	out, _ := protocol.Encode(*frame)
	r.session.SendToUser(raw.UserID, out)
}

func (r *Router) handleHeartbeat(connID uint) {
	conn := r.session.GetConn(connID)
	if conn != nil {
		pong := protocol.Frame{Cmd: protocol.CmdHeartbeat}
		data, _ := protocol.Encode(pong)
		conn.Send(data)
	}
}

func (r *Router) handleAIBotMsg(raw *ws.RawMessage, frame *protocol.Frame) {
	out, _ := protocol.Encode(*frame)
	r.session.SendToUser(raw.UserID, out)
}

func (r *Router) sendError(connID uint, msg string) {
	conn := r.session.GetConn(connID)
	if conn != nil {
		errFrame := protocol.Frame{Cmd: protocol.CmdError, Content: msg}
		data, _ := protocol.Encode(errFrame)
		conn.Send(data)
	}
}
