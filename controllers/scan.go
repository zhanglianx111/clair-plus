package controllers

import (
	"github.com/astaxie/beego"
	"scanImage/client"
	"scanImage/clair"
	"github.com/astaxie/beego/logs"
)

type ScanController struct {
	beego.Controller
}

// @Title Get
// @Description get layer
// @Success 200
// @router / [get]
func (s *ScanController) GetLayer() {

	repository := "library/mysql"
	tag := "5.5"

	//获取token
	token, err := client.GetClient().GetToken(repository)
	if err != nil {
		logs.Error("获取token失败:", err)
		return
	}

	//调用harbor api，拿到manifest
	manifest, err := client.GetClient().GetManifest(repository, tag)
	if err != nil {
		logs.Error("获取manifest失败:", err)
		return
	}
	if manifest.Manifest.MediaType == "" {
		logs.Error("manifest为空")
		return
	}

	//通过manifest获取layers，扫描image并获取漏洞
	scanedLayer, err := clair.GetClairHandler().ScanImage(manifest, token.Token, repository)
	if err != nil {
		logs.Error("扫描images失败:", err)
		return
	}

	s.Data["json"] = scanedLayer
	s.ServeJSON()

}

// @Title post
// @Description post
// @Success 200
// @Failure 403
// @router / [post]
func (s *ScanController) PostLayer() {

}