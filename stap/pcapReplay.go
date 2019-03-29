package stap

import (
    "github.com/astaxie/beego/logs"
    // "godoc.org/golang.org/x/crypto/ssh"
    "os"
    "os/exec"
    // "strings"
    // "regexp"
  	// "owlhnode/utils"
  	// "owlhnode/database"
	  "io/ioutil"
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

func Pcap_replay()() {
	logs.Notice("PCAP_REPLAY")
	inQueue := "/usr/share/owlh/in_queue/"
	outQueue := "/usr/share/owlh/out_queue/"
	logs.Debug("Inside the PcapReplay, just before the loop")
	for{
		files, _ := ioutil.ReadDir(inQueue)
		if len(files) == 0 {
			logs.Error("Error Pcap_replay reading files: No files")
			time.Sleep(time.Second * 10)
		}
		x := 0
		// logs.Info(len(files))
		for a, f := range files{
			x += 1
			logs.Info("Reading local count..>"+strconv.Itoa(x))
			logs.Info("Reading for index 's'-->"+strconv.Itoa(a))
			logs.Debug("Pcap_Replay-->"+f.Name())
			cmd := "tcpreplay -i enp0s3 -t -l 1 "+inQueue+f.Name()
			output, err := exec.Command("bash", "-c", cmd).Output()
			logs.Info(string(output))
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