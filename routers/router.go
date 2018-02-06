// @APIVersion 1.0.0
// @Title beego Test API
// @Description beego has a very cool tools to autogenerate documents for your API
// @Contact astaxie@gmail.com
// @TermsOfServiceUrl http://beego.me/
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
	"github.com/astaxie/beego"

	"github.com/zhanglianx111/clair-plus/controllers"
)

func init() {
	/*
		ns := beego.NewNamespace("/v1",
			beego.NSNamespace("/scan",
				beego.NSPost("/:namespace/:repository/:tag",
					&controllers.GetLayer),
			),
		)
	*/
	beego.Router("/v1/scan/:namespace/:repository/:tag", &controllers.ScanController{}, "post:PostLayer")
	beego.Router("/v1/scan/", &controllers.ScanController{}, "get:GetLay")
	beego.Router("/v1/scan/manifest", &controllers.ScanController{}, "get:GetLayerManifest")
	beego.Router("/v1/scan/tags/:namespace/:repository/:tag", &controllers.ScanController{}, "get:GetRepoTags")

	beego.Run(":8080")
}
