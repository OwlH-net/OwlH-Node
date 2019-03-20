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
    "sync"
    // "time"
    // "strconv"
)

var waitGroup sync.WaitGroup
var status bool
// var data chan string


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
		//logs.Warn("Key--> "+z+" Value..> "+v)
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
        logs.Error("GetAllServers stap Error executing query: %s", err.Error())
        return nil, err
    }
	for rows.Next() {
        if err = rows.Scan(&uniqueid, &param, &value); err != nil {
            logs.Error("GetAllServers stap -- Can't read query result: %s", err.Error())
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
        logs.Error("GetAllServers stap Error query execution: %s", err.Error())
        return nil, err
    }
	for rows.Next() {
        if err = rows.Scan(&uniqueid, &param, &value); err != nil {
            logs.Error("GetAllServers stap -- Can't read the query result: %s", err.Error())
            return nil, err
        }
        if allservers[uniqueid] == nil { allservers[uniqueid] = map[string]string{}}
		allservers[uniqueid][param]=value
		
		//logs.Info(param+"   ----   "+value)
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
	stap["stapStatus"] = false
	
	sql := "select stap_value from stap where stap_param = \"status\";"
    rows, err := ndb.Sdb.Query(sql)

    if err != nil {
        logs.Info("PingStap Query Error immediately after retrieve data %s",err.Error())
        return stap
    }
	defer rows.Close()
    if rows.Next() {
        err := rows.Scan(&res)
        if err != nil {
            logs.Info("Stap query error %s",err.Error())
            return stap
        }
        if res=="true"{
            stap["stapStatus"]=true
        }else{
            stap["stapStatus"]=false
        }

        logs.Info("Stap status-->")
        logs.Info(stap)
    }else if uuid != "" {
        logs.Info("Put Stap status INSERT")
        insertStap, err := ndb.Sdb.Prepare("insert into stap (stap_uniqueid, stap_param, stap_value) values (?,?,?);")
        _, err = insertStap.Exec(&uuid, "status", "false")  
        defer insertStap.Close()
        if (err != nil){
            logs.Info("Error Insert uuid !=")
            return stap
        }
    }
    return stap
}


//Run stap main server on node
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

    //call Concurrency StapInit
    StapInit()


    // //Start concurrency for stap servers.
    // status = true
    // var serverOnUUID string
    // for{
    //     rows, _ := ndb.Sdb.Query("select server_uniqueid from servers where server_param = \"status\" and server_value = \"true\";")
    //     for rows.Next(){
    //         rows.Scan(&serverOnUUID)
    //         waitGroup.Add(1)
    //         logs.Info("Worker added No:"+serverOnUUID)
    //         //data <- ("Worker "+serverOnUUID)
    //         go worker(serverOnUUID)
    //     }
    //     defer rows.Close()
    //     waitGroup.Wait()
    //     if (!status) {
    //         logs.Info("Close Concurrency cLoop")
    //         break
    //     }
    // }
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
    status = false
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
    if err != nil {
        logs.Error("PingServerStap Query Error immediately after retrieve data %s",err.Error())
        return stap
    }
    logs.Info("After rows Query")
	defer rows.Close()
    if rows.Next() {
        err := rows.Scan(&res)
        if err != nil {
            logs.Error("Stap query error %s",err.Error())
            return stap
        }
        logs.Info("Stap status exists on stap DB --> Value: "+res)
        if res=="true"{stap["stapStatus"]=true}else{stap["stapStatus"]=false}
        logs.Info("PingServerStap status-->")
        logs.Info(stap)
        return stap
    }
    return stap
}

func GetStapUUID()(uuid string){
    //database connection
	if ndb.Sdb == nil {
        logs.Error("UUID stap -- Can't access to database")
        return ""
    } 
    
    var res string
	sql := "select stap_uniqueid from stap;"
    rows, err := ndb.Sdb.Query(sql)
    if err != nil {
        logs.Info("GetStapUUID Error immediately after retrieve data %s",err.Error())
        return ""
    }
    defer rows.Close()
    if rows.Next() {
        err := rows.Scan(&res)
        if err != nil {
            logs.Info("Stap query error %s",err.Error())
            return ""
        }
        return res
    }
    return ""
}


// //Goroutine for concurrency with Stap servers
// func worker(uuid string){
//     logs.Info("Starting Worker")
//     defer func() {
// 		logs.Info("Destroying worker "+uuid)
// 		waitGroup.Done()
//     }()
//     for {
//         // value, err := <-data
//         logs.Warn("UUID: "+uuid+" --- Sleep for 1 second")
//         time.Sleep(time.Second * 3)
//         break
//     }    
    
// }