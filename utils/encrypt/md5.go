package encrypt

import (
	"crypto/md5"
	"fmt"
	"strings"
)

func Md5(str string) string {
	hasher := md5.New()
	hasher.Write([]byte(str))
	return fmt.Sprintf("%x", hasher.Sum(nil))
}
func MD5(str string) string {
	return strings.ToUpper(Md5(str))
}

// EncryptPassword 使用MD5和盐值加密密码
func EncryptPassword(password string, salt string) (encryptPassword string) {
	return MD5(password + salt)
}

// ValidPassword 检验密码
func ValidPassword(password, salt string, sqlPassword string) bool {
	return MD5(password+salt) == sqlPassword
}
