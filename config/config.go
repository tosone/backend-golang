package config

import (
	"fmt"

	"github.com/caarlos0/env"
)

var (
	// IP 地址
	IP = "127.0.0.1"

	// PORT 端口
	PORT = "8080"

	// MongoURL mongo 连接字符串
	MongoURL = ""

	// MongoDatabase 链接 mongo 的数据库
	MongoDatabase = "tosone"

	// PriKeyPath 私钥
	PriKeyPath = "keys/app.rsa"

	// PubKeyPath 公钥
	PubKeyPath = "keys/app.rsa.pub"

	// PasswordSalt 用户密码的盐值
	PasswordSalt = "tosone"

	// SessionExpire 最长的 session 过期时间
	SessionExpire int64 = 24 * 60 * 60

	// SessionSecret Session Secret
	SessionSecret = "9787581d51ca21f452512ce58d98ceb4"
)

// config 配置
type config struct {
	MongoHost string `env:"MongoHost" envDefault:"localhost"`
	MongoPort string `env:"MongoPort" envDefault:"3000"`
	MongoDB   string `env:"MongoDB"`
	MongoUser string `env:"MongoUser"`
	MongoPass string `env:"MongoPass"`
}

func init() {
	var cfg config
	err := env.Parse(&cfg)
	if err != nil {
		fmt.Printf("%+v\n", err)
	}
	MongoURL = fmt.Sprintf("mongodb://%s:%s@%s:%s/%s", cfg.MongoUser, cfg.MongoPass, cfg.MongoHost, cfg.MongoPort, "admin")
}
