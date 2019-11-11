package configuration

import (
    // "encoding/json"
    // "strconv"
    "github.com/astaxie/beego/logs"
    "database/sql"
    // "io/ioutil"
    // "io"
    // "errors"
    "owlhnode/utils"
    "os"
    // "time"
    // "os/exec"
    // "fmt"
    // "crypto/rand"
    _ "github.com/mattn/go-sqlite3"
)

type Dbconfig struct {
    Dbname          string
    Dbconn          string
    Dbpath          string
    Dbtables        []Table
}

type Table struct {
    Tconn           string
    Tname           string
    Tcreate         string
    Tfields         []Field
}

type Field struct {
    Fname           string
    Finsert         string
}

var DBCONFIG        []Dbconfig


func MainCheck()(cancontinue bool){
    dbs := []string{"monitorConn","stapConn","pluginConn","nodeConn"}
    for db := range dbs {
        ok := CheckDB(dbs[db])
        if !ok {
            return false
        }
    }

    var table Table

    table.Tname = "plugins"
    table.Tconn = "pluginConn"
    table.Tcreate = "CREATE TABLE plugins (plugin_id integer PRIMARY KEY AUTOINCREMENT,plugin_uniqueid text NOT NULL,plugin_param text NOT NULL,plugin_value text NOT NULL)"

    ok := CheckTable(table)
    if !ok {
        return false
    }

    table.Tname = "knownports"
    table.Tconn = "pluginConn"
    table.Tcreate = "CREATE TABLE knownports (kp_id integer PRIMARY KEY AUTOINCREMENT,kp_uniqueid text NOT NULL,kp_param text NOT NULL,kp_value text NOT NULL)"
    ok = CheckTable(table)
    if !ok {
        return false
    }

    table.Tname = "analyzer"
    table.Tconn = "pluginConn"
    table.Tcreate = "CREATE TABLE analyzer (analyzer_id integer PRIMARY KEY AUTOINCREMENT,analyzer_uniqueid text NOT NULL,analyzer_param text NOT NULL,analyzer_value text NOT NULL)"
    ok = CheckTable(table)
    if !ok {
        return false
    }

    table.Tname = "mainconf"
    table.Tconn = "pluginConn"
    table.Tcreate = "CREATE TABLE mainconf (main_id integer PRIMARY KEY AUTOINCREMENT,main_uniqueid text NOT NULL,main_param text NOT NULL,main_value text NOT NULL)"
    ok = CheckTable(table)
    if !ok {
        return false
    }

    table.Tname = "cluster"
    table.Tconn = "pluginConn"
    table.Tcreate = "CREATE TABLE cluster (cluster_id integer PRIMARY KEY AUTOINCREMENT,cluster_uniqueid text NOT NULL,cluster_param text NOT NULL,cluster_value text NOT NULL)"
    ok = CheckTable(table)
    if !ok {
        return false
    }

    table.Tname = "files"
    table.Tconn = "monitorConn"
    table.Tcreate = "CREATE TABLE files (file_id integer PRIMARY KEY AUTOINCREMENT,file_uniqueid text NOT NULL,file_param text NOT NULL,file_value text NOT NULL)"
    ok = CheckTable(table)
    if !ok {
        return false
    }

    table.Tname = "nodeconfig"
    table.Tconn = "nodeConn"
    table.Tcreate = "CREATE TABLE nodeconfig (config_id integer PRIMARY KEY AUTOINCREMENT,config_uniqueid text NOT NULL,config_param text NOT NULL,config_value text NOT NULL)"
    ok = CheckTable(table)
    if !ok {
        return false
    }

    table.Tname = "dataflow"
    table.Tconn = "nodeConn"
    table.Tcreate = "CREATE TABLE dataflow (flow_id integer PRIMARY KEY AUTOINCREMENT,flow_uniqueid text NOT NULL,flow_param text NOT NULL,flow_value text NOT NULL)"
    ok = CheckTable(table)
    if !ok {
        return false
    }

    table.Tname = "changerecord"
    table.Tconn = "nodeConn"
    table.Tcreate = "CREATE TABLE changerecord (control_id integer PRIMARY KEY AUTOINCREMENT,control_uniqueid text NOT NULL,control_param text NOT NULL,control_value text NOT NULL);"
    ok = CheckTable(table)
    if !ok {
        return false
    }

    table.Tname = "node"
    table.Tconn = "nodeConn"
    table.Tcreate = "CREATE TABLE node (node_id integer PRIMARY KEY AUTOINCREMENT,node_uniqueid text NOT NULL,node_param text NOT NULL,node_value text NOT NULL);"
    ok = CheckTable(table)
    if !ok {
        return false
    }

    table.Tname = "incidents"
    table.Tconn = "nodeConn"
    table.Tcreate = "CREATE TABLE incidents (incidents_id integer PRIMARY KEY AUTOINCREMENT,incidents_uniqueid text NOT NULL,incidents_param text NOT NULL,incidents_value text NOT NULL);"
    ok = CheckTable(table)
    if !ok {
        return false
    }

    table.Tname = "servers"
    table.Tconn = "stapConn"
    table.Tcreate = "CREATE TABLE servers (server_id integer PRIMARY KEY AUTOINCREMENT,server_uniqueid text NOT NULL,server_param text NOT NULL,server_value text NOT NULL);"
    ok = CheckTable(table)
    if !ok {
        return false
    }

    table.Tname = "stap"
    table.Tconn = "stapConn"
    table.Tcreate = "CREATE TABLE stap (stap_id integer PRIMARY KEY AUTOINCREMENT,stap_uniqueid text NOT NULL,stap_param text NOT NULL,stap_value text NOT NULL);"
    ok = CheckTable(table)
    if !ok {
        return false
    }


    return true
}

func checkFields()(ok bool){
// plugins -             UUID - "suricata", param - "previousStatus", data["suricata"]["status"]
// plugins - "analyzer" - select analyzer_value from analyzer where analyzer_param='status'"
}

func CheckDB(conn string)(ok bool) {
    loadDataSQL := map[string]map[string]string{}
    loadDataSQL[conn] = map[string]string{}
    loadDataSQL[conn]["path"] = ""
    loadDataSQL, err := utils.GetConf(loadDataSQL)
    if err != nil {
        logs.Error("Configuration -> Can't get "+conn+" path from main.conf")
        return false
    }
    dbpath := loadDataSQL[conn]["path"]

    exists := DbExists(dbpath)

    if exists {
        logs.Warn("Configuration -> db "+dbpath+" exists")
        return true
    } else {
        logs.Warn("Configuration -> db "+dbpath+" does not exist, ... Creating")
        err = DbCreate(dbpath)
        if err != nil {
            return false
        }
    }
    return true
}

func CheckTable(table Table)(ok bool){
    loadDataSQL := map[string]map[string]string{}
    loadDataSQL[table.Tconn] = map[string]string{}
    loadDataSQL[table.Tconn]["path"] = ""
    loadDataSQL, err := utils.GetConf(loadDataSQL)
    if err != nil {
        logs.Error("Configuration -> Can't get "+table.Tconn+" path from main.conf")
        return false
    }
    dbpath := loadDataSQL[table.Tconn]["path"]

    exists := TableExists(dbpath, table.Tname)
    if !exists {
        logs.Warn("Configuration -> Table "+table.Tname+" doesn't exist on DB "+table.Tconn+" ...Creating")
        created := TableCreate(table.Tconn,table.Tname,table.Tcreate)
        if !created {
            return false
        }
        return true
    }

    logs.Info("Configuration -> Table "+table.Tname+" exists on DB "+table.Tconn)
    return true
}



func DbExists(db string)(exists bool){
    if _, err := os.Stat(db); os.IsNotExist(err) {
        logs.Error("Configuration -> Check DB -> db " + db + " not found -> err: " + err.Error())
        return false
    }else{
        dblink, err := sql.Open("sqlite3", db)
        if err != nil {
            logs.Error("Configuration -> Check DB -> db " + db + " can't be opened -> err: "+err.Error())
            return false
        }
        defer dblink.Close()
        return true
    }
    logs.Error("Configuration -> Check DB -> db " + db + " something went wrong, can't find a reason")
    return false
}

func TableExists(db string, table string)(exists bool){
    dblink, err := sql.Open("sqlite3", db)
    if err != nil {
        logs.Error("Configuration -> Check Table -> db " + db + " can't open -> err: "+err.Error())
        return false
    }
    defer dblink.Close()
    qry := "SELECT name FROM sqlite_master WHERE type='table' AND name=$1;"
    row := dblink.QueryRow(qry, table)

    var tablename string
    switch err := row.Scan(&tablename); err {
    case sql.ErrNoRows:
        return false
    case nil:
        return true
    default:
        return false
    }

    return true
}

func RecordExists(  db string, 
                    table string, 
                    uuid_field string, 
                    uuid_value string, 
                    param_field string, 
                    param_value string,
                    value_field string)(exists bool) {
    dblink, err := sql.Open("sqlite3", db)
    if err != nil {
        logs.Error("Configuration -> Check Record -> db " + db + " can't open -> err: "+err.Error())
        return false
    }
    defer dblink.Close()
    sql := "select "+value_field+" from "+table+" where "+uuid_field+"='"+uuid_value+"' AND "+param_field+"='"+param_value+"'";
    rows, err := dblink.Query(sql)
    if err != nil {
        logs.Error("Configuration -> Check Record -> "+db+" record "+table+"->"+uuid_value+"->"+param_value+" can't query -> err: "+err.Error())
        return false
    }
    defer rows.Close()
    for rows.Next() {
        return true
    } 
    logs.Error("Configuration -> Check Record -> "+db+" record "+table+"->"+uuid_value+"->"+param_value+" doesn't exists -> err: "+err.Error())
    return false
}

func TableCreate(conn string, tablename string, create string)(ok bool){
    logs.Info("Configuration -> Creating table "+tablename+" in "+conn)
    loadDataSQL := map[string]map[string]string{}
    loadDataSQL[conn] = map[string]string{}
    loadDataSQL[conn]["path"] = ""
    loadDataSQL, err := utils.GetConf(loadDataSQL)
    if err != nil {
        logs.Error("Configuration -> Can't get "+conn+" path from main.conf -> "+err.Error())
        return false
    }
    dbpath := loadDataSQL[conn]["path"]
    db, err := sql.Open("sqlite3",dbpath)
    if err != nil {
        logs.Error("Configuration -> "+dbpath+" Open Failed -> err: "+err.Error())
        return false
    }
    _, err = db.Exec(create)
    if err != nil {
        logs.Error("Configuration -> Creating table " +tablename + " failed -> err: "+err.Error())
        return false
    }
    return true 
}


func TableCreatePluginPlugins()(ok bool){
    conn := "pluginConn"
    tablename := "plugins"
    create := "CREATE TABLE plugins ("+
                "plugin_id integer PRIMARY KEY AUTOINCREMENT,"+
                "plugin_uniqueid text NOT NULL,"+
                "plugin_param text NOT NULL,"+
                "plugin_value text NOT NULL)"
    return TableCreate(conn, tablename, create)
}




func DbCreate(db string)(err error) {
    logs.Warn ("Configuration -> Creating DB file -> "+db)
    _, err = os.OpenFile(db, os.O_CREATE, 0644)
    if err != nil {
        logs.Error("Configuration -> Creating DB File "+ db +" err: "+err.Error())
        return err
    }
    return nil
}
