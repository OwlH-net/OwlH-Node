package models

import (
    "owlhnode/file"
    "owlhnode/changeControl"
    "github.com/astaxie/beego/logs"
)

func SendFile(filename string, username string) (data map[string]string, err error) {
    data,err = file.SendFile(filename)
    //changecontrol.ChangeControlInsertData(err, "SendFile", username)    
    return data,err
}

func SaveFile(data map[string]string, username string) (err error) {
    cc := data
    logs.Info("============")
    logs.Info("FILE - SaveFile")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(data,"action")
    delete(data,"controller")
    delete(data,"router")

    logs.Info("SaveFile into Node file")
    err = file.SaveFile(data)

    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }
    cc["username"] = username
    cc["actionDescription"] = "Save file"

    changecontrol.InsertChangeControl(cc)

    //changecontrol.ChangeControlInsertData(err, "SaveFile")    
    return err
}

func GetAllFiles(username string) (data map[string]string, err error) {
    logs.Info("GetAllFiles into Node file")
    data,err = file.GetAllFiles()
    //changecontrol.ChangeControlInsertData(err, "GetAllFiles")    
    return data,err
}

func ReloadFilesData(username string) (data map[string]map[string]string, err error) {
    data,err = file.ReloadFilesData()
    //changecontrol.ChangeControlInsertData(err, "ReloadFilesData")    
    return data,err
}