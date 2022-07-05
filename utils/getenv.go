package utils

import "os"

func GetenvOr(key, default_val string) string {
	value := os.Getenv(key)
	if value == "" {
		value = default_val
	}
	return value
}
