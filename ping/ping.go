package ping

import (
    "errors"
    "github.com/astaxie/beego/logs"
    "os"
    "os/exec"
    "owlhnode/suricata"
    "owlhnode/utils"
    "strings"
    "strconv"
    // "owlhnode/zeek"
    "owlhnode/database"
)

type Zeek struct {
    Path    bool              `json:"path"`
    Rol     string            `json:"role"`
    Bin     bool              `json:"bin"`
    Action  string            `json:"action"`
    Running []ZeekNodeStatus  `json:"running"`
    Mode    string            `json:"mode"`
    Managed bool              `json:"managed"`
    Nodes   []ZeekNode        `json:"nodes"`
    Extra   map[string]string `json:"extra"`
}

type ZeekKeys struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}

type ZeekNode struct {
    Name       string     `json:"name"`
    Host       string     `json:"host"`
    Status     string     `json:"status"`
    Type       string     `json:"type"`
    NInterface string     `json:"interface"`
    Pid        string     `json:"pid"`
    Started    string     `json:"started"`
    Extra      []ZeekKeys `json:"extra"`
}

type ZeekNodeStatus struct {
    Status string `json:"status"`
    Nodes  int    `json:"nodes"`
}

func PingService() (err error) {
    dstPath, err := utils.GetKeyValueString("service", "dstPath")
    if err != nil {
        logs.Error("ping/PingService -- Error getting service data: " + err.Error())
        return err
    }
    file, err := utils.GetKeyValueString("service", "file")
    if err != nil {
        logs.Error("ping/PingService -- Error getting service data: " + err.Error())
        return err
    }

    if _, err := os.Stat(dstPath + file); os.IsNotExist(err) {
        return errors.New("Service don't exists")
    } else {
        logs.Info("OwlHnode service already exists")
        return nil
    }
}

func DeployService() (err error) {
    dstPath, err := utils.GetKeyValueString("service", "dstPath")
    if err != nil {
        logs.Error("ping/DeployService -- Error getting deploy service data: " + err.Error())
        return err
    }
    file, err := utils.GetKeyValueString("service", "file")
    if err != nil {
        logs.Error("ping/DeployService -- Error getting deploy service data: " + err.Error())
        return err
    }
    origPath, err := utils.GetKeyValueString("service", "origPath")
    if err != nil {
        logs.Error("ping/DeployService -- Error getting deploy service data: " + err.Error())
        return err
    }
    reload, err := utils.GetKeyValueString("service", "reload")
    if err != nil {
        logs.Error("ping/DeployService -- Error getting deploy service data: " + err.Error())
        return err
    }
    enable, err := utils.GetKeyValueString("service", "enable")
    if err != nil {
        logs.Error("ping/DeployService -- Error getting deploy service data: " + err.Error())
        return err
    }
    param, err := utils.GetKeyValueString("execute", "param")
    if err != nil {
        logs.Error("ping/DeployService Error getting data from main.conf")
        return err
    }
    command, err := utils.GetKeyValueString("execute", "command")
    if err != nil {
        logs.Error("ping/DeployService Error getting data from main.conf")
        return err
    }

    if _, err := os.Stat(dstPath + file); os.IsNotExist(err) {
        // //copy file
        err = utils.CopyFile(dstPath, origPath, file, 0)
        if err != nil {
            logs.Error("ping/Copy Error Copying file: " + err.Error())
            return err
        }

        // //exec reload
        _, err = exec.Command(command, param, reload).Output()
        if err != nil {
            logs.Error("utils.PingService Error reload service: " + err.Error())
            return err
        }

        // //exec enable
        _, err = exec.Command(command, param, enable).Output()
        if err != nil {
            logs.Error("utils.PingService Error enabling service: " + err.Error())
            return err
        }

        // //return nil
        logs.Info("OwlHnode service deployed successfully!")
        return nil
    } else {
        logs.Info("OwlHnode service already exists")
        return nil
    }
}

func GetMainconfData() (data map[string]map[string]string, err error) {
    main, err := ndb.GetMainconfData()
    if err != nil {
        logs.Error("ping/GetMainconfData error getting GetMainconfData values: " + err.Error())
        return nil, err
    }
    if main["suricata"]["status"] == "" {
        err = ndb.InsertGetMainconfData("suricata", "status", "disabled")
        if err != nil {
            logs.Error("ping/GetMainconfData error creating Suricata main data: " + err.Error())
            return nil, err
        }
    }
    if main["zeek"]["status"] == "" {
        err = ndb.InsertGetMainconfData("zeek", "status", "disabled")
        if err != nil {
            logs.Error("ping/GetMainconfData error creating Zeek main data: " + err.Error())
            return nil, err
        }
    }

    main, err = ndb.GetMainconfData()
    if err != nil {
        logs.Error("ping/GetMainconfData error getting GetMainconfData values: " + err.Error())
        return nil, err
    }

    //Ping Suricata for check if is installed
    isInstalled,err := suricata.Installed()
    if isInstalled["path"] == false || isInstalled["bin"] == false {
        main["suricata"]["status"] = ""
    }

    //******************************
    //Check if Zeek is installed is not necessary
    //Is checked by other way
    //******************************

    return main, err
}

func PingPluginsNode() (data map[string]map[string]string, err error) {
    bck, err := utils.GetKeyValueString("suricata", "backup")
    if err != nil {
        logs.Error("ping/PingService -- Error getting suricata service data: " + err.Error())
        return nil, err
    }
    stapConnNetSoc, err := utils.GetKeyValueString("execute", "stapConnNetSoc")
    if err != nil {
        logs.Error("ping/PingPluginsNode Error getting data from main.conf")
        return nil, err
    }
    stapConn, err := utils.GetKeyValueString("execute", "stapConn")
    if err != nil {
        logs.Error("ping/PingPluginsNode Error getting data from main.conf")
        return nil, err
    }
    param, err := utils.GetKeyValueString("execute", "param")
    if err != nil {
        logs.Error("ping/PingPluginsNode Error getting data from main.conf")
        return nil, err
    }
    suriPID, err := utils.GetKeyValueString("execute", "suriPID")
    if err != nil {
        logs.Error("ping/PingPluginsNode Error getting data from main.conf")
        return nil, err
    }
    socatPID, err := utils.GetKeyValueString("execute", "socatPID")
    if err != nil {
        logs.Error("ping/PingPluginsNode Error getting data from main.conf")
        return nil, err
    }
    openSSL, err := utils.GetKeyValueString("execute", "openSSL")
    if err != nil {
        logs.Error("ping/PingPluginsNode Error getting data from main.conf")
        return nil, err
    }
    pidID, err := utils.GetKeyValueString("execute", "pidID")
    if err != nil {
        logs.Error("ping/PingPluginsNode Error getting data from main.conf")
        return nil, err
    }
    command, err := utils.GetKeyValueString("execute", "command")
    if err != nil {
        logs.Error("ping/PingPluginsNode Error getting data from main.conf")
        return nil, err
    }
    check, err := utils.GetKeyValueString("execute", "check")
    if err != nil {
        logs.Error("ping/PingPluginsNode Error getting data from main.conf")
        return nil, err
    }
    checkTCPDUMP, err := utils.GetKeyValueString("stap", "checkTCPDUMP")
    if err != nil {
        logs.Error("ping/PingPluginsNode Error getting data from main.conf")
        return nil, err
    }
    checkTCPREPLAY, err := utils.GetKeyValueString("stap", "checkTCPREPLAY")
    if err != nil {
        logs.Error("ping/PingPluginsNode Error getting data from main.conf")
        return nil, err
    }
    checkSOCAT, err := utils.GetKeyValueString("stap", "checkSOCAT")
    if err != nil {
        logs.Error("ping/PingPluginsNode Error getting data from main.conf")
        return nil, err
    }
    greenMax, err := utils.GetKeyValueInt("stap", "greenMax")
    if err != nil {
        logs.Error("ping/PingPluginsNode Error getting data from main.conf")
        return nil, err
    }
    greenMin, err := utils.GetKeyValueInt("stap", "greenMin")
    if err != nil {
        logs.Error("ping/PingPluginsNode Error getting data from main.conf")
        return nil, err
    }
    yellowMax, err := utils.GetKeyValueInt("stap", "yellowMax")
    if err != nil {
        logs.Error("ping/PingPluginsNode Error getting data from main.conf")
        return nil, err
    }
    yellowMin, err := utils.GetKeyValueInt("stap", "yellowMin")
    if err != nil {
        logs.Error("ping/PingPluginsNode Error getting data from main.conf")
        return nil, err
    }

    //get plugins
    allPlugins, err := ndb.GetPlugins()
    if err != nil {
        logs.Error("ping/GetMainconfData error getting GetPlugins values: " + err.Error())
        return nil, err
    }

    //check if tcpdump, pcapreplay and socat are installed
    allPlugins["installed"] = map[string]string{}
    checkTcpdump, err := exec.Command(check, checkTCPDUMP).Output()
    if len(checkTcpdump) > 0 {
        allPlugins["installed"]["checkTcpdump"] = "true"
    } else {
        allPlugins["installed"]["checkTcpdump"] = "false"
    }

    checkTcpreplay, err := exec.Command(check, checkTCPREPLAY).Output()
    if len(checkTcpreplay) > 0 {
        allPlugins["installed"]["checkTcpreplay"] = "true"
    } else {
        allPlugins["installed"]["checkTcpreplay"] = "false"
    }

    checkSocat, err := exec.Command(check, checkSOCAT).Output()
    if len(checkSocat) > 0 {
        allPlugins["installed"]["checkSocat"] = "true"
    } else {
        allPlugins["installed"]["checkSocat"] = "false"
    }

    for x := range allPlugins {
        if allPlugins[x]["type"] == "suricata" {
            path, err := utils.GetKeyValueString("suricataRuleset", "path")
            if err != nil {
                logs.Error("ping/PingPluginsNode Error getting data from main.conf: " + err.Error())
                return nil, err
            }
            fileToEdit, err := utils.GetKeyValueString("suricataRuleset", "file")
            if err != nil {
                logs.Error("ping/PingPluginsNode Error getting data from main.conf: " + err.Error())
                return nil, err
            }

            //check if ruleset exists locally            
            fileName := strings.Replace(fileToEdit, "<NAME>", allPlugins[x]["localRulesetName"], -1)            
            if _, err := os.Stat(path+fileName); os.IsNotExist(err) {
                allPlugins[x]["rulesetSync"] = "false"
            }else{
                allPlugins[x]["rulesetSync"] = "true"
            }

            //check if process is running even though database status is disabled
            pid, err := exec.Command(command, param, strings.Replace(suriPID, "<ID>", "grep "+x+" |", -1)).Output()
            if err != nil {
                logs.Error("ping/PingPluginsNode Checking suricata PID: " + err.Error())
            }
            if strings.Split(string(pid), "\n")[0] == "" {
                allPlugins[x]["running"] = "false"
            } else {
                allPlugins[x]["running"] = "true"
            }
        }
        if allPlugins[x]["status"] == "enabled" && allPlugins[x]["type"] == "suricata" { 
            // change pid file name
            if _, err := os.Stat(bck + x + "-pidfile.pid"); os.IsNotExist(err) {
                err = suricata.StopSuricataService(x, allPlugins[x]["status"])
                if err != nil {
                    logs.Error("ping/PingPluginsNode pidfile doesn't exist. Error stopping suricata for launch again: " + err.Error())
                    return nil, err
                }
                err = suricata.LaunchSuricataService(x, allPlugins[x]["interface"])
                if err != nil {
                    logs.Error("ping/PingPluginsNode pidfile doesn't exist. Error launching suricata again: " + err.Error())
                    return nil, err
                }
            }        
        }

        //check if process is running even though database status is enabled
        if (allPlugins[x]["type"] == "socket-pcap" || allPlugins[x]["type"] == "socket-network") && allPlugins[x]["pid"] != "none" {
            pid, err := exec.Command(command, param, strings.Replace(socatPID, "<PORT>", allPlugins[x]["port"], -1)).Output()
            if err != nil {
                logs.Error("ping/PingPluginsNode Checking STAP PID: " + err.Error())
            }
            if strings.Split(string(pid), "\n")[0] == "" {
                allPlugins[x]["running"] = "false"
            } else {
                allPlugins[x]["running"] = "true"
            }
        }
        //check if process is running even though database status is enabled
        if allPlugins[x]["type"] == "network-socket" && allPlugins[x]["pid"] != "none" {
            val := strings.Replace(openSSL, "<COLLECTOR>", allPlugins[x]["collector"], -1)
            allValues := strings.Replace(val, "<PORT>", allPlugins[x]["port"], -1)

            pid, err := exec.Command(command, param, allValues).Output()
            if err != nil {
                logs.Error("ping/PingPluginsNode Checking STAP network-socket PID: " + err.Error())
            }
            pids := strings.Split(string(pid), "\n")
            exists := false
            for q := range pids {
                if pids[q] == allPlugins[x]["pid"] {
                    exists = true
                }
            }
            if !exists {
                allPlugins[x]["running"] = "false"
            } else {
                allPlugins[x]["running"] = "true"
            }
        }

        //get all stap connections 
        if allPlugins[x]["type"] == "network-socket" {                    
            ipReplace:=strings.Replace(stapConnNetSoc, "<IP>", allPlugins[x]["ip"], -1)
            portReplace:=strings.Replace(ipReplace, "<PORT>", allPlugins[x]["port"], -1)
            data, err := exec.Command(command, param, portReplace).Output()
            if err != nil {logs.Error("ping/PingPluginsNode getting STAP connections: " + err.Error())}    
            allPlugins[x]["connections"] = string(data)            
        }else if allPlugins[x]["type"] == "socket-network" || allPlugins[x]["type"] == "socket-pcap"{
            data, err := exec.Command(command, param, strings.Replace(stapConn, "<PORT>", allPlugins[x]["port"], -1)).Output()
            if err != nil {logs.Error("ping/PingPluginsNode getting STAP connections: " + err.Error())}
            allPlugins[x]["connections"] = string(data)                    
        } 

        if allPlugins[x]["type"] == "network-socket" || allPlugins[x]["type"] == "socket-network" || allPlugins[x]["type"] == "socket-pcap"{            
            //split connections
            splitted := strings.Split(allPlugins[x]["connections"], "\n")
            var dataConn []string
            for _,val := range splitted {
                if val != "" {
                    dataConn = append(dataConn,  val)
                }
            }
    
            //get number of connections
            allPlugins[x]["connectionsCount"] = strconv.Itoa(len(dataConn))
    
            //check clients umbral
            if len(dataConn) <= greenMax && len(dataConn) >= greenMin {
                allPlugins[x]["connectionsColor"] = "success"
            }else if (len(dataConn) <= yellowMax) && (len(dataConn) >= yellowMin) {
                allPlugins[x]["connectionsColor"] = "warning"            
            }else{
                allPlugins[x]["connectionsColor"] = "danger"
            }
        }          
    }

    //get suricata values that are not in the database
    var avoidUUIDS string
    for f := range allPlugins {
        if allPlugins[f]["type"] == "suricata" {
            avoidUUIDS = avoidUUIDS + "grep -v " + f + " |"
        }
    }

    com, err := exec.Command(command, param, strings.Replace(suriPID, "<ID>", avoidUUIDS, -1)).Output()
    if err != nil {
        logs.Error("PingPluginsNode error getting suricata PIDs: " + err.Error())
    }
    pidValue := strings.Split(string(com), "\n")
    for pid := range pidValue {
        if pidValue[pid] != "" {
            fullCommand, err := exec.Command(command, param, strings.Replace(pidID, "<PID>", pidValue[pid], -1)).Output()
            if err != nil {
                logs.Error("PingPluginsNode error getting suricata shell full command: " + err.Error())
            }

            existsPid := false
            for f := range allPlugins {
                if allPlugins[f]["type"] == "suricata" && allPlugins[f]["pid"] == pidValue[pid] {
                    existsPid = true
                }
            }
            if !existsPid {
                uuid := utils.Generate()
                pluginNotControlled := make(map[string]string)
                pluginNotControlled["type"] = "suricata"
                pluginNotControlled["pid"] = pidValue[pid]
                pluginNotControlled["command"] = string(fullCommand)
                allPlugins[uuid] = pluginNotControlled
            }
        }
    }

    return allPlugins, err
}

func UpdateNodeData(data map[string]map[string]string) (err error) {
    for x, y := range data {
        for y := range y {
            err = ndb.UpdateNodeData(x, y, data[x][y])
            if err != nil {
                logs.Error("UpdateNodeData Error updating node data: " + err.Error())
                return err
            }
        }
    }

    return nil
}

func SaveNodeInformation(anode map[string]map[string]string) (err error) {
    // nodeData, err := ndb.GetNodeData()
    // if err != nil {
    //     logs.Error("SaveNodeInformation Error getting node data: " + err.Error())
    //     return err
    // }
    for x := range anode {
        err = ndb.DeleteNodeInformation()
        if err != nil {
            logs.Error("SaveNodeInformation Error updating node values: " + err.Error())
            return err
        }
        err = ndb.InsertNodeData(x, "ip", anode[x]["ip"])
        if err != nil {
            logs.Error("SaveNodeInformation Error inserting node ip: " + err.Error())
            return err
        }
        err = ndb.InsertNodeData(x, "name", anode[x]["name"])
        if err != nil {
            logs.Error("SaveNodeInformation Error inserting node name: " + err.Error())
            return err
        }
        err = ndb.InsertNodeData(x, "port", anode[x]["port"])
        if err != nil {
            logs.Error("SaveNodeInformation Error inserting node port: " + err.Error())
            return err
        }
    }

    return nil
}

func DeleteNode(masterID string) (err error) {
    //delete node information
    err = ndb.DeleteNodeInformation()
    if err != nil {
        logs.Error("DeleteNode Error deleting node data: " + err.Error())
        return err
    }

    masters, err := ndb.GetMasters()
    for x:= range masters{
        if masters[x]["master"] == masterID{
            err = ndb.DeleteMastersInformation(x)
            if err != nil {
                logs.Error("DeleteNode Error deleting master data: " + err.Error())
                return err
            }
        }
    }

    return err
}
