package stap
import (
    "github.com/astaxie/beego/logs"
    // "os"
    // "os/exec"
    // "strings"
    // "regexp"
  	"owlhnode/utils"
  	"owlhnode/database"
	//   "io/ioutil"
	  //"errors"
      //"encoding/json"
    //   "time"
    //   "strconv"
    //   "errors"
      //"ssh.CleintConfig"
    //   "code.google.com/p/go.crypto/ssh"
    //   "sync"
    // "runtime"
    // "math/rand"
)


func RunCMD(cmd string, ssh string)(status bool, data string){
    try:
        stdin, stdout, stderr = ssh.exec_command(cmd)
        output = ""
        for l in stdout :
            output = l.strip()
        for l in stderr:
            print("stderr : %s" % l.strip())
        return True, output
    except Exception as inst:
        flogger("Oops!  there was a problem: %s" % str(inst),"WARNING")
	return False, ""
}




func owl_connect(stapServer map[string]string)(alive bool, ssh string){
    //read public key
    loadData := map[string]map[string]string{}
	loadData["stapPubKey"] = map[string]string{}
	loadData["stapPubKey"]["publicKey"] = ""
    loadData = utils.GetConf(loadData)    
    key := loadData["stapPubKey"]["stapPubKey"]
    pk, _:=ioutil.ReadFile(key)

    logs.Info(pk)

    return true, "ssh"


    // //ParsePrivateKey
    // signer, err := ssh.ParsePrivateKey(pk)

    // //Declare ssh config
    // sshConfig := &ssh.ClientConfig{
    //     User: stapServer["ip"]
    //     Auth: []ssh.AuthMethod{
    //         ssh.PublicKey(signer),
    //     },
    // }
    
    // //Configure session with ssh config
    // client, err := ssh.Dial("tcp", stapServer["ip"]+":22", sshConfig)
    // if err != nil {
	// 	return false, ""
    // }

    // //Launch new session
    // session, err := client.NewSession()
	// if err != nil {
	// 	client.Close()
	// 	return false, ""
    // }
    
    // return true, session
}

func GetStatusSnifferSSH(owlh map[string]string, ssh string)(status bool, pid string, cpu string, mem string ){
	logs.Info("Checking if sniffer is working in "+owlh["name"]+" - "+owlh["ip"])
	cmd:="top -b -n1 | grep -v sudo | grep tcpdump | awk \'{print $1 "," $9 "," $10}\'"
	output := ""
	status, output := RunCMD(cmd, ssh)
	if regexp.MustCompile(`"\d+,\d+\.\d+,\d+\.\d+"`+output) != nil{
		pid, cpu, mem := strings.Split(output,",")
		if pid != nil {
			logs.Info("Sniffer is working in "+owlh["name"]+" - "+owlh["ip"])
			return true, pid, cpu, mem
		}
	}
	logs.Info("Sniffer is NOT working in "+owlh["name"]+" - "+owlh["ip"])
	return false, "", "", ""
}

func CheckOwlhAliveSSh()(){
	
}

func GetStatusStorageSSh()(){
	
}

func RunSnifferSSh()(){
	
}

func StopSnifferSSh()(){
	
}

func GetFileListSSH()(){
	
}

func OwnerOwlhSSH()(){
	
}

func TransportFileSSH()(){
	
}

func RemoveFileSSH()(){
	
}