// @APIVersion 0.1.0
// @Title OwlH Node API
// @Description OwlH node API
// @Contact support@owlh.net
// @TermsOfServiceUrl http://www.owlh.net
// @License Apache 2.0
// @LicenseUrl http://www.apache.org/licenses/LICENSE-2.0.html
package routers

import (
    "github.com/astaxie/beego"
    "owlhnode/controllers"
)

func init() {
    ns := beego.NewNamespace("/node",
        beego.NSNamespace("/suricata",
            beego.NSInclude(
                &controllers.SuricataController{},
            ),
        ),
        beego.NSNamespace("/zeek",
            beego.NSInclude(
                &controllers.ZeekController{},
            ),
        ),
        beego.NSNamespace("/hwaddmng",
            beego.NSInclude(
                &controllers.HwaddmngController{},
            ),
        ),
        beego.NSNamespace("/group",
            beego.NSInclude(
                &controllers.GroupController{},
            ),
        ),
        beego.NSNamespace("/wazuh",
            beego.NSInclude(
                &controllers.WazuhController{},
            ),
        ),
        beego.NSNamespace("/file",
            beego.NSInclude(
                &controllers.FileController{},
            ),
        ),
        beego.NSNamespace("/stap",
            beego.NSInclude(
                &controllers.StapController{},
            ),
        ),
        beego.NSNamespace("/ping",
            beego.NSInclude(
                &controllers.PingController{},
            ),
        ),
        beego.NSNamespace("/collector",
            beego.NSInclude(
                &controllers.CollectorController{},
            ),
        ),
        beego.NSNamespace("/ports",
            beego.NSInclude(
                &controllers.PortsController{},
            ),
        ),
        beego.NSNamespace("/analyzer",
            beego.NSInclude(
                &controllers.AnalyzerController{},
            ),
        ),
        beego.NSNamespace("/deploy",
            beego.NSInclude(
                &controllers.DeployController{},
            ),
        ),
        beego.NSNamespace("/dataflow",
            beego.NSInclude(
                &controllers.DataflowController{},
            ),
        ),
        beego.NSNamespace("/net",
            beego.NSInclude(
                &controllers.NetController{},
            ),
        ),
        beego.NSNamespace("/monitor",
            beego.NSInclude(
                &controllers.MonitorController{},
            ),
        ),
        beego.NSNamespace("/plugin",
            beego.NSInclude(
                &controllers.PluginController{},
            ),
        ),
        beego.NSNamespace("/changecontrol",
            beego.NSInclude(
                &controllers.ChangecontrolController{},
            ),
        ),
        beego.NSNamespace("/incidents",
            beego.NSInclude(
                &controllers.IncidentslController{},
            ),
        ),
        beego.NSNamespace("/autentication",
            beego.NSInclude(
                &controllers.AutenticationController{},
            ),
        ),
        beego.NSNamespace("/authentication",
            beego.NSInclude(
                &controllers.AutenticationController{},
            ),
        ),
    )
    beego.AddNamespace(ns)
}
