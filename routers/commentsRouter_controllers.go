package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["github.com/zhanglianx111/clair-plus/controllers:ScanController"] = append(beego.GlobalControllerRouter["github.com/zhanglianx111/clair-plus/controllers:ScanController"],
		beego.ControllerComments{
			Method: "GetLay",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/zhanglianx111/clair-plus/controllers:ScanController"] = append(beego.GlobalControllerRouter["github.com/zhanglianx111/clair-plus/controllers:ScanController"],
		beego.ControllerComments{
			Method: "PostLayer",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/zhanglianx111/clair-plus/controllers:ScanController"] = append(beego.GlobalControllerRouter["github.com/zhanglianx111/clair-plus/controllers:ScanController"],
		beego.ControllerComments{
			Method: "GetLayer",
			Router: `/:namespace/:repository/:tag`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/zhanglianx111/clair-plus/controllers:ScanController"] = append(beego.GlobalControllerRouter["github.com/zhanglianx111/clair-plus/controllers:ScanController"],
		beego.ControllerComments{
			Method: "GetLayerManifest",
			Router: `/manifest`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/zhanglianx111/clair-plus/controllers:ScanController"] = append(beego.GlobalControllerRouter["github.com/zhanglianx111/clair-plus/controllers:ScanController"],
		beego.ControllerComments{
			Method: "Getos",
			Router: `/getos`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["github.com/zhanglianx111/clair-plus/controllers:ScanController"] = append(beego.GlobalControllerRouter["github.com/zhanglianx111/clair-plus/controllers:ScanController"],
		beego.ControllerComments{
			Method: "GetRepoTags",
			Router: `/tags/:namespace/:repository/:tag`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

}
