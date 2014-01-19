package models

import (
	"github.com/astaxie/beego/orm"
	//"strconv"
	"time"
)

type Message struct {
	MessageId int64     `orm:"pk;column(MessageId)"`
	FatherId  *User     `orm:"null;rel(fk);column(FatherId)"`
	UserId    *User     `orm:"null;rel(fk);column(UserId)"`
	Body      string    `orm:"size(512);column(Body)"`
	Time      time.Time `orm:"auto_now_add;type(datetime);column(Time)"`
	Enabled   int64     `orm:"null;type(datetime);column(Enabled)"`
}

func (u *Message) TableName() string {
	return TableName("Message")
}

func (u *Message) Insert() error {
	if _, err := orm.NewOrm().Insert(u); err != nil {
		return err
	}
	return nil
}

func (u *Message) Read(fields ...string) error {
	if err := orm.NewOrm().Read(u, fields...); err != nil {
		return err
	}
	return nil
}

func (u *Message) Update(fields ...string) error {
	if _, err := orm.NewOrm().Update(u, fields...); err != nil {
		return err
	}
	return nil
}
func (u *Message) Delete() error {
	if _, err := orm.NewOrm().Delete(u); err != nil {
		return err
	}
	return nil
}

func (u *Message) Query() orm.QuerySeter {
	return orm.NewOrm().QueryTable(u)
}
