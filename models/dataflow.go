package models 

import (
    "owlhnode/dataflow"
    "owlhnode/changeControl"
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
    //changecontrol.ChangeControlInsertData(err, "ChangeDataflowValues")    
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

    err = dataflow.SaveSocketToNetwork(anode)
    //changecontrol.ChangeControlInsertData(err, "SaveSocketToNetwork")    
    return err
}

func LoadDataflowValues()(data map[string]map[string]string, err error) {
    data, err = dataflow.LoadDataflowValues()
    //changecontrol.ChangeControlInsertData(err, "LoadDataflowValues")    
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
    //changecontrol.ChangeControlInsertData(err, "SaveNewLocal")    
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
    //changecontrol.ChangeControlInsertData(err, "SaveVxLAN")    
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
    //changecontrol.ChangeControlInsertData(err, "SaveSocketToNetworkSelected")    
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
    //changecontrol.ChangeControlInsertData(err, "DeleteDataFlowValueSelected")    
    return err
}