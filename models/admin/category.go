package admin

import (
	"github.com/astaxie/beego/orm"
	"goblog/models"
)

//查询类型
func CategoryGet() (categoryData []orm.Params, err error) {
	category := new(models.Category)
	o := orm.NewOrm()
	qs := o.QueryTable(category)
	_, err = qs.OrderBy("id").Values(&categoryData, "id", "title")

	return categoryData, err
}
