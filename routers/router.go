package routers

import (
	"testNews/controllers"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func init() {
	beego.InsertFilter("/article/*", beego.BeforeExec, filterFunc)
	beego.Router("/", &controllers.MainController{})
	beego.Router("/register", &controllers.UserController{}, "get:ShowRegister;post:HandleRegister")
	beego.Router("/login", &controllers.UserController{}, "get:ShowLogin;post:HandleLogin")
	beego.Router("/article/quit", &controllers.UserController{}, "get:HandleQuit")

	beego.Router("/article/index", &controllers.ArticleController{}, "get:ShowIndex")
	beego.Router("/article/add", &controllers.ArticleController{}, "get:ShowAdd;post:HandleAdd")
	beego.Router("/article/content", &controllers.ArticleController{}, "get:ShowContent")
	beego.Router("/article/update", &controllers.ArticleController{}, "get:ShowUpdate;post:HandleUpdate")
	beego.Router("/article/delete", &controllers.ArticleController{}, "get:HandleDelete")
	beego.Router("/article/addType", &controllers.ArticleController{}, "get:ShowAddType;post:HandleShowType")
	beego.Router("/article/deleteType", &controllers.ArticleController{}, "get:DeleteType")

	beego.Router("/redis", &controllers.GoRedis{}, "get:ShowGet")
}

var filterFunc = func(ctx *context.Context) {
	userName := ctx.Input.Session("userName")
	if userName == nil {
		ctx.Redirect(302, "/login")
		return
	}
}
