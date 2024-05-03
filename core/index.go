package core

import (
	"gin_chat/cmd"
	"gin_chat/config"
	"gin_chat/global"
	"gin_chat/router"
	"gin_chat/utils/res"
	"gin_chat/utils/validators"
)

func Init() {
	InitConfig()              // 初始化配置
	global.Log = InitLogger() // 初始化日志
	res.ReadErrorCodeJson()   // 初始化错误码
	global.DB = InitGorm()    // 初始化MySQL
	global.RDB = InitRedis()  // 初始化Redis

	// 终端执行 go main.go -db时就会触发下面的方法
	option := cmd.Parse()
	if cmd.IsStopWeb(&option) {
		cmd.SwitchOption(&option)
		return
	}

	gin := router.InitRouter()

	// 注册自定义校验器
	validators.RegisterPhoneValidators()
	validators.LoginCodeValidate()

	addr := config.Addr()
	global.Log.Infof("服务监听地址为: %v", addr)
	err := gin.Run(addr)
	if err != nil {
		global.Log.Errorf("服务启动失败, err: %v", err)
		return
	}
}
