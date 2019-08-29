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

func GetStatusAnalyzer()(data string, err error){
	var value string

	sql := "select analyzer_value from analyzer where analyzer_param='status'";
	rows, err := Pdb.Query(sql)
	if err != nil {
		logs.Error("GetStatusAnalyzer Pdb.Query Error : %s", err.Error())
		return "", err
	}
	for rows.Next() {
		if err = rows.Scan(&value); err != nil {
            logs.Error("GetStatusAnalyzer -- Query return error: %s", err.Error())
            return "", err
		}
	} 
	return value,nil
}

func UpdateAnalyzer(uuid string, param string, value string)(err error){
	updateAnalyzerNode, err := Pdb.Prepare("update analyzer set analyzer_value = ? where analyzer_uniqueid = ? and analyzer_param = ?;")
	if (err != nil){ logs.Error("updateAnalyzerNode UPDATE prepare error: "+err.Error()); return err}

	_, err = updateAnalyzerNode.Exec(&value, &uuid, &param)
	if (err != nil){ logs.Error("updateAnalyzerNode UPDATE exec error: "+err.Error()); return err}

	defer updateAnalyzerNode.Close()
	
	return nil
}

func InsertPluginService(uuid string, param string, value string)(err error){
	updateAnalyzerNode, err := Pdb.Prepare("insert into plugins(plugin_uniqueid, plugin_param, plugin_value) values (?,?,?);")
	if (err != nil){ logs.Error("InsertPluginService INSERT prepare error: "+err.Error()); return err}

	_, err = updateAnalyzerNode.Exec(&uuid, &param, &value)
	if (err != nil){ logs.Error("InsertPluginService INSERT exec error: "+err.Error()); return err}

	defer updateAnalyzerNode.Close()
	
	return nil
}

func GetServices(service string)(path map[string]map[string]string, err error){
	var serviceValues = map[string]map[string]string{}
    var uniqid string
    var param string
	var value string
	var uuidSource string
	rowsQuery, err := Pdb.Query("select plugin_uniqueid from plugins where plugin_value='"+service+"';")
	if err != nil {
		logs.Error("GetServices Pdb.Query Error : %s", err.Error())
		return nil, err
	}
	defer rowsQuery.Close()
	for rowsQuery.Next() {
		if err = rowsQuery.Scan(&uuidSource); err != nil { logs.Error("GetServices -- Query return error: %s", err.Error()); return nil, err}
		logs.Warn(uuidSource)
		logs.Warn(uuidSource)
		logs.Warn(uuidSource)
		logs.Warn(uuidSource)
		logs.Warn(uuidSource)
		logs.Warn(uuidSource)
		
		sql := "select plugin_uniqueid,plugin_param,plugin_value from plugins where plugin_uniqueid='"+uuidSource+"';";
		
		rows, err := Pdb.Query(sql)
		if err != nil {
			logs.Error("GetServices Pdb.Query Error : %s", err.Error())
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			if err = rows.Scan(&service); err != nil {
				logs.Error("GetServices -- Query return error: %s", err.Error())
				return nil, err
			}
			if serviceValues[uniqid] == nil { serviceValues[uniqid] = map[string]string{}}
			serviceValues[uniqid][param]=value
		}
	} 
	return serviceValues,nil
}