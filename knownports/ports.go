package knownports

import (
    "github.com/astaxie/beego/logs"
	"errors"
	"owlhnode/database"
)

func ShowPorts() (data map[string]string, err error) {

	var value string
	var allKnownPorts = map[string]string{}

	//database connection
	if ndb.Pdb == nil {
        logs.Error("ShowPorts knownports -- Can't access to database")
        return nil,errors.New("ShowPorts knownports -- Can't access to database")
	} 
	//query and make map[]map[]
	sql := "select kp_value from knownports where kp_param='portprot';"   
	rows, err := ndb.Pdb.Query(sql)
	defer rows.Close()
    if err != nil {
        logs.Error("ShowPorts knownports Error executing query: %s", err.Error())
        return nil, err
    }
	for rows.Next() {
        if err = rows.Scan(&value); err != nil {
            logs.Error("ShowPorts knownports -- Can't read query result: %s", err.Error())
            return nil, err
        }
        allKnownPorts[value]=value
	} 
	return allKnownPorts, nil
}

func PingPorts() (data map[string]map[string]string ,err error) {
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
	sql := "select plugin_uniqueid, plugin_param, plugin_value from plugins;"   
	rows, err := ndb.Pdb.Query(sql)
	defer rows.Close()
    if err != nil {
        logs.Error("ShowPorts knownports Error executing query: %s", err.Error())
        return nil, err
    }
	for rows.Next() {
        if err = rows.Scan(&uniqueid, &param, &value); err != nil {
            logs.Error("LoadPorts knownports -- Can't read query result: %s", err.Error())
            return nil, err
        }
        if allKnownPorts[uniqueid] == nil { allKnownPorts[uniqueid] = map[string]string{}}
        allKnownPorts[uniqueid][param]=value
	} 
	return allKnownPorts, nil
}

func ChangeStatus(anode map[string]string) (err error) {
	value := anode["status"]
	protoportUpdate, err := ndb.Pdb.Prepare("update plugins set plugin_value = ? where plugin_param = ? and plugin_uniqueid = ?")
	defer protoportUpdate.Close()
	_, err = protoportUpdate.Exec(&value, "status", "knownports")
	if err != nil {
		logs.Error("ChangeMode --> update error-> %s", err.Error())
		return err
	}
	if value == "Enabled" {
		Init()
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
