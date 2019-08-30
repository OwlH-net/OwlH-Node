package knownports

import (
    "github.com/astaxie/beego/logs"
	"errors"
	"owlhnode/database"
)

func ShowPorts() (data map[string]map[string]string, err error) {

	var uniqueid string
	var param string
	var value string
	var allKnownPorts = map[string]map[string]string{}

	//database connection
	if ndb.Pdb == nil {
        logs.Error("ShowPorts knownports -- Can't access to database")
        return nil,errors.New("ShowPorts knownports -- Can't access to database")
	} 
	//query and make map[]map[]
	sql := "select kp_uniqueid, kp_param, kp_value from knownports;"   
	rows, err := ndb.Pdb.Query(sql)
	defer rows.Close()
    if err != nil {
        logs.Error("ShowPorts knownports Error executing query: %s", err.Error())
        return nil, err
    }
	for rows.Next() {
        if err = rows.Scan(&uniqueid, &param, &value); err != nil {
            logs.Error("ShowPorts knownports -- Can't read query result: %s", err.Error())
			return nil, err
        }
        if allKnownPorts[uniqueid] == nil { allKnownPorts[uniqueid] = map[string]string{}}
        allKnownPorts[uniqueid][param]=value
	} 
	return allKnownPorts, nil
}

func PingPluginsNode() (data map[string]map[string]string ,err error) {
	var uniqueid string
	var param string
	var value string
	var allKnownPorts = map[string]map[string]string{}

	//database connection
	if ndb.Pdb == nil {
        logs.Error("PingPluginsNode knownports -- Can't access to database")
        return nil,errors.New("PingPluginsNode knownports -- Can't access to database")
	} 

	//query and make map[]map[]
	sql := "select plugin_uniqueid, plugin_param, plugin_value from plugins;"   
	rows, err := ndb.Pdb.Query(sql)
	defer rows.Close()
    if err != nil {
        logs.Error("PingPluginsNode knownports Error executing query: %s", err.Error())
        return nil, err
    }
	for rows.Next() {
        if err = rows.Scan(&uniqueid, &param, &value); err != nil {
            logs.Error("PingPluginsNode knownports -- Can't read query result: %s", err.Error())
            return nil, err
        }
        if allKnownPorts[uniqueid] == nil { allKnownPorts[uniqueid] = map[string]string{}}
        allKnownPorts[uniqueid][param]=value
	} 
	return allKnownPorts, nil
}

func ChangeStatus(anode map[string]string) (err error) {
	value := anode["status"]
	plugin:= anode["plugin"]
	protoportUpdate, err := ndb.Pdb.Prepare("update plugins set plugin_value = ? where plugin_param = ? and plugin_uniqueid = ?")
	defer protoportUpdate.Close()
	_, err = protoportUpdate.Exec(&value, "status", &plugin)
	if err != nil {
		logs.Error("ChangeMode --> update error-> %s", err.Error())
		return err
	}
	return nil
}

func ChangeMode(anode map[string]string) (err error) {
	value := anode["mode"]
	protoportUpdate, err := ndb.Pdb.Prepare("update plugins set plugin_value = ? where plugin_param = ? and plugin_uniqueid = ?")
	defer protoportUpdate.Close()
	_, err = protoportUpdate.Exec(&value, "mode", "knownports")
	if err != nil {
		logs.Error("ChangeMode --> update error-> %s", err.Error())
		return err
	}
	return nil
}

func DeletePorts(ports map[string]string) (err error) {
	for id := range ports {
		protoportUpdate, err := ndb.Pdb.Prepare("delete from knownports where kp_uniqueid = ?")
		defer protoportUpdate.Close()
		_, err = protoportUpdate.Exec(&id)
		if err != nil {
			logs.Error("DeletePorts --> update error-> %s", err.Error())
			return err
		}
	}
	anode := make(map[string]string)
	anode["plugin"]="knownports"
	anode["status"]="Reload"
	_ = ChangeStatus(anode)

	return nil
}

func DeleteAllPorts() (err error) {
	protoportUpdate, err := ndb.Pdb.Prepare("delete from knownports")
	defer protoportUpdate.Close()
	_, err = protoportUpdate.Exec()
	if err != nil {
		logs.Error("DeleteAllPorts --> update error-> %s", err.Error())
		return err
	}

	anode := make(map[string]string)
	anode["plugin"]="knownports"
	anode["status"]="Reload"
	_ = ChangeStatus(anode)

	return nil
}