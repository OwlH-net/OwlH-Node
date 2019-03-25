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
	//   "bytes"
    //   "strconv"
    //   "errors"
      //"ssh.CleintConfig"
    //   "code.google.com/p/go.crypto/ssh"
    //   "sync"
    // "runtime"
    // "math/rand"
    "golang.org/x/crypto/ssh"
)

func RunCMD(cmd string, sshValue *ssh.Session)(status bool, data string){
	logs.Info(cmd)
	var outputSSH []byte
	var err error
	outputSSH, err = sshValue.Output(cmd)

	// var stdoutBuf bytes.Buffer
	// var stdinBuf bytes.Buffer
	// sshValue.Stdout = &stdoutBuf
	// sshValue.Stdin = &stdinBuf
	// err := sshValue.Run(cmd)
	// logs.Info("Std in --> "+stdinBuf.String())    

	// stdin, _ := sshValue.StdinPipe()
	// stdout, _ := sshValue.StdoutPipe()
	// stdout, err := sshValue.Run(cmd)
	if err!=nil{
		logs.Error("Error executing SSH commands: "+err.Error())    
		return false, ""
	}else{
		logs.Info("No error executing command through SSH connection")    
		return true, string(outputSSH)
	}
}

func owlh_connect(stapServer map[string]string)(alive bool, sshValue *ssh.Session){
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

    //Launch new session
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

func GetStatusStorageSSh(owlh map[string]string, sshSession *ssh.Session, folder string)(status bool, path string, storage string){
	logs.Info("Checking if storage is OK in "+owlh["name"]+" - "+owlh["ip"])
	cmd:="df -h "+folder+" --output=source,pcent | grep -v Filesystem | awk '{print $1 \",\" $2}'"
	output := ""
	status, output = RunCMD(cmd, sshSession)
	logs.Info("Output Run CMD -->"+output)
	if regexp.MustCompile(`[^,]+,\d+%`+output+`;`) != nil{
		splitValue := strings.Split(output,",")
		// // // //splitValue := strings.Fields(output)
		logs.Info("sPLITvALUE: "+splitValue[0]+" ///// "+splitValue[1])
		// path = splitValue[0]
		// storage = splitValue[1]
		// logs.Info("Device: "+path+" Percentage: "+storage)
		return true, splitValue[0], splitValue[1]
	}
	return false, "", ""
	// return true, "/etc/","18%"
}

func RunSnifferSSH(owlh map[string]string, sshSession *ssh.Session)(){
	logs.Info("Launching Sniffer...")
	cmd := "nohup sudo tcpdump -i "+owlh["interface"]+" -G "+owlh["capture"]+" -w "+owlh["pcap_path"]+"`hostname`-%%y%%m%%d%%H%%M%%S.pcap -F "+owlh["filter_path"]+" -z "+owlh["user"]+" >/dev/null 2>&1 &"
	status, output := RunCMD(cmd, sshSession)
	if status {
		logs.Info("Sniffer is Running for server "+owlh["name"]+" - "+owlh["ip"]+". Output--> "+output)	
	}
}

func StopSnifferSSH(sshSession *ssh.Session)(){
	logs.Info("Stopping Sniffer...")
	cmd := "ps -ef | grep -v grep | grep tcpdump | awk '{printf $2}'"
    output := ""
	_, output = RunCMD(cmd,sshSession)
	if regexp.MustCompile(`^\d+$,`+output) != nil {
		cmd := "sudo kill -9 "+output
		_, _ = RunCMD(cmd,sshSession)
		logs.Info("Sniffer has been stopped succesfully")

	}else{
		logs.Info("Sniffer is still running...")
	}
}

// func GetFileListSSH()(){
	
// }

// func OwnerOwlhSSH()(){
	
// }

// func TransportFileSSH()(){
	
// }

// func RemoveFileSSH()(){
	
// }