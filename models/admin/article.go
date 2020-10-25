package admin

import (
	"github.com/astaxie/beego/orm"
	"goblog/models"
)

//查询文章数据
func ArticleGetAll(page, pageSize int64) ([]models.Article, int64, error) {
	var articleData []models.Article
	var count int64
	offset := (page - 1) * pageSize
	//article := new(Article)
	o := orm.NewOrm()
	qs := o.QueryTable("article")
	_, err := qs.OrderBy("-id").Limit(pageSize).Offset(offset).All(&articleData)
	count, _ = qs.Count()
	return articleData, count, err

}
