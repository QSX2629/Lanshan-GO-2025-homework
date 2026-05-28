package ws

import (
	"testing"
	"time"
)

func TestNewConn(t *testing.T) {
	// Test struct creation without a real websocket.
	c := &Conn{
		ID:     1,
		UserID: 100,
		send:   make(chan []byte, 256),
	}
	if c.ID != 1 {
		t.Errorf("ID = %d, want 1", c.ID)
	}
	if c.UserID != 100 {
		t.Errorf("UserID = %d, want 100", c.UserID)
	}
}

func TestConnSendJSON(t *testing.T) {
	// Test that SendJSON produces valid JSON.
	c := &Conn{
		ID:     1,
		UserID: 100,
		send:   make(chan []byte, 256),
		closed: false,
	}

	type testMsg struct {
		Text string `json:"text"`
	}

	// SendJSON will try to write to send channel.
	// Since nobody reads from send in test, we consume it.
	go func() {
		<-c.send // Consume the message so channel doesn't block
	}()

	if err := c.SendJSON(testMsg{Text: "hello"}); err != nil {
		t.Fatalf("SendJSON() error = %v", err)
	}
	time.Sleep(10 * time.Millisecond)
}

func TestConnSendClosed(t *testing.T) {
	c := &Conn{
		ID:     2,
		UserID: 200,
		send:   make(chan []byte, 256),
		closed: true,
	}
	// Send to closed conn should not panic.
	c.Send([]byte("test"))
}

func TestConnClose(t *testing.T) {
	c := &Conn{
		ID:     3,
		UserID: 300,
		send:   make(chan []byte, 1),
	}
	c.Close()
	// Double close should not panic.
	c.Close()
}

func TestRawMessage(t *testing.T) {
	rm := RawMessage{ConnID: 1, UserID: 100, Data: []byte("test")}
	if rm.ConnID != 1 || rm.UserID != 100 {
		t.Error("RawMessage fields mismatch")
	}
}
