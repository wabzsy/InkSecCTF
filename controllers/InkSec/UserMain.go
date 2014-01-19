package InkSec

import (
	"InkSec/models"
	"fmt"
	"github.com/astaxie/beego/validation"
	"strconv"
	"strings"
	"time"
)

type UserController struct {
	BaseController
	Invite models.InviteCode
}

func (this *UserController) Prepare() {
	this.BaseController.Prepare()
	fmt.Println("UserPrepare()")
}

func (this *UserController) Main() {
	action := this.Input().Get("a")
	if action == "login" {
		this.Login()
	} else if action == "reg" {
		this.Register()
	} else if action == "logout" {
		this.Logout()
	} else if action == "test" {
		this.Temp()
	} else {
		this.Index()
	}
}

func (this *UserController) Register() {
	User := make(map[string]interface{})
	Errmsg := make(map[string]interface{})
	if this.IsPost() {
		UserName := strings.TrimSpace(this.GetString("UserName"))
		Password := strings.TrimSpace(this.GetString("Password"))
		ConfirmPassword := strings.TrimSpace(this.GetString("ConfirmPassword"))
		Email := strings.TrimSpace(this.GetString("Email"))
		InviteKey := strings.TrimSpace(this.GetString("InviteKey"))

		User["UserName"] = UserName
		User["Password"] = Password
		User["ConfirmPassword"] = ConfirmPassword
		User["Email"] = Email
		User["InviteKey"] = InviteKey

		valid := validation.Validation{}

		if v := valid.Required(UserName, "UserName"); !v.Ok {
			Errmsg["UserName"] = "请输入用户名"
		} else if v := valid.MaxSize(UserName, 16, "UserName"); !v.Ok {
			Errmsg["UserName"] = "用户名长度不能大于16个字符"
		} else if v := valid.MinSize(UserName, 3, "UserName"); !v.Ok {
			Errmsg["UserName"] = "用户名长度不能小于3个字符"
		} else if this.User.Exist(UserName) {
			Errmsg["UserName"] = "用户名已存在"
		}

		if v := valid.Required(Password, "Password"); !v.Ok {
			Errmsg["Password"] = "请输入密码"
		} else if v := valid.MaxSize(Password, 18, "Password"); !v.Ok {
			Errmsg["Password"] = "密码长度不能大于18个字符"
		} else if v := valid.MinSize(Password, 6, "Password"); !v.Ok {
			Errmsg["Password"] = "密码长度不能小于6个字符"
		}

		if v := valid.Required(ConfirmPassword, "ConfirmPassword"); !v.Ok {
			Errmsg["ConfirmPassword"] = "请确认您输入密码"
		} else if Password != ConfirmPassword {
			Errmsg["ConfirmPassword"] = "两次输入的密码不一致"
		}

		if v := valid.Required(Email, "Email"); !v.Ok {
			Errmsg["Email"] = "请输入Email地址"
		} else if v := valid.Email(Email, "Email"); !v.Ok || MailCheck(Email) {
			Errmsg["Email"] = "Email地址无效"
		} else if this.User.MailExist(Email) {
			Errmsg["Email"] = "邮箱已被使用"
		}

		if Invite && len(InviteKey) == 0 {
			Errmsg["InviteKey"] = "请输入邀请码"
		}
		if len(InviteKey) != 32 && len(InviteKey) != 0 {
			Errmsg["InviteKey"] = "邀请码不合法"
		} else if len(InviteKey) == 32 && !this.Invite.VerifyInviteCode(InviteKey) {
			Errmsg["InviteKey"] = "邀请码不存在或已被使用"
		}

		if len(Errmsg) == 0 {
			this.User.UserName = UserName
			this.User.Password = models.MD5(Password)
			this.User.Email = Email
			this.User.LoginIP = this.GetClientIP()
			this.User.Enabled = 1
			this.User.SignSum = 1
			this.Invite.UserId = &this.User
			this.Invite.Used = 1
			this.Invite.UsedTime = time.Now()
			if err := this.User.Insert(); err != nil {
				Errmsg["Other"] = GetErrCode(err)
			} else if err := SendVerifyMail(&this.User); err != nil {
				Errmsg["Other"] = GetErrCode(err)
			} else if InviteKey != "" {
				if err := this.Invite.Update(); err != nil {
					Errmsg["Other"] = GetErrCode(err)
				}
			}
			if len(Errmsg) == 0 {
				this.SaveSession()
				this.Redirect("/user/", 302)
			}
		}
	}
	this.Data["User"] = User
	this.Data["Errmsg"] = Errmsg
	this.Show("Register")
}

func (this *UserController) Index() {
	this.IsLogin()
	User := make(map[string]interface{})

	User["UserId"] = this.User.UserId
	User["UserName"] = this.User.UserName
	User["Email"] = this.User.Email

	if this.EmailIsVerified() {
		User["EmailVerify"] = "已验证"
	} else {
		User["EmailVerify"] = "未验证"
	}

	User["LoginIP"] = this.User.LoginIP
	User["LoginTime"] = this.User.LoginTime
	User["RegTime"] = this.User.RegTime
	User["QQ"] = this.User.QQ
	User["Rank"] = this.User.Rank
	User["GoldCoin"] = this.User.GoldCoin
	User["SignSum"] = this.User.SignSum
	User["School"] = this.User.School
	User["StuId"] = this.User.StuId

	this.Data["User"] = User
	this.Show("Index")
}

func (this *UserController) Login() {
	User := make(map[string]interface{})
	Errmsg := make(map[string]interface{})
	if this.IsPost() {
		UserName := strings.TrimSpace(this.GetString("UserName"))
		Password := strings.TrimSpace(this.GetString("Password"))

		User["UserName"] = UserName

		valid := validation.Validation{}

		if v := valid.Required(UserName, "UserName"); !v.Ok {
			Errmsg["UserName"] = "请输入用户名"
		} else if v := valid.MaxSize(UserName, 16, "UserName"); !v.Ok {
			Errmsg["UserName"] = "用户名长度不能大于16个字符"
		} else if v := valid.MinSize(UserName, 3, "UserName"); !v.Ok {
			Errmsg["UserName"] = "用户名长度不能小于3个字符"
		} else if !this.User.Exist(UserName) {
			Errmsg["UserName"] = "用户名不存在"
		}

		if v := valid.Required(Password, "Password"); !v.Ok {
			Errmsg["Password"] = "请输入密码"
		} else if v := valid.MaxSize(Password, 18, "Password"); !v.Ok {
			Errmsg["Password"] = "密码长度不能大于18个字符"
		} else if v := valid.MinSize(Password, 6, "Password"); !v.Ok {
			Errmsg["Password"] = "密码长度不能小于6个字符"
		}

		if len(Errmsg) == 0 {
			this.User.UserName = UserName
			Password = models.MD5(Password)
			if err := this.User.Read("UserName"); err != nil {
				Errmsg["Other"] = GetErrCode(err)
			} else {
				if Password != this.User.Password {
					Errmsg["Other"] = "用户名或密码错误!"
				} else {
					this.User.LoginIP = this.GetClientIP()
					this.User.LoginTime = time.Now()
					this.User.SignSum = this.User.SignSum + 1
					this.User.Update()
					this.SaveSession()
					this.Redirect("/user/", 302)
				}
			}
		}
	}
	this.Data["User"] = User
	this.Data["Errmsg"] = Errmsg
	this.Show("Login")
}

func (this *UserController) Logout() {
	this.ClearSession()
	this.Redirect("/user/", 302)
}

func (this *UserController) Temp() {
	Errmsg := make(map[string]interface{})
	/*
		this.User.UserName = "wabzsy"
		this.User.Read("UserName")
		Code := this.Invite.GenerateInviteCode(&this.User, 5)
		for i := 0; i < len(Code); i++ {
			Errmsg["code"] += Code[i]
		}
	*/
	this.Data["Errmsg"] = Errmsg
	this.Show("Temp")
}

func (this *UserController) Verify() string {
	Msg := "Success"
	code := strings.Split(Base64Decode(strings.TrimSpace(this.GetString("code"))), "|")
	if len(code) != 3 {
		Msg = "验证失败!<!--Err:AuthCode错误-->"
	} else {
		uIdStr, Key, timeStampStr := code[0], code[1], code[2]
		uId, idErr := strconv.ParseInt(uIdStr, 10, 64)
		timeStamp, timeErr := strconv.ParseInt(timeStampStr, 10, 64)
		if idErr != nil || timeErr != nil {
			Msg = "验证失败!<!--Err:uId或时间戳转换失败-->"
		} else if len(Key) != 32 {
			Msg = "验证失败!<!--Err:AuthKey不合法-->"
		} else if len(timeStampStr) != 19 {
			Msg = "验证失败!<!--Err:时间戳不合法-->"
		} else if time.Now().UnixNano() > timeStamp+(72*60*60*1000000000) {
			Msg = "验证失败!<!--Err:认证超时-->"
		} else {
			this.User.UserId = uId
			if err := this.User.Read("UserId"); err != nil {
				Msg = "验证失败!<!--Err:UserId错误-->"
			} else if Key != models.MD5(this.User.Email+this.User.UserName+this.User.Password) {
				Msg = "验证失败!<!--Err:验证信息有误-->"
			}
		}
	}
	return Msg
}

func (this *UserController) VerifyEmail() {
	Result := this.Verify()
	if this.IsPost() {
		Result = "验证失败!<!--Err:非GET请求-->"
	} else if Result == "Success" {
		if this.User.EmailVerify != 0 {
			Result = "验证失败!<!--Err:重复验证-->"
		} else {
			this.User.EmailVerify = 1
			this.User.Update("EmailVerify")
		}
	}
	this.Data["Msg"] = Result
	this.Show("VerifyEmail")
}

func (this *UserController) EmailIsVerified() bool {
	return (this.User.EmailVerify == 1)
}

func (this *UserController) Modify() {
	this.IsLogin()
	action := this.Input().Get("a")
	if action == "passwd" {
		this.ModifyPassword()
	} else if action == "img" {
		this.ModifyHeadImage()
	} else {
		this.ModifyInfo()
	}
}

func (this *UserController) ModifyInfo() {
	this.IsLogin()
	User := make(map[string]interface{})
	Errmsg := make(map[string]interface{})
	if this.IsPost() {
		QQ := strings.TrimSpace(this.GetString("QQ"))
		Phone := strings.TrimSpace(this.GetString("Phone"))
		Email := strings.TrimSpace(this.GetString("Email"))
		School := strings.TrimSpace(this.GetString("School"))
		StuId := strings.TrimSpace(this.GetString("StuId"))

		User["UserName"] = this.User.UserName
		User["Email"] = Email
		User["QQ"] = QQ
		User["School"] = School
		User["StuId"] = StuId
		User["Phone"] = Phone

		valid := validation.Validation{}

		if v := valid.Required(Email, "Email"); !v.Ok {
			Errmsg["Email"] = "请输入Email地址"
		} else if v := valid.Email(Email, "Email"); !v.Ok || MailCheck(Email) {
			Errmsg["Email"] = "Email地址无效"
		} else if this.User.Email != Email && this.User.MailExist(Email) {
			Errmsg["Email"] = "邮箱已被使用"
		}

		if len(QQ) != 0 {
			if v := valid.Numeric(QQ, "QQ"); !v.Ok {
				Errmsg["QQ"] = "请输入数字"
			} else if v := valid.MinSize(QQ, 5, "QQ"); !v.Ok {
				Errmsg["QQ"] = "长度不合法"
			} else if v := valid.MaxSize(QQ, 11, "QQ"); !v.Ok {
				Errmsg["QQ"] = "长度不合法"
			} else {
				this.User.QQ = QQ
			}
		}

		if len(Phone) != 0 {
			if v := valid.Phone(Phone, "Phone"); !v.Ok {
				Errmsg["Phone"] = "电话号码不合法"
			} else {
				this.User.Phone = Phone
			}
		}

		if len(School) != 0 {
			if v := valid.MaxSize(School, 64, "School"); !v.Ok {
				Errmsg["School"] = "长度不合法"
			} else {
				this.User.School = School
				if len(StuId) != 0 {
					if v := valid.MinSize(StuId, 5, "StuId"); !v.Ok {
						Errmsg["StuId"] = "长度不合法"
					} else if v := valid.MaxSize(StuId, 16, "StuId"); !v.Ok {
						Errmsg["StuId"] = "长度不合法"
					} else {
						this.User.StuId = StuId
					}
				}
			}
		}

		if len(Errmsg) == 0 {
			if this.User.Email != Email {
				this.User.EmailVerify = 0
				this.User.Email = Email
			}
			if err := this.User.Update(); err != nil {
				Errmsg["Other"] = GetErrCode(err)
			} else if !this.EmailIsVerified() {
				if err := SendVerifyMail(&this.User); err != nil {
					Errmsg["Other"] = GetErrCode(err)
				}
			}

			if len(Errmsg) == 0 {
				//this.SaveSession()
				this.Redirect("/user/modify/", 302)
			}
		}
	} else {
		User["UserName"] = this.User.UserName
		User["Email"] = this.User.Email
		User["QQ"] = this.User.QQ
		User["School"] = this.User.School
		User["StuId"] = this.User.StuId
		User["Phone"] = this.User.Phone
	}
	this.Data["User"] = User
	this.Data["Errmsg"] = Errmsg
	this.Show("ModifyInfo")
}

func (this *UserController) ModifyPassword() {
	this.IsLogin()
	User := make(map[string]interface{})
	Errmsg := make(map[string]interface{})
	User["UserName"] = this.User.UserName
	if this.IsPost() {
		OldPassword := strings.TrimSpace(this.GetString("OldPassword"))
		Password := strings.TrimSpace(this.GetString("Password"))
		ConfirmPassword := strings.TrimSpace(this.GetString("ConfirmPassword"))

		valid := validation.Validation{}

		if v := valid.Required(OldPassword, "OldPassword"); !v.Ok {
			Errmsg["OldPassword"] = "请输入原密码"
		} else if v := valid.MaxSize(OldPassword, 18, "OldPassword"); !v.Ok {
			Errmsg["OldPassword"] = "密码长度不能大于18个字符"
		} else if v := valid.MinSize(OldPassword, 6, "OldPassword"); !v.Ok {
			Errmsg["OldPassword"] = "密码长度不能小于6个字符"
		} else if this.User.Password != models.MD5(OldPassword) {
			Errmsg["OldPassword"] = "原密码输入错误!"
		}

		if v := valid.Required(Password, "Password"); !v.Ok {
			Errmsg["Password"] = "请输入新密码"
		} else if v := valid.MaxSize(Password, 18, "Password"); !v.Ok {
			Errmsg["Password"] = "密码长度不能大于18个字符"
		} else if v := valid.MinSize(Password, 6, "Password"); !v.Ok {
			Errmsg["Password"] = "密码长度不能小于6个字符"
		}

		if v := valid.Required(ConfirmPassword, "ConfirmPassword"); !v.Ok {
			Errmsg["ConfirmPassword"] = "请确认您输入密码"
		} else if Password != ConfirmPassword {
			Errmsg["ConfirmPassword"] = "两次输入的密码不一致"
		}

		if len(Errmsg) == 0 {
			this.User.Password = models.MD5(Password)
			if err := this.User.Update("Password"); err != nil {
				Errmsg["Other"] = GetErrCode(err)
			}
			if len(Errmsg) == 0 {
				//this.SaveSession()
				this.Redirect("/user/", 302)
			}
		}
	}
	this.Data["User"] = User
	this.Data["Errmsg"] = Errmsg
	this.Show("ModifyPassword")
}

func (this *UserController) ModifyHeadImage() {
	this.IsLogin()
}
func (this *UserController) ForgotPassword() {
	//this.IsLogin()
	User := make(map[string]interface{})
	Errmsg := make(map[string]interface{})
	if this.IsPost() {
		UserName := strings.TrimSpace(this.GetString("UserName"))
		Email := strings.TrimSpace(this.GetString("Email"))

		User["UserName"] = UserName
		User["Email"] = Email

		valid := validation.Validation{}

		if v := valid.Required(UserName, "UserName"); !v.Ok {
			Errmsg["UserName"] = "请输入用户名"
		} else if v := valid.MaxSize(UserName, 16, "UserName"); !v.Ok {
			Errmsg["UserName"] = "用户名长度不能大于16个字符"
		} else if v := valid.MinSize(UserName, 3, "UserName"); !v.Ok {
			Errmsg["UserName"] = "用户名长度不能小于3个字符"
		} else if !this.User.Exist(UserName) {
			Errmsg["UserName"] = "用户名不存在"
		}

		if v := valid.Required(Email, "Email"); !v.Ok {
			Errmsg["Email"] = "请输入Email地址"
		} else if v := valid.Email(Email, "Email"); !v.Ok || MailCheck(Email) {
			Errmsg["Email"] = "Email地址无效"
		} else if !this.User.MailExist(Email) {
			Errmsg["Email"] = "邮箱不存在"
		}

		if len(Errmsg) == 0 {
			this.User.UserName = UserName
			if err := this.User.Read("UserName"); err != nil {
				Errmsg["Other"] = GetErrCode(err)
			} else if this.User.Email != Email {
				Errmsg["Other"] = "用户名或邮箱输入错误"
			} else if !this.EmailIsVerified() {
				Errmsg["UserName"] = "用户邮箱未验证"
			}

			if len(Errmsg) == 0 {
				if err := SendResetMail(&this.User); err != nil {
					Errmsg["Other"] = GetErrCode(err)
				} else {
					//跳转到提示已发送界面
				}
			}
		}
	}
	this.Data["User"] = User
	this.Data["Errmsg"] = Errmsg
	this.Show("ForgotPassword")
}

func (this *UserController) VerifyReset() {
	Result := this.Verify()
	User := make(map[string]interface{})
	Errmsg := make(map[string]interface{})
	if Result == "Success" {
		if this.IsPost() {
			UserName := strings.TrimSpace(this.GetString("UserName"))
			Password := strings.TrimSpace(this.GetString("Password"))
			ConfirmPassword := strings.TrimSpace(this.GetString("ConfirmPassword"))
			Email := strings.TrimSpace(this.GetString("Email"))

			User["UserName"] = UserName
			User["Password"] = Password
			User["ConfirmPassword"] = ConfirmPassword
			User["Email"] = Email

			valid := validation.Validation{}

			if v := valid.Required(UserName, "UserName"); !v.Ok {
				Errmsg["UserName"] = "请输入用户名"
			} else if v := valid.MaxSize(UserName, 16, "UserName"); !v.Ok {
				Errmsg["UserName"] = "用户名长度不能大于16个字符"
			} else if v := valid.MinSize(UserName, 3, "UserName"); !v.Ok {
				Errmsg["UserName"] = "用户名长度不能小于3个字符"
			} else if this.User.UserName != UserName {
				Errmsg["UserName"] = "用户名输入错误"
			}

			if v := valid.Required(Password, "Password"); !v.Ok {
				Errmsg["Password"] = "请输入密码"
			} else if v := valid.MaxSize(Password, 18, "Password"); !v.Ok {
				Errmsg["Password"] = "密码长度不能大于18个字符"
			} else if v := valid.MinSize(Password, 6, "Password"); !v.Ok {
				Errmsg["Password"] = "密码长度不能小于6个字符"
			}

			if v := valid.Required(ConfirmPassword, "ConfirmPassword"); !v.Ok {
				Errmsg["ConfirmPassword"] = "请确认您输入密码"
			} else if Password != ConfirmPassword {
				Errmsg["ConfirmPassword"] = "两次输入的密码不一致"
			}

			if v := valid.Required(Email, "Email"); !v.Ok {
				Errmsg["Email"] = "请输入Email地址"
			} else if v := valid.Email(Email, "Email"); !v.Ok || MailCheck(Email) {
				Errmsg["Email"] = "Email地址无效"
			} else if this.User.Email != Email {
				Errmsg["Email"] = "邮箱输入错误"
			}

			if len(Errmsg) == 0 {
				this.User.Password = models.MD5(Password)
				this.User.LoginTime = time.Now()
				if err := this.User.Update(); err != nil {
					Errmsg["Other"] = GetErrCode(err)
				}
				if len(Errmsg) == 0 {
					this.SaveSession()
					this.Redirect("/user/", 302)
				}
			}
		}
		this.Data["User"] = User
		this.Data["Errmsg"] = Errmsg
		this.Show("ResetPassword")
	} else {
		this.Data["Msg"] = Result
		this.Show("VerifyEmail")
		return
	}
}
