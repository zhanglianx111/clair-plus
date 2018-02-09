package clair

import (
	"errors"
	"strings"

	"github.com/astaxie/beego/logs"
	"github.com/coreos/clair/api/v1"

	"github.com/astaxie/beego"
	"github.com/kr/text"
	"github.com/zhanglianx111/clair-plus/client"
	"github.com/zhanglianx111/clair-plus/models"
	"sync"
)

type ClairInterface interface {
	ScanAndGetFeatures(repository string, tag string) (scanedLayer v1.LayerEnvelope, err error)
	GetWebPortVulner(reposiroty string, tag string) (vulner models.Vulner, err error)
}

var clair *clairHandler
var once sync.Once

type clairHandler struct {
}

func GetClairHandler() ClairInterface {

	once.Do(func() {
		clair = &clairHandler{}
	})

	return clair
}

func (c *clairHandler) GetWebPortVulner(reposiroty string, tag string) (vulner models.Vulner, err error) {

	layer, err := c.ScanAndGetFeatures(reposiroty, tag)
	if err != nil {
		return
	}
	if layer.Layer == nil {
		return
	}

	var vulnerabilities = make([]models.VulnerabilityInfo, 0)
	for _, feature := range layer.Layer.Features {

		if len(feature.Vulnerabilities) > 0 {
			for _, vulnerability := range feature.Vulnerabilities {
				severity := vulnerability.Severity
				vulnerabilities = append(vulnerabilities, models.VulnerabilityInfo{vulnerability, feature, severity})
			}
		}
	}

	harborURL := strings.Split(beego.AppConfig.String("harborURL"), "//")[1]
	vulner.ImageName = harborURL + "/" + reposiroty + ":" + tag

	for _, vulnerabilityInfo := range vulnerabilities {

		vulnerability := vulnerabilityInfo.Vulnerability
		feature := vulnerabilityInfo.Feature
		v := models.V{}

		v.Name = vulnerability.Name
		v.Severity = vulnerabilityInfo.Severity

		if vulnerability.Description != "" {
			v.Description = text.Indent(text.Wrap(vulnerability.Description, 80), "\t")
		}

		v.Package = models.Package{
			Name:    feature.Name,
			Version: feature.Version,
		}

		if vulnerability.Link != "" {
			v.Link = vulnerability.Link
		}

		if vulnerability.FixedBy != "" {
			v.FixedByVersion = vulnerability.FixedBy
		}

		v.Layer = feature.AddedBy

		vulner.Vulners = append(vulner.Vulners, v)
	}

	return
}

func (c *clairHandler) ScanAndGetFeatures(repository string, tag string) (scanedLayer v1.LayerEnvelope, err error) {

	//判断repository是否存在
	isExit, err := client.GetClient().IsRepoTagExist(repository, tag)
	if err != nil {
		logs.Error("获取repository tags失败: ", err)
		return
	}
	if !isExit {
		logs.Error("repository:" + repository + ":" + tag + " 不存在")
		return
	}
	logs.Debug("repository: " + repository + ":" + tag + " 存在")

	//获取token
	token, err := client.GetClient().GetToken(repository)
	if err != nil {
		logs.Error("获取token失败: ", err)
		return
	}
	//logs.Info("token:", token.Token)

	//调用harbor api，拿到manifest
	manifest, err := client.GetClient().GetManifest(repository, tag)
	if err != nil {
		logs.Error("获取manifest失败: ", err)
		return
	}
	if manifest.Manifest.MediaType == "" {
		logs.Error("manifest为空")
		return
	}
	//logs.Info("manifest:", manifest)

	//通过manifest获取layers，扫描image并获取漏洞
	scanedLayer, err = scanImage(manifest, token.Token, repository)
	if err != nil {
		logs.Error("扫描images失败: ", err)
		return
	}

	return
}

func scanImage(manifest models.ManifestObj, token string, repoName string) (scanedLayer v1.LayerEnvelope, err error) {

	//获取layer
	layers := getLayers(manifest.Manifest.Layers)

	if layers == nil {
		err = errors.New("获取的layers为空")
		return
	}

	//逐层扫描layer
	for _, layer := range layers {

		err = client.GetClient().ScanLayer(layer, repoName, token)
		if err != nil {
			logs.Error("扫描"+repoName+"的layer失败: ", err)
			return
		}
	}
	logs.Debug(repoName + "调用clair post layer api 成功")

	//获取扫描后的漏洞
	imageDigestIndex := len(layers)
	scanedLayer, err = client.GetClient().GetLayerVulnerabilities(layers[imageDigestIndex-1].Name)
	if err != nil {
		logs.Error("获取layer漏洞失败: ", err)
		return
	}

	logs.Info("扫描image成功: ", repoName)
	return
}

func getLayers(manifestLayer []models.Layer) (layers []models.ClairLayer) {

	if manifestLayer == nil {
		logs.Error("manifest不可以为空")
		return
	}

	for index := 0; index < len(manifestLayer); index++ {
		digest := manifestLayer[index].Digest
		name := strings.TrimPrefix(digest, "sha256:")
		var parentName string
		if index == 0 {
			parentName = ""
		} else {
			parentName = strings.TrimPrefix(manifestLayer[index-1].Digest, "sha256:")
		}

		clairLayer := models.ClairLayer{
			Name:       name,
			Digest:     digest,
			ParentName: parentName,
		}

		layers = append(layers, clairLayer)
	}

	logs.Debug("manifest解析layers成功")
	return
}
