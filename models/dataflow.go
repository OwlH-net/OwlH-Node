package models 

import (
    "owlhnode/dataflow"
    "owlhnode/changeControl"
    "github.com/astaxie/beego/logs")

func ChangeDataflowValues(anode map[string]string, username string) (err error) {
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

    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }

    cc["actionDescription"] = "Change Traffic transport values"


    changecontrol.InsertChangeControl(cc)
    //changecontrol.ChangeControlInsertData(err, "ChangeDataflowValues")    
    return err
}

func SaveSocketToNetwork(anode map[string]string, username string) (err error) {
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

    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }

    cc["actionDescription"] = "Set Socket to Network details"

    changecontrol.InsertChangeControl(cc)

    //changecontrol.ChangeControlInsertData(err, "SaveSocketToNetwork")    
    return err
}

func LoadDataflowValues(username string)(data map[string]map[string]string, err error) {
    data, err = dataflow.LoadDataflowValues()
    //changecontrol.ChangeControlInsertData(err, "LoadDataflowValues")    
    return data, err
}

func SaveNewLocal(anode map[string]string, username string)(err error) {
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

    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }

    cc["actionDescription"] = "Set traffic transport values"

    changecontrol.InsertChangeControl(cc)
    //changecontrol.ChangeControlInsertData(err, "SaveNewLocal")    
    return err
}

func SaveVxLAN(anode map[string]string, username string)(err error) {
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
    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }

    cc["actionDescription"] = "Set VxLAN interface config"

    changecontrol.InsertChangeControl(cc)
    //changecontrol.ChangeControlInsertData(err, "SaveVxLAN")    
    return err
}

func SaveSocketToNetworkSelected(anode map[string]string, username string)(err error) {
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
    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }

    cc["actionDescription"] = "Set Socket to Network"

    changecontrol.InsertChangeControl(cc)


    //changecontrol.ChangeControlInsertData(err, "SaveSocketToNetworkSelected")    
    return err
}

func DeleteDataFlowValueSelected(anode map[string]string, username string)(err error) {
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
    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }
    
    cc["username"] = username
    cc["actionDescription"] = "Delete Traffic Transport configuration"

    changecontrol.InsertChangeControl(cc)

    //changecontrol.ChangeControlInsertData(err, "DeleteDataFlowValueSelected")    
    return err
}