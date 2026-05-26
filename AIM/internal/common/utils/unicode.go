package utils

import "unicode"

func IsEmpty(s string) bool {
	return len(s) == 0
}
func IsNotEmpty(s string) bool {
	return !IsEmpty(s)
}
func HashChinese(s string) bool {
	for _, ch := range s {
		if unicode.Is(unicode.Han, ch) {
			return true
		}
	}
	return false
}
