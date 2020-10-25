package models

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	redisClient "goblog/service/redis"
	"log"
	"strconv"
	"time"
)

type Article struct {
	Id       int       `json:"id";orm:"pk"`
	Title    string    `json:"title"`
	Desc     string    `json:"desc"`
	Type     int       `json:"type"`
	Img      string    `json:"img"`
	Content  string    `json:"content"`
	Status   int       `json:"status"`
	C_time   time.Time `json:"c_time";form:"-"`
	Sort     int       `json:"sort"`
	Isbanner int       `json:"isbanner";form:"-"`
}

type ArticleData struct {
	Id     int       `json:"id"`
	Title  string    `json:"title"`
	Desc   string    `json:"desc"`
	Type   int       `json:"type"`
	Img    string    `json:"img"`
	C_time time.Time `json:"c_time"`
	Sort   int       `json:"sort"`
}

func init() {
	//注册文章模型
	orm.RegisterModel(new(Article))
}

//根据文章ID获取信息
func RedisArticleGet(id int) (Article, error) {
	var article Article
	var err error
	redisKey := "article:id" + strconv.Itoa(id)
	//连接redis
	conn, _ := redisClient.PoolConnect()
	defer conn.Close()
	//判断key是否存在
	exists, _ := redis.Bool(conn.Do("exists", redisKey))
	if exists {
		values, _ := redis.Values(conn.Do("hgetall", redisKey))
		err = redis.ScanStruct(values, &article)
	} else {
		//没有redis查询mysql

		o := orm.NewOrm()
		err := o.QueryTable(article).Filter("id", id).One(&article)
		if err == nil {
			//写入redis
			_, err := conn.Do("hmset", redis.Args{}.Add(redisKey).AddFlat(article))
			if err == nil {
				conn.Do("expire", redisKey, 86400)
			}
		}

	}
	return article, err
}

//查询文章所有数据
func ArticleGetAll(page int64, typeId int) ([]Article, int64, error) {
	var (
		//articleData []ArticleData
		article []Article
		//articleData []Article
		err error
	)
	limit, _ := beego.AppConfig.Int64("pageSize")
	offset := (page - 1) * limit // 偏移量
	// 获取 QuerySeter 对象，article 为表名
	//article := new(Article)
	o := orm.NewOrm()
	qs := o.QueryTable("article")
	qs = qs.Filter("status", 1)
	if typeId != 0 {
		qs = qs.Filter("type", typeId)
	}
	// 获取数据
	_, err = qs.OrderBy("-id").Limit(limit).Offset(offset).All(&article)

	//_, err = qs.OrderBy("-id").Limit(limit).Offset(offset).All(&articleData)

	// 统计
	count, _ := qs.Count()

	return article, count, err
}

//查询文章所有数据 redis
func RedisArticleGetAll(page int64, typeId int) ([]*Article, int64, error) {
	var (
		article []*Article
		count   int64
		//err     error
	)
	limit, _ := beego.AppConfig.Int64("pageSize")
	offset := (page - 1) * limit // 偏移量

	//查询redis中是否有数据，没有就在数据库查询

	conn, err := redisClient.PoolConnect()
	if err != nil {
		logs.Error("redis连接失败！")
	}

	defer conn.Close()
	redisKey := "articleList:typeId:" + strconv.Itoa(typeId) + "limit" + strconv.Itoa(int(limit))
	redisCount := "articleCount:typeId" + strconv.Itoa(typeId)
	//判断rediskey是否已存在
	exists, err := redis.Bool(conn.Do("EXISTS ", redisKey))
	if exists {
		fmt.Println("走redis")
		count, _ = redis.Int64(conn.Do("get", redisCount))
		values, _ := redis.Values(conn.Do("lrange", redisKey, "0", "-1"))
		var articleInfo *Article
		for _, v := range values {
			err = json.Unmarshal(v.([]byte), &articleInfo)
			if err == nil {
				article = append(article, articleInfo)
			}
		}
	} else {
		fmt.Println("没有走redis")
		//不存在去数据库查询
		// 获取 QuerySeter 对象，article 为表名
		//article := new(Article)
		o := orm.NewOrm()
		qs := o.QueryTable("article")
		qs = qs.Filter("status", 1)
		if typeId != 0 {
			qs = qs.Filter("type", typeId)
		}
		// 获取数据
		_, err := qs.OrderBy("-id").Limit(limit).Offset(offset).All(&article)
		log.Printf("err的类型是%T,err=%v \n", err, err)
		// 统计
		count, _ = qs.Count()
		//遍历数据把信息json化保存
		if err == nil {
			for _, v := range article {
				jsonValue, err := json.Marshal(v)
				if err == nil {
					//保存redis
					conn.Do("rpush", redisKey, jsonValue)
				}
			}
			_, err = conn.Do("set", redisCount, count)
			if err != nil {
				logs.Error("redis写入失败！")
			}
			conn.Do("expire", redisKey, 86400)
			conn.Do("expire", redisCount, 86400)
		} else {
			return article, count, err
		}

	}

	return article, count, err
}

//设置文章阅读数量
func SetArticleClick(id int) {
	redisKey := "articleClick"
	conn, err := redisClient.PoolConnect()
	defer conn.Close()
	if err != nil {
		logs.Error(err.Error())
	}
	//判断是否存在key
	t, _ := redis.Bool(conn.Do("exists", redisKey))
	if t {
		//判断member是否存在
		score, _ := redis.Int(conn.Do("ZSCORE", redisKey, id))
		if score > 0 {
			//存在值加1
			conn.Do("ZINCRBY", redisKey, 1, id)
		} else {
			//不存在就设置
			conn.Do("zadd", redisKey, 1, id)
		}
	} else {
		conn.Do("zadd", redisKey, 1, id)
	}

}
