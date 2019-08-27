package stap

import (
    "github.com/astaxie/beego/logs"
    "os/exec"
      "owlhnode/utils"
    "io/ioutil"
    "time"
)

func Pcap_replay()() {
    var err error
    //load in_, out_queue and interface from main.conf
    loadStap := map[string]map[string]string{}
    loadStap["stap"] = map[string]string{}
    loadStap["stap"]["in_queue"] = ""
    loadStap["stap"]["out_queue"] = ""
    loadStap["stap"]["interface"] = ""
    loadStap["stap"]["keepPCAP"] = ""
    loadStap,err = utils.GetConf(loadStap)
    inQueue := loadStap["stap"]["in_queue"]
    outQueue := loadStap["stap"]["out_queue"]
    stapInterface := loadStap["stap"]["interface"]
    keepPCAP := loadStap["stap"]["keepPCAP"]
    if err != nil {
        logs.Error("Pcap_replay Error getting data from main.conf")
    }
    
    //check Stap status
    stapStatus := make(map[string]bool)
    stapStatus,err = PingStap("")
    if err != nil {
        logs.Error("Waiting 60 seconds: Error doing ping to STAP : "+err.Error())
        time.Sleep(time.Second * 60)
    }
    
    //while stap == true, infinite loop will be active
    for stapStatus["stapStatus"]{
        //checking stap for each loop
        stapStatus, err = PingStap("")
        if err != nil {
            logs.Error("Waiting 60 seconds: Error doing ping to STAP : "+err.Error())
            time.Sleep(time.Second * 60)
            continue
        }

        //checking in_queue path
        files, err := ioutil.ReadDir(inQueue)
        if err != nil {
            logs.Error("Error reading in_queue path: "+err.Error())
            time.Sleep(time.Second * 60)
            continue
        }
        
        //check files in remote path
        if len(files) == 0 {
            logs.Error("Error Pcap_replay reading files: No files")
            time.Sleep(time.Second * 10)
            continue
        }
        
        //if there are files in remote path, use tcpreplay
        for _, f := range files{
            logs.Debug("Pcap_Replay-->"+f.Name())
            cmd := "tcpreplay -i "+stapInterface+" -t -l 1 "+inQueue+f.Name()
            _, err := exec.Command("bash", "-c", cmd).Output()
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