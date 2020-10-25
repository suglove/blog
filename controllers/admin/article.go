package admin

import (
	"fmt"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"github.com/gomodule/redigo/redis"
	"goblog/models"
	"goblog/models/admin"
	redisClient "goblog/service/redis"
	"goblog/utils"
	"time"
)

type ArticleController struct {
	BaseController
}

type Article struct {
	Id      int    `json:"id"`
	Title   string `json:"title"`
	Desc    string `json:"desc"`
	Type    int    `json:"type"`
	Img     string `json:"img"`
	Content string `json:"content"`
	Status  int    `json:"status"`
	Sort    int    `json:"sort"`
}

type EasyData struct {
	Total int64            `json:"total"`
	Rows  []models.Article `json:"rows"`
}

func (this *ArticleController) Index() {

	category, err := admin.CategoryGet()
	if err != nil {
		fmt.Println(err.Error())
	}

	this.Data["category"] = category

	this.TplName = "admin/article/index.html"
}

func (this *ArticleController) List() {
	var article []models.Article
	var page int64
	var count int64
	var err error
	var pageSize int64
	var category []orm.Params

	page, _ = this.GetInt64("page", 1)
	pageSize, _ = this.GetInt64("rows", 2)
	article, count, err = admin.ArticleGetAll(page, pageSize)
	if err != nil {
		this.Rsp(false, err.Error())
		return
	}

	//fmt.Printf("cate的类型%T,值==%V\n", category, category)
	if err != nil {
		this.Rsp(false, err.Error())
		return
	}
	fmt.Println(count)
	this.Data["category"] = category
	this.Data["json"] = &map[string]interface{}{"total": count, "rows": &article, "category": category}
	this.ServeJSON()
	//this.Ctx.WriteString("11")
}

func (this *ArticleController) GetCategory() {

	category, err := admin.CategoryGet()
	//fmt.Printf("cate的类型%T,值==%V\n", category, category)
	if err != nil {
		this.Rsp(false, err.Error())
		return
	}
	this.Data["json"] = &map[string]interface{}{"status": true, "data": category}
	this.ServeJSON()
}

//添加文章
func (this *ArticleController) Add() {
	//file := this.GetFile()
	var article models.Article

	article.Title = this.GetString("title")
	article.Desc = this.GetString("desc")
	article.Content = this.GetString("content")
	article.Img = this.GetString("img")
	article.Type, _ = this.GetInt("type")
	article.Status, _ = this.GetInt("status")
	article.Sort, _ = this.GetInt("sort", 1)
	article.Isbanner, _ = this.GetInt("isbanner", 0)
	article.C_time = time.Now()
	o := orm.NewOrm()
	_, err := o.Insert(&article)
	if err != nil {
		this.Rsp(false, err.Error())

	} else {
		//删除redis
		//redisKey := "articleList:typeId:*"
		conn, _ := redisClient.PoolConnect()
		defer conn.Close()
		//redis.
		articleListKeys, err := redis.Values(conn.Do("keys", "articleList:typeId*"))
		if err != nil {
			logs.Error(err.Error())
		} else {
			conn.Do("del", articleListKeys...)
		}
		articleCountKeys, err := redis.Values(conn.Do("keys", "articleCount:typeId*"))
		if err != nil {
			logs.Error(err.Error())
		} else {
			conn.Do("del", articleCountKeys...)
		}
		//t, _ := redis.Bool(conn.Do("del", redisKey))
		//fmt.Println("t====", t)
		this.Rsp(true, "添加成功！")

	}

}

//编辑文章
func (this *ArticleController) Update() {
	var article models.Article
	o := orm.NewOrm()
	article.Id, _ = this.GetInt("id")
	if o.Read(&article) == nil {
		article.Title = this.GetString("title")
		article.Desc = this.GetString("desc")
		article.Content = this.GetString("content")
		article.Img = this.GetString("img")
		article.Type, _ = this.GetInt("type")
		article.Status, _ = this.GetInt("status")
		article.Sort, _ = this.GetInt("sort")
		if _, err := o.Update(&article, "title", "desc", "content",
			"img", "type", "status", "sort"); err != nil {
			this.Rsp(false, "修改失败！")
		} else {
			this.Rsp(true, "修改成功！")
		}
	} else {
		this.Rsp(false, "修改失败！")
	}

}

//获取文章信息
func (this *ArticleController) Get() {
	id, _ := this.GetInt("Id")
	o := orm.NewOrm()
	article := models.Article{Id: id}
	err := o.Read(&article)
	if err != nil {
		this.Rsp(false, err.Error())
	} else {
		this.Data["json"] = &map[string]interface{}{"code": 1, "msg": "查询成功！", "data": article}
		this.ServeJSON()
	}
}

//删除文章
func (this *ArticleController) Delete() {

	id, _ := this.GetInt("Id")
	fmt.Println("id======", id)
	o := orm.NewOrm()
	_, err := o.Delete(&models.Article{Id: id})

	if err != nil {
		this.Rsp(false, err.Error())
	} else {
		//删除redis
		//redisKey := "articleList:typeId:*"
		conn, _ := redisClient.PoolConnect()
		defer conn.Close()
		//redis.
		articleListKeys, err := redis.Values(conn.Do("keys", "articleList:typeId*"))
		if err != nil {
			logs.Error(err.Error())
		} else {
			conn.Do("del", articleListKeys...)
		}
		articleCountKeys, err := redis.Values(conn.Do("keys", "articleCount:typeId*"))
		if err != nil {
			logs.Error(err.Error())
		} else {
			conn.Do("del", articleCountKeys...)
		}
		this.Rsp(true, "删除成功！")
	}
}

//图片上传
func (this *ArticleController) UploadImg() {
	f, h, err := this.GetFile("info_upload_img")
	defer f.Close()
	if err != nil {
		this.Rsp(false, err.Error())
		return
	}

	folderPath, imgName, err := utils.UploadImg(h.Filename, "info_upload_img")
	if err != nil {
		this.Rsp(false, err.Error())
		return
	}
	path := folderPath + "/" + imgName
	fmt.Println("path====", path)
	this.SaveToFile("info_upload_img", path)
	data := "http://127.0.0.1:8080/" + folderPath + "/" + imgName
	this.Data["json"] = &map[string]interface{}{"code": true, "data": data, "imgPath": path}
	this.ServeJSON()
}

//图片上传封装
