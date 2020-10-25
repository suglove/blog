package admin

import "github.com/astaxie/beego/orm"

type Admin struct {
	Id        int    `json:"id"`
	Name      string `json:"name"`
	Password  string `json:"password"`
	Roleid    int    `json:"roleid"`
	Rolename  string `json:"rolename"`
	Realname  string `json:"realname"`
	Nickname  string `json:"nickname"`
	Email     string `json:"email"`
	Logintime int    `json:"logintime"`
	Loginip   string `json:"loginip"`
	Addtime   int    `json:"addtime"`
	Addpeople string `json:"addpeople"`
	Status    int    `json:"status"`
}

func init() {
	//注册
	orm.RegisterModel(new(Admin))
}

//查询用户信息
func GetAdminInfo(name string) (Admin, error) {
	var admin Admin
	o := orm.NewOrm()
	err := o.QueryTable("admin").Filter("name", name).One(&admin)
	return admin, err
}
