package initialize

import (
	"go.uber.org/zap"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"platform/global"
)

func Gorm() *gorm.DB {
	return GormMysql()
}

func MysqlTables(db *gorm.DB) {
	err := db.AutoMigrate()
	if err != nil {
		global.Log.Error("初始化数据表失败", zap.Any("err", err))
		os.Exit(0)
	}
	global.Log.Info("初始化数据表成功")
}

func GormMysql() *gorm.DB {
	m := global.Config.Mysql
	if m.Dbname == "" {
		return nil
	}
	dsn := m.Username + ":" + m.Password + "@tcp(" + m.Host + ":" + m.Port + ")/" + m.Dbname + "?" + m.Config
	mysqlConfig := mysql.Config{
		DSN:                       dsn,
		DefaultStringSize:         191,
		DisableDatetimePrecision:  true,
		DontSupportRenameIndex:    true,
		DontSupportRenameColumn:   true,
		SkipInitializeWithVersion: false,
	}
	if db, err := gorm.Open(mysql.New(mysqlConfig), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	}); err != nil {
		global.Log.Info("连接数据库失败")
		return nil
	} else {
		global.Log.Info("连接数据库成功")
		return db
	}
}
