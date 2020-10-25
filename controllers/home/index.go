package home

import (
	"fmt"
	"github.com/astaxie/beego"
	"goblog/models"
	"goblog/utils"
)

type IndexController struct {
	BaseController
}

func (this *IndexController) Index() {
	//声明article表
	//var article []*models.Article
	var page int64
	//var counts int64
	typeId, _ := this.GetInt("id", 0)

	page, _ = this.GetInt64(":page", 1)

	limit, _ := beego.AppConfig.Int64("pageSize")
	//page = 1

	article, counts, _ := models.RedisArticleGetAll(page, typeId)
	//fmt.Printf("类型==%T,,article的值==%v", article, len(article))
	num := len(article)
	if num == 0 {
		this.Abort("404")
	}

	//拼装点击数
	for _, v := range article {
		v.Sort = utils.GetClick(v.Id)
	}

	//fmt.Println("article==", article)
	this.Data["type"] = typeId
	this.Data["list"] = article
	this.Data["motto"] = beego.AppConfig.String("motto")

	this.Data["Paginator"] = utils.GenPaginator(page, limit, counts)
	//this.Ctx.WriteString("111")
	this.TplName = "home/index.html"

}

func (b *IndexController) Show() {
	fmt.Printf("我是show!")
}
