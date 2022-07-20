package utils

import "os"

func GetenvOr(key, defaultVal string) string {
	value := os.Getenv(key)
	if value == "" {
		value = defaultVal
	}
	return value
}
