package routers

import (
	//"InkSec/controllers/admin"
	"InkSec/controllers/InkSec"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

func init() {
	if beego.AppConfig.String("runmode") == "dev" {
		orm.Debug = true
	}
	//首页
	beego.Router("/", &InkSec.IndexController{})
	beego.Router("/index.php", &InkSec.IndexController{})
	//beego.Router("/test", &InkSec.BaseController{}, "*:BaseMain")
	//用户相关

	beego.Router("/user/", &InkSec.UserController{}, "*:Main")
	beego.Router("/user/index.php", &InkSec.UserController{}, "*:Main")
	beego.Router("/user/verify/", &InkSec.UserController{}, "*:VerifyEmail")
	beego.Router("/user/verify/index.php", &InkSec.UserController{}, "*:VerifyEmail")
	beego.Router("/user/modify/", &InkSec.UserController{}, "*:Modify")
	beego.Router("/user/modify/index.php", &InkSec.UserController{}, "*:Modify")
	beego.Router("/user/forgot/", &InkSec.UserController{}, "*:ForgotPassword")
	beego.Router("/user/forgot/index.php", &InkSec.UserController{}, "*:ForgotPassword")
	beego.Router("/user/reset/", &InkSec.UserController{}, "*:VerifyReset")
	beego.Router("/user/reset/index.php", &InkSec.UserController{}, "*:VerifyReset")

	//查看排行
	//beego.Router("/rank/", &InkSec.MainController{})
	//beego.Router("/rank/index.php", &InkSec.MainController{})
	//留言板
	beego.Router("/chat/", &InkSec.ChatController{}, "*:Main")
	beego.Router("/chat/index.php", &InkSec.ChatController{}, "*:Main")
	//下载
	//beego.Router("/download.php", &InkSec.MainController{})
	//CTF
	//beego.Router("/ctf/", &InkSec.MainController{})
	//beego.Router("/ctf/:id:int/", &InkSec.MainController{})
	//后台管理
	//beego.Router("/InkSec_Admin/", &InkSec.MainController{})
	//beego.Router("/InkSec_Admin/index.php", &InkSec.MainController{})
}
