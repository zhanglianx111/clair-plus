package client

import (
	"encoding/json"
	"errors"
	"strings"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"github.com/astaxie/beego/logs"

	"github.com/coreos/clair/api/v1"
	"github.com/zhanglianx111/clair-plus/models"
	"time"
)

var harborURL string
var clairURL string
var checkCycle int64
var harborVersion float64

type client struct {
}

type ClientInterface interface {
	GetManifest(repoName string, tag string) (manifest models.ManifestObj, err error)
	ScanLayer(layer models.ClairLayer, repository string, token string) (err error)
	GetLayerVulnerabilities(layerName string) (scanedLayer v1.LayerEnvelope, err error)
	GetToken(repository string) (token models.Token, err error)
}

func GetClient() ClientInterface {
	return &client{}
}

func init() {

	harborURL = beego.AppConfig.String("harborURL")
	clairURL = beego.AppConfig.String("clairURL")
	checkCycle = beego.AppConfig.DefaultInt64("checkCycle", 2)
	harborVersion = beego.AppConfig.DefaultFloat("harborVersion", 0.4)

	//周期性验证harbor与clair的健康状态
	go func() {
		ticker := time.NewTicker(time.Minute * (time.Duration(checkCycle)))

		for range ticker.C {
			go checkHarborHealthy()
			go checkClairHealthy()
		}
	}()
}

func (c *client) GetManifest(repoName string, tag string) (manifest models.ManifestObj, err error) {

	//调用harbor api获取image的manifest
	var req *httplib.BeegoHTTPRequest

	if harborVersion == 0.4 {
		req = httplib.Get(buildOldHarborGetManifestURL(repoName, tag))
	} else if harborVersion == 1.2 {
		req = httplib.Get(buildHarborGetManifestURL(repoName, tag))
	} else {
		logs.Error("Harbor版本不存在")
		return
	}

	req.Header("Accept", " application/vnd.docker.distribution.manifest.v2+json")

	resp, err := req.String()
	if err != nil {
		return
	}
	logs.Info("获取" + repoName + "的manifest，成功")

	err = json.Unmarshal([]byte(resp), &manifest)
	if err != nil {
		return
	}
	return
}

func (c *client) ScanLayer(layer models.ClairLayer, repository string, token string) (err error) {

	//构建clair访问harbor registry的token
	header := make(map[string]string)
	header["Authorization"] = "Bearer " + token

	//构建clair官方的layer数据结构
	payload := v1.LayerEnvelope{
		Layer: &v1.Layer{
			Name:       layer.Name,
			Path:       buildHarborGetBlobURL(repository, layer.Digest),
			Headers:    header,
			ParentName: layer.ParentName,
			Format:     "Docker",
		},
	}

	//调用clair扫描api
	req := httplib.Post(buildClairPostLayerURL())
	req.Header("Content-Type", "application/json")

	//将layers以json格式放到body中
	req.JSONBody(payload)

	str, err := req.String()
	if err != nil {
		return err
	}
	if strings.Contains(str, "Error") {
		return errors.New(str)
	}

	//logs.Info("调用post扫描结果:",str)

	return err
}

func (c *client) GetLayerVulnerabilities(layerName string) (scanedLayer v1.LayerEnvelope, err error) {

	req := httplib.Get(buildClairGetLayerFeaturesURL(layerName))
	resp, err := req.String()
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &scanedLayer)
	if err != nil {
		return
	}

	logs.Info("调用clair get layer api 成功")
	return
}

func (c *client) GetToken(repository string) (token models.Token, err error) {

	//调用harbor api获取token
	var req *httplib.BeegoHTTPRequest

	if harborVersion == 0.4 {
		req = httplib.Get(buildOldHarborGetTokenURL(repository))
	} else if harborVersion == 1.2 {
		req = httplib.Get(buildOldHarborGetTokenURL(repository))
	} else {
		logs.Error("Harbor版本不存在")
		return
	}

	req.SetBasicAuth("admin", "12345")

	resp, err := req.String()
	if err != nil {
		return
	}

	err = json.Unmarshal([]byte(resp), &token)
	if err != nil {
		return
	}

	logs.Info("获取token 成功")
	return
}

func checkHarborHealthy() {

	req := httplib.Get(buildHarborGetSysInfoURL())
	_, err := req.String()

	if err != nil {
		logs.Error("Harbor状态异常:", err)
	} else {
		logs.Info("Harbor状态正常")
	}
}

func checkClairHealthy() {

	req := httplib.Get(buildClairGetNamespaceURL())
	_, err := req.String()

	if err != nil {
		logs.Error("Clair状态异常:", err)
	} else {
		logs.Info("Clair状态正常")
	}
}

func buildHarborGetManifestURL(repository string, tag string) string {
	return harborURL + "/api/repositories/" + repository + "/tags/" + tag + "/manifest"
}

func buildOldHarborGetManifestURL(repository string, tag string) string {
	return harborURL + "/api/repositories/manifests?repo_name=" + repository + "&tag=" + tag
}

func buildHarborGetBlobURL(repository string, digest string) string {
	return harborURL + "/v2/" + repository + "/blobs/" + digest
}

func buildHarborGetTokenURL(repository string) string {
	return harborURL + "/service/token?account=admin&service=harbor-registry&scope=repository:" + repository + ":pull"
}

func buildOldHarborGetTokenURL(repository string) string {
	return harborURL + "/service/token?account=admin&service=token-service&scope=repository:" + repository + ":pull"
}

func buildClairPostLayerURL() string {
	return clairURL + "/v1/layers"
}

func buildClairGetLayerFeaturesURL(layerName string) string {
	return clairURL + "/v1/layers/" + layerName + "?vulnerabilities"
}

func buildHarborGetSysInfoURL() string {
	return harborURL + "/api/systeminfo"
}

func buildClairGetNamespaceURL() string {
	return clairURL + "/v1/namespaces"
}
