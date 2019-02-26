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
	var voidArray map[string]string
	sendBackArray := make(map[string]string)
    
    //create map and obtain file
    loadData := map[string]map[string]string{}
	loadData["files"] = map[string]string{}
	loadData["files"][file] = ""
	loadData = utils.GetConf(loadData)
	    
    //save url from file selected and open file
    fileConfPath := loadData["files"][file]
	fileReaded, err := ioutil.ReadFile(fileConfPath) // just pass the file name
    if err != nil {
		logs.Info("eRROR")
        return voidArray,err
    }
	
	sendBackArray["fileContent"] = string(fileReaded)
    sendBackArray["fileName"] = file

	return sendBackArray, err
}

//read changed file, make a backup and save into file
func SaveFile(file map[string]string)(err error){

    logs.Info("entro SaveFile->__   ")
    //Get full path
    loadData := map[string]map[string]string{}
    loadData["files"] = map[string]string{}
    loadData["files"][file["file"]] = ""
    loadData = utils.GetConf(loadData)

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