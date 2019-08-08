package dataflow

import (
    "github.com/astaxie/beego/logs"
    "owlhnode/database"
    "owlhnode/utils"
)

func ChangeDataflowValues(anode map[string]string)(err error) {
	err = ndb.ChangeDataflowValues(anode["FlowUUID"], anode["param"], anode["value"])
	if (err != nil){ logs.Error("ChangeDataflowValues UPDATE error for update dataflow values: "+err.Error()); return err}
    return nil
}

func LoadDataflowValues()(data map[string]map[string]string, err error) {
	data,err = ndb.LoadDataflowValues()
	if (err != nil){ logs.Error("LoadDataflowValues error loading dataflow values: "+err.Error()); return nil,err}
    return data,nil
}

func SaveSocketToNetwork(anode map[string]string) (err error) {
	uuid := utils.Generate()
	err = ndb.InsertDataflowValues(uuid, "name", anode["name"])
	if (err != nil){ logs.Error("SaveSocketToNetwork error inserting name dataflow values: "+err.Error()); return err}
	err = ndb.InsertDataflowValues(uuid, "cert", anode["cert"])
	if (err != nil){ logs.Error("SaveSocketToNetwork error inserting cert dataflow values: "+err.Error()); return err}
	err = ndb.InsertDataflowValues(uuid, "interface", anode["interface"])
	if (err != nil){ logs.Error("SaveSocketToNetwork error inserting interface dataflow values: "+err.Error()); return err}
	err = ndb.InsertDataflowValues(uuid, "port", anode["port"])
	if (err != nil){ logs.Error("SaveSocketToNetwork error inserting port dataflow values: "+err.Error()); return err}
	err = ndb.InsertDataflowValues(uuid, "type", "sockettonetwork")
	if (err != nil){ logs.Error("SaveSocketToNetwork error inserting sockettonetwork dataflow values: "+err.Error()); return err}
    
    return nil
}

func SaveNewLocal(anode map[string]string)(err error) {
	uuid := utils.Generate()
	err = ndb.InsertDataflowValues(uuid, "name", anode["name"])
	if (err != nil){ logs.Error("SaveSocketToNetwork error inserting name dataflow values: "+err.Error()); return err}
	err = ndb.InsertDataflowValues(uuid, "mtu", anode["mtu"])
	if (err != nil){ logs.Error("SaveSocketToNetwork error inserting cert dataflow values: "+err.Error()); return err}
	err = ndb.InsertDataflowValues(uuid, "type", "networknewlocal")
	if (err != nil){ logs.Error("SaveSocketToNetwork error inserting sockettonetwork dataflow values: "+err.Error()); return err}
    
    return nil
}

func SaveVxLAN(anode map[string]string)(err error) {
	logs.Info("Saving new local...")
    logs.Notice(anode)
    return nil
}

func SaveSocketToNetworkSelected(anode map[string]string)(err error) {
    logs.Notice(anode["uuidNode"])
    return nil
}

func DeleteSocketToNetworkSelected(anode map[string]string)(err error) {
	err = ndb.DeleteSocketToNetworkSelected(anode["uuidNode"])
	if (err != nil){ logs.Error("DeleteSocketToNetworkSelected error deleting a socket->network: "+err.Error()); return err}
    return nil
}