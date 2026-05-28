package errcode

import "testing"

func TestCode_Message(t *testing.T) {
	tests := []struct {
		code Code
		want string
	}{
		{OK, "ok"},
		{InvalidParam, "invalid parameter"},
		{Unauthorized, "authentication required"},
		{Forbidden, "permission denied"},
		{NotFound, "resource not found"},
		{Conflict, "resource already exists"},
		{TooManyRequest, "too many requests"},
		{UserNotFound, "user not found"},
		{UserExists, "user already exists"},
		{WrongPassword, "wrong password"},
		{TokenExpired, "token expired"},
		{TokenInvalid, "token invalid"},
		{MsgNotFound, "message not found"},
		{MsgSendFailed, "message send failed"},
		{MsgTypeUnsupport, "message type unsupported"},
		{NotFriend, "not a friend"},
		{AlreadyFriend, "already a friend"},
		{GroupNotFound, "group not found"},
		{NotGroupMember, "not a group member"},
		{GroupBanned, "user is banned in group"},
		{AIProviderError, "AI provider error"},
		{AIQuotaExceeded, "AI quota exceeded"},
		{AIKeyInvalid, "user API key invalid"},
		{FileTooLarge, "file too large"},
		{FileTypeInvalid, "file type not allowed"},
		{UploadFailed, "upload failed"},
		{InternalError, "internal server error"},
		{DBError, "database error"},
		{RPCError, "RPC call error"},
	}

	for _, tt := range tests {
		if got := tt.code.Message(); got != tt.want {
			t.Errorf("Code(%d).Message() = %q, want %q", tt.code, got, tt.want)
		}
	}
}

func TestCode_Unknown(t *testing.T) {
	if got := Code(9999).Message(); got != "unknown error" {
		t.Errorf("unknown code = %q, want \"unknown error\"", got)
	}
}

func TestCodeValues(t *testing.T) {
	// Verify no duplicate values
	seen := make(map[Code]bool)
	var duplicates []Code
	for _, c := range allCodes() {
		if seen[c] {
			duplicates = append(duplicates, c)
		}
		seen[c] = true
	}
	if len(duplicates) > 0 {
		t.Errorf("duplicate error codes found: %v", duplicates)
	}
}

func allCodes() []Code {
	return []Code{
		OK, InvalidParam, Unauthorized, Forbidden, NotFound, Conflict, TooManyRequest,
		UserNotFound, UserExists, WrongPassword, TokenExpired, TokenInvalid,
		MsgNotFound, MsgSendFailed, MsgTypeUnsupport,
		NotFriend, AlreadyFriend, GroupNotFound, NotGroupMember, GroupBanned,
		AIProviderError, AIQuotaExceeded, AIKeyInvalid,
		FileTooLarge, FileTypeInvalid, UploadFailed,
		InternalError, DBError, RPCError,
	}
}
