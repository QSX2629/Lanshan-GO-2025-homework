package session

import (
	"testing"

	"github.com/aim/aim/internal/gateway/ws"
)

func newTestConn(id, userID uint) *ws.Conn {
	return &ws.Conn{ID: id, UserID: userID}
}

func TestManager_AddAndRemove(t *testing.T) {
	m := NewManager()
	c1 := newTestConn(1, 100)
	m.Add(c1)
	m.Add(newTestConn(2, 100))
	m.Add(newTestConn(3, 200))

	if !m.IsOnline(100) {
		t.Error("user 100 should be online")
	}
	if !m.IsOnline(200) {
		t.Error("user 200 should be online")
	}
	if m.IsOnline(999) {
		t.Error("user 999 should not be online")
	}

	ids := m.GetUserConnIDs(100)
	if len(ids) != 2 {
		t.Errorf("len(ids) = %d, want 2", len(ids))
	}

	m.Remove(1, 100)
	if !m.IsOnline(100) {
		t.Error("user 100 should still be online")
	}

	m.Remove(2, 100)
	if m.IsOnline(100) {
		t.Error("user 100 should be offline")
	}
}

func TestManager_GetConn(t *testing.T) {
	m := NewManager()
	c1 := newTestConn(1, 100)
	m.Add(c1)

	if got := m.GetConn(1); got != c1 {
		t.Error("GetConn should return the same pointer")
	}
	if m.GetConn(999) != nil {
		t.Error("GetConn should return nil for unknown")
	}
}

func TestManager_SendToUser(t *testing.T) {
	m := NewManager()
	c1 := newTestConn(1, 100)
	m.Add(c1)
	m.SendToUser(100, []byte("hello"))
}

func TestManager_Broadcast(t *testing.T) {
	m := NewManager()
	m.Add(newTestConn(1, 100))
	m.Add(newTestConn(2, 200))
	m.Broadcast([]byte("broadcast"))
}

func TestManager_OnlineUsers(t *testing.T) {
	m := NewManager()
	if n := m.OnlineUsers(); n != 0 {
		t.Errorf("OnlineUsers = %d, want 0", n)
	}
	m.Add(newTestConn(1, 100))
	m.Add(newTestConn(2, 200))
	if n := m.OnlineUsers(); n != 2 {
		t.Errorf("OnlineUsers = %d, want 2", n)
	}
}
