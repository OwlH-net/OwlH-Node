package stap

import (
    "github.com/astaxie/beego/logs"
    // "os"
    // "os/exec"
    // "strings"
    // "regexp"
  	// "owlhnode/stap"
  	"owlhnode/database"
	  // "io/ioutil"
	  //"errors"
      //"encoding/json"
      "time"
    //   "sync"
    // "runtime"
)

const MaxWorkers = 4

func StapInit()(){
    // stapStatus := make(map[string]bool)
    // stapStatus = PingStap(uuid)
    
    //for{
        logs.Info("Init Controller Working")
        // channel := make(chan int)
        // var countServers int
        // rowsCount, _ := ndb.Sdb.Query("select count(*) from servers where server_param = \"status\" and server_value = \"true\";")
        // rowsCount.Scan(&countServers)

        // logs.Error("Before Channel load")
        // for i:=0;i<countServers;i++{
        //     channel<-i
        //     logs.Info("No elements on channel: "+string(i))
        // }
        
        // logs.Error("Before concurrency starts")
        // for j:=0;j<MaxWorkers;j++{
        //     waitGroup.Add(1)
        //     go Controller(channel)
        // }
        // close(channel)
        // waitGroup.Wait()
        // defer rowsCount.Close()
        // logs.Error("Ending...")




        //waitGroup.Add(1)
        go Controller()
        // waitGroup.Add(1)
        // go Controller(uuid, stapStatus)
        // stapStatus = PingStap(uuid)
        // if !stapStatus["stapStatus"]{
        //     break
        // }
        //time.Sleep(time.Second * 2)
    //}
    //waitGroup.Wait()
    
}

func Controller()() {
    //runtime.GOMAXPROCS(runtime.NumCPU())
     
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

    stapStatus := make(map[string]bool)
    stapStatus = PingStap("")
    jobs := make(chan string, 100)
    results := make(chan string, 100)
    //results := make(chan string, 100)
    var serverOnUUID string
    
    // for w := 1; w <= MaxWorkers; w++ { //Number of workers
    //     logs.Info("loop workers ",w)
    //     go serverTask(w, jobs)
    // }

    // currentWorkers := 0
    // for i:=0 ; i< MaxWorkers ; i++{
    //     for rows.Next(){
    //         rows.Scan(&serverOnUUID)
    //         server = true
    //         for server {
    //             if currentWorkers < MaxWorkers {
    //                 waitGroup.Add(1)
    //                 logs.Error("Worker added No:"+serverOnUUID)
    //                 //data <- ("Worker "+serverOnUUID)
    //                 go Worker(serverOnUUID)
    //                 currentWorkers ++
    //                 server = false
    //             }else{
    //                 //¿?
    //             }
    //             currentWorkers --
                
    //         }
            
    //     }
    // }

    for stapStatus["stapStatus"]{
        for w := 1; w <= MaxWorkers; w++ { //Number of workers
            logs.Info("loop workers ",w)
            go serverTask(w, jobs, results)
        }
        logs.Error("Inside the infinite loop")
        rows, _ := ndb.Sdb.Query("select server_uniqueid from servers where server_param = \"status\" and server_value = \"true\";")
        for rows.Next(){
            rows.Scan(&serverOnUUID)
            logs.Warn(serverOnUUID)
            jobs <- serverOnUUID
        }
        //collect all severs with ON status
        // for server := range servers{
            //     jobs <- uuid
            // }
            
        stapStatus = PingStap("")
        time.Sleep(time.Second * 2)
    }
        
        for a := 1; a <= 5; a++ {
            <-results
        }
        close(jobs)
    }

func serverTask(id int, jobs <-chan string, results <-chan string) {

    //do all things
    for j := range jobs {
        logs.Warn("Doing something ",id)
        results <- j
    }
    // alive, ssh = CheckOwlhAlive(owlh)
    // if alive{
    //     //FlockLogger(">>> as owlh name -> "+owlh["name"]+" is alive with check status")
    //     running, status_ok := GetStatusSniffer(owlh, ssh)
    //     //FlockLogger(">>> Running "+running+", Status "+status_ok)
    //     if running {
    //         if !status_ok{
    //             StopSniffer(owlh, ssh)
    //         }
    //     } else if status_ok {
    //         RunSniffer(owlh, ssh)
    //     }
    //     GetFileList(owlh, ssh)
    //     ssh.Close()
    // }
}
