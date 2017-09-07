package register

import (
	"github.com/tosone/backend-golang/service/register"
	"gopkg.in/kataras/iris.v8"
)

// Index 登陆注册
func Index(app *iris.Application) {
	app.Post("/login", register.Login)
	app.Post("/register", register.Register)
}
