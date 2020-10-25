package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
)

type MainController struct {
	beego.Controller
}

func (c *MainController) Get() {
	scheme := "http://"
	if c.Ctx.Request.TLS != nil {
		scheme = "https://"
	}
	host := scheme + c.Ctx.Request.Host
	c.Data["host"] = host
	c.Data["webName"] = beego.AppConfig.String("webName")
	c.Redirect("home/index", 302)
	//c.TplName = "index.html"
}

func (this *MainController) URLMapping() {
	this.Mapping("Index", this.Index)
}

func (this *MainController) Index() {
	fmt.Println("111")
	this.Ctx.WriteString("测试")
}
