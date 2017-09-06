package mongo

import (
	"log"

	"github.com/tosone/backend-golang/config"
	"gopkg.in/mgo.v2"
)

// MSession mongo Session
var MSession *mgo.Session

// MgoDb 链接数据库的对象
type MgoDb struct {
	Session *mgo.Session
	Db      *mgo.Database
	Col     *mgo.Collection
}

func init() {
	if MSession == nil {
		var err error
		MSession, err = mgo.Dial(config.MongoURL)
		if err != nil {
			log.Println(err)
		}
		MSession.SetMode(mgo.Monotonic, true)
	}
}

// Init 初始化
func (mdb *MgoDb) Init() *mgo.Session {
	mdb.Session = MSession.Copy()
	mdb.Db = mdb.Session.DB(config.MongoDatabase)
	return mdb.Session
}

// C 选择集合
func (mdb *MgoDb) C(collection string) *mgo.Collection {
	return mdb.Db.C(collection)
}

// Close 关闭数据库
func (mdb *MgoDb) Close() bool {
	defer mdb.Session.Close()
	return true
}

// DropoDb 删除数据库
func (mdb *MgoDb) DropoDb() error {
	err := mdb.Db.DropDatabase()
	if err != nil {
		return err
	}
	return nil
}

// RemoveAll 移除所有的集合
func (mdb *MgoDb) RemoveAll(collection string) bool {
	mdb.Db.C(collection).RemoveAll(nil)
	mdb.Col = mdb.Db.C(collection)
	return true
}

// Index 返回所有的序号
func (mdb *MgoDb) Index(collection string, keys []string) bool {
	index := mgo.Index{
		Key:        keys,
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}
	err := mdb.Db.C(collection).EnsureIndex(index)
	if err != nil {
		log.Println(err)
		return false
	}
	return true
}

// IsDup 检查是否连接
func (mdb *MgoDb) IsDup(err error) bool {
	if mgo.IsDup(err) {
		return true
	}
	return false
}
