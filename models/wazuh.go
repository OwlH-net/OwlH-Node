package models

import (
    "owlhnode/wazuh"
   "owlhnode/changeControl"
    "github.com/astaxie/beego/logs"
)


func GetWazuh() (status map[string]bool, err error) {
    //changecontrol.ChangeControlInsertData(err, "GetWazuh")    
    return wazuh.Installed()
}

func RunWazuh() (data string, err error) {
    logs.Info("Run RunWazuh system into node server")
    data,err = wazuh.RunWazuh()
    var cc = make(map[string]string)

    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }

    cc["action"] = "PUT"
    cc["controller"] = "WAZUH"
    cc["router"] = "@router /RunWazuh [put]"

    cc["actionDescription"] = "Start Wazuh"
    cc["actionReult"] = data

    changecontrol.InsertChangeControl(cc)
    // changecontrol.ChangeControlInsertData(err, "RunWazuh")    
    return data,err
}

func StopWazuh() (data string, err error) {
    logs.Info("Stops StopWazuh system into node server")
    data,err = wazuh.StopWazuh()
    var cc = make(map[string]string)
    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }

    cc["action"] = "PUT"
    cc["controller"] = "WAZUH"
    cc["router"] = "@router /StopWazuh [put]"

    cc["actionDescription"] = "Stop Wazuh"
    cc["actionReult"] = data

    changecontrol.InsertChangeControl(cc)
    // changecontrol.ChangeControlInsertData(err, "StopWazuh")    
    return data,err
}

func PingWazuhFiles() (files map[int]map[string]string, err error) {
    files, err = wazuh.PingWazuhFiles()
    // changecontrol.ChangeControlInsertData(err, "PingWazuhFiles")    
    return files ,err
}

func DeleteWazuhFile(file map[string]interface{})(err error) {
    cc := file
    logs.Info("============")
    logs.Info("WAZUH - DeleteWazuhFile")
    for key :=range cc {
        logs.Info(key +" -> ")
    }
    delete(file,"action")
    delete(file,"controller")
    delete(file,"router")

    err = wazuh.ModifyWazuhFile(file)

    // if err!=nil { 
    //     cc["actionStatus"] = "error"
    //     cc["errorDescription"] = err.Error()
    // }else{
    //     cc["actionStatus"] = "success"
    // }

    // cc["actionDescription"] = "Delete Wazuh file to monitor"

    // changecontrol.InsertChangeControl(cc)
    // changecontrol.ChangeControlInsertData(err, "DeleteWazuhFile")    
    return err
}

func AddWazuhFile(file map[string]interface{})(err error) {
    cc := file
    logs.Info("============")
    logs.Info("WAZUH - AddWazuhFile")
    for key :=range cc {
        logs.Info(key +" -> ")
    }
    delete(file,"action")
    delete(file,"controller")
    delete(file,"router")
    
    err = wazuh.ModifyWazuhFile(file)
    // changecontrol.ChangeControlInsertData(err, "AddWazuhFile")    
    return err
}

func LoadFileLastLines(file map[string]string)(data map[string]string, err error) {
    data, err = wazuh.LoadFileLastLines(file)
    // changecontrol.ChangeControlInsertData(err, "LoadFileLastLines")    
    return data, err
}

func SaveFileContentWazuh(file map[string]string)(err error) {
    cc := file
    logs.Info("============")
    logs.Info("WAZUH - SaveFileContentWazuh")
    for key :=range cc {
        logs.Info(key +" -> ")
    }
    delete(file,"action")
    delete(file,"controller")
    delete(file,"router")
    
     err = wazuh.SaveFileContentWazuh(file)

         if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }

    cc["actionDescription"] = "Save Wazuh config"

    changecontrol.InsertChangeControl(cc)
    // changecontrol.ChangeControlInsertData(err, "SaveFileContentWazuh")    
    return  err
}