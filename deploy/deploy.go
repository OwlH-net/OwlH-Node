package deploy

import (
    "github.com/astaxie/beego/logs"
    "owlhnode/utils"
    "os/exec"
    "os"
)

func DeployNode(anode map[string]string)(err error) { 
    deployElement, err := utils.GetKeyValueString("deploy", anode["value"])
    if err != nil { logs.Error("DeployNode Error getting data from main.conf"); return err}

    _,err = exec.Command("bash", "-c", deployElement).Output()
    if err != nil{logs.Error("utils.BackupFullPath Error exec cmd command: "+err.Error()); return err}
    
    return nil
}

func CheckDeployFiles()(anode map[string]string){
    check := map[string]string{}
    
    suricata, err := utils.GetKeyValueString("deploy", "suricata"); if err != nil { logs.Error("CheckDeployFiles Error getting suricata from main.conf"); return nil}
    zeek, err := utils.GetKeyValueString("deploy", "zeek"); if err != nil { logs.Error("CheckDeployFiles Error getting zeek from main.conf"); return nil}
    moloch, err := utils.GetKeyValueString("deploy", "moloch"); if err != nil { logs.Error("CheckDeployFiles Error getting moloch from main.conf"); return nil}
    iface, err := utils.GetKeyValueString("deploy", "interface"); if err != nil { logs.Error("CheckDeployFiles Error getting interface from main.conf"); return nil}
    firewall, err := utils.GetKeyValueString("deploy", "firewall"); if err != nil { logs.Error("CheckDeployFiles Error getting firewall from main.conf"); return nil}

    if _, err := os.Stat(suricata); os.IsNotExist(err) { check["suricata"]="false" } else { check["suricata"]="true" }
    if _, err := os.Stat(zeek); os.IsNotExist(err) { check["zeek"]="false" } else { check["zeek"]="true" }
    if _, err := os.Stat(moloch); os.IsNotExist(err) { check["moloch"]="false" } else { check["moloch"]="true" }
    if _, err := os.Stat(iface); os.IsNotExist(err) { check["interface"]="false" } else { check["interface"]="true" }
    if _, err := os.Stat(firewall); os.IsNotExist(err) { check["firewall"]="false" } else { check["firewall"]="true" }

    return check
}