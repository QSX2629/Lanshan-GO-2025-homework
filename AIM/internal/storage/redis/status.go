package redis

import "time"

const OnlineKey = "im:user:online"

func GetUserCometID(userID string) string {
	return RDB.HGet(Ctx, OnlineKey, userID).Val()
}

func IsUserOnline(userID string) bool {
	ok, _ := RDB.HExists(Ctx, OnlineKey, userID).Result()
	return ok
}
func SetOnline(uid string) error {
	return RDB.Set(Ctx, "online:"+uid, "1", 30*time.Minute).Err()
}

func SetOffline(uid string) error {
	return RDB.Del(Ctx, "online:"+uid).Err()
}
