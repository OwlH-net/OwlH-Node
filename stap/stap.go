package stap

import (
    "github.com/astaxie/beego/logs"
    "owlhnode/utils"
    "owlhnode/database"
    "errors"
    "encoding/json"
    "sync"
    "strconv"
)

var waitGroup sync.WaitGroup
var status bool
// var data chan string


func AddServer(elem map[string]string) (err error){
    logs.Warn(elem)
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
    if err != nil {
        logs.Error("ERROR stap/AddServer INSERT new server ip: "+err.Error())
        return err
    }	

    //insert name into server database
    logs.Info("stap/AddServer INSERT new server name")
    insertServerName, err := ndb.Sdb.Prepare("insert into servers (server_uniqueid, server_param, server_value) values (?,?,?);")
    _, err = insertServerName.Exec(&uuidServer, "name", &nodeName)  
    defer insertServerName.Close()
    if err != nil {
        logs.Error("ERROR stap/AddServer INSERT new server name: "+err.Error())
        return err
    }	
    
    //Load default server data from main.conf
    jsonDefaultData,err := utils.LoadDefaultServerData("software TAP PULL mode conf")
    logs.Info("File readed !!!")
    logs.Info(jsonDefaultData)
    jsonData := make(map[string]string)

    //parse raw data into byte array
    jsonDataArray := []byte(jsonDefaultData["fileContent"])
    err = json.Unmarshal(jsonDataArray, &jsonData)
    if err != nil {
        logs.Error("ERROR unmarshal json on AddServer func: "+err.Error())
        return err
    }		

    //insert default data into server database
    for z,v := range jsonData{
        logs.Warn(z+"   ---   "+v)
        insertServerName, err := ndb.Sdb.Prepare("insert into servers (server_uniqueid, server_param, server_value) values (?,?,?);")
        _, err = insertServerName.Exec(&uuidServer, &z, &v)  
        defer insertServerName.Close()
        if err != nil {
            logs.Error("ERROR INSERT servers: "+err.Error())
            return err
        }
    }

    return nil
}

func GetAllServers()(data map[string]map[string]string, err error){
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
    defer rows.Close()
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
    return allservers, nil
} 

func GetServer(serveruuid string)(data map[string]map[string]string, err error){
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
    defer rows.Close()
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
    } 
    return allservers, nil
} 



func PingStap(uuid string) (isIt map[string]bool, err error){
    //database connection
    if ndb.Sdb == nil {
        logs.Error("PingStap stap -- Can't access to database")
        return nil, errors.New("PingStap stap -- Can't access to database")
    } 
    var res string
    stap := make(map[string]bool)
    stap["stapStatus"] = false
    
    sql := "select stap_value from stap where stap_param = \"status\";"
    rows, err := ndb.Sdb.Query(sql)
    if err != nil {
        logs.Info("PingStap Query Error immediately after retrieve data %s",err.Error())
        return stap,err
    }
    defer rows.Close()
    if rows.Next() {
        err := rows.Scan(&res)
        if err != nil {
            logs.Info("Stap query error %s",err.Error())
            return stap, err
        }
        if res=="true"{
            stap["stapStatus"]=true
        }else{
            stap["stapStatus"]=false
        }
    }else if uuid != "" {
        logs.Info("Put Stap status INSERT")
        insertStap, err := ndb.Sdb.Prepare("insert into stap (stap_uniqueid, stap_param, stap_value) values (?,?,?);")
        _, err = insertStap.Exec(&uuid, "status", "false")  
        defer insertStap.Close()
        if err != nil{
            logs.Info("Error Insert uuid !=")
            return stap,err
        }
    }
    logs.Debug("checking stap value "+strconv.FormatBool(stap["stapStatus"]))
    return stap,nil
}


//Run stap main server on node
func RunStap(uuid string)(data string, err error){
    //database connection
    if ndb.Sdb == nil {
        logs.Error("RunStap stap -- Can't access to database")
        return "",errors.New("RunStap stap -- Can't access to database")
    } 
    logs.Info("Starting Stap...  "+uuid)
    updateStap, err := ndb.Sdb.Prepare("update stap set stap_value = ? where stap_param = ?;")
    _, err = updateStap.Exec("true", "status")  
    if (err != nil){
        logs.Error("Error updating RunStap: "+err.Error())
        return "", err
    }
    defer updateStap.Close()

    //call Concurrency StapInit
    StapInit()

    return "Stap is Running!", err
}

//Stop Stap
func StopStap(uuid string)(data string, err error){
    //database connection
    if ndb.Sdb == nil {
        logs.Error("StopStap stap -- Can't access to database")
        return "",errors.New("StopStap stap -- Can't access to database")
    } 
    logs.Info("Stopping Stap..."+uuid)
    updateStap, err := ndb.Sdb.Prepare("update stap set stap_value = ? where stap_param = ?;")
    if (err != nil){
        logs.Error ("Stopping Stap Update Prepare Error -> "+ err.Error())
        return "", err
    }
    
    _, err = updateStap.Exec("false", "status")  
    if (err != nil){
        logs.Error ("Stopping Stap Update Error -> "+ err.Error())
        return "", err
    }
    defer updateStap.Close()
    status = false
    return "Stap is stopped", err
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
        return "", errors.New("StopStapServer stap -- Can't access to database")
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
//Error Stap specific server
func ErrorStapServer(serveruuid string)(data string, err error){
    //database connection
    if ndb.Sdb == nil {
        logs.Error("ErrorStapServer stap -- Can't access to database")
        return "", errors.New("ErrorStapServer stap -- Can't access to database")
    } 
    logs.Info("Stopping specific Stap server...")
    updateErrorStapServer, err := ndb.Sdb.Prepare("update servers set server_value = ? where server_uniqueid = ? and server_param = ?;")
    _, err = updateErrorStapServer.Exec("error", &serveruuid, "status")  
    defer updateErrorStapServer.Close()
    if (err != nil){
        logs.Error("ErrorStapServer error updating: %s", err.Error())
        return "", err
    }
    return "Stap specific server status is error now", err
}

func PingServerStap(server string) (isIt map[string]string, err error){
    //database connection
    if ndb.Sdb == nil {
        logs.Error("PingServerStap stap -- Can't access to database")
        return nil, errors.New("PingServerStap stap -- Can't access to database")
    } 
    var res string
    stap := make(map[string]string)
    stap["stapStatus"] = "false"
    
    sql := "select server_value from servers where server_uniqueid = \""+server+"\" and server_param = \"status\";"
    logs.Info("PingServerStap select for check if exist query sql %s",sql)
    rows, err := ndb.Sdb.Query(sql)
    defer rows.Close()

    if err != nil {
        logs.Error("PingServerStap Query Error immediately after retrieve data %s",err.Error())
        return stap, err
    }
    if rows.Next() {
        err := rows.Scan(&res)
        if err != nil {
            logs.Error("Stap query error %s",err.Error())
            return stap, err
        }
        logs.Info("Stap status exists on stap DB --> Value: "+res)
        if res=="true"{
            stap["stapStatus"]="true"
        }else if res=="false"{
            stap["stapStatus"]="false"
        }else if res == "error"{
            stap["stapStatus"]="error"
        }else{
            stap["stapStatus"]="N/A"
        }
    }
    return stap, nil
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
    defer rows.Close()
    
    if err != nil {
        logs.Info("GetStapUUID Error immediately after retrieve data %s",err.Error())
        return ""
    }
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

//delete specific stap server
func DeleteStapServer(serveruuid string)(data string, err error){
    if ndb.Sdb == nil {
        logs.Error("DeleteStapServer stap -- Can't access to database")
        return "",errors.New("DeleteStapServer stap -- Can't access to database")
    } 
    sql := "delete from servers where server_uniqueid = ?;"
    deleteStartStapServer, err := ndb.Sdb.Prepare(sql)
    _, err = deleteStartStapServer.Exec(&serveruuid)  
    if (err != nil){
        logs.Error("DeleteStapServer error deleting: %s", err.Error())
        return "", err
    }
    defer deleteStartStapServer.Close()
    return "Stap specific server is deleted!", err
}

//Edit params for stap servers
func EditStapServer(data map[string]string) (err error){
    param := data["param"]
    value := data["value"]
    server := data["server"]
    if ndb.Sdb == nil {
        logs.Error("EditStapServer stap -- Can't access to database")
        return errors.New("EditStapServer stap -- Can't access to database")
    } 
    sql := "update servers set server_value = ? where server_uniqueid = ? and server_param = ?;"
    editStapServer, err := ndb.Sdb.Prepare(sql)
    _, err = editStapServer.Exec(&value, &server, &param)  
    defer editStapServer.Close()
    if (err != nil){
        logs.Error("EditStapServer error updating: %s", err.Error())
        return err
    }
    return nil
}