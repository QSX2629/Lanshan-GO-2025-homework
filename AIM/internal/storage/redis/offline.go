package redis

import "encoding/json"

const OfflinePrefix = "im:offline:"

func SaveOfflineMessage(userID string, msg any) {
	key := OfflinePrefix + userID
	data, _ := json.Marshal(msg)
	RDB.LPush(Ctx, key, data)
}

func GetOfflineMessages(userID string) [][]byte {
	key := OfflinePrefix + userID
	list, _ := RDB.LRange(Ctx, key, 0, -1).Result()
	RDB.Del(Ctx, key)

	res := make([][]byte, len(list))
	for i, v := range list {
		res[i] = []byte(v)
	}
	return res
}
