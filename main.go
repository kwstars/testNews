package main

import (
	_ "testNews/routers"
	"github.com/astaxie/beego"
	_ "testNews/models"
)

func main() {
	beego.AddFuncMap("pre", getPre)
	beego.AddFuncMap("next", getNext)
	beego.Run()
}

func getPre(pageIndex int) int {
	if pageIndex-1 <= 0 {
		return pageIndex
	}
	return pageIndex - 1
}

func getNext(pageIndex, pageCount int) int {
	if pageIndex+1 > pageCount {
		return pageIndex
	}
	return pageIndex + 1
}
