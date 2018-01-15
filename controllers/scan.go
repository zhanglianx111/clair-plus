package controllers

import (
	"github.com/astaxie/beego"
	"scanImage/clair"
	"github.com/astaxie/beego/logs"
	"scanImage/client"
)

type ScanController struct {
	beego.Controller
}

// @Title Get
// @Description get layer
// @Success 200
// @router / [get]
func (s *ScanController) GetLayer() {

	repository := "library/openldap"
	tag := "1.1.9"

	scanedLayer, err := clair.GetClairHandler().ScanAndGetFeatures(repository, tag)
	if err != nil {
		logs.Error("扫描images失败:", err)
		return
	}

	s.Data["json"] = scanedLayer
	s.ServeJSON()
}

// @Title Get
// @Description get layer
// @Success 200
// @router /manifest [get]
func (s *ScanController) GetLayerManifest() {

	repository := "library/openldap"
	tag := "1.1.9"

	manifest, err := client.GetClient().GetManifest(repository, tag)
	if err != nil {
		logs.Error("获取manifest失败:", err)
		return
	}

	s.Data["json"] = manifest
	s.ServeJSON()
}

// @Title post
// @Description post
// @Success 200
// @Failure 403
// @router / [post]
func (s *ScanController) PostLayer() {

}