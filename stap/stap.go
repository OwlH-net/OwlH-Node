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