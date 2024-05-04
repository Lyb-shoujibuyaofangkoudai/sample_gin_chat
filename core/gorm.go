package core

import (
	"fmt"
	"gin_chat/config"
	"gin_chat/global"
	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"time"
)

func InitGorm() *gorm.DB {
	mysqlSettings := viper.GetStringMapString("mysql")
	if mysqlSettings["host"] == "" {
		global.Log.Warnf("未配置数据库信息")
		return nil
	}
	mySqlLogger := logger.New(
		//log.New(os.Stdout, "\r\n", log.LstdFlags),
		&LogrusGormLogger{global.Log},
		logger.Config{
			SlowThreshold: time.Second, // Slow SQL threshold
			LogLevel:      logger.Info,
			Colorful:      true,
		},
	)
	db, err := gorm.Open(mysql.Open(config.Dsn()), &gorm.Config{
		Logger: mySqlLogger,
	})
	if err != nil {
		global.Log.Fatal(fmt.Sprintf("数据库连接失败: %s", err))
	}

	sqlDB, err := db.DB()
	if err != nil {
		panic(err)
	}

	sqlDB.SetMaxIdleConns(10)               // 设置连接池中的最大闲置连接数
	sqlDB.SetMaxOpenConns(100)              // 设置连接池最大连接数
	sqlDB.SetConnMaxLifetime(time.Hour * 4) // 设置连接池最大生存时间，不能超过mysql的wait_timeout
	global.Log.Infof("数据库连接成功")
	return db
}
