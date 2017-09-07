package router

import (
	"github.com/tosone/backend-golang/router/register"
	"gopkg.in/kataras/iris.v8"
)

// Index 入口文件
func Index(app *iris.Application) {
	register.Index(app)
}
