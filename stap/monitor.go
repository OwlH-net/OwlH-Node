package stap
import (
    "github.com/astaxie/beego/logs"
    // "os"
    // "os/exec"
    // "strings"
	"regexp"
  	// "owlhnode/utils"
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
  	stapServer := GetStapServerInformation(uuid)
  	logs.Info("Stap Server Task "+stapServer["name"]+" -- "+stapServer["ip"])
  	alive, sshSession = owl_connect(stapServer)
  	if alive{
	    logs.Info("ALIVE Stap Server Task "+stapServer["name"]+" -- "+stapServer["ip"])
	    return true, sshSession
  	}
  	logs.Error("NOT ALIVE Stap Server Task "+stapServer["name"]+" -- "+stapServer["ip"])
  	return false, nil
}
      
func GetStatusSniffer(uuid string, sshSession *ssh.Session)(running string, status bool){
	owlh := GetStapServerInformation(uuid)
	logs.Info("Checking "+owlh["name"]+" - "+owlh["ip"]+" Sniffer status")
	  
	status, pid, cpu, mem := GetStatusSnifferSSH(owlh, sshSession)
	cpuStatus := GetStatusCPU(owlh,cpu)
  	memStatus := GetStatusMEM(owlh,mem)
	storageStatus := GetStatusStorage(owlh,ssh)
	
	logs.Error("Checking "+owlh["name"]+" - "+owlh["ip"]+" CPU: "+cpuStatus+" MEM: "+memStatus+" STORAGE: "+storageStatus)
  	if cpuStatus && memStatus && storageStatus {
  	    return "IS RUNNING", true
	}
	return "NOT running", false
}

func GetStatusCPU(owlh map[string]string, cpu string)(status bool){
	var validCPU = regexp.MustCompile(`(\d+)`+cpu+`;`)
	if validCPU == nil{
		return false
	}
	localCPU, _ := strconv.ParseFloat(cpu, 64)
	ddbbCPU, _ := strconv.ParseFloat(owlh["max_cpu"], 64)
	logs.Error("Check CPU for "+owlh["name"]+" - "+owlh["ip"])
	if localCPU>ddbbCPU{
		logs.Error("SNIFFER -> Too much CPU on "+owlh["name"]+" - "+owlh["ip"])
		return false
	}
	return true
}
func GetStatusMEM(owlh map[string]string, mem string)(status bool){
	var validMEM = regexp.MustCompile(`(\d+)`+mem+`;`)
	if validMEM == nil{
		return false
	}
	localMEM, _ := strconv.ParseFloat(mem, 64)
	ddbbMEM, _ := strconv.ParseFloat(owlh["max_mem"], 64)
	logs.Error("Check MEM for "+owlh["name"]+" - "+owlh["ip"])
	if localMEM>ddbbMEM{
		logs.Error("SNIFFER -> Too much MEM on "+owlh["name"]+" - "+owlh["ip"])
		return false
	}
	return true
}	
func GetStatusStorage(owlh map[string]string, sshSession *ssh.Session)(status bool){
	var pcapPath = owlh["pcap_path"]
	logs.Info("SNIFFER -> Too much MEM on "+owlh["name"]+" - "+owlh["ip"])
	status, storage, path := GetStatusStorageSSh(owlh,sshSession, pcapPath)
	if status {
		localStorage, _ := strconv.ParseFloat(storage, 64)
		ddbbStorage, _ := strconv.ParseFloat(owlh["max_storage"], 64)
		if localStorage > ddbbStorage {
			logs.Error("SNIFFER -> Too much STORAGE on "+owlh["name"]+" - "+owlh["ip"])
			return false
		}
		return true
	}
	return false
}	


// func GetStatusStorage(owlh map[string]string, ssh string)(status bool){
// 	pcapPath := conf("pcap_path") //Get pcapPath data from json
// 	logs.Error("Check STORAGE for "+owlh["name"]+" - "+owlh["ip"]+" at pcap path: "+pcapPath)
// 	data, storage, path := GetStatusStorageSSH(owlh, ssh, pcapPath)
// 	if status {
// 		if strconv.ParseInt(storage, 32) > strconv.ParseInt("Conf.Get_Item", 32) {
// 			logs.Error("Too much STORAGE used on "+owlh["name"]+" - "+owlh["ip"]+"  ---  size: "+storage+" path: "+path)
// 			return false
// 		}
// 		return true
// 	}
// 	return false
// }
			

// func RunSniffer(owlh map[string]string, ssh string)(){
// 	logs.Error("Running Sniffer on "+owlh["name"]+" - "+owlh["ip"])
// 	RunSnifferSSH(owlh, ssh, conf("default"), conf("capture_time"), conf("pcap_path"), conf("filter_path"), conf("owlh_user"))
// }

// func StopSniffer(owlh map[string]string, ssh string)(){
// 	logs.Error("Stop Sniffer on "+owlh["name"]+" - "+owlh["ip"])
// 	StopSnifferSSH(owlh, ssh)
// }

// func GetFileList(owlh map[string]string, ssh string)(){
// 	logs.Error("Get file list "+owlh["name"]+" - "+owlh["ip"])
// 	file_list := GetFileListSSH(owlh, ssh, conf("pcap_path"))
// 	sftp := openSftpSSH()
// 	for file := range file_list {
// 		if regexp.MustCompile(`\.pcap`+file+`;`){
// 			OwnerOwlh(owlh, ssh, conf("pcap_path")+file)
// 			TransportFile(owlh, sftp, conf("pcap_path")+file, conf("local_pcap_path")+file)
// 			RemoveFile(owlh, sftp, conf("pcap_path")+file)
// 		}
// 	}
// }

// func OwnerOwlh(owlh map[string]string, ssh string, fileRemote string)(){
// 	logs.Error("Set "+owlh["name"]+" - "+owlh["ip"]+" as owner of file "+conf("owlh_user"))
// 	OwnerOwlhSSH(owlh, ssh, fileRemote, conf("owlh_user"))
// }

// func TransportFile(owlh map[string]string, file string, sftp string, local_path)(){
// 	logs.Error("Get file "+local_path+" from "+owlh["name"]+" - "+owlh["ip"])
// 	TransportFileSSH(owlh, sftp, file, local_path)
// }

// func RemoveFile(owlh map[string]string, sftp string, file string)(){
// 	logs.Error("Remove file "+local_path+" from "+owlh["name"]+" - "+owlh["ip"])
// 	RemoveFileSSH(owlh, ssh, file)
// }







func GetStapServerInformation(uuid string)(serverData map[string]string){
	var param string
	var value string
	stapServer := make(map[string]string)
	logs.Info("CheckOwlhAlive, creating data map for uuid: "+uuid)
	// ip, err := ndb.Sdb.Query("select server_param,server_value from servers where server_param = \"ip\" and server_uniqueid = \""+uuid+"\";")
	ip, err := ndb.Sdb.Query("select server_param,server_value from servers where server_uniqueid = \""+uuid+"\";")
	for ip.Next(){
		if err = ip.Scan(&param, &value); err!=nil {
			logs.Error("Worker Concurrency. Error creating data Map: "+err.Error())
			return nil
		}
		stapServer[param]=value
	}
	return stapServer
}