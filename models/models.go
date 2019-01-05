package models

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/astaxie/beego/orm"
	"time"
)

type User struct {
	Id          int
	Name        string
	Password    string
	ArticleType []*ArticleType `orm:"reverse(many)"`
	Artilce     []*Article     `orm:"reverse(many)"`
}

type Article struct {
	Id    int    `orm:"auto;pk"`
	Title string `orm:"size(60)"`
	//Ctime   time.Time `orm:"auto_now_add;type(datetime)"`
	Time        time.Time    `orm:"type(datetime);auto_now"`
	Count       int          `orm:"default(0)"`
	Content     string       `orm:"type(text)"`
	Img         string       `orm:"null"`
	User        []*User      `orm:"rel(m2m)"`
	ArticleType *ArticleType `orm:"rel(fk)"`
}

type ArticleType struct {
	Id       int        `orm:auto;pk`
	TypeName string     `orm:"size(20);unique;null;on_delete(do_nothing)"`
	Articles []*Article `orm:"reverse(many)"`
	User     []*User    `orm:"rel(m2m)"`
}

func init() {
	// [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
	//orm.RegisterDataBase("default","mysql","root:123456@tcp(127.0.0.1:3306)/test?charset=utf8")
	orm.RegisterDataBase("default", "mysql", "root:123456@tcp(127.0.0.1:3306)/testNews?charset=utf8")
	orm.RegisterModel(new(User), new(Article), new(ArticleType))
	orm.RunSyncdb("default", false, true)
}
