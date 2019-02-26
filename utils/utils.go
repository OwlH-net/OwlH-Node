package utils

import (
    //"owlhnode/models"
    "encoding/json"
    "strconv"
    //"github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
    "io/ioutil"
    // "io"
    // "strings"
    "os"
    "time"
    "os/exec"
    "fmt"
)

//Read map data
//leer json del fichero para obtener el path del bpf
func GetConf(loadData map[string]map[string]string)(loadDataReturn map[string]map[string]string) { 
    confFilePath := "/etc/owlh/conf/main.conf"
    jsonPathBpf, err := ioutil.ReadFile(confFilePath)
    if err != nil {
        logs.Error("utils/GetConf -> can't open Conf file -> " + confFilePath)
        return nil
    }

    var anode map[string]map[string]string
    json.Unmarshal(jsonPathBpf, &anode)

    logs.Error("|................|")
    for k,y := range loadData { 
        for y,_ := range y {
            if v, ok := anode[k][y]; ok {
                loadData[k][y] = v
            }else{
                loadData[k][y] = "None"
            }
        }
    }
    
    return loadData
}

func UpdateBPFFile(path string, file string, bpf string) (err error) {
    //delete file content
    err = os.Truncate(path+file, 0)
	if err != nil {
		logs.Info(err)
    }
    //write new bpf content
    newBPF, err := os.OpenFile(path+file, os.O_RDWR|os.O_APPEND|os.O_CREATE, 0660)
    if err != nil {
		logs.Info(err)
        os.Exit(-1)
        return err
	}
	defer newBPF.Close()
	fmt.Fprintf(newBPF, "%s\n", bpf)

    return nil
}

func BackupFullPath(path string) (err error) { 
    t := time.Now()

    destFolder := path+"-"+strconv.FormatInt(t.Unix(), 10)
    cpCmd := exec.Command("cp", path, destFolder)

    err = cpCmd.Run()
    if err != nil{
        logs.Info ("Erro exec cmd command")
        return err
    }

    return nil
}
    

func BackupFile(path string, fileName string) (err error) { 
    t := time.Now()
    newFile := fileName+"-"+strconv.FormatInt(t.Unix(), 10)
    srcFolder := path+fileName
    destFolder := path+newFile
    cpCmd := exec.Command("cp", srcFolder, destFolder)
    err = cpCmd.Run()
    if err != nil{
        logs.Info ("Erro exec cmd command")
        return err
    }

    return nil
}

func WriteNewDataOnFile(path string, data []byte)(err error){
    
    logs.Info("WriteNewDataOnFile  path->__   "+path)
    logs.Info("WriteNewDataOnFile  data->__   "+string(data))

    err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
        logs.Info("Error WriteNewData")
		return err
	}

    return nil
}

//leer ficheros
func GetConfFiles()(loadDataReturn map[string]string, err error) { 
    confFilePath := "/etc/owlh/conf/main.conf"
    JSONconf, err := ioutil.ReadFile(confFilePath)
    if err != nil {
        logs.Error("utils/GetConf -> can't open Conf file -> " + confFilePath)
        return nil, err
    }

    var anode map[string]map[string]string
    json.Unmarshal(JSONconf, &anode)
    
    logs.Info("MOSTRAR ficheros MAIN.CONF")
    logs.Info(anode["files"])

    return anode["files"], nil
}