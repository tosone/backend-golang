package register

import (
	"log"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	"github.com/satori/go.uuid"
	"github.com/tosone/backend-golang/config"
	"github.com/tosone/backend-golang/model"
	"github.com/tosone/backend-golang/mongo"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/kataras/iris.v8"
	"gopkg.in/mgo.v2/bson"
)

type registerForm struct {
	Name     string `form:"name"`
	Password string `form:"password"`
}

// Login 注册
func Login(ctx *iris.Context) {
	var loginForm registerForm
	ctx.ReadJSON(&loginForm)
	DB := mongo.MgoDb{}
	DB.Init()
	defer DB.Close()
	var userInfo model.UserRegisterForm
	if err := DB.C("test").Find(bson.M{"name": loginForm.Name}).One(&userInfo); err != nil {
		log.Println(err)
		ctx.JSON(405, model.ResponseInfo{Status: 405, Info: "MongoDB is error, please retry again."})
		return
	}
	if err := checkPasswordHash(loginForm.Password+userInfo.Salt, userInfo.Hash); err != nil {
		log.Println(err)
		ctx.JSON(405, model.ResponseInfo{Status: 405, Info: "Password or Username is wrong."})
		return
	}
	mongo.RedisPool.Get().Do("Set", uuid.NewV4().String(), "EX", config.SessionExpire)
	token := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"sid": uuid.NewV4().String(),
		"exp": time.Now().Unix() + config.SessionExpire,
	})
	signedString, err := token.SignedString(config.SessionSecret)
	if err != nil {
		log.Println(err)
	}
	ctx.JSON(iris.StatusOK, map[string]string{"authenticate": signedString})
}

// Register 注册
func Register(ctx *iris.Context) {
	var data registerForm
	ctx.ReadJSON(&data)
	var hash string
	var err error
	salt := uuid.NewV4().String()
	if hash, err = hashPassword(data.Password + salt); err != nil {
		log.Println(err)
	}
	userInfo := model.UserRegisterForm{
		ID:   bson.NewObjectId(),
		Name: data.Name,
		Hash: hash,
		Salt: salt,
	}
	DB := mongo.MgoDb{}
	DB.Init()
	DB.C("test").Insert(userInfo)
	defer DB.Close()
	ctx.JSON(iris.StatusOK, userInfo)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err
}
