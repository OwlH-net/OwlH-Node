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


func GetConf()(path string, file string) { //leer json del fichero para obtener el path del bpf
    
    //crear json con bpf: path y bpf: file
    //crear directorios y fichero
    confFilePath := "/etc/owlh/conf/main.conf"
    jsonPathBpf, err := ioutil.ReadFile(confFilePath)
    if err != nil {
        logs.Error ("utils/GetConf -> can't open Conf file -> " + confFilePath)
        return "", ""
    }
    //confFile:=string(jsonPathBpf)

    var anode map[string]string
    json.Unmarshal(jsonPathBpf, &anode)

    logs.Error ("utils.GetConf  || path --> "+anode["path"]+" file -->"+anode["file"])


    return anode["path"], anode["file"]



/*

    confFile, err := os.Open(confFilePath)
    if err != nil {
        logs.Error ("utils/GetConf -> can't open Conf file -> " + confFilePath)
        return ""
    }
    defer confFile.Close()

    var anode map[string]string
    json.Unmarshal(confFile, &anode)

    return anode["bpf"]
*/

    /*
    byteValue, _ := ioutil.ReadAll(confFile)

    var config map[string]string
    json.Unmarshal([]byte(byteValue), &config)

    if value, exists := config[param]; exists {
        return value
    } else {
        return "ERROR"
    }
    */
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

    logs.Info ("NODE:UTILS.GO //  PATH  -->" + path)
    logs.Info ("NODE:UTILS.GO //  FILE  -->" + file)

    newFile := file+"-"+strconv.FormatInt(t.Unix(), 10)
    logs.Info ("NODE:UTILS.GO // NEW FILE NAME -->" + newFile)

    srcFolder := path+file
    destFolder := path+newFile
    cpCmd := exec.Command("cp", srcFolder, destFolder)
    err = cpCmd.Run()
    if err != nil{
        logs.Info ("Erro exec cmd command")
    }

    return err

/*
    //reader
    fileReader, err := os.Open(path+file)
    if (err != nil) { 
        logs.Info ("Error io.Reader for make a Backup file")
        return err
    }
    //writer


    old := data, err := ioutil.ReadFile(path+file)
    new := path+newFile

    err = io.Copy(old, new)
    return err
    
    sourceFileStat, err := os.Stat(path+file)
    if err != nil {
            return err
    }
    in, err := os.Open(path+file)
    if err != nil {
        return err
    }
    
    defer file.Close()
    t := time.Now()
    rename := file+"-"+strconv.FormatInt(t.Unix(), 10)

    dst := os.Rename(path+file, path+rename)

    out, err := os.Create(dst)
    if err != nil {
        return err
    }

    defer out.Close()

    nBytes, err := io.Copy(destination, source)
    return nBytes, err
    */
}