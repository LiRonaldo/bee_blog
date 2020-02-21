package routers

import (
	"bee_blog/controllers"
	"github.com/astaxie/beego"
)

func init() {
	beego.AutoRouter(&controllers.AdminController{})
}
