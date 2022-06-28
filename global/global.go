package global

import (
	"github.com/go-redis/redis"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"platform/config"
)

var (
	Db     *gorm.DB
	Config config.App
	Log    *zap.SugaredLogger
	Vp     *viper.Viper
	Redis  *redis.Client
)
