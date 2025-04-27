package ndb

import (
	"database/sql"
	"errors"
	"os"

	"github.com/OwlH-net/OwlH-Node/utils"
	"github.com/astaxie/beego/logs"
	_ "github.com/mattn/go-sqlite3"
)

var (
	Pdb *sql.DB
)

func PConn() {
	var err error
	path, err := utils.GetKeyValueString("pluginConn", "path")
	if err != nil {
		logs.Error("PConn Error getting data from main.conf")
	}
	cmd, err := utils.GetKeyValueString("pluginConn", "cmd")
	if err != nil {
		logs.Error("PConn Error getting data from main.conf")
	}
	_, err = os.Stat(path)
	if err != nil {
		panic("Fail opening plugins.db from path: " + path + "  --  " + err.Error())
	}
	Pdb, err = sql.Open(cmd, path)
	if err != nil {
		logs.Error("Pdb/plugins -- plugins.db Open Failed: " + err.Error())
	} else {
		logs.Info("Pdb/plugins -- plugins.db -> sql.Open, plugins.db Ready")
	}
}

func GetStatusAnalyzer() (data string, err error) {
	var value string

	sql := "select analyzer_value from analyzer where analyzer_param='status'"
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
	return value, nil
}

func UpdateAnalyzer(uuid string, param string, value string) (err error) {
	updateAnalyzerNode, err := Pdb.Prepare("update analyzer set analyzer_value = ? where analyzer_uniqueid = ? and analyzer_param = ?;")
	if err != nil {
		logs.Error("updateAnalyzerNode UPDATE prepare error: " + err.Error())
		return err
	}

	_, err = updateAnalyzerNode.Exec(&value, &uuid, &param)
	if err != nil {
		logs.Error("updateAnalyzerNode UPDATE exec error: " + err.Error())
		return err
	}

	defer updateAnalyzerNode.Close()

	return nil
}

func InsertPluginService(uuid string, param string, value string) (err error) {
	insertPlugin, err := Pdb.Prepare("insert into plugins(plugin_uniqueid, plugin_param, plugin_value) values (?,?,?);")
	if err != nil {
		logs.Error("InsertPluginService INSERT prepare error: " + err.Error())
		return err
	}

	_, err = insertPlugin.Exec(&uuid, &param, &value)
	if err != nil {
		logs.Error("InsertPluginService INSERT exec error: " + err.Error())
		return err
	}

	defer insertPlugin.Close()

	return nil
}

func GetServices(service string) (path map[string]map[string]string, err error) {
	var serviceValues = map[string]map[string]string{}
	var uniqid string
	var param string
	var value string
	var uuidSource string
	rowsQuery, err := Pdb.Query("select plugin_uniqueid from plugins where plugin_value='" + service + "';")
	if err != nil {
		logs.Error("GetServices Pdb.Query Error : %s", err.Error())
		return nil, err
	}
	defer rowsQuery.Close()
	for rowsQuery.Next() {
		if err = rowsQuery.Scan(&uuidSource); err != nil {
			logs.Error("GetServices -- UUID Query return error: %s", err.Error())
			return nil, err
		}

		sql := "select plugin_uniqueid,plugin_param,plugin_value from plugins where plugin_uniqueid='" + uuidSource + "';"

		rows, err := Pdb.Query(sql)
		if err != nil {
			logs.Error("GetServices - get uuid items error - Pdb.Query Error : %s", err.Error())
			return nil, err
		}
		defer rows.Close()
		for rows.Next() {
			if err = rows.Scan(&uniqid, &param, &value); err != nil {
				logs.Error("GetServices -- get uuid values scan  error: %s", err.Error())
				return nil, err
			}
			if serviceValues[uniqid] == nil {
				serviceValues[uniqid] = map[string]string{}
			}
			serviceValues[uniqid][param] = value
		}
	}
	return serviceValues, nil
}

func GetPlugins() (path map[string]map[string]string, err error) {
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
		if err = rowsQuery.Scan(&uniqid, &param, &value); err != nil {
			logs.Error("GetPlugins -- Query return error: %s", err.Error())
			return nil, err
		}
		if serviceValues[uniqid] == nil {
			serviceValues[uniqid] = map[string]string{}
		}
		serviceValues[uniqid][param] = value
	}

	return serviceValues, nil
}

func GetPluginsByParam(uniqueid string, param string) (value string, err error) {
	rowsQuery, err := Pdb.Query("select plugin_value from plugins where plugin_uniqueid = '" + uniqueid + "' and plugin_param = '" + param + "';")
	if err != nil {
		logs.Error("GetPlugins Pdb.Query Error : %s", err.Error())
		return "", err
	}

	defer rowsQuery.Close()
	for rowsQuery.Next() {
		if err = rowsQuery.Scan(&value); err != nil {
			logs.Error("GetPluginsByParam -- Query return error: %s", err.Error())
			return "error", err
		}
	}
	return value, nil
}

func GetMainconfParam(uniqid string, param string) (value string, err error) {
	sql := "select main_value from mainconf where main_uniqueid='" + uniqid + "' and main_param='" + param + "'"
	rows, err := Pdb.Query(sql)
	if err != nil {
		logs.Error("GetMainconfData Pdb.Query Error : %s", err.Error())
		return "", err
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&value); err != nil {
			logs.Error("GetMainconfData -- Query return error: %s", err.Error())
			return "", err
		}
	}
	return value, nil

}

func GetMainconfData() (values map[string]map[string]string, err error) {
	var mainconfValues = map[string]map[string]string{}
	var uniqid string
	var param string
	var value string

	sql := "select main_uniqueid, main_param, main_value from mainconf;"

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
		if mainconfValues[uniqid] == nil {
			mainconfValues[uniqid] = map[string]string{}
		}
		mainconfValues[uniqid][param] = value
	}
	return mainconfValues, nil
}

func InsertGetMainconfData(uuid string, param string, value string) (err error) {
	insertPlugin, err := Pdb.Prepare("insert into mainconf(main_uniqueid, main_param, main_value) values (?,?,?);")
	if err != nil {
		logs.Error("InsertGetMainconfData INSERT prepare error: " + err.Error())
		return err
	}

	_, err = insertPlugin.Exec(&uuid, &param, &value)
	if err != nil {
		logs.Error("InsertGetMainconfData INSERT exec error: " + err.Error())
		return err
	}

	defer insertPlugin.Close()

	return nil
}

func UpdatePluginValue(uuid string, param string, value string) (err error) {
	UpdatePluginValueNode, err := Pdb.Prepare("update plugins set plugin_value = ? where plugin_uniqueid = ? and plugin_param = ?;")
	if err != nil {
		logs.Error("UpdatePluginValue UPDATE prepare error: " + err.Error())
		return err
	}

	_, err = UpdatePluginValueNode.Exec(&value, &uuid, &param)
	if err != nil {
		logs.Error("UpdatePluginValue UPDATE exec error: " + err.Error())
		return err
	}

	defer UpdatePluginValueNode.Close()

	return nil
}

func UpdateMainconfValue(uuid string, param string, value string) (err error) {
	currentvalue, err := GetMainconfParam(uuid, param)
	if err != nil || currentvalue == "" {
		InsertGetMainconfData(uuid, param, value)
	} else {
		UpdateMainconfValueNode, err := Pdb.Prepare("update mainconf set main_value = ? where main_uniqueid = ? and main_param = ?;")
		if err != nil {
			logs.Error("UpdateMainconfValue UPDATE prepare error: " + err.Error())
			return err
		}

		_, err = UpdateMainconfValueNode.Exec(&value, &uuid, &param)
		if err != nil {
			logs.Error("UpdateMainconfValue exec error: " + err.Error())
			return err
		}

		defer UpdateMainconfValueNode.Close()
	}

	return nil
}

func DeleteService(uuid string) (err error) {
	DeleteServiceNode, err := Pdb.Prepare("delete from plugins where plugin_uniqueid = ?;")
	if err != nil {
		logs.Error("DeleteService UPDATE prepare error: " + err.Error())
		return err
	}

	_, err = DeleteServiceNode.Exec(&uuid)
	if err != nil {
		logs.Error("DeleteService exec error: " + err.Error())
		return err
	}

	defer DeleteServiceNode.Close()

	return nil
}

func InsertClusterData(uuid string, param string, value string) (err error) {
	insertPlugin, err := Pdb.Prepare("insert into cluster(cluster_uniqueid, cluster_param, cluster_value) values (?,?,?);")
	if err != nil {
		logs.Error("InsertClusterData INSERT prepare error: " + err.Error())
		return err
	}

	_, err = insertPlugin.Exec(&uuid, &param, &value)
	if err != nil {
		logs.Error("InsertClusterData INSERT exec error: " + err.Error())
		return err
	}

	defer insertPlugin.Close()

	return nil
}

func CountDBEntries(clusterType string) (count int, err error) {
	rows, err := Pdb.Query("select count(*) from cluster where cluster_value like '%" + clusterType + "%';")
	if err != nil {
		logs.Error("CountDBEntries Pdb.Query Error : %s", err.Error())
		return -1, err
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&count); err != nil {
			logs.Error("CountDBEntries -- Query return error: %s", err.Error())
			return -1, err
		}
	}
	return count, nil
}

func GetClusterData() (path map[string]map[string]string, err error) {
	var clusterValues = map[string]map[string]string{}
	var uniqid string
	var param string
	var value string

	sql := "select cluster_uniqueid, cluster_param, cluster_value from cluster;"

	rows, err := Pdb.Query(sql)
	if err != nil {
		logs.Error("GetClusterData Pdb.Query Error : %s", err.Error())
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		if err = rows.Scan(&uniqid, &param, &value); err != nil {
			logs.Error("GetClusterData -- Query return error: %s", err.Error())
			return nil, err
		}
		if clusterValues[uniqid] == nil {
			clusterValues[uniqid] = map[string]string{}
		}
		clusterValues[uniqid][param] = value
	}
	return clusterValues, nil
}

func UpdateClusterValue(uuid string, param string, value string) (err error) {
	UpdateClusterValueNode, err := Pdb.Prepare("update cluster set cluster_value = ? where cluster_uniqueid = ? and cluster_param = ?;")
	if err != nil {
		logs.Error("UpdateClusterValue UPDATE prepare error: " + err.Error())
		return err
	}

	_, err = UpdateClusterValueNode.Exec(&value, &uuid, &param)
	if err != nil {
		logs.Error("UpdateClusterValue exec error: " + err.Error())
		return err
	}

	defer UpdateClusterValueNode.Close()

	return nil
}

func DeleteClusterValue(uuid string) (err error) {
	DeleteServiceNode, err := Pdb.Prepare("delete from cluster where cluster_uniqueid = ?;")
	if err != nil {
		logs.Error("DeleteClusterValue DELETE prepare error: " + err.Error())
		return err
	}

	_, err = DeleteServiceNode.Exec(&uuid)
	if err != nil {
		logs.Error("DeleteClusterValue exec error: " + err.Error())
		return err
	}

	defer DeleteServiceNode.Close()

	return nil
}

func DeleteAllClusters() (err error) {
	DeleteServiceNode, err := Pdb.Prepare("delete from cluster;")
	if err != nil {
		logs.Error("DeleteAllClusters DELETE prepare error: " + err.Error())
		return err
	}

	_, err = DeleteServiceNode.Exec()
	if err != nil {
		logs.Error("DeleteAllClusters exec error: " + err.Error())
		return err
	}

	defer DeleteServiceNode.Close()

	return nil
}

func LoadPortsData() (data map[string]map[string]string, err error) {
	var uniqueid string
	var param string
	var value string
	var allKnownPorts = map[string]map[string]string{}

	//database connection
	if Pdb == nil {
		logs.Error("LoadPorts knownports -- Can't access to database")
		return nil, errors.New("LoadPorts knownports -- Can't access to database")
	}
	//query and make map[]map[]
	sql := "select kp_uniqueid, kp_param, kp_value from knownports;"
	rows, err := Pdb.Query(sql)
	defer rows.Close()
	if err != nil {
		logs.Error("LoadPorts knownports Error executing query: %s", err.Error())
		return nil, err
	}
	for rows.Next() {
		if err = rows.Scan(&uniqueid, &param, &value); err != nil {
			logs.Error("LoadPorts knownports -- Can't read query result: %s", err.Error())
			return nil, err
		}
		if allKnownPorts[uniqueid] == nil {
			allKnownPorts[uniqueid] = map[string]string{}
		}
		allKnownPorts[uniqueid][param] = value
	}
	return allKnownPorts, nil
}

func CheckParamKnownports(param string) (data string, err error) {
	var res string
	res = "Disabled"
	sql := "select plugin_value from plugins where plugin_uniqueid = 'knownports' and plugin_param='" + param + "'"
	rows, err := Pdb.Query(sql)
	defer rows.Close()
	if err != nil {
		logs.Error("knownports CheckParamKnownports Error executing query: %s", err.Error())
		return "", err
	}
	if rows.Next() {
		if err = rows.Scan(&res); err != nil {
			logs.Error("knownports CheckParamKnownports -- Can't read query result: %s", err.Error())
			return "", err
		}
	}
	return res, nil
}

func InsertKnownports(uuid string, param string, value string) (err error) {
	insertKP, err := Pdb.Prepare("insert into knownports (kp_uniqueid, kp_param, kp_value) values (?,?,?);")
	_, err = insertKP.Exec(&uuid, &param, &value)
	defer insertKP.Close()
	if err != nil {
		logs.Error("Error InsertKnownports: " + err.Error())
		return err
	}
	return nil
}

func UpdateKnownports(uuid string, param string, value string) (err error) {
	UpdateKnownportsNode, err := Pdb.Prepare("update knownports set kp_value = ? where kp_param = ? and kp_uniqueid = ?;")
	if err != nil {
		logs.Error("UpdateKnownports UPDATE prepare error: " + err.Error())
		return err
	}

	_, err = UpdateKnownportsNode.Exec(&value, &uuid, &param)
	if err != nil {
		logs.Error("UpdateKnownports UPDATE exec error: " + err.Error())
		return err
	}

	defer UpdateKnownportsNode.Close()

	return nil
}

func DeletePort(ports map[string]string) (err error) {
	for id := range ports {
		protoportDelete, err := Pdb.Prepare("delete from knownports where kp_uniqueid = ?")
		defer protoportDelete.Close()
		_, err = protoportDelete.Exec(&id)
		if err != nil {
			logs.Error("DeletePort --> delete error-> %s", err.Error())
			return err
		}
	}
	return nil
}

func DeleteAllPorts() (err error) {
	protoportDelete, err := Pdb.Prepare("delete from knownports;")
	defer protoportDelete.Close()
	_, err = protoportDelete.Exec()
	if err != nil {
		logs.Error("DeleteAllPorts --> delete error-> %s", err.Error())
		return err
	}

	return nil
}

func InsertPluginCommand(uuid string, param string, value string) (err error) {
	insertPlugin, err := Pdb.Prepare("insert into commands(cmd_uniqueid, cmd_param, cmd_value) values (?,?,?);")
	if err != nil {
		logs.Error("InsertPluginCommand INSERT prepare error: " + err.Error())
		return err
	}

	_, err = insertPlugin.Exec(&uuid, &param, &value)
	if err != nil {
		logs.Error("InsertPluginCommand INSERT exec error: " + err.Error())
		return err
	}

	defer insertPlugin.Close()

	return nil
}

func GetPluginCommands() (data map[string]map[string]string, err error) {
	var uniqueid string
	var param string
	var value string
	var allCommands = map[string]map[string]string{}

	//database connection
	if Pdb == nil {
		logs.Error("GetPluginCommands -- Can't access to database")
		return nil, errors.New("GetPluginCommands -- Can't access to database")
	}
	//query and make map[]map[]
	sql := "select cmd_uniqueid, cmd_param, cmd_value from commands;"
	rows, err := Pdb.Query(sql)
	defer rows.Close()
	if err != nil {
		logs.Error("GetPluginCommands Error executing query: %s", err.Error())
		return nil, err
	}
	for rows.Next() {
		if err = rows.Scan(&uniqueid, &param, &value); err != nil {
			logs.Error("GetPluginCommands -- Can't read query result: %s", err.Error())
			return nil, err
		}
		if allCommands[uniqueid] == nil {
			allCommands[uniqueid] = map[string]string{}
		}
		allCommands[uniqueid][param] = value
	}
	return allCommands, nil
}
