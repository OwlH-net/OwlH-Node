package zeek

import (
    "github.com/astaxie/beego/logs"
    "os"
    "os/exec"
    "strings"
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
	DeployZeek["zeekDeploy"] = map[string]string{}
    DeployZeek["zeekDeploy"]["cmd"] = ""
    DeployZeek["zeekDeploy"]["param"] = ""
    DeployZeek["zeekDeploy"]["command"] = ""
    DeployZeek,err = utils.GetConf(DeployZeek)    
    cmd := DeployZeek["zeekDeploy"]["cmd"]
    param := DeployZeek["zeekDeploy"]["param"]
    command := DeployZeek["zeekDeploy"]["command"]
	if err != nil {
		logs.Error("DeployZeek Error getting data from main.conf: "+err.Error())
	}
	// _,err = exec.Command(command, param, cmd).Output()
	_,err = exec.Command(cmd).Output()
    if err != nil {
        logs.Error("Error deploying zeek: "+err.Error())
        return err
    }
    return nil
}