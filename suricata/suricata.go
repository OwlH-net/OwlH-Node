package suricata

import (
    "fmt"
    "github.com/astaxie/beego/logs"
    "os"
    "os/exec"
    "owlhnode/database"
    "owlhnode/utils"
    "regexp"
    "strings"
    "time"
    // "encoding/json"
    "errors"
    "io/ioutil"
    "strconv"
    // "encoding/base64"
    // "crypto/sha256"
)

//Retrieve suricata path from main.conf
func suriPath() (exists bool) {
    var err error
    path, err := utils.GetKeyValueString("suriPath", "path")
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
    cmd, err := utils.GetKeyValueString("suriBin", "cmd")
    if err != nil {
        logs.Error("suriBin Error getting data from main.conf")
    }
    param, err := utils.GetKeyValueString("suriBin", "param")
    if err != nil {
        logs.Error("suriBin Error getting data from main.conf")
    }

    out, err := exec.Command(cmd, param).Output()
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
func SuriRunning() (running bool) {
    var err error
    cmd, err := utils.GetKeyValueString("suriRunning", "cmd")
    if err != nil {
        logs.Error("suriRunning Error getting data from main.conf")
    }
    param, err := utils.GetKeyValueString("suriRunning", "param")
    if err != nil {
        logs.Error("suriRunning Error getting data from main.conf")
    }
    command, err := utils.GetKeyValueString("suriRunning", "command")
    if err != nil {
        logs.Error("suriRunning Error getting data from main.conf")
    }

    out, err := exec.Command(command, param, cmd).Output()
    if err == nil {
        // if strings.Contains(string(out), "suricata") {
        spid := regexp.MustCompile("[0-9]+")
        pid := spid.FindAllString(string(out), 1)
        if len(pid) <= 0 || pid == nil {
            return false
        }
        logs.Info("Suricata PID -> %s", pid[0])
        return true
        // }
    }
    logs.Error("Suricata isn't running " + string(out))
    return false
}

//if suricata, path and bin are true, then suricata is installed and running
func Installed() (isIt map[string]bool, err error) {
    suricata := make(map[string]bool)
    suricata["path"] = suriPath()
    suricata["bin"] = suriBin()
    suricata["running"] = SuriRunning()
    if suricata["path"] || suricata["bin"] || suricata["running"] {
        logs.Info("Suricata installed and running")
        return suricata, nil
    } else {
        logs.Error("Suricata isn't present or not running")
        return suricata, errors.New("Suricata isn't present or not running")
    }
}

// func GetBPF()(bpf string, err error) {
//     logs.Info("Set Suricata BPF -- Making Map")
//     loadData := map[string]map[string]string{}
//     loadData["suricataBPF"] = map[string]string{}
//     loadData["suricataBPF"]["pathBPF"] = ""
//     loadData["suricataBPF"]["fileBPF"] = ""
//     loadData,err = utils.GetConf(loadData)
//     path := loadData["suricataBPF"]["pathBPF"]
//     file := loadData["suricataBPF"]["fileBPF"]
//     if err != nil {
//         logs.Error("GetBPF Error getting data from main.conf: "+err.Error())
//         return "",err
//     }

//     //read filter.bpf
//     bpfByte, err := ioutil.ReadFile(path+file) // just pass the file name
//     if err != nil {
//         logs.Error("GetBPF Error getting data from filter.bpf: "+err.Error())
//         return "",err
//     }
//     return string(bpfByte),nil
// }

//set BPF for suricata
func SetBPF(n map[string]string) (err error) {
    //historical log
    uuid := utils.Generate()
    currentTime := time.Now()
    timeFormated := currentTime.Format("2006-01-02T15:04:05")
    _ = ndb.InsertPluginCommand(uuid, "date", timeFormated)
    _ = ndb.InsertPluginCommand(uuid, "id", n["service"])
    _ = ndb.InsertPluginCommand(uuid, "type", "Suricata")
    _ = ndb.InsertPluginCommand(uuid, "action", "setBPF")
    _ = ndb.InsertPluginCommand(uuid, "description", "Save data into BPF file")
    path, err := utils.GetKeyValueString("suricataBPF", "pathBPF")
    if err != nil {
        logs.Error("SetBPF Error getting data from main.conf: " + err.Error())
        return err
    }
    file, err := utils.GetKeyValueString("suricataBPF", "fileBPF")
    if err != nil {
        logs.Error("SetBPF Error getting data from main.conf: " + err.Error())
        return err
    }

    //save bpf into specific suricata service database
    err = ndb.UpdatePluginValue(n["service"], "bpf", n["value"])
    if err != nil {
        logs.Error("SetBPF Error updating plugin.db: " + err.Error())
        return err
    }

    //check if exists
    if _, err = os.Stat(path + n["service"] + "-" + file); os.IsNotExist(err) {
        err = ioutil.WriteFile(path+n["service"]+"-"+file, []byte(n["value"]), 0644)
        if err != nil {
            _ = ndb.InsertPluginCommand(uuid, "status", "Error")
            _ = ndb.InsertPluginCommand(uuid, "output", err.Error())
            logs.Error("Error writting data into BPF file: " + err.Error())
            return err
        }
    } else {
        //make backup file
        err = utils.BackupFile(path, n["service"]+"-"+file)
        if err != nil {
            _ = ndb.InsertPluginCommand(uuid, "status", "Error")
            _ = ndb.InsertPluginCommand(uuid, "output", err.Error())
            logs.Error("Error creating BPF backup: " + err.Error())
            return err
        }

        //write bpf into the file
        err = utils.UpdateBPFFile(path, n["service"]+"-"+file, n["value"])
        if err != nil {
            _ = ndb.InsertPluginCommand(uuid, "status", "Error")
            _ = ndb.InsertPluginCommand(uuid, "output", err.Error())
            logs.Error("Error UpdateBPFFile: " + err.Error())
            return err
        }
    }
    _ = ndb.InsertPluginCommand(uuid, "status", "Success")
    _ = ndb.InsertPluginCommand(uuid, "output", "New BPF set successfully")
    return nil
}

//Retrieve data, make a backup file and write the new data on the original file
func SyncRulesetFromMaster(file map[string][]byte) (err error) {

    //historical log
    uuid := utils.Generate()
    currentTime := time.Now()
    timeFormated := currentTime.Format("2006-01-02T15:04:05")
    _ = ndb.InsertPluginCommand(uuid, "date", timeFormated)
    _ = ndb.InsertPluginCommand(uuid, "type", "Suricata")
    _ = ndb.InsertPluginCommand(uuid, "action", "SyncRulesetFromMaster")
    _ = ndb.InsertPluginCommand(uuid, "description", "Sync ruleset from master")

    if file["data"] == nil || len(file["data"]) <= 0 {
        return errors.New("SyncRulesetFromMaster error: Can't Synchronize empty ruleset")
    }

    path, err := utils.GetKeyValueString("suricataRuleset", "path")
    if err != nil {
        logs.Error("SyncRulesetFromMaster Error getting data from main.conf: " + err.Error())
        return err
    }
    fileToEdit, err := utils.GetKeyValueString("suricataRuleset", "file")
    if err != nil {
        logs.Error("SyncRulesetFromMaster Error getting data from main.conf: " + err.Error())
        return err
    }

    //replace file by name
    fileName := strings.Replace(string(file["rulesetName"]), " ", "-", -1)
    plug := strings.Replace(fileToEdit, "<NAME>", fileName, -1)
    logs.Info("SURICATA - RULESET SYNC -> ruleset name -> %s, full ruleset name -> %s", fileName, plug)

    //create owlh.rules backup
    err = utils.BackupFile(path, plug)
    if err != nil {
        _ = ndb.InsertPluginCommand(uuid, "status", "Error")
        _ = ndb.InsertPluginCommand(uuid, "output", "SyncRulesetFromMaster Error doing Backup: "+err.Error())
        logs.Error("SURICATA - RULESET SYNC -> Error creating ruleset %s backup: %s", plug, err.Error())
        return err
    }

    //write new data into owlh.rules file
    err = utils.WriteNewDataOnFile(path+plug, file["data"])
    if err != nil {
        _ = ndb.InsertPluginCommand(uuid, "status", "Error")
        _ = ndb.InsertPluginCommand(uuid, "output", "SyncRulesetFromMaster Error writing new ruleset data into file: "+err.Error())
        logs.Error("SURICATA - RULESET SYNC -> Error writting data into  file: " + err.Error())
        return err
    }

    // /usr/local/bin/suricatasc -c reload-rules /var/run/suricata/suricata-command.socket
    //SuricataRulesetReload
    logs.Info("SURICATA - RULSET SYNC -> please reload suricatas runing ruleset -> %s", fileName)
    ReloadSuricatas(fileName)

    _ = ndb.InsertPluginCommand(uuid, "status", "Success")
    _ = ndb.InsertPluginCommand(uuid, "output", "Ruleset sync successfully")

    //update plugin sync status
    // err = ndb.UpdatePluginValue()

    return nil
}

func ReloadSuricatas(rulesetName string) (err error) {

    fullpidfile, err := utils.GetKeyValueString("suricata", "fullpidfile")
    if err != nil {
        logs.Error("SyncRulesetFromMaster Error getting fullpidfile data from main.conf: " + err.Error())
        return err
    }
    suricatasc, err := utils.GetKeyValueString("SuricataRulesetReload", "suricatasc")
    if err != nil {
        logs.Error("SyncRulesetFromMaster Error getting data from main.conf: " + err.Error())
        return err
    }
    param, err := utils.GetKeyValueString("SuricataRulesetReload", "param")
    if err != nil {
        logs.Error("SyncRulesetFromMaster Error getting data from main.conf: " + err.Error())
        return err
    }
    reloads, err := utils.GetKeyValueString("SuricataRulesetReload", "reload")
    if err != nil {
        logs.Error("SyncRulesetFromMaster Error getting data from main.conf: " + err.Error())
        return err
    }

    reloadSignal, err := utils.GetKeyValueString("suricata", "reloadsignal")
    if err != nil {
        logs.Error("SyncRulesetFromMaster Error getting suricata-reloadsignal data from main.conf: " + err.Error())
        return err
    }

    Suricatas, _ := GetSuricataServices()
    logs.Info("Suricata Services -> %+v", Suricatas)

    for eachSuricata := range Suricatas {
        if Suricatas[eachSuricata]["localRulesetName"] == rulesetName {
            logs.Info("SURICATA - RULESET SYNC -> reload uuid %s", eachSuricata)
            pidfile := strings.Replace(fullpidfile, "<ID>", eachSuricata, -1)
            if SuriRunning() {
                logs.Info("SURICATA - RULESET SYNC -> reload Suricata name -> %s", Suricatas[eachSuricata]["name"])
                logs.Info("SURICATA - RULESET SYNC -> sc - %s, param - %s, reload command - %s, pidfile - %s", suricatasc, param, reloads, pidfile)
                cmdOutput, err := exec.Command(suricatasc, param, reloads, pidfile).Output()
                if err != nil {
                    logs.Warning("SURICATA - RULESET SYNC -> Can't reload Suricata using SC command: %s (%s)", cmdOutput, err.Error())
                    cmdOutput, err = exec.Command("kill", reloadSignal, Suricatas[eachSuricata]["pid"]).Output()
                    logs.Info("SURICATA - RULESET SYNC -> running kill signal to reload pid -> %s", Suricatas[eachSuricata]["pid"])
                    if err != nil {
                        logs.Warning("SURICATA - RULESET SYNC -> Can't reload Suricata using kill signal: %s (%s)", cmdOutput, err.Error())
                        return err
                    }
                }
            }
        }
    }
    return nil
}

//Run suricata
func RunSuricata() (data string, err error) {
    // Deprecated use StartSuricataMainConf
    var ddata = map[string]string{}
    err = StartSuricataMainConf(ddata)
    return "", err

    // cmd, err := utils.GetKeyValueString("suriStart", "start")
    // if err != nil {
    //     logs.Error("RunSuricata Error getting data from main.conf: " + err.Error())
    //     return "", err
    // }
    // param, err := utils.GetKeyValueString("suriStart", "param")
    // if err != nil {
    //     logs.Error("RunSuricata Error getting data from main.conf: " + err.Error())
    //     return "", err
    // }
    // command, err := utils.GetKeyValueString("suriStart", "command")
    // if err != nil {
    //     logs.Error("RunSuricata Error getting data from main.conf: " + err.Error())
    //     return "", err
    // }

    // _, err = exec.Command(command, param, cmd).Output()
    // if err != nil {
    //     logs.Error("Error launching suricata: " + err.Error())
    //     return "", err
    // }
    // return "Suricata system is ON!", nil
}

//Stop suricata
func StopSuricata() (data string, err error) {
    //historical log
    uuid := utils.Generate()
    currentTime := time.Now()
    timeFormated := currentTime.Format("2006-01-02T15:04:05")
    _ = ndb.InsertPluginCommand(uuid, "date", timeFormated)
    _ = ndb.InsertPluginCommand(uuid, "type", "Suricata")
    _ = ndb.InsertPluginCommand(uuid, "action", "StopSuricata")
    _ = ndb.InsertPluginCommand(uuid, "description", "Stop suricata")

    // //Retrieve path for suricata.
    cmd, err := utils.GetKeyValueString("suriStop", "stop")
    if err != nil {
        logs.Error("StopSuricata Error getting data from main.conf")
    }
    param, err := utils.GetKeyValueString("suriStop", "param")
    if err != nil {
        logs.Error("StopSuricata Error getting data from main.conf")
    }
    command, err := utils.GetKeyValueString("suriStop", "command")
    if err != nil {
        logs.Error("StopSuricata Error getting data from main.conf")
    }

    _, err = exec.Command(command, param, cmd).Output()
    if err != nil {
        _ = ndb.InsertPluginCommand(uuid, "status", "Error")
        _ = ndb.InsertPluginCommand(uuid, "output", "SyncRulesetFromMaster restarting Suricata: "+err.Error())
        _ = ndb.InsertPluginCommand(uuid, "command", command+param+cmd)
        logs.Error("Error stopping suricata: " + err.Error())
        return "", err
    }

    _ = ndb.InsertPluginCommand(uuid, "status", "Success")
    _ = ndb.InsertPluginCommand(uuid, "output", "Suricata is stopped successfully")
    return "Suricata stopped ", nil
}

func GetSuricataServices() (data map[string]map[string]string, err error) {
    data, err = ndb.GetServices("suricata")
    if err != nil {
        logs.Error("GetSuricataServices Reading services Error: " + err.Error())
        return nil, err
    }
    return data, nil
}

func SaveConfigFile(files map[string][]byte) (err error) {

    for x := range files {
        if _, err := os.Stat(x); os.IsNotExist(err) {
            os.MkdirAll(x, os.ModePerm)
        }

        //write and create tar.gzip file
        err = ioutil.WriteFile(x+"/file.tar.gzip", files[x], 0644)
        if err != nil {
            logs.Error("SaveConfigFile Error: " + err.Error())
            return err
        }
        logs.Info("config tar file written to %s", x)
        //unzip
        err = utils.ExtractFile(x+"/file.tar.gzip", x)
        if err != nil {
            logs.Error("SaveConfigFile ExtractFile Error: " + err.Error())
            return err
        }
        //remove zip file
        os.Remove(x + "/file.tar.gzip")
    }
    return nil
}

func StartSuricataMainConf(anode map[string]string) (err error) {
    //historical log
    uuid := utils.Generate()
    currentTime := time.Now()
    timeFormated := currentTime.Format("2006-01-02T15:04:05")
    _ = ndb.InsertPluginCommand(uuid, "date", timeFormated)
    _ = ndb.InsertPluginCommand(uuid, "type", "Suricata")
    _ = ndb.InsertPluginCommand(uuid, "action", "StartSuricataMainConf")
    _ = ndb.InsertPluginCommand(uuid, "description", "Start Suricata mainconf")

    cmd, err := utils.GetKeyValueString("suricata", "start")
    if err != nil {
        logs.Error("StartSuricataMainConf Error getting data from main.conf")
    }
    param, err := utils.GetKeyValueString("suricata", "param")
    if err != nil {
        logs.Error("StartSuricataMainConf Error getting data from main.conf")
    }
    command, err := utils.GetKeyValueString("suricata", "command")
    if err != nil {
        logs.Error("StartSuricataMainConf Error getting data from main.conf")
    }

    // err = utils.RunCommand(cmd, param)
    _, err = exec.Command(command, param, cmd).Output()
    if err != nil {
        _ = ndb.InsertPluginCommand(uuid, "status", "Error")
        _ = ndb.InsertPluginCommand(uuid, "output", "StartSuricataMainConf start Suricata mainconf error: "+err.Error())
        _ = ndb.InsertPluginCommand(uuid, "command", command+param+cmd)
        logs.Error("StartSuricataMainConf/Error starting suricata: " + err.Error())
        return err
    }

    return nil
}
func StopSuricataMainConf(anode map[string]string) (err error) {
    //historical log
    uuid := utils.Generate()
    currentTime := time.Now()
    timeFormated := currentTime.Format("2006-01-02T15:04:05")
    _ = ndb.InsertPluginCommand(uuid, "date", timeFormated)
    _ = ndb.InsertPluginCommand(uuid, "type", "Suricata")
    _ = ndb.InsertPluginCommand(uuid, "action", "StopSuricataMainConf")
    _ = ndb.InsertPluginCommand(uuid, "description", "Stop Suricata mainconf")

    cmd, err := utils.GetKeyValueString("suricata", "stop")
    if err != nil {
        logs.Error("StopSuricataMainConf Error getting data from main.conf")
    }
    param, err := utils.GetKeyValueString("suricata", "param")
    if err != nil {
        logs.Error("StopSuricataMainConf Error getting data from main.conf")
    }
    command, err := utils.GetKeyValueString("suricata", "command")
    if err != nil {
        logs.Error("StopSuricataMainConf Error getting data from main.conf")
    }

    _, err = exec.Command(command, param, cmd).Output()
    if err != nil {
        _ = ndb.InsertPluginCommand(uuid, "status", "Error")
        _ = ndb.InsertPluginCommand(uuid, "output", "StopSuricataMainConf stop Suricata mainconf error: "+err.Error())
        _ = ndb.InsertPluginCommand(uuid, "command", command+param+cmd)
        logs.Error("StopSuricataMainConf/Error stopping suricata: " + err.Error())
        return err
    }
    _ = ndb.InsertPluginCommand(uuid, "status", "Success")
    _ = ndb.InsertPluginCommand(uuid, "output", "StopSuricataMainConf stopped successfully")
    return nil
}
func KillSuricataMainConf(anode map[string]string) (err error) {
    //historical log
    uuid := utils.Generate()
    currentTime := time.Now()
    timeFormated := currentTime.Format("2006-01-02T15:04:05")
    _ = ndb.InsertPluginCommand(uuid, "date", timeFormated)
    _ = ndb.InsertPluginCommand(uuid, "type", "Suricata")
    _ = ndb.InsertPluginCommand(uuid, "action", "KillSuricataMainConf")
    _ = ndb.InsertPluginCommand(uuid, "description", "Kill Suricata mainconf")

    cmd, err := utils.GetKeyValueString("suricata", "kill")
    if err != nil {
        logs.Error("KillSuricataMainConf Error getting data from main.conf")
    }
    param, err := utils.GetKeyValueString("suricata", "param")
    if err != nil {
        logs.Error("KillSuricataMainConf Error getting data from main.conf")
    }
    command, err := utils.GetKeyValueString("suricata", "command")
    if err != nil {
        logs.Error("KillSuricataMainConf Error getting data from main.conf")
    }

    _, err = exec.Command(command, param, cmd+" "+anode["pid"]).Output()
    if err != nil {
        _ = ndb.InsertPluginCommand(uuid, "status", "Error")
        _ = ndb.InsertPluginCommand(uuid, "output", "KillSuricataMainConf kill Suricata mainconf error: "+err.Error())
        _ = ndb.InsertPluginCommand(uuid, "command", command+param+cmd+" "+anode["pid"])
        logs.Error("KillSuricataMainConf/Error starting suricata from main conf: " + err.Error())
        return err
    }
    _ = ndb.InsertPluginCommand(uuid, "status", "Success")
    _ = ndb.InsertPluginCommand(uuid, "output", "KillSuricataMainConf killed successfully")
    return nil
}
func ReloadSuricataMainConf(anode map[string]string) (err error) {
    //historical log
    uuid := utils.Generate()
    currentTime := time.Now()
    timeFormated := currentTime.Format("2006-01-02T15:04:05")
    _ = ndb.InsertPluginCommand(uuid, "date", timeFormated)
    _ = ndb.InsertPluginCommand(uuid, "type", "Suricata")
    _ = ndb.InsertPluginCommand(uuid, "action", "ReloadSuricataMainConf")
    _ = ndb.InsertPluginCommand(uuid, "description", "Reload Suricata mainconf")

    cmd, err := utils.GetKeyValueString("suricata", "reload")
    if err != nil {
        logs.Error("ReloadSuricataMainConf Error getting data from main.conf")
    }
    param, err := utils.GetKeyValueString("suricata", "param")
    if err != nil {
        logs.Error("ReloadSuricataMainConf Error getting data from main.conf")
    }
    command, err := utils.GetKeyValueString("suricata", "command")
    if err != nil {
        logs.Error("ReloadSuricataMainConf Error getting data from main.conf")
    }

    _, err = exec.Command(command, param, cmd+" "+anode["pid"]).Output()
    if err != nil {
        _ = ndb.InsertPluginCommand(uuid, "status", "Error")
        _ = ndb.InsertPluginCommand(uuid, "output", "ReloadSuricataMainConf reload Suricata mainconf error: "+err.Error())
        _ = ndb.InsertPluginCommand(uuid, "command", command+param+cmd+" "+anode["pid"])
        logs.Error("ReloadSuricataMainConf/Error starting suricata from main conf: " + err.Error())
        return err
    }
    _ = ndb.InsertPluginCommand(uuid, "status", "Success")
    _ = ndb.InsertPluginCommand(uuid, "output", "ReloadSuricataMainConf reloaded successfully")
    return nil
}

func GetMD5files(files map[string]map[string]string) (data map[string]map[string]string, err error) {
    var MD5data = map[string]map[string]string{}

    for x := range files {
        if MD5data[x] == nil {
            MD5data[x] = map[string]string{}
        }

        md5, err := utils.CalculateMD5(files[x]["nodepath"] + files[x]["path"])
        MD5data[x]["path"] = files[x]["path"]
        if err != nil {
            MD5data[x]["md5"] = ""
        } else {
            MD5data[x]["md5"] = md5
        }
    }

    return MD5data, err
}

// func LaunchSuricataServiceOld(uuid string, iface string) (err error) {
//     //historical log
//     uuidLog := utils.Generate()
//     currentTime := time.Now()
//     timeFormated := currentTime.Format("2006-01-02T15:04:05")
//     _ = ndb.InsertPluginCommand(uuidLog, "date", timeFormated)
//     _ = ndb.InsertPluginCommand(uuidLog, "id", uuid)
//     _ = ndb.InsertPluginCommand(uuidLog, "type", "Suricata")
//     _ = ndb.InsertPluginCommand(uuidLog, "action", "LaunchSuricataServiceOld")
//     _ = ndb.InsertPluginCommand(uuidLog, "description", "Launch suricata service")

//     fullpidfile, err := utils.GetKeyValueString("suricata", "fullpidfile")
//     if err != nil {
//         logs.Error("LaunchSuricataService Error getting data from main.conf: " + err.Error())
//     }
//     suricata, err := utils.GetKeyValueString("suricata", "suricata")
//     if err != nil {
//         logs.Error("LaunchSuricataService Error getting data from main.conf: " + err.Error())
//     }
//     suricata_config, err := utils.GetKeyValueString("suricata", "suricata_config")
//     if err != nil {
//         logs.Error("LaunchSuricataService Error getting data from main.conf: " + err.Error())
//     }

//     mainConfData, err := ndb.GetMainconfData()
//     if mainConfData["suricata"]["status"] == "disabled" {
//         return nil
//     }

//     allPlugins, err := ndb.GetPlugins()

//     if allPlugins[uuid]["configFile"] != "" {
//         suricata_config = allPlugins[uuid]["configFile"]
//     } else if suricata_config == "" {
//         str := fmt.Sprintf("SURICATA - Start Suricata - missing suricata configuration file, please review default value in main.conf, or configFile property of Suricata %s ", allPlugins[uuid]["name"])
//         logs.Error(str)
//         return errors.New(str)
//     }

//     // bpfFilter := ""
//     suricata_iface := ""
//     if allPlugins[uuid]["interface"] != "" {
//         suricata_iface = allPlugins[uuid]["interface"]
//     } else if iface != "" {
//         suricata_iface = iface
//     } else {
//         str := "SURICATA - Start Suricata - no interface defined - aborting"
//         logs.Error(str)
//         return errors.New(str)
//     }

//     suricata_pidfile := ""
//     if fullpidfile != "" {
//         suricata_pidfile = strings.Replace(fullpidfile, "<ID>", uuid, -1)
//     } else {
//         suricata_pidfile = strings.Replace("/var/run/suricata/<ID>-pidfile.pid", "<ID>", uuid, -1)
//     }

//     args := []string{}
//     args = append(args, "-D")
//     args = append(args, "-i")
//     args = append(args, suricata_iface)
//     args = append(args, "-c")
//     args = append(args, suricata_config)
//     args = append(args, "--pidfile")
//     args = append(args, suricata_pidfile)

//     if allPlugins[uuid]["bpfFile"] != "" {
//         args = append(args, "-F")
//         args = append(args, allPlugins[uuid]["bpfFile"])
//     } else if allPlugins[uuid]["bpf"] != "" {
//         args = append(args, allPlugins[uuid]["bpf"])
//     }

//     err = os.Remove(suricata_pidfile)
//     if err != nil {
//         logs.Error("SURICATA - Cannot remove pid file %s -> %s", suricata_pidfile, err.Error())
//     }

//     cmd := exec.Command(suricata, args...)

//     stdoutStderr, err := cmd.CombinedOutput()
//     if err != nil {
//         logs.Error(err)
//     }
//     logs.Debug("out -> %v", string(stdoutStderr))
//     // err = cmd.Run()
//     if err != nil {
//         //error launching suricata
//         // logs.Error(stdoutStderr.String())
//         _ = ndb.InsertPluginCommand(uuidLog, "status", "Error")
//         _ = ndb.InsertPluginCommand(uuidLog, "output", "LaunchSuricataService Launch Suricata mainconf error")
//         vals := strings.Join(args, ",")
//         _ = ndb.InsertPluginCommand(uuidLog, "command", suricata+vals)

//         logs.Error("plugin/LaunchSuricataService error launching Suricata: " + err.Error())
//         return errors.New("Error Launching suricata - " + err.Error())
//     } else {
//         time.Sleep(time.Second * 1)
//         //read file
//         currentpid, err := os.Open(suricata_pidfile)
//         if err != nil {
//             logs.Error("plugin/LaunchSuricataService error openning Suricata: " + err.Error())
//             return err
//         }
//         defer currentpid.Close()
//         pid, err := ioutil.ReadAll(currentpid)
//         PidNumber := strings.Split(string(pid), "\n")

//         //save pid to db
//         err = ndb.UpdatePluginValue(uuid, "pid", PidNumber[0])
//         if err != nil {
//             logs.Error("plugin/LaunchSuricataService error updating pid at DB: " + err.Error())
//             return err
//         }

//         //change DB status
//         err = ndb.UpdatePluginValue(uuid, "previousStatus", "enabled")
//         if err != nil {
//             logs.Error("plugin/LaunchSuricataService error: " + err.Error())
//             return err
//         }

//         //change DB status
//         err = ndb.UpdatePluginValue(uuid, "status", "enabled")
//         if err != nil {
//             logs.Error("plugin/LaunchSuricataService error: " + err.Error())
//             return err
//         }
//     }
//     _ = ndb.InsertPluginCommand(uuidLog, "status", "Success")
//     _ = ndb.InsertPluginCommand(uuidLog, "output", "LaunchSuricataService reloaded successfully")

//     return nil
// }

func LaunchSuricataService(uuid string, iface string) (err error) {
    //historical log
    uuidLog := utils.Generate()
    currentTime := time.Now()
    timeFormated := currentTime.Format("2006-01-02T15:04:05")
    _ = ndb.InsertPluginCommand(uuidLog, "date", timeFormated)
    _ = ndb.InsertPluginCommand(uuidLog, "id", uuid)
    _ = ndb.InsertPluginCommand(uuidLog, "type", "Suricata")
    _ = ndb.InsertPluginCommand(uuidLog, "action", "LaunchSuricataService")
    _ = ndb.InsertPluginCommand(uuidLog, "description", "Launch suricata service")

    logs.Info("SURICATA -> Launch suricata %s, if needed", uuid)

    fullpidfile, err := utils.GetKeyValueString("suricata", "fullpidfile")
    if err != nil {
        logs.Error("LaunchSuricataService Error getting data from main.conf -> suricata / fullpidfile: " + err.Error())
    }
    rulesetPath, err := utils.GetKeyValueString("suricataRuleset", "path")
    if err != nil {
        logs.Error("LaunchSuricataService Error getting data from main.conf: " + err.Error())
    }
    suricata, err := utils.GetKeyValueString("suricata", "suricata")
    if err != nil {
        logs.Error("LaunchSuricataService Error getting data from main.conf: " + err.Error())
    }
    suricata_config, err := utils.GetKeyValueString("suricata", "suricata_config")
    if err != nil {
        logs.Error("LaunchSuricataService Error getting data from main.conf: " + err.Error())
    }

    mainConfData, err := ndb.GetMainconfData()
    if mainConfData["suricata"]["status"] == "disabled" {
        logs.Info("Suricata is disabled by main.conf - won't run any suricata service")
        return nil
    }

    allPlugins, err := ndb.GetPlugins()

    args := []string{}
    args = append(args, "-D")

    if allPlugins[uuid]["configFile"] != "" {
        suricata_config = allPlugins[uuid]["configFile"]
        args = append(args, "-c")
        args = append(args, suricata_config)
    } else if suricata_config == "" {
        _ = ndb.InsertPluginCommand(uuidLog, "status", "Error")
        _ = ndb.InsertPluginCommand(uuidLog, "output", "SURICATA - Start Suricata - missing suricata configuration file, please review default value in main.conf, or configFile property of Suricata "+allPlugins[uuid]["name"])
        str := fmt.Sprintf("SURICATA - Start Suricata - missing suricata configuration file, please review default value in main.conf, or configFile property of Suricata %s ", allPlugins[uuid]["name"])
        logs.Error(str)
        return errors.New(str)
    }

    suricata_iface := ""
    if allPlugins[uuid]["interface"] != "" {
        suricata_iface = allPlugins[uuid]["interface"]
    } else if iface != "" {
        suricata_iface = iface
    }

    if suricata_iface != "" {
        args = append(args, "-i")
        args = append(args, suricata_iface)
    }

    suricata_pidfile := ""
    if fullpidfile != "" {
        suricata_pidfile = strings.Replace(fullpidfile, "<ID>", uuid, -1)
        args = append(args, "--pidfile")
        args = append(args, suricata_pidfile)
    } else {
        suricata_pidfile = strings.Replace("/var/run/suricata/<ID>-pidfile.pid", "<ID>", uuid, -1)
    }

    suricata_ruleset_name := ""
    if allPlugins[uuid]["localRulesetName"] != "" {
        suricata_ruleset_name = rulesetPath + allPlugins[uuid]["localRulesetName"] + ".rules"
        args = append(args, "-S")
        args = append(args, suricata_ruleset_name)
    }

    if allPlugins[uuid]["bpfFile"] != "" {
        args = append(args, "-F")
        args = append(args, allPlugins[uuid]["bpfFile"])
    } else if allPlugins[uuid]["bpf"] != "" {
        args = append(args, allPlugins[uuid]["bpf"])
    }

    logs.Info("SURICATA -> running configuration Test")
    result, err := SuricataConfigurationTest(uuid)
    if err != nil {
        _ = ndb.InsertPluginCommand(uuidLog, "status", "Error")
        _ = ndb.InsertPluginCommand(uuidLog, "output", "Suricata configuration check Error, Can't start Suricata")
        logs.Error("Suricata configuration check Error, Can't start Suricata")
        logs.Error(result)
        return errors.New(result["error"])
    }

    err = os.Remove(suricata_pidfile)
    if err != nil {
        _ = ndb.InsertPluginCommand(uuidLog, "status", "Error")
        _ = ndb.InsertPluginCommand(uuidLog, "output", "Suricata configuration Error removing PIDfile")
    }

    vals := strings.Join(args, ",")
    _ = ndb.InsertPluginCommand(uuidLog, "command", suricata+vals)

    logs.Info("SURICATA -> Running Suricata with args %+v", args)
    cmd := exec.Command(suricata, args...)

    stdoutStderr, err := cmd.CombinedOutput()
    logs.Debug("out -> %v", string(stdoutStderr))
    if err != nil {
        _ = ndb.InsertPluginCommand(uuidLog, "status", "Error")
        _ = ndb.InsertPluginCommand(uuidLog, "output", "plugin/LaunchSuricataService error launching Suricata: "+err.Error())
        logs.Error("plugin/LaunchSuricataService error launching Suricata: " + err.Error())
        return errors.New("Error Launching suricata - " + err.Error())
    } else {
        time.Sleep(time.Second * 10)
        logs.Info("SURICATA -> Started ")

        //change DB status
        logs.Info("SURICATA -> UPDATE previous status")
        err = ndb.UpdatePluginValue(uuid, "previousStatus", "enabled")
        if err != nil {
            logs.Error("plugin/LaunchSuricataService set previous status error : " + err.Error())
            return err
        }

        //change DB status
        logs.Info("SURICATA -> UPDATE current status")
        err = ndb.UpdatePluginValue(uuid, "status", "enabled")
        if err != nil {
            logs.Error("plugin/LaunchSuricataService set current status error: " + err.Error())
            return err
        }

        //read file
        currentpid, err := os.Open(suricata_pidfile)
        if err != nil {
            logs.Error("plugin/LaunchSuricataService error openning Suricata: " + err.Error())
            return err
        }
        defer currentpid.Close()
        pid, err := ioutil.ReadAll(currentpid)
        PidNumber := strings.Split(string(pid), "\n")

        logs.Info("SURICATA -> Started, new PID %v", PidNumber)

        //save pid to db
        err = ndb.UpdatePluginValue(uuid, "pid", PidNumber[0])
        if err != nil {
            logs.Error("plugin/LaunchSuricataService error updating pid at DB: " + err.Error())
            return err
        }

    }
    logs.Info("SURICATA -> Launch Service done.")

    _ = ndb.InsertPluginCommand(uuidLog, "status", "Success")
    _ = ndb.InsertPluginCommand(uuidLog, "output", "LaunchSuricataService reloaded successfully")
    return nil
}

func StopSuricataService(uuid string, status string) (err error) {
    //historical log
    uuidLog := utils.Generate()
    currentTime := time.Now()
    timeFormated := currentTime.Format("2006-01-02T15:04:05")
    _ = ndb.InsertPluginCommand(uuidLog, "date", timeFormated)
    _ = ndb.InsertPluginCommand(uuidLog, "id", uuid)
    _ = ndb.InsertPluginCommand(uuidLog, "type", "Suricata")
    _ = ndb.InsertPluginCommand(uuidLog, "action", "StopSuricataService")
    _ = ndb.InsertPluginCommand(uuidLog, "description", "Stop suricata service")

    suricataBackup, err := utils.GetKeyValueString("suricata", "backup")
    if err != nil {
        logs.Error("StopSuricataService Error getting data from main.conf: " + err.Error())
    }
    suricataPidfile, err := utils.GetKeyValueString("suricata", "pidfile")
    if err != nil {
        logs.Error("StopSuricataService Error getting data from main.conf: " + err.Error())
    }
    //pid
    allPlugins, err := ndb.GetPlugins()

    //kill suricata process
    PidInt, _ := strconv.Atoi(strings.Trim(string(allPlugins[uuid]["pid"]), "\n"))
    process, _ := os.FindProcess(PidInt)
    err = process.Kill()
    if err != nil {
        _ = ndb.InsertPluginCommand(uuidLog, "status", "Error")
        _ = ndb.InsertPluginCommand(uuidLog, "output", "StopSuricataServicekill process error: "+err.Error())
        _ = ndb.InsertPluginCommand(uuidLog, "command", strconv.Itoa(PidInt))
    }

    //delete pid file
    err = os.Remove(suricataBackup + uuid + "-" + suricataPidfile)
    if err != nil {
        _ = ndb.InsertPluginCommand(uuidLog, "status", "Error")
        _ = ndb.InsertPluginCommand(uuidLog, "output", "StopSuricataService remove PID file error: "+err.Error())
        _ = ndb.InsertPluginCommand(uuidLog, "command", suricataBackup+uuid+"-"+suricataPidfile)
    }

    //change DB pid
    // err = ndb.UpdatePluginValue(uuid, "pid", "none")
    // if err != nil {
    //     logs.Error("plugin/StopSuricataService error updating pid at DB: " + err.Error())
    //     return err
    // }

    //change DB status
    err = ndb.UpdatePluginValue(uuid, "previousStatus", status)
    if err != nil {
        logs.Error("plugin/StopSuricataService error: " + err.Error())
        return err
    }

    //change DB status
    err = ndb.UpdatePluginValue(uuid, "status", "disabled")
    if err != nil {
        logs.Error("plugin/StopSuricataService error: " + err.Error())
        return err
    }

    _ = ndb.InsertPluginCommand(uuidLog, "status", "Success")
    _ = ndb.InsertPluginCommand(uuidLog, "output", "StopSuricataService stopped successfully")
    return nil
}

func SuricataConfigurationTest(uuid string) (responseBack map[string]string, err error) {

    var response = map[string]string{}

    suricata_config, err := utils.GetKeyValueString("suricata", "suricata_config")
    if err != nil {
        logs.Error("SuricataConfigurationTest Error getting data from main.conf: " + err.Error())
    }
    suricata, err := utils.GetKeyValueString("suricata", "suricata")
    if err != nil {
        logs.Error("SuricataConfigurationTest Error getting data from main.conf: " + err.Error())
    }

    suricata_interface := ""
    if uuid != "" {
        allPlugins, err := ndb.GetPlugins()
        if err != nil {
            logs.Error("SuricataConfigurationTest Error getting plugin list from DB: " + err.Error())
        } else if allPlugins[uuid]["configFile"] != "" {
            suricata_config = allPlugins[uuid]["configFile"]
            suricata_interface = allPlugins[uuid]["interface"]
        }
    }

    if suricata_config == "" {
        str := fmt.Sprintf("SURICATA - check Suricata configuration - missing suricata configuration file, please review default value in main.conf, or configFile property of Suricata")
        logs.Error(str)
        return response, errors.New(str)
    }

    rulesetPath, err := utils.GetKeyValueString("suricataRuleset", "path")
    if err != nil {
        logs.Error("SURICATA - Check Suricata configuration - Error getting rulesetpath from main.conf: " + err.Error())
    }

    args := []string{}
    args = append(args, "-c")
    args = append(args, suricata_config)

    args = append(args, "-i")
    args = append(args, suricata_interface)

    if allPlugins[uuid]["localRulesetName"] != "" {
        suricata_ruleset_name = rulesetPath + allPlugins[uuid]["localRulesetName"] + ".rules"
        args = append(args, "-S")
        args = append(args, suricata_ruleset_name)
    }

    args = append(args, "-T")

    cmd := exec.Command(suricata, args...)

    response["ack"] = "true"

    stdoutStderr, err := cmd.CombinedOutput()
    if err != nil {
        logs.Error(err)
        response["ack"] = "false"
        response["error"] = err.Error()
    }
    response["result"] = string(stdoutStderr)
    logs.Debug("out -> %v", string(stdoutStderr))
    return response, err
}

func SuricataDumpCurrentConfig() (responseBack map[string]string, err error) {
    var response = map[string]string{}

    suricata, err := utils.GetKeyValueString("suricata", "suricata")
    if err != nil {
        logs.Error("SuricataConfigurationTest Error getting data from main.conf: " + err.Error())
    }

    args := []string{}
    args = append(args, "--dump")

    cmd := exec.Command(suricata, args...)

    response["ack"] = "true"

    stdoutStderr, err := cmd.CombinedOutput()
    if err != nil {
        logs.Error(err)
        response["ack"] = "false"
        response["error"] = err.Error()
    }
    response["result"] = string(stdoutStderr)
    logs.Debug("out -> %v", string(stdoutStderr))
    return response, err

}

func SuricataVersion() (version map[string]string, err error) {
    var response = map[string]string{}

    suricata, err := utils.GetKeyValueString("suricata", "suricata")
    if err != nil {
        logs.Error("SuricataVersion Error getting suricata binary data from main.conf: " + err.Error())
    }

    args := []string{}
    args = append(args, "-V")

    cmd := exec.Command(suricata, args...)

    response["ack"] = "true"

    stdoutStderr, err := cmd.CombinedOutput()
    if err != nil {
        logs.Error("Suricata Version Error -> %s", err.Error())
        response["ack"] = "false"
        response["error"] = err.Error()
    }
    response["version_string"] = string(stdoutStderr)
    re := regexp.MustCompile("(\\d+\\.\\d+\\.\\d+)")
    match := re.FindStringSubmatch(string(stdoutStderr))
    if len(match) > 1 {
        response["version"] = match[1]
    }

    logs.Debug("Suricata Version output -> %v", string(stdoutStderr))

    return response, err
}

func GetSuricataRulesets() (version map[string]map[string]string, err error) {
    plugins, err := ndb.GetPlugins()
    if err != nil {
        logs.Error("suricata/GetSuricataRulesets error: " + err.Error())
        return nil, err
    }

    var serviceValues = map[string]map[string]string{}

    for x := range plugins {
        if plugins[x]["type"] == "suricata" {
            if serviceValues[x] == nil {
                serviceValues[x] = map[string]string{}
            }
            serviceValues[x]["name"] = plugins[x]["name"]
            serviceValues[x]["ruleset"] = plugins[x]["rulesetName"]
        }
    }

    return serviceValues, err
}
