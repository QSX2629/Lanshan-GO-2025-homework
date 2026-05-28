package router

import (
	"testing"

	"github.com/aim/aim/internal/gateway/session"
	"github.com/aim/aim/internal/gateway/ws"
	"github.com/aim/aim/internal/pkg/protocol"
)

func TestRouter_Heartbeat(t *testing.T) {
	sm := session.NewManager()
	// Add a test conn that can receive the pong.
	c := &ws.Conn{ID: 1, UserID: 100}
	sm.Add(c)

	r := NewRouter(sm)
	raw := &ws.RawMessage{
		ConnID: 1,
		UserID: 100,
		Data:   mustEncode(protocol.Frame{Cmd: protocol.CmdHeartbeat}),
	}
	r.Route(raw)
	// No panic = pass.
}

func TestRouter_InvalidFrame(t *testing.T) {
	sm := session.NewManager()
	c := &ws.Conn{ID: 1, UserID: 100}
	sm.Add(c)

	r := NewRouter(sm)
	r.Route(&ws.RawMessage{ConnID: 1, UserID: 100, Data: []byte("bad json")})
	// No panic = pass.
}

func TestRouter_ChatMsg(t *testing.T) {
	sm := session.NewManager()
	c := &ws.Conn{ID: 1, UserID: 100}
	sm.Add(c)

	r := NewRouter(sm)
	raw := &ws.RawMessage{
		ConnID: 1,
		UserID: 100,
		Data: mustEncode(protocol.Frame{
			Cmd:     protocol.CmdChatMsg,
			Seq:     1,
			From:    "100",
			To:      "200",
			MsgType: protocol.TypeText,
			Content: "hello",
		}),
	}
	r.Route(raw)
}

func TestRouter_ChatMsg_MissingTarget(t *testing.T) {
	sm := session.NewManager()
	c := &ws.Conn{ID: 1, UserID: 100}
	sm.Add(c)

	r := NewRouter(sm)
	r.Route(&ws.RawMessage{
		ConnID: 1,
		UserID: 100,
		Data: mustEncode(protocol.Frame{
			Cmd:     protocol.CmdChatMsg,
			MsgType: protocol.TypeText,
			Content: "hello",
		}),
	})
}

func TestRouter_UnknownCmd(t *testing.T) {
	sm := session.NewManager()
	c := &ws.Conn{ID: 1, UserID: 100}
	sm.Add(c)

	r := NewRouter(sm)
	r.Route(&ws.RawMessage{
		ConnID: 1,
		UserID: 100,
		Data:   mustEncode(protocol.Frame{Cmd: protocol.CmdType("unknown.cmd")}),
	})
}

func TestRouter_SendError(t *testing.T) {
	sm := session.NewManager()
	c := &ws.Conn{ID: 1, UserID: 100}
	sm.Add(c)

	r := NewRouter(sm)
	r.sendError(1, "test error")
	r.sendError(999, "no conn") // Should not panic with unknown conn.
}

func mustEncode(f protocol.Frame) []byte {
	data, _ := protocol.Encode(f)
	return data
}
