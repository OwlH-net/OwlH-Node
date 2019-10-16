package ping

import (
    "github.com/astaxie/beego/logs"
	"os"
	"errors"
	"os/exec"
	"owlhnode/utils"
	"owlhnode/plugin"
	"owlhnode/database"
)

func PingService()(err error) {
	stapCollector := map[string]map[string]string{}
    stapCollector["service"] = map[string]string{}
    stapCollector["service"]["dstPath"] = ""
    stapCollector["service"]["file"] = ""
	stapCollector,err = utils.GetConf(stapCollector)
	if err != nil {logs.Error("ping/PingService -- Error GetConf service data: "+err.Error()); return err}
	dstPath := stapCollector["service"]["dstPath"]
	file := stapCollector["service"]["file"]

	if _, err := os.Stat(dstPath+file); os.IsNotExist(err) {
		return errors.New("Service don't exists")
	}else{
		logs.Info("OwlHnode service already exists")
		return nil
	}
}

func DeployService()(err error) {
	stapCollector := map[string]map[string]string{}
    stapCollector["service"] = map[string]string{}
    stapCollector["service"]["dstPath"] = ""
    stapCollector["service"]["file"] = ""
    stapCollector["service"]["origPath"] = ""
    stapCollector["service"]["reload"] = ""
    stapCollector["service"]["enable"] = ""
	stapCollector,err = utils.GetConf(stapCollector)
	if err != nil {logs.Error("ping/PingService -- Error GetConf service data: "+err.Error()); return err}
	dstPath := stapCollector["service"]["dstPath"]
	file := stapCollector["service"]["file"]
	origPath := stapCollector["service"]["origPath"]
	reload := stapCollector["service"]["reload"]
	enable := stapCollector["service"]["enable"]

	if _, err := os.Stat(dstPath+file); os.IsNotExist(err) {
		// //copy file
		err = utils.CopyFile(dstPath, origPath, file, 0)
		if err != nil {logs.Error("ping/Copy Error Copying file: "+err.Error()); return err}
	
		// //exec reload
		_,err = exec.Command("bash", "-c", reload).Output()
		if err != nil{logs.Error("utils.PingService Error reload service: "+err.Error()); return err}

		// //exec enable
		_,err = exec.Command("bash", "-c", enable).Output()
		if err != nil{logs.Error("utils.PingService Error enabling service: "+err.Error()); return err}

		// //return nil
		logs.Info("OwlHnode service deployed successfully!")
		return nil
	}else{
		logs.Info("OwlHnode service already exists")
		return nil
	}
}

func GetMainconfData()(data map[string]map[string]string, err error) {
	main,err := ndb.GetMainconfData(); if err != nil {logs.Error("ping/GetMainconfData error getting GetMainconfData values: "+err.Error()); return nil, err}
	if main["suricata"]["status"] == "" {
		err = ndb.InsertGetMainconfData("suricata","status","disabled"); if err != nil {logs.Error("ping/GetMainconfData error creating Suricata main data: "+err.Error()); return nil, err}
	}
	if main["zeek"]["status"] == "" {
 		err = ndb.InsertGetMainconfData("zeek","status","disabled"); if err != nil {logs.Error("ping/GetMainconfData error creating Zeek main data: "+err.Error()); return nil, err}
	}
	
	main,err = ndb.GetMainconfData(); if err != nil {logs.Error("ping/GetMainconfData error getting GetMainconfData values: "+err.Error()); return nil, err}

    return main,err
}

func PingPluginsNode() (data map[string]map[string]string ,err error) {
	allPlugin,err := ndb.GetPlugins()
	if err != nil {logs.Error("ping/GetMainconfData error getting GetPlugins values: "+err.Error()); return nil, err}

	for x := range allPlugin {
		if allPlugin[x]["status"] == "enabled" && allPlugin[x]["type"] == "suricata"{

			if _, err := os.Stat("/var/run/suricata/"+x+"-pidfile.pid"); os.IsNotExist(err) {		
				err = plugin.StopSuricataService(x, allPlugin[x]["status"])
    			if err != nil {logs.Error("ping/PingPluginsNode pidfile doesn't exist. Error stopping suricata for launch again: "+err.Error()); return nil,err}
				err = plugin.LaunchSuricataService(x, allPlugin[x]["interface"])
    			if err != nil {logs.Error("ping/PingPluginsNode pidfile doesn't exist. Error launching suricata again: "+err.Error()); return nil, err}
			}
		}
	}
	return allPlugin,err
}