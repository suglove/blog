package models

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	redisClient "goblog/service/redis"
	"time"
)

type Category struct {
	Id     int       `json:"id"`
	Title  string    `json:"title"`
	Status string    `json:"status"`
	C_time time.Time `json:"c_time"`
	Sort   int       `json:"sort"`
}

func init() {
	//注册文章模型
	orm.RegisterModel(new(Category))
}

//根据类型数据
func CategoryGetAll() ([]Category, error) {
	o := orm.NewOrm()
	var categoryData []Category
	category := new(Category)
	qs := o.QueryTable(category)
	qs = qs.Filter("status", 1)
	_, err := qs.All(&categoryData, "id", "title")
	return categoryData, err
}

//类型数据 redis
func RedisCategoryGetAll() ([]Category, error) {
	var categoryData []Category
	var err error
	//连接redis
	conn, _ := redisClient.PoolConnect()
	defer conn.Close()

	redisKey := "ategoryList"
	//判断key是否存在
	exists, _ := redis.Bool(conn.Do("exists", redisKey))
	if exists {
		values, _ := redis.Values(conn.Do("lrange", redisKey, 0, -1))
		var CategoryInfo Category
		for _, v := range values {
			err := json.Unmarshal(v.([]byte), &CategoryInfo)
			if err == nil {
				categoryData = append(categoryData, CategoryInfo)
			}
		}
	} else {
		//不存在查询mysql
		o := orm.NewOrm()

		category := new(Category)
		qs := o.QueryTable(category)
		qs = qs.Filter("status", 1)
		_, err := qs.All(&categoryData, "id", "title")
		if err != nil {
			logs.Error(err)
		}
		//遍历数据

		for _, v := range categoryData {
			jsonValue, err := json.Marshal(v)
			if err == nil {
				//保存redis
				conn.Do("rpush", redisKey, jsonValue)
			}
		}
		conn.Do("expire", redisKey, 86400)
	}

	return categoryData, err
}
