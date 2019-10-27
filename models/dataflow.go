package models 

import (
    "owlhnode/dataflow"
//    "owlhnode/changeControl"
    "github.com/astaxie/beego/logs")

func ChangeDataflowValues(anode map[string]string) (err error) {
    cc := anode
    logs.Info("============")
    logs.Info("DATAFLOW - ChangeDataflowValues")
    for key :=range anode {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")
    
    err = dataflow.ChangeDataflowValues(anode)
    return err
}

func SaveSocketToNetwork(anode map[string]string) (err error) {
    cc := anode
    logs.Info("============")
    logs.Info("DATAFLOW - ChangeDataflowValues")
    for key :=range anode {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")
    
    err = dataflow.ChangeDataflowValues(anode)
    return err

	err = dataflow.SaveSocketToNetwork(anode)
	return err
}

func LoadDataflowValues()(data map[string]map[string]string, err error) {
    data, err = dataflow.LoadDataflowValues()
    return data, err
}

func SaveNewLocal(anode map[string]string)(err error) {
    cc := anode
    logs.Info("============")
    logs.Info("DATAFLOW - SaveNewLocal")
    for key :=range anode {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")
    
    err = dataflow.ChangeDataflowValues(anode)
    return err

    err = dataflow.SaveNewLocal(anode)
    return err
}

func SaveVxLAN(anode map[string]string)(err error) {
    cc := anode
    logs.Info("============")
    logs.Info("DATAFLOW - SaveVxLAN")
    for key :=range anode {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")
    
    err = dataflow.SaveVxLAN(anode)
    return err
}

func SaveSocketToNetworkSelected(anode map[string]string)(err error) {
    cc := anode
    logs.Info("============")
    logs.Info("DATAFLOW - SaveSocketToNetworkSelected")
    for key :=range anode {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")
    
    err = dataflow.SaveSocketToNetworkSelected(anode)
    return err
}

func DeleteDataFlowValueSelected(anode map[string]string)(err error) {
    cc := anode
    logs.Info("============")
    logs.Info("DATAFLOW - DeleteDataFlowValueSelected")
    for key :=range anode {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")
    
    err = dataflow.DeleteDataFlowValueSelected(anode)
    return err
}