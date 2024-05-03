package cmd

import (
	"flag"
	"gin_chat/global"
)

type Option struct {
	DB bool
}

// Parse 解析命令行： go run main.go -db
func Parse() Option {
	// go run main.go -db
	db := flag.Bool("db", false, "初始化数据库")
	// 解析命令
	flag.Parse()
	return Option{
		DB: *db,
	}
}

func IsStopWeb(option *Option) bool {
	if option.DB {
		global.Log.Infof("进行数据库迁移不执行gin服务")
		return true
	}
	return false
}

func SwitchOption(option *Option) {
	if option.DB {
		// 迁移数据库
		MakeMigration()
	}
}
