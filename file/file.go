package file

import (
    "github.com/astaxie/beego/logs"
    "owlhnode/utils"
    "owlhnode/wazuh"
    "io/ioutil"
    "os"
    "strconv"
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
        logs.Error("SaveFile. Error doing backup with function BackupFullPath: "+err.Error())
        return err
    }

    //make byte array for save the file modified
    bytearray := []byte(file["content"])
    err = utils.WriteNewDataOnFile(loadData["files"][file["file"]], bytearray)
    if err != nil {
		logs.Error("SaveFile. Error doing backup with function WriteNewDataOnFile: "+err.Error())
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

func ReloadFilesData() (data map[string]map[string]string, err error) {
    sendBackArray := map[string]map[string]string{}
    sendBackArray["wazuh"] = map[string]string{}
    sendBackArray["analyzer"] = map[string]string{}
    
    //ANALYZER
    loadData := map[string]map[string]string{}
	loadData["node"] = map[string]string{}
	loadData["node"]["alertLog"] = ""
    loadData,err = utils.GetConf(loadData)
    alertLog := loadData["node"]["alertLog"]
	if err != nil { logs.Error("ReloadFilesData Error getting data from main.conf"); return nil,err}
    
    // size = fi.Size()
    fi, err := os.Stat(alertLog);
    if err != nil { logs.Error("ReloadFilesData.Error reading analyzer file: "+err.Error())}
    sendBackArray["analyzer"]["size"] = strconv.FormatInt(fi.Size(), 10)
    sendBackArray["analyzer"]["path"] = alertLog
    
    //WAZUH
    wazuhFiles,err := wazuh.PingWazuhFiles()
    for x,_ := range wazuhFiles{
        sendBackArray["wazuh"][wazuhFiles[x]["path"]] = wazuhFiles[x]["size"]
    }

    logs.Notice(sendBackArray)
    return sendBackArray,nil
}