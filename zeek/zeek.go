package zeek

import (
    "github.com/astaxie/beego/logs"
    "os"
    "os/exec"
    "strings"
    "strconv"
    "owlhnode/utils"
    "owlhnode/database"
    "errors"
    "regexp"
    "bufio"
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


func ZeekPath() (exists bool) {
    var err error
    path, err := utils.GetKeyValueString("zeek", "zeekpath")
    if err != nil {
        logs.Error("ZeekPath Error getting data from main.conf: "+err.Error())
        return false
    }

    if _, err := os.Stat(path); os.IsNotExist(err) {
        logs.Error("Zeek is not installed on "+path+"."+err.Error())
        return false
    }
    return true
}

func ZeekBin() (exists bool) {
    var err error
    zeekctl, err := utils.GetKeyValueString("zeek", "zeekctl")
    if err != nil {
        logs.Error("ZeekBin Error getting data from main.conf: "+err.Error())
        return false
    }

    _,err = os.Stat(zeekctl)
    if err != nil {
        logs.Error("Zeek OS path err: "+err.Error())
        return false
    }

    if os.IsNotExist(err){
        logs.Error("Zeek path not exist: "+err.Error())
        return false
    }
    return true
}

func ZeekRunning() (running bool) {
    return false
    // DEPRECATED
    // var err error
    // //_, err = ZeekStatus()
    // //Retrieve running for zeek.
    // loadDataZeekRunning := map[string]map[string]string{}
    // loadDataZeekRunning["loadDataZeekRunning"] = map[string]string{}
    // loadDataZeekRunning["loadDataZeekRunning"]["cmd"] = ""
    // loadDataZeekRunning["loadDataZeekRunning"]["param"] = ""
    // loadDataZeekRunning["loadDataZeekRunning"]["command"] = ""
    // loadDataZeekRunning,err = utils.GetConf(loadDataZeekRunning)    
    // cmd := loadDataZeekRunning["loadDataZeekRunning"]["cmd"]
    // param := loadDataZeekRunning["loadDataZeekRunning"]["param"]
    // command := loadDataZeekRunning["loadDataZeekRunning"]["command"]
    // if err != nil {
    //     logs.Error("ZeekRunning Error getting data from main.conf")
    //     return false
    // }
    // out, err := exec.Command(command, param, cmd).Output()
    // if err != nil {
    //     logs.Error("Zeek is NOT running: "+err.Error())
    //     return false
    // }
    // logs.Error("String out zeek Running: "+string(out))
    // if strings.Contains(string(out), "running") {
    //     logs.Info("Zeek is now running: "+string(out))
    //     return true
    // }
    // return false    
}

func ZeekStatus() (zeekstatus []ZeekNode, err error) {
    zeekctl, err := utils.GetKeyValueString("zeek", "zeekctl")
    if err != nil {logs.Error(errors.New("ZEEK STATUS -> Error getting Status command from main.conf"))}
    currentStatus, err := utils.GetKeyValueString("zeek", "currentstatus")
    if err != nil {logs.Error(errors.New("ZEEK STATUS -> Error getting Status command from main.conf"))}

    logs.Info("Zeek CTL -> %s", zeekctl)
    logs.Info("Zeek currentstatus -> %s", currentStatus)
    output, err:= exec.Command(zeekctl, currentStatus).Output()
    if err != nil {logs.Error(errors.New("ZEEK STATUS -> Error running status command -> " + err.Error()))}

    nodes := []ZeekNode{}
    outputlines := strings.Split(string(output),"\n")
    for outputline := range outputlines {
        line := strings.Fields(outputlines[outputline])
        if len(line) > 2 {
            if strings.Contains(line[1], "manager") || strings.Contains(line[1], "logger") || strings.Contains(line[1], "proxy") || strings.Contains(line[1], "worker") || strings.Contains(line[1], "standalone") {
                node := ZeekNode{}
                node.Name = line[0]
                node.Type = line[1]
                node.Host = line[2]
                node.Status = line[3]
                if len(line) > 4 {
                    node.Pid = line[4]
                }
                nodes = append(nodes, node)
            } else {
                logs.Error ("Zeek -> status output: "+ outputlines[outputline])
            }
        }
    }
    logs.Info(nodes)
    return nodes, nil 
}


func GetZeek()(zeek Zeek) {
    zeek.Path = ZeekPath()
    zeek.Bin = ZeekBin()
    zeek.Mode = ZeekMode()
    zeek.Managed = ZeekManaged()
    nodes, err := ZeekStatus()
    zeek.Nodes = nodes
    if err != nil {
        logs.Info(err.Error())
    }
    for node := range nodes {
        updated := false
        for nstatus := range zeek.Running {
            if zeek.Running[nstatus].Status == nodes[node].Status {
                updated = true
                zeek.Running[nstatus].Nodes++
                break
            }
        }
        if !updated {
            newStatus := ZeekNodeStatus{}
            newStatus.Status = nodes[node].Status
            newStatus.Nodes = 1
            zeek.Running = append(zeek.Running, newStatus)
        }
    }
    return zeek
}

func SetZeek(zeekdata Zeek)(newzeekdata Zeek, err error) {
    for node := range zeekdata.Nodes {
        logs.Warn("=============")
        logs.Warn("name - "+zeekdata.Nodes[node].Name)
        logs.Warn("interface - "+zeekdata.Nodes[node].NInterface)
        logs.Warn("host - "+zeekdata.Nodes[node].Host)
        logs.Warn("type - "+zeekdata.Nodes[node].Type)
        logs.Warn("=============")
        logs.Warn("======= EXTRA ========")
        for extra := range zeekdata.Nodes[node].Extra {
            logs.Warn(zeekdata.Nodes[node].Extra[extra])
            logs.Warn("key - "+zeekdata.Nodes[node].Extra[extra].Key + " -- " + zeekdata.Nodes[node].Extra[extra].Value )
        }
        logs.Warn("======= EXTRA ========")
    }

    newzeekdata = GetZeek()
    return newzeekdata, nil
}

func ZeekMode()(mode string) {
    nodeConfig, err := utils.GetKeyValueString("zeek", "nodeconfig")
    if err != nil {logs.Error(err); return}

    file, err := os.Open(nodeConfig)
    if err != nil {logs.Error(err); return}
    defer file.Close()

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        var init = regexp.MustCompile(`type=`)
        owlhInit := init.FindStringSubmatch(scanner.Text())
        if owlhInit != nil {
            val := strings.Trim(scanner.Text(), "type=")
            //save current mode into DB
            err = ndb.UpdateMainconfValue("zeek","mode",val) //uuid, param, value
            if err != nil {logs.Error("Error Zeek Mode Get current Mode: "+err.Error()); return "Error: Zeek Mode PUT current Mode at DB: "+err.Error()}
        }
    }

    //get current mode into DB
    currentmode, err := ndb.GetMainconfParam("zeek","mode")
    if err != nil {logs.Error("Error Zeek Mode Get current Mode: "+err.Error()); return "Error: Zeek Mode Get current Mode at DB: "+err.Error()}

    logs.Info("ZEEK -> current Mode: " + currentmode)
    return currentmode
}

func ZeekManaged()(ismanaged bool) {
    currentstatus, err := ndb.GetMainconfParam("zeek","status")
    if err != nil {
        logs.Error("Error Zeek Mode Get current Status: "+err.Error())
        return false
    }
    logs.Info("ZEEK -> current Managed Mode: " + currentstatus)
    if currentstatus == "enabled" {return true}
    return false 
}

//Run zeek
func RunZeek()(data string, err error){   
    zeekctl, err := utils.GetKeyValueString("zeek", "zeekctl")
    if err != nil {logs.Error("RunZeek Error getting data from main.conf: "+err.Error()); return "", err}
    zeekStart, err := utils.GetKeyValueString("zeek", "start")
    if err != nil {logs.Error("RunZeek Error getting data from main.conf: "+err.Error()); return "", err}

    err = utils.RunCommand(zeekctl,zeekStart)
    if err != nil {logs.Error("Error starting zeek: "+err.Error()); return "",err} 

    return "Zeek system is on",nil
}

//Start Zeek
func StartZeek(action string)(data string, err error){
    // ACTION - Start or Deploy
    getaction := "deploy"
    if action != "" { getaction = action } 
    logs.Warn("Starting Zeek by deploy")
    cmd, err := utils.GetKeyValueString("zeek", "zeekctl")
    if err != nil {logs.Error("StartZeek Error getting data from main.conf: "+err.Error()); return "", err}
    realaction, err := utils.GetKeyValueString("zeek", getaction)
    if err != nil {logs.Error("StartZeek Error getting data from main.conf: "+err.Error()); return "", err}

    output,err := exec.Command(cmd, realaction).Output()
    if err != nil {
        logs.Error("Error launching StartZeek: "+err.Error())
        return "",err
    }
    return string(output), nil
}

func StartingZeek()(err error){   
    cmd, err := utils.GetKeyValueString("zeek", "zeekctl")
    if err != nil {logs.Error("StartingZeek Error getting data from main.conf: "+err.Error())}
    start, err := utils.GetKeyValueString("zeek", "start")
    if err != nil {logs.Error("StartingZeek Error getting data from main.conf: "+err.Error())}

    err = utils.RunCommand(cmd,start)
    if err != nil {logs.Error("Error deploying zeek: "+err.Error()); return err}

    return nil
}

// //Stop zeek
func StopZeek()(data string, err error){    
    cmd, err := utils.GetKeyValueString("zeek", "zeekctl")
    if err != nil {logs.Error("StopZeek Error getting data from main.conf: "+err.Error())}
    stop, err := utils.GetKeyValueString("zeek", "stop")
    if err != nil {logs.Error("StopZeek Error getting data from main.conf: "+err.Error())}

    err = utils.RunCommand(cmd, stop)
    if err != nil {logs.Error("Error deploying zeek: "+err.Error()); return "", err}

    return "Zeek stopped ",nil
}
//Stop zeek
func StopingZeek()(err error){
    cmd, err := utils.GetKeyValueString("zeek", "zeekctl")
    if err != nil {logs.Error("StopingZeek Error getting data from main.conf: "+err.Error())}
    stop, err := utils.GetKeyValueString("zeek", "stop")
    if err != nil {logs.Error("StopingZeek Error getting data from main.conf: "+err.Error())}

    err = utils.RunCommand(cmd, stop)
    if err != nil {logs.Error("Error deploying zeek: "+err.Error()); return err}

    return nil
}

//Deploy zeek
func DeployZeek()(err error){
    zeekctl, err := utils.GetKeyValueString("zeek", "zeekctl")
    if err != nil {logs.Error("DeployZeek Error getting data from main.conf: "+err.Error())}
    deploy, err := utils.GetKeyValueString("zeek", "deploy")
    if err != nil {logs.Error("DeployZeek Error getting data from main.conf: "+err.Error())}

    err = utils.RunCommand(zeekctl, deploy)
    if err != nil {logs.Error("Error deploying zeek: "+err.Error()); return err}

    return nil
}

func ChangeZeekMode(anode map[string]string) (err error) {
    err = ndb.UpdateMainconfValue("zeek", "mode", anode["mode"])
    if err != nil {logs.Error("Error ChangeZeekMode: "+err.Error()); return err}
    SyncCluster(nil,"standalone")
    return err
}

func AddClusterValue(anode map[string]string) (err error) {
    count,err := ndb.CountDBEntries(anode["type"]); if err != nil {logs.Error("Error AddClusterValue type: "+err.Error()); return err}
    count++
    err = ndb.InsertClusterData(anode["type"]+"-"+strconv.Itoa(count), "type", anode["type"]); if err != nil {logs.Error("Error AddClusterValue type: "+err.Error()); return err}
    err = ndb.InsertClusterData(anode["type"]+"-"+strconv.Itoa(count), "host", anode["host"]); if err != nil {logs.Error("Error1 AddClusterValue host: "+err.Error()); return err}
    if anode["type"] == "worker"{
        err = ndb.InsertClusterData(anode["type"]+"-"+strconv.Itoa(count), "interface", anode["interface"]); if err != nil {logs.Error("Error AddClusterValue interface: "+err.Error()); return err}
    }
    return err
}

func PingCluster()(data map[string]map[string]string, err error) {
    data,err = ndb.GetClusterData(); if err != nil {logs.Error("Error Zeek/PingCluster: "+err.Error()); return nil,err}
    return data,err
}

func EditClusterValue(anode map[string]string) (err error) {
    err = ndb.UpdateClusterValue(anode["type"], "host", anode["host"]); if err != nil {logs.Error("Error Zeek/EditClusterValue: "+err.Error()); return err}
    if anode["cluster"] == "worker"{
        err = ndb.UpdateClusterValue(anode["type"], "interface", anode["interface"]); if err != nil {logs.Error("Error Zeek/EditClusterValue: "+err.Error()); return err}
    }
    return err
}

func DeleteClusterValue(anode map[string]string) (err error) {
    err = ndb.DeleteClusterValue(anode["type"]); if err != nil {logs.Error("Error Zeek/DeleteClusterValue: "+err.Error()); return err}
    //change indentifier
    countWorker := 1
    countProxy := 1
    data,err := ndb.GetClusterData(); if err != nil {logs.Error("Error Zeek/DeleteClusterValue: "+err.Error()); return err}
    err = ndb.DeleteAllClusters(); if err != nil {logs.Error("Error Zeek/DeleteClusterValue: "+err.Error()); return err}
    
    for id,_ := range data {
        if id == "manager" || id == "logger"{
            err = ndb.InsertClusterData(id, "host", data[id]["host"]); if err != nil {logs.Error("Error DeleteClusterValue manager: "+err.Error()); return err}
        }else{
            if data[id]["type"] == "worker" {
                err = ndb.InsertClusterData(data[id]["type"]+"-"+strconv.Itoa(countWorker), "type", data[id]["type"]); if err != nil {logs.Error("Error DeleteClusterValue type: "+err.Error()); return err}
                err = ndb.InsertClusterData(data[id]["type"]+"-"+strconv.Itoa(countWorker), "host", data[id]["host"]); if err != nil {logs.Error("Error DeleteClusterValue host: "+err.Error()); return err}
                err = ndb.InsertClusterData(data[id]["type"]+"-"+strconv.Itoa(countWorker), "interface", data[id]["interface"]); if err != nil {logs.Error("Error DeleteClusterValue type: "+err.Error()); return err}
                countWorker++
            }else{
                err = ndb.InsertClusterData(data[id]["type"]+"-"+strconv.Itoa(countProxy), "type", data[id]["type"]); if err != nil {logs.Error("Error DeleteClusterValue type: "+err.Error()); return err}
                err = ndb.InsertClusterData(data[id]["type"]+"-"+strconv.Itoa(countProxy), "host", data[id]["host"]); if err != nil {logs.Error("Error DeleteClusterValue host: "+err.Error()); return err}
                countProxy++
            }
        }
    }

    return err
}

func SyncCluster(anode map[string]string, clusterType string) (err error) {            
    path, err := utils.GetKeyValueString("zeek", "nodeconfig")
    if err != nil {logs.Error("SyncCluster Error readding main.conf: "+err.Error())}
    
    h := 0
    fileContent := make(map[int]string)

    if clusterType == "standalone" {
        fileContent[h] = "[bro]"; h++
        fileContent[h] = "type=standalone"; h++
        fileContent[h] = "host=localhost"; h++
        fileContent[h] = "interface="+anode["value"]; h++
    }else if clusterType == "cluster" {
        data,err := ndb.GetClusterData(); if err != nil {logs.Error("Error Zeek/SyncCluster: "+err.Error()); return err}
        
        for t := range data{
            if t == "logger"{
                fileContent[h] = "[logger]"; h++
                fileContent[h] = "type=logger"; h++
                fileContent[h] = "host="+data[t]["host"]; h++
                fileContent[h] = ""; h++
            }else if t == "manager"{
                fileContent[h] = "[manager]"; h++
                fileContent[h] = "type=manager"; h++
                fileContent[h] = "host="+data[t]["host"]; h++
                fileContent[h] = ""; h++
            }else if data[t]["type"] == "proxy"{
                fileContent[h] = "["+t+"]"; h++
                fileContent[h] = "type="+data[t]["type"]; h++
                fileContent[h] = "host="+data[t]["host"]; h++
                fileContent[h] = ""; h++
            }else if data[t]["type"] == "worker"{
                fileContent[h] = "["+t+"]"; h++
                fileContent[h] = "type="+data[t]["type"]; h++
                fileContent[h] = "host="+data[t]["host"]; h++
                fileContent[h] = "interface="+data[t]["interface"]; h++
                fileContent[h] = ""; h++
            }
        }
    }

    saveIntoFile, err := os.OpenFile(path , os.O_CREATE|os.O_APPEND|os.O_WRONLY, os.ModeAppend)
    if err != nil {logs.Error("Error SyncCluster readding file: "+err.Error()); return err}
    defer saveIntoFile.Close()
    saveIntoFile.Truncate(0)
    saveIntoFile.Seek(0,0)
    for x:=0 ; x < h ; x++{
        _, err = saveIntoFile.WriteAt([]byte(fileContent[x]+"\n"), 0) // Write at 0 beginning
        if err != nil {logs.Error("SyncCluster failed writing to file: %s", err); return err}
    }

    return err
}

func SavePolicyFiles(files map[string]map[string][]byte)(err error){
    for nodePath, file := range files {
        //check path
        if _, err := os.Stat(nodePath); os.IsNotExist(err) {
            os.MkdirAll(nodePath, os.ModePerm)
        }

        for file,_ := range file {            
            err = utils.WriteNewDataOnFile(nodePath+"/"+file, files[nodePath][file])
            if err != nil{logs.Error("SavePolicyFiles Error writting data into "+nodePath+"/"+file+" file: "+err.Error()); return err}
        }
    }
    return nil
}

func SyncClusterFile(anode map[string][]byte) (err error) {
    path, err := utils.GetKeyValueString("zeek", "nodeconfig")
    if err != nil {logs.Error("SyncClusterFile Error readding main.conf: "+err.Error())}

    err = utils.WriteNewDataOnFile(path, anode["data"])
    if err != nil{logs.Error("zeek/SyncClusterFile Error writting cluster file content: "+err.Error()); return err}
    return err
}

// func SaveZeekValues(anode map[string]string) (err error) {
//     plugins,err := ndb.GetPlugins()
//     for x := range plugins{
//         if plugins[x]["type"] == "zeek"{
//             if anode["param"] == "nodeConfig"{err = ndb.UpdatePluginValue(x, anode["param"], anode["nodeConfig"]); if err != nil{logs.Error("zeek/SaveZeekValues Error writting node content: "+err.Error()); return err}}            
//             if anode["param"] == "networksConfig"{err = ndb.UpdatePluginValue(x, anode["param"], anode["networksConfig"]); if err != nil{logs.Error("zeek/SaveZeekValues Error writting networks content: "+err.Error()); return err}}            
//             if anode["param"] == "policies"{
//                 err = ndb.UpdatePluginValue(x, "policiesMaster", anode["policiesMaster"]); if err != nil{logs.Error("zeek/SaveZeekValues Error writting policies content: "+err.Error()); return err}
//                 err = ndb.UpdatePluginValue(x, "policiesNode", anode["policiesNode"]); if err != nil{logs.Error("zeek/SaveZeekValues Error writting policies content: "+err.Error()); return err}
//             }            
//             if anode["param"] == "variables"{
//                 err = ndb.UpdatePluginValue(x, "variables1", anode["variables1"]); if err != nil{logs.Error("zeek/SaveZeekValues Error writting variables content: "+err.Error()); return err}
//                 err = ndb.UpdatePluginValue(x, "variables2", anode["variables2"]); if err != nil{logs.Error("zeek/SaveZeekValues Error writting variables content: "+err.Error()); return err}
//             }            
//         }
//     }
//     return err
// }

// func LaunchZeekMainConf(anode map[string]string) (err error) {
//     zeekPath := map[string]map[string]string{}
//     zeekPath["zeek"] = map[string]string{}
//     zeekPath["zeek"]["command"] = ""
//     zeekPath["zeek"]["param"] = ""
//     zeekPath["zeek"]["zeekctl"] = ""
//     zeekPath["zeek"][anode["param"]] = ""
//     zeekPath,err = utils.GetConf(zeekPath)
//     if err != nil {logs.Error("LaunchZeekMainConf Error readding GetConf: "+err.Error())}
//     command := zeekPath["zeek"]["command"]
//     param := zeekPath["zeek"]["param"]
//     path := zeekPath["zeek"]["zeekctl"]
//     cmd := zeekPath["zeek"][anode["param"]]
//     if err != nil{logs.Error("zeek/LaunchZeekMainConf Error getting main.conf file content: "+err.Error()); return err}

//     _,err = exec.Command(command, param, path, cmd).Output()
//     if err != nil {logs.Error("zeek/LaunchZeekMainCon Error starting Zeek from main conf: "+err.Error());return err}

//     return err
// }

func SyncZeekValues(anode map[string]string) (err error) {

    for x := range anode{
        if x == "nodeConfig"{
            path, err := utils.GetKeyValueString("zeek", "nodeconfig")
            if err != nil {logs.Error("zeek/SyncZeekValues Error readding main.conf: "+err.Error())}
        
            err = utils.BackupFullPath(path)
            if err != nil{logs.Error("zeek/SyncZeekValues Error backing up node.cfg file before overwrite: "+err.Error()); return err}
            err = utils.WriteNewDataOnFile(path, []byte(anode["nodeConfig"]))
            if err != nil{logs.Error("zeek/SyncZeekValues Error writting new file content: "+err.Error()); return err}

        }else if x == "networksConfig"{
            path, err := utils.GetKeyValueString("zeek", "networkconfig")
            if err != nil {logs.Error("zeek/SyncZeekValues Error readding main.conf: "+err.Error())}
        
            err = utils.BackupFullPath(path)
            if err != nil{logs.Error("zeek/SyncZeekValues Error backing up networks.cfg file before overwrite: "+err.Error()); return err}
            err = utils.WriteNewDataOnFile(path, []byte(anode["networksConfig"]))
            if err != nil{logs.Error("zeek/SyncZeekValues Error writting new file content: "+err.Error()); return err}

        }else{
            if _, err := os.Stat(anode["dst"]); os.IsNotExist(err) {
                logs.Error("zeek/SyncZeekValues Destiny file down't exists: "+err.Error()); return err
            }
            err = utils.BackupFullPath(anode["dst"])
            if err != nil{logs.Error("zeek/SyncZeekValues Error backing up file before overwrite: "+err.Error()); return err}
            if x == "policiesMaster"{
                err = utils.WriteNewDataOnFile(anode["dst"], []byte(anode["policiesMaster"]))
                if err != nil{logs.Error("zeek/SyncZeekValues Error writting new file content: "+err.Error()); return err}
            }
            if x == "variables1"{
                err = utils.WriteNewDataOnFile(anode["dst"], []byte(anode["variables1"]))
                if err != nil{logs.Error("zeek/SyncZeekValues Error writting new file content: "+err.Error()); return err}
            }
        }
    }
    
    return err
}