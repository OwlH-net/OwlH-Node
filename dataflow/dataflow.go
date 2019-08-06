package dataflow

import (
    "github.com/astaxie/beego/logs"
    "owlhnode/database"
)

func ChangeDataflowValues(anode map[string]string)(err error) {
	err = ndb.ChangeDataflowValues(anode["FlowUUID"], anode["param"], anode["value"])
	if (err != nil){ logs.Error("ChangeDataflowValues UPDATE error for update dataflow values: "+err.Error()); return err}
    return nil
}

func LoadDataflowValues()(data map[string]map[string]string, err error) {
	data,err = ndb.LoadDataflowValues()
	if (err != nil){ logs.Error("ChangeDataflowValues UPDATE error for update dataflow values: "+err.Error()); return nil,err}
    return data,nil
}

func SaveSocketToNetwork(anode map[string]string) (err error) {
    err = ndb.ChangeDataflowValues("sockettonetwork", "interface", anode["interface"])
	if (err != nil){ logs.Error("ChangeDataflowValues UPDATE error for update dataflow values: "+err.Error()); return err}
    err = ndb.ChangeDataflowValues("sockettonetwork", "cert", anode["cert"])
	if (err != nil){ logs.Error("ChangeDataflowValues UPDATE error for update dataflow values: "+err.Error()); return err}
    err = ndb.ChangeDataflowValues("sockettonetwork", "port", anode["port"])
	if (err != nil){ logs.Error("ChangeDataflowValues UPDATE error for update dataflow values: "+err.Error()); return err}
    return nil
}

// func LoadSocketToNetwork() (data map[string]map[string]string, err error) {
//     err = ndb.LoadSocketToNetwork()
// 	if (err != nil){ logs.Error("ChangeDataflowValues UPDATE error for update dataflow values: "+err.Error()); return nil,err}
//     return nil
// }