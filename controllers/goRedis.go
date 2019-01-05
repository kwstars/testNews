package controllers

import (
	"github.com/astaxie/beego"
)

type GoRedis struct {
	beego.Controller
}

func (this *GoRedis) ShowGet() {
/*	conn, err := redis.Dial("tcp", ":6379")
	//defer conn.Close()
	if err != nil {
		beego.Error("redis connect err", err)
		return
	}*/

	//resp, err := conn.Do("mset", "class2", "test01","class3","test02")
	//resp, err = redis.String(resp, err)
	//beego.Info("回复值=", resp)

	//resp, err := conn.Do("lrange", "l1", "0", "-1")
	//re, err := redis.Values(resp, err)


}
