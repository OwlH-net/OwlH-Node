package zeek

import (
    "github.com/astaxie/beego/logs"
    "os"
    "os/exec"
    "strings"
    // "regexp"
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
    // loadDataZeekBin["loadDataZeekBin"]["cmd"] = ""
    // loadDataZeekBin["loadDataZeekBin"]["param"] = ""
    loadDataZeekBin["loadDataZeekBin"]["bin"] = ""
    loadDataZeekBin = utils.GetConf(loadDataZeekBin)    
    // cmd := loadDataZeekBin["loadDataZeekBin"]["cmd"]
    // param := loadDataZeekBin["loadDataZeekBin"]["param"]
    bin := loadDataZeekBin["loadDataZeekBin"]["bin"]

    //out, err := exec.Command("broctl","-V").Output()
    // out, err := exec.Command(cmd,param).Output()
    _,err := os.Stat(bin)
    if err != nil {
        logs.Info("Zeek OS path err")
        return false
    }
    if os.IsNotExist(err){
        logs.Info("Zeek path not exist")
        return false
    }
    // if err == nil {
    //     if strings.Contains(string(out), "Zeek version") {
    //         logs.Info("Zeek binario existe -> " + string(out))
    //         return true
    //     }
    // }
    logs.Error("Zeek bin exist")
    return true
}

func ZeekRunning() (running bool) {
    //Retrieve running for zeek.
	loadDataZeekRunning := map[string]map[string]string{}
	loadDataZeekRunning["loadDataZeekRunning"] = map[string]string{}
    loadDataZeekRunning["loadDataZeekRunning"]["cmd"] = ""
    loadDataZeekRunning["loadDataZeekRunning"]["param"] = ""
    loadDataZeekRunning["loadDataZeekRunning"]["command"] = ""
    loadDataZeekRunning = utils.GetConf(loadDataZeekRunning)    
    cmd := loadDataZeekRunning["loadDataZeekRunning"]["cmd"]
    param := loadDataZeekRunning["loadDataZeekRunning"]["param"]
    command := loadDataZeekRunning["loadDataZeekRunning"]["command"]

    //cmd := "ps -ef | grep zeek | grep -v grep | grep -v sudo | awk '{print $8 \" \" $2}' "
    //out, err := exec.Command("bash", "-c", cmd).Output()
    out, err := exec.Command(command, param, cmd).Output()
    if err == nil {
        if strings.Contains(string(out), "running") {
            // spid := regexp.MustCompile("[0-9]+")
            // pid := spid.FindAllString(string(out),1)
            // logs.Info("Zeek is on execution -> " + string(out))
            // logs.Info("Zeek PID exec -> %s", pid[0])
            logs.Info("Zeek is running --> "+string(out))
            return true
        }
    }
    logs.Error("Zeek is NOT running -> " + string(out))
    return false
}

func Installed() (isIt map[string]bool){
    zeek := make(map[string]bool)
    //zeek = false
    zeek["path"] = ZeekPath()
    zeek["bin"] = ZeekBin()
    zeek["running"] = ZeekRunning()
    logs.Warn("ZEEK --> ")
    logs.Warn(zeek)
    if zeek["Path"] || zeek["Bin"] || zeek["Running"]  {
        logs.Info("Zeek installed and running")
        return zeek
    } else {
        logs.Error("Zeek isn't present or not running")
        return zeek
    }
}

//Run zeek
func RunZeek()(data string, err error){

    // //Retrieve path for RunZeek.
    StartZeek := map[string]map[string]string{}
    StartZeek["zeekStart"] = map[string]string{}
    StartZeek["zeekStart"]["start"] = ""
    StartZeek["zeekStart"]["param"] = ""
    StartZeek["zeekStart"]["command"] = ""
    StartZeek = utils.GetConf(StartZeek)    
    cmd := StartZeek["zeekStart"]["start"]
    param := StartZeek["zeekStart"]["param"]
    command := StartZeek["zeekStart"]["command"]

    logs.Info("Loading parameters Node GetConf")
    logs.Info("Loading parameters Node GetConf")
    logs.Info("Loading parameters Node GetConf")
    logs.Info("Loading parameters Node GetConf")
    logs.Info("Loading parameters Node GetConf")
    logs.Info(cmd)
    logs.Info(param)
    logs.Info(command)

    out,err := exec.Command(command, param, cmd).Output()
    logs.Info(string(out))
    if err != nil {
        logs.Error("Error launching zeekStart: "+err.Error())
        return "",err
    }
    return "zeekStart system is on",nil
}

//Stop zeek
func StopZeek()(data string, err error){

    // //Retrieve path for zeek.
    StopZeek := map[string]map[string]string{}
	StopZeek["zeekStop"] = map[string]string{}
    StopZeek["zeekStop"]["stop"] = ""
    StopZeek["zeekStop"]["param"] = ""
    StopZeek["zeekStop"]["command"] = ""
    StopZeek = utils.GetConf(StopZeek)    
    cmd := StopZeek["zeekStop"]["stop"]
    param := StopZeek["zeekStop"]["param"]
    command := StopZeek["zeekStop"]["command"]

    _,err = exec.Command(command, param, cmd).Output()
    if err != nil {
        logs.Error("Error stopping zeek: "+err.Error())
        return "",err
    }
    return "Zeek stopped ",nil
}