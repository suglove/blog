package admin

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/session"
	admin2 "goblog/models/admin"
)

var globalSessions *session.Manager

type BaseController struct {
	beego.Controller
	User admin2.Admin
}

func init() {
	sessionConfig := &session.ManagerConfig{
		CookieName:      "gosessionid",
		EnableSetCookie: true,
		Gclifetime:      3600,
		Maxlifetime:     3600,
		Secure:          false,
		CookieLifeTime:  3600,
		ProviderConfig:  "./tmp",
	}
	globalSessions, _ = session.NewManager("memory", sessionConfig)
	go globalSessions.GC()
}

//检查用户是否登录
func (this *BaseController) Prepare() {

	//查询session
	user := this.GetSession("User")
	if user != nil {
		this.User = user.(admin2.Admin)
		this.Data["user"] = user
	} else {
		//fmt.Println("我继承了这个方法")
		this.Redirect("/admin/login", 302)
	}

}

//公共方法库
func (this *BaseController) Rsp(status bool, str string) {
	this.Data["json"] = &map[string]interface{}{"code": status, "msg": str}
	this.ServeJSON()
}
