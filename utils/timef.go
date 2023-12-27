package utils

import (
	"time"
)

// 20230101101010
func TimeFormat() string {
	layout := "20060102150405"
	return time.Now().Format(layout)
}
