package ping

import (
    "github.com/astaxie/beego/logs"
	"os"
	"errors"
	"os/exec"
	"strings"
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
	allPlugins,err := ndb.GetPlugins()
	if err != nil {logs.Error("ping/GetMainconfData error getting GetPlugins values: "+err.Error()); return nil, err}

	for x := range allPlugins {
		if allPlugins[x]["status"] == "enabled" && allPlugins[x]["type"] == "suricata"{

			if _, err := os.Stat("/var/run/suricata/"+x+"-pidfile.pid"); os.IsNotExist(err) {		
				err = plugin.StopSuricataService(x, allPlugins[x]["status"])
    			if err != nil {logs.Error("ping/PingPluginsNode pidfile doesn't exist. Error stopping suricata for launch again: "+err.Error()); return nil,err}
				err = plugin.LaunchSuricataService(x, allPlugins[x]["interface"])
    			if err != nil {logs.Error("ping/PingPluginsNode pidfile doesn't exist. Error launching suricata again: "+err.Error()); return nil, err}
			}
		}
	}

	//get suricata values that are not in the database
	var avoidUUIDS string
    for f := range allPlugins {
        if allPlugins[f]["type"] == "suricata" {
            avoidUUIDS = avoidUUIDS + "grep -v "+f+" | "
        }
    }

    com, err := exec.Command("bash","-c","ps -ef | grep suricata | "+avoidUUIDS+" grep -v grep | awk '{print $2}'").Output()
    if err != nil {logs.Error("PingPluginsNode error getting suricata shell launched: "+err.Error())}
    pidValue := strings.Split(string(com), "\n")
    for pid := range pidValue{
        if pidValue[pid] != "" {
            fullCommand, err := exec.Command("bash","-c","ps -ef | grep "+pidValue[pid]+" | grep -v grep").Output()
            if err != nil {logs.Error("PingPluginsNode error getting suricata shell full command: "+err.Error())}

            existsPid := false
            for f := range allPlugins {
                if allPlugins[f]["type"] == "suricata" && allPlugins[f]["pid"] == pidValue[pid] {
                    existsPid = true
                }
            }
            if !existsPid{
				uuid := utils.Generate()
				pluginNotControlled := make(map[string]string)
				pluginNotControlled["type"] = "suricata" 
				pluginNotControlled["pid"] = pidValue[pid]
				pluginNotControlled["command"] = string(fullCommand)
				allPlugins[uuid] = pluginNotControlled
            }
        }
    }

	return allPlugins,err
}

func UpdateNodeData(data map[string]map[string]string)(err error) {
	var action string
	currentData, err := ndb.GetNodeData()	

	if len(currentData) == 0{
		action = "insert"
	}else{
		action = "update"
	}

	for x,y := range data {
		for y,_ := range y {
			if action == "insert"{
				err = ndb.InsertNodeData(x,y,data[x][y])
				if err != nil { logs.Error("Error inserting node data: "+err.Error()); return err }
			}else if action == "update"{
				err = ndb.UpdateNodeData(x,y,data[x][y])
				if err != nil { logs.Error("Error updating node data: "+err.Error()); return err }
			}
		}
	}

	return nil
}