package models

import (
	"crypto/md5"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
)

var dbpfx string

func init() {
	dbhost := string(RsaDecrypt(beego.AppConfig.String("dbHost")))
	dbuser := string(RsaDecrypt(beego.AppConfig.String("dbUser")))
	dbpass := string(RsaDecrypt(beego.AppConfig.String("dbPass")))
	dbname := string(RsaDecrypt(beego.AppConfig.String("dbName")))
	dbpfx = string(RsaDecrypt(beego.AppConfig.String("dbPfx")))
	dbport := beego.AppConfig.String("dbPort")
	if dbport == "" {
		dbport = "3306"
	}
	dsn := dbuser + ":" + dbpass + "@tcp(" + dbhost + ":" + dbport + ")/" + dbname + "?charset=utf8"
	orm.RegisterDataBase("default", "mysql", dsn)
	orm.RegisterModel(new(User), new(InviteCode), new(Message))
	fmt.Println("ModelsInit()")
}

func MD5(buf string) string {
	hash := md5.New()
	hash.Write([]byte("InkSec#" + buf + "@CTF"))
	return fmt.Sprintf("%x", hash.Sum(nil))
}

func TableName(str string) string {
	return fmt.Sprintf("%s%s", dbpfx, str)
}

func checkErr(err error) {
	if err != nil {
		panic(err)
	}
}
