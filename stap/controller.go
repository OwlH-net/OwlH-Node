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
)

const MaxWorkers = 4


func StapInit()(){
    logs.Info("Init Controller Working")
    go Controller()
}

func Controller()() {                     
    var serverOnUUID string
    var countServers string
    stapStatus := make(map[string]bool)
    stapStatus = PingStap("")
    numServers, _ := ndb.Sdb.Query("select count(*) from servers where server_param = \"status\" and server_value = \"true\";")

    //load number of servers with status = true
    for numServers.Next(){
        numServers.Scan(&countServers) 
    }
    logs.Info("Number of servers ON --> "+countServers)
    i, _ := strconv.Atoi(countServers)
    jobs := make(chan string, i)  
    res := make(chan string,i)  
    
    //create workers 
    for w := 1; w <= MaxWorkers; w++ {             
        logs.Info("loop workers ",w)
        go serverTask(w, jobs, res)
    }
    
    //add UUID servers to jobs channel
    if stapStatus["stapStatus"]{
        rows, _ := ndb.Sdb.Query("select server_uniqueid from servers where server_param = \"status\" and server_value = \"true\";")
        for rows.Next(){
            rows.Scan(&serverOnUUID)
            logs.Warn("Reading query UUID --> "+serverOnUUID) 
            jobs <- serverOnUUID
        }
    }

    //add dinamically to channel the server who had finished their works
    for stapStatus["stapStatus"]{
        uuid := <-res
        jobs <- uuid 
        stapStatus = PingStap("")
    }
    logs.Info("Closing workers")
}

func serverTask(id int, jobs <-chan string, res chan<- string) {
    for uuid:=range jobs{
        alive,sshStat := CheckOwlhAlive(uuid)
        if alive {
            logs.Info("Status SSH Session: True")
			running, status := GetStatusSniffer(uuid, sshStat)
			if status {
				logs.Info(running)
			}
        }else{
            logs.Info("Status SSH Session: False")
        }
        time.Sleep(time.Second * 5)
		res <- uuid
		defer sshStat.Close()
    }
}