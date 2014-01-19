package InkSec

import (
	"InkSec/models"
	"fmt"
	//"github.com/astaxie/beego"
)

type ChatController struct {
	BaseController
	Msg models.Message
}

func (this *ChatController) Prepare() {
	this.BaseController.Prepare()
	fmt.Println("ChatPrepare()")
}

func (this *ChatController) Main() {
	//this.IsLogin()
	this.Show("Index")
}
