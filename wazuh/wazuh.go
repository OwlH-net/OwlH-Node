package wazuh

import (
    "github.com/astaxie/beego/logs"
    "os"
    "os/exec"
    "strings"
    "regexp"
)

func WazuhPath() (exists bool) {
    if _, err := os.Stat("/var/ossec"); os.IsNotExist(err) {
        logs.Error("Wazuh no esta instalado, al menos la carpeta /var/ossec no existe")
        return false
    }
    return true
}

func WazuhBin() (exists bool) {
    out, err := exec.Command("/var/ossec/bin/ossec-control","-V").Output()
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
    cmd := "ps -ef | grep ossec | grep -v grep | grep -v sudo | awk '{print $8 \" \" $2}' "
    out, err := exec.Command("bash", "-c", cmd).Output()
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