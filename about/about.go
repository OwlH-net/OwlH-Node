package about

import (
	"github.com/astaxie/beego/logs"
)

type AboutST struct {
}

var ab AboutST

var version = "OwlH Node 01082020.1900"

func about() (aboutStruct AboutST) {
	logs.Info("About -> get node details")
	logs.Info("Node Version -> %s", version)
	logs.Info("Node Name ->")
	logs.Info("Node IP ->")
	//suricataVersion, versionError := suricata.SuricataVersion()
	//logs.Info("Suricata Version -> %+v", suricataVersion)
	return ab
}
