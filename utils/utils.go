package utils

import (
	"fmt"
	"github.com/mssola/user_agent"
	"golang.org/x/exp/rand"
	"regexp"
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

// GetUserDeviceInfo 获取用户设备信息
func GetUserDeviceInfo(userAgent string) string {
	ua := user_agent.New(userAgent)
	browser, version := ua.Browser()
	osInfo := ua.OSInfo()
	return fmt.Sprintf("%s - %s - %s", browser, version, osInfo.FullName)
}

// CheckStrIsPhone 使用正则表达式校验字符串是否手机
func CheckStrIsPhone(str string) bool {
	// 手机正则表达式
	pattern := `^(13[0-9]|14[5-9]|15[0-3,5-9]|16[6]|17[0-8]|18[0-9]|19[8,9])\d{8}$`
	// 使用正则表达式进行匹配
	return regexp.MustCompile(pattern).MatchString(str)
}
