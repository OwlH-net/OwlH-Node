package ndb

import (
    "github.com/astaxie/beego/logs"
    "database/sql"
	"os"
	"owlhnode/utils"
	_ "github.com/mattn/go-sqlite3"
)

var (
    Pdb *sql.DB
)

func PConn() {
    var err error
	loadDataSQL := map[string]map[string]string{}
	loadDataSQL["pluginConn"] = map[string]string{}
	loadDataSQL["pluginConn"]["path"] = ""
	loadDataSQL["pluginConn"]["cmd"] = "" 
	loadDataSQL, err = utils.GetConf(loadDataSQL)    
    path := loadDataSQL["pluginConn"]["path"]
	cmd := loadDataSQL["pluginConn"]["cmd"]
	if err != nil {
		logs.Error("PConn Error getting data from main.conf")
	}
	_, err = os.Stat(path) 
	if err != nil {
		panic("Fail opening plugins.db from path: "+path+"  --  "+err.Error())
	}	
	Pdb, err = sql.Open(cmd,path)
    if err != nil {
        logs.Error("Pdb/plugins -- plugins.db Open Failed: "+err.Error())
    }else {
		logs.Info("Pdb/plugins -- plugins.db -> sql.Open, plugins.db Ready") 
	}
}