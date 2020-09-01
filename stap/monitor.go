package stap

import (
    "github.com/astaxie/beego/logs"
    "golang.org/x/crypto/ssh"
    "owlhnode/database"
    "regexp"
    "strconv"
    "strings"
)

//check whether owlh is alive
func CheckOwlhAlive(uuid string) (alive bool, sshSession *ssh.Session) {
    alive, sshSession = owlh_connect(uuid)
    logs.Info("Stap Server Task with uuid: " + uuid)
    if alive {
        logs.Info("ALIVE Stap Server " + uuid)
        return true, sshSession
    }
    logs.Error("NOT ALIVE Stap Server " + uuid)
    return false, nil
}

//check sniffer status through CPU, MEM and STORAGE status
func GetStatusSniffer(uuid string) (running bool, status bool) {
    owlh, err := ndb.GetStapServerInformation(uuid)
    if err != nil {
        logs.Error("Error retrieving stap server information")
    }
    logs.Info("Checking Sniffer status for uuid: " + uuid)

    running, pid, cpu, mem := GetStatusSnifferSSH(uuid)
    cpuStatus := GetStatusCPU(owlh, cpu, uuid)
    memStatus := GetStatusMEM(owlh, mem, uuid)
    storageStatus := GetStatusStorage(owlh, uuid)

    logs.Alert("Checking " + owlh["name"] + " - " + owlh["ip"] + " - PID:" + pid + " CPU: " + strconv.FormatBool(cpuStatus) + " MEM: " + strconv.FormatBool(memStatus) + " STORAGE: " + strconv.FormatBool(storageStatus))
    if cpuStatus && memStatus && storageStatus {
        return running, true
    }
    return running, false
}

//check CPU status
func GetStatusCPU(owlh map[string]string, cpu string, uuid string) (status bool) {
    var validCPU = regexp.MustCompile(`(\d+)` + cpu + `;`)
    if validCPU == nil {
        return false
    }
    localCPU, _ := strconv.ParseFloat(cpu, 64)
    ddbbCPU, _ := strconv.ParseFloat(owlh["max_cpu"], 64)
    logs.Info("Check CPU for " + owlh["name"] + " - " + owlh["ip"])
    if localCPU > ddbbCPU {
        logs.Error("SNIFFER -> Too much CPU on " + owlh["name"] + " - " + owlh["ip"])
        StopSniffer(uuid)
        return false
    }
    return true
}

//check MEM status
func GetStatusMEM(owlh map[string]string, mem string, uuid string) (status bool) {
    var validMEM = regexp.MustCompile(`(\d+)` + mem + `;`)
    if validMEM == nil {
        return false
    }
    localMEM, _ := strconv.ParseFloat(mem, 64)
    ddbbMEM, _ := strconv.ParseFloat(owlh["max_mem"], 64)
    logs.Info("Check MEM for " + owlh["name"] + " - " + owlh["ip"])
    if localMEM > ddbbMEM {
        logs.Error("SNIFFER -> Too much MEM on " + owlh["name"] + " - " + owlh["ip"])
        StopSniffer(uuid)
        return false
    }
    return true
}

//check STORAGE status
func GetStatusStorage(owlh map[string]string, uuid string) (status bool) {
    var pcapPath = owlh["pcap_path"]
    status, path, storage := GetStatusStorageSSh(uuid, pcapPath)
    if status {
        localFunc := strings.Replace(storage, "%", "", -1)
        localStorage, _ := strconv.ParseFloat(localFunc, 64)
        ddbbStorage, _ := strconv.ParseFloat(owlh["max_storage"], 64)
        if localStorage > ddbbStorage {
            logs.Error("SNIFFER -> Too much STORAGE on " + owlh["name"] + " - " + owlh["ip"] + " on path: " + path)
            StopSniffer(uuid)
            return false
        }
        logs.Error("SNIFFER -> STORAGE status is correct!!")
        return true
    }
    return false
}

//launch the sniffer
func RunSniffer(uuid string) {
    RunSnifferSSH(uuid)
}

//stop the sniffer
func StopSniffer(uuid string) {
    StopSnifferSSH(uuid)
}

//get pcap file list forfrom remote node
func GetFileList(uuid string) {
    owlh, err := ndb.GetStapServerInformation(uuid)
    if err != nil {
        logs.Error("Error retrieving stap server information")
    }
    logs.Error("Get file list for " + owlh["name"] + " - " + owlh["ip"])
    file_list := GetFileListSSH(uuid, owlh, owlh["pcap_path"])
    for file := range file_list {
        var validOutput = regexp.MustCompile(`\.pcap+`)
        if validOutput.MatchString(file_list[file]) {
            logs.Notice("Change remote file owned")
            OwnerOwlh(uuid, owlh, file_list[file])

            logs.Notice("Copy files using sftp command and remove it!!")
            TransportFile(uuid, owlh, file_list[file])

            logs.Warn("File list completed!")
        }
    }

}

//change pcap remote file owner
func OwnerOwlh(uuid string, owlh map[string]string, fileRemote string) {
    logs.Info("Set " + owlh["name"] + " - " + owlh["ip"] + " as owner of file " + owlh["owlh_user"])
    OwnerOwlhSSH(uuid, owlh, fileRemote)
}

//send remote file to local machine
func TransportFile(uuid string, owlh map[string]string, file string) {
    logs.Info("Get file " + owlh["local_pcap_path"] + " from " + owlh["name"] + " - " + owlh["ip"])
    TransportFileSSH(uuid, owlh, file)
}
