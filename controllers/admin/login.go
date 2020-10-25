package admin

import (
	"github.com/astaxie/beego"
	admin2 "goblog/models/admin"
	"goblog/utils"
)

// LoginController operations for Login
type LoginController struct {
	beego.Controller
}

func (this *LoginController) Index() {
	//判断session是否有值
	user := this.GetSession("User")
	if user != nil {
		this.Redirect("/admin/index", 302)
		return
	}
	this.TplName = "admin/login.html"
}

//登录操作
func (this *LoginController) Login() {

	name := this.GetString("username")
	password := this.GetString("password")
	//查询是否有该用户
	//var admin admin2.Admin
	user, err := admin2.GetAdminInfo(name)
	password = utils.PasswdMd5(password)
	if err != nil {
		this.Rsp(1, err.Error())
	}
	if user.Password != password {
		this.Rsp(1, "你输入的密码错误")
	}
	//数据存入session
	this.SetSession("User", user)
	this.Redirect("/admin/index", 302)

}

//退出登录
func (this *LoginController) Logout() {
	//删除session
	this.DelSession("User")
	this.Redirect("/admin/login", 302)
}

func (this *LoginController) Rsp(status int, str string) {
	this.Data["json"] = &map[string]interface{}{"code": status, "msg": str}
	this.ServeJSON()
}
