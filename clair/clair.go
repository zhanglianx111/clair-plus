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
	ScanImage(manifest models.ManifestObj, token string, repoName string) (scanedLayer v1.LayerEnvelope, err error)
}

func GetClairHandler() ClairInterface {
	return &clairHandler{}
}

func (c *clairHandler) ScanImage(manifest models.ManifestObj, token string, repoName string) (scanedLayer v1.LayerEnvelope, err error) {

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
			logs.Error("扫描layer失败:", err)
			return
		}
	}
	logs.Info("调用clair post layer api 成功")

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