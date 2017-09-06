package main

import (
	"fmt"

	"github.com/iris-contrib/middleware/cors"
	"github.com/tosone/backend-golang/config"
	"github.com/tosone/backend-golang/router"
	"gopkg.in/kataras/iris.v8"
)

var (
	// BuildStamp BuildStamp
	BuildStamp = "Nothing Provided."
	// GitHash GitHash
	GitHash = "Nothing Provided."
)

func main() {
	fmt.Printf("Git Commit Hash: %s\n", GitHash)
	fmt.Printf("UTC Build Time: %s\n", BuildStamp)
	app := iris.New()
	app.Adapt(
		iris.DevLogger(),
		cors.New(cors.Options{
			AllowedOrigins:   []string{"*"},
			AllowCredentials: true,
		}))
	router.Index(app)

	app.Listen(config.IP + ":" + config.PORT)
}
