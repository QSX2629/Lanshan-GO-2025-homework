// Package session manages the mapping between users and their WebSocket connections.
package session

import (
	"sync"

	"github.com/aim/aim/internal/gateway/ws"
)

// Manager tracks user sessions and their active connections.
// A single user may have multiple connections (different devices).
type Manager struct {
	mu    sync.RWMutex
	conns map[uint]*ws.Conn      // connID -> conn
	user  map[uint]map[uint]bool // userID -> set of connIDs
}

// NewManager creates a new session Manager.
func NewManager() *Manager {
	return &Manager{
		conns: make(map[uint]*ws.Conn),
		user:  make(map[uint]map[uint]bool),
	}
}

// Add registers a new connection for a user.
func (m *Manager) Add(conn *ws.Conn) {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.conns[conn.ID] = conn

	if m.user[conn.UserID] == nil {
		m.user[conn.UserID] = make(map[uint]bool)
	}
	m.user[conn.UserID][conn.ID] = true
}

// Remove unregisters a connection.
func (m *Manager) Remove(connID, userID uint) {
	m.mu.Lock()
	defer m.mu.Unlock()

	delete(m.conns, connID)

	if m.user[userID] != nil {
		delete(m.user[userID], connID)
		if len(m.user[userID]) == 0 {
			delete(m.user, userID)
		}
	}
}

// GetConn returns a connection by ID.
func (m *Manager) GetConn(connID uint) *ws.Conn {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.conns[connID]
}

// GetUserConnIDs returns all connection IDs for a user.
func (m *Manager) GetUserConnIDs(userID uint) []uint {
	m.mu.RLock()
	defer m.mu.RUnlock()

	connSet := m.user[userID]
	ids := make([]uint, 0, len(connSet))
	for id := range connSet {
		ids = append(ids, id)
	}
	return ids
}

// IsOnline reports whether a user has any active connections.
func (m *Manager) IsOnline(userID uint) bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.user[userID]) > 0
}

// SendToUser sends a message to all connections of a user.
func (m *Manager) SendToUser(userID uint, data []byte) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for connID := range m.user[userID] {
		if conn, ok := m.conns[connID]; ok {
			conn.Send(data)
		}
	}
}

// Broadcast sends a message to all connected users.
func (m *Manager) Broadcast(data []byte) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for _, conn := range m.conns {
		conn.Send(data)
	}
}

// OnlineUsers returns the count of online users.
func (m *Manager) OnlineUsers() int {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return len(m.user)
}
