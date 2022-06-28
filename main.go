package main

import (
	"github.com/kataras/iris/v12"
	"platform/core"
	"platform/global"
	"platform/initialize"
)

// @title 权限控制系统 API
// @version 1.0
// @description 权限控制系统API文档
// @BasePath /
func main() {
	global.Vp = core.Viper()
	global.Log = core.Zap()
	global.Db = initialize.Gorm()
	if global.Db != nil {
		initialize.MysqlTables(global.Db) // 初始化表
		// 程序结束前关闭数据库链接
		db, _ := global.Db.DB()
		defer db.Close()
	}

	app := iris.Default()
	if err := core.Run(app); err != nil {
		global.Log.Fatal(err)
	}
}
