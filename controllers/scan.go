package controllers

import (
	"github.com/astaxie/beego"
	"scanImage/clair"
	"github.com/astaxie/beego/logs"
	"scanImage/client"
	"time"
)

type ScanController struct {
	beego.Controller
}

// @Title Get
// @Description get layer
// @Success 200
// @router /:namespace/:repository/:tag [get]
func (s *ScanController) GetLayer() {

	beginTime := time.Now()

	ParamsMap := s.Ctx.Input.Params()

	namespace := ParamsMap[":namespace"]
	repository := ParamsMap[":repository"]
	tag := ParamsMap[":tag"]
	repo := namespace + "/" + repository

	scanedLayer, err := clair.GetClairHandler().ScanAndGetFeatures(repo, tag)
	if err != nil {
		logs.Error("扫描images失败:", err)
		return
	}

	elapsed := time.Since(beginTime)
	logs.Info("执行时间:",elapsed)

	s.Data["json"] = scanedLayer
	s.ServeJSON()
}

// @Title Get
// @Description get layer
// @Success 200
// @router / [get]
func (s *ScanController) GetLay() {

	beginTime := time.Now()


	m := make(map[string]string)
	m["library/tomcat"] = "9.0"
	m["library/golang"] = "1.7.3"
	m["library/centos"] = "7"
	m["library/openldap"] = "1.1.9"
	m["library/elasticsearch"] = "2.3.5"
	m["library/php"] = "7.1.13"


	for k, v := range m {
		go func(k, v string) {
			bTime := time.Now()
			_, err := clair.GetClairHandler().ScanAndGetFeatures(k, v)
			if err != nil {
				logs.Error("扫描images失败:", err)
				return
			}
			elap := time.Since(bTime)
			logs.Info(k + "的执行时间:",elap)
		}(k, v)
	}

	elapsss := time.Since(beginTime)
	logs.Info("总执行时间:",elapsss)


	s.Data["json"] = elapsss
	s.ServeJSON()
}

// @Title Get
// @Description get layer
// @Success 200
// @router /manifest [get]
func (s *ScanController) GetLayerManifest() {

	repository := "library/python"
	tag := "3.5"

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