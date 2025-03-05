package utils

import "os"

func GetEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		return ""
	}
	return value
}
