package models

import (
	"github.com/astaxie/beego/orm"
	"strconv"
	"time"
)

type InviteCode struct {
	InviteCodeId int64     `orm:"pk;column(InviteCodeId)"`
	Code         string    `orm:"unique;size(32);column(Code)"`
	Used         int64     `orm:"column(Used)"`
	UserId       *User     `orm:"null;rel(fk);column(UserId)"`
	FromUserId   *User     `orm:"null;rel(fk);column(FromUserId)"`
	CreateTime   time.Time `orm:"auto_now_add;type(datetime);column(CreateTime)"`
	UsedTime     time.Time `orm:"null;type(datetime);column(UsedTime)"`
}

func (u *InviteCode) TableName() string {
	return TableName("InviteCode")
}

func (u *InviteCode) Insert() error {
	if _, err := orm.NewOrm().Insert(u); err != nil {
		return err
	}
	return nil
}

func (u *InviteCode) Read(fields ...string) error {
	if err := orm.NewOrm().Read(u, fields...); err != nil {
		return err
	}
	return nil
}

func (u *InviteCode) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(u, fields...); err != nil {
		return err
	}
	return nil
}
func (u *InviteCode) Delete() error {
	if _, err := orm.NewOrm().Delete(u); err != nil {
		return err
	}
	return nil
}

func (u *InviteCode) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(u)
}

func (u *InviteCode) VerifyInviteCode(Code string) bool {
	u.Code = Code
	err := orm.NewOrm().Read(u, "Code")
	if err == orm.ErrNoRows {
		return false
	} else if u.Used != 0 {
		return false
	}
	return true
}

func (u *InviteCode) GenerateInviteCode(user *User, num int) map[int]string {
	Code := make(map[int]string)
	for i := 0; i < num; i++ {
		invi := new(InviteCode)
		invi.Code = MD5(strconv.FormatInt(time.Now().UnixNano(), 10))
		invi.FromUserId = user
		invi.Insert()
		Code[i] = invi.Code
	}
	return Code
}
