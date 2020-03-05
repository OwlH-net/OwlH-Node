package ping

import (
    "github.com/astaxie/beego/logs"
    "os"
    "errors"
    "os/exec"
    "strings"
    "owlhnode/utils"
    "owlhnode/plugin"
    "owlhnode/zeek"
    "owlhnode/database"
)

type Zeek struct {
    Path        bool                `json:"path"`
    Rol         string              `json:"role"`
    Bin         bool                `json:"bin"`
    Action      string              `json:"action"`
    Running     []ZeekNodeStatus    `json:"running"`
    Mode        string              `json:"mode"`
    Managed     bool                `json:"managed"`
    Nodes       []ZeekNode          `json:"nodes"`
    Extra       map[string]string   `json:"extra"`
}

type ZeekKeys struct {
    Key         string              `json:"key"`
    Value       string              `json:"value"`
}

type ZeekNode struct {
    Name        string              `json:"name"`
    Host        string              `json:"host"`
    Status      string              `json:"status"`
    Type        string              `json:"type"`
    NInterface  string              `json:"interface"`
    Pid         string              `json:"pid"`
    Started     string              `json:"started"`
    Extra       []ZeekKeys          `json:"extra"`
}

type ZeekNodeStatus struct {
    Status      string              `json:"status"`
    Nodes       int                 `json:"nodes"`
}

func PingService()(err error) {  
    dstPath, err := utils.GetKeyValueString("service", "dstPath")
    if err != nil {logs.Error("ping/PingService -- Error getting service data: "+err.Error()); return err}
    file, err := utils.GetKeyValueString("service", "file")
    if err != nil {logs.Error("ping/PingService -- Error getting service data: "+err.Error()); return err}

    if _, err := os.Stat(dstPath+file); os.IsNotExist(err) {
        return errors.New("Service don't exists")
    }else{
        logs.Info("OwlHnode service already exists")
        return nil
    }
}

func DeployService()(err error) {
    dstPath, err := utils.GetKeyValueString("service", "dstPath")
    if err != nil {logs.Error("ping/DeployService -- Error getting deploy service data: "+err.Error()); return err}
    file, err := utils.GetKeyValueString("service", "file")
    if err != nil {logs.Error("ping/DeployService -- Error getting deploy service data: "+err.Error()); return err}
    origPath, err := utils.GetKeyValueString("service", "origPath")
    if err != nil {logs.Error("ping/DeployService -- Error getting deploy service data: "+err.Error()); return err}
    reload, err := utils.GetKeyValueString("service", "reload")
    if err != nil {logs.Error("ping/DeployService -- Error getting deploy service data: "+err.Error()); return err}
    enable, err := utils.GetKeyValueString("service", "enable")
    if err != nil {logs.Error("ping/DeployService -- Error getting deploy service data: "+err.Error()); return err}
    

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
    bck, err := utils.GetKeyValueString("suricata", "backup")
    if err != nil {logs.Error("ping/PingService -- Error getting suricata service data: "+err.Error()); return nil, err}

    allPlugins,err := ndb.GetPlugins()
    if err != nil {logs.Error("ping/GetMainconfData error getting GetPlugins values: "+err.Error()); return nil, err}

    for x := range allPlugins {
        if allPlugins[x]["status"] == "enabled" && allPlugins[x]["type"] == "suricata"{
            // change pid file name 
            if _, err := os.Stat(bck+x+"-pidfile.pid"); os.IsNotExist(err) {        
                err = plugin.StopSuricataService(x, allPlugins[x]["status"])
                if err != nil {logs.Error("ping/PingPluginsNode pidfile doesn't exist. Error stopping suricata for launch again: "+err.Error()); return nil,err}
                err = plugin.LaunchSuricataService(x, allPlugins[x]["interface"])
                if err != nil {logs.Error("ping/PingPluginsNode pidfile doesn't exist. Error launching suricata again: "+err.Error()); return nil, err}
            }
            //check if process is running even though database status is enabled
            pid, err := exec.Command("bash","-c","ps -ef | grep suricata | grep "+x+" | grep -v grep | awk '{print $2}'").Output()
            if err != nil {logs.Error("ping/PingPluginsNode Checking suricata PID: "+err.Error())}
            if strings.Split(string(pid), "\n")[0] == "" {
                allPlugins[x]["running"] = "false"
            }else{
                allPlugins[x]["running"] = "true"
            }
        }
        //check if process is running even though database status is enabled
        if allPlugins[x]["type"] == "zeek" && allPlugins[x]["pid"] != "none"{            
            zk := zeek.GetZeek()
            for node := range zk.Nodes {
                if zk.Nodes[node].Status != "running"{
                    allPlugins[x]["running"] = "false"
                }else{
                    allPlugins[x]["running"] = "true"
                }
            }
        }

        //check if process is running even though database status is enabled
        if (allPlugins[x]["type"] == "socket-pcap" || allPlugins[x]["type"] == "socket-network") && allPlugins[x]["pid"] != "none"{
            pid, err := exec.Command("bash","-c","ps -ef | grep socat | grep OPENSSL-LISTEN:"+allPlugins[x]["port"]+" | grep -v grep | awk '{print $2}'").Output()
            if err != nil {logs.Error("ping/PingPluginsNode Checking STAP PID: "+err.Error())}
            if strings.Split(string(pid), "\n")[0] == "" {
                allPlugins[x]["running"] = "false"
            }else{
                allPlugins[x]["running"] = "true"
            }
        }
        //check if process is running even though database status is enabled        
        if allPlugins[x]["type"] == "network-socket" && allPlugins[x]["pid"] != "none"{
            pid, err := exec.Command("bash","-c","ps -ef | grep OPENSSL:"+allPlugins[x]["collector"]+":"+allPlugins[x]["port"]+" | grep -v grep | awk '{print $2}'").Output()
            if err != nil {logs.Error("ping/PingPluginsNode Checking STAP network-socket PID: "+err.Error())}
            pids := strings.Split(string(pid), "\n")
            exists := false
            for q := range pids {
                if pids[q] == allPlugins[x]["pid"] { exists = true}
            }
            if !exists{
                allPlugins[x]["running"] = "false"
            }else{
                allPlugins[x]["running"] = "true"
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

    com, err := exec.Command("bash","-c","ps -ef | grep \"suricata \" | "+avoidUUIDS+" grep -v grep | awk '{print $2}'").Output()
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
    for x,y := range data {
        for y,_ := range y {
            err = ndb.UpdateNodeData(x,y,data[x][y])
            if err != nil { logs.Error("UpdateNodeData Error updating node data: "+err.Error()); return err }
        }
    }

    return nil
}

func SaveNodeInformation(anode map[string]map[string]string)(err error) {
    // nodeData, err := ndb.GetNodeData()
    if err != nil { logs.Error("SaveNodeInformation Error getting node data: "+err.Error()); return err }
    for x := range anode {
        err = ndb.DeleteNodeInformation()
        if err != nil { logs.Error("SaveNodeInformation Error updating node values: "+err.Error()); return err }
        err = ndb.InsertNodeData(x,"ip",anode[x]["ip"])
        if err != nil { logs.Error("SaveNodeInformation Error inserting node ip: "+err.Error()); return err }
        err = ndb.InsertNodeData(x,"name",anode[x]["name"])
        if err != nil { logs.Error("SaveNodeInformation Error inserting node name: "+err.Error()); return err }
        err = ndb.InsertNodeData(x,"port",anode[x]["port"])
        if err != nil { logs.Error("SaveNodeInformation Error inserting node port: "+err.Error()); return err }                     
        
    }

    return nil
}

func DeleteNode(masterID string)(err error) {
    //delete node information
    err = ndb.DeleteNodeInformation()
    if err != nil { logs.Error("DeleteNode Error deleting node data: "+err.Error()); return err }
    err = ndb.DeleteMastesInformation(masterID)
    if err != nil { logs.Error("DeleteNode Error deleting master data: "+err.Error()); return err }
    
    return err
}