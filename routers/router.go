package routers

import (
	"github.com/astaxie/beego"
	"goblog/controllers"
	"goblog/controllers/admin"
	"goblog/controllers/home"
)

func init() {

	//home模块
	//初始化 namespace
	ns :=
		beego.NewNamespace("/home",
			/*beego.NSInclude(
			&home.IndexController{},
			),*/
			beego.NSRouter("/index/:page([0-9]*)", &home.IndexController{}, "get:Index"),
			beego.NSRouter("/info/:id([0-9]*)", &home.InfoController{}, "get:Index"),
			beego.NSRouter("/category/id/:id([0-9]*)/page/:page([0-9]*)", &home.CategoryController{}, "get:Index"),
		)
	//注册 namespace
	beego.AddNamespace(ns)

	beego.Router("/", &controllers.MainController{})

	//admin模块

	adminNs :=
		beego.NewNamespace("/admin",
			//登录
			beego.NSRouter("/login",
				&admin.LoginController{}, "get:Index"),
			//登录
			beego.NSRouter("/login",
				&admin.LoginController{}, "post:Login"),
			//退出
			beego.NSRouter("/logout",
				&admin.LoginController{}, "get:Logout"),
			//IndexController
			beego.NSRouter("/index",
				&admin.IndexController{}, "get:Index"),
			beego.NSRouter("/index/welcome",
				&admin.IndexController{}, "get:Welcome"),
			//ArticleController
			beego.NSRouter("/article",
				&admin.ArticleController{}, "get:Index"),
			//list
			beego.NSRouter("/article/list",
				&admin.ArticleController{}, "*:List"),
			//添加
			beego.NSRouter("/article/add",
				&admin.ArticleController{}, "POST:Add"),
			//编辑
			beego.NSRouter("/article/edit",
				&admin.ArticleController{}, "POST:Update"),
			//编辑
			beego.NSRouter("/article/get",
				&admin.ArticleController{}, "POST:Get"),
			//删除
			beego.NSRouter("/article/del",
				&admin.ArticleController{}, "POST:Delete"),
			//删除
			beego.NSRouter("/article/cate",
				&admin.ArticleController{}, "*:GetCategory"),
			//上传图片
			beego.NSRouter("/article/uploadImg",
				&admin.ArticleController{}, "*:UploadImg"),
		)

	//注册 namespace
	beego.AddNamespace(adminNs)
}
