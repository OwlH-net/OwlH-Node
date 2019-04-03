package file

import (
    "github.com/astaxie/beego/logs"
    // "os"
    // "os/exec"
    // "strings"
    // "regexp"
    "owlhnode/utils"
    "io/ioutil"
    // "encoding/json"
)

//read file and send back to webpage
func SendFile(file string)(data map[string]string, err error){
	//var voidArray map[string]string
	sendBackArray := make(map[string]string)
    
    //create map and obtain file
    loadData := map[string]map[string]string{}
	loadData["files"] = map[string]string{}
	loadData["files"][file] = ""
	loadData,err = utils.GetConf(loadData)
	if err != nil {
        logs.Error("Error getting path and BPF from main.conf")
    }
	    
    //save url from file selected and open file
    fileConfPath := loadData["files"][file]
	fileReaded, err := ioutil.ReadFile(fileConfPath) // just pass the file name
    if err != nil {
		logs.Info("Error reading file for path: "+fileConfPath)
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
        logs.Error("Error getting path and BPF from main.conf")
    }

    //make file backup before overwrite
    err = utils.BackupFullPath(loadData["files"][file["file"]])
    logs.Info("after bck SaveFile->__   ")
    if err != nil {
        logs.Info("Backup error->__   ")
        return err
    }

    //make byte array for save the file modified
    bytearray := []byte(file["content"])
    err = utils.WriteNewDataOnFile(loadData["files"][file["file"]], bytearray)
    if err != nil {
        return err
    }
    return err
}

func GetAllFiles()(data map[string]string, err error){

    // var files []string     
    var returnedData map[string]string
	returnedData,err = utils.GetConfFiles()
	if err != nil {
        logs.Error("Error getting path and BPF from main.conf")
    }
    //var files []string
    // for k,_ := range returnedData { 
    //     files = append(files, k)
    //     //files = append(files,k)
    // }
    // result := make(map[string][]string)
    // result["fileNames"] = files
    logs.Info("GetAllFiles -> returing file names")
    return returnedData, err

}