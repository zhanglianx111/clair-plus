package clair

import (
	"scanImage/models"
	"strings"
	"scanImage/client"
	"github.com/coreos/clair/api/v1"
	"github.com/astaxie/beego/logs"
	"errors"
)

type clairHandler struct {

}

type ClairInterface interface {
	ScanAndGetFeatures(repository string, tag string) (scanedLayer v1.LayerEnvelope, err error)
}

func GetClairHandler() ClairInterface {
	return &clairHandler{}
}

func (c *clairHandler) ScanAndGetFeatures(repository string, tag string) (scanedLayer v1.LayerEnvelope, err error) {

	//获取token
	token, err := client.GetClient().GetToken(repository)
	if err != nil {
		logs.Error("获取token失败:", err)
		return
	}
	//logs.Info("token:", token.Token)

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
	//logs.Info("manifest:", manifest)

	//通过manifest获取layers，扫描image并获取漏洞
	scanedLayer, err = scanImage(manifest, token.Token, repository)
	if err != nil {
		logs.Error("扫描images失败:", err)
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
			logs.Error("扫描" + repoName + "的layer失败:", err)
			return
		}
	}
	logs.Info(repoName + "调用clair post layer api 成功")

	//获取扫描后的漏洞
	imageDigestIndex := len(layers)
	scanedLayer, err = client.GetClient().GetLayerVulnerabilities(layers[imageDigestIndex - 1].Name)
	if err != nil {
		logs.Error("获取layer漏洞失败:", err)
		return
	}

	logs.Info("扫描image成功")
	return
}

func getLayers(manifestLayer []models.Layer) (layers []models.ClairLayer) {

	if manifestLayer == nil {
		logs.Error("manifest不可以为空")
		return
	}

	for index := 0 ; index < len(manifestLayer) ; index ++ {
		digest := manifestLayer[index].Digest
		name := strings.TrimPrefix(digest, "sha256:")
		var parentName string
		if index == 0 {
			parentName = ""
		}else {
			parentName = strings.TrimPrefix(manifestLayer[index - 1].Digest, "sha256:")
		}

		clairLayer := models.ClairLayer{
			Name: name,
			Digest: digest,
			ParentName: parentName,
		}

		layers = append(layers, clairLayer)
	}

	logs.Info("manifest解析layers成功")
	return
}