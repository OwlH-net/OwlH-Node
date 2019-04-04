package zeek

import (
    "github.com/astaxie/beego/logs"
    "os"
    "os/exec"
    "strings"
    // "regexp"
	"owlhnode/utils"
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
		logs.Error("ZeekPath Error getting data from main.conf")
	}

    if _, err := os.Stat(path); os.IsNotExist(err) {
        logs.Error("Zeek is not installed on "+path+".")
        return false
    }
    return true
}

func ZeekBin() (exists bool) {
	var err error
    //Retrieve bin for wazuh.
	loadDataZeekBin := map[string]map[string]string{}
	loadDataZeekBin["loadDataZeekBin"] = map[string]string{}
    // loadDataZeekBin["loadDataZeekBin"]["cmd"] = ""
    // loadDataZeekBin["loadDataZeekBin"]["param"] = ""
    loadDataZeekBin["loadDataZeekBin"]["bin"] = ""
    loadDataZeekBin,err = utils.GetConf(loadDataZeekBin)    
    // cmd := loadDataZeekBin["loadDataZeekBin"]["cmd"]
    // param := loadDataZeekBin["loadDataZeekBin"]["param"]
    bin := loadDataZeekBin["loadDataZeekBin"]["bin"]
	if err != nil {
		logs.Error("ZeekBin Error getting data from main.conf")
	}

    //out, err := exec.Command("broctl","-V").Output()
    // out, err := exec.Command(cmd,param).Output()
    _,err = os.Stat(bin)
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
	}

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

func Installed() (isIt map[string]bool, err error){
    zeek := make(map[string]bool)
    //zeek = false
    zeek["path"] = ZeekPath()
    zeek["bin"] = ZeekBin()
    zeek["running"] = ZeekRunning()
    logs.Info("ZEEK --> ")
	logs.Info(zeek)

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

    // //Retrieve path for RunZeek.
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
		logs.Error("RunZeek Error getting data from main.conf")
	}

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
    StopZeek,err = utils.GetConf(StopZeek)    
    cmd := StopZeek["zeekStop"]["stop"]
    param := StopZeek["zeekStop"]["param"]
    command := StopZeek["zeekStop"]["command"]
	if err != nil {
		logs.Error("StopZeek Error getting data from main.conf")
	}
	
    _,err = exec.Command(command, param, cmd).Output()
    if err != nil {
        logs.Error("Error stopping zeek: "+err.Error())
        return "",err
    }
    return "Zeek stopped ",nil
}