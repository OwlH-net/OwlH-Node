package zeek

import (
    "github.com/astaxie/beego/logs"
    "os"
    "os/exec"
    "strings"
    "regexp"
    "owlhnode/utils"
)

func ZeekPath() (exists bool) {
    //Retrieve path for wazuh.
	loadDataZeekPath := map[string]map[string]string{}
	loadDataZeekPath["loadDataZeekPath"] = map[string]string{}
	loadDataZeekPath["loadDataZeekPath"]["path"] = ""
    loadDataZeekPath = utils.GetConf(loadDataZeekPath)    
    path := loadDataZeekPath["loadDataZeekPath"]["path"]

    //if _, err := os.Stat("/etc/zeek"); os.IsNotExist(err) {
    if _, err := os.Stat(path); os.IsNotExist(err) {
        logs.Error("Zeek no esta instalado, al menos la carpeta /etc/zeek no existe")
        return false
    }
    return true
}

func ZeekBin() (exists bool) {
    //Retrieve bin for wazuh.
	loadDataZeekBin := map[string]map[string]string{}
	loadDataZeekBin["loadDataZeekBin"] = map[string]string{}
    loadDataZeekBin["loadDataZeekBin"]["cmd"] = ""
    loadDataZeekBin["loadDataZeekBin"]["param"] = ""
    loadDataZeekBin = utils.GetConf(loadDataZeekBin)    
    cmd := loadDataZeekBin["loadDataZeekBin"]["cmd"]
    param := loadDataZeekBin["loadDataZeekBin"]["param"]

    //out, err := exec.Command("broctl","-V").Output()
    out, err := exec.Command(cmd,param).Output()
    if err == nil {
        if strings.Contains(string(out), "Zeek version") {
            logs.Info("Zeek binario existe -> " + string(out))
            return true
        }
    }
    logs.Error("Zeek binario NO existe")
    return false
}

func ZeekRunning() (running bool) {
    //Retrieve running for suricata.
	loadDatasuriZeekRunning := map[string]map[string]string{}
	loadDatasuriZeekRunning["loadDatasuriZeekRunning"] = map[string]string{}
    loadDatasuriZeekRunning["loadDatasuriZeekRunning"]["cmd"] = ""
    loadDatasuriZeekRunning["loadDatasuriZeekRunning"]["param"] = ""
    loadDatasuriZeekRunning["loadDatasuriZeekRunning"]["command"] = ""
    loadDatasuriZeekRunning = utils.GetConf(loadDatasuriZeekRunning)    
    cmd := loadDatasuriZeekRunning["loadDatasuriZeekRunning"]["cmd"]
    param := loadDatasuriZeekRunning["loadDatasuriZeekRunning"]["param"]
    command := loadDatasuriZeekRunning["loadDatasuriZeekRunning"]["command"]

    //cmd := "ps -ef | grep zeek | grep -v grep | grep -v sudo | awk '{print $8 \" \" $2}' "
    //out, err := exec.Command("bash", "-c", cmd).Output()
    out, err := exec.Command(command, param, cmd).Output()
    if err == nil {
        if strings.Contains(string(out), "zeek") {
            spid := regexp.MustCompile("[0-9]+")
            pid := spid.FindAllString(string(out),1)
            logs.Info("Zeek esta en ejecucion -> " + string(out))
            logs.Info("Zeek esta en ejecucion PID -> %s", pid[0])
            return true
        }
    }
    logs.Error("Zeek NO esta en ejecucion -> " + string(out))
    return false
}

func Installed() (isIt bool){
    var zeek bool
    zeek = false
    zeek = ZeekPath()
    zeek = ZeekBin()
    zeek = ZeekRunning()
    if zeek {
        logs.Info("Zeek Existe")
        return true
    } else {
        logs.Error("Zeek no existe")
    }
    return false
}