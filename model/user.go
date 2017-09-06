package model

import "gopkg.in/mgo.v2/bson"

// UserRegisterForm 用户登录注册的信息
type UserRegisterForm struct {
	ID   bson.ObjectId `json:"id" bson:"_id"`
	Name string        `json:"name" bson:"name"`
	Salt string        `json:"salt" bson:"salt"`
	Hash string        `json:"hash" bson:"hash"`
}
