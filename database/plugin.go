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
	insertPlugin, err := Pdb.Prepare("insert into plugins(plugin_uniqueid, plugin_param, plugin_value) values (?,?,?);")
	if (err != nil){ logs.Error("InsertPluginService INSERT prepare error: "+err.Error()); return err}

	_, err = insertPlugin.Exec(&uuid, &param, &value)
	if (err != nil){ logs.Error("InsertPluginService INSERT exec error: "+err.Error()); return err}

	defer insertPlugin.Close()
	
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

func GetPlugins()(path map[string]map[string]string, err error){
	var serviceValues = map[string]map[string]string{}
    var uniqid string
    var param string
	var value string
	rowsQuery, err := Pdb.Query("select plugin_uniqueid, plugin_param, plugin_value from plugins;")
	if err != nil {
		logs.Error("GetPlugins Pdb.Query Error : %s", err.Error())
		return nil, err
	}
	defer rowsQuery.Close()
	for rowsQuery.Next() {
		if err = rowsQuery.Scan(&uniqid, &param, &value); err != nil { logs.Error("GetPlugins -- Query return error: %s", err.Error()); return nil, err}

		if serviceValues[uniqid] == nil { serviceValues[uniqid] = map[string]string{}}
		serviceValues[uniqid][param]=value
	} 
	return serviceValues,nil
}

func GetPluginsByParam(uniqueid string, param string)(value string, err error){
	rowsQuery, err := Pdb.Query("select plugin_value from plugins where plugin_uniqueid = '"+uniqueid+"' and plugin_param = '"+param+"';")
	if err != nil {logs.Error("GetPlugins Pdb.Query Error : %s", err.Error()); return "", err}

	defer rowsQuery.Close()
	for rowsQuery.Next() {
		if err = rowsQuery.Scan(&value); err != nil { logs.Error("GetPluginsByParam -- Query return error: %s", err.Error()); return "error", err} 
	} 
	return value,nil
}

func GetMainconfData()(path map[string]map[string]string, err error){
	var mainconfValues = map[string]map[string]string{}
    var uniqid string
    var param string
	var value string

	sql := "select main_uniqueid, main_param, main_value from mainconf;";
	
	rows, err := Pdb.Query(sql)
	if err != nil {
		logs.Error("GetMainconfData Pdb.Query Error : %s", err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&uniqid, &param, &value); err != nil {
            logs.Error("GetMainconfData -- Query return error: %s", err.Error())
            return nil, err
		}
        if mainconfValues[uniqid] == nil { mainconfValues[uniqid] = map[string]string{}}
        mainconfValues[uniqid][param]=value
	} 
	return mainconfValues,nil
}

func UpdatePluginValue(uuid string, param string, value string)(err error){
	UpdatePluginValueNode, err := Pdb.Prepare("update plugins set plugin_value = ? where plugin_uniqueid = ? and plugin_param = ?;")
	if (err != nil){ logs.Error("UpdatePluginValue UPDATE prepare error: "+err.Error()); return err}

	_, err = UpdatePluginValueNode.Exec(&value, &uuid, &param)
	if (err != nil){ logs.Error("UpdatePluginValue UPDATE exec error: "+err.Error()); return err}

	defer UpdatePluginValueNode.Close()
	
	return nil
}

func UpdateMainconfValue(uuid string, param string, value string)(err error){
	UpdateMainconfValueNode, err := Pdb.Prepare("update mainconf set main_value = ? where main_uniqueid = ? and main_param = ?;")
	if (err != nil){ logs.Error("UpdateMainconfValue UPDATE prepare error: "+err.Error()); return err}

	_, err = UpdateMainconfValueNode.Exec(&value, &uuid, &param)
	if (err != nil){ logs.Error("UpdateMainconfValue exec error: "+err.Error()); return err}

	defer UpdateMainconfValueNode.Close()
	
	return nil
}

func DeleteService(uuid string)(err error){
	DeleteServiceNode, err := Pdb.Prepare("delete from plugins where plugin_uniqueid = ?;")
	if (err != nil){ logs.Error("DeleteService UPDATE prepare error: "+err.Error()); return err}

	_, err = DeleteServiceNode.Exec(&uuid)
	if (err != nil){ logs.Error("DeleteService exec error: "+err.Error()); return err}

	defer DeleteServiceNode.Close()
	
	return nil
}