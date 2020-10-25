package redisClient

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/gomodule/redigo/redis"
	"time"
)

//初始化redis连接池
var RedisPool *redis.Pool

func PoolConnect() (redis.Conn, error) {
	RedisPool = &redis.Pool{
		MaxIdle:     5000,              //最大空闲连接数
		MaxActive:   10000,             //最大连接数
		IdleTimeout: 180 * time.Second, //空闲连接超时时间
		Wait:        true,              //超过最大连接数时，是等待还是报错
		Dial: func() (redis.Conn, error) {
			redisDb := beego.AppConfig.String("redisdb")
			return redis.Dial("tcp", redisDb)
		},
	}

	//连接redis
	conn := RedisPool.Get()
	//defer conn.Close()
	//ping测试redis是否连接成功
	_, err := conn.Do("PING")
	if err != nil {
		logs.Error("redis连接失败，err:%v", err)
		//return
	}
	return RedisPool.Get(), err
}
