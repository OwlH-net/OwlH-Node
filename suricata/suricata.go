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
            logs.Error("Error writting data into BPF file: " + err.Error())
            return err
        }
    } else {
        //make backup file
        err = utils.BackupFile(path, n["service"]+"-"+file)
        if err != nil {
            logs.Error("Error creating BPF backup: " + err.Error())
            return err
        }

        //write bpf into the file
        err = utils.UpdateBPFFile(path, n["service"]+"-"+file, n["value"])
        if err != nil {
            logs.Error("Error UpdateBPFFile: " + err.Error())
            return err
        }
    }

    return nil
}

//Retrieve data, make a backup file and write the new data on the original file
func SyncRulesetFromMaster(file map[string][]byte) (err error) {
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

    // //get name from db
    // pluginName, err := ndb.GetPluginsByParam(string(file["service"]), "name")
    // if err != nil {
    //     logs.Error("SyncRulesetFromMaster Error getting data from database: " + err.Error())
    //     return err
    // }

    //replace file by name
    fileName := strings.Replace(string(file["rulesetName"]), " ", "-", -1)
    plug := strings.Replace(fileToEdit, "<NAME>", fileName, -1)

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
    if SuriRunning() {
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
        socket, err := utils.GetKeyValueString("SuricataRulesetReload", "socket")
        if err != nil {
            logs.Error("SyncRulesetFromMaster Error getting data from main.conf: " + err.Error())
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
        logs.Error("Error stopping suricata: " + err.Error())
        return "", err
    }
    return "Suricata stopped ", nil
}

func GetSuricataServices() (data map[string]map[string]string, err error) {
    data, err = ndb.GetServices("suricata")
    if err != nil {
        logs.Error("GetSuricataServices Error: " + err.Error())
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

        //remove tar.gz file
        // os.Remove(x+"/file.tar.gzip")
    }
    return nil
}

func StartSuricataMainConf(anode map[string]string) (err error) {
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
        logs.Error("StartSuricataMainConf/Error starting suricata: " + err.Error())
        return err
    }
    return nil
}
func StopSuricataMainConf(anode map[string]string) (err error) {
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
    // err = utils.RunCommand(cmd, param)
    if err != nil {
        logs.Error("StopSuricataMainConf/Error stopping suricata: " + err.Error())
        return err
    }
    return nil
}
func KillSuricataMainConf(anode map[string]string) (err error) {
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
        logs.Error("KillSuricataMainConf/Error starting suricata from main conf: " + err.Error())
        return err
    }
    return nil
}
func ReloadSuricataMainConf(anode map[string]string) (err error) {
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
        logs.Error("ReloadSuricataMainConf/Error starting suricata from main conf: " + err.Error())
        return err
    }
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

func LaunchSuricataService(uuid string, iface string) (err error) {
    fullpidfile, err := utils.GetKeyValueString("suricata", "fullpidfile")
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
        return nil
    }

    allPlugins, err := ndb.GetPlugins()

    if allPlugins[uuid]["configFile"] != "" {
        suricata_config = allPlugins[uuid]["configFile"]
    } else if suricata_config == "" {
        str := fmt.Sprintf("SURICATA - Start Suricata - missing suricata configuration file, please review default value in main.conf, or configFile property of Suricata %s ", allPlugins[uuid]["name"])
        logs.Error(str)
        return errors.New(str)
    }

    // bpfFilter := ""
    suricata_iface := ""
    if allPlugins[uuid]["interface"] != "" {
        suricata_iface = allPlugins[uuid]["interface"]
    } else if iface != "" {
        suricata_iface = iface
    } else {
        str := "SURICATA - Start Suricata - no interface defined - aborting"
        logs.Error(str)
        return errors.New(str)
    }

    suricata_pidfile := ""
    if fullpidfile != "" {
        suricata_pidfile = strings.Replace(fullpidfile, "<ID>", uuid, -1)
    } else {
        suricata_pidfile = strings.Replace("/var/run/suricata/<ID>-pidfile.pid", "<ID>", uuid, -1)
    }

    args := []string{}
    args = append(args, "-D")
    args = append(args, "-i")
    args = append(args, suricata_iface)
    args = append(args, "-c")
    args = append(args, suricata_config)
    args = append(args, "--pidfile")
    args = append(args, suricata_pidfile)

    if allPlugins[uuid]["bpfFile"] != "" {
        args = append(args, "-F")
        args = append(args, allPlugins[uuid]["bpfFile"])
    } else if allPlugins[uuid]["bpf"] != "" {
        args = append(args, allPlugins[uuid]["bpf"])
    }

    err = os.Remove(suricata_pidfile)
    if err != nil {
        logs.Error("SURICATA - Cannot remove pid file %s -> %s", suricata_pidfile, err.Error())
    }

    cmd := exec.Command(suricata, args...)

    stdoutStderr, err := cmd.CombinedOutput()
    if err != nil {
        logs.Error(err)
    }
    logs.Debug("out -> %v", string(stdoutStderr))
    // err = cmd.Run()
    if err != nil {
        //error launching suricata
        // logs.Error(stdoutStderr.String())
        logs.Error("plugin/LaunchSuricataService error launching Suricata: " + err.Error())
        return errors.New("Error Launching suricata - " + err.Error())
    } else {
        time.Sleep(time.Second * 1)
        //read file
        currentpid, err := os.Open(suricata_pidfile)
        if err != nil {
            logs.Error("plugin/LaunchSuricataService error openning Suricata: " + err.Error())
            return err
        }
        defer currentpid.Close()
        pid, err := ioutil.ReadAll(currentpid)
        PidNumber := strings.Split(string(pid), "\n")

        //save pid to db
        err = ndb.UpdatePluginValue(uuid, "pid", PidNumber[0])
        if err != nil {
            logs.Error("plugin/LaunchSuricataService error updating pid at DB: " + err.Error())
            return err
        }

        //change DB status
        err = ndb.UpdatePluginValue(uuid, "previousStatus", "enabled")
        if err != nil {
            logs.Error("plugin/LaunchSuricataService error: " + err.Error())
            return err
        }

        //change DB status
        err = ndb.UpdatePluginValue(uuid, "status", "enabled")
        if err != nil {
            logs.Error("plugin/LaunchSuricataService error: " + err.Error())
            return err
        }
    }
    return nil
}

func LaunchSuricataServiceWithData(uuid string, iface string) (err error) {
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
        args = append(args, "-s")
        args = append(args, suricata_ruleset_name)
    }

    if allPlugins[uuid]["bpfFile"] != "" {
        args = append(args, "-F")
        args = append(args, allPlugins[uuid]["bpfFile"])
    } else if allPlugins[uuid]["bpf"] != "" {
        args = append(args, allPlugins[uuid]["bpf"])
    }

    err = os.Remove(suricata_pidfile)

    cmd := exec.Command(suricata, args...)

    stdoutStderr, err := cmd.CombinedOutput()
    if err != nil {
        logs.Error(err)
    }
    logs.Debug("out -> %v", string(stdoutStderr))
    if err != nil {
        logs.Error("plugin/LaunchSuricataService error launching Suricata: " + err.Error())
        return errors.New("Error Launching suricata - " + err.Error())
    } else {
        time.Sleep(time.Second * 1)
        //read file
        currentpid, err := os.Open(suricata_pidfile)
        if err != nil {
            logs.Error("plugin/LaunchSuricataService error openning Suricata: " + err.Error())
            return err
        }
        defer currentpid.Close()
        pid, err := ioutil.ReadAll(currentpid)
        PidNumber := strings.Split(string(pid), "\n")

        //save pid to db
        err = ndb.UpdatePluginValue(uuid, "pid", PidNumber[0])
        if err != nil {
            logs.Error("plugin/LaunchSuricataService error updating pid at DB: " + err.Error())
            return err
        }

        //change DB status
        err = ndb.UpdatePluginValue(uuid, "previousStatus", "enabled")
        if err != nil {
            logs.Error("plugin/LaunchSuricataService error: " + err.Error())
            return err
        }

        //change DB status
        err = ndb.UpdatePluginValue(uuid, "status", "enabled")
        if err != nil {
            logs.Error("plugin/LaunchSuricataService error: " + err.Error())
            return err
        }
    }
    return nil
}

func StopSuricataService(uuid string, status string) (err error) {
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
    _ = process.Kill()

    //delete pid file
    _ = os.Remove(suricataBackup + uuid + "-" + suricataPidfile)

    //change DB pid
    err = ndb.UpdatePluginValue(uuid, "pid", "none")
    if err != nil {
        logs.Error("plugin/StopSuricataService error updating pid at DB: " + err.Error())
        return err
    }

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

    return nil
}
