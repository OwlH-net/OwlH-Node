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
    Monitordb *sql.DB
)

func MConn() {
    var err error
    loadDataSQL := map[string]map[string]string{}
    loadDataSQL["monitorConn"] = map[string]string{}
    loadDataSQL["monitorConn"]["path"] = ""
    loadDataSQL["monitorConn"]["cmd"] = "" 
    loadDataSQL, err = utils.GetConf(loadDataSQL)    
    path := loadDataSQL["monitorConn"]["path"]
    cmd := loadDataSQL["monitorConn"]["cmd"]
    if err != nil {
        logs.Error("MConn Error getting data from main.conf")
    }
    _, err = os.Stat(path) 
    if err != nil {
        panic("Fail opening monitor.db from path: "+path+"  --  "+err.Error())
    }    
    Monitordb, err = sql.Open(cmd,path)
    if err != nil {
        logs.Error("Monitordb/database -- monitor.db Open Failed: "+err.Error())
    }else {
        logs.Info("Monitordb/database -- monitor.db -> sql.Open, monitor.db Ready") 
    }
}

func InsertMonitorValue(uuid string, param string, value string)(err error){
    if Monitordb == nil {logs.Error("no access to database monitor");return errors.New("no access to database monitor")}
    
    stmt, err := Monitordb.Prepare("insert into files (file_uniqueid, file_param, file_value) values (?,?,?);")
    if err != nil {logs.Error("InsertMonitorValue Prepare error: %s", err.Error());return err}
    
    _, err = stmt.Exec(&uuid, &param, &value)
    if err != nil {logs.Error("InsertMonitorValue Execute error: %s", err.Error());return err}
    
    return nil
}

func LoadMonitorFiles()(data map[string]map[string]string, err error){
    var pingData = map[string]map[string]string{}
    var uniqid string
    var param string
    var value string

    sql := "select file_uniqueid, file_param, file_value from files;";
    
    rows, err := Monitordb.Query(sql)
    if err != nil {
        logs.Error("LoadDataflowValues Monitordb.Query Error : %s", err.Error())
        return nil, err
    }
    defer rows.Close()
    for rows.Next() {
        if err = rows.Scan(&uniqid, &param, &value); err != nil {
            logs.Error("LoadDataflowValues -- Query return error: %s", err.Error())
            return nil, err
        }
        if pingData[uniqid] == nil { pingData[uniqid] = map[string]string{}}
        pingData[uniqid][param]=value
    } 
    return pingData,nil
}

func DeleteMonitorFile(uuid string)(err error){
    deleteFile, err := Monitordb.Prepare("delete from files where file_uniqueid = ?;")
    _, err = deleteFile.Exec(&uuid)
    defer deleteFile.Close()
    if err != nil {logs.Error("DeleteMonitorFile ERROR deleting: "+err.Error());return err}
    return nil
}