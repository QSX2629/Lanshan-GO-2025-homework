package push

import (
	"AIM/internal/comet/connection"
	"AIM/internal/comet/protocol"
	"encoding/json"
)

// PushToUser 给单个用户推送
func PushToUser(uid string, msg *protocol.Message) {
	conn := connection.GetConn(uid)
	if conn == nil {
		return
	}
	data, _ := json.Marshal(msg)
	conn.SendMessage(data)
}

// PushToGroup 给一组用户推送（只负责推送，不查数据库）
func PushToGroup(uids []string, msg *protocol.Message) {
	for _, uid := range uids {
		PushToUser(uid, msg)
	}
}
