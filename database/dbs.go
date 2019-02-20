package ndb

import (
    "github.com/astaxie/beego/logs"
    "database/sql"
    "strconv"
    "owlhnode/utils"
//    "fmt"
//   "time"
    _ "github.com/mattn/go-sqlite3"
)

var (
    Db *sql.DB
)

type Master struct {
    Id       string
    Name     string
    Ip       string
    Port     int
}

func init() {
    logs.Info ("DB Init ()")
}

func Conn() {
    logs.Info("DB -> sql.Open, let's try to be Ready")
    var err error
    logs.Info("Set Suricata BPF -- Making Map")

    //Retrieve path and command for open sql.
	loadDataSQL := map[string]map[string]string{}
	loadDataSQL["dbsConn"] = map[string]string{}
	loadDataSQL["dbsConn"]["path"] = ""
	loadDataSQL["dbsConn"]["cmd"] = "" 
    loadDataSQL = utils.GetConf(loadDataSQL)    
    path := loadDataSQL["dbsConn"]["path"]
    cmd := loadDataSQL["dbsConn"]["cmd"]

    //Db, err = sql.Open("sqlite3", "/etc/owlh/databases/node.db")
    Db, err = sql.Open(cmd, path)
    if err != nil {
        panic("DB Open Failed ")
    }
    logs.Info("DB -> sql.Open, DB Ready") 
}


func Get_master() (err error) {
    logs.Info("DB -> Get Master")
    var (
        id string
        name string
        ip string
        port int
    )

    if Db != nil {
        sql := "SELECT * FROM master WHERE master_id=1;"
        row := Db.QueryRow(sql) //$
        logs.Info ("DB -> Row %s", row) //$
        err = row.Scan(&id, &name, &ip, &port)

        if err != nil {
            logs.Warn("DB -> Can't read database")
        }
        logs.Info("DB -> Master : id - %s, name - %s, ip - %s, port - %d", id,name,ip,port)
    } else {
        panic("DB -> there is no DB")
    }
    return nil
}



func Insert_master(mname string, mip string, mport int) (err error) {
    logs.Info("DB -> Insert Master")
    logs.Info("DB -> name - %s, ip - %s, port - %d", mname, mip, mport)

    if Db != nil {
        err = Db.Ping()
        if err != nil {
            logs.Info ("DB -> Can't access DB")
            return err
        }
        stmt, err := Db.Prepare("INSERT INTO master(master_name, master_ip, master_port) values(?,?,?)")
        if err != nil {
            panic("DB Prepare failed")
        }
        logs.Info("DB -> db.Prepare, err = %s", err)

        res, err := stmt.Exec(mname, mip, strconv.Itoa(mport))
        logs.Info("DB -> info", res)
        if err != nil {
            panic("DB bad Query ")
            return err
        }
    } else {
        panic("DB -> there is no DB")
    }
    return nil
}

func Close() {
    Db.Close()
}