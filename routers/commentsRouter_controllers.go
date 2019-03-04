package routers

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context/param"
)

func init() {

    beego.GlobalControllerRouter["owlhnode/controllers:FileController"] = append(beego.GlobalControllerRouter["owlhnode/controllers:FileController"],
        beego.ControllerComments{
            Method: "SaveFile",
            Router: `/`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["owlhnode/controllers:FileController"] = append(beego.GlobalControllerRouter["owlhnode/controllers:FileController"],
        beego.ControllerComments{
            Method: "GetAllFiles",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["owlhnode/controllers:FileController"] = append(beego.GlobalControllerRouter["owlhnode/controllers:FileController"],
        beego.ControllerComments{
            Method: "SendFile",
            Router: `/:fileName`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["owlhnode/controllers:StapController"] = append(beego.GlobalControllerRouter["owlhnode/controllers:StapController"],
        beego.ControllerComments{
            Method: "AddServer",
            Router: `/`,
            AllowHTTPMethods: []string{"post"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["owlhnode/controllers:StapController"] = append(beego.GlobalControllerRouter["owlhnode/controllers:StapController"],
        beego.ControllerComments{
            Method: "GetAllServers",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["owlhnode/controllers:StapController"] = append(beego.GlobalControllerRouter["owlhnode/controllers:StapController"],
        beego.ControllerComments{
            Method: "GetServer",
            Router: `/server/:uuid`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["owlhnode/controllers:SuricataController"] = append(beego.GlobalControllerRouter["owlhnode/controllers:SuricataController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["owlhnode/controllers:SuricataController"] = append(beego.GlobalControllerRouter["owlhnode/controllers:SuricataController"],
        beego.ControllerComments{
            Method: "SetBPF",
            Router: `/bpf`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["owlhnode/controllers:SuricataController"] = append(beego.GlobalControllerRouter["owlhnode/controllers:SuricataController"],
        beego.ControllerComments{
            Method: "RetrieveFile",
            Router: `/retrieve`,
            AllowHTTPMethods: []string{"put"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["owlhnode/controllers:WazuhController"] = append(beego.GlobalControllerRouter["owlhnode/controllers:WazuhController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

    beego.GlobalControllerRouter["owlhnode/controllers:ZeekController"] = append(beego.GlobalControllerRouter["owlhnode/controllers:ZeekController"],
        beego.ControllerComments{
            Method: "Get",
            Router: `/`,
            AllowHTTPMethods: []string{"get"},
            MethodParams: param.Make(),
            Filters: nil,
            Params: nil})

}
