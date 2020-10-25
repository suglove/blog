package home

import (
	"github.com/astaxie/beego"
	"goblog/models"
	"goblog/utils"
)

type CategoryController struct {
	BaseController
}

func (this *CategoryController) Index() {
	//声明article表
	//var article []models.Article
	var page int64
	//var counts int64
	typeId, _ := this.GetInt(":id", 0)

	page, _ = this.GetInt64(":page", 1)

	limit, _ := beego.AppConfig.Int64("pageSize")
	//page = 1

	article, counts, _ := models.RedisArticleGetAll(page, typeId)

	num := len(article)
	if num == 0 {

		this.Abort("404")
	}

	//拼装点击数
	for _, v := range article {
		v.Sort = utils.GetClick(v.Id)
	}

	this.Data["type"] = typeId
	this.Data["list"] = article

	this.Data["Paginator"] = utils.GenPaginator(page, limit, counts)
	//this.Ctx.WriteString("111")
	this.TplName = "home/category/index.html"

}
