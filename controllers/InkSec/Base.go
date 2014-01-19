package InkSec

import (
	"InkSec/models"
	"encoding/base64"
	"fmt"
	"github.com/astaxie/beego"
	"net/smtp"
	"strconv"
	"strings"
	"time"
)

type BaseController struct {
	beego.Controller
	moduleName     string
	controllerName string
	actionName     string
	User           models.User
}

var (
	Invite      bool
	EmailVerify bool
	WebDomain   string
	SMTP_HOST   string
	SMTP_USER   string
	SMTP_PASS   string
	SMTP_PORT   string
)

const (
	base64Table = "defghijklmnopqrstuvw:01abc23xyz45ABC67XYZ89$DEFGHIJKLMNOPQRSTUVW"
)

func init() {
	if beego.AppConfig.String("Invite") == "true" {
		Invite = true
	}
	if beego.AppConfig.String("EmailVerify") == "true" {
		EmailVerify = true
	}

	WebDomain = beego.AppConfig.String("WebDomain")
	SMTP_HOST = beego.AppConfig.String("SMTP_HOST")
	SMTP_USER = beego.AppConfig.String("SMTP_USER")
	SMTP_PASS = beego.AppConfig.String("SMTP_PASS")
	SMTP_PORT = beego.AppConfig.String("SMTP_PORT")
	if SMTP_PORT == "" {
		SMTP_PORT = "25"
	}

	//smtp_host := string(RsaDecrypt(beego.AppConfig.String("SMTP_HOST")))
	//smtp_user = string(RsaDecrypt(beego.AppConfig.String("SMTP_USER")))
	//smtp_pass := string(RsaDecrypt(beego.AppConfig.String("SMTP_PASS")))
	fmt.Println("InkSec Package Init()")
}

func (this *BaseController) Prepare() {
	controllerName, actionName := this.GetControllerAndAction()
	this.moduleName = "InkSec"
	this.controllerName = controllerName[0 : len(controllerName)-10]
	this.actionName = actionName
	fmt.Println("BaseController Prepare()")
	//fmt.Printf("%s\n%s\n%s\n%s\n", controllerName, actionName, this.controllerName, this.actionName)
}

func (this *BaseController) IsLogin() {
	if !this.Auth() {
		this.Redirect("/user/?a=login", 302)
	}
}

func (this *BaseController) Auth() bool {
	//Session
	UserId := this.GetSession("UserId")
	if UserId != nil {
		UserName := this.GetSession("UserName")
		Password := this.GetSession("Password")
		LoginIP := this.GetSession("LoginIP")
		LoginTime := this.GetSession("LoginTime")
		this.User.UserId = UserId.(int64)
		this.User.Read("UserId")
		if this.User.UserName == UserName && this.User.Password == Password && this.User.LoginIP == LoginIP && this.User.LoginTime.Format("2006-01-02 15:04:05") == LoginTime {
			//fmt.Printf("%s\n", "sess_true")
			return true
		}
	}
	//Cookies
	Cookies := strings.Split(Base64Decode(this.Ctx.GetCookie("InkSec_Auth")), "|")
	if len(Cookies) == 3 {
		UserId = Cookies[0]
		AuthKey := Cookies[1]
		Password := Cookies[2]
		uId, err := strconv.ParseInt(UserId.(string), 10, 64)
		if err == nil {
			this.User.UserId = uId
			this.User.Read("UserId")
			if Password == models.MD5(this.User.Password) && AuthKey == models.MD5(this.User.UserName+this.User.LoginTime.Format("2006-01-02 15:04:05")) {
				//fmt.Printf("%s\n", "cook_true")
				return true
			}
		}
	}
	return false
}

func (this *BaseController) SaveSession() {
	Remember := strings.TrimSpace(this.GetString("Remember"))
	UserId := strconv.FormatInt(this.User.UserId, 10)
	AuthKey := models.MD5(this.User.UserName + this.User.LoginTime.Format("2006-01-02 15:04:05"))
	Password := models.MD5(this.User.Password)
	if Remember == "yes" {
		this.Ctx.SetCookie("InkSec_Auth", Base64Encode(UserId+"|"+AuthKey+"|"+Password), 31*86400, "/")
	} else {
		this.Ctx.SetCookie("InkSec_Auth", Base64Encode(UserId+"|"+AuthKey+"|"+Password), 0, "/")
	}
	this.SetSession("UserId", this.User.UserId)
	this.SetSession("UserName", this.User.UserName)
	this.SetSession("Password", this.User.Password)
	this.SetSession("LoginIP", this.User.LoginIP)
	this.SetSession("LoginTime", this.User.LoginTime.Format("2006-01-02 15:04:05"))
}

func (this *BaseController) ClearSession() {
	this.DelSession("UserId")
	this.DelSession("UserName")
	this.DelSession("Password")
	this.DelSession("LoginIP")
	this.DelSession("LoginTime")
	this.DestroySession()
	this.Ctx.SetCookie("InkSec_Auth", "", -1, "/")
}

func (this *BaseController) GetClientIP() string {
	return this.Ctx.Input.IP()
}

func (this *BaseController) IsPost() bool {
	return this.Ctx.Request.Method == "POST"
}

func (this *BaseController) Show(tpl ...string) {
	folder := this.controllerName
	file := this.actionName
	if len(tpl) != 0 {
		file = tpl[0]
	}
	//fmt.Printf("%s\n%s\n", folder, file)
	//this.Layout = folder + "/layout.html"
	this.TplNames = folder + "/" + file + ".html"
}

//////////////////////////////////////////////////////////////////////////////
func GetErrCode(e error) string {
	return e.Error()
	index := strings.Index(e.Error(), ":")
	if index == -1 {
		return "未知错误"
	}
	return e.Error()[0:index]
}

func MailCheck(s string) bool {
	unsafe := []string{"'", "`", "~", "!", "$", "%", "^", "&", "*", "(", ")", "-", "+", "=", "[", "]", "{", "}", "\\", ";", "\"", ":", "<", ">", "/"}
	for i := 0; i < len(unsafe); i++ {
		if strings.Index(s, unsafe[i]) != -1 {
			return true
		}
	}
	return false
}

func StdBase64Encode(text string) string {
	result := base64.StdEncoding.EncodeToString([]byte(text))
	return result
}

func StdBase64Decode(base64Data string) string {
	result, _ := base64.StdEncoding.DecodeString(base64Data)
	return string(result)
}

func Base64Encode(text string) string {
	result := base64.NewEncoding(base64Table).EncodeToString([]byte(text))
	return result
}

func Base64Decode(base64Data string) string {
	result, _ := base64.NewEncoding(base64Table).DecodeString(base64Data)
	return string(result)
}

func SendVerifyMail(u *models.User) error {
	Auth := smtp.PlainAuth("", SMTP_USER, SMTP_PASS, SMTP_HOST)
	VerifyUrl := "http://" + WebDomain + "/user/verify/?code=" + Base64Encode(strconv.FormatInt(u.UserId, 10)+"|"+models.MD5(u.Email+u.UserName+u.Password)+"|"+strconv.FormatInt(time.Now().UnixNano(), 10))
	ContactUsUrl := "http://" + WebDomain + "/feedback/"
	Content_Type := "Content-Type: text/html; charset=UTF-8"
	Subject := "InkSec CTF用户激活邮件"
	Body := `
<div name="InkSec">
	<div style="height:36px;background-color:#E3FEFF;font-size:20px;padding:15px 0 0 25px;line-height:22px;border:1px solid #ddd;">
		<strong>InkSec CTF 用户邮箱验证</strong>
	</div>
	<div style="padding-top:20px;border:1px solid #ddd;border-top:none;background-color:#fff">
		<p style="text-indent:2em;padding-bottom:12px;width:620px;font-size:12px;font-family:'宋体';line-height:20px;padding-left: 36px">
		感谢您注册InkSec CTF帐号！
		</p>
		<p style="text-indent:2em;padding-bottom:12px;width:620px;font-size:12px;font-family:'宋体';line-height:20px;padding-left: 36px">
		<strong> 您的帐号为：
			<span style="color:#00f;">` + u.UserName + `</span>
		</strong>
		</p>
		<p style="text-indent:2em;padding-bottom:12px;width:620px;font-size:12px;font-family:'宋体';line-height:20px;padding-left: 36px"> 请点击以下链接完成帐号注册：<br>
			<a href="` + VerifyUrl + `" target="_blank">` + VerifyUrl + `</a><br>
			(如链接无法点击，可以将此链接复制到浏览器地址栏打开页面) </p>
		<h3 style="height:28px;font-size:12px;padding-left:36px"> 温馨提示：</h3>
		<p style="text-indent:2em;padding-bottom:12px;width:620px;font-size:12px;font-family:'宋体';line-height:20px;padding-left: 36px"> 以上链接仅72小时内有效，过期后需重新验证，请尽快点击哦！</p>
		<p style="text-indent:2em;padding-bottom:12px;width:620px;font-size:12px;font-family:'宋体';line-height:20px;padding-left: 36px"> 本邮件为系统自动生成，请勿回复。如有问题，请<a href="` + ContactUsUrl + `" target="_blank">联系我们</a></p>
		<div style="padding-top:13px;height:30px;background-color:#E5F2FB;color:#666;text-align:center;font-family:Arial;">
			Copyright &copy; 2014 InkSec. All Rights Reserved
		</div>
	</div>
</div>
	`
	Msg := []byte("To: " + u.Email + "\r\nFrom: " + "InkSec CTF" + "<" + SMTP_USER + ">\r\nSubject: " + Subject + "\r\n" + Content_Type + "\r\n\r\n" + Body)
	To := []string{u.Email}
	Err := smtp.SendMail(SMTP_HOST+":"+SMTP_PORT, Auth, SMTP_USER, To, Msg)
	return Err
}

func SendResetMail(u *models.User) error {
	Auth := smtp.PlainAuth("", SMTP_USER, SMTP_PASS, SMTP_HOST)
	VerifyUrl := "http://" + WebDomain + "/user/reset/?code=" + Base64Encode(strconv.FormatInt(u.UserId, 10)+"|"+models.MD5(u.Email+u.UserName+u.Password)+"|"+strconv.FormatInt(time.Now().UnixNano(), 10))
	ContactUsUrl := "http://" + WebDomain + "/feedback/"
	Content_Type := "Content-Type: text/html; charset=UTF-8"
	Subject := "InkSec CTF用户密码重置邮件"
	Body := `
<div name="InkSec">
	<div style="height:36px;background-color:#E3FEFF;font-size:20px;padding:15px 0 0 25px;line-height:22px;border:1px solid #ddd;">
		<strong>InkSec CTF 用户密码重置</strong>
	</div>
	<div style="padding-top:20px;border:1px solid #ddd;border-top:none;background-color:#fff">
		<p style="text-indent:2em;padding-bottom:12px;width:620px;font-size:12px;font-family:'宋体';line-height:20px;padding-left: 36px">
		感谢您使用InkSec CTF！
		</p>
		<p style="text-indent:2em;padding-bottom:12px;width:620px;font-size:12px;font-family:'宋体';line-height:20px;padding-left: 36px">
		<strong> 您的帐号为：
			<span style="color:#00f;">` + u.UserName + `</span>
		</strong>
		</p>
		<p style="text-indent:2em;padding-bottom:12px;width:620px;font-size:12px;font-family:'宋体';line-height:20px;padding-left: 36px"> 请点击以下链接进行密码重置：<br>
			<a href="` + VerifyUrl + `" target="_blank">` + VerifyUrl + `</a><br>
			(如链接无法点击，可以将此链接复制到浏览器地址栏打开页面) </p>
		<h3 style="height:28px;font-size:12px;padding-left:36px"> 温馨提示：</h3>
		<p style="text-indent:2em;padding-bottom:12px;width:620px;font-size:12px;font-family:'宋体';line-height:20px;padding-left: 36px"> 以上链接仅72小时内有效，过期后需重新验证，请尽快点击哦！</p>
		<p style="text-indent:2em;padding-bottom:12px;width:620px;font-size:12px;font-family:'宋体';line-height:20px;padding-left: 36px"> 本邮件为系统自动生成，请勿回复。如有问题，请<a href="` + ContactUsUrl + `" target="_blank">联系我们</a></p>
		<div style="padding-top:13px;height:30px;background-color:#E5F2FB;color:#666;text-align:center;font-family:Arial;">
			Copyright &copy; 2014 InkSec. All Rights Reserved
		</div>
	</div>
</div>
	`
	Msg := []byte("To: " + u.Email + "\r\nFrom: " + "InkSec CTF" + "<" + SMTP_USER + ">\r\nSubject: " + Subject + "\r\n" + Content_Type + "\r\n\r\n" + Body)
	To := []string{u.Email}
	Err := smtp.SendMail(SMTP_HOST+":"+SMTP_PORT, Auth, SMTP_USER, To, Msg)
	return Err
}
