package client

import (
	"github.com/astaxie/beego/httplib"
	"github.com/coreos/clair/api/v1"
	"scanImage/models"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/config"
)

var harborURL string
var clairURL string

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

	urlConf, err :=  config.NewConfig("ini", "conf/url.conf")
	if err != nil {
		logs.Error("解析url配置文件出错:", err)
	}

	harborURL = urlConf.String("harborURL")
	clairURL = urlConf.String("clairURL")
}

func (c *client) GetManifest(repoName string, tag string) (manifest models.ManifestObj, err error) {

	//调用harbor api获取image的manifest
	req := httplib.Get(buildHarborGetManifestURL(repoName, tag))
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
		Layer : &v1.Layer{
			Name : layer.Name,
			Path : buildHarborGetBlobURL(repository, layer.Digest),
			Headers	: header,
			ParentName : layer.ParentName,
			Format : "Docker",
		},
	}

	//调用clair扫描api
	req := httplib.Post(buildClairPostLayerURL())
	req.Header("Content-Type", "application/json")

	//将layers以json格式放到body中
	req.JSONBody(payload)

	_, err = req.String()
	if err != nil {
		return err
	}

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
	req := httplib.Get(buildHarborGetTokenURL(repository))
	req.SetBasicAuth("admin", "Harbor12345")

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

func buildHarborGetManifestURL(repository string, tag string) string {
	return harborURL + "/api/repositories/" + repository + "/tags/" + tag + "/manifest"
}

func buildHarborGetBlobURL(repository string, digest string) string {
	return harborURL + "/v2/" + repository + "/blobs/" + digest
}

func buildHarborGetTokenURL(repository string) string {
	return harborURL + "/service/token?account=admin&service=harbor-registry&scope=repository:"+ repository + ":pull"
}

func buildClairPostLayerURL() string {
	return clairURL + "/v1/layers"
}

func buildClairGetLayerFeaturesURL(layerName string) string {
	return clairURL + "/v1/layers/" + layerName + "?vulnerabilities"
}