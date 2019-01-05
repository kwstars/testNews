package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"testNews/models"
)

type UserController struct {
	beego.Controller
}

func (this *UserController) ShowRegister() {
	this.TplName = "register.html"
}

func (this *UserController) HandleRegister() {
	userName := this.GetString("userName")
	userPassword := this.GetString("password")

	if userName == "" || userPassword == "" {
		this.Data["errMsg"] = "Register failed, username or password are not allowed to be empty"
		this.TplName = "register.html"
		return
	}

	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	user.Password = userPassword
	id, err := o.Insert(&user)
	beego.Info("User primary key id = ", id)
	if err != nil {
		this.Data["errMsg"] = "Registered user failed"
		this.TplName = "register.html"
		return
	}

	this.Ctx.WriteString("Registered successfully")
}

func (this *UserController) ShowLogin() {
	userName := this.Ctx.GetCookie("userName")
	beego.Info("userName=", userName)
	if userName != "" {
		this.Data["userName"] = userName
		this.Data["checked"] = "checked"
	} else {
		this.Data["username"] = userName
		this.Data["checked"] = ""
	}
	this.TplName = "login.html"
}

func (this *UserController) HandleLogin() {
	userName := this.GetString("userName")
	userPassword := this.GetString("password")

	if userName == "" || userPassword == "" {
		this.Data["errMsg"] = "Username and password cannot be empty"
		this.TplName = "login.html"
		return
	}

	o := orm.NewOrm()
	var user models.User
	user.Name = userName
	err := o.Read(&user, "Name")
	if err != nil {
		this.Data["errMsg"] = "No user found"
		this.TplName = "login.html"
		return
	}
	beego.Info("User password = ", user.Password)

	if user.Password != userPassword {
		this.Data["errMsg"] = "The user password is incorrect"
		this.TplName = "login.html"
		return
	}

	remember := this.GetString("remember")
	if remember == "on" {
		this.Ctx.SetCookie("userName", userName, 3600*24)
	} else {
		this.Ctx.SetCookie("userName", userName, -1)
	}

	this.SetSession("userName", userName)
	this.Redirect("/article/index", 302)
}

func (this *UserController) HandleQuit() {
	this.DelSession("userName")
	this.Redirect("/login", 302)
}
