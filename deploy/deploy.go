package deploy

import (
    "github.com/astaxie/beego/logs"
    "owlhnode/utils"
    "os/exec"
    "os"
)

func DeployNode(anode map[string]string)(err error) {
    loadData := map[string]map[string]string{}
    loadData["deploy"] = map[string]string{}
    loadData["deploy"][anode["value"]] = ""
    loadData,err = utils.GetConf(loadData)
    if err != nil { logs.Error("SendFile Error getting data from main.conf"); return err}
    deployElement := loadData["deploy"][anode["value"]]
    
    _,err = exec.Command("bash", "-c", deployElement).Output()
    if err != nil{logs.Error("utils.BackupFullPath Error exec cmd command: "+err.Error()); return err}
    
    return nil
}

func CheckDeployFiles()(anode map[string]string){
    loadData := map[string]map[string]string{}
    check := map[string]string{}
    loadData["deploy"] = map[string]string{}
    loadData["deploy"]["suricata"] = ""
    loadData["deploy"]["zeek"] = ""
    loadData["deploy"]["moloch"] = ""
    loadData["deploy"]["interface"] = ""
    loadData["deploy"]["firewall"] = ""
    loadData,err := utils.GetConf(loadData)
    if err != nil { logs.Error("SendFile Error getting data from main.conf"); return nil}

    for x:= range loadData["deploy"]{
        if _, err := os.Stat(loadData["deploy"][x]); os.IsNotExist(err) { check[x]="false" } else { check[x]="true" }
    }

    return check
}