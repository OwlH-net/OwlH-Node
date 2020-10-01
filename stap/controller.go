package stap

import (
    "github.com/astaxie/beego/logs"
    "strings"
    "regexp"
    "owlhnode/database"
    "owlhnode/utils"
    "time"
    "strconv"
    "runtime"
)

//init PcapReplay and Controller for concurrency
func StapInit()(){
    go Pcap_replay()
    go Controller()
}

//concurrency for Software TAP
func Controller()() {   
    logs.Info("Init Controller Working")                
    var serverOnUUID string
    var err error
    stapStatus := make(map[string]bool)
    stapStatus, err = PingStap("")
    if err != nil {
        logs.Error("Error doing ping to STAP : "+err.Error())
        logs.Error("Waiting 60 seconds...")

        t,err := utils.GetKeyValueString("loop", "stap")
        if err != nil {logs.Error("Search Error: Cannot load node information.")}
        tDuration, err := strconv.Atoi(t)
        time.Sleep(time.Second * time.Duration(tDuration))
    }

    //load number of servers with status = true
    var countServers string
    numServers, err := ndb.Sdb.Query("select count(*) from servers where server_param = \"status\" and server_value = \"true\";")
    if err != nil {
        logs.Error("Error query counting stap servers : "+err.Error())
    }
    defer numServers.Close()
    for numServers.Next(){
        numServers.Scan(&countServers) 
    }
    
    logs.Info("Number of servers ON --> "+countServers)
    i, err := strconv.Atoi(countServers)
    if err != nil {
        logs.Error("Error converting to int the number of stap servers : "+err.Error())
    }
    jobs := make(chan string, i)  
    res := make(chan string,i)  
    isWorking := false
    
    //add UUID servers to jobs channel
    if stapStatus["stapStatus"]{
        //number of cores -1 for concurrency
        var MaxWorkers int
        if runtime.GOMAXPROCS(runtime.NumCPU()) == 1 {
            MaxWorkers := 1
            logs.Info("CORE FOR CONCURRENCY: "+strconv.Itoa(MaxWorkers))
        }else{    
            MaxWorkers := runtime.GOMAXPROCS(runtime.NumCPU())-1
            logs.Info(strconv.Itoa(MaxWorkers)+" CORE FOR CONCURRENCY")
        }
        isWorking = true

        //create workers 
        for w := 0; w <= MaxWorkers; w++ {             
            logs.Info("loop workers ",w)
            go serverTask(w, jobs, res)
        }
        rows, err := ndb.Sdb.Query("select server_uniqueid from servers where server_param = \"status\" and server_value = \"true\";")
        if err != nil {
            logs.Error("Error query counting stap servers : "+err.Error())
        }
        defer rows.Close()
        for rows.Next(){
            rows.Scan(&serverOnUUID)
            logs.Warn("Reading query UUID --> "+serverOnUUID) 
            jobs <- serverOnUUID
        }
    }

    //add dinamically to channel the server who had finished their works
    var validOutput = regexp.MustCompile(`error:`)
    for stapStatus["stapStatus"]{
        uuid := <-res
        if validOutput.MatchString(uuid){
            splitValue := strings.Split(uuid,":")
            ErrorStapServer(splitValue[1])
        }else{
            jobs <- uuid 
        }
        stapStatus,err = PingStap("")
    }
    
    
    //Kill Servers when STAP stops
    if isWorking{
        rowsKillStap, err := ndb.Sdb.Query("select server_uniqueid from servers where server_param = \"status\" and server_value = \"true\";")
        if err != nil {
            logs.Error("Error query status stap servers for stop all servers: "+err.Error())
        }
        defer rowsKillStap.Close()
        logs.Info("Killing servers with status == True") 
        for rowsKillStap.Next(){
            rowsKillStap.Scan(&serverOnUUID)
            StopSniffer(serverOnUUID)
        }
    }
    close(jobs)
    logs.Info("Workers Closed")
}

//Launch a task for each worker with the server alocated on jobs channel
func serverTask(id int, jobs <-chan string, res chan<- string) {
    for uuid:=range jobs{
        alive,_ := CheckOwlhAlive(uuid)
        if alive {
            logs.Alert("Status  Session: True")
            running, status := GetStatusSniffer(uuid)
            if running {
                logs.Info("TCPDUMP is running!!")
                if !status {
                    logs.Info("Something is wrong with the system...")
                    StopSniffer(uuid)
                }
            }else if status{
                logs.Info("Start Sniffer!!")
                RunSniffer(uuid)
            }
            GetFileList(uuid)
            res <- uuid
        }else{
            logs.Info("Status SSH Session: False")
            res <- "error:"+uuid
            
        }
    }
}
