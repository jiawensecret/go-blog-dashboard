package initialize

import (
	"github.com/iris-contrib/swagger/v12"
	"github.com/iris-contrib/swagger/v12/swaggerFiles"
	"github.com/kataras/iris/v12"
)

func Router(app *iris.Application) {
	app.Get("/swagger/*", swagger.WrapHandler(swaggerFiles.Handler))
}
