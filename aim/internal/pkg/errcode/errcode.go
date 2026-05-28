// Package errcode defines unified error codes for the AIM system.
package errcode

// Code is a machine-readable error code.
type Code int

const (
	// 2xx: Success
	OK Code = 200 // Success

	// 4xx: Client errors
	InvalidParam   Code = 400 // Invalid request parameter
	Unauthorized   Code = 401 // Authentication required
	Forbidden      Code = 403 // Permission denied
	NotFound       Code = 404 // Resource not found
	Conflict       Code = 409 // Resource already exists
	TooManyRequest Code = 429 // Rate limit exceeded

	// 4xxx specific
	UserNotFound  Code = 4001 // User not found
	UserExists    Code = 4002 // User already exists
	WrongPassword Code = 4003 // Wrong password
	TokenExpired  Code = 4004 // Token expired
	TokenInvalid  Code = 4005 // Token invalid

	MsgNotFound      Code = 4101 // Message not found
	MsgSendFailed    Code = 4102 // Message send failed
	MsgTypeUnsupport Code = 4103 // Message type unsupported

	NotFriend      Code = 4201 // Not a friend
	AlreadyFriend  Code = 4202 // Already a friend
	GroupNotFound  Code = 4203 // Group not found
	NotGroupMember Code = 4204 // Not a group member
	GroupBanned    Code = 4205 // User is banned in group

	AIProviderError Code = 4301 // AI provider error
	AIQuotaExceeded Code = 4302 // AI quota exceeded
	AIKeyInvalid    Code = 4303 // User API key invalid

	FileTooLarge    Code = 4401 // File too large
	FileTypeInvalid Code = 4402 // File type not allowed
	UploadFailed    Code = 4403 // Upload failed

	// 5xx: Server errors
	InternalError Code = 500 // Internal server error
	DBError       Code = 501 // Database error
	RPCError      Code = 502 // RPC call error
)

// Message returns the default human-readable message for the code.
func (c Code) Message() string {
	switch c {
	case OK:
		return "ok"
	case InvalidParam:
		return "invalid parameter"
	case Unauthorized:
		return "authentication required"
	case Forbidden:
		return "permission denied"
	case NotFound:
		return "resource not found"
	case Conflict:
		return "resource already exists"
	case TooManyRequest:
		return "too many requests"
	case UserNotFound:
		return "user not found"
	case UserExists:
		return "user already exists"
	case WrongPassword:
		return "wrong password"
	case TokenExpired:
		return "token expired"
	case TokenInvalid:
		return "token invalid"
	case MsgNotFound:
		return "message not found"
	case MsgSendFailed:
		return "message send failed"
	case MsgTypeUnsupport:
		return "message type unsupported"
	case NotFriend:
		return "not a friend"
	case AlreadyFriend:
		return "already a friend"
	case GroupNotFound:
		return "group not found"
	case NotGroupMember:
		return "not a group member"
	case GroupBanned:
		return "user is banned in group"
	case AIProviderError:
		return "AI provider error"
	case AIQuotaExceeded:
		return "AI quota exceeded"
	case AIKeyInvalid:
		return "user API key invalid"
	case FileTooLarge:
		return "file too large"
	case FileTypeInvalid:
		return "file type not allowed"
	case UploadFailed:
		return "upload failed"
	case InternalError:
		return "internal server error"
	case DBError:
		return "database error"
	case RPCError:
		return "RPC call error"
	default:
		return "unknown error"
	}
}
