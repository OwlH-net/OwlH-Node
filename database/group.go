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
    Groupdb *sql.DB
)

func GConn() {
    var err error
    loadDataSQL := map[string]map[string]string{}
    loadDataSQL["groupConn"] = map[string]string{}
    loadDataSQL["groupConn"]["path"] = ""
    loadDataSQL["groupConn"]["cmd"] = "" 
    loadDataSQL, err = utils.GetConf(loadDataSQL)
    path := loadDataSQL["groupConn"]["path"]
    cmd := loadDataSQL["groupConn"]["cmd"]
    if err != nil {
        logs.Error("Gconn Error getting data from main.conf")
    }
    _, err = os.Stat(path) 
    if err != nil {
        panic("Fail opening group.db from path: "+path+"  --  "+err.Error())
    }    
    Groupdb, err = sql.Open(cmd,path)
    if err != nil {
        logs.Error("Groupdb/group -- group.db Open Failed: "+err.Error())
    }else {
        logs.Info("Groupdb/group -- group.db -> sql.Open, group.db Ready") 
    }
}

func InsertSuricataGroupValue(uuid string, param string, value string)(err error){
    if Groupdb == nil {logs.Error("no access to database group");return errors.New("no access to database group")}
    
    stmt, err := Groupdb.Prepare("insert into suricata (suri_uniqueid, suri_param, suri_value) values (?,?,?);")
    if err != nil {logs.Error("InsertSuricataGroupValue Prepare error: %s", err.Error());return err}
    
    _, err = stmt.Exec(&uuid, &param, &value)
    if err != nil {logs.Error("InsertSuricataGroupValue Execute error: %s", err.Error());return err}
    
    return nil
}

func UpdateSuricataGroupValue(uuid string, param string, value string)(err error){
    if Groupdb == nil {logs.Error("no access to database group");return errors.New("no access to database group")}
    
    gr, err := Groupdb.Prepare("update suricata set suri_value=? where suri_uniqueid = ? and suri_param = ?;")
    if err != nil {logs.Error("updateSuricataGroupValue Prepare error: %s", err.Error());return err}
    
    _, err = gr.Exec(&value, &uuid, &param)
    defer gr.Close()
    if err != nil {logs.Error("updateSuricataGroupValue Execute error: %s", err.Error());return err}
    
    return nil
}

func GetAllGroupData()(data map[string]map[string]string, err error){
    var pingData = map[string]map[string]string{}
    var uniqid string
    var param string
    var value string

    sql := "select suri_uniqueid, suri_param, suri_value from suricata;";
    
    rows, err := Groupdb.Query(sql)
    if err != nil { logs.Error("GetAllGroupData Groupdb.Query Error : %s", err.Error()); return nil, err}

    defer rows.Close()
    for rows.Next() {
        if err = rows.Scan(&uniqid, &param, &value); err != nil { logs.Error("GetAllGroupData -- Query return error: %s", err.Error()); return nil, err}

        if pingData[uniqid] == nil { pingData[uniqid] = map[string]string{}}
        pingData[uniqid][param]=value
    } 
    return pingData,nil
}