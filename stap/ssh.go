package stap

import (
    "github.com/astaxie/beego/logs"
    "github.com/pkg/sftp"
    "golang.org/x/crypto/ssh"
    "io/ioutil"
    "os"
    "owlhnode/database"
    "owlhnode/utils"
    "regexp"
    "strconv"
    "strings"
)

//Run a command over the active ssh session
func RunCMD(uuid string, cmd string) (status bool, data string) {
    logs.Info("Connect ssh from RunCMD")
    var CMDoutput []byte

    //check if server is alive
    alive, sshValue := owlh_connect(uuid)
    if !alive {
        return false, ""
    }
    logs.Info("Command to exec --> " + cmd)
    CMDoutput, err := sshValue.Output(cmd)
    sshValue.Close()
    if err != nil {
        logs.Error("Error executing SSH commands: " + err.Error())
        return false, ""
    } else {
        logs.Info("Command output --> " + string(CMDoutput))
        logs.Info("No error executing command through SSH connection")
        return true, string(CMDoutput)
    }
}

//connect to remote server throught ssh
func owlh_connect(uuid string) (alive bool, sshValue *ssh.Session) {
    var err error
    userSSH, err := utils.GetKeyValueString("stapPubKey", "user")
    if err != nil {
        logs.Error("owlh_connect Error getting data from main.conf")
    }
    cert, err := utils.GetKeyValueString("stapPubKey", "cert")
    if err != nil {
        logs.Error("owlh_connect Error getting data from main.conf")
    }

    owlh, err := ndb.GetStapServerInformation(uuid)
    if err != nil {
        logs.Error("Error retrieving stap server information")
    }
    logs.Info("Name: " + owlh["name"] + " IP: " + owlh["ip"])

    // //Declare ssh config
    sshConfig := &ssh.ClientConfig{
        User: userSSH,
        Auth: []ssh.AuthMethod{
            PublicKeyFile(cert),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    }
    logs.Warn("SSH Config declared!! -->" + owlh["ip"])
    client, err := ssh.Dial("tcp", owlh["ip"]+":22", sshConfig)
    if err != nil {
        logs.Error("SSH Dial Error: " + err.Error())
        return false, nil
    }

    logs.Warn("SSH starting session")

    //Launch new session
    session, err := client.NewSession()
    if err != nil {
        logs.Error("New Session Error: " + err.Error())
        session.Close()
        return false, nil
    }
    logs.Info("New session has been established")
    return true, session //return session
}

//read public key for ssh configuration
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

//Get sniffer status of CPU, MEM and PID
func GetStatusSnifferSSH(uuid string) (status bool, pid string, cpu string, mem string) {
    cmd := "top -b -n1 | grep -v sudo | grep tcpdump | awk '{print $1 \",\" $9 \",\" $10}'"
    output := ""
    status, output = RunCMD(uuid, cmd)
    logs.Info("OUTPUT FOR REGEXP: " + output)
    var validOutput = regexp.MustCompile(`\d+,\d+\.\d+,\d+\.\d+`)
    if validOutput.MatchString(output) {
        logs.Info("Inside Validation REGEXP")
        splitValue := strings.Split(output, ",")
        pid = splitValue[0]
        cpu = splitValue[1]
        mem = splitValue[2]
        if pid != "" {
            logs.Info("Sniffer PID: " + pid + " is working")
            return true, pid, cpu, mem
        }
    }
    logs.Info("Sniffer is NOT working for " + uuid)
    return false, "", "", ""
}

//Get status, path and storage available of the remote server
func GetStatusStorageSSh(uuid string, folder string) (status bool, path string, storage string) {
    logs.Info("Checking if storage is OK in " + uuid)
    cmd := "df -h " + folder + " --output=source,pcent | grep -v Filesystem | awk '{print $1 \",\" $2}'"
    output := ""
    status, output = RunCMD(uuid, cmd)
    logs.Info("Output Run CMD -->" + output)
    var validOutput = regexp.MustCompile(`[^,]+,\d+%`)
    if validOutput.MatchString(output) {
        splitValue := strings.Split(output, ",")
        return true, splitValue[0], splitValue[1]
    }
    return false, "", ""
}

//launch tcpdump when server is launched
func RunSnifferSSH(uuid string) {
    logs.Info("Launching Sniffer...")
    owlh, err := ndb.GetStapServerInformation(uuid)
    if err != nil {
        logs.Error("Error retrieving stap server information")
    }
    cmd := "nohup sudo tcpdump -i " + owlh["default_interface"] + " -G " + owlh["capture_time"] + " -w " + owlh["pcap_path"] + "`hostname`-%y%m%d%H%M%S.pcap -F " + owlh["filter_path"] + " -z " + owlh["owlh_user"] + " >/dev/null 2>&1 &"
    status, output := RunCMD(uuid, cmd)
    if status {
        logs.Info("Sniffer is Running for server " + owlh["name"] + " - " + owlh["ip"] + ". Output--> " + output)
    }
}

//kill PID of tcpdump when server is stopped
func StopSnifferSSH(uuid string) {
    logs.Info(uuid + "Stopping Sniffer...")
    cmd := "ps -ef | grep -v grep | grep tcpdump | awk '{print $2}'"
    output := ""
    _, output = RunCMD(uuid, cmd)
    logs.Info("NUMBER OF PID -->" + output)
    splitValue := strings.Split(output, "\n")
    for value := range splitValue {
        var validOutput = regexp.MustCompile(`^\d+`)
        if validOutput.MatchString(splitValue[value]) {
            cmd := "sudo kill -9 " + splitValue[value]
            _, _ = RunCMD(uuid, cmd)
            logs.Info("Sniffer has been stopped succesfully")
        } else {
            logs.Error("Sniffer is still running...")
        }
    }
}

//get all pcap files older than a minute
func GetFileListSSH(uuid string, owlh map[string]string, path string) (list []string) {
    cmd := "find " + owlh["pcap_path"] + "*.pcap -maxdepth 0 -type f -mmin +1|sed 's#.*/##'| awk '{printf $1 \",\"}'"
    output := ""
    _, output = RunCMD(uuid, cmd)
    splitValue := strings.Split(output, ",")
    return splitValue
}

//change owner of remote files
func OwnerOwlhSSH(uuid string, owlh map[string]string, fileRemote string) {
    logs.Info(uuid + " Change remote file owner: " + fileRemote)
    var validOutput = regexp.MustCompile(`\.pcap+`)
    if validOutput.MatchString(fileRemote) {
        cmd := "sudo chown " + owlh["owlh_user"] + " " + owlh["pcap_path"] + fileRemote
        logs.Alert("Command for change owner file--> " + cmd)
        RunCMD(uuid, cmd)
    }
}

//use  sftp to transfer files from remote to local machine
func TransportFileSSH(uuid string, owlh map[string]string, file string) {
    logs.Info(uuid + " Transport remote file: " + file + " to local machine")
    var validOutput = regexp.MustCompile(`\.pcap+`)
    if validOutput.MatchString(file) {
        status := SftpCMD(uuid, owlh["pcap_path"]+file, owlh["local_pcap_path"]+file)
        logs.Info("Output CMD TransportFileSSH: " + file + " - " + strconv.FormatBool(status))
    }
}

//When a files is copied from remote to local machine, the remote file is removed
func RemoveFileSSH(uuid string, owlh map[string]string, file string) {
    logs.Info("Remove remote files " + uuid)
    var validOutput = regexp.MustCompile(`\.pcap+`)
    if validOutput.MatchString(file) {
        cmd := "sudo rm " + file + ""
        _, output := RunCMD(uuid, cmd)
        logs.Info("Output CMD RemoveFileSSH: " + output)
    }
}

//Copy remote file by sftp and remove them
func SftpCMD(uuid string, srcFile string, dstFile string) (status bool) {
    logs.Info("Connect sftp SftpCMD")
    alive, sshClient := owlh_connect_client(uuid)
    if !alive {
        return false
    }

    logs.Info("NewClient")
    sftp, err := sftp.NewClient(sshClient)
    if err != nil {
        logs.Error(err)
        return false
    }
    defer sftp.Close()

    // Open the source file
    logs.Info("Open Remote File")
    remoteFile, err := sftp.Open(srcFile)
    if err != nil {
        logs.Error(err)
        return false
    }
    defer remoteFile.Close()

    // Create the destination file
    logs.Info("Create DST file")
    localFile, err := os.Create(dstFile)
    if err != nil {
        logs.Error(err)
        return false
    }
    defer localFile.Close()

    // Copy the file
    logs.Info("Copy file")
    remoteFile.WriteTo(localFile)

    //remove remote file
    logs.Info("Delete remote file")
    err = sftp.Remove(srcFile)
    if err != nil {
        logs.Error(err)
        return false
    }
    return true
}

//get certs from main.conf and create a ssh connection for make sftp
func owlh_connect_client(uuid string) (alive bool, sshClient *ssh.Client) {
    var err error
    userSSH, err := utils.GetKeyValueString("stapPubKey", "user")
    if err != nil {
        logs.Error("owlh_connect_client Error getting data from main.conf")
    }
    cert, err := utils.GetKeyValueString("stapPubKey", "cert")
    if err != nil {
        logs.Error("owlh_connect_client Error getting data from main.conf")
    }

    owlh, err := ndb.GetStapServerInformation(uuid)
    if err != nil {
        logs.Error("Error retrieving stap server information")
    }
    // //Declare ssh config
    sshConfig := &ssh.ClientConfig{
        User: userSSH,
        Auth: []ssh.AuthMethod{
            PublicKeyFile(cert),
        },
        HostKeyCallback: ssh.InsecureIgnoreHostKey(),
    }
    client, err := ssh.Dial("tcp", owlh["ip"]+":22", sshConfig)
    if err != nil {
        logs.Error("SSH Dial Error: " + err.Error())
        return false, nil
    }
    return true, client
}
