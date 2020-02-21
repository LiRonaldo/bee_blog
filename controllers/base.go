package controllers

import (
	"bee_blog/models"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strings"
)

type baseController struct {
	beego.Controller
	controllerName string
	actionName     string
	o              orm.Ormer
}

func (p *baseController) Prepare() {
	controllerName, actionName := p.GetControllerAndAction()
	p.controllerName = strings.ToLower(controllerName[0 : len(controllerName)-10])
	p.actionName = strings.ToLower(actionName)
	p.o = orm.NewOrm()
	if strings.ToLower(p.controllerName) == "admin" && strings.ToLower(p.actionName) != "login" {
		if p.GetSession("user") == nil {
			p.History("未登录", "/admin/login")
			//p.Ctx.WriteString(p.controllerName +"==="+ p.actionName)
		}
	}

	//初始化前台页面相关元素
	if strings.ToLower(p.controllerName) == "blog" {

		p.Data["actionName"] = strings.ToLower(actionName)
		var result []*models.Config
		p.o.QueryTable(new(models.Config).TableName()).All(&result)
		configs := make(map[string]string)
		for _, v := range result {
			configs[v.Name] = v.Value
		}
		p.Data["config"] = configs
	}

}

func (b *baseController) History(msg string, url string) {
	if url == "" {
		b.Ctx.WriteString("<script>alert('" + msg + "');window.history.go(-1);</script>")
		b.StopRun()
	} else {
		b.Redirect(url, 302)
	}
}
