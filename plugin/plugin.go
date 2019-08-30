package plugin

import (
    "github.com/astaxie/beego/logs"	
	"owlhnode/database"
)

func ChangeServiceStatus(anode map[string]string)(err error) {
	err = ndb.UpdatePluginValue(anode["server"],anode["param"],anode["status"])
	if err != nil {logs.Error("plugin/ChangeServiceStatus error: "+err.Error()); return err}
    return err
}

func ChangeMainServiceStatus(anode map[string]string)(err error) {
	err = ndb.UpdateMainconfValue(anode["service"],anode["param"],anode["status"])
	if err != nil {logs.Error("plugin/ChangeMainServiceStatus error: "+err.Error()); return err}
    return err
}