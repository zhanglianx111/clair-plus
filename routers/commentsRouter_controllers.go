package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

	beego.GlobalControllerRouter["scanImage/controllers:ScanController"] = append(beego.GlobalControllerRouter["scanImage/controllers:ScanController"],
		beego.ControllerComments{
			Method: "GetLayer",
			Router: `/`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["scanImage/controllers:ScanController"] = append(beego.GlobalControllerRouter["scanImage/controllers:ScanController"],
		beego.ControllerComments{
			Method: "PostLayer",
			Router: `/`,
			AllowHTTPMethods: []string{"post"},
			MethodParams: param.Make(),
			Params: nil})

	beego.GlobalControllerRouter["scanImage/controllers:ScanController"] = append(beego.GlobalControllerRouter["scanImage/controllers:ScanController"],
		beego.ControllerComments{
			Method: "GetLayerManifest",
			Router: `/manifest`,
			AllowHTTPMethods: []string{"get"},
			MethodParams: param.Make(),
			Params: nil})

}
