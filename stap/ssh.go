package stap
import (
    "github.com/astaxie/beego/logs"
    // "os"
    // "os/exec"
    // "strings"
    // "regexp"
  	"owlhnode/utils"
  	// "owlhnode/database"
	  "io/ioutil"
	  //"errors"
      //"encoding/json"
      "time"
    //   "strconv"
    //   "errors"
      //"ssh.CleintConfig"
    //   "code.google.com/p/go.crypto/ssh"
    //   "sync"
    // "runtime"
    // "math/rand"
    "golang.org/x/crypto/ssh"
)


// func RunCMD(cmd string, ssh string)(status bool, data string){
//     try:
//         stdin, stdout, stderr = ssh.exec_command(cmd)
//         output = ""
//         for l in stdout :
//             output = l.strip()
//         for l in stderr:
//             print("stderr : %s" % l.strip())
//         return True, output
//     except Exception as inst:
//         flogger("Oops!  there was a problem: %s" % str(inst),"WARNING")
// 	return False, ""
// }




func owl_connect(stapServer map[string]string)(alive bool, sshValue *ssh.Session){
    
    // var hostKey ssh.PublicKey
    //read public key
    loadData := map[string]map[string]string{}
	loadData["stapPubKey"] = map[string]string{}
    loadData["stapPubKey"]["user"] = ""
    loadData["stapPubKey"]["cert"] = ""
    loadData = utils.GetConf(loadData)    
    userSSH := loadData["stapPubKey"]["user"]
    cert := loadData["stapPubKey"]["cert"]
    //pk, _:=ioutil.ReadFile(cert)

    //return true, "ssh"

    logs.Warn("Trying to Connect SSH")    

    // //ParsePrivateKey
    //signer, err := ssh.ParsePrivateKey(pk)

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
    // //Configure session with ssh config
    // client, err := ssh.Dial("tcp", stapServer["ip"]+":22", sshConfig)
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
    defer session.Close()
    logs.Info("New session has been established")
    logs.Info("Time to close the session")
    //session.Close()
    
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

// func GetStatusSnifferSSH(owlh map[string]string, ssh string)(status bool, pid string, cpu string, mem string ){
// 	logs.Info("Checking if sniffer is working in "+owlh["name"]+" - "+owlh["ip"])
// 	cmd:="top -b -n1 | grep -v sudo | grep tcpdump | awk \'{print $1 "," $9 "," $10}\'"
// 	output := ""
// 	status, output := RunCMD(cmd, ssh)
// 	if regexp.MustCompile(`"\d+,\d+\.\d+,\d+\.\d+"`+output) != nil{
// 		pid, cpu, mem := strings.Split(output,",")
// 		if pid != nil {
// 			logs.Info("Sniffer is working in "+owlh["name"]+" - "+owlh["ip"])
// 			return true, pid, cpu, mem
// 		}
// 	}
// 	logs.Info("Sniffer is NOT working in "+owlh["name"]+" - "+owlh["ip"])
// 	return false, "", "", ""
// }

// func CheckOwlhAliveSSh()(){
	
// }

// func GetStatusStorageSSh()(){
	
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