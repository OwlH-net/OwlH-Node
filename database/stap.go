package ndb

import (
    "github.com/astaxie/beego/logs"
    "database/sql"
    "os"
    "errors"
    "owlhnode/utils"
    _ "github.com/mattn/go-sqlite3"
)

var (
    Sdb *sql.DB
)

func SConn() {
    var err error
    loadDataSQL := map[string]map[string]string{}
    loadDataSQL["stapConn"] = map[string]string{}
    loadDataSQL["stapConn"]["path"] = ""
    loadDataSQL["stapConn"]["cmd"] = "" 
    loadDataSQL, err = utils.GetConf(loadDataSQL)
    path := loadDataSQL["stapConn"]["path"]
    cmd := loadDataSQL["stapConn"]["cmd"]
    if err != nil {
        logs.Error("Sconn Error getting data from main.conf")
    }
    _, err = os.Stat(path) 
    if err != nil {
        panic("Fail opening servers.db from path: "+path+"  --  "+err.Error())
    }    
    Sdb, err = sql.Open(cmd,path)
    if err != nil {
        logs.Error("sdb/stap -- servers.db Open Failed: "+err.Error())
    }else {
        logs.Info("sdb/stap -- servers.db -> sql.Open, servers.db Ready") 
    }
}

func GetStapServerInformation(uuid string)(serverData map[string]string, err error){
    var param string
    var value string
    stapServer := make(map[string]string)
    logs.Info("Creating data map for uuid: "+uuid)
    uuidParams, err := Sdb.Query("select server_param,server_value from servers where server_uniqueid = \""+uuid+"\";")
    defer uuidParams.Close()
    for uuidParams.Next(){
        if err = uuidParams.Scan(&param, &value); err!=nil {
            logs.Error("Error creating data Map: "+err.Error())
            return nil, err
        }
        stapServer[param]=value
    }
    return stapServer, nil
}

func AddServer(uuid string, param string, value string)(err error){
    if Sdb == nil {
        logs.Error("AddServer stap -- Can't access to database")
        return errors.New("AddServer stap -- Can't access to database")
    } 
    insertServerIP, err := Sdb.Prepare("insert into servers (server_uniqueid, server_param, server_value) values (?,?,?);")
    _, err = insertServerIP.Exec(&uuid, &param, &value)  
    defer insertServerIP.Close()
    if err != nil {
        logs.Error("ERROR stap/AddServer INSERT new server: "+err.Error())
        return err
    }  
    return nil
}

func GetAllServers()(data map[string]map[string]string, err error){
    var uniqueid string
    var param string
    var value string
    var allservers = map[string]map[string]string{}

    //database connection
    if Sdb == nil {
        logs.Error("GetAllServers stap -- Can't access to database")
        return nil,errors.New("GetAllServers stap -- Can't access to database")
    } 
    //query and make map[]map[]
    sql := "select server_uniqueid, server_param, server_value from servers;"
    rows, err := Sdb.Query(sql)
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
    if Sdb == nil {
        logs.Error("GetAllServers stap -- Can't access to database")
        return nil,errors.New("GetAllServers stap -- Can't access to database")
    } 
    //query and make map[]map[]
    sql := "select server_uniqueid, server_param, server_value from servers where server_uniqueid = '"+serveruuid+"';"
    rows, err := Sdb.Query(sql)
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

func UpdateStapByparam(param string, value string)(err error){
    if Sdb == nil {
        logs.Error("UpdateStapByparam stap -- Can't access to database")
        return errors.New("UpdateStapByparam stap -- Can't access to database")
    } 
    updateStapNode, err := Pdb.Prepare("update stap set stap_value = ? where stap_param = ?;")
    if (err != nil){ logs.Error("UpdateStapByparam UPDATE prepare error: "+err.Error()); return err}

    _, err = updateStapNode.Exec(&value, &param)
    if (err != nil){ logs.Error("UpdateStapByparam UPDATE exec error: "+err.Error()); return err}

    defer updateStapNode.Close()
    
    return nil
}

func UpdateStap(uuid string, param string, value string)(err error){
    if Sdb == nil {
        logs.Error("UpdateStapByparam stap -- Can't access to database")
        return errors.New("UpdateStapByparam stap -- Can't access to database")
    } 
    updateStapNode, err := Pdb.Prepare("update servers set server_value = ? where server_uniqueid = ? and server_param = ?;")
    if (err != nil){ logs.Error("UpdateStapByparam UPDATE prepare error: "+err.Error()); return err}

    _, err = updateStapNode.Exec(&value, &uuid, &param)
    if (err != nil){ logs.Error("UpdateStapByparam UPDATE exec error: "+err.Error()); return err}

    defer updateStapNode.Close()
    
    return nil
}

func GetStapUUID()(uuid string){
    //database connection
    if Sdb == nil {
        logs.Error("UUID stap -- Can't access to database")
        return ""
    } 
    
    var res string
    sql := "select stap_uniqueid from stap;"
    rows, err := Sdb.Query(sql)
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

func DeleteServer(uuid string)(err error){
    DeleteServiceNode, err := Pdb.Prepare("delete from servers where server_uniqueid = ?;")
    if (err != nil){ logs.Error("DeleteService UPDATE prepare error: "+err.Error()); return err}

    _, err = DeleteServiceNode.Exec(&uuid)
    if (err != nil){ logs.Error("DeleteService exec error: "+err.Error()); return err}

    defer DeleteServiceNode.Close()
    
    return nil
}

func PingServerStap(server string) (isIt map[string]string, err error){
    //database connection
    if Sdb == nil {
        logs.Error("PingServerStap stap -- Can't access to database")
        return nil, errors.New("PingServerStap stap -- Can't access to database")
    } 
    var res string
    stap := make(map[string]string)
    stap["stapStatus"] = "false"
    
    sql := "select server_value from servers where server_uniqueid = \""+server+"\" and server_param = \"status\";"
    logs.Info("PingServerStap select for check if exist query sql %s",sql)
    rows, err := Sdb.Query(sql)
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