package utils

import (
    //"owlhnode/models"
    "encoding/json"
    "strconv"
    //"github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
    "io/ioutil"
    //"io"
    "os"
    "time"
    "os/exec"
    "fmt"
)


func GetConf(loadData map[string]map[string]string)(loadDataReturn map[string]map[string]string) { //leer json del fichero para obtener el path del bpf
    confFilePath := "/etc/owlh/conf/main.conf"
    jsonPathBpf, err := ioutil.ReadFile(confFilePath)
    if err != nil {
        logs.Error("utils/GetConf -> can't open Conf file -> " + confFilePath)
        return nil
    }

    var anode map[string]map[string]string
    json.Unmarshal(jsonPathBpf, &anode)

    logs.Debug(anode)

    logs.Error("|................|")
    for k,y := range loadData { 
        for y,_ := range y {
            if v, ok := anode[k][y]; ok {
                logs.Debug(k+"-"+y+"-"+v)
                logs.Info(anode[k][y])
                logs.Notice(loadData[k][y])
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

func BackupFile(path string, file string) (err error) { 
    t := time.Now()

    newFile := file+"-"+strconv.FormatInt(t.Unix(), 10)
    logs.Info ("NODE:UTILS.GO // NEW FILE NAME -->" + newFile)

    srcFolder := path+file
    destFolder := path+newFile
    logs.Info ("NODE:UTILS.GO // OLD FILE NAME -->" + srcFolder)
    logs.Info ("NODE:UTILS.GO // DST FILE NAME -->" + destFolder)
    cpCmd := exec.Command("cp", srcFolder, destFolder)
    err = cpCmd.Run()
    if err != nil{
        logs.Info ("Erro exec cmd command")
    }

    return err
}