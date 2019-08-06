package dataflow

import (
    "github.com/astaxie/beego/logs"
    "owlhnode/database"
)

func ChangeDataflowValues(anode map[string]string)(err error) {
	err = ndb.ChangeDataflowValues(anode)
	if (err != nil){ logs.Error("ChangeDataflowValues UPDATE error for update dataflow values: "+err.Error()); return err}
    return nil
}

func LoadDataflowValues()(data map[string]map[string]string, err error) {
	data,err = ndb.LoadDataflowValues()
	if (err != nil){ logs.Error("ChangeDataflowValues UPDATE error for update dataflow values: "+err.Error()); return nil,err}
    return data,nil
}
