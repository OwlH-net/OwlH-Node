package stap
import (
    "github.com/astaxie/beego/logs"
    // "os"
    // "os/exec"
    // "strings"
	"regexp"
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


func CheckOwlhAlive(uuid string)(alive bool, ssh string){
    var param string
    var value string
    stapServer := make(map[string]string)
    logs.Info("CheckOwlhAlive, creating data map for uuid: "+uuid)

    ip, err := ndb.Sdb.Query("select server_param,server_value from servers where server_param = \"ip\" and server_uniqueid = \""+uuid+"\";")
    for ip.Next(){
        if err = ip.Scan(&param, &value); err!=nil {
            logs.Error("Worker Concurrency. Error creating data Map: "+err.Error())
            return false,""
        }
		stapServer[param]=value
    } 
    logs.Info("Stap Server Task "+stapServer["name"]+" -- "+stapServer["ip"])
    alive, ssh = owl_connect(stapServer)
    if alive{
		logs.Info("ALIVE Stap Server Task "+stapServer["name"]+" -- "+stapServer["ip"])
        return true, ssh
	}
	logs.Error("NOT ALIVE Stap Server Task "+stapServer["name"]+" -- "+stapServer["ip"])
	return false, ""
}

func GetStatusCPU(owlh map[string]string, cpu string)(status bool){
	var validCPU = regexp.MustCompile(`(\d+)`+cpu+`;`)
	if validCPU == nil{
		return true
	}
	logs.Error("Check CPU for "+owlh["name"]+" - "+owlh["ip"])
	max_cpu := conf("max_cpu")//Get max_cpu data from json
	if strconv.ParseFloat(cpu, 32)>strconv.ParseFloat(max_cpu, 32){
		logs.Error("SNIFFER -> Too much CPU on "+owlh["name"]+" - "+owlh["ip"])
		return false
	}
	return true
}

func GetStatusMEM(owlh map[string]string, mem string)(status bool){
	var validCPU = regexp.MustCompile(`(\d+)`+mem+`;`)
	if validCPU == nil{
		return true
	}
	logs.Error("Check CPU for "+owlh["name"]+" - "+owlh["ip"])
	max_mem :=  conf("max_mem")//Get max_cpu data from json
	if strconv.ParseFloat(cpu, 32)>strconv.ParseFloat(max_mem, 32){
		logs.Error("SNIFFER -> Too much MEM on "+owlh["name"]+" - "+owlh["ip"])
		return false
	}
	return true
}

func GetStatusStorage(owlh map[string]string, ssh string)(status bool){
	pcapPath := conf("pcap_path") //Get pcapPath data from json
	logs.Error("Check STORAGE for "+owlh["name"]+" - "+owlh["ip"]+" at pcap path: "+pcapPath)
	data, storage, path := GetStatusStorageSSH(owlh, ssh, pcapPath)
	if status {
		if strconv.ParseInt(storage, 32) > strconv.ParseInt("Conf.Get_Item", 32) {
			logs.Error("Too much STORAGE used on "+owlh["name"]+" - "+owlh["ip"]+"  ---  size: "+storage+" path: "+path)
			return false
		}
		return true
	}
	return false
}

func GetStatusSniffer(owlh map[string]string, ssh string)(running string, status bool){
	logs.Error("Checking "+owlh["name"]+" - "+owlh["ip"]+" Sniffer status")
	running, pid, cpu, mem := GetStatusSnifferSSH(owlh, ssh)
	cpuStatus := GetStatusCPU(owlh,cpu)
	memStatus := GetStatusMEM(owlh,mem)
	storageStatus := GetStatusStorage(owlh,ssh)
	logs.Error("Checking "+owlh["name"]+" - "+owlh["ip"]+" CPU: "+cpuStatus+" MEM: "+memStatus+" STORAGE: "+storageStatus)
	if cpuStatus && memStatus && storageStatus {
		return running, true
	}
	return running, false
}

func RunSniffer(owlh map[string]string, ssh string)(){
	logs.Error("Running Sniffer on "+owlh["name"]+" - "+owlh["ip"])
	RunSnifferSSH(owlh, ssh, conf("default"), conf("capture_time"), conf("pcap_path"), conf("filter_path"), conf("owlh_user"))
}

func StopSniffer(owlh map[string]string, ssh string)(){
	logs.Error("Stop Sniffer on "+owlh["name"]+" - "+owlh["ip"])
	StopSnifferSSH(owlh, ssh)
}

func GetFileList(owlh map[string]string, ssh string)(){
	logs.Error("Get file list "+owlh["name"]+" - "+owlh["ip"])
	file_list := GetFileListSSH(owlh, ssh, conf("pcap_path"))
	sftp := openSftpSSH()
	for file := range file_list {
		if regexp.MustCompile(`\.pcap`+file+`;`){
			OwnerOwlh(owlh, ssh, conf("pcap_path")+file)
			TransportFile(owlh, sftp, conf("pcap_path")+file, conf("local_pcap_path")+file)
			RemoveFile(owlh, sftp, conf("pcap_path")+file)
		}
	}
}

func OwnerOwlh(owlh map[string]string, ssh string, fileRemote string)(){
	logs.Error("Set "+owlh["name"]+" - "+owlh["ip"]+" as owner of file "+conf("owlh_user"))
	OwnerOwlhSSH(owlh, ssh, fileRemote, conf("owlh_user"))
}
func TransportFile(owlh map[string]string, file string, sftp string, local_path)(){
	logs.Error("Get file "+local_path+" from "+owlh["name"]+" - "+owlh["ip"])
	TransportFileSSH(owlh, sftp, file, local_path)
}
func RemoveFile(owlh map[string]string, sftp string, file string)(){
	logs.Error("Remove file "+local_path+" from "+owlh["name"]+" - "+owlh["ip"])
	RemoveFileSSH(owlh, ssh, file)
}