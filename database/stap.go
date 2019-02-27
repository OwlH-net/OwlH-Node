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
	Sdb, err = sql.Open("sqlite3", "database/servers.db")
    if err != nil {
        panic("sdb/servers -- DB Open Failed")
    }
    logs.Info("sdb/servers -- DB -> sql.Open, DB Ready") 
}