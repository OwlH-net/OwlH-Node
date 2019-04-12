package ndb

import (
    "github.com/astaxie/beego/logs"
    "database/sql"
	"os"
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

// func Close() {
// 	Sdb.Close()
// }