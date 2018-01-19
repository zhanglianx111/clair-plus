package controllers

import (
	"encoding/json"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"

	"github.com/zhanglianx111/clair-plus/clair"
	"github.com/zhanglianx111/clair-plus/client"
	"github.com/zhanglianx111/clair-plus/models"
	"github.com/zhanglianx111/clair-plus/mq"
)

type ScanController struct {
	beego.Controller
}

var Queue mq.Mqer

func init() {
	Queue = new(mq.RedisMq)
	err := Queue.NewMq("tasks1", "service", "tcp", "localhost:6379", 1)
	if err != nil {
		panic(err)
	}

	go Queue.NewConsumer("consumer")
}

// @Title Get
// @Description get layer
// @Success 200
// @router /:namespace/:repository/:tag [get]
func (s *ScanController) GetLayer() {
	var result string

	ParamsMap := s.Ctx.Input.Params()
	namespace := ParamsMap[":namespace"]
	repository := ParamsMap[":repository"]
	tag := ParamsMap[":tag"]
	repo := namespace + "/" + repository

	image := models.Image{
		Repo: repo,
		Tag:  tag,
	}

	r, err := json.Marshal(image)
	if err != nil {
		logs.Error("marshal failed: ", repo, tag)
	}
	// send request into mq
	b := Queue.SendBytes(r)
	if !b {
		result = "insert request into mq failed"
	} else {
		result = "ok"
	}
	/*
		scanedLayer, err := clair.GetClairHandler().ScanAndGetFeatures(repo, tag)
		if err != nil {
			logs.Error("扫描images失败:", err)
			return
		}

		elapsed := time.Since(beginTime)
		logs.Info("执行时间:", elapsed)

		s.Data["json"] = scanedLayer
	*/
	s.Data["json"] = result
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
