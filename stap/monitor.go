package stap
import (
    "github.com/astaxie/beego/logs"
    // "os"
    // "os/exec"
    "strings"
	"regexp"
  	"owlhnode/database"
	//   "io/ioutil"
	  //"errors"
      //"encoding/json"
    //   "time"
      "strconv"
    //   "errors"
      //"ssh.CleintConfig"
    //   "code.google.com/p/go.crypto/ssh"
    //   "sync"
    // "runtime"
    // "math/rand"
    "golang.org/x/crypto/ssh"
)


func CheckOwlhAlive(uuid string)(alive bool, sshSession *ssh.Session){
	// owlh := ndb.GetStapServerInformation(uuid)
  	alive, sshSession = owlh_connect(uuid)
  	logs.Info("Stap Server Task with uuid: "+uuid)
  	if alive{
	    logs.Info("ALIVE Stap Server "+uuid)
	    return true, sshSession
  	}
  	logs.Error("NOT ALIVE Stap Server "+uuid)
  	return false, nil
}
      
func GetStatusSniffer(uuid string)(running bool, status bool){
	owlh := ndb.GetStapServerInformation(uuid)
	logs.Info("Checking Sniffer status for uuid: "+uuid)
	  
	running, pid, cpu, mem := GetStatusSnifferSSH(uuid)
	cpuStatus := GetStatusCPU(owlh,cpu,uuid)
	memStatus := GetStatusMEM(owlh,mem,uuid)
	storageStatus := GetStatusStorage(owlh,uuid)

	logs.Alert("Checking "+owlh["name"]+" - "+owlh["ip"]+" - PID:"+pid+" CPU: "+strconv.FormatBool(cpuStatus)+" MEM: "+strconv.FormatBool(memStatus)+" STORAGE: "+strconv.FormatBool(storageStatus))
  	if cpuStatus && memStatus && storageStatus {
  	    return running, true
	}
	return running, false
}

func GetStatusCPU(owlh map[string]string, cpu string, uuid string)(status bool){
	var validCPU = regexp.MustCompile(`(\d+)`+cpu+`;`)
	if validCPU == nil{
		return false
	}
	localCPU, _ := strconv.ParseFloat(cpu, 64)
	ddbbCPU, _ := strconv.ParseFloat(owlh["max_cpu"], 64)
	logs.Error("Check CPU for "+owlh["name"]+" - "+owlh["ip"])
	if localCPU>ddbbCPU{
		logs.Error("SNIFFER -> Too much CPU on "+owlh["name"]+" - "+owlh["ip"])
		StopSniffer(uuid)
		return false
	}
	return true
}

func GetStatusMEM(owlh map[string]string, mem string, uuid string)(status bool){
	var validMEM = regexp.MustCompile(`(\d+)`+mem+`;`)
	if validMEM == nil{
		return false
	}
	localMEM, _ := strconv.ParseFloat(mem, 64)
	ddbbMEM, _ := strconv.ParseFloat(owlh["max_mem"], 64)
	logs.Error("Check MEM for "+owlh["name"]+" - "+owlh["ip"])
	if localMEM>ddbbMEM{
		logs.Error("SNIFFER -> Too much MEM on "+owlh["name"]+" - "+owlh["ip"])
		StopSniffer(uuid)
		return false
	}
	return true
}	

func GetStatusStorage(owlh map[string]string, uuid string)(status bool){
	var pcapPath = owlh["pcap_path"]
	status, path, storage := GetStatusStorageSSh(uuid, pcapPath)
	if status {
		localFunc := strings.Replace(storage, "%", "", -1)
		localStorage, _ := strconv.ParseFloat(localFunc, 64)
		ddbbStorage, _ := strconv.ParseFloat(owlh["max_storage"], 64)
		if localStorage > ddbbStorage {
			logs.Error("SNIFFER -> Too much STORAGE on "+owlh["name"]+" - "+owlh["ip"]+" on path: "+path)
			StopSniffer(uuid)
			return false
		}
		logs.Error("SNIFFER -> STORAGE status is correct!!")
		return true
	}
	return false
}	

func RunSniffer(uuid string)(){
	// RunSnifferSSH(owlh, ssh, conf("default"), conf("capture_time"), conf("pcap_path"), conf("filter_path"), conf("owlh_user"))
	RunSnifferSSH(uuid)
}

func StopSniffer(uuid string)(){
	StopSnifferSSH(uuid)
}

func GetFileList(uuid string)(){
	owlh := ndb.GetStapServerInformation(uuid)
	logs.Error("Get file list for "+owlh["name"]+" - "+owlh["ip"])
	file_list := GetFileListSSH(uuid,owlh, owlh["pcap_path"])
	// sftp := openSftpSSH()
	for file := range file_list {
		var validOutput = regexp.MustCompile(`\.pcap+`)
		if validOutput.MatchString(file_list[file]) {
		logs.Debug(file_list[file])
		logs.Info("Change remote file owned")
		OwnerOwlh(uuid, owlh, file_list[file])
		
		logs.Info("Copy full directory using SCP command")
		TransportFile(uuid, owlh, file_list[file])
		
		// logs.Info("Delete remote files")
		// RemoveFile(uuid, owlh, file_list[file])

		logs.Warn("File list completed!")
		}
	}
}

func OwnerOwlh(uuid string, owlh map[string]string, fileRemote string)(){
	logs.Error("Set "+owlh["name"]+" - "+owlh["ip"]+" as owner of file "+owlh["owlh_user"])
	OwnerOwlhSSH(uuid, owlh, fileRemote)
}

func TransportFile(uuid string, owlh map[string]string, file string)(){
	logs.Error("Get file "+owlh["local_pcap_path"]+" from "+owlh["name"]+" - "+owlh["ip"])
	TransportFileSSH(uuid, owlh, file)
}

func RemoveFile(uuid string, owlh map[string]string, file string)(){
	logs.Error("Remove file "+owlh["local_pcap_path"]+" from "+owlh["name"]+" - "+owlh["ip"])
	RemoveFileSSH(uuid, owlh, file)
}







