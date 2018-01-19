package controllers

import (
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/zhanglianx111/clair-plus/clair"
	"github.com/zhanglianx111/clair-plus/client"
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
	logs.Info("执行时间:", elapsed)

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
	m["appstore/sonarqube"] = "latest"
	m["appstore/jetty"] = "latest"
	m["appstore/redis"] = "latest"
	m["appstore/httpd"] = "latest"
	m["appstore/gitlab-ce"] = "latest"
	m["appstore/rabbitmq"] = "3.6.6"
	m["appstore/postgres"] = "latest"
	m["appstore/etcd"] = "caas"
	m["appstore/redmine"] = "latest"
	m["appstore/wordpress"] = "latest"
	m["appstore/joomla"] = "latest"
	m["appstore/magento"] = "alexcheng"
	m["appstore/durpal"] = "latest"
	m["fanbc/redis"] = "1.0"
	m["chrju/etcd"] = "4.0"
	m["library/ldap"] = "1.0"
	m["bitnami/ghost"] = "1.14"
	m["library/centos"] = "7.2.1511"

	for k, v := range m {
		go func(k, v string) {
			bTime := time.Now()
			_, err := clair.GetClairHandler().ScanAndGetFeatures(k, v)
			if err != nil {
				logs.Error("扫描images失败:", err)
				return
			}
			elap := time.Since(bTime)
			logs.Info(k+"的执行时间:", elap)
		}(k, v)
	}

	elapsss := time.Since(beginTime)
	logs.Info("总执行时间:", elapsss)

	s.Data["json"] = elapsss
	s.ServeJSON()
}

// @Title Get
// @Description get layer
// @Success 200
// @router /manifest [get]
func (s *ScanController) GetLayerManifest() {

	repository := "chrju/etcd"
	tag := "4.0"

	manifest, err := client.GetClient().GetManifest(repository, tag)
	if err != nil {
		logs.Error("获取manifest失败:", err)
		return
	}

	s.Data["json"] = manifest
	s.ServeJSON()
}

// @Title Get
// @Description get layer
// @Success 200
// @router /tags/:namespace/:repository/:tag [get]
func (s *ScanController) GetRepoTags() {

	ParamsMap := s.Ctx.Input.Params()

	namespace := ParamsMap[":namespace"]
	repository := ParamsMap[":repository"]
	tag := ParamsMap[":tag"]

	repo := namespace + "/" + repository

	isExist, err := client.GetClient().IsRepoTagExist(repo, tag)
	if err != nil {
		logs.Error("失败:", err)
		return
	}

	s.Data["json"] = isExist
	s.ServeJSON()
}

// @Title post
// @Description post
// @Success 200
// @Failure 403
// @router / [post]
func (s *ScanController) PostLayer() {

}