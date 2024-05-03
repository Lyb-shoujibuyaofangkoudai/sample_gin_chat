package utils

import (
	"fmt"
	"github.com/mssola/user_agent"
	"golang.org/x/exp/rand"
	"time"
)

// GenerateSalt 创建一个随机盐值
func GenerateSalt(length int) string {
	const alphanumeric = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	bytes := make([]byte, length)
	rand.Seed(uint64(time.Now().UnixNano()))
	for i := range bytes {
		bytes[i] = alphanumeric[rand.Intn(len(alphanumeric))]
	}
	return string(bytes)
}

func GetUserDeviceInfo(userAgent string) string {
	ua := user_agent.New(userAgent)
	browser, version := ua.Browser()
	osInfo := ua.OSInfo()
	return fmt.Sprintf("%s - %s - %s", browser, version, osInfo.FullName)
}
