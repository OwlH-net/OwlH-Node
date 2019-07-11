package knownports

import (
	"github.com/astaxie/beego/logs"
	"owlhnode/database"
	"owlhnode/utils"
	"github.com/hpcloud/tail"
	"errors"
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
    Data      Data `json:"data"`
}

type Data struct {
    Dstport     string `json:"dstport"`
	Proto       string `json:"proto"`
	Times		int `json:"times"`
    Signature   Signature `json:"alert"`
}

type Signature struct {
    Signature       string `json:"signature"`
    Signature_id    string `json:"signature_id"`
}

type LastAlert struct {
    PortProto       string 
    Last		    time.Time
    Counter		    int
}


func Init(){
	logs.Debug("-------------------------INIT KNOWNPORTS-------------------------")
	go NewPorts()
}

func GetStatus()(){
	for {
		_, err := CheckParamKnownports("status")
		_, err = CheckParamKnownports("mode")
		if err != nil {
			logs.Error("CheckParamKnownports Error: "+err.Error())
		}
		time.Sleep(time.Second * 20)
	}
}

func NewPorts()(){
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
			ChangeStatus(anode)
		}

		t, err := tail.TailFile(file, tail.Config{Follow: true, Location: &tail.SeekInfo{0, 2}})
		if err != nil {
			logs.Error("TailFile Error: %s", err.Error())
			Status = "Disabled"
			break
		}
		
		portsData, err := LoadPortsData()

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
		for line := range t.Lines {
			Status, err = CheckParamKnownports("status")
			Mode, err = CheckParamKnownports("mode")
			if err != nil {
				logs.Error("CheckParamKnownports Error: "+err.Error())
			}
			if Status != "Enabled"{
				break
			}
			var protoportRegexp = regexp.MustCompile(`"id.resp_h":"(\d+\.\d+\.\d+\.\d+)","id.resp_p":(\d+),"proto":"(\w+)"`)
			portProtocol := protoportRegexp.FindStringSubmatch(line.Text)
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
				notPortprotLearn := false
				for x := range portsData {

					if portsData[x]["portprot"] == protoport{
						// logs.Warn(portsData[x]["portprot"]+"     /--------/     "+protoport+" -------------------> UDAPTE")
						timeNow := time.Now() 
						value := strconv.FormatInt(timeNow.Unix(), 10)
						notPortprotLearn = true
						protoportUpdate, err := ndb.Pdb.Prepare("update knownports set kp_value = ? where kp_param = ? and kp_uniqueid = ?")
						_, err = protoportUpdate.Exec(&value, "last", &x)
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
					logs.Notice(portsData[uuid]["portprot"]+"     /--------/     "+protoport+" -------------------> INSERT")
					
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
				notPortprotProd := false
				for x := range portsData { 
					if portsData[x]["portprot"] == protoport{
						notPortprotProd = true
						t := time.Now()
						value := strconv.FormatInt(t.Unix(), 10)
						protoportUpdate, err := ndb.Pdb.Prepare("update knownports set kp_value = ? where kp_param = ? and kp_uniqueid = ?")
						_, err = protoportUpdate.Exec(&value, "last", &x)
						if err != nil {
							logs.Error("PRODUCTION MODE --> knownports[last] update error-> %s", err.Error())
							Status = "Disabled"
						}
						break
					}
				}
				if !notPortprotProd {
					// logs.Debug("MODE PRODUCTION: port and port do NOT exist into DB. Port/Protocol: "+protoport)				

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
							// logs.Notice("create alert - " +protoport + "/"+strconv.Itoa(alerted.Counter))
							createAlert = true
							alerted.Last = time.Now()
							alerted.Counter = 0
						} else {
							alerted.Counter += 1
							// logs.Debug("do not alert yet - "+protoport+ "/"+strconv.Itoa(alerted.Counter))
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
		t.Cleanup()
		t.Stop()		
		defer t.file.Close()
	}
	logs.Info("Knownports main loop: Exit")
}

func LoadPortsData()(data map[string]map[string]string, err error){
	var uniqueid string
	var param string
	var value string
	var allKnownPorts = map[string]map[string]string{}

	//database connection
	if ndb.Pdb == nil {
        logs.Error("LoadPorts knownports -- Can't access to database")
        return nil,errors.New("LoadPorts knownports -- Can't access to database")
	} 
	//query and make map[]map[]
	sql := "select kp_uniqueid, kp_param, kp_value from knownports;"
	rows, err := ndb.Pdb.Query(sql)
	defer rows.Close()
    if err != nil {
        logs.Error("LoadPorts knownports Error executing query: %s", err.Error())
        return nil, err
    }
	for rows.Next() {
        if err = rows.Scan(&uniqueid, &param, &value); err != nil {
            logs.Error("LoadPorts knownports -- Can't read query result: %s", err.Error())
            return nil, err
        }
        if allKnownPorts[uniqueid] == nil { allKnownPorts[uniqueid] = map[string]string{}}
        allKnownPorts[uniqueid][param]=value
	} 
	return allKnownPorts, nil
}

func CheckParamKnownports(param string)(data string, err error){
	var res string
	sql := "select plugin_value from plugins where plugin_uniqueid = 'knownports' and plugin_param='"+param+"'"
	rows, err := ndb.Pdb.Query(sql)
	defer rows.Close()
	if err != nil {
		logs.Error("knownports CheckParamKnownports Error executing query: %s", err.Error())
		return "",err
	}
	if rows.Next() {
		if err = rows.Scan(&res); err != nil {
			logs.Error("knownports CheckParamKnownports -- Can't read query result: %s", err.Error())
			return "",err
		}
	} 
	return res, nil
}

func InsertknownportsElements(uuid string, param string, value string)(err error){
	insertKP, err := ndb.Pdb.Prepare("insert into knownports (kp_uniqueid, kp_param, kp_value) values (?,?,?);")
	_, err = insertKP.Exec(&uuid, &param, &value)  
	defer insertKP.Close()
	if err != nil{
		logs.Error("Error InsertknownportsElements: "+err.Error())
		return err
	}
	return nil
}