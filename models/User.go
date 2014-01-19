package models

import (
	"github.com/astaxie/beego/orm"
	"time"
)

type User struct {
	UserId      int64     `orm:"pk;column(UserId)"`
	UserName    string    `orm:"unique;size(32);column(UserName)"`
	Password    string    `orm:"size(32)"`
	LoginIP     string    `orm:"null;size(32);column(LoginIP)"`
	RegTime     time.Time `orm:"auto_now_add;type(datetime);column(RegTime)"`
	LoginTime   time.Time `orm:"auto_now_add;type(datetime);column(LoginTime)"`
	Email       string    `orm:"unique;size(64)"`
	EmailVerify int64     `orm:"column(EmailVerify)"`
	QQ          string    `orm:"null;size(16);column(QQ)"`
	Phone       string    `orm:"null;size(32);column(Phone)"`
	Rank        int64     `orm:"column(Rank)"`
	GoldCoin    int64     `orm:"column(GoldCoin)"`
	SignSum     int64     `orm:"column(SignSum)"`
	UserType    int64     `orm:"column(UserType)"`
	School      string    `orm:"null;size(128)"`
	StuId       string    `orm:"null;size(16);column(StuId)"`
	Enabled     int64     `orm:"column(Enabled)"`
}

func (u *User) TableName() string {
	return TableName("User")
}

func (u *User) Insert() error {
	if _, err := orm.NewOrm().Insert(u); err != nil {
		return err
	}
	u.Read("UserName")
	return nil
}

func (u *User) Read(fields ...string) error {
	if err := orm.NewOrm().Read(u, fields...); err != nil {
		return err
	}
	return nil
}

func (u *User) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(u, fields...); err != nil {
		return err
	}
	return nil
}
func (u *User) Delete() error {
	if _, err := orm.NewOrm().Delete(u); err != nil {
		return err
	}
	return nil
}

func (u *User) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(u)
}

func (u *User) Exist(UserName string) bool {
	user := new(User)
	user.UserName = UserName
	err := orm.NewOrm().Read(user, "UserName")
	if err == orm.ErrNoRows {
		return false
	}
	return true
}

func (u *User) MailExist(email string) bool {
	user := new(User)
	user.Email = email
	err := orm.NewOrm().Read(user, "Email")
	if err == orm.ErrNoRows {
		return false
	}
	return true
}
