package home

import (
	"goblog/models"
)

type InfoController struct {
	BaseController
}

func (this *InfoController) Index() {
	var article models.Article
	//var err error
	id, _ := this.GetInt(":id")
	//设置文章点击数
	models.SetArticleClick(id)
	article, _ = models.RedisArticleGet(id)
	if article.Id == 0 {
		//没有数据返回404
		this.Abort("404")
	}
	//设置点击数到redis

	this.Data["article"] = article

	//this.Data["json"] = article
	//this.ServeJSON()
	this.TplName = "home/info/index.html"
}
