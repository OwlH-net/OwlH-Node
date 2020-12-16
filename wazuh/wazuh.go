package wazuh

import (
    "bufio"
    "errors"
    "github.com/astaxie/beego/logs"
    "os"
    "os/exec"
    "regexp"
    "strings"
    // "io/ioutil"
    "encoding/json"
    "strconv"
    // "bytes"
    "owlhnode/database"
    "owlhnode/utils"
    "path/filepath"
    "time"
)

func WazuhPath() (exists bool) {
    var err error
    path, err := utils.GetKeyValueString("loadDataWazuhPath", "path")
    if err != nil {
        logs.Error("WazuhPath Error getting data from main.conf")
    }

    if _, err := os.Stat(path); os.IsNotExist(err) {
        logs.Error("Wazuh is not installed, at least, /var/ossec folder does not exist")
        return false
    }
    return true
}

func WazuhBin() (exists bool) {
    var err error
    bin, err := utils.GetKeyValueString("loadDataWazuhBin", "bin")
    if err != nil {
        logs.Error("WazuhBin Error getting data from main.conf")
    }
    if _, err := os.Stat(bin); os.IsNotExist(err) {
        logs.Error("Wazuh bin does not exist")
        return false
    }
    return true
}

func WazuhRunning() (running bool) {
    var err error
    cmd, err := utils.GetKeyValueString("loadDataWazuhRunning", "cmd")
    if err != nil {
        logs.Error("WazuhRunning Error getting data from main.conf")
    }
    param, err := utils.GetKeyValueString("loadDataWazuhRunning", "param")
    if err != nil {
        logs.Error("WazuhRunning Error getting data from main.conf")
    }
    command, err := utils.GetKeyValueString("loadDataWazuhRunning", "command")
    if err != nil {
        logs.Error("WazuhRunning Error getting data from main.conf")
    }

    out, err := exec.Command(command, param, cmd).Output()
    if err == nil {
        if strings.Contains(string(out), "is running") {
            return true
        }
    }
    return false
}

func Installed() (isIt map[string]bool, err error) {
    wazuh := make(map[string]bool)
    wazuh["path"] = WazuhPath()
    wazuh["bin"] = WazuhBin()
    wazuh["running"] = WazuhRunning()

    if wazuh["path"] || wazuh["bin"] || wazuh["running"] {
        return wazuh, nil
    } else {
        return wazuh, errors.New("Wazuh isn't present or not running")
    }
}

//Run wazuh
func RunWazuh() (data string, err error) {
    cmd, err := utils.GetKeyValueString("wazuhStart", "start")
    if err != nil {
        logs.Error("RunWazuh Error getting data from main.conf")
        return "", err
    }
    param, err := utils.GetKeyValueString("wazuhStart", "param")
    if err != nil {
        logs.Error("RunWazuh Error getting data from main.conf")
        return "", err
    }
    command, err := utils.GetKeyValueString("wazuhStart", "command")
    if err != nil {
        logs.Error("RunWazuh Error getting data from main.conf")
        return "", err
    }

    _, err = exec.Command(command, param, cmd).Output()
    if err != nil {
        logs.Error("Error launching wazuh: " + err.Error())
        return "", err
    }
    return "wazuh system is on", nil
}

//Stop wazuh
func StopWazuh() (data string, err error) {
    cmd, err := utils.GetKeyValueString("wazuhStop", "stop")
    if err != nil {
        logs.Error("StopWazuh Error getting data from main.conf")
        return "", err
    }
    param, err := utils.GetKeyValueString("wazuhStop", "param")
    if err != nil {
        logs.Error("StopWazuh Error getting data from main.conf")
        return "", err
    }
    command, err := utils.GetKeyValueString("wazuhStop", "command")
    if err != nil {
        logs.Error("StopWazuh Error getting data from main.conf")
        return "", err
    }

    _, err = exec.Command(command, param, cmd).Output()
    if err != nil {
        logs.Error("Error stopping Wazuh: " + err.Error())
        return "", err
    }

    return "Wazuh stopped ", nil
}

func PingWazuhFiles() (files map[string]map[string]string, err error) {
    file, err := os.Open("/var/ossec/etc/ossec.conf")
    if err != nil {
        logs.Error(err)
        return nil, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    isInit := false
    isEnd := false
    filesPath := make(map[string]map[string]string)
    count := 0
    var size int64
    for scanner.Scan() {
        var init = regexp.MustCompile(`<!-- OWLH INIT -->`)
        var end = regexp.MustCompile(`<!-- OWLH END -->`)
        owlhInit := init.FindStringSubmatch(scanner.Text())
        owlhEnd := end.FindStringSubmatch(scanner.Text())
        if owlhInit != nil {
            isInit = true
            continue
        }
        if owlhEnd != nil {
            isEnd = true
        }
        if isInit && !isEnd {
            var locationPath = regexp.MustCompile(`<location>([^"]+)<\/location>`)
            locationFound := locationPath.FindStringSubmatch(scanner.Text())
            if locationFound != nil {
                fi, err := os.Stat(locationFound[1])
                if err != nil {
                    logs.Error("Error PingWazuhFiles checking file size: " + err.Error())
                    size = -1
                } else {
                    size = fi.Size()
                }

                if filesPath[strconv.Itoa(count)] == nil {
                    filesPath[strconv.Itoa(count)] = map[string]string{}
                }
                filesPath[strconv.Itoa(count)]["path"] = locationFound[1]
                filesPath[strconv.Itoa(count)]["size"] = strconv.FormatInt(size, 10)

                //Add new Incident
                //Add incident to database
                if size > 1024 {
                    uuid := utils.Generate()
                    nodeData, err := ndb.GetNodeData()
                    if err != nil {
                        logs.Error("Error creating incident for Wazuh: " + err.Error())
                        return nil, err
                    }
                    currentTime := time.Now()
                    timeFormated := currentTime.Format("2006-01-02T15:04:05")
                    err = ndb.PutIncidentNode(uuid, "date", timeFormated)
                    err = ndb.PutIncidentNode(uuid, "desc", "This is the description")
                    err = ndb.PutIncidentNode(uuid, "status", "new")                         // new, open, closed, delayed
                    err = ndb.PutIncidentNode(uuid, "level", "info")                         // warning, info or danger
                    err = ndb.PutIncidentNode(uuid, "filePath", locationFound[1])            // warning, info or danger
                    err = ndb.PutIncidentNode(uuid, "fileSize", strconv.FormatInt(size, 10)) // warning, info or danger
                    for x := range nodeData {
                        err = ndb.PutIncidentNode(uuid, "deviceName", nodeData[x]["name"])
                        err = ndb.PutIncidentNode(uuid, "deviceIp", nodeData[x]["ip"])
                        err = ndb.PutIncidentNode(uuid, "devicePort", nodeData[x]["port"])
                        err = ndb.PutIncidentNode(uuid, "deviceUuid", x)
                    }
                }
            }
        }
        count++
    }
    if err := scanner.Err(); err != nil {
        logs.Error(err)
    }

    return filesPath, err
}

type AllFiles struct {
    UUID  string   `json:"uuid"`
    Paths []string `json:"paths"`
}

func ModifyWazuhFile(anode map[string]interface{}) (err error) {
    ossec, err := utils.GetKeyValueString("loadDataWazuhPath", "ossec")
    if err != nil {
        logs.Error("ModifyWazuhFile Error getting data from main.conf")
        return err
    }
    ossecFile, err := utils.GetKeyValueString("loadDataWazuhPath", "file")
    if err != nil {
        logs.Error("ModifyWazuhFile Error getting data from main.conf")
        return err
    }
    ossecFileFullPath, err := utils.GetKeyValueString("loadDataWazuhPath", "fullPath")
    if err != nil {
        logs.Error("ModifyWazuhFile Error getting data from main.conf")
        return err
    }

    receivedWazuhFiles := AllFiles{}
    byteData, _ := json.Marshal(anode)
    json.Unmarshal(byteData, &receivedWazuhFiles)

    // //check gor empty array
    // if receivedWazuhFiles.Paths == nil || len(receivedWazuhFiles.Paths) <= 0 {
    //     return errors.New("ModifyWazuhFile error - New content is empty")
    // }
    //backup file
    err = utils.BackupFile(ossecFileFullPath, ossecFile)
    if err != nil {
        logs.Error("ModifyWazuhFile error: " + err.Error())
        return errors.New("ModifyWazuhFile error - Backup failed")
    }

    //modify ossec file
    file, err := os.Open(ossec)
    if err != nil {
        logs.Error("Error ModifyWazuhFile readding file: " + err.Error())
        return err
    }
    defer file.Close()

    // var buf bytes.Buffer
    isInit := false
    isEnd := false
    isSecondEnd := false
    scanner := bufio.NewScanner(file)
    var h int
    h = 0
    fileContent := make(map[int]string)

    for scanner.Scan() {
        var init = regexp.MustCompile(`<!-- OWLH INIT -->`)
        var end = regexp.MustCompile(`<!-- OWLH END -->`)
        owlhInit := init.FindStringSubmatch(scanner.Text())
        owlhEnd := end.FindStringSubmatch(scanner.Text())

        if owlhInit != nil {
            isInit = true
            fileContent[h] = "<!-- OWLH INIT -->"
            h++
            fileContent[h] = "<ossec_config>"
            h++
        }
        if owlhEnd != nil {
            isEnd = true
        }
        if isInit && !isEnd {

            for x := range receivedWazuhFiles.Paths {
                fileContent[h] = "\t<localfile>"
                h++
                fileContent[h] = "\t\t<log_format>syslog</log_format>"
                h++
                fileContent[h] = "\t\t<location>" + receivedWazuhFiles.Paths[x] + "</location>  "
                h++
                fileContent[h] = "\t</localfile>"
                h++
            }
            isEnd = true
            fileContent[h] = "</ossec_config>"
            h++
            fileContent[h] = "<!-- OWLH END -->"
        } else if isInit && isEnd && !isSecondEnd {
            var secondEnd = regexp.MustCompile(`<!-- OWLH END -->`)
            secondOwlhEnd := secondEnd.FindStringSubmatch(scanner.Text())
            if secondOwlhEnd == nil {
                continue
            } else if secondOwlhEnd != nil && !isSecondEnd {
                isSecondEnd = true
                // h++
                continue
            }

        } else {
            fileContent[h] = scanner.Text()
        }
        h++
    }

    if !isInit {
        fileContent[h] = "<!-- OWLH INIT -->"
        h++
        fileContent[h] = "<ossec_config>"
        h++
        for d := range receivedWazuhFiles.Paths {
            fileContent[h] = "\t<localfile>"
            h++
            fileContent[h] = "\t\t<log_format>syslog</log_format>"
            h++
            fileContent[h] = "\t\t<location>" + receivedWazuhFiles.Paths[d] + "</location>  "
            h++
            fileContent[h] = "\t</localfile>"
            h++
        }
        fileContent[h] = "</ossec_config>"
        h++
        fileContent[h] = "<!-- OWLH END -->"
        h++
    }

    if err := scanner.Err(); err != nil {
        logs.Error("ModifyWazuhFile. Scanner file error: " + err.Error())
        return err
    }

    saveIntoFile, err := os.OpenFile(ossec, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        logs.Error("Error ModifyWazuhFile readding file: " + err.Error())
        return err
    }

    //move backup for check error .bck

    defer saveIntoFile.Close()
    saveIntoFile.Truncate(0)
    saveIntoFile.Seek(0, 0)
    for x := 0; x < h; x++ {
        _, err = saveIntoFile.WriteAt([]byte(fileContent[x]+"\n"), 0) // Write at 0 beginning
        if err != nil {
            logs.Error("ModifyWazuhFile failed writing to file: %s", err)
            return err
        }
    }

    return err
}

func LoadFileLastLines(file map[string]string) (data map[string]string, err error) {
    values, err := utils.LoadFileLastLines(file)
    if err != nil {
        logs.Error(err.Error())
        return nil, err
    }

    return values, nil
}

// func LoadFileLastLines(file map[string]string)(data map[string]string, err error) {
//     command, err := utils.GetKeyValueString("execute", "command")
//     if err != nil {logs.Error("Error getting data from main.conf: "+err.Error())}
//     param, err := utils.GetKeyValueString("execute", "param")
//     if err != nil {logs.Error("Error getting data from main.conf: "+err.Error())}

//     linesResult := make(map[string]string)

//     if file["number"] != "none"{
//         lines,err := exec.Command(command, param, "tail -"+file["number"]+" "+file["path"]).Output()
//         if err != nil{logs.Error("LoadFileLastLines Error retrieving last lines of the path "+file["path"]+": "+err.Error()); return nil,err}

//         linesResult["result"] = string(lines)
//     }else{
//         fileReaded, err := ioutil.ReadFile(file["path"]) // just pass the file name
//         if err != nil {logs.Error("Error reading Wazuh file for path: "+file["path"]); return nil,err}

//         linesResult["result"] = string(fileReaded)
//     }

//     return linesResult, err
// }

func SaveFileContentWazuh(file map[string]string) (err error) {
    bytearray := []byte(file["content"])

    //check for empty file content
    if bytearray == nil || len(bytearray) <= 0 {
        logs.Error("SaveFileContentWazuh error - New content is empty")
        return errors.New("SaveFileContentWazuh error - New content is empty")
    }

    //split path
    fileName := filepath.Base(file["path"])
    fileDir := filepath.Dir(file["path"])

    //backup file
    err = utils.BackupFile(fileDir+"/", fileName)
    if err != nil {
        logs.Error("ModifyWazuhFile error: " + err.Error())
        return errors.New("ModifyWazuhFile error - Backup failed")
    }

    // err = utils.BackupFullPath(file["path"])
    // if err != nil {logs.Info("SaveFileContentWazuh. Error doing backup with function BackupFullPath: "+err.Error()); return err}

    //make byte array for save the file modified
    err = utils.WriteNewDataOnFile(file["path"], bytearray)
    if err != nil {
        logs.Info("SaveFileContentWazuh error doing backup with function WriteNewDataOnFile: " + err.Error())
        return err
    }

    return nil
}

func GetWazuhAgentID() (data utils.WazuhData, err error) {

    var wazuhData utils.WazuhData

    clientPath, err := utils.GetKeyValueString("loadDataWazuhPath", "clientPath")
    if err != nil {
        logs.Error("ModifyWazuhFile Error getting data from main.conf")
        return wazuhData, err
    }

    if _, err := os.Stat(clientPath); os.IsNotExist(err) {
        logs.Error("Wazuh client file does not exist")
        return wazuhData, err
    }

    file, err := os.Open(clientPath)
    if err != nil {
        logs.Error(err)
        return wazuhData, err
    }
    defer file.Close()

    scanner := bufio.NewScanner(file)
    var line []string
    for scanner.Scan() {
        line = strings.Split(scanner.Text(), " ")
    }
    wazuhData.ID = line[0]
    wazuhData.Ip = line[2]
    wazuhData.Name = line[1]
    file.Close()

    logs.Info("wazuh agent data --- %+v", wazuhData)
    return wazuhData, nil
}
