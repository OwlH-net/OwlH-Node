package file

import (
	"io/ioutil"
	"os"
	"strconv"

	"github.com/OwlH-net/OwlH-Node/utils"
	"github.com/OwlH-net/OwlH-Node/wazuh"
	"github.com/astaxie/beego/logs"
)

// read file and send back to webpage
func SendFile(file string) (data map[string]string, err error) {
	sendBackArray := make(map[string]string)
	var key string
	var fileName string
	if file == "node.cfg" || file == "networks.cfg" || file == "zeekctl.cfg" {
		if file == "node.cfg" {
			key = "nodeconfig"
		} else if file == "networks.cfg" {
			key = "networkconfig"
		} else {
			key = "zeekctlconfig"
		}

		fileName, err = utils.GetKeyValueString("zeek", key)
		if err != nil {
			logs.Error("SendFile Error getting data from main.conf")
			return nil, err
		}
	} else {
		//create map and obtain file
		fileName, err = utils.GetKeyValueString("files", file)
		if err != nil {
			logs.Error("SendFile Error getting data from main.conf")
			return nil, err
		}
	}

	//save url from file selected and open file
	fileReaded, err := ioutil.ReadFile(fileName) // just pass the file name
	if err != nil {
		logs.Error("Error reading file for path: " + fileName)
		return nil, err
	}

	sendBackArray["fileContent"] = string(fileReaded)
	sendBackArray["fileName"] = file

	return sendBackArray, nil
}

// read changed file, make a backup and save into file
func SaveFile(file map[string]string) (err error) {
	//Get full path
	fileName, err := utils.GetKeyValueString("files", file["file"])
	if err != nil {
		logs.Error("SaveFile Error getting data from main.conf")
	}

	//make file backup before overwrite
	err = utils.BackupFullPath(fileName)
	if err != nil {
		logs.Error("SaveFile. Error doing backup with function BackupFullPath: " + err.Error())
		return err
	}

	//make byte array for save the file modified
	bytearray := []byte(file["content"])
	err = utils.WriteNewDataOnFile(fileName, bytearray)
	if err != nil {
		logs.Error("SaveFile. Error doing backup with function WriteNewDataOnFile: " + err.Error())
		return err
	}
	return nil
}

func GetAllFiles() (data map[string]string, err error) {

	var returnedData map[string]string
	returnedData, err = utils.GetConfFiles()
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
	alertLog, err := utils.GetKeyValueString("node", "alertLog")
	if err != nil {
		logs.Error("ReloadFilesData Error getting data from main.conf")
		return nil, err
	}

	// size = fi.Size()
	fi, err := os.Stat(alertLog)
	if err != nil {
		logs.Error("ReloadFilesData.Error reading analyzer file: " + err.Error())
	}
	sendBackArray["analyzer"]["size"] = strconv.FormatInt(fi.Size(), 10)
	sendBackArray["analyzer"]["path"] = alertLog

	//WAZUH
	wazuhFiles, err := wazuh.PingWazuhFiles()
	for x, _ := range wazuhFiles {
		sendBackArray["wazuh"][wazuhFiles[x]["path"]] = wazuhFiles[x]["size"]
	}

	return sendBackArray, nil
}

func RotateFile(path string) {
	// numberOfFiles := 5
	// compress := true

	// oldpath := path+".4.tar.gz"
	// if existsFile(oldpath) {
	//     renameFile(oldpath, path+".5.tar.gz")
	// }
	// oldpath = path+".3.tar.gz"
	// if existsFile(oldpath) {
	//     renameFile(oldpath, path+".4.tar.gz")
	// }
	// oldpath = path+".2.tar.gz"
	// if existsFile(oldpath) {
	//     renameFile(oldpath, path+".3.tar.gz")
	// }
	// oldpath = path+".1"
	// if existsFile(oldpath) {
	//     renameFile(oldpath, path+".5.tar.gz")
	// }
	// os.Rename(src, dst)
}
