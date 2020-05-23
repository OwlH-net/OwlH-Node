package zeek

import (
    "bufio"
    "encoding/json"
    "errors"
    "fmt"
    "github.com/astaxie/beego/logs"
    "os"
    "os/exec"
    "owlhnode/database"
    "owlhnode/utils"
    "regexp"
    "strconv"
    "strings"
    "time"
)

type Diag map[string]Diagnodes
type Diagnodes map[string][]string

var Diagnosis = make(Diag)

type NKeys map[string]string
type Nodes map[string]NKeys
type ZeekNodeCfg struct {
    IsCluster bool             `json:"iscluster"`
    Manager   bool             `json:"ismanager"`
    Managerip string           `json:"managerip"`
    Nodes     map[string]NKeys `json:"nodes"`
}

var GlobalZeekCFG = ZeekNodeCfg{}

type Nodekey struct {
    Key   string
    Value string
}

type NodeConfig struct {
    Keys []Nodekey
}

var NodesConfig = map[string]NodeConfig{}

type Zeek struct {
    Path      bool              `json:"path"`
    Rol       string            `json:"role"`
    Bin       bool              `json:"bin"`
    Action    string            `json:"action"`
    Running   []ZeekNodeStatus  `json:"running"`
    Mode      string            `json:"mode"`
    Managed   bool              `json:"managed"`
    Manager   bool              `json:"manager"`
    ManagerIP string            `json:"managerip"`
    Nodes     []ZeekNode        `json:"nodes"`
    Extra     map[string]string `json:"extra"`
}

type ZeekKeys struct {
    Key   string `json:"key"`
    Value string `json:"value"`
}

type ZeekNode struct {
    Name       string     `json:"name"`
    Host       string     `json:"host"`
    Status     string     `json:"status"`
    Type       string     `json:"type"`
    NInterface string     `json:"interface"`
    Pid        string     `json:"pid"`
    Started    string     `json:"started"`
    Extra      []ZeekKeys `json:"extra"`
}

type ZeekNodeStatus struct {
    Status string `json:"status"`
    Nodes  int    `json:"nodes"`
}

type ConfigZeek struct {
    Verbose bool `json:"verbose"`
    Managed bool `json:"managed"`
}

var zeekConfig ConfigZeek

func logDefaulting(value, err, vdefault string) {
    logs.Warn("ZEEK - Error getting '%s' value from main.conf: %s", value, err)
    logs.Warn("ZEEK - Defaulting '%s' value to -> %s", value, vdefault)
}

func logDefaultingGrpKey(grp, value, err, vdefault string) {
    logs.Warn("ZEEK - Error getting '%s'-'%s' value from main.conf: %s", grp, value, err)
    logs.Warn("ZEEK - Defaulting '%s'-%s' value to -> %s", grp, value, vdefault)
}

func canIManage() bool {
    Loadconfig()
    if GlobalZeekCFG.IsCluster && GlobalZeekCFG.Manager {
        return true
    }
    if !GlobalZeekCFG.IsCluster {
        return true
    }
    return false
}

func Loadconfig() {
    verbose, _ := ndb.GetMainconfParam("zeek", "verbose")
    if verbose == "" {
        bverbose, err := utils.GetKeyValueBool("zeek", "verbose")
        if err != nil {
            zeekConfig.Verbose = true
            logDefaultingGrpKey("zeek", "verbose", err.Error(), "true")
        } else {
            zeekConfig.Verbose = bverbose
        }
    } else {
        zeekConfig.Verbose, _ = strconv.ParseBool(verbose)
    }

    // managed, _ := ndb.GetMainconfParam("zeek", "managed")
    managed, _ := ndb.GetMainconfParam("zeek", "status")
    if managed == "" {
        bmanaged, err := utils.GetKeyValueBool("zeek", "managed")
        if err != nil {
            zeekConfig.Managed = true
            logDefaultingGrpKey("zeek", "managed", err.Error(), "true")
        } else {
            zeekConfig.Managed = bmanaged
        }
    } else {
        if managed == "enabled" || managed == "Enabled" {
            zeekConfig.Managed = true
        } else {
            zeekConfig.Managed = false
        }
    }
    ReadZeekNodeData()
}

func ReadZeekNodeData() {
    nodeConfig, err := utils.GetKeyValueString("zeek", "nodeconfig")
    if err != nil {
        nodeConfig = "/usr/local/zeek/etc/node.cfg"
        logDefaultingGrpKey("zeek", "nodeconig", err.Error(), "/usr/local/zeek/etc/node.cfg")
    }

    file, err := os.Open(nodeConfig)
    if err != nil {
        logs.Error(err)
        return
    }
    defer file.Close()

    currentnode := ""
    GlobalZeekCFG.Nodes = make(Nodes)

    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        var comment = regexp.MustCompile(`^#`)
        if comment.MatchString(scanner.Text()) || scanner.Text() == "" {
            continue
        }

        var nodename = regexp.MustCompile(`^\[([^\]]+)]`)
        if val := nodename.FindStringSubmatch(scanner.Text()); val != nil {
            currentnode = val[1]
            if _, ok := GlobalZeekCFG.Nodes[currentnode]; !ok {
                GlobalZeekCFG.Nodes[currentnode] = make(NKeys)
            }
        }

        var nodeparam = regexp.MustCompile(`([^=]+)=(.*)`)
        if val := nodeparam.FindStringSubmatch(scanner.Text()); val != nil {
            if val[1] == "type" {
                if val[2] == "manager" {
                    GlobalZeekCFG.IsCluster = true
                } else if val[2] == "standalone" {
                    GlobalZeekCFG.IsCluster = false
                }
            }
            GlobalZeekCFG.Nodes[currentnode][val[1]] = val[2]
        }
    }
    GlobalZeekCFG.Manager, GlobalZeekCFG.Managerip = AmITheManager()
}

func ManageZeekFromNode() (canbedone bool, err error) {
    if !zeekConfig.Managed {
        str := "ZEEK Management - We don't manage zeek automatically from node. User will do by API calls"
        logs.Warn(str)
        return false, errors.New(str)
    }

    if !canIManage() {
        str := fmt.Sprintf("ZEEK Management - This node belongs to a Cluster, but is not the manager, manager is at %s ", GlobalZeekCFG.Managerip)
        logs.Warn(str)
        return false, errors.New(str)
    }
    return true, nil
}

func StoppingZeekAtNodeEnd() error {
    if iCan, err := ManageZeekFromNode(); !iCan {
        return err
    }

    output, err := StopZeek()
    if err != nil {
        logs.Error("ERROR - ZEEK - STOP at node EXIT -> %s", err.Error())
    }
    if zeekConfig.Verbose {
        logs.Info("ZEEK - Stop at node EXIT - output -> %s", output)
    }
    return err
}

func StartingZeekAtNodeInit() error {
    if iCan, err := ManageZeekFromNode(); !iCan {
        return err
    }

    if isZeekRunning() {
        return nil
    }

    output, err := StartZeek("deploy")
    if err != nil {
        logs.Error("ERROR - ZEEK - DEPLOY At Node Start -> %s", err.Error())
    }
    if zeekConfig.Verbose {
        logs.Info("ZEEK - DEPLOY At node start -> %s", output)
    }
    return err
}

func AmITheManager() (manager bool, ip string) {
    for node := range GlobalZeekCFG.Nodes {
        if GlobalZeekCFG.Nodes[node]["type"] == "manager" {
            hostip := GlobalZeekCFG.Nodes[node]["host"]
            GlobalZeekCFG.Manager = utils.IsLocalAddress(hostip)
            GlobalZeekCFG.Managerip = hostip
            return GlobalZeekCFG.Manager, GlobalZeekCFG.Managerip
        }
    }
    return false, ""
}

func ZeekPath() (exists bool) {
    if zeekConfig.Verbose {
        logs.Info("Zeek getting zeek path")
    }
    var err error
    path, err := utils.GetKeyValueString("zeek", "zeekpath")
    if err != nil {
        logs.Error("ZeekPath Error getting data from main.conf: " + err.Error())
        return false
    }

    if _, err := os.Stat(path); os.IsNotExist(err) {
        logs.Error("Zeek path %s can't be found : %s", path, err.Error())
        return false
    }
    return true
}

func ZeekBin() (exists bool) {
    if zeekConfig.Verbose {
        logs.Info("Zeek getting zeek binary file")
    }
    var err error
    zeekctl, err := utils.GetKeyValueString("zeek", "zeekctl")
    if err != nil {
        logDefaultingGrpKey("zeek", "zeekctl", err.Error(), "/usr/local/zeek/bin/zeekctl")
        zeekctl = "/usr/local/zeek/bin/zeekctl"
    }

    _, err = os.Stat(zeekctl)
    if err != nil {
        logs.Error("Zeek OS path err: " + err.Error())
        return false
    }

    if os.IsNotExist(err) {
        logs.Error("Zeek path not exist: " + err.Error())
        return false
    }
    return true
}

func ZeekRunning() (running bool) {
    logs.Info("Is Zeek running? nothing done - DEPRECATED!!!!!!")
    return false
}

func isZeekRunning() bool {
    nodes, _ := ZeekStatus()
    for node := range nodes {
        if nodes[node].Status != "running" {
            return false
        }
    }
    return true
}

func ZeekStatus() (zeekstatus []ZeekNode, err error) {
    if zeekConfig.Verbose {
        logs.Info("Zeek checking zeek status")
    }
    if !canIManage() {
        str := fmt.Sprintf("ZEEK Management - This node belongs to a Cluster, but is not the manager, manager is at %s ", GlobalZeekCFG.Managerip)
        logs.Warn(str)
        return zeekstatus, errors.New(str)
    }

    zeekctl, err := utils.GetKeyValueString("zeek", "zeekctl")
    if err != nil {
        logDefaultingGrpKey("zeek", "zeekctl", err.Error(), "/usr/local/zeek/bin/zeekctl")
        zeekctl = "/usr/local/zeek/bin/zeekctl"
    }

    currentStatus, err := utils.GetKeyValueString("zeek", "status")
    if err != nil {
        logDefaultingGrpKey("zeek", "status", err.Error(), "status")
        currentStatus = "status"
    }

    output, _ := exec.Command(zeekctl, currentStatus).Output()
    if len(output) == 0 {
        str := "ZEEK STATUS -> There is no output from status command"
        logs.Error(str)
        return zeekstatus, errors.New(str)
    }

    nodes := []ZeekNode{}
    outputlines := strings.Split(string(output), "\n")
    for outputline := range outputlines {
        line := strings.Fields(outputlines[outputline])
        if len(line) > 1 {
            if strings.Contains(line[1], "manager") || strings.Contains(line[1], "logger") || strings.Contains(line[1], "proxy") || strings.Contains(line[1], "worker") || strings.Contains(line[1], "standalone") {
                node := ZeekNode{}
                node.Name = line[0]
                node.Type = line[1]
                node.Host = line[2]
                node.Status = line[3]
                if len(line) > 4 {
                    node.Pid = line[4]
                    node.Started = strings.Join(line[5:], " ")
                }
                nodes = append(nodes, node)
            }
        }
    }
    return nodes, nil
}

func ZeekCurrentStatus() (status string, err error) {
    zeek, err := ZeekStatus()
    if err != nil {
        return "", err
    }
    for x := range zeek {
        status = zeek[x].Status
    }
    return status, nil
}

func GetZeek() (zeek Zeek, err error) {
    Loadconfig()
    if zeekConfig.Verbose {
        logs.Info("ZEEK - Getting zeek values")
    }

    zeek.Path = ZeekPath()
    zeek.Bin = ZeekBin()
    if !zeek.Path || !zeek.Bin {
        return zeek, errors.New("Zeek path or binary does not exist")
    }
    if GlobalZeekCFG.IsCluster {
        zeek.Mode = "cluster"
    } else {
        zeek.Mode = "standalone"
    }
    zeek.Managed = zeekConfig.Managed
    zeek.ManagerIP = GlobalZeekCFG.Managerip
    zeek.Manager = GlobalZeekCFG.Manager
    nodes, err := ZeekStatus()
    if err != nil {
        logs.Error(err.Error())
    }
    zeek.Nodes = nodes

    for node := range nodes {
        updated := false
        for nstatus := range zeek.Running {
            if zeek.Running[nstatus].Status == nodes[node].Status {
                updated = true
                zeek.Running[nstatus].Nodes++
                break
            }
        }
        if !updated {
            newStatus := ZeekNodeStatus{}
            newStatus.Status = nodes[node].Status
            newStatus.Nodes = 1
            zeek.Running = append(zeek.Running, newStatus)
        }
    }
    logs.Debug("%+v", zeek)
    return zeek, nil
}

func SetZeek(zeekdata Zeek) (newzeekdata Zeek, err error) {
    if zeekConfig.Verbose {
        logs.Info("Zeek set zeek values")
    }

    // convert to useful call
    if zeekConfig.Verbose {
        for node := range zeekdata.Nodes {
            logs.Warn("=============")
            logs.Warn("name - " + zeekdata.Nodes[node].Name)
            logs.Warn("interface - " + zeekdata.Nodes[node].NInterface)
            logs.Warn("host - " + zeekdata.Nodes[node].Host)
            logs.Warn("type - " + zeekdata.Nodes[node].Type)
            logs.Warn("=============")
            logs.Warn("======= EXTRA ========")
            for extra := range zeekdata.Nodes[node].Extra {
                logs.Warn(zeekdata.Nodes[node].Extra[extra])
                logs.Warn("key - " + zeekdata.Nodes[node].Extra[extra].Key + " -- " + zeekdata.Nodes[node].Extra[extra].Value)
            }
            logs.Warn("======= EXTRA ========")
        }
    }

    newzeekdata, err = GetZeek()
    if err != nil {
        return newzeekdata, err
    }

    return newzeekdata, nil
}

func ZeekMode() (mode string) {
    if zeekConfig.Verbose {
        logs.Info("Zeek Mode Function DEPRECATED")
    }

    return ""
}

func ZeekManaged() (ismanaged bool) {
    Loadconfig()
    return zeekConfig.Managed
}

//Run zeek
func RunZeek() (data string, err error) {
    if zeekConfig.Verbose {
        logs.Info("ZEEK - Start Zeek...")
    }
    if !canIManage() {
        str := fmt.Sprintf("ZEEK Management - This node belongs to a Cluster, but is not the manager, manager is at %s ", GlobalZeekCFG.Managerip)
        logs.Warn(str)
        return "", errors.New(str)
    }
    zeekctl, err := utils.GetKeyValueString("zeek", "zeekctl")
    if err != nil {
        logs.Error("RunZeek Error getting data from main.conf: " + err.Error())
        return "", err
    }
    zeekStart, err := utils.GetKeyValueString("zeek", "start")
    if err != nil {
        logs.Error("RunZeek Error getting data from main.conf: " + err.Error())
        return "", err
    }

    // err = utils.RunCommand(zeekctl,zeekStart)
    cmd := exec.Command(zeekctl, zeekStart)
    err = cmd.Run()
    if err != nil {
        logs.Error("Error starting zeek: " + err.Error())
        return "", err
    }

    //update mainconf status
    err = ndb.UpdateMainconfValue("zeek", "previousStatus", "start")
    if err != nil {
        logs.Error("RunZeek error changing mainconf database previous status: " + err.Error())
        return "", err
    }

    return "Zeek system is on", nil
}

//Start Zeek
func StartZeek(action string) (data string, err error) {
    if zeekConfig.Verbose {
        logs.Info("ZEEK - Start Zeek by action", action)
    }
    if !canIManage() {
        str := fmt.Sprintf("ZEEK Management - This node belongs to a Cluster, but is not the manager, manager is at %s ", GlobalZeekCFG.Managerip)
        logs.Warn(str)
        return "", errors.New(str)
    }
    // ACTION - Start or Deploy
    getaction := "deploy"
    if action != "" {
        getaction = action
    }
    if zeekConfig.Verbose {
        logs.Warn("Starting Zeek by action -> %s", action)
    }
    zeekctl, err := utils.GetKeyValueString("zeek", "zeekctl")
    if err != nil {
        logs.Error("StartZeek Error getting data from main.conf: " + err.Error())
        return "", err
    }
    realaction, err := utils.GetKeyValueString("zeek", getaction)
    if err != nil {
        logs.Error("StartZeek Error getting data from main.conf: " + err.Error())
        return "", err
    }

    output, err := exec.Command(zeekctl, realaction).Output()
    if err != nil {
        logs.Error("Error launching StartZeek: " + err.Error())
        return "", err
    }
    if zeekConfig.Verbose {
        logs.Info(string(output))
    }
    str := "ZEEK - Start Zeek command done. Check Status"
    if zeekConfig.Verbose {
        logs.Info(str)
    }
    return string(output), nil
}

func StartingZeek() (err error) {
    if zeekConfig.Verbose {
        logs.Info("Starting Zeek...")
    }
    if !canIManage() {
        str := fmt.Sprintf("ZEEK Management - This node belongs to a Cluster, but is not the manager, manager is at %s ", GlobalZeekCFG.Managerip)
        logs.Warn(str)
        return errors.New(str)
    }
    zeekctl, err := utils.GetKeyValueString("zeek", "zeekctl")
    if err != nil {
        logs.Error("StartingZeek Error getting data from main.conf: " + err.Error())
    }
    start, err := utils.GetKeyValueString("zeek", "start")
    if err != nil {
        logs.Error("StartingZeek Error getting data from main.conf: " + err.Error())
    }

    // err = utils.RunCommand(cmd,start)
    // err = utils.StartCommand(cmd,start)
    // cmd := exec.Command(zeekctl, start)
    // err = cmd.Run()
    // if err != nil {
    //     logs.Error("Error deploying zeek: " + err.Error())
    //     return err
    // }

    output, err := exec.Command(zeekctl, start).Output()
    if err != nil {
        logs.Error("Error Starting Zeek: " + err.Error())
        return err
    }
    if zeekConfig.Verbose {
        logs.Info(string(output))
    }

    //update mainconf status
    err = ndb.UpdateMainconfValue("zeek", "previousStatus", "start")
    if err != nil {
        logs.Error("StartingZeek error changing mainconf database previous status: " + err.Error())
        return err
    }

    return nil
}

// //Stop zeek
func StopZeek() (data string, err error) {
    if zeekConfig.Verbose {
        logs.Info("ZEEK - Stop Zeek...")
    }
    if !canIManage() {
        str := fmt.Sprintf("ZEEK Management - This node belongs to a Cluster, but is not the manager, manager is at %s ", GlobalZeekCFG.Managerip)
        logs.Warn(str)
        return "", errors.New(str)
    }
    zeekctl, err := utils.GetKeyValueString("zeek", "zeekctl")
    if err != nil {
        logs.Error("ZEEK StopZeek Error getting 'zeekctl' key from main.conf: %s", err.Error())
        zeekctl = "/usr/local/zeek/bin/zeekctl"
    }
    stop, err := utils.GetKeyValueString("zeek", "stop")
    if err != nil {
        logs.Error("StopZeek Error getting 'stop' key from main.conf: %s", err.Error())
        stop = "stop"
    }

    // err = utils.RunCommand(cmd, stop)
    // cmd := exec.Command(zeekctl, stop)
    // err = cmd.Run()
    // if err != nil {
    //     logs.Error("ERROR - ZEEK - Stopping zeek: %s", err.Error())
    //     return "", err
    // }

    output, err := exec.Command(zeekctl, stop).Output()
    if err != nil {
        logs.Error("ERROR - ZEEK - Stopping zeek: %s", err.Error())
        return "", err
    }
    if zeekConfig.Verbose {
        logs.Info(string(output))
    }

    //update mainconf status
    err = ndb.UpdateMainconfValue("zeek", "previousStatus", "stop")
    if err != nil {
        logs.Error("ZEEK - Can't update Zeek previous Status in mainconf db -> %s ", err.Error())
        return "", err
    }

    str := "ZEEK - Stop Zeek command done. Check Status"
    if zeekConfig.Verbose {
        logs.Info(str)
    }
    return str, nil
}

//Stop zeek
func StoppingZeek() (err error) {
    if zeekConfig.Verbose {
        logs.Info("Stopping Zeek...")
    }
    if !canIManage() {
        str := fmt.Sprintf("ZEEK Management - This node belongs to a Cluster, but is not the manager, manager is at %s ", GlobalZeekCFG.Managerip)
        logs.Warn(str)
        return errors.New(str)
    }

    zeekctl, err := utils.GetKeyValueString("zeek", "zeekctl")
    if err != nil {
        logs.Error("StoppingZeek Error getting binary data from main.conf: " + err.Error() + ". Defaulting...")
        zeekctl = "/usr/local/zeek/bin/zeekctl"
    }

    stop, err := utils.GetKeyValueString("zeek", "stop")
    if err != nil {
        logs.Error("StoppingZeek Error getting zeek stop parameter data from main.conf: " + err.Error() + ". Defaulting...")
        stop = "stop"
    }
    t, err := utils.GetKeyValueString("zeek", "wait")
    if err != nil {
        logs.Error("StoppingZeek Error getting zeek wait parameter data from main.conf: " + err.Error() + ". Defaulting...")
        t = "5"
    }
    tDuration, err := strconv.Atoi(t)

    cmd := exec.Command(zeekctl, stop)
    cmd.Run()

    time.Sleep(time.Second * time.Duration(tDuration))

    cmd = exec.Command(zeekctl, stop)
    cmd.Run()

    time.Sleep(time.Second * time.Duration(tDuration))

    return nil
}

func ChangeZeekPreviousStatus() (err error) {
    return nil
    err = ndb.UpdateMainconfValue("zeek", "previousStatus", "stop")
    if err != nil {
        logs.Error("ChangeMainServiceStatus error changing mainconf database status: " + err.Error())
        return err
    }
    return nil
}

//Deploy zeek
func DeployZeek() (err error) {
    if zeekConfig.Verbose {
        logs.Info("Deploying Zeek...")
    }
    if !canIManage() {
        str := fmt.Sprintf("ZEEK Management - This node belongs to a Cluster, but is not the manager, manager is at %s ", GlobalZeekCFG.Managerip)
        logs.Warn(str)
        return errors.New(str)
    }
    zeekctl, err := utils.GetKeyValueString("zeek", "zeekctl")
    if err != nil {
        logs.Error("DeployZeek Error getting data from main.conf: " + err.Error())
    }
    deploy, err := utils.GetKeyValueString("zeek", "deploy")
    if err != nil {
        logs.Error("DeployZeek Error getting data from main.conf: " + err.Error())
    }

    output, err := exec.Command(zeekctl, deploy).Output()
    if err != nil {
        logs.Error("ERROR - ZEEK - Stopping zeek: %s", err.Error())
        return err
    }
    if zeekConfig.Verbose {
        logs.Info(string(output))
    }

    //update mainconf status
    err = ndb.UpdateMainconfValue("zeek", "previousStatus", "start")
    if err != nil {
        logs.Error("DeployZeek error changing mainconf database previous status: " + err.Error())
        return err
    }

    return nil
}

func ChangeZeekMode(anode map[string]string) (err error) {
    if zeekConfig.Verbose {
        logs.Info("Changing Zeek mode...")
    }
    err = ndb.UpdateMainconfValue("zeek", "mode", anode["mode"])
    if err != nil {
        logs.Error("Error ChangeZeekMode: " + err.Error())
        return err
    }
    SyncCluster(nil, "standalone")
    return err
}

func AddClusterValue(anode map[string]string) (err error) {
    if zeekConfig.Verbose {
        logs.Info("Zeek add new cluster value")
    }
    count, err := ndb.CountDBEntries(anode["type"])
    if err != nil {
        logs.Error("Error AddClusterValue type: " + err.Error())
        return err
    }
    count++
    err = ndb.InsertClusterData(anode["type"]+"-"+strconv.Itoa(count), "type", anode["type"])
    if err != nil {
        logs.Error("Error AddClusterValue type: " + err.Error())
        return err
    }
    err = ndb.InsertClusterData(anode["type"]+"-"+strconv.Itoa(count), "host", anode["host"])
    if err != nil {
        logs.Error("Error1 AddClusterValue host: " + err.Error())
        return err
    }
    if anode["type"] == "worker" {
        err = ndb.InsertClusterData(anode["type"]+"-"+strconv.Itoa(count), "interface", anode["interface"])
        if err != nil {
            logs.Error("Error AddClusterValue interface: " + err.Error())
            return err
        }
    }
    return err
}

func PingCluster() (data map[string]map[string]string, err error) {
    if zeekConfig.Verbose {
        logs.Info("Zeek get cluster values")
    }
    data, err = ndb.GetClusterData()
    if err != nil {
        logs.Error("Error Zeek/PingCluster: " + err.Error())
        return nil, err
    }
    return data, err
}

func EditClusterValue(anode map[string]string) (err error) {
    if zeekConfig.Verbose {
        logs.Info("Zeek edit cluster values")
    }
    err = ndb.UpdateClusterValue(anode["type"], "host", anode["host"])
    if err != nil {
        logs.Error("Error Zeek/EditClusterValue: " + err.Error())
        return err
    }
    if anode["cluster"] == "worker" {
        err = ndb.UpdateClusterValue(anode["type"], "interface", anode["interface"])
        if err != nil {
            logs.Error("Error Zeek/EditClusterValue: " + err.Error())
            return err
        }
    }
    return err
}

func DeleteClusterValue(anode map[string]string) (err error) {
    if zeekConfig.Verbose {
        logs.Info("Zeek delete cluster values")
    }
    err = ndb.DeleteClusterValue(anode["type"])
    if err != nil {
        logs.Error("Error Zeek/DeleteClusterValue: " + err.Error())
        return err
    }
    //change indentifier
    countWorker := 1
    countProxy := 1
    data, err := ndb.GetClusterData()
    if err != nil {
        logs.Error("Error Zeek/DeleteClusterValue: " + err.Error())
        return err
    }
    err = ndb.DeleteAllClusters()
    if err != nil {
        logs.Error("Error Zeek/DeleteClusterValue: " + err.Error())
        return err
    }

    for id := range data {
        if id == "manager" || id == "logger" {
            err = ndb.InsertClusterData(id, "host", data[id]["host"])
            if err != nil {
                logs.Error("Error DeleteClusterValue manager: " + err.Error())
                return err
            }
        } else {
            if data[id]["type"] == "worker" {
                err = ndb.InsertClusterData(data[id]["type"]+"-"+strconv.Itoa(countWorker), "type", data[id]["type"])
                if err != nil {
                    logs.Error("Error DeleteClusterValue type: " + err.Error())
                    return err
                }
                err = ndb.InsertClusterData(data[id]["type"]+"-"+strconv.Itoa(countWorker), "host", data[id]["host"])
                if err != nil {
                    logs.Error("Error DeleteClusterValue host: " + err.Error())
                    return err
                }
                err = ndb.InsertClusterData(data[id]["type"]+"-"+strconv.Itoa(countWorker), "interface", data[id]["interface"])
                if err != nil {
                    logs.Error("Error DeleteClusterValue type: " + err.Error())
                    return err
                }
                countWorker++
            } else {
                err = ndb.InsertClusterData(data[id]["type"]+"-"+strconv.Itoa(countProxy), "type", data[id]["type"])
                if err != nil {
                    logs.Error("Error DeleteClusterValue type: " + err.Error())
                    return err
                }
                err = ndb.InsertClusterData(data[id]["type"]+"-"+strconv.Itoa(countProxy), "host", data[id]["host"])
                if err != nil {
                    logs.Error("Error DeleteClusterValue host: " + err.Error())
                    return err
                }
                countProxy++
            }
        }
    }

    return err
}

func SyncCluster(anode map[string]string, clusterType string) (err error) {
    if zeekConfig.Verbose {
        logs.Info("Zeek synchronizing cluster values")
    }
    path, err := utils.GetKeyValueString("zeek", "nodeconfig")
    if err != nil {
        logs.Error("SyncCluster Error readding main.conf: " + err.Error())
    }

    h := 0
    fileContent := make(map[int]string)

    if clusterType == "standalone" {
        fileContent[h] = "[zeek]"
        h++
        fileContent[h] = "type=standalone"
        h++
        fileContent[h] = "host=localhost"
        h++
        fileContent[h] = "interface=" + anode["value"]
        h++
    } else if clusterType == "cluster" {
        data, err := ndb.GetClusterData()
        if err != nil {
            logs.Error("Error Zeek/SyncCluster: " + err.Error())
            return err
        }

        for t := range data {
            if t == "logger" {
                fileContent[h] = "[logger]"
                h++
                fileContent[h] = "type=logger"
                h++
                fileContent[h] = "host=" + data[t]["host"]
                h++
                fileContent[h] = ""
                h++
            } else if t == "manager" {
                fileContent[h] = "[manager]"
                h++
                fileContent[h] = "type=manager"
                h++
                fileContent[h] = "host=" + data[t]["host"]
                h++
                fileContent[h] = ""
                h++
            } else if data[t]["type"] == "proxy" {
                fileContent[h] = "[" + t + "]"
                h++
                fileContent[h] = "type=" + data[t]["type"]
                h++
                fileContent[h] = "host=" + data[t]["host"]
                h++
                fileContent[h] = ""
                h++
            } else if data[t]["type"] == "worker" {
                fileContent[h] = "[" + t + "]"
                h++
                fileContent[h] = "type=" + data[t]["type"]
                h++
                fileContent[h] = "host=" + data[t]["host"]
                h++
                fileContent[h] = "interface=" + data[t]["interface"]
                h++
                fileContent[h] = ""
                h++
            }
        }
    }

    saveIntoFile, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
    if err != nil {
        logs.Error("Error SyncCluster readding file: " + err.Error())
        return err
    }
    defer saveIntoFile.Close()
    saveIntoFile.Truncate(0)
    saveIntoFile.Seek(0, 0)
    for x := 0; x < h; x++ {
        _, err = saveIntoFile.WriteAt([]byte(fileContent[x]+"\n"), 0) // Write at 0 beginning
        if err != nil {
            logs.Error("SyncCluster failed writing to file: %s", err)
            return err
        }
    }

    return err
}

func SavePolicyFiles(files map[string]map[string][]byte) (err error) {
    if zeekConfig.Verbose {
        logs.Info("Zeek save policy files")
    }
    for nodePath, file := range files {
        //check path
        if _, err := os.Stat(nodePath); os.IsNotExist(err) {
            os.MkdirAll(nodePath, os.ModePerm)
        }

        for file := range file {
            err = utils.WriteNewDataOnFile(nodePath+"/"+file, files[nodePath][file])
            if err != nil {
                logs.Error("SavePolicyFiles Error writting data into " + nodePath + "/" + file + " file: " + err.Error())
                return err
            }
        }
    }
    return nil
}

func SyncClusterFile(anode map[string][]byte) (err error) {
    path, err := utils.GetKeyValueString("zeek", "nodeconfig")
    if err != nil {
        logs.Error("SyncClusterFile Error readding main.conf: " + err.Error())
    }

    err = utils.WriteNewDataOnFile(path, anode["data"])
    if err != nil {
        logs.Error("zeek/SyncClusterFile Error writting cluster file content: " + err.Error())
        return err
    }

    return nil
}

func SyncZeekValues(anode map[string]string) (err error) {
    if zeekConfig.Verbose {
        logs.Info("Zeek synchronize values")
    }
    for x := range anode {
        if x == "nodeConfig" {
            path, err := utils.GetKeyValueString("zeek", "nodeconfig")
            if err != nil {
                logs.Error("zeek/SyncZeekValues Error readding main.conf: " + err.Error())
            }

            err = utils.BackupFullPath(path)
            if err != nil {
                logs.Error("zeek/SyncZeekValues Error backing up node.cfg file before overwrite: " + err.Error())
                return err
            }
            err = utils.WriteNewDataOnFile(path, []byte(anode["nodeConfig"]))
            if err != nil {
                logs.Error("zeek/SyncZeekValues Error writting new file content: " + err.Error())
                return err
            }

        } else if x == "networksConfig" {
            path, err := utils.GetKeyValueString("zeek", "networkconfig")
            if err != nil {
                logs.Error("zeek/SyncZeekValues Error readding main.conf: " + err.Error())
            }

            err = utils.BackupFullPath(path)
            if err != nil {
                logs.Error("zeek/SyncZeekValues Error backing up networks.cfg file before overwrite: " + err.Error())
                return err
            }
            err = utils.WriteNewDataOnFile(path, []byte(anode["networksConfig"]))
            if err != nil {
                logs.Error("zeek/SyncZeekValues Error writting new file content: " + err.Error())
                return err
            }

        } else {
            if _, err := os.Stat(anode["dst"]); os.IsNotExist(err) {
                logs.Error("zeek/SyncZeekValues Destiny file down't exists: " + err.Error())
                return err
            }
            err = utils.BackupFullPath(anode["dst"])
            if err != nil {
                logs.Error("zeek/SyncZeekValues Error backing up file before overwrite: " + err.Error())
                return err
            }
            if x == "policiesMaster" {
                err = utils.WriteNewDataOnFile(anode["dst"], []byte(anode["policiesMaster"]))
                if err != nil {
                    logs.Error("zeek/SyncZeekValues Error writting new file content: " + err.Error())
                    return err
                }
            }
            if x == "variables1" {
                err = utils.WriteNewDataOnFile(anode["dst"], []byte(anode["variables1"]))
                if err != nil {
                    logs.Error("zeek/SyncZeekValues Error writting new file content: " + err.Error())
                    return err
                }
            }
        }
    }

    return err
}

func DiagZeek() (data map[string]string, err error) {
    diag, err := utils.GetKeyValueString("zeek", "diag")
    if err != nil {
        diag = "diag"
        logDefaultingGrpKey("zeek", "diag", err.Error(), "diag")
    }

    zeekctl, err := utils.GetKeyValueString("zeek", "zeekctl")
    if err != nil {
        zeekctl = "/usr/local/zeek/bin/zeekctl"
        logDefaultingGrpKey("zeek", "zeekctl", err.Error(), "/usr/local/zeek/bin/zeekctl")
    }

    output, err := exec.Command(zeekctl, diag).Output()
    linesResult := make(map[string]string)
    linesResult["result"] = string(output)
    linesResult["ack"] = "true"

    if err != nil {
        logs.Error("ERROR - ZEEK - zeek diag: %s", err.Error())
        linesResult["ack"] = "false"
        linesResult["error"] = err.Error()
    }
    if zeekConfig.Verbose {
        logs.Info(string(output))
    }
    parseDiag(string(output))
    return linesResult, err
}

func parseDiag(output string) (data Diag) {
    lines := strings.Split(output, "\n")
    currentnode := "warnings"
    currentitem := ""
    Diagnosis[currentnode] = make(Diagnodes)

    for line := range lines {
        if lines[line] == "" {
            continue
        }
        logs.Debug("zeek - parse diag - parsing line - %s", lines[line])
        var nodename = regexp.MustCompile(`^\[([^\]]+)]`)
        if val := nodename.FindStringSubmatch(lines[line]); val != nil {
            logs.Debug("zeek - parse diag - parsing node - %s", lines[line])

            currentnode = val[1]
            if _, ok := Diagnosis[currentnode]; !ok {
                Diagnosis[currentnode] = make(Diagnodes)
            }
            continue
        }

        var nodeparam = regexp.MustCompile(`==== (.*)`)
        if val := nodeparam.FindStringSubmatch(lines[line]); val != nil {
            logs.Debug("zeek - parse diag - parsing item - %s", lines[line])

            currentitem = val[1]
            Diagnosis[currentnode][currentitem] = []string{}
            continue
        }
        logs.Debug("zeek - parse diag - parsing content - %s", lines[line])

        Diagnosis[currentnode][currentitem] = append(Diagnosis[currentnode][currentitem], lines[line])
    }
    logs.Debug(Diagnosis)
    byteData, _ := json.Marshal(Diagnosis)
    logs.Debug(string(byteData))
    return Diagnosis
}

func Init() {
    Loadconfig()

}
