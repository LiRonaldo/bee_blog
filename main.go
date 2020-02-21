package main

import (
	"bee_blog/models"
	_ "bee_blog/routers"
	"github.com/astaxie/beego"
)

func init() {
	models.Init()
	beego.BConfig.WebConfig.Session.SessionOn = true

}
func main() {
	beego.Run()
}
