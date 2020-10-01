package stap

import (
    "github.com/astaxie/beego/logs"
    "os/exec"
    "owlhnode/utils"  
    "strconv"
    "io/ioutil"
    "strings"
    "time"
)

func Pcap_replay()() {
    var err error
    //load in_, out_queue and interface from main.conf
    if err != nil {logs.Error("Pcap_replay Error getting data from main.conf")}
    inQueue, err := utils.GetKeyValueString("stap", "in_queue")
    if err != nil {logs.Error("Pcap_replay Error getting data from main.conf")}
    outQueue, err := utils.GetKeyValueString("stap", "out_queue")
    if err != nil {logs.Error("Pcap_replay Error getting data from main.conf")}
    stapInterface, err := utils.GetKeyValueString("stap", "interface")
    if err != nil {logs.Error("Pcap_replay Error getting data from main.conf")}
    keepPCAP, err := utils.GetKeyValueString("stap", "keepPCAP")
    if err != nil {logs.Error("Pcap_replay Error getting data from main.conf")}
    command, err := utils.GetKeyValueString("execute", "command")  
    if err != nil {logs.Error("Error getting data from main.conf: "+err.Error())}
    param, err := utils.GetKeyValueString("execute", "param")  
    if err != nil {logs.Error("Error getting data from main.conf: "+err.Error())}
    tcpreplay, err := utils.GetKeyValueString("stap", "tcpreplay")  

    //check Stap status
    stapStatus := make(map[string]bool)
    stapStatus,err = PingStap("")
    if err != nil {
        logs.Error("Waiting 60 seconds: Error doing ping to STAP : "+err.Error())
        t,err := utils.GetKeyValueString("loop", "Pcap_replay")
        if err != nil {logs.Error("Search Error: Cannot load node information.")}
        tDuration, err := strconv.Atoi(t)

        time.Sleep(time.Second * time.Duration(tDuration))
    }
    
    //while stap == true, infinite loop will be active
    for stapStatus["stapStatus"]{
        //checking stap for each loop
        stapStatus, err = PingStap("")
        if err != nil {
            logs.Error("Waiting 60 seconds: Error doing ping to STAP : "+err.Error())
            t,err := utils.GetKeyValueString("loop", "stapStatus")
            if err != nil {logs.Error("Search Error: Cannot load node information.")}
            tDuration, err := strconv.Atoi(t)
            time.Sleep(time.Second * time.Duration(tDuration))
            continue
        }

        //checking in_queue path
        files, err := ioutil.ReadDir(inQueue)
        if err != nil {
            logs.Error("Error reading in_queue path: "+err.Error())
            t,err := utils.GetKeyValueString("loop", "ReadDir")
            if err != nil {logs.Error("Search Error: Cannot load node information.")}
            tDuration, err := strconv.Atoi(t)
            time.Sleep(time.Second * time.Duration(tDuration))
            continue
        }
        
        //check files in remote path
        if len(files) == 0 {
            logs.Error("Error Pcap_replay reading files: No files")
            t,err := utils.GetKeyValueString("loop", "remote")
            if err != nil {logs.Error("Search Error: Cannot load node information.")}
            tDuration, err := strconv.Atoi(t)
            time.Sleep(time.Second * time.Duration(tDuration))
            continue
        }
        
        //if there are files in remote path, use tcpreplay
        for _, f := range files{
            // cmd := "tcpreplay -i "+stapInterface+" -t -l 1 "+inQueue+f.Name()
            iface := strings.Replace(tcpreplay, "<IFACE>", stapInterface, -1)
            allValues := strings.Replace(iface, "<NAME>", inQueue+f.Name(), -1)

            _, err := exec.Command(command, param, allValues).Output()
            if err != nil{
                logs.Error("Error exec cmd command "+err.Error())
            }
            if keepPCAP == "false" {
                err = utils.CopyFile(outQueue, inQueue, f.Name(), 1000)
                if err != nil{
                    logs.Error("Error moving file "+err.Error())
                }
            }
            err = utils.RemoveFile(inQueue, f.Name())
            if err != nil {
                logs.Error("Error removing file "+err.Error())
            }
        }
    }
}