package stap

import (
    "github.com/astaxie/beego/logs"
    // "godoc.org/golang.org/x/crypto/ssh"
    "os"
    "os/exec"
    // "strings"
    // "regexp"
  	"owlhnode/utils"
  	// "owlhnode/database"
	  "io/ioutil"
      //"encoding/json"
      "time"
    //   "strconv"
    //   "errors"
      //"ssh.CleintConfig"
    //   "code.google.com/p/go.crypto/ssh"
    //   "sync"
    // "runtime"
	// "math/rand"
	// "golang.org/x/crypto/ssh"  
)

func Pcap_replay()() {
	var err error
	loadStap := map[string]map[string]string{}
	loadStap["stap"] = map[string]string{}
    loadStap["stap"]["in_queue"] = ""
	loadStap["stap"]["out_queue"] = ""
	loadStap["stap"]["interface"] = ""
	loadStap,err = utils.GetConf(loadStap)
    inQueue := loadStap["stap"]["in_queue"]
	outQueue := loadStap["stap"]["out_queue"]
	stapInterface := loadStap["stap"]["interface"]
	if err != nil {logs.Error("Error getting path and BPF from main.conf")}
	
	stapStatus := make(map[string]bool)
	stapStatus,err = PingStap("")

	if err != nil {
		logs.Error("Waiting 60 seconds: Error doing ping to STAP : "+err.Error())
		time.Sleep(time.Second * 60)
	}
	
	logs.Debug("Inside the PcapReplay, just before the loop")
	for stapStatus["stapStatus"]{
		//checking stap
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
		if len(files) == 0 {
			logs.Error("Error Pcap_replay reading files: No files")
			time.Sleep(time.Second * 10)
			continue
		}
		for _, f := range files{
			logs.Debug("Pcap_Replay-->"+f.Name())
			cmd := "tcpreplay -i "+stapInterface+" -t -l 1 "+inQueue+f.Name()
			_, err := exec.Command("bash", "-c", cmd).Output()
			if err != nil{
				logs.Error("Error exec cmd command "+err.Error())
			}
			err = os.Rename(inQueue+f.Name(), outQueue+f.Name())
			if err != nil{
				logs.Error("Error moving file "+err.Error())
			}
		}
	}
}