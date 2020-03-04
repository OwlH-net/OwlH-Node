package knownports

import (
    "github.com/astaxie/beego/logs"
    "owlhnode/database"
    "owlhnode/utils"
    // "owlhnode/analyzer"
    "regexp"
    "strconv"
    "encoding/json"
    "time"
    "os"
    "net"
)

var Status string
var Mode string

type newPortAlert struct {
    Data             Data         `json:"data"`
}

type Data struct {
    Dstport          string       `json:"dstport"`
    Proto            string       `json:"proto"`
    Times            int          `json:"times"`
    Signature        Signature    `json:"alert"`
}

type Signature struct {
    Signature        string       `json:"signature"`
    Signature_id     string       `json:"signature_id"`
}

type LastAlert struct {
    PortProto        string 
    Last             time.Time
    Counter          int
}


func Init(){
    go NewPorts()
}

func GetStatus()(){
    logs.Info("KNOWNPORTS GETSTATUS --- ")
    return
    for {
        _, err := CheckParamKnownports("status")
        if err != nil {
            logs.Error("CheckParamKnownports Status Error: "+err.Error())
        }
        _, err = CheckParamKnownports("mode")
        if err != nil {
            logs.Error("CheckParamKnownports Mode Error: "+err.Error())
        }
        time.Sleep(time.Second * 20)
    }
}

func NewPorts()(){
    return
    var err error
    loadPorts := map[string]map[string]string{}
    loadPorts["knownports"] = map[string]string{}
    loadPorts["knownports"]["file"] = ""
    loadPorts["knownports"]["timeToAlert"] = ""
    loadPorts,err = utils.GetConf(loadPorts)
    file := loadPorts["knownports"]["file"]
    timeout := loadPorts["knownports"]["timeToAlert"]
    if err != nil {
        logs.Error("loadPorts Error getting data from main.conf: "+err.Error())
        return
    }

    Status, err = CheckParamKnownports("status")
    Mode, err = CheckParamKnownports("mode")
    if err != nil {
        logs.Error("CheckParamKnownports Error: "+err.Error())
    }
    for{
        if _,err := os.Stat(file); os.IsNotExist(err){
            logs.Info("KNOWNPORTS -- Waiting file...")

            logs.Info("WAITING FILE KNOWNPORTS")
            time.Sleep(time.Second * 60) 
        }else{
            break
        }
    }
    for Status != "Disabled"{
        if Status == "Reload"{
            anode := make(map[string]string)
            anode["plugin"]="knownports"
            anode["status"]="Enabled"
            err = ChangeStatus(anode)
            if err!=nil {logs.Error("Error changing status from Realod status: "+err.Error())}
        }
        
        newuuid := utils.Generate()
        logs.Info(newuuid + ": starting analyzer for Knwon ports")
        // analyzer.Registerchannel(newuuid)
        portsData, err := LoadPortsData()
        if err!=nil {logs.Error("Error LoadPortsData: "+err.Error())}

        alertList := map[string]LastAlert{}

        loadHomenet := map[string]map[string][]string{}
        loadHomenet["Node"] = map[string][]string{}
        loadHomenet["Node"]["homenet"] = nil
        loadHomenet,err = utils.GetConfArray(loadHomenet)
        IpNet := loadHomenet["Node"]["homenet"]

        if err != nil {
            logs.Error("LoadPortsData NewPorts Error: %s", err.Error())
            Status = "Disabled"
        }
        for {
            line := ""
            Status, err = CheckParamKnownports("status"); if err != nil {logs.Error("CheckParamKnownports Error: "+err.Error())}
            Mode, err = CheckParamKnownports("mode"); if err != nil {logs.Error("CheckParamKnownports Error: "+err.Error())}
            logs.Notice("STATUS knownports: "+Status+"   //   Mode: "+Mode)
            if Status != "Enabled"{
                break
            }
            var protoportRegexp = regexp.MustCompile(`"id.resp_h":"(\d+\.\d+\.\d+\.\d+)","id.resp_p":(\d+),"proto":"(\w+)"`)
            portProtocol := protoportRegexp.FindStringSubmatch(line)
            if portProtocol== nil {continue}

            dstip := portProtocol[1]
            dstipnew, _, _ := net.ParseCIDR(dstip+"/32")
            dstport := portProtocol[2]
            proto := portProtocol[3]
            var protoport = dstport+"/"+proto

            netLocal := false
            for currentNet := range IpNet{
                _, localNet, err := net.ParseCIDR(IpNet[currentNet])
                if err != nil {
                    logs.Error("localNet currentNet Error: "+err.Error())
                }
                if localNet.Contains(dstipnew){
                    logs.Error("dstip is local: "+dstip)
                    netLocal = true
                    break
                }
            }
            if netLocal {
                continue
            }

            logs.Info("dstip is NOT local: "+dstip)
            if Mode == "Learning"{
                logs.Warn("LEARNING MODE")
                notPortprotLearn := false
                for x := range portsData {

                    if portsData[x]["portprot"] == protoport{
                        timeNow := time.Now() 
                        value := strconv.FormatInt(timeNow.Unix(), 10)
                        notPortprotLearn = true
                        err = ndb.UpdateKnownports(value, "last", x)
                        if err != nil {
                            logs.Error("LEARNING MODE --> knownports[last] update error-> %s", err.Error())
                            Status = "Disabled"
                        }
                        break
                    }
                    // if notPortprotLearn == true {break}
                }
                if !notPortprotLearn{
                    uuid := utils.Generate()
                    timeNow := time.Now() 
                    value := strconv.FormatInt(timeNow.Unix(), 10)
                    
                    //insert into MAP portsData
                    // logs.Error(portsData)
                    portsData[uuid] = map[string]string{}
                    portsData[uuid]["port"] = dstport
                    portsData[uuid]["protocol"] = proto
                    portsData[uuid]["portprot"] = dstport+"/"+proto
                    portsData[uuid]["first"] = value
                    portsData[uuid]["last"] = value

                    //insert to DB
                    err = InsertknownportsElements(uuid, "port", dstport)
                    err = InsertknownportsElements(uuid, "protocol", proto)
                    err = InsertknownportsElements(uuid, "portprot", dstport+"/"+proto)
                    err = InsertknownportsElements(uuid, "first", value)
                    err = InsertknownportsElements(uuid, "last", value)
                    if err != nil {
                        logs.Error("knownports insert error-> %s", err.Error())
                        Status = "Disabled"
                    }
                }
            }else{
                logs.Warn("NOT LEARNING MODE")
                notPortprotProd := false
                for x := range portsData { 
                    if portsData[x]["portprot"] == protoport{
                        notPortprotProd = true
                        t := time.Now()
                        value := strconv.FormatInt(t.Unix(), 10)
                        err = ndb.UpdateKnownports(value, "last", x)
                        if err != nil {
                            logs.Error("PRODUCTION MODE --> knownports[last] update error-> %s", err.Error())
                            Status = "Disabled"
                        }
                        break
                    }
                }
                if !notPortprotProd {
                    createAlert := false
                    counter := 0 
                    alerted := LastAlert{}

                    if _, ok := alertList[protoport]; !ok {
                        alerted.PortProto = protoport
                        alerted.Last = time.Now()
                        alerted.Counter = 1
                        counter = 1

                        alertList[protoport] = alerted
                        // logs.Error("first time port alert - " +protoport)
                        createAlert = true
                        
                    } else {
                        alerted = alertList[protoport]
                        tm, _ := strconv.Atoi(timeout)

                        counter = alerted.Counter
                        if time.Now().After(alerted.Last.Add(time.Second*time.Duration(tm))) {
                            createAlert = true
                            alerted.Last = time.Now()
                            alerted.Counter = 0
                        } else {
                            alerted.Counter += 1
                        }
                        alertList[protoport] = alerted
                    }
                    if createAlert {
                        alert := newPortAlert{}
                        data := Data{}
                        signature := Signature{}

                        signature.Signature = "OwlH UNKNOWN PORT - new port detected - "+protoport
                        signature.Signature_id = "8000101"

                        data.Dstport = dstport
                        data.Proto = proto
                        data.Times = counter
                        data.Signature = signature

                        alert.Data = data
                        alertOutput, err := json.Marshal(alert)
                        if err != nil {
                            logs.Error("Marshal Error creating JSON alert output at Production Knownports: %s", err.Error())
                            Status = "Disabled"
                        }

                        AlertLog := map[string]map[string]string{}
                        AlertLog["Node"] = map[string]string{}
                        AlertLog["Node"]["AlertLog"] = ""
                        AlertLog,err = utils.GetConf(AlertLog)
                        AlertLogJson := AlertLog["Node"]["AlertLog"]
                        if err != nil {
                            logs.Error("AlertLog Error getting data from main.conf: "+err.Error())
                            return
                        }
                        err = utils.WriteNewDataOnFile(AlertLogJson, alertOutput)
                        if err != nil {
                            logs.Error("Error creating JSON alert at Production Knownports: %s", err.Error())
                            Status = "Disabled"
                        }
                    }
                }
            }
        }
        Status, err = CheckParamKnownports("status")
        Mode, err = CheckParamKnownports("mode")
        //t.Cleanup()
        //t.Stop()
    }
    logs.Info("Knownports main loop: Exit")
}

func LoadPortsData()(data map[string]map[string]string, err error){
    values,err := ndb.LoadPortsData()
    if err != nil {logs.Error("knownports LoadPortsData -- Can't read query result: %s", err.Error()); return nil,err}
    return values, nil
}

func CheckParamKnownports(param string)(data string, err error){
    values, err := ndb.CheckParamKnownports(param)
    if err != nil  {logs.Error("knownports CheckParamKnownports -- Can't read query result: %s", err.Error()); return "",err} 
    return values, nil
}

func InsertknownportsElements(uuid string, param string, value string)(err error){
    err = ndb.InsertKnownports(uuid, param, value)
    if err != nil{logs.Error("Error InsertknownportsElements: "+err.Error()); return err}
    return nil
}