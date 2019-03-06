package stap

import (
    "github.com/astaxie/beego/logs"
    // "os"
    // "os/exec"
    // "strings"
    // "regexp"
	"owlhnode/utils"
	"owlhnode/database"
	// "io/ioutil"
	"errors"
	"encoding/json"
)

func AddServer(elem map[string]string) (err error){
	nodeName:= elem["nodeName"]
	nodeIP:= elem["nodeIP"]
	uuidServer := utils.Generate()
	
	//database connection
	if ndb.Sdb == nil {
        logs.Error("AddServer stap -- Can't access to database")
        return errors.New("AddServer stap -- Can't access to database")
	} 
	
	//insert ip into server database
	logs.Info("stap/AddServer INSERT new server ip")
	insertServerIP, err := ndb.Sdb.Prepare("insert into servers (server_uniqueid, server_param, server_value) values (?,?,?);")
	_, err = insertServerIP.Exec(&uuidServer, "ip", &nodeIP)  
	defer insertServerIP.Close()
	if err != nil {return err}	

	//insert name into server database
	logs.Info("stap/AddServer INSERT new server name")
	insertServerName, err := ndb.Sdb.Prepare("insert into servers (server_uniqueid, server_param, server_value) values (?,?,?);")
	_, err = insertServerName.Exec(&uuidServer, "name", &nodeName)  
	defer insertServerName.Close()
	if err != nil {return err}	
	
	//Load default server data from main.conf
	jsonDefaultData,err := utils.LoadDefaultServerData("defaults.json")
	logs.Info("Fichero Leido!")
	logs.Info(jsonDefaultData)
	jsonData := make(map[string]string)
	//parse raw data into byte array
	jsonDataArray := []byte(jsonDefaultData["fileContent"])
	err = json.Unmarshal(jsonDataArray, &jsonData)
	if err != nil {return err}	
	//insert default data into server database
	for z,v := range jsonData{
		logs.Warn("Key--> "+z+" Value..> "+v)
		insertServerName, err := ndb.Sdb.Prepare("insert into servers (server_uniqueid, server_param, server_value) values (?,?,?);")
		_, err = insertServerName.Exec(&uuidServer, &z, &v)  
		defer insertServerName.Close()
		if err != nil {
			return err
		}
	}

	return nil
}

func GetAllServers()(data *map[string]map[string]string, err error){
	var uniqueid string
	var param string
	var value string
	var allservers = map[string]map[string]string{}

	//database connection
	if ndb.Sdb == nil {
        logs.Error("GetAllServers stap -- Can't access to database")
        return nil,errors.New("GetAllServers stap -- Can't access to database")
	} 
	//query and make map[]map[]
	sql := "select server_uniqueid, server_param, server_value from servers;"
    rows, err := ndb.Sdb.Query(sql)
    if err != nil {
        logs.Error("GetAllServers stap Error al ejecutar la query: %s", err.Error())
        return nil, err
    }
	for rows.Next() {
        if err = rows.Scan(&uniqueid, &param, &value); err != nil {
            logs.Error("GetAllServers stap -- No hemos podido leer del resultado de la query: %s", err.Error())
            return nil, err
        }
        if allservers[uniqueid] == nil { allservers[uniqueid] = map[string]string{}}
        allservers[uniqueid][param]=value
	} 
	return &allservers, nil
} 

func GetServer(serveruuid string)(data *map[string]map[string]string, err error){
	var uniqueid string
	var param string
	var value string
	var allservers = map[string]map[string]string{}

	//database connection
	if ndb.Sdb == nil {
        logs.Error("GetAllServers stap -- Can't access to database")
        return nil,errors.New("GetAllServers stap -- Can't access to database")
	} 
	//query and make map[]map[]
	sql := "select server_uniqueid, server_param, server_value from servers where server_uniqueid = '"+serveruuid+"';"
    rows, err := ndb.Sdb.Query(sql)
    if err != nil {
        logs.Error("GetAllServers stap Error al ejecutar la query: %s", err.Error())
        return nil, err
    }
	for rows.Next() {
        if err = rows.Scan(&uniqueid, &param, &value); err != nil {
            logs.Error("GetAllServers stap -- No hemos podido leer del resultado de la query: %s", err.Error())
            return nil, err
        }
        if allservers[uniqueid] == nil { allservers[uniqueid] = map[string]string{}}
		allservers[uniqueid][param]=value
		
		logs.Info(param+"   ----   "+value)
	} 
	return &allservers, nil
} 



func PingStap() (isIt map[string]bool){
    stap := make(map[string]bool)
    //stap = false
	stap["stapStatus"] = ""
	
	sql := "select node_value from stap where node_uniqueid = \""+uuid+"\" and node_param = \"status\";"
    logs.Info("Stap select for check if exist query sql %s",sql)
    rows, err := ndb.Db.Query(sql)
    if err != nil {
        logs.Info("Stap query error %s",err.Error())
        return "", err
    }
	defer rows.Close()
    if rows.Next() {
        //rows.Close()
        logs.Info("Stap UPDATE")
        updtbpf, err := ndb.Db.Prepare("update stap set node_value = ? where node_uniqueid = ? and node_param = ?;")

        if (err != nil){
            logs.Info("Put BPF Suricata prepare UPDATE -- "+err.Error())
            return "", err
        }
        _, err = updtbpf.Exec(&jsonbpf, &jsonnid, bpftext)  
        defer updtbpf.Close()      

        return "SuccessUpdate", err
    }else{
        logs.Info("Put BPF Suricata res INSERT")
        indtbpf, err := ndb.Db.Prepare("insert into nodes (node_uniqueid, node_param, node_value) values (?,?,?);")
        _, err = indtbpf.Exec(&jsonnid, bpftext, &jsonbpf)  
        defer indtbpf.Close()
        if (err != nil){
            return "", err
        }
        return "SuccessInsert", err
    }
    return "Error", errors.New("Put SuricataBPF -- There is no defined BPF")
	



    logs.Warn("stap --> ")
    logs.Warn(stap)
}


//Run stap
func RunStap()(data string, err error){

    // //Retrieve path for Stap.
    StartStap := map[string]map[string]string{}
    StartStap["stapStart"] = map[string]string{}
    StartStap["stapStart"]["start"] = ""
    StartStap["stapStart"]["param"] = ""
    StartStap["stapStart"]["command"] = ""
    StartStap = utils.GetConf(StartStap)    
    cmd := StartStap["stapStart"]["start"]
    param := StartStap["stapStart"]["param"]
    command := StartStap["stapStart"]["command"]

    out,err := exec.Command(command, param, cmd).Output()
    logs.Info(string(out))
    if err != nil {
        logs.Error("Error launching Stap: "+err.Error())
        return "",err
    }
    return "Stap system is on",nil
}

//Stop Stap
func StopStap()(data string, err error){

    // //Retrieve path for Stap.
    StopStap := map[string]map[string]string{}
	StopStap["stapStop"] = map[string]string{}
    StopStap["stapStop"]["stop"] = ""
    StopStap["stapStop"]["param"] = ""
    StopStap["stapStop"]["command"] = ""
    StopStap = utils.GetConf(StopStap)    
    cmd := StopStap["stapStop"]["stop"]
    param := StopStap["stapStop"]["param"]
    command := StopStap["stapStop"]["command"]

    _,err = exec.Command(command, param, cmd).Output()
    if err != nil {
        logs.Error("Error stopping Stap: "+err.Error())
        return "",err
    }
    return "Stap stopped ",nil
}