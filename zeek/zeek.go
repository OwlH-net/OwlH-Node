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
)

func ZeekPath() (exists bool) {
	var err error
    //Retrieve path for wazuh.
	loadDataZeekPath := map[string]map[string]string{}
	loadDataZeekPath["loadDataZeekPath"] = map[string]string{}
	loadDataZeekPath["loadDataZeekPath"]["path"] = ""
    loadDataZeekPath,err = utils.GetConf(loadDataZeekPath)    
    path := loadDataZeekPath["loadDataZeekPath"]["path"]
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
    //Retrieve bin for wazuh.
	loadDataZeekBin := map[string]map[string]string{}
	loadDataZeekBin["loadDataZeekBin"] = map[string]string{}
    loadDataZeekBin["loadDataZeekBin"]["bin"] = ""
    loadDataZeekBin,err = utils.GetConf(loadDataZeekBin)    
    bin := loadDataZeekBin["loadDataZeekBin"]["bin"]
	if err != nil {
		logs.Error("ZeekBin Error getting data from main.conf: "+err.Error())
		return false
	}
    _,err = os.Stat(bin)
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
	var err error
    //Retrieve running for zeek.
	loadDataZeekRunning := map[string]map[string]string{}
	loadDataZeekRunning["loadDataZeekRunning"] = map[string]string{}
    loadDataZeekRunning["loadDataZeekRunning"]["cmd"] = ""
    loadDataZeekRunning["loadDataZeekRunning"]["param"] = ""
    loadDataZeekRunning["loadDataZeekRunning"]["command"] = ""
    loadDataZeekRunning,err = utils.GetConf(loadDataZeekRunning)    
    cmd := loadDataZeekRunning["loadDataZeekRunning"]["cmd"]
    param := loadDataZeekRunning["loadDataZeekRunning"]["param"]
    command := loadDataZeekRunning["loadDataZeekRunning"]["command"]
	if err != nil {
		logs.Error("ZeekRunning Error getting data from main.conf")
		return false
	}
	out, err := exec.Command(command, param, cmd).Output()
	if err != nil {
		logs.Error("Zeek is NOT running: "+err.Error())
		return false
	}
	logs.Error("String out zeek Running: "+string(out))
	if strings.Contains(string(out), "running") {
		logs.Info("Zeek is now running: "+string(out))
		return true
	}
	return false    
}

func Installed() (isIt map[string]bool, err error){
    zeek := make(map[string]bool)
    zeek["path"] = ZeekPath()
    zeek["bin"] = ZeekBin()
	zeek["running"] = ZeekRunning()

    if zeek["path"] || zeek["bin"] || zeek["running"]  {
        logs.Info("Zeek installed and running")
        return zeek, nil
    } else {
        logs.Error("Zeek isn't present or not running")
        return zeek, errors.New("Zeek isn't present or not running")
    }
}

//Run zeek
func RunZeek()(data string, err error){
	//Retrieve path for RunZeek.
	logs.Warn("RunZeek")
    StartZeek := map[string]map[string]string{}
    StartZeek["zeekStart"] = map[string]string{}
    StartZeek["zeekStart"]["start"] = ""
    StartZeek["zeekStart"]["param"] = ""
    StartZeek["zeekStart"]["command"] = ""
    StartZeek,err = utils.GetConf(StartZeek)    
    cmd := StartZeek["zeekStart"]["start"]
    param := StartZeek["zeekStart"]["param"]
    command := StartZeek["zeekStart"]["command"]
	if err != nil {
		logs.Error("RunZeek Error getting data from main.conf: "+err.Error())
		return "", err
	}
    _,err = exec.Command(command, param, cmd).Output()
    if err != nil {
        logs.Error("Error launching zeekStart: "+err.Error())
        return "",err
    }
    return "zeekStart system is on",nil
}

//Stop zeek
func StopZeek()(data string, err error){
    //Retrieve path for zeek.
    StopZeek := map[string]map[string]string{}
	StopZeek["zeekStop"] = map[string]string{}
    StopZeek["zeekStop"]["stop"] = ""
    StopZeek["zeekStop"]["param"] = ""
    StopZeek["zeekStop"]["command"] = ""
    StopZeek,err = utils.GetConf(StopZeek)    
    cmd := StopZeek["zeekStop"]["stop"]
    param := StopZeek["zeekStop"]["param"]
    command := StopZeek["zeekStop"]["command"]
	if err != nil {
		logs.Error("StopZeek Error getting data from main.conf: "+err.Error())
	}
	_,err = exec.Command(command, param, cmd).Output()
    if err != nil {
        logs.Error("Error stopping zeek: "+err.Error())
        return "",err
    }
    return "Zeek stopped ",nil
}

//Stop zeek
func DeployZeek()(err error){
    //Retrieve path for zeek.
    DeployZeek := map[string]map[string]string{}
    DeployZeek["zeek"] = map[string]string{}
    DeployZeek["zeek"]["deploy"] = ""
    DeployZeek["zeek"]["zeekctl"] = ""
    //DeployZeek["zeekDeploy"]["command"] = ""
    DeployZeek,err = utils.GetConf(DeployZeek)    
    //cmd := DeployZeek["zeekDeploy"]["cmd"]
    //param := DeployZeek["zeekDeploy"]["param"],
    //command := DeployZeek["zeekDeploy"]["command"]
    if err != nil {
        logs.Error("DeployZeek Error getting data from main.conf: "+err.Error())
    }
    err = utils.RunCommand(DeployZeek["zeek"]["zeekctl"],DeployZeek["zeek"]["deploy"])
    if err != nil {logs.Error("Error deploying zeek: "+err.Error()); return err}

    return nil
}

func ChangeZeekMode(anode map[string]string) (err error) {
    err = ndb.UpdateMainconfValue("zeek", "mode", anode["mode"])
    if err != nil {logs.Error("Error ChangeZeekMode: "+err.Error()); return err}
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
    zeekPath := map[string]map[string]string{}
    zeekPath["loadDataZeekPath"] = map[string]string{}
    zeekPath["loadDataZeekPath"]["nodeConfig"] = ""
	zeekPath,err = utils.GetConf(zeekPath)
	if err != nil {logs.Error("SyncCluster Error readding GetConf: "+err.Error())}
    path := zeekPath["loadDataZeekPath"]["nodeConfig"]
    
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