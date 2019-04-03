package ndb

import (
    "github.com/astaxie/beego/logs"
    "database/sql"
//    "fmt"
//   "time"
    _ "github.com/mattn/go-sqlite3"
    // "errors"
    "owlhnode/utils"
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
    loadDataSQL = utils.GetConf(loadDataSQL)    
    path := loadDataSQL["stapConn"]["path"]
    cmd := loadDataSQL["stapConn"]["cmd"]
   
    //Sdb, err = sql.Open("sqlite3", "database/servers.db")
    Sdb, err = sql.Open(cmd,path)
    if err != nil {
        panic("sdb/servers -- DB Open Failed")
    }
    logs.Info("sdb/servers -- DB -> sql.Open, DB Ready") 
}

func GetStapServerInformation(uuid string)(serverData map[string]string, err error){
	var param string
	var value string
	stapServer := make(map[string]string)
	logs.Info("Creating data map for uuid: "+uuid)
	// ip, err := ndb.Sdb.Query("select server_param,server_value from servers where server_param = \"ip\" and server_uniqueid = \""+uuid+"\";")
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