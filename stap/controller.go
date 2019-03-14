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
      "strconv"
    //   "sync"
    // "runtime"
    // "math/rand"
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
                        
    var serverOnUUID string
    var countServers string
    stapStatus := make(map[string]bool)
    stapStatus = PingStap("")
    numServers, _ := ndb.Sdb.Query("select count(*) from servers where server_param = \"status\" and server_value = \"true\";")
    for numServers.Next(){
        numServers.Scan(&countServers) 
    }
    logs.Info("Number of servers ON --> "+countServers)
    i, _ := strconv.Atoi(countServers)
    jobs := make(chan string, i)
    
    
    
    //crear workers 
    for w := 1; w <= MaxWorkers; w++ {             
        logs.Info("loop workers ",w)
        go serverTask(w, jobs)
    }

    //loop
    for stapStatus["stapStatus"]{
        //query
        rows, _ := ndb.Sdb.Query("select server_uniqueid from servers where server_param = \"status\" and server_value = \"true\";")
        for rows.Next(){
            rows.Scan(&serverOnUUID)
            logs.Warn("Reading query UUID --> "+serverOnUUID)
            jobs <- serverOnUUID
        }
        stapStatus = PingStap("")
        //time.Sleep(time.Second * 5)
    }
        close(jobs)
}

func serverTask(id int, jobs <-chan string) {
    //var jobServer map[string]string
    for job := range jobs {
        var jobServer string
        uuid := job
        logs.Info(" ---- Stap Server Task - looking for servers ---- "+uuid)

        ip, _ := ndb.Sdb.Query("select server_value from servers where server_param = \"ip\" and server_uniqueid = \""+uuid+"\";")
        for ip.Next(){
            ip.Scan(&jobServer)
        }
    
        //delay := rand.Intn(15)
        time.Sleep(time.Second * 5)
        logs.Info("--------------"+string(jobServer))
        <-jobs
        //alive, ssh := CheckOwlhAlive(job)
    }
}

func CheckOwlhAlive(owlh map[string]string)(data bool, ssh string){
    logs.Info("Stap Server Task "+owlh["name"]+" -- "+owlh["ip"])
    alive, ssh = owl_connect(owlh)
    if alive{
		logs.Info("ALIVE Stap Server Task "+owlh["name"]+" -- "+owlh["ip"])
        return true, ssh
	}
	logs.Error("NOT ALIVE Stap Server Task "+owlh["name"]+" -- "+owlh["ip"])
	return false, ""
}

func owl_connect(owlh map[string]string)(data bool, ssh string){
    //connect by ssh to server
    
}