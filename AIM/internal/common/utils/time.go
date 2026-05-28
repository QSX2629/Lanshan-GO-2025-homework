package utils

import "time"

// 格式化时间
func FormatTime(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}
func NowTime() string {
	return FormatTime(time.Now())
}
func NowUnix() int64 {
	return time.Now().Unix()
}
