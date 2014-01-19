package InkSec

import (
//"fmt"
//"github.com/astaxie/beego"
)

type IndexController struct {
	BaseController
}

/*
func init() {
	CtrlName = "index"
	fmt.Println("IndexControllerInit()")
}
*/
func (this *IndexController) Get() {
	//this.Data["text"] = "index"
	this.Show("Index")
}

/*
func Template(Name string) string {
	return CtrlName + "/" + Name + ".html"
}
*/
