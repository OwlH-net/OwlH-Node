package models

import (
    "owlhnode/file"
    "owlhnode/changeControl"
    "github.com/astaxie/beego/logs")

func SendFile(filename string) (data map[string]string, err error) {
    logs.Info("SendFile into Node file")
    data,err = file.SendFile(filename)
    //changecontrol.ChangeControlInsertData(err, "SendFile")    
    return data,err
}

func SaveFile(data map[string]string) (err error) {
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
    //changecontrol.ChangeControlInsertData(err, "SaveFile")    
    return err
}

func GetAllFiles() (data map[string]string, err error) {
    logs.Info("GetAllFiles into Node file")
    data,err = file.GetAllFiles()
    //changecontrol.ChangeControlInsertData(err, "GetAllFiles")    
    return data,err
}

func ReloadFilesData() (data map[string]map[string]string, err error) {
    data,err = file.ReloadFilesData()
    //changecontrol.ChangeControlInsertData(err, "ReloadFilesData")    
    return data,err
}