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
	
	out,err := exec.Command("bash", "-c", deployElement).Output()
	// err = cpCmd.Run()
	
	logs.Notice(string(out))

	if err != nil{logs.Error("utils.BackupFullPath Error exec cmd command: "+err.Error()); return err}
	
    return nil
}

func CheckDeployFiles()(anode map[string]string){
	loadData := map[string]map[string]string{}
	check := map[string]string{}
	loadData["deploy"] = map[string]string{}
	loadData["deploy"]["suricata"] = ""
	loadData["deploy"]["zeek"] = ""
	loadData["deploy"]["owlhmoloch"] = ""
	loadData["deploy"]["owlhinterface"] = ""
	loadData["deploy"]["owlhfirewall"] = ""
	loadData,err := utils.GetConf(loadData)
	if err != nil { logs.Error("SendFile Error getting data from main.conf"); return nil}
	suricata := loadData["deploy"]["suricata"]
	zeek := loadData["deploy"]["zeek"]
	moloch := loadData["deploy"]["owlhmoloch"]
	iface := loadData["deploy"]["owlhinterface"]
	firewall := loadData["deploy"]["owlhfirewall"]

	if _, err := os.Stat(suricata); os.IsNotExist(err) { check["suricata"]="false" } else { check["suricata"]="true" }
	if _, err := os.Stat(zeek); os.IsNotExist(err) { check["zeek"]="false" } else { check["zeek"]="true" }
	if _, err := os.Stat(moloch); os.IsNotExist(err) { check["moloch"]="false" } else { check["moloch"]="true" }
	if _, err := os.Stat(iface); os.IsNotExist(err) { check["interface"]="false" } else { check["interface"]="true" }
	if _, err := os.Stat(firewall); os.IsNotExist(err) { check["firewall"]="false" } else { check["firewall"]="true" }

	return check
}