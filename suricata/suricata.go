package suricata

import (
    "github.com/astaxie/beego/logs"
    "os"
    "os/exec"
    "strings"
    "regexp"
    "owlhnode/utils"
    // "io/ioutil"
)

func suriPath() (exists bool) {
    //Retrieve path for suricata.
	loadDatasuriPath := map[string]map[string]string{}
	loadDatasuriPath["suriPath"] = map[string]string{}
	loadDatasuriPath["suriPath"]["path"] = ""
    loadDatasuriPath = utils.GetConf(loadDatasuriPath)    
    path := loadDatasuriPath["suriPath"]["path"]

    if _, err := os.Stat(path); os.IsNotExist(err) {
        logs.Error("Suricata not installed, at least folder /etc/suricata dosn't exist")
        return false
    }
    return true
}

func suriBin() (exists bool) {
    //Retrieve path for suricata.
	loadDatasuriBin := map[string]map[string]string{}
	loadDatasuriBin["suriBin"] = map[string]string{}
    loadDatasuriBin["suriBin"]["cmd"] = ""
    loadDatasuriBin["suriBin"]["param"] = ""
    loadDatasuriBin = utils.GetConf(loadDatasuriBin)    
    cmd := loadDatasuriBin["suriBin"]["cmd"]
    param := loadDatasuriBin["suriBin"]["param"]

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

func suriRunning() (running bool) {
    //Retrieve path for suricata.
	loadDatasuriRunning := map[string]map[string]string{}
	loadDatasuriRunning["suriRunning"] = map[string]string{}
    loadDatasuriRunning["suriRunning"]["cmd"] = ""
    loadDatasuriRunning["suriRunning"]["param"] = ""
    loadDatasuriRunning["suriRunning"]["command"] = ""
    loadDatasuriRunning = utils.GetConf(loadDatasuriRunning)    
    cmd := loadDatasuriRunning["suriRunning"]["cmd"]
    param := loadDatasuriRunning["suriRunning"]["param"]
    command := loadDatasuriRunning["suriRunning"]["command"]

    //cmd := "ps -ef | grep suricata | grep -v grep | grep -v sudo | awk '{print $8 \" \" $2}' "
    //out, err := exec.Command("bash", "-c", cmd).Output()
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
/*
func GetBPF()(currentBPF string) {
    utils.GetConf("bpfPath")
    return ""
}
*/
func SetBPF(n map[string]string)(bpf string, err error) {
    //read path
    logs.Info("Set Suricata BPF -- Making Map")

	loadData := map[string]map[string]string{}
	loadData["suricataBPF"] = map[string]string{}
	loadData["suricataBPF"]["pathBPF"] = ""
	loadData["suricataBPF"]["fileBPF"] = "" 
    loadData = utils.GetConf(loadData)    

    path := loadData["suricataBPF"]["pathBPF"]
    file := loadData["suricataBPF"]["fileBPF"]

    //make backup file
    err = utils.BackupFile(path, file)
    if err != nil{
        return "",err    
    }

    //write bpf into the file
    textbpf := n["bpf"]
    err = utils.UpdateBPFFile(path, file, textbpf)
    if err != nil{
        return "",err    
    }

    return bpf, nil
}

//Retrieve data, make a backup file and write the new data on the original file
func RetrieveFile(file map[string][]byte)(err error){
    fileRetrieved := file["data"]
    path := "/etc/owlh/suricata/ruleset/"
    fileToEdit := "owlh.rules"
    
    err = utils.BackupFile(path, fileToEdit)
    if err != nil{
        return err    
    }
    
    err = utils.WriteNewDataOnFile(path+fileToEdit, fileRetrieved)
    if err != nil{
        return err    
    }

    return nil
}