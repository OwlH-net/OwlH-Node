package group

import (
    "github.com/astaxie/beego/logs"
	"owlhnode/database"
	"owlhnode/utils"
	"os/exec"
)

func SyncSuricataGroupValues(data map[string]string) (err error) {
	//update suricata values from group
	err = ndb.UpdateSuricataGroupValue("suricata","name",data["name"]); if err != nil {logs.Error("SyncSuricataGroupValues name insert error: %s", err.Error());return err}
	err = ndb.UpdateSuricataGroupValue("suricata","configFile",data["configFile"]); if err != nil {logs.Error("SyncSuricataGroupValues configFile insert error: %s", err.Error());return err}
	err = ndb.UpdateSuricataGroupValue("suricata","BPFfile",data["BPFfile"]); if err != nil {logs.Error("SyncSuricataGroupValues BPFfile insert error: %s", err.Error());return err}
	err = ndb.UpdateSuricataGroupValue("suricata","interface",data["interface"]); if err != nil {logs.Error("SyncSuricataGroupValues interface insert error: %s", err.Error());return err}
	err = ndb.UpdateSuricataGroupValue("suricata","BPFrule",data["BPFrule"]); if err != nil {logs.Error("SyncSuricataGroupValues BPFrule insert error: %s", err.Error());return err}
	err = ndb.UpdateSuricataGroupValue("suricata","commandLine",data["commandLine"]); if err != nil {logs.Error("SyncSuricataGroupValues commandLine insert error: %s", err.Error());return err}

	return nil
}

func SuricataGroupService(data map[string]string) (err error) {
	group,err := ndb.GetAllGroupData()
	if err != nil {logs.Error("SuricataGroupService Error getting Suricata group values: %s", err.Error());return err}

	if data["action"] == "start"{
		for x,_ := range group{
			if x == data["uuid"]{

				suricata, err := utils.GetKeyValueString("suricata", "suricata")
				if err != nil { logs.Error("DeployNode Error getting data from main.conf"); return err}
				detail, err := utils.GetKeyValueString("execute", "detail")
				if err != nil { logs.Error("DeployNode Error getting data from main.conf"); return err}
				param, err := utils.GetKeyValueString("execute", "param")
				if err != nil { logs.Error("DeployNode Error getting data from main.conf"); return err}
				info, err := utils.GetKeyValueString("execute", "i")
				if err != nil { logs.Error("DeployNode Error getting data from main.conf"); return err}
				flag, err := utils.GetKeyValueString("execute", "flag")
				if err != nil { logs.Error("DeployNode Error getting data from main.conf"); return err}

				cmd := exec.Command(suricata, detail, param, group[x]["configFile"], info, group[x]["interface"], flag,  group[x]["BPFfile"])
				err = cmd.Run()
				if err != nil {
					logs.Error("group/SuricataGroupService error launching Suricata: "+err.Error());
					return err
				}
			}
		}
	}else{
		//stop suricata
	}

	return nil
}