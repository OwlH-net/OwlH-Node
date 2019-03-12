// package flock

// func RunCMD(cmd, ssh)(data bool, output string){
// 	stdin, stdout, stderr := exec.Command(cmd)
// 	output = ""
// 	for l := range stdout {
// 		output,err := l.strip()
// 		if err != nil {
// 			FlockLogger("Oops! there was a problem stdout range")
// 			return false, ""
// 		}
// 	}
// 	for l := range stderr {
// 		logs.Err("Error RunCMD: "+l.strip())
// 	}
// 	return true, output

// }

// func RunCMD_Bg(cmd, ssh)(data bool, out string){
// 	channel, err := ssh.get_transport().open_session()
// 	if err != nil {
// 		FlockLogger("Oops! there was a problem RunCMD_Bg")
// 		return false, ""
// 	}
// 	FlockLogger("will run command: %s" % cmd,"INFO")
// 	stdin, stdout, stderr, err := channel.exec_command(cmd)
// 	if err != nil {
// 		FlockLogger("Oops! there was a problem RunCMD_Bg command")
// 		return false, ""
// 	}
// 	FlockLogger("will run command DONE: %s" % cmd,"INFO")
// 	output := ""
// 	for l in stdout{
// 		output = l.strip()
// 		print("stdout : %s" % l.strip())
// 	}
		
// 	for l in stderr{
// 		print("stderr : %s" % l.strip())
// 	}
// 	channel.close()
// 	return true, output
// }

// func OwlhConnect(owlh)(data bool, out string){
// 	owlh_user:=GetItem("owlh_user")
//     owlh_key:=GetItem("owlh_user_key")
// 	owlh_ip:=owlh["ip"]

// 	ssh, err := paramiko.SSHClient()
// 	if err != nil {
// 		FlockLogger("Oops! there was a problem SSHClient")
// 		return false, ""
// 	}
// 	err:= ssh.set_missing_host_key_policy( paramiko.AutoAddPolicy() )
// 	if err != nil {
// 		FlockLogger("Oops! there was a problem set_missing_host_key_policy")
// 		return false, ""
// 	}
// 	err:= ssh.connect(owlh_ip, owlh_user, owlh_key)
// 	if err != nil {
// 		FlockLogger("Oops! there was a problem  ssh.connect")
// 		return false, ""
// 	}
// }


// func GetStatusSnifferSSH(owlh, ssh)(running bool, pid string, cpu string, mem string){
// 	FlockLogger("Check sniffer in "+owlh["name"]+" ip "+owlh["ip"])
// 	cmd:="top -b -n1 | grep -v sudo | grep tcpdump | awk \'{print $1 "," $9 "," $10}\'"
//     output := ""
// 	status, output := run_cmd(cmd, ssh)
// 	if _, err := os.Stat("'\'d+,'\'d+'\'.'\'d+,'\'d+'\'.'\'d+",str(output)); err == nil{
//         pid, cpu, mem := output.split(",")
//         if pid != nil {
//         	flogger("sniffer is working in owlh %s (%s) with pid %s, CPU %s, MEM %s",owlh["name"], owlh["ip"],pid, cpu, mem)
// 			return true, pid, cpu, mem
// 		}
//     flogger("sniffer NOT working in owlh" +owlh["name"]+" - "+owlh["ip"])
// 	return false, "", "", ""
// }

// func CheckOwlhAlive(owlh)(data bool, ssh string){
//     FlockLogger("check if owl %s (%s) is alive" % (owlh["name"], owlh["ip"]))
//     alive, ssh = owl_connect(owlh)
//     if alive{
// 		FlockLogger("owlh is alive" + owlh["name"]+" - "+ owlh["ip"])
//         return true, ssh
// 	}
// 	FlockLogger("owlh is NOT alive" + owlh["name"]+" - "+ owlh["ip"])
// 	return false, ""
// }

// func GetStatusStorage(owlh, ssh, folder)(data bool, gr1 string, gr2 string){
//     FlockLogger("Check status " + owlh["name"]+" - "+ owlh["ip"])
//     cmd:="df -h %s --output=source,pcent | grep -v Filesystem | awk \'{print $1","$2}\'' % folder"
//     output := ""
//     status, output := run_cmd(cmd, ssh)
//     if re.search("[^,]+,\'d+%",str(output)){
// 		regx := re.match("([^,]+),(\d+)%",output)
// 		FlockLogger("PCAP storage used in " + owlh["name"]+" - "+ owlh["ip"])    
//         return true, regx.group(2), regx.group(1)
// 	}    
// 	FlockLogger("Can't fix " + owlh["name"]+" - "+ owlh["ip"])    
// 	return false, "", ""
// }

// func RunSniffer(owlh, ssh, interface, capture, pcapPath, filterPath, user)(){
// 	FlockLogger("Starting traffic collector " + owlh["name"]+" - "+ owlh["ip"])
//     cmd := "nohup sudo tcpdump -i %s -G %s -w %s`hostname`-%%y%%m%%d%%H%%M%%S.pcap -F %s -z %s >/dev/null 2>&1 &' % (interface, capture, pcap_path, filter_path, user)"
//     output := ""
// 	RunCMD_Bg(cmd, ssh)
// }

// func StopSniffer(owlh, ssh)(){
// 	cmd := "ps -ef | grep -v grep | grep tcpdump | awk '{printf $2\",\"}'"
//     output := ""
//     status, output = RunCMD(cmd,ssh)
// 	if regexp.MustCompile(`[^,]+,`+output+`,`){
// 		output := re.sub(","," ",output)
//         FlockLogger("Stopping traffic collector " + owlh["name"]+" - "+ owlh["ip"])
//         cmd := "sudo kill -9 " + output
//         status, output := RunCMD(cmd,ssh)
//         return status, output
// 	}
// 	FlockLogger("Stopping traffic collector => Nothing to stop... " + owlh["name"]+" - "+ owlh["ip"])
// }

// func GetFileList(owlh, ssh, folder)(){
//     cmd := "find %s*.pcap -maxdepth 0 -type f -mmin +1|sed 's#.*/##'| awk '{printf $1\",\"}'"+folder
//     status, output := RunCMD(cmd,ssh)
//     files := output.split(",")
// 	return files
// }

// func OwnerOwlhSSH(owlh, ssh, fileRemote, User)(){
// 	if  _, err := os.Stat("'\'.pcap",fileRemote); err == nil{
// 		FlockLogger("Settings as owner " + owlh["name"]+" - "+ owlh["ip"])
//         cmd := "sudo chown "+user+" "+fileRemote
// 		status, output := RunCMD(cmd,ssh)
// 	}
// }

// func TransportFileSSH(owlh, sftp, fileRemote, fileLocal)(){
//     if  _, err := os.Stat("'\'.pcap",fileRemote); err == nil{
// 		FlockLogger("Collecting " fileRemote" - "+ fileLocal)
// 		sftp.get(file_remote, file_local)
// 	}
// }

// func RemoveFileSSH(owlh, sftp, fileRemote)(){
//     if  _, err := os.Stat("'\'.pcap",fileRemote); err == nil{
// 		FlockLogger("Cleaning " + owlh["name"]+" - "+ fileLocal)
// 		sftp.get(file_remote, file_local)
// 	}
// }

// func Nothing()(){
//     scp_opt=""
//     cmd="scp -q " + scp_opt + " -o NumberOfPasswordPrompts=1 -o StrictHostKeyChecking=no "+test_script+" root@"+priv_ip+":~/; echo $? done."
//     logs.Info("\n test 2\n cmd "+cmd+"\n")
//     RunCMD(cmd)

//     scp_opt="-v"
//     cmd="scp -q " + scp_opt + " -o NumberOfPasswordPrompts=1 -o StrictHostKeyChecking=no "+test_script+" root@"+priv_ip+":~/; echo $? done."
//     logs.Info("\n test 3\n cmd "+cmd+"\n")
// 	RunCMD(cmd)
// }