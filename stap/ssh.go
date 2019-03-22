package stap
import (
    "github.com/astaxie/beego/logs"
    // "os"
    // "os/exec"
    "strings"
    "regexp"
  	"owlhnode/utils"
  	// "owlhnode/database"
	  "io/ioutil"
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
    "golang.org/x/crypto/ssh"
)

func RunCMD(cmd string, sshValue *ssh.Session)(status bool, data string){
	// err := sshValue.Run(cmd)
	var outputSSH []byte
	outputSSH, err := sshValue.CombinedOutput(cmd)
	if err!=nil{
		logs.Error("Error executing SSH commands: "+err.Error())    
		return false, ""
	}else{
		logs.Info("No error executing command through SSH connection")    
		return true, string(outputSSH)
	}
}

func owl_connect(stapServer map[string]string)(alive bool, sshValue *ssh.Session){
    loadData := map[string]map[string]string{}
	loadData["stapPubKey"] = map[string]string{}
    loadData["stapPubKey"]["user"] = ""
    loadData["stapPubKey"]["cert"] = ""
    loadData = utils.GetConf(loadData)    
    userSSH := loadData["stapPubKey"]["user"]
    cert := loadData["stapPubKey"]["cert"]

    logs.Warn("Trying to Connect SSH")    

    // //Declare ssh config
    sshConfig := &ssh.ClientConfig{
        User: userSSH,
        Auth: []ssh.AuthMethod{
            PublicKeyFile(cert),
            // PublicKey(pk),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
        Timeout: time.Duration(20)*time.Second,
    }
    logs.Warn("SSH Config declared!! -->"+stapServer["ip"])    
    client, err := ssh.Dial("tcp", stapServer["ip"]+":22", sshConfig)
    if err != nil {
        logs.Error("SSH Dial Error: "+err.Error())    
		return false, nil
    }

    logs.Warn("SSH starting session")    

    // //Launch new session
    session, err := client.NewSession()
	if err != nil {
        logs.Info("New Session Error: "+err.Error())
		session.Close()
		return false, nil
    }
    // defer session.Close()
    logs.Info("New session has been established")
    return true, session //return session
}

func PublicKeyFile(file string) ssh.AuthMethod {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return ssh.PublicKeys(key)
}

func GetStatusSnifferSSH(owlh map[string]string, sshSession *ssh.Session)(status bool, pid string, cpu string, mem string ){
	logs.Info("Checking if sniffer is working in "+owlh["name"]+" - "+owlh["ip"])
	cmd:="top -b -n1 | grep -v sudo | grep tcpdump | awk '{print $1 \",\" $9 \",\" $10}'"
	output := ""
	status, output = RunCMD(cmd, sshSession)

	if regexp.MustCompile(`\d+,\d+\.\d+,\d+\.\d+`+output+`;`) != nil{
		splitValue := strings.Split(output,",")
		pid = splitValue[0]
		cpu = splitValue[1]
		mem = splitValue[2]
		if pid != "" {
			logs.Info("Sniffer PID: "+pid+" is working in "+owlh["name"]+" - "+owlh["ip"])
			return true, pid, cpu, mem
		}
	}
	logs.Info("Sniffer is NOT working for "+owlh["name"]+" - "+owlh["ip"])
	return false, "","",""
}

func GetStatusStorageSSh(owlh map[string]string, sshSession *ssh.Session, folder string)(status bool, storage string, path string){
	logs.Info("Checking if storage is OK in "+owlh["name"]+" - "+owlh["ip"])
	//COMPLETE
}

// func CheckOwlhAliveSSh()(){
	
// }


// func RunSnifferSSh()(){
	
// }

// func StopSnifferSSh()(){
	
// }

// func GetFileListSSH()(){
	
// }

// func OwnerOwlhSSH()(){
	
// }

// func TransportFileSSH()(){
	
// }

// func RemoveFileSSH()(){
	
// }