package protocol

import (
	"encoding/json"
	"testing"
)

func TestEncodeDecode(t *testing.T) {
	original := Frame{
		Cmd:     CmdChatMsg,
		Seq:     1,
		From:    "user1",
		To:      "user2",
		MsgType: TypeText,
		Content: `{"text":"hello"}`,
	}

	data, err := Encode(original)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}

	decoded, err := Decode(data)
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	if decoded.Cmd != original.Cmd {
		t.Errorf("Cmd = %q, want %q", decoded.Cmd, original.Cmd)
	}
	if decoded.Seq != original.Seq {
		t.Errorf("Seq = %d, want %d", decoded.Seq, original.Seq)
	}
	if decoded.From != original.From {
		t.Errorf("From = %q, want %q", decoded.From, original.From)
	}
	if decoded.Content != original.Content {
		t.Errorf("Content = %q, want %q", decoded.Content, original.Content)
	}
}

func TestDecodeInvalid(t *testing.T) {
	_, err := Decode([]byte("not json"))
	if err == nil {
		t.Error("Decode() expected error for invalid JSON")
	}
}

func TestValidMsgType(t *testing.T) {
	tests := []struct {
		typ  MsgType
		want bool
	}{
		{TypeText, true},
		{TypeImage, true},
		{TypeFile, true},
		{TypeVoice, true},
		{TypeEvent, true},
		{MsgType("video"), false},
		{MsgType(""), false},
	}

	for _, tt := range tests {
		if got := ValidMsgType(tt.typ); got != tt.want {
			t.Errorf("ValidMsgType(%q) = %v, want %v", tt.typ, got, tt.want)
		}
	}
}

func TestFrameWithExtra(t *testing.T) {
	extra := json.RawMessage(`{"key":"value"}`)
	f := Frame{
		Cmd:   CmdTyping,
		From:  "user1",
		To:    "group1",
		Extra: extra,
	}

	data, err := Encode(f)
	if err != nil {
		t.Fatalf("Encode() error = %v", err)
	}

	decoded, err := Decode(data)
	if err != nil {
		t.Fatalf("Decode() error = %v", err)
	}

	if string(decoded.Extra) != string(extra) {
		t.Errorf("Extra = %s, want %s", decoded.Extra, extra)
	}
}
