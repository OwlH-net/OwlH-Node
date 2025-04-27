package models

import (
	changecontrol "github.com/OwlH-net/OwlH-Node/changeControl"
	"github.com/OwlH-net/OwlH-Node/group"
	"github.com/astaxie/beego/logs"
)

func SyncSuricataGroupValues(data map[string]string, username string) (err error) {
	cc := data
	logs.Info("============")
	logs.Info("GROUP - SyncSuricataGroupValues")
	for key := range cc {
		logs.Info(key + " -> " + cc[key])
	}
	delete(data, "action")
	delete(data, "controller")
	delete(data, "router")

	logs.Info("Sync Suricata group values")
	err = group.SyncSuricataGroupValues(data)

	if err != nil {
		cc["actionStatus"] = "error"
		cc["errorDescription"] = err.Error()
	} else {
		cc["actionStatus"] = "success"
	}
	cc["username"] = username
	cc["actionDescription"] = "Sync Suricata Group Values"

	changecontrol.InsertChangeControl(cc)
	return err
}

func SuricataGroupService(data map[string]string, username string) (err error) {
	cc := data
	logs.Info("============")
	logs.Info("GROUP - SuricataGroupService")
	for key := range cc {
		logs.Info(key + " -> " + cc[key])
	}
	delete(data, "action")
	delete(data, "controller")
	delete(data, "router")

	logs.Info(data["action"] + " Suricata group values")
	err = group.SuricataGroupService(data)

	if err != nil {
		cc["actionStatus"] = "error"
		cc["errorDescription"] = err.Error()
	} else {
		cc["actionStatus"] = "success"
	}
	cc["username"] = username
	cc["actionDescription"] = data["action"] + " Suricata Group Values"

	changecontrol.InsertChangeControl(cc)
	return err
}

func SyncGroupRulesetToNode(data map[string][]byte, username string) (err error) {
	// cc := data
	// logs.Info("============")
	// logs.Info("GROUP - SyncGroupRulesetToNode")
	// for key :=range cc {
	//     logs.Info(key +" -> "+ cc[key])
	// }
	// delete(data,"action")
	// delete(data,"controller")
	// delete(data,"router")

	// logs.Info(data["action"]+" Group ruleset")
	err = group.SyncGroupRulesetToNode(data)

	// if err!=nil {
	//     cc["actionStatus"] = "error"
	//     cc["errorDescription"] = err.Error()
	// }else{
	//     cc["actionStatus"] = "success"
	// }
	// cc["username"] = username
	// cc["actionDescription"] = data["action"]+" Group ruleset"

	// changecontrol.InsertChangeControl(cc)
	return err
}
