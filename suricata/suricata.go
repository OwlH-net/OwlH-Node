package suricata

import (
    "github.com/astaxie/beego/logs"
    "os"
    "os/exec"
    "strings"
    "regexp"
	"owlhnode/utils"
	"errors"
    // "fmt"
    // "io/ioutil"
)

//Retrieve suricata path from main.conf
func suriPath() (exists bool) {
	var err error
	loadDatasuriPath := map[string]map[string]string{}
	loadDatasuriPath["suriPath"] = map[string]string{}
	loadDatasuriPath["suriPath"]["path"] = ""
    loadDatasuriPath,err = utils.GetConf(loadDatasuriPath)    
	path := loadDatasuriPath["suriPath"]["path"]
	if err != nil {
		logs.Error("suriPath Error getting data from main.conf")
	}

    if _, err := os.Stat(path); os.IsNotExist(err) {
        logs.Error("Suricata not installed, at least folder /etc/suricata dosn't exist")
        return false
    }
    return true
}

//Retrieve suricata binary files path from main.conf
func suriBin() (exists bool) {
	var err error
	loadDatasuriBin := map[string]map[string]string{}
	loadDatasuriBin["suriBin"] = map[string]string{}
    loadDatasuriBin["suriBin"]["cmd"] = ""
    loadDatasuriBin["suriBin"]["param"] = ""
    loadDatasuriBin,err = utils.GetConf(loadDatasuriBin)    
    cmd := loadDatasuriBin["suriBin"]["cmd"]
    param := loadDatasuriBin["suriBin"]["param"]
	if err != nil {
		logs.Error("suriBin Error getting data from main.conf")
	}

    out, err := exec.Command(cmd,param).Output()
    if err == nil {
        if strings.Contains(string(out), "Suricata version") {
            logs.Info("Suricata installed -> " + string(out))
            return true
        }
    }
    logs.Error("Suricata Suricata not installed")
    return false
}

//Check whether Suricata is running
func suriRunning() (running bool) {
	var err error
	loadDatasuriRunning := map[string]map[string]string{}
	loadDatasuriRunning["suriRunning"] = map[string]string{}
    loadDatasuriRunning["suriRunning"]["cmd"] = ""
    loadDatasuriRunning["suriRunning"]["param"] = ""
    loadDatasuriRunning["suriRunning"]["command"] = ""
    loadDatasuriRunning,err = utils.GetConf(loadDatasuriRunning)    
    cmd := loadDatasuriRunning["suriRunning"]["cmd"]
    param := loadDatasuriRunning["suriRunning"]["param"]
    command := loadDatasuriRunning["suriRunning"]["command"]
	if err != nil {
		logs.Error("suriRunning Error getting data from main.conf")
	}

	out, err := exec.Command(command, param, cmd).Output()
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

//if suricata, path and bin are true, then suricata is installed and running
func Installed() (isIt map[string]bool, err error){
    suricata := make(map[string]bool)
    //suricata = false
    suricata["path"] = suriPath()
    suricata["bin"] = suriBin()
    suricata["running"] = suriRunning()
    
    logs.Info("SURICATA")
    logs.Info(suricata)

    if suricata["path"] || suricata["bin"] || suricata["running"] {
        logs.Info("Suricata installed and running")
        return suricata, nil
    } else {
        logs.Error("Suricata isn't present or not running")
        return suricata, errors.New("Suricata isn't present or not running")
    }   
}
/*
func GetBPF()(currentBPF string) {
    utils.GetConf("bpfPath")
    return ""
}
*/

//set BPF for suricata
func SetBPF(n map[string]string)(bpf string, err error) {
    logs.Info("Set Suricata BPF -- Making Map")
	loadData := map[string]map[string]string{}
	loadData["suricataBPF"] = map[string]string{}
	loadData["suricataBPF"]["pathBPF"] = ""
	loadData["suricataBPF"]["fileBPF"] = "" 
    loadData,err = utils.GetConf(loadData)    
    path := loadData["suricataBPF"]["pathBPF"]
    file := loadData["suricataBPF"]["fileBPF"]
	if err != nil {
		logs.Error("SetBPF Error getting data from main.conf")
	}

    //make backup file
    err = utils.BackupFile(path, file)
    if err != nil{
		logs.Error("Error creating BPF backup...")
        return "",err    
    }

    //write bpf into the file
    textbpf := n["bpf"]
    err = utils.UpdateBPFFile(path, file, textbpf)
    if err != nil{
		logs.Error("Error writting data into BPF file...")
        return "",err    
    }
    return bpf, nil
}

//Retrieve data, make a backup file and write the new data on the original file
func RetrieveFile(file map[string][]byte)(err error){
    fileRetrieved := file["data"]
    path := "/etc/owlh/suricata/ruleset/"
    fileToEdit := "owlh.rules"
	
	//create owlh.rules backup
    err = utils.BackupFile(path, fileToEdit)
    if err != nil{
		logs.Error("Error creating owlh.rules backup...")
        return err    
    }
	
	//write new data into owlh.rules file
    err = utils.WriteNewDataOnFile(path+fileToEdit, fileRetrieved)
    if err != nil{
		logs.Error("Error writting data into owlh.rules file...")
        return err    
    }
    return nil
}

//Run suricata
func RunSuricata()(data string, err error){
    StartSuricata := map[string]map[string]string{}
    StartSuricata["suriStart"] = map[string]string{}
    StartSuricata["suriStart"]["start"] = ""
    StartSuricata["suriStart"]["param"] = ""
    StartSuricata["suriStart"]["command"] = ""
    StartSuricata,err = utils.GetConf(StartSuricata)    
    cmd := StartSuricata["suriStart"]["start"]
    param := StartSuricata["suriStart"]["param"]
    command := StartSuricata["suriStart"]["command"]
	if err != nil {
		logs.Error("RunSuricata Error getting data from main.conf")
	}

    out,err := exec.Command(command, param, cmd).Output()
    logs.Info(string(out))
    if err != nil {
        logs.Error("Error launching suricata: "+err.Error())
        return "",err
    }
    return "Suricata system is on",nil
}

//Stop suricata
func StopSuricata()(data string, err error){
    // //Retrieve path for suricata.
    StopSuricata := map[string]map[string]string{}
	StopSuricata["suriStop"] = map[string]string{}
    StopSuricata["suriStop"]["stop"] = ""
    StopSuricata["suriStop"]["param"] = ""
    StopSuricata["suriStop"]["command"] = ""
    StopSuricata,err = utils.GetConf(StopSuricata)    
    cmd := StopSuricata["suriStop"]["stop"]
    param := StopSuricata["suriStop"]["param"]
    command := StopSuricata["suriStop"]["command"]
	if err != nil {
		logs.Error("StopSuricata Error getting data from main.conf")
	}
	
    _,err = exec.Command(command, param, cmd).Output()
    if err != nil {
        logs.Error("Error stopping suricata: "+err.Error())
        return "",err
    }
    return "Suricata stopped ",nil
}