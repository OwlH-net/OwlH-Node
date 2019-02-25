package file

import (
    "github.com/astaxie/beego/logs"
    // "os"
    // "os/exec"
    // "strings"
    // "regexp"
    "owlhnode/utils"
    "io/ioutil"
)

//read file and send back to webpage
func SendFile(file string)(data map[string]string, err error){
	var voidArray map[string]string

	sendBackArray := make(map[string]string)

	logs.Info("Intro")

    //create map and obtain file
    loadData := map[string]map[string]string{}
	loadData["files"] = map[string]string{}
	loadData["files"][file] = ""
	loadData = utils.GetConf(loadData)
	


    
    //save url from file selected and open file
    fileConfPath := loadData["files"][file]
    logs.Warn(fileConfPath)
	URLFile, err := ioutil.ReadFile(fileConfPath) // just pass the file name
    if err != nil {
		logs.Info("eRROR")
        return voidArray,err
    }
	
	sendBackArray["fileContent"] = string(URLFile)
	logs.Info(file+" // "+string(URLFile))
	sendBackArray["fileName"] = file

	logs.Warn("	tras lectura ")
	logs.Info(sendBackArray["fileContent"])
	logs.Info(sendBackArray["fileName"])

	logs.Info("outro")


	return sendBackArray, err
}

//read changed file, make a backup and save into file
func SaveFile(file map[string]string)(err error){
    path := "/etc/owlh/conf/"
    fileName := "main.conf"
    fileRetrieved := file["file"]

    //make file backup before overwrite
    err = utils.BackupFile(path, fileName)
    if err != nil {
        return err
    }

    //make byte array for save the file modified
    bytearray := []byte(fileRetrieved)
    err = utils.WriteNewDataOnFile(path+fileName, bytearray)
    if err != nil {
        return err
    }
    return err
}