package group

import (
    "github.com/astaxie/beego/logs"
    "os/exec"
    "errors"
    "strings"
    "owlhnode/database"
    "owlhnode/suricata"
    "owlhnode/utils"
)

func SyncSuricataGroupValues(data map[string]string) (err error) {
    //update suricata values from group
    err = ndb.UpdateSuricataGroupValue("suricata", "name", data["name"])
    if err != nil {
        logs.Error("SyncSuricataGroupValues name insert error: %s", err.Error())
        return err
    }
    err = ndb.UpdateSuricataGroupValue("suricata", "configFile", data["configFile"])
    if err != nil {
        logs.Error("SyncSuricataGroupValues configFile insert error: %s", err.Error())
        return err
    }
    err = ndb.UpdateSuricataGroupValue("suricata", "BPFfile", data["BPFfile"])
    if err != nil {
        logs.Error("SyncSuricataGroupValues BPFfile insert error: %s", err.Error())
        return err
    }
    err = ndb.UpdateSuricataGroupValue("suricata", "interface", data["interface"])
    if err != nil {
        logs.Error("SyncSuricataGroupValues interface insert error: %s", err.Error())
        return err
    }
    err = ndb.UpdateSuricataGroupValue("suricata", "BPFrule", data["BPFrule"])
    if err != nil {
        logs.Error("SyncSuricataGroupValues BPFrule insert error: %s", err.Error())
        return err
    }
    err = ndb.UpdateSuricataGroupValue("suricata", "commandLine", data["commandLine"])
    if err != nil {
        logs.Error("SyncSuricataGroupValues commandLine insert error: %s", err.Error())
        return err
    }

    return nil
}

func SuricataGroupService(data map[string]string) (err error) {
    group, err := ndb.GetAllGroupData()
    if err != nil {
        logs.Error("SuricataGroupService Error getting Suricata group values: %s", err.Error())
        return err
    }

    if data["action"] == "start" {
        for x := range group {
            if x == data["uuid"] {

                suricata, err := utils.GetKeyValueString("suricata", "suricata")
                if err != nil {
                    logs.Error("DeployNode Error getting data from main.conf")
                    return err
                }
                param, err := utils.GetKeyValueString("execute", "param")
                if err != nil {
                    logs.Error("DeployNode Error getting data from main.conf")
                    return err
                }

                cmd := exec.Command(suricata, "-D", param, group[x]["configFile"], "-i", group[x]["interface"], "-F", group[x]["BPFfile"])
                err = cmd.Run()
                if err != nil {
                    logs.Error("group/SuricataGroupService error launching Suricata: " + err.Error())
                    return err
                }
            }
        }
    } else {
        //stop suricata
    }

    return nil
}

func SyncGroupRulesetToNode(file map[string][]byte) (err error) {
    if file["data"] == nil || len(file["data"]) <= 0 {
        return errors.New("SyncGroupRulesetToNode error: Can't Synchronize empty ruleset")
    }

    path, err := utils.GetKeyValueString("suricataRuleset", "path")
    if err != nil {
        logs.Error("SyncGroupRulesetToNode Error getting data from main.conf: " + err.Error())
        return err
    }
    fileToEdit, err := utils.GetKeyValueString("suricataRuleset", "file")
    if err != nil {
        logs.Error("SyncGroupRulesetToNode Error getting data from main.conf: " + err.Error())
        return err
    }

    //replace file by name
    plug := strings.Replace(fileToEdit, "<NAME>", string(file["name"]), -1)

    //create owlh.rules backup
    err = utils.BackupFile(path, plug)
    if err != nil {
        logs.Error("Error creating owlh.rules backup: " + err.Error())
        return err
    }

    //write new data into owlh.rules file
    err = utils.WriteNewDataOnFile(path+plug, file["data"])
    if err != nil {
        logs.Error("Error writting data into owlh.rules file: " + err.Error())
        return err
    }
    // /usr/local/bin/suricatasc -c reload-rules /var/run/suricata/suricata-command.socket
    //SuricataRulesetReload
    if suricata.SuriRunning() {
        suricatasc, err := utils.GetKeyValueString("SuricataRulesetReload", "suricatasc")
        if err != nil {
            logs.Error("suriRunning Error getting data from main.conf: " + err.Error())
            return err
        }
        param, err := utils.GetKeyValueString("SuricataRulesetReload", "param")
        if err != nil {
            logs.Error("suriRunning Error getting data from main.conf: " + err.Error())
            return err
        }
        reloads, err := utils.GetKeyValueString("SuricataRulesetReload", "reload")
        if err != nil {
            logs.Error("suriRunning Error getting data from main.conf: " + err.Error())
            return err
        }
        socket, err := utils.GetKeyValueString("SuricataRulesetReload", "socket")
        if err != nil {
            logs.Error("suriRunning Error getting data from main.conf: " + err.Error())
            return err
        }

        _, err = exec.Command(suricatasc, param, reloads, socket).Output()
        if err != nil {
            logs.Error("Error executing command in SyncRulesetFromMaster function: " + err.Error())
            return err
        }
    }

    return nil
}