package knownports

import (
	"github.com/astaxie/beego/logs"
	"owlhnode/database"
	"owlhnode/utils"
	"github.com/hpcloud/tail"
	"errors"
	"regexp"
	"strconv"
	"time"
)

func Init()(){
	// res, err := CheckParamKnownports()
	// if res = "Enabled"{
	go NewPorts()		
	// }
}

func NewPorts()(){
	var err error
	flag := true
	loadPorts := map[string]map[string]string{}
	loadPorts["knownports"] = map[string]string{}
    loadPorts["knownports"]["file"] = ""
	loadPorts,err = utils.GetConf(loadPorts)
    file := loadPorts["knownports"]["file"]
	if err != nil {
		logs.Error("loadPorts Error getting data from main.conf: "+err.Error())
		flag = false
	}

	
	for flag == true {
		status, err := CheckParamKnownports("status")
		mode, err := CheckParamKnownports("mode")
		if err != nil {
			logs.Error("CheckParamKnownports Error: "+err.Error())
			flag = false
			continue
		}
		if status == "Enabled"{
			t, err := tail.TailFile(file, tail.Config{Follow: true})
			if err != nil {
				logs.Error("TailFile Error: %s", err.Error())
				flag = false
			}
			for line := range t.Lines {
				portsData, err := LoadPortsData()
				logs.Debug(portsData)
				if err != nil {
					logs.Error("LoadPortsData NewPorts Error: %s", err.Error())
					flag = false
				}
				var protoportRegexp = regexp.MustCompile(`"id.resp_p":(\d+),"proto":"(\w+)"`)
				portProtocol := protoportRegexp.FindStringSubmatch(line.Text)
				var protoport = portProtocol[1]+"/"+portProtocol[2]
				if mode == "Learning"{
					notPortprotLearn := false
					for x := range portsData {
						if portsData[x]["portprot"] == protoport{
							logs.Warn(portsData[x]["portprot"]+"     /--------/     "+protoport+" -------------------> UDAPTE")
							t := time.Now() 
							value := strconv.FormatInt(t.Unix(), 10)
							notPortprotLearn = true
							protoportUpdate, err := ndb.Pdb.Prepare("update knownports set kp_value = ? where kp_param = ? and kp_uniqueid = ?")
							_, err = protoportUpdate.Exec(&value, "last", &x)
							if err != nil {
								logs.Error("LEARNING MODE --> knownports[last] update error-> %s", err.Error())
								flag = false
							}
						}
					}
					if !notPortprotLearn{
						uuid := utils.Generate()
						t := time.Now() 
						value := strconv.FormatInt(t.Unix(), 10)
						logs.Notice(portsData[uuid]["portprot"]+"     /--------/     "+protoport+" -------------------> INSERT")
						err = InsertknownportsElements(uuid, "port", portProtocol[1])
						err = InsertknownportsElements(uuid, "protocol", portProtocol[2])
						err = InsertknownportsElements(uuid, "portprot", portProtocol[1]+"/"+portProtocol[2])
						err = InsertknownportsElements(uuid, "first", value)
						err = InsertknownportsElements(uuid, "last", value)
						if err != nil {
							logs.Error("knownports insert error-> %s", err.Error())
							flag = false
						}
					}
				}else{
					notPortprotProd := false
					for x := range portsData { 
						if portsData[x]["portprot"] == protoport{
							logs.Critical("MODE PRODUCTION: portprot NOT exist into DB...")
							notPortprotProd = true
							t := time.Now()
							value := strconv.FormatInt(t.Unix(), 10)
							protoportUpdate, err := ndb.Pdb.Prepare("update knownports set kp_value = ? where kp_param = ? and kp_uniqueid = ?")
							_, err = protoportUpdate.Exec(&value, "last", &x)
							if err != nil {
								logs.Error("PRODUCTION MODE --> knownports[last] update error-> %s", err.Error())
								flag = false
							}
						}
					}
					if !notPortprotProd {
						logs.Debug("MODE PRODUCTION: portprot NOT exist into DB...")
					}
				}
			}
		}else {
			flag = false
		}
	}
	logs.Error("NewPorts: Concurrent loop error")
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
	sql := "select plugin_value from plugins where plugin_uniqueid = '0000-00-00-00-000000' and plugin_param='"+param+"'"
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
            logs.Info("Error Insert uuid !=")
            return err
		}
		return nil
}