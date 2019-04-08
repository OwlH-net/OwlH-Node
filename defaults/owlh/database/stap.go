package ndb

import (
    "github.com/astaxie/beego/logs"
    "database/sql"
//    "fmt"
//   "time"
    _ "github.com/mattn/go-sqlite3"
    //"errors"
)

var (
    Sdb *sql.DB
)

func SConn() {
    var err error

    //Retrieve path and command for open sql.
	loadDataSQL := map[string]map[string]string{}
	loadDataSQL["stapConn"] = map[string]string{}
	loadDataSQL["stapConn"]["path"] = ""
	loadDataSQL["stapConn"]["cmd"] = "" 
	loadDataSQL,err = utils.GetConf(loadDataSQL)    
    path := loadDataSQL["stapConn"]["path"]
    cmd := loadDataSQL["stapConn"]["cmd"]
	if err != nil {
		logs.Error("Sconn Error getting data from main.conf")
	}
	
    //Sdb, err = sql.Open("sqlite3", "database/servers.db")
    Sdb, err = sql.Open(cmd,path)
    if err != nil {
        panic("sdb/servers -- DB Open Failed")
    }
    logs.Info("sdb/servers -- DB -> sql.Open, DB Ready") 
}