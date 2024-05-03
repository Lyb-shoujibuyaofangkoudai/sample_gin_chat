package models

import (
	"gin_chat/global"
	"gorm.io/gorm"
	"time"
)

// UserBasic 用户基础信息表
type UserBasic struct {
	gorm.Model
	Name          string    `gorm:"size:32" json:"name"`                                                  //用户名
	Password      string    `gorm:"size:128" json:"password"`                                             //密码
	Phone         string    `gorm:"size:20;uniqueIndex"  valid:"matches(^1[3-9]{1}\\d{9}$)" json:"phone"` //手机号
	Email         string    `gorm:"size:128;uniqueIndex" valid:"email" json:"email"`                      //邮箱
	Avatar        string    `gorm:"size:128" json:"avatar"`                                               //头像
	Identity      string    `json:"identity"`                                                             // 用户的唯一标识
	ClientIp      string    `json:"client_ip"`                                                            // 登陆ip
	ClientPort    string    `json:"client_port"`                                                          // 端口
	Salt          string    `json:"salt"`                                                                 // 密码盐
	LoginTime     time.Time `json:"login_time"`                                                           // 登陆时间
	HeartbeatTime time.Time `json:"heartbeat_time"`                                                       // 心跳时间
	LoginOutTime  time.Time `gorm:"column:login_out_time" json:"login_out_time"`                          // 退出登录时间
	IsLogout      bool      `json:"is_logout"`                                                            // 是否退出
	DeviceInfo    string    `json:"device_info"`                                                          // 用户设备信息
}

func (*UserBasic) TableName() string {
	return "user_basic"
}

// CreateUser 创建用户
func (user *UserBasic) CreateUser() error {
	return global.DB.Create(user).Error
}

func (user *UserBasic) FindUserByStruct(selectSlice []string) error {
	if len(selectSlice) != 0 {
		return global.DB.Select(selectSlice).Where(user).First(user).Error
	}
	return global.DB.Where(user).First(user).Error
}

type UserInfo struct {
	gorm.Model
	UserID    uint   `json:"user_id"`
	Signature string `gorm:"size:512" json:"signature"`
	Gender    int    `gorm:"default:1,comment:1男2女" json:"gender"`
}
