package stap
import (
    "github.com/astaxie/beego/logs"
    "os"
    // "os/exec"
    "strings"
    "regexp"
  	"owlhnode/utils"
  	"owlhnode/database"
	  "io/ioutil"
	  //"errors"
      //"encoding/json"
	//   "time"
	//   "bytes"
      "strconv"
    //   "errors"
      //"ssh.CleintConfig"
    //   "code.google.com/p/go.crypto/ssh"
    //   "sync"
    // "runtime"
    // "math/rand"
	"golang.org/x/crypto/ssh"
	"github.com/pkg/sftp"
	//"github.com/tmc/scp"
)

func RunCMD(uuid string, cmd string)(status bool, data string){
	logs.Info("Connect ssh from RunCMD")
	var CMDoutput []byte
	alive, sshValue := owlh_connect(uuid)
	if !alive{
		//sshValue.Close()
		return false, ""
	}
	logs.Info("Command to exec --> "+cmd)
	CMDoutput, err := sshValue.Output(cmd)
	sshValue.Close()
	if err!=nil{
		logs.Error("Error executing SSH commands: "+err.Error())
		return false, ""
	}else{
		logs.Info("Command output --> "+string(CMDoutput))
		logs.Info("No error executing command through SSH connection")
		return true, string(CMDoutput)
	}
}

func owlh_connect(uuid string)(alive bool, sshValue *ssh.Session){
    loadData := map[string]map[string]string{}
	loadData["stapPubKey"] = map[string]string{}
    loadData["stapPubKey"]["user"] = ""
    loadData["stapPubKey"]["cert"] = ""
    loadData = utils.GetConf(loadData)
    userSSH := loadData["stapPubKey"]["user"]
    cert := loadData["stapPubKey"]["cert"]

	owlh,err := ndb.GetStapServerInformation(uuid)
	if err != nil {
		logs.Error("Error retrieving stap server information")
	}
	logs.Info("Name: "+owlh["name"]+" IP: "+owlh["ip"])

    // //Declare ssh config
    sshConfig := &ssh.ClientConfig{
        User: userSSH,
        Auth: []ssh.AuthMethod{
            PublicKeyFile(cert),
            // PublicKey(pk),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
        //Timeout: time.Duration(10)*time.Second,
    }
    logs.Warn("SSH Config declared!! -->"+owlh["ip"])
    client, err := ssh.Dial("tcp", owlh["ip"]+":22", sshConfig)
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

func GetStatusSnifferSSH(uuid string)(status bool, pid string, cpu string, mem string ){
	// logs.Info("Checking if sniffer is working in "+owlh["name"]+" - "+owlh["ip"])
	cmd:="top -b -n1 | grep -v sudo | grep tcpdump | awk '{print $1 \",\" $9 \",\" $10}'"
	output := ""
	status, output = RunCMD(uuid,cmd)
	logs.Info("OUTPUT FOR REGEXP: "+output)
	var validOutput = regexp.MustCompile(`\d+,\d+\.\d+,\d+\.\d+`)
	if validOutput.MatchString(output){
		logs.Info("Inside Validation REGEXP")
		splitValue := strings.Split(output,",")
		pid = splitValue[0]
		cpu = splitValue[1]
		mem = splitValue[2]
		if pid != "" {
			// logs.Info("Sniffer PID: "+pid+" is working in "+owlh["name"]+" - "+owlh["ip"])
			logs.Info("Sniffer PID: "+pid+" is working")
			return true, pid, cpu, mem
		}
	}
	logs.Info("Sniffer is NOT working for "+uuid)
	return false, "","",""
}

func GetStatusStorageSSh(uuid string, folder string)(status bool, path string, storage string){
	logs.Info("Checking if storage is OK in "+uuid)
	cmd:="df -h "+folder+" --output=source,pcent | grep -v Filesystem | awk '{print $1 \",\" $2}'"
	output := ""
	status, output = RunCMD(uuid,cmd)
	logs.Info("Output Run CMD -->"+output)
	var validOutput = regexp.MustCompile(`[^,]+,\d+%`)
	if validOutput.MatchString(output) {
	// if regexp.MustCompile(`[^,]+,\d+%`+output+`;`) != nil{
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

func RunSnifferSSH(uuid string)(){
	logs.Info("Launching Sniffer...")
	owlh,err := ndb.GetStapServerInformation(uuid)
	if err != nil {
		logs.Error("Error retrieving stap server information")
	}
	cmd := "nohup sudo tcpdump -i "+owlh["default_interface"]+" -G "+owlh["capture_time"]+" -w "+owlh["pcap_path"]+"`hostname`-%y%m%d%H%M%S.pcap -F "+owlh["filter_path"]+" -z "+owlh["owlh_user"]+" >/dev/null 2>&1 &"
	//logs.Debug(cmd)
	status, output := RunCMD(uuid,cmd)
	if status {
		logs.Info("Sniffer is Running for server "+owlh["name"]+" - "+owlh["ip"]+". Output--> "+output)
	}
}

func StopSnifferSSH(uuid string)(){
	logs.Info(uuid+"Stopping Sniffer...")
	cmd := "ps -ef | grep -v grep | grep tcpdump | awk '{print $2}'"
    output := ""
	_, output = RunCMD(uuid,cmd)
	logs.Info("NUMBER OF PID -->"+output)
	splitValue := strings.Split(output,"\n")
	for value := range splitValue{
		var validOutput = regexp.MustCompile(`^\d+`)
		if validOutput.MatchString(splitValue[value]) {
			cmd := "sudo kill -9 "+splitValue[value]
			_, _ = RunCMD(uuid,cmd)
			logs.Info("Sniffer has been stopped succesfully")
		}else{
			logs.Error("Sniffer is still running...")
		}
	}
}

func GetFileListSSH(uuid string, owlh map[string]string, path string)(list []string){
	//logs.Debug(uuid+" Getting file list...")
	cmd := "find "+owlh["pcap_path"]+"*.pcap -maxdepth 0 -type f -mmin +1|sed 's#.*/##'| awk '{printf $1 \",\"}'"
    output := ""
	_, output = RunCMD(uuid,cmd)

	splitValue := strings.Split(output,",")
	return splitValue
}

func OwnerOwlhSSH(uuid string, owlh map[string]string, fileRemote string)(){
	logs.Info(uuid+" Change remote file owner: "+ fileRemote)
	var validOutput = regexp.MustCompile(`\.pcap+`)
	if validOutput.MatchString(fileRemote) {
		// cmd := "sudo chown "+owlh["owlh_user"]+" "+fileRemote
		cmd := "sudo chown "+owlh["owlh_user"]+" "+owlh["pcap_path"]+fileRemote
		logs.Alert("Command for change owner file--> "+cmd)
		RunCMD(uuid,cmd)
	}
}

func TransportFileSSH(uuid string, owlh map[string]string, file string)(){
	logs.Info(uuid+" Transport remote file: "+ file+" to local machine")
	var validOutput = regexp.MustCompile(`\.pcap+`)
	if validOutput.MatchString(file) {

		//exec scp local

		//cmd := "sudo scp -r "+owlh["owlh_user"]+"@"+owlh["pcap_path"]+":"+owlh["pcap_path"]+" "+owlh["local_pcap_path"]+""
		status := SftpCMD(uuid, owlh["pcap_path"]+file, owlh["local_pcap_path"]+file)
		logs.Info("Output CMD TransportFileSSH: "+file+" - "+strconv.FormatBool(status))
	}
}

func RemoveFileSSH(uuid string, owlh map[string]string, file string)(){
	logs.Info(uuid+" Remove remote files")
	var validOutput = regexp.MustCompile(`\.pcap+`)
	if validOutput.MatchString(file) {
		cmd := "sudo rm "+file+""
		_, output := RunCMD(uuid,cmd)
		logs.Info("Output CMD RemoveFileSSH: "+ output)
	}
}


func SftpCMD(uuid string, srcFile string, dstFile string)(status bool){
	logs.Notice("SFTP_CMD")
	logs.Info("Connect sftp SftpCMD")
	alive, sshClient:= owlh_connect_client(uuid)
	if !alive{
		return false
	}

	logs.Notice("NewClient")
	sftp, err := sftp.NewClient(sshClient)
	if err != nil {
		logs.Error(err)
		return false
	}
	defer sftp.Close()

	// Open the source file
	logs.Notice("Open Remote File")
	remoteFile, err := sftp.Open(srcFile)
	if err != nil {
		logs.Error(err)
		return false
	}
	defer remoteFile.Close()

	// Create the destination file
	logs.Notice("Create DST file")
	localFile, err := os.Create(dstFile)
	if err != nil {
		logs.Error(err)
		return false
	}
	defer localFile.Close()

	// Copy the file
	logs.Notice("Copy file")
	remoteFile.WriteTo(localFile)

	//remove remote file
	logs.Notice("Delete remote file")
	err = sftp.Remove(srcFile)
	if err != nil {
		logs.Error(err)
		return false
	}
	
	return true
}

func owlh_connect_client(uuid string)(alive bool, sshClient *ssh.Client){
    loadData := map[string]map[string]string{}
	loadData["stapPubKey"] = map[string]string{}
    loadData["stapPubKey"]["user"] = ""
    loadData["stapPubKey"]["cert"] = ""
    loadData = utils.GetConf(loadData)
    userSSH := loadData["stapPubKey"]["user"]
    cert := loadData["stapPubKey"]["cert"]

	owlh,err := ndb.GetStapServerInformation(uuid)
	if err != nil {
		logs.Error("Error retrieving stap server information")
	}
	logs.Info("Name: "+owlh["name"]+" IP: "+owlh["ip"])

    // //Declare ssh config
    sshConfig := &ssh.ClientConfig{
        User: userSSH,
        Auth: []ssh.AuthMethod{
            PublicKeyFile(cert),
            // PublicKey(pk),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
        //Timeout: time.Duration(10)*time.Second,
    }
    logs.Warn("SSH Config declared!! -->"+owlh["ip"])
    client, err := ssh.Dial("tcp", owlh["ip"]+":22", sshConfig)
    if err != nil {
        logs.Error("SSH Dial Error: "+err.Error())
		return false, nil
    }
	return true, client
}