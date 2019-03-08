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
	logs.Info("File readed !!!")
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



func PingStap(uuid string) (isIt map[string]bool){
    //database connection
	if ndb.Sdb == nil {
        logs.Error("PingStap stap -- Can't access to database")
        return nil
    } 
    var res string
    stap := make(map[string]bool)
    //stap = false
	stap["stapStatus"] = false
	
	sql := "select stap_value from stap where stap_uniqueid = \""+uuid+"\" and stap_param = \"status\";"
    logs.Info("Stap select for check if exist query sql %s",sql)
    rows, err := ndb.Sdb.Query(sql)
    // rows, err := ndb.Db.Query(sql).Scan(&res)
    // logs.Info("Value returned query Stap-->"+res)
    if err != nil {
        logs.Info("Query Error immediately after retrieve data %s",err.Error())
        return stap
    }
    logs.Info("After rows Query")
	defer rows.Close()
    if rows.Next() {
        err := rows.Scan(&res)
        if err != nil {
            logs.Info("Stap query error %s",err.Error())
            return stap
        }
        logs.Info("Stap status exists on stap DB --> Value: "+res)
        if res=="true"{stap["stapStatus"]=true}else{stap["stapStatus"]=false}
        logs.Warn("Stap status-->")
        logs.Warn(stap)
        return stap
    }else{
        logs.Info("Put Stap status INSERT")
        insertStap, err := ndb.Sdb.Prepare("insert into stap (stap_uniqueid, stap_param, stap_value) values (?,?,?);")
        _, err = insertStap.Exec(&uuid, "status", "false")  
        defer insertStap.Close()
        if (err != nil){
            return stap
        }
        return stap
    }
}


//Run stap
func RunStap(uuid string)(data string, err error){
    //database connection
	if ndb.Sdb == nil {
        logs.Error("RunStap stap -- Can't access to database")
        return "",errors.New("RunStap stap -- Can't access to database")
	} 
    logs.Info("Starting Stap...")
    //insertStap, err := ndb.Sdb.Prepare("insert into stap (stap_uniqueid, stap_param, stap_value) values (?,?,?);")
    updateStap, err := ndb.Sdb.Prepare("update stap set stap_value = ? where stap_uniqueid = ? and stap_param = ?;")
    _, err = updateStap.Exec("true", &uuid, "status")  
    defer updateStap.Close()
    if (err != nil){
        return "", err
    }
    return "Stap is Running!", err
}

//Stop Stap
func StopStap(uuid string)(data string, err error){
    //database connection
	if ndb.Sdb == nil {
        logs.Error("StopStap stap -- Can't access to database")
        return "",errors.New("StopStap stap -- Can't access to database")
	} 
    logs.Info("Stopping Stap...")
    updateStap, err := ndb.Sdb.Prepare("update stap set stap_value = ? where stap_uniqueid = ? and stap_param = ?;")
    _, err = updateStap.Exec("false", &uuid, "status")  
    defer updateStap.Close()
    if (err != nil){
        return "", err
    }
    return "Stap is stoped", err
}

//Run stap specific server
func RunStapServer(serveruuid string)(data string, err error){
    //database connection
	if ndb.Sdb == nil {
        logs.Error("RunStapServer stap -- Can't access to database")
        return "",errors.New("RunStapServer stap -- Can't access to database")
    } 
    logs.Error("Server uuid "+serveruuid)
    logs.Info("Starting specific Stap server...")
    //insertStap, err := ndb.Sdb.Prepare("insert into stap (stap_uniqueid, stap_param, stap_value) values (?,?,?);")
    sql := "update servers set server_value = ? where server_uniqueid = ? and server_param = ?;"
    updateStartStapServer, err := ndb.Sdb.Prepare(sql)
    _, err = updateStartStapServer.Exec("true", &serveruuid, "status")  
    
    defer updateStartStapServer.Close()
    if (err != nil){
        logs.Error("RunStapServer error updating: %s", err.Error())
        return "", err
    }
    logs.Error("RunStapServer finishing function at Node")
    return "Stap specific server is Running!", err
}

//Stop Stap specific server
func StopStapServer(serveruuid string)(data string, err error){
    //database connection
	if ndb.Sdb == nil {
        logs.Error("StopStapServer stap -- Can't access to database")
        return "",errors.New("StopStapServer stap -- Can't access to database")
	} 
    logs.Info("Stopping specific Stap server...")
    updateStopStapServer, err := ndb.Sdb.Prepare("update servers set server_value = ? where server_uniqueid = ? and server_param = ?;")
    _, err = updateStopStapServer.Exec("false", &serveruuid, "status")  
    defer updateStopStapServer.Close()
    if (err != nil){
        logs.Error("StopStapServer error updating: %s", err.Error())
        return "", err
    }
    return "Stap specific server  is stoped", err
}

func PingServerStap(server string) (isIt map[string]bool){
    //database connection
	if ndb.Sdb == nil {
        logs.Error("PingServerStap stap -- Can't access to database")
        return nil
    } 
    var res string
    stap := make(map[string]bool)
	stap["stapStatus"] = false
	
	sql := "select server_value from servers where server_uniqueid = \""+server+"\" and server_param = \"status\";"
    logs.Info("PingServerStap select for check if exist query sql %s",sql)
    rows, err := ndb.Sdb.Query(sql)
    // rows, err := ndb.Db.Query(sql).Scan(&res)
    // logs.Info("Value returned query PingServerStap-->"+res)
    if err != nil {
        logs.Info("Query Error immediately after retrieve data %s",err.Error())
        return stap
    }
    logs.Info("After rows Query")
	defer rows.Close()
    if rows.Next() {
        err := rows.Scan(&res)
        if err != nil {
            logs.Info("Stap query error %s",err.Error())
            return stap
        }
        logs.Info("Stap status exists on stap DB --> Value: "+res)
        if res=="true"{stap["stapStatus"]=true}else{stap["stapStatus"]=false}
        logs.Warn("PingServerStap status-->")
        logs.Warn(stap)
        return stap
    }
    return stap
    // else{
    //     logs.Info("Put Stap status INSERT")
    //     insertStap, err := ndb.Sdb.Prepare("insert into stap (stap_uniqueid, stap_param, stap_value) values (?,?,?);")
    //     _, err = insertStap.Exec(&server, "status", "false")  
    //     defer insertStap.Close()
    //     if (err != nil){
    //         return stap
    //     }
    //     return stap
    // }
}