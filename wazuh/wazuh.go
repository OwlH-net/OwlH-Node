package wazuh

import (
    "github.com/astaxie/beego/logs"
    "os"
    "os/exec"
    "strings"
    "regexp"
    "owlhnode/utils"
)

func WazuhPath() (exists bool) {
    //Retrieve path for wazuh.
	loadDataWazuhPath := map[string]map[string]string{}
	loadDataWazuhPath["loadDataWazuhPath"] = map[string]string{}
	loadDataWazuhPath["loadDataWazuhPath"]["path"] = ""
    loadDataWazuhPath = utils.GetConf(loadDataWazuhPath)    
    path := loadDataWazuhPath["loadDataWazuhPath"]["path"]
    
    if _, err := os.Stat(path); os.IsNotExist(err) {
        logs.Error("Wazuh no esta instalado, al menos la carpeta /var/ossec no existe")
        return false
    }
    return true
}

func WazuhBin() (exists bool) {
    //Retrieve bin for wazuh.
	loadDataWazuhBin := map[string]map[string]string{}
	loadDataWazuhBin["loadDataWazuhBin"] = map[string]string{}
    loadDataWazuhBin["loadDataWazuhBin"]["path"] = ""
    loadDataWazuhBin["loadDataWazuhBin"]["param"] = ""
    loadDataWazuhBin = utils.GetConf(loadDataWazuhBin)    
    path := loadDataWazuhBin["loadDataWazuhBin"]["path"]
    param := loadDataWazuhBin["loadDataWazuhBin"]["param"]

    out, err := exec.Command(path,param).Output()
    if err == nil {
        if strings.Contains(string(out), "Wazuh") {
            logs.Info("Wazuh binario existe -> " + string(out))
            return true
        }
    }
    logs.Error("Wazuh binario NO existe")
    return false
}

func WazuhRunning() (running bool) {
    //Retrieve running for suricata.
	loadDatasuriWazuhRunning := map[string]map[string]string{}
	loadDatasuriWazuhRunning["loadDatasuriWazuhRunning"] = map[string]string{}
    loadDatasuriWazuhRunning["loadDatasuriWazuhRunning"]["cmd"] = ""
    loadDatasuriWazuhRunning["loadDatasuriWazuhRunning"]["param"] = ""
    loadDatasuriWazuhRunning["loadDatasuriWazuhRunning"]["command"] = ""
    loadDatasuriWazuhRunning = utils.GetConf(loadDatasuriWazuhRunning)    
    cmd := loadDatasuriWazuhRunning["loadDatasuriWazuhRunning"]["cmd"]
    param := loadDatasuriWazuhRunning["loadDatasuriWazuhRunning"]["param"]
    command := loadDatasuriWazuhRunning["loadDatasuriWazuhRunning"]["command"]

    //cmd := "ps -ef | grep ossec | grep -v grep | grep -v sudo | awk '{print $8 \" \" $2}' "
    //out, err := exec.Command("bash", "-c", cmd).Output()
    out, err := exec.Command(command, param, cmd).Output()
    if err == nil {
        if strings.Contains(string(out), "ossec") {
            spid := regexp.MustCompile("[0-9]+")
            pid := spid.FindAllString(string(out),1)
            logs.Info("Wazuh esta en ejecucion -> " + string(out))
            logs.Info("Wazuh esta en ejecucion PID -> %s", pid[0])
            return true
        }
    }
    logs.Error("Wazuh NO esta en ejecucion -> " + string(out))
    return false
}

func Installed() (isIt bool){
    var Wazuh bool
    Wazuh = false
    Wazuh = WazuhPath()
    Wazuh = WazuhBin()
    Wazuh = WazuhRunning()
    if Wazuh {
        logs.Info("Wazuh Existe")
        return true
    } else {
        logs.Error("Wazuh no existe")
    }
    return false
}