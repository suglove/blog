package main

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	_ "goblog/routers"
	"goblog/utils"
)

//初始化日志
func InitLogs() (err error) {

	config := make(map[string]interface{})
	config["filename"] = "./logs/error.log"
	config["level"] = 7
	configStr, err := json.Marshal(config)
	if err != nil {
		fmt.Println(err)
		return
	}
	logs.SetLogger(logs.AdapterFile, string(configStr))
	return
}

func main() {
	InitLogs()
	if beego.BConfig.RunMode == "dev" {
		// 在beego.Run之前要配置swagger
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}

func init() {
	mysqlInit()
	beego.AddFuncMap("IndexAddOne", utils.IndexAddOne)
	beego.AddFuncMap("IndexDecrOne", utils.IndexDecrOne)
}

//数据库初始化
func mysqlInit() {
	//连接mysql
	defaultdb := beego.AppConfig.String("defaultdb")
	orm.RegisterDriver("mysql", orm.DRMySQL)

	orm.RegisterDataBase("default", "mysql", defaultdb)
}
