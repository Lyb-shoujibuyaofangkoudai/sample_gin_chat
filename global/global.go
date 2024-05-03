package global

import (
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

var (
	DB  *gorm.DB
	RDB *redis.Client
	Log *logrus.Logger
)
