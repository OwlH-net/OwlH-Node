package stap

import (
    "github.com/astaxie/beego/logs"
    // "os"
    // "os/exec"
    // "strings"
    // "regexp"
  	// "owlhnode/stap"
  	//"owlhnode/database"
	  // "io/ioutil"
	  //"errors"
      //"encoding/json"
    //   "time"
      "sync"
)

func Init(uuid string)(){
    stapStatus := make(map[string]bool)
    stapStatus = PingStap(uuid)
    for{
        sync.WaitGroup.Add(1)
        go Controller(uuid)
        waitGroup.Wait()
        if !stapStatus["stapStatus"]{
            break
        }
    }
}

func Controller(uuid string, stapStatus map[string]bool)() {
     
    //check if server stap is enabled in our config
    //check if server is reacheble
    //    if disabled - be sure the stap if off
    //    check if server status (CPU, RAM, STORAGE) is ok
    //    check stap status - stop or start as needed.
    //collect pcap files
    
    //1- ping Stap Node
        //2-Lanzar Worker
            //2.1 - Leer Stap servers 
                //3 Por cada SV
                    //3.1 Crear worker
                        //3.1.1 Verificaciones



    logs.Warn(stapStatus["stapStatus"])

    for{
        sync.WaitGroup.Add(1)
        defer func() {
            logs.Info("Destroying worker "+uuid)
            waitGroup.Done()
        }()
        logs.Warn("Stap node is true")
        time.Sleep(time.Second * 1)
        stapStatus = PingStap(uuid)
        break
    }

}



