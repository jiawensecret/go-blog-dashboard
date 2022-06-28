package core

import (
	"context"
	"fmt"
	"github.com/kataras/iris/v12"
	"os"
	"os/signal"
	"platform/global"
	"platform/initialize"
	"syscall"
	"time"
)

func Run(app *iris.Application) error {
	// 路由映射处理
	initialize.Router(app)

	go func() {
		global.Log.Infof("服务监听端口: %s", global.Config.System.Port)
		if err := app.Listen(fmt.Sprintf(":%s", global.Config.System.Port)); err != nil {
			global.Log.Fatalf("服务启动失败: ", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGINT)

	<-quit

	ctx, shutdown := context.WithTimeout(context.Background(), 5*time.Second)
	defer shutdown()

	global.Log.Info("服务已正确退出")
	return app.Shutdown(ctx)
}
