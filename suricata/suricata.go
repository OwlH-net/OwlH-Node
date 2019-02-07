package suricata

import (
    "github.com/astaxie/beego/logs"
    "os"
    "os/exec"
    "strings"
    "regexp"
)

func suriPath() (exists bool) {
    if _, err := os.Stat("/etc/suricata"); os.IsNotExist(err) {
        logs.Error("Suricata no esta instalado, al menos la carpeta /etc/suricata no existe")
        return false
    }
    return true
}

func suriBin() (exists bool) {
    out, err := exec.Command("suricata","-V").Output()
    if err == nil {
        if strings.Contains(string(out), "Suricata version") {
            logs.Info("Suricata binario existe -> " + string(out))
            return true
        }
    }
    logs.Error("Suricata binario NO existe")
    return false
}

func suriRunning() (running bool) {
    cmd := "ps -ef | grep suricata | grep -v grep | grep -v sudo | awk '{print $8 \" \" $2}' "
    out, err := exec.Command("bash", "-c", cmd).Output()
    if err == nil {
        if strings.Contains(string(out), "suricata") {
            spid := regexp.MustCompile("[0-9]+")
            pid := spid.FindAllString(string(out),1)
            logs.Info("Suricata esta en ejecucion -> " + string(out))
            logs.Info("Suricata esta en ejecucion PID -> %s", pid[0])
            return true
        }
    }
    logs.Error("Suricata NO esta en ejecucion -> " + string(out))
    return false
}

func Installed() (isIt bool){
    var suricata bool
    suricata = false
    suricata = suriPath()
    suricata = suriBin()
    suricata = suriRunning()
    if suricata {
        logs.Info("Suricata Existe")
        return true
    } else {
        logs.Error("Suricata no existe")
    }
    return false
}