// package flock
// import (
// 	"encoding/json"
// 	"github.com/astaxie/beego/logs"
// 	"regexp"
// )


// func CheckOwlhAlive(owlh)(alive, ssh){
// 	FlockLogger("check owlh "+owlh["name"]+" is alive")
// 	alive, ssh:= CheckOwlhAlive(owlh)
// 	return alive, ssh
// }

// func GetStatusCPU(owlh, cpu)(data bool){
// 	if !regexp.MustCompile(`(\d+)`+cpu+`;`){
// 		return true
// 	}
// 	FlockLogger("check owl %s CPU",owlh["name"])
// 	max_cpu := GetItem("max_cpu")
// 	if cpu > max_cpu{
// 		FlockLogger("SNIFFER -> too much cpu on owlh "+owlh["name"]+" WARNING")
// 		return false
// 	}
// 	return true
// }

// func GetStatusMem(owlh, mem)(data bool){
// 	if !regexp.MustCompile(`(\d+)`+mem+`;`) {
// 		return true
// 	}
// 	FlockLogger("check owlh "+owlh["name"]+" MEM")
// 	max_mem := GetItem("max_mem")
// 	if mem > max_mem{
// 		FlockLogger("SNIFFER -> too much mem on owlh "+owlh["name"]+" WARNING")
// 		return false
// 	}
// 	return true
// }

// func GetStatusSniffer(owlh, ssh)(running string, data bool){
// 	FlockLogger("Check "+owlh["name"]+" sniffer status")
// 	running, pid, cpu, mem := GetStatusSnifferSSH(owlh, ssh)
// 	cpustatus := GetStatusCPU(owlh, cpu)
// 	memstatus := GetStatusMem(owlh, mem)
// 	storagestatus := GetStatusStorage(owlh, ssh)
// 	FlockLogger("Check "+owlh["name"]+" CPU: "+cpustatus+" MEM: "+memstatus+" STORAGE: "+storagestatus+" sniffer status")
// 	if cpustatus && memstatus && storagestatus {
// 		return running, true
// 	}
// 	return running, false
// }

// func GetStatusStorage(owlh, ssh)(){
// 	pcap_path = GetItem("pcap_path")
// 	FlockLogger("Check "+owlh["name"]+" storage -> pcap path")
// 	status, storage, path := GetStatusStorage(owlh, ssh, pcap_path)
//     if status{
// 		if storage > GetItem("max_storage"){
// 			FlockLogger("PCAP -> too much storage used "+owlh["name"]+" WARNING")
// 			return false
// 		}
// 		return true
// 	}
// 	return false
// }

// func RunSniffer(owlh,ssh)(){
// 	FlockLogger("Run sniffer "+owlh["name"])
// 	RunSniffer(owlh, ssh, GetItem("default_interface"),GetItem("capture_time"),GetItem("pcap_path"),GetItem("filter_path"),GetItem("owlh_user"))
// }

// func StopSniffer(owlh, ssh)(){
// 	FlockLogger("Stop sniffer "+owlh["name"])
// 	StopSniffer(owlh, ssh)
// }

// func GetFileList(owlh, ssh)(){
// 	FlockLogger("Get file list "+owlh["name"])
// 	file_list := GetFileList(owlh, ssh, GetItem("pcap_path"))
// 	sftp := ssh.open_sftp()
// 	for file := range file_list {
// 		if _, err := os.Stat("'\'.pcap",file); err == nil{
// 			OwnerOwlhMonitor(owlh, ssh, GetItem("pcap_path")+file)
// 		}
// 	}
// }

// func OwnerOwlhMonitor(owlh, ssh, fileRemote)(){
// 	FlockLogger("Set "+owlh["name"]+" as owner of file "+fileRemote)
// 	OwnerOwlhSSH(owlh, ssh, GetItem("pcap_path")+file)
// }

// func TransportFileMonitor(owlh, sftp, file, localPath)(){
// 	FlockLogger("Get "+owlh["name"]+" file")
// 	TransportFileSSH(owlh, ssh, file, localPath)
// }

// func RemoveFileMonitor(owlh, sftp, file)(){
// 	FlockLogger("Remove "+owlh["name"]+" file")
// 	RemoveFileSSH(owlh, sftp, file)
// }