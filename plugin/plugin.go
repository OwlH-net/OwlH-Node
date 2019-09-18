package plugin

import (
    "github.com/astaxie/beego/logs"
	"owlhnode/database"
	"owlhnode/zeek"
    // "owlhnode/suricata"
    "os/exec"
    "bytes"
    "errors"
    "os"
    "strconv"
    "strings"
    "io/ioutil"
	"owlhnode/utils"
)

func ChangeServiceStatus(anode map[string]string)(err error) {
    allPlugins,err := ndb.GetPlugins()
    if anode["type"] == "suricata"{
        if anode["status"] == "enabled"{
            for x := range allPlugins {
                //get all db values and check if any
                if allPlugins[x]["pid"] != "none" && allPlugins[x]["interface"] == anode["interface"] && allPlugins[x]["status"] == "enabled" && x != anode["server"]{
                    logs.Error("Can't launch more than one suricata with same interface. Please, select other interface.")
                    return errors.New("Can't launch more than one suricata with same interface. Please, select other interface.")
                }
            }
            err = LaunchSuricataService(anode["server"], anode["interface"])
            if err != nil {logs.Error("LaunchSuricataService status Error: "+err.Error()); return err}
        }else if anode["status"] == "disabled"{
            err = StopSuricataService(anode["server"], anode["status"])
            if err != nil {logs.Error("StopSuricataService status Error: "+err.Error()); return err}

        }
    } else if anode["type"] == "zeek"{
        mainConfData, err := ndb.GetMainconfData()
        if (mainConfData["zeek"]["status"] == "disabled"){ return nil }        
        if anode["status"] == "enabled"{
            err = zeek.DeployZeek()
            if err != nil {logs.Error("plugin/ChangeServiceStatus error deploying zeek: "+err.Error()); return err}

            err = ndb.UpdatePluginValue(anode["server"],"previousStatus","none")
            if err != nil {logs.Error("plugin/ChangeServiceStatus error updating zeek previousStatus to none: "+err.Error()); return err}

            err = ndb.UpdatePluginValue(anode["server"],"status","enabled")
            if err != nil {logs.Error("plugin/ChangeServiceStatus error updating zeek status to enabled: "+err.Error()); return err}
        } else if anode["status"] == "disabled"{
            data, err := zeek.StopZeek(); logs.Error(data)
            if err != nil {logs.Error("plugin/ChangeServiceStatus error deploying zeek: "+err.Error()); return err}

            err = ndb.UpdatePluginValue(anode["server"],"previousStatus",anode["status"])
            if err != nil {logs.Error("plugin/ChangeServiceStatus error updating zeek previousStatus to status: "+err.Error()); return err}

            err = ndb.UpdatePluginValue(anode["server"],"status","disabled")
            if err != nil {logs.Error("plugin/ChangeServiceStatus error updating zeek status to disabled: "+err.Error()); return err}
        }

    }
    return err
}

func ChangeMainServiceStatus(anode map[string]string)(err error) {
    err = ndb.UpdateMainconfValue(anode["service"],anode["param"],anode["status"])
    if err != nil {logs.Error("plugin/ChangeMainServiceStatus error: "+err.Error()); return err}

    allPlugins,err := ndb.GetPlugins()
    if anode["service"] == "suricata" {
        for x := range allPlugins {
            if anode["status"] == "disabled"{
                if allPlugins[x]["status"] == "enabled" && allPlugins[x]["type"] == "suricata"{
                    err = StopSuricataService(x, allPlugins[x]["status"])
                    if err != nil {logs.Error("StopSuricataService status Error: "+err.Error()); return err}
                }
            }else if anode["status"] == "enabled"{
                if allPlugins[x]["previousStatus"] == "enabled" && allPlugins[x]["type"] == "suricata"{
                    err = LaunchSuricataService(x, allPlugins[x]["interface"])
                    if err != nil {logs.Error("LaunchSuricataService status Error: "+err.Error()); return err}
                }
            }
        }
    } else if anode["service"] == "zeek" {
        for x := range allPlugins {
            if anode["status"] == "disabled"{
                if allPlugins[x]["status"] == "enabled" && allPlugins[x]["type"] == "zeek"{
                    err = ndb.UpdatePluginValue(x,"previousStatus","enabled")
                    if err != nil {logs.Error("plugin/ChangeMainServiceStatus error updating pid at DB: "+err.Error()); return err}

                    err = ndb.UpdatePluginValue(x,"status","disabled")
                    if err != nil {logs.Error("plugin/ChangeMainServiceStatus error updating pid at DB: "+err.Error()); return err}
                    
                    data, err := zeek.StopZeek(); logs.Error(data)
                    if err != nil {logs.Error("plugin/ChangeMainServiceStatus error deploying zeek: "+err.Error()); return err}
                }
            }else if anode["status"] == "enabled"{
                if allPlugins[x]["previousStatus"] == "enabled" && allPlugins[x]["type"] == "zeek"{
                    err = ndb.UpdatePluginValue(x,"previousStatus","none")
                    if err != nil {logs.Error("plugin/ChangeMainServiceStatus error updating pid at DB: "+err.Error()); return err}

                    err = ndb.UpdatePluginValue(x,"status","enabled")
                    if err != nil {logs.Error("plugin/ChangeMainServiceStatus error updating pid at DB: "+err.Error()); return err}

                    err = zeek.DeployZeek()
                    if err != nil {logs.Error("plugin/ChangeMainServiceStatus error deploying zeek: "+err.Error()); return err}
                }
            }
        }
    }

    return err
}

func DeleteService(anode map[string]string)(err error) {
	err = ndb.DeleteService(anode["server"])
    if err != nil {logs.Error("plugin/DeleteService error: "+err.Error()); return err}

    if _, err := os.Stat("/etc/suricata/bpf/"+anode["server"]+" - filter.bpf"); !os.IsNotExist(err) {
        err = os.Remove("/etc/suricata/bpf/"+anode["server"]+" - filter.bpf")
        if err != nil {logs.Error("plugin/SaveSuricataInterface error deleting a pid file: "+err.Error())}
    }

    return err
}

func AddPluginService(anode map[string]string) (err error) {
    uuid := utils.Generate()    
    if anode["type"] == "socket-network"{
        err = ndb.InsertPluginService(uuid, "node", anode["uuid"]); if err != nil {logs.Error("InsertPluginService node Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "name", anode["name"]); if err != nil {logs.Error("InsertPluginService name Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "type", anode["type"]); if err != nil {logs.Error("InsertPluginService type Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "interface", anode["interface"]); if err != nil {logs.Error("InsertPluginService interface Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "port", anode["port"]); if err != nil {logs.Error("InsertPluginService port Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "cert", anode["cert"]); if err != nil {logs.Error("InsertPluginService certtificate Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "pid", "none"); if err != nil {logs.Error("InsertPluginService pid Error: "+err.Error()); return err}
    }
    if anode["type"] == "socket-pcap"{
        err = ndb.InsertPluginService(uuid, "node", anode["uuid"]); if err != nil {logs.Error("InsertPluginService node Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "name", anode["name"]); if err != nil {logs.Error("InsertPluginService name Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "type", anode["type"]); if err != nil {logs.Error("InsertPluginService type Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "interface", anode["interface"]); if err != nil {logs.Error("InsertPluginService interface Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "port", anode["port"]); if err != nil {logs.Error("InsertPluginService port Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "cert", anode["cert"]); if err != nil {logs.Error("InsertPluginService certtificate Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "pcap-path", anode["pcap-path"]); if err != nil {logs.Error("InsertPluginService pcap-path Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "pcap-prefix", anode["pcap-prefix"]); if err != nil {logs.Error("InsertPluginService pcap-prefix Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "bpf", anode["bpf"]); if err != nil {logs.Error("InsertPluginService bpf Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "pid", "none"); if err != nil {logs.Error("InsertPluginService pid Error: "+err.Error()); return err}
    }
    if anode["type"] == "network-socket"{
        err = ndb.InsertPluginService(uuid, "node", anode["uuid"]); if err != nil {logs.Error("InsertPluginService node Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "name", anode["name"]); if err != nil {logs.Error("InsertPluginService name Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "type", anode["type"]); if err != nil {logs.Error("InsertPluginService type Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "interface", anode["interface"]); if err != nil {logs.Error("InsertPluginService interface Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "port", anode["port"]); if err != nil {logs.Error("InsertPluginService port Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "cert", anode["cert"]); if err != nil {logs.Error("InsertPluginService certtificate Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "collector", anode["collector"]); if err != nil {logs.Error("InsertPluginService collector Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "bpf", anode["bpf"]); if err != nil {logs.Error("InsertPluginService bpf Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "pid", "none"); if err != nil {logs.Error("InsertPluginService pid Error: "+err.Error()); return err}
    }
    if anode["type"] == "zeek"{
        allPlugins,err := ndb.GetPlugins()
        for x := range allPlugins{
            if allPlugins[x]["type"] == "zeek"{ return nil}
        }
        err = ndb.InsertPluginService(uuid, "node", anode["uuid"]); if err != nil {logs.Error("InsertPluginService node Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "name", anode["name"]); if err != nil {logs.Error("InsertPluginService name Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "type", anode["type"]); if err != nil {logs.Error("InsertPluginService type Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "status", "disabled"); if err != nil {logs.Error("InsertPluginService status Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "previousStatus", "none"); if err != nil {logs.Error("InsertPluginService previousStatus Error: "+err.Error()); return err}
    }
    if anode["type"] == "suricata"{
        err = ndb.InsertPluginService(uuid, "node", anode["uuid"]); if err != nil {logs.Error("InsertPluginService node Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "name", anode["name"]); if err != nil {logs.Error("InsertPluginService name Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "type", anode["type"]); if err != nil {logs.Error("InsertPluginService type Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "status", "disabled"); if err != nil {logs.Error("InsertPluginService status Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "previousStatus", "none"); if err != nil {logs.Error("InsertPluginService previousStatus Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "interface", ""); if err != nil {logs.Error("InsertPluginService interface Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "bpf", ""); if err != nil {logs.Error("InsertPluginService bpf Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "ruleset", ""); if err != nil {logs.Error("InsertPluginService ruleset Error: "+err.Error()); return err}
        err = ndb.InsertPluginService(uuid, "pid", "none"); if err != nil {logs.Error("InsertPluginService ruleset Error: "+err.Error()); return err}
    }

    return nil
}

func SaveSuricataInterface(anode map[string]string)(err error) {
    err = ndb.UpdatePluginValue(anode["service"],anode["param"],anode["interface"])
    if err != nil {logs.Error("plugin/SaveSuricataInterface error: "+err.Error()); return err}
    return err
}

func LaunchSuricataService(uuid string, iface string)(err error){

    mainConfData, err := ndb.GetMainconfData()
    if (mainConfData["suricata"]["status"] == "disabled"){ return nil }

    cmd := exec.Command("suricata", "-D", "-c", "/etc/suricata/suricata.yaml", "-i", iface, "-F", "/etc/suricata/bpf/"+uuid+" - filter.bpf" ,"--pidfile", "/var/run/suricata/"+uuid+"-pidfile.pid")
    var stdBuffer bytes.Buffer
    cmd.Stderr = &stdBuffer

    err = cmd.Run()
    if err != nil {
        logs.Error(stdBuffer.String())
        logs.Error("plugin/LaunchSuricataService error launching Suricata: "+err.Error());
        //delete pid file
        err = os.Remove("/var/run/suricata/"+uuid+"-pidfile.pid")
        if err != nil {logs.Error("plugin/SaveSuricataInterface error deleting a pid file: "+err.Error())}
    }else{
        //read file
        currentpid, err := os.Open("/var/run/suricata/"+uuid+"-pidfile.pid")
        if err != nil {logs.Error("plugin/LaunchSuricataService error openning Suricata: "+err.Error()); return err}
        defer currentpid.Close()
        pid, err := ioutil.ReadAll(currentpid)

        //save pid to db
        err = ndb.UpdatePluginValue(uuid,"pid",string(pid))
        if err != nil {logs.Error("plugin/SaveSuricataInterface error updating pid at DB: "+err.Error()); return err}

        //change DB status
        err = ndb.UpdatePluginValue(uuid,"previousStatus","none")
        if err != nil {logs.Error("plugin/LaunchSuricataService error: "+err.Error()); return err}

        //change DB status
        err = ndb.UpdatePluginValue(uuid,"status","enabled")
        if err != nil {logs.Error("plugin/LaunchSuricataService error: "+err.Error()); return err}
    }
    return nil
}

func ModifyStapValues(anode map[string]string)(err error) {
    if anode["type"] == "zeek"{
        err = ndb.UpdatePluginValue(anode["service"],"name",anode["name"]); if err != nil {logs.Error("ModifyStapValues zeek Error: "+err.Error()); return err}
        err = zeek.DeployZeek()
        if err != nil {logs.Error("plugin/ModifyStapValues error deploying zeek: "+err.Error()); return err}
    }else if anode["type"] == "socket-network"{
        err = ndb.UpdatePluginValue(anode["service"],"name",anode["name"]); if err != nil {logs.Error("ModifyStapValues socket-network Error: "+err.Error()); return err}
        err = ndb.UpdatePluginValue(anode["service"],"port",anode["port"]) ; if err != nil {logs.Error("ModifyStapValues socket-network Error: "+err.Error()); return err}
        err = ndb.UpdatePluginValue(anode["service"],"cert",anode["cert"]) ; if err != nil {logs.Error("ModifyStapValues socket-network Error: "+err.Error()); return err}
    }else if anode["type"] == "socket-pcap"{
        err = ndb.UpdatePluginValue(anode["service"],"name",anode["name"]) ; if err != nil {logs.Error("ModifyStapValues socket-pcap Error: "+err.Error()); return err}
        err = ndb.UpdatePluginValue(anode["service"],"port",anode["port"]) ; if err != nil {logs.Error("ModifyStapValues socket-pcap Error: "+err.Error()); return err}
        err = ndb.UpdatePluginValue(anode["service"],"cert",anode["cert"]) ; if err != nil {logs.Error("ModifyStapValues socket-pcap Error: "+err.Error()); return err}
        err = ndb.UpdatePluginValue(anode["service"],"pcap-path",anode["pcap-path"]) ; if err != nil {logs.Error("ModifyStapValues socket-pcap Error: "+err.Error()); return err}
        err = ndb.UpdatePluginValue(anode["service"],"pcap-prefix",anode["pcap-prefix"]) ; if err != nil {logs.Error("ModifyStapValues socket-pcap Error: "+err.Error()); return err}
    }else if anode["type"] == "network-socket"{
        err = ndb.UpdatePluginValue(anode["service"],"name",anode["name"]) ; if err != nil {logs.Error("ModifyStapValues network-socket Error: "+err.Error()); return err}
        err = ndb.UpdatePluginValue(anode["service"],"port",anode["port"]) ; if err != nil {logs.Error("ModifyStapValues network-socket Error: "+err.Error()); return err}
        err = ndb.UpdatePluginValue(anode["service"],"cert",anode["cert"])  ; if err != nil {logs.Error("ModifyStapValues network-socket Error: "+err.Error()); return err}
        err = ndb.UpdatePluginValue(anode["service"],"collector",anode["collector"]) ; if err != nil {logs.Error("ModifyStapValues network-socket Error: "+err.Error()); return err}
    }
    return nil
}

func StopSuricataService(uuid string, status string)(err error){
    //pid
    // currentpid, err := os.Open("/var/run/suricata/"+uuid+"-pidfile.pid")
    // if err != nil {logs.Error("plugin/ChangeServiceStatus error reading Suricata pid: "+err.Error()); return err}
    // defer currentpid.Close()
    // pid, err := ioutil.ReadAll(currentpid)
    allPlugin,err := ndb.GetPlugins()

    //kill suricata process
    PidInt,_ := strconv.Atoi(strings.Trim(string(allPlugin[uuid]["pid"]), "\n"))
    process, _ := os.FindProcess(PidInt)
    _ = process.Kill()
    // if err != nil {logs.Error("plugin/StopSuricataService error killing Suricata process: "+err.Error()); return err}

    //delete pid file
    _ = os.Remove("/var/run/suricata/"+uuid+"-pidfile.pid")
    // if err != nil {logs.Error("plugin/SaveSuricataInterface error deleting a pid file: "+err.Error()); return err}

    //change DB pid
    err = ndb.UpdatePluginValue(uuid,"pid","none")
    if err != nil {logs.Error("plugin/SaveSuricataInterface error updating pid at DB: "+err.Error()); return err}

    //change DB status
    err = ndb.UpdatePluginValue(uuid,"previousStatus",status)
    if err != nil {logs.Error("plugin/StopSuricataService error: "+err.Error()); return err}

    //change DB status
    err = ndb.UpdatePluginValue(uuid,"status","disabled")
    if err != nil {logs.Error("plugin/StopSuricataService error: "+err.Error()); return err}

    return nil
}

func DeployStapService(anode map[string]string)(err error) { 
    allPlugins,err := ndb.GetPlugins()
    
    if anode["type"] == "socket-network" {
        pid, err := exec.Command("bash","-c","ps -ef | grep socat | grep OPENSSL-LISTEN:"+allPlugins[anode["service"]]["port"]+" | grep -v bash | awk '{print $2}'").Output()
        if err != nil {logs.Error("DeployStapService deploy socket-network Error: "+err.Error()); return err}
        pidValue := strings.Split(string(pid), "\n")
        if pidValue[0] != "" {
            return nil
        }

        cmd := exec.Command("bash","-c","/usr/bin/socat -d OPENSSL-LISTEN:"+allPlugins[anode["service"]]["port"]+",reuseaddr,pf=ip4,fork,cert="+allPlugins[anode["service"]]["cert"]+",verify=0 SYSTEM:\"tcpreplay -t -i "+allPlugins[anode["service"]]["interface"]+" -\" &")
        var errores bytes.Buffer
        cmd.Stdout = &errores
        err = cmd.Start()
        if err != nil {logs.Error("DeployStapService deploying Error: "+err.Error()); return err}        


        pid, err = exec.Command("bash","-c","ps -ef | grep socat | grep OPENSSL-LISTEN:"+allPlugins[anode["service"]]["port"]+" | grep -v bash | awk '{print $2}'").Output()
        if err != nil {logs.Error("DeployStapService deploy socket-network Error: "+err.Error()); return err}
        pidValue = strings.Split(string(pid), "\n")
        if pidValue[0] != "" {
            err = ndb.UpdatePluginValue(anode["service"],"pid",pidValue[0]); if err != nil {logs.Error("DeployStapService change pid to value Error: "+err.Error()); return err}
        }else{
            return nil
        }
    }else if anode["type"] == "socket-pcap" {
        pid, err := exec.Command("bash","-c","ps -ef | grep socat | grep OPENSSL-LISTEN:"+allPlugins[anode["service"]]["port"]+" | grep -v bash | awk '{print $2}'").Output()
        if err != nil {logs.Error("DeployStapService deploy socket-network Error: "+err.Error()); return err}
        pidValue := strings.Split(string(pid), "\n")
        if pidValue[0] != "" {
            return nil
        }

        cmd := exec.Command("bash","-c","/usr/bin/socat -d OPENSSL-LISTEN:"+allPlugins[anode["service"]]["port"]+",reuseaddr,pf=ip4,fork,cert="+allPlugins[anode["service"]]["cert"]+",verify=0 SYSTEM:\"tcpdump -n -r - -s 0 -G 50 -W 100 -w "+allPlugins[anode["service"]]["pcap-path"]+allPlugins[anode["service"]]["pcap-prefix"]+"%d%m%Y%H%M%S.pcap "+allPlugins[anode["service"]]["bpf"]+"\" &")
        var errores bytes.Buffer
        cmd.Stdout = &errores
        err = cmd.Start()
        if err != nil {logs.Error("DeployStapService deploying Error: "+err.Error()); return err}        

        pid, err = exec.Command("bash","-c","ps -ef | grep socat | grep OPENSSL-LISTEN:"+allPlugins[anode["service"]]["port"]+" | grep -v bash | awk '{print $2}'").Output()
        if err != nil {logs.Error("DeployStapService deploy socket-network Error: "+err.Error()); return err}
        pidValue = strings.Split(string(pid), "\n")
        if pidValue[0] != "" {
            err = ndb.UpdatePluginValue(anode["service"],"pid",pidValue[0]); if err != nil {logs.Error("DeployStapService change pid to value Error: "+err.Error()); return err}
        }else{
            return nil
        }
    }else if anode["type"] == "network-socket" {
        for x := range allPlugins{
            if allPlugins[x]["type"] == anode["type"] && allPlugins[x]["collector"] == anode["collector"] && allPlugins[x]["port"] == anode["port"] && allPlugins[x]["interface"] == anode["interface"] && allPlugins[x]["pid"] != "none"{
                logs.Error("This network-socket has been deployed yet.")
                return errors.New("This network-socket has been deployed yet.")
            }
        }

        cmd := exec.Command("bash","-c","/usr/sbin/tcpdump -n -i "+allPlugins[anode["service"]]["interface"]+" -s 0 -w - "+allPlugins[anode["service"]]["bpf"]+" | /usr/bin/socat - OPENSSL:"+allPlugins[anode["service"]]["collector"]+":"+allPlugins[anode["service"]]["port"]+",cert="+allPlugins[anode["service"]]["cert"]+",verify=0,forever,retry=10,interval=5 &")
        err = cmd.Start()
        if err != nil {logs.Error("DeployStapService deploying Error: "+err.Error()); return err}

        pid, err := exec.Command("bash","-c","ps -ef | grep OPENSSL:"+allPlugins[anode["service"]]["collector"]+":"+allPlugins[anode["service"]]["port"]+" | grep -v bash | awk '{print $2}'").Output()
        if err != nil {logs.Error("DeployStapService deploy socket-pcap Error: "+err.Error()); return err}
        pidValue := strings.Split(string(pid), "\n")
        if pidValue[0] != "" {
            err = ndb.UpdatePluginValue(anode["service"],"pid",pidValue[0]); if err != nil {logs.Error("DeployStapService change pid to value Error: "+err.Error()); return err}
        }else{
            return nil
        }
    }
    
    return nil
}

func StopStapService(anode map[string]string)(err error) {
    allPlugins,err := ndb.GetPlugins()
    pidToInt,err := strconv.Atoi(allPlugins[anode["service"]]["pid"])
    if err != nil {logs.Error("DeployStapService pid to int error: "+err.Error()); return err}
    process, err := os.FindProcess(pidToInt)
    if err != nil {logs.Error("DeployStapService pid process not found: "+err.Error()); return err}
    err = process.Kill()
    if err != nil {logs.Error("DeployStapService Kill pid process Error: "+err.Error()); return err}
    err = ndb.UpdatePluginValue(anode["service"],"pid","none") ; if err != nil {logs.Error("DeployStapService change pid to none Error: "+err.Error()); return err}

    return nil
}