package controllers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"fmt"
)

type ScanController struct {
	beego.Controller
}

// @Title Get
// @Description get layer
// @Success 200
// @router / [get]
func (s *ScanController) GetLayer() {
	req := httplib.Get("http://beego.me/")

	str, err := req.String()

	if err != nil {
		fmt.Println(err)
	}

	fmt.Println(str)

	fmt.Println("==============")

	fmt.Println(req.Response())

}

// @Title post
// @Description post
// @Success 200
// @Failure 403
// @router / [post]
func (s *ScanController) PostLayer() {

}