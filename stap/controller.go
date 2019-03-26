package stap

import (
    "github.com/astaxie/beego/logs"
    // "godoc.org/golang.org/x/crypto/ssh"
    // "os"
    // "os/exec"
    // "strings"
    // "regexp"
  	// "owlhnode/utils"
  	"owlhnode/database"
	//   "io/ioutil"
	  //"errors"
      //"encoding/json"
      "time"
      "strconv"
    //   "errors"
      //"ssh.CleintConfig"
    //   "code.google.com/p/go.crypto/ssh"
    //   "sync"
    // "runtime"
	// "math/rand"
	// "golang.org/x/crypto/ssh"  
)

const MaxWorkers = 4


func StapInit()(){
    logs.Info("Init Controller Working")
    go Controller()
}

func Controller()() {                     
    var serverOnUUID string
    stapStatus := make(map[string]bool)
    stapStatus = PingStap("")
    var countServers string
    numServers, _ := ndb.Sdb.Query("select count(*) from servers where server_param = \"status\" and server_value = \"true\";")
	defer numServers.Close()
    //load number of servers with status = true
    for numServers.Next(){
        numServers.Scan(&countServers) 
	}
	
    logs.Info("Number of servers ON --> "+countServers)
    i, _ := strconv.Atoi(countServers)
    jobs := make(chan string, i)  
	res := make(chan string,i)  
	isWorking := false
    
    // //create workers 
    // for w := 1; w <= MaxWorkers; w++ {             
    //     logs.Info("loop workers ",w)
    //     go serverTask(w, jobs, res)
    // }
    
    //add UUID servers to jobs channel
	//logs.Debug("Checking Stap server Status before launch goroutines-->"+strconv.FormatBool(stapStatus["stapStatus"]))
    if stapStatus["stapStatus"]{
		isWorking = true
		//create workers 
		for w := 1; w <= MaxWorkers; w++ {             
			logs.Info("loop workers ",w)
			go serverTask(w, jobs, res)
		}
		rows, _ := ndb.Sdb.Query("select server_uniqueid from servers where server_param = \"status\" and server_value = \"true\";")
		defer rows.Close()
        for rows.Next(){
            rows.Scan(&serverOnUUID)
            logs.Warn("Reading query UUID --> "+serverOnUUID) 
            jobs <- serverOnUUID
        }
    }

    //add dinamically to channel the server who had finished their works
    for stapStatus["stapStatus"]{
		var countServers string
		numServers, _ := ndb.Sdb.Query("select count(*) from servers where server_param = \"status\" and server_value = \"true\";")
		defer numServers.Close()
		//load number of servers with status = true
		for numServers.Next(){
			numServers.Scan(&countServers) 
		}
		logs.Debug("NEW LOOP WITH---------------------------------->"+countServers)
		uuid := <-res
        jobs <- uuid 
        stapStatus = PingStap("")
	}
	
	//Kill Servers when STAP stops
	if isWorking{
		rowsKillStap, _ := ndb.Sdb.Query("select server_uniqueid from servers where server_param = \"status\" and server_value = \"true\";")
		defer rowsKillStap.Close()
		logs.Info("Killing servers with status == True") 
		for rowsKillStap.Next(){
			rowsKillStap.Scan(&serverOnUUID)
			// owlh := ndb.GetStapServerInformation(serverOnUUID)
			StopSniffer(serverOnUUID)
		}
	}
	
	close(jobs)
	close(res)
    logs.Info("Workers Closed")
}

func serverTask(id int, jobs <-chan string, res chan<- string) {
    for uuid:=range jobs{
		alive,_ := CheckOwlhAlive(uuid)
        if alive {
            logs.Info("Status  Session: True")
			running, status := GetStatusSniffer(uuid)
			if status {
				logs.Info(running)
			}
        }else{
            logs.Info("Status SSH Session: False")
        }
        time.Sleep(time.Second * 2)
		res <- uuid
    }
}
