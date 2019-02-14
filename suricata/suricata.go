package suricata

import (
    "github.com/astaxie/beego/logs"
    "os"
    "os/exec"
    "strings"
    "regexp"
    "owlhnode/utils"
)

func suriPath() (exists bool) {
    if _, err := os.Stat("/etc/suricata"); os.IsNotExist(err) {
        logs.Error("Suricata not installed, at least folder /etc/suricata dosn't exist")
        return false
    }
    return true
}

func suriBin() (exists bool) {
    out, err := exec.Command("suricata","-V").Output()
    if err == nil {
        if strings.Contains(string(out), "Suricata version") {
            logs.Info("Suricata installed -> " + string(out))
            return true
        }
    }
    logs.Error("Suricata Suricata not installed")
    return false
}

func suriRunning() (running bool) {
    cmd := "ps -ef | grep suricata | grep -v grep | grep -v sudo | awk '{print $8 \" \" $2}' "
    out, err := exec.Command("bash", "-c", cmd).Output()
    if err == nil {
        if strings.Contains(string(out), "suricata") {
            spid := regexp.MustCompile("[0-9]+")
            pid := spid.FindAllString(string(out),1)
            logs.Info("Suricata is running -> " + string(out))
            logs.Info("Suricata PID -> %s", pid[0])
            return true
        }
    }
    logs.Error("Suricata isn't running " + string(out))
    return false
}

func Installed() (isIt bool){
    var suricata bool
    suricata = false
    suricata = suriPath()
    suricata = suriBin()
    suricata = suriRunning()
    if suricata {
        logs.Info("Suricata installed and running")
        return true
    } else {
        logs.Error("Suricata isn't present or not running")
    }
    return false
}

func GetBPF()(currentBPF string) {
    utils.GetConf("bpfPath")
    return ""
}

func SetBPF(n map[string]string)(bpf string, err error) {
    //utils.GetConf("bpfPath")
    logs.Info("log de N[BPF]-- "+n["bpf"])
    err = utils.BackupFile("/etc/owlh/suricata/", "filter.bpf")
    if err != nil{
        return "",err    
    }

    textbpf := n["bpf"]

    err = utils.UpdateBPFFile("/etc/owlh/suricata/", "filter.bpf", textbpf)
    if err != nil{
        return "",err    
    }

    return bpf, nil
}