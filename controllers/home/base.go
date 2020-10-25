package home

import (
	"github.com/astaxie/beego"
	"github.com/gomodule/redigo/redis"
	"goblog/models"
	redisClient "goblog/service/redis"
)

type BaseController struct {
	beego.Controller
}

func (this *BaseController) Prepare() {

	//统计网站的PV数
	conn, err := redisClient.PoolConnect()
	defer conn.Close()
	if err == nil {
		//判断redsi是否存在
		redisKey := "PV"
		t, _ := redis.Bool(conn.Do("exists", redisKey))
		if t {
			//存在
			conn.Do("incr", redisKey)
		} else {
			//不存在
			conn.Do("set", redisKey, 1)
		}

	}
	//统计网站
	webName := beego.AppConfig.String("webName")
	this.Data["webName"] = webName
	var category []models.Category
	category, _ = models.RedisCategoryGetAll()
	this.Data["category"] = category
}
