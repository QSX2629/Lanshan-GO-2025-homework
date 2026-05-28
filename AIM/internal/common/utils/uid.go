package utils

import (
	"strings"

	"github.com/google/uuid"
)

// 生成id
func UUID() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")
}

// 生成消息id
func MsgID() string {
	return "msg_" + UUID()
}

// 生成用户id
func UserID() string {
	return "user_" + UUID()
}

// 生成群id
func GroupID() string {
	return "group_" + UUID()
}
