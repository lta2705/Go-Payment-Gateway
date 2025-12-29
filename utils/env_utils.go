package utils

import (
	"os"
	"strconv"
)

func String(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok {
		return val
	}
	return defaultVal
}

func Int(key string, defaultVal int) int {
	valStr := String(key, "")
	if val, err := strconv.Atoi(valStr); err == nil {
		return val
	}
	return defaultVal
}

func Bool(key string, defaultVal bool) bool {
	valStr := String(key, "")
	if val, err := strconv.ParseBool(valStr); err == nil {
		return val
	}
	return defaultVal
}
