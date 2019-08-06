package file

import (
    "github.com/astaxie/beego/logs"
    "owlhnode/utils"
    "io/ioutil"
)

//read file and send back to webpage
func SendFile(file string)(data map[string]string, err error){
	sendBackArray := make(map[string]string)
    
    //create map and obtain file
    loadData := map[string]map[string]string{}
	loadData["files"] = map[string]string{}
	loadData["files"][file] = ""
	loadData,err = utils.GetConf(loadData)
	if err != nil { logs.Error("SendFile Error getting data from main.conf"); return nil,err}
	    
    //save url from file selected and open file
    fileConfPath := loadData["files"][file]
	fileReaded, err := ioutil.ReadFile(fileConfPath) // just pass the file name
    if err != nil {
		logs.Error("Error reading file for path: "+fileConfPath)
        return nil,err
    }
	
	sendBackArray["fileContent"] = string(fileReaded)
    sendBackArray["fileName"] = file

	return sendBackArray, nil
}

//read changed file, make a backup and save into file
func SaveFile(file map[string]string)(err error){
    //Get full path
    loadData := map[string]map[string]string{}
    loadData["files"] = map[string]string{}
    loadData["files"][file["file"]] = ""
	loadData,err = utils.GetConf(loadData)
	if err != nil {
        logs.Error("SaveFile Error getting data from main.conf")
    }

    //make file backup before overwrite
    err = utils.BackupFullPath(loadData["files"][file["file"]])
    if err != nil {
        logs.Info("SaveFile. Error doing backup with function BackupFullPath: "+err.Error())
        return err
    }

    //make byte array for save the file modified
    bytearray := []byte(file["content"])
    err = utils.WriteNewDataOnFile(loadData["files"][file["file"]], bytearray)
    if err != nil {
		logs.Info("SaveFile. Error doing backup with function WriteNewDataOnFile: "+err.Error())
        return err
    }
    return nil
}

func GetAllFiles()(data map[string]string, err error){

    var returnedData map[string]string
	returnedData,err = utils.GetConfFiles()
	if err != nil {
        logs.Error("Error getting data from main.conf")
    }
    logs.Info("GetAllFiles -> returing file names")
    return returnedData, err

}