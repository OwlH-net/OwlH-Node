package group

import (
    "github.com/astaxie/beego/logs"
	"owlhnode/database"
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
				cmd := exec.Command("suricata", "-D", "-c", group[x]["configFile"], "-i", group[x]["interface"], "-F",  group[x]["BPFfile"])
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