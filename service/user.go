package service

import (
	"context"
	"fmt"
	"gin_chat/global"
	"gin_chat/models"
	"gin_chat/utils"
	"gin_chat/utils/encrypt"
	"gin_chat/utils/res"
	"gin_chat/utils/validators"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"golang.org/x/exp/rand"
	"strconv"
	"time"
)

// GetIndexView
// @Summary      获取首页数据
// @Description  返回首页json公用数据
// @Tags         首页数据
// @Produce      json
// @Router       /api/index [get]
// @success       200 {object}  string "{"code":200,"data":{"message":"pong"},"msg":"ok"}"
func GetIndexView(c *gin.Context) {
	c.JSON(200, gin.H{
		"message": "pong",
	})
}

type UserRegister struct {
	Name       string `form:"name" binding:"required"`
	Password   string `form:"password" binding:"required"`
	RePassword string `form:"re_password" binding:"required,eqfield=Password" msg:"两次密码输入不同"`
	Phone      string `form:"phone" binding:"required,phone" msg:"手机号码格式不正确"`
	Email      string `form:"email" binding:"required,email" msg:"邮箱格式不正确"`
	Avatar     string `form:"avatar"`
}

// Register
// @Summary      创建用户
// @Description  注册新用户
// @Tags         用户
// @Accept       json
// @Produce      json
// @Param        name  body      string  true  "用户名"
// @Param        password  body      string  true  "密码"
// @Param        re_password  body      string  true  "确认密码"
// @Param        phone  body      string  true  "手机号码"
// @Param        email  body      string  true  "邮箱"
// @Param        avatar  body      string  true  "头像"
// @Success      200  {object}  res.Response
// @Router       /api/register [post]
func Register(c *gin.Context) {
	// 获取用户设备信息
	userAgent := c.Request.UserAgent()
	if userAgent == "" {
		res.FailWithMsg("非法请求", c)
		return
	}
	userInfo := UserRegister{}
	err := c.ShouldBind(&userInfo)
	if err != nil {
		res.FailWithMsg(validators.GetValidMsg(err, &userInfo), c)
		return
	}

	// 校验用户是否已经注册
	u := FindUserByPhone(userInfo.Phone)
	if u.ID != 0 {
		res.FailWithMsg("用户已注册", c)
		return
	}

	// 使用加密安全的随机数生成器
	src := rand.NewSource(uint64(time.Now().UnixNano()))
	r := rand.New(src)
	// 生成1到62之间的随机整数
	seed := r.Intn(62) + 1
	// 生成盐值
	salt := utils.GenerateSalt(seed)
	// 密码加密
	userInfo.Password = encrypt.EncryptPassword(userInfo.Password, salt)

	// 生成用户唯一标识
	identity := encrypt.MD5(fmt.Sprintf("%d", strconv.FormatInt(time.Now().UnixNano(), 10)+userInfo.Phone))

	deviceInfo := utils.GetUserDeviceInfo(userAgent)

	user := models.UserBasic{
		Name:          userInfo.Name,
		Password:      userInfo.Password,
		Phone:         userInfo.Phone,
		Email:         userInfo.Email,
		Avatar:        userInfo.Avatar,
		DeviceInfo:    deviceInfo,
		Salt:          salt,
		Identity:      identity,
		LoginTime:     time.Now(),
		HeartbeatTime: time.Now(),
		LoginOutTime:  time.Now(),
	}
	err = user.CreateUser()
	if err != nil {
		global.Log.Errorf("注册失败：%v", err.Error())
		res.FailWithMsg("注册失败", c)
		return
	}
	res.OkWithMsg("注册成功", c)
}

func FindUserByPhone(phone string) (user models.UserBasic) {
	user = models.UserBasic{}
	if global.DB.Where("phone = ?", phone).First(&user).Error != nil {
		return
	}
	return user
}

type UserLogin struct {
	Account      string `form:"account" binding:"required" msg:"账号不能为空，账号为手机号码或邮箱"`
	Password     string `form:"password" binding:"required" msg:"密码不能为空"`
	SignInMethod string `form:"sign_in_method" binding:"required,oneof=email phone password" msg:"登录方式不能为空"`
	Code         string `form:"code" binding:"loginCode" msg:"验证码不能为空"`
}

// Login
// @Summary 用户登录
// @Description 用户登录并获取令牌
// @Tags 用户
// @Accept json
// @Produce json
// @Param email body LoginRequest true "邮箱"
// @Success 200 {object} LoginResponse "成功响应"
// @Failure 400 {object} Error "请求错误"
// @Router /users/login [post]

func Login(c *gin.Context) {
	loginInfo := UserLogin{}
	err := c.ShouldBind(&loginInfo)
	if err != nil {
		res.FailWithMsg(validators.GetValidMsg(err, &loginInfo), c)
		return
	}
	user := models.UserBasic{
		Email: loginInfo.Account,
	}
	selectSlice := []string{
		"name",
		"password",
		"phone",
		"email",
		"avatar",
		"identity",
		"device_info",
	}
	err = user.FindUserByStruct(selectSlice)
	if err != nil {
		if err.Error() == "record not found" {
			user.Email = ""
			user.Phone = loginInfo.Account
			err = user.FindUserByStruct(selectSlice)
			if err != nil {
				res.FailWithMsg("该用户不存在", c)
				return
			}
		}
	}
	rdbUserTokenKey := fmt.Sprintf("user_%s", user.Identity)
	oToken := global.RDB.Get(context.Background(), rdbUserTokenKey).Val()
	if oToken != "" {
		myClaims, err := utils.ParseTokenRs256(oToken)
		global.Log.Infof("解析token: %v", myClaims)
		if err != nil {
			res.FailWithMsg("token已过期", c)
			return
		}
		expirationTime, err := myClaims.GetExpirationTime()
		global.Log.Infof("过期时间: %v", expirationTime)
		if err == nil {
			res.OkWithData(user, c)
			return
		}
		global.Log.Errorf("解析token失败: %v", err)
	}

	var token string
	token, err = utils.GenerateTokenUsingRS256(int(user.ID), user.Name)
	if err != nil {
		res.FailWithMsg(err.Error(), c)
		return
	}
	global.RDB.Set(context.Background(), rdbUserTokenKey, token, time.Duration(viper.GetInt64("jwt.expire_time"))*time.Hour)
	res.OkWithData(user, c)
}
