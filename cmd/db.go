package cmd

import (
	"gin_chat/global"
	"gin_chat/models"
)

func MakeMigration() {
	var err error

	global.Log.Infof("开始迁移数据库")
	err = global.DB.Set("gorm:table_options", "ENGINE=InnoDB CHARSET=utf8mb4").AutoMigrate(
		&models.UserBasic{},
	)
	if err != nil {
		global.Log.Error("数据库迁移失败: %v", err)
		return
	}
	global.Log.Info("数据库迁移成功")

}
