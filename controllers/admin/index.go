package admin

import (
	"github.com/gomodule/redigo/redis"
	redisClient "goblog/service/redis"
)

type IndexController struct {
	BaseController
}

func (this *IndexController) Index() {
	this.TplName = "admin/index.html"
}

func (this *IndexController) Welcome() {
	conn, err := redisClient.PoolConnect()
	redisKey := "PV"
	var pv int
	if err == nil {
		pv, _ = redis.Int(conn.Do("get", redisKey))

	}
	this.Data["pv"] = pv
	this.TplName = "admin/welcome.html"
}
