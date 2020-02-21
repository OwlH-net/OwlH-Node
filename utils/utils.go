package utils

import (
    "encoding/json"
    "strconv"
    "github.com/astaxie/beego/logs"
    "io/ioutil"
    "io"
    "errors"
    "os"
    "time"
    "os/exec"
    "fmt"
    "path/filepath"
    "strings"
    "crypto/rand"
)


//read data from main.conf
func GetConf(loadData map[string]map[string]string)(loadDataReturn map[string]map[string]string, err error) { 
    confFilePath := "conf/main.conf"
    jsonPath, err := ioutil.ReadFile(confFilePath)
    if err != nil {
        logs.Error("utils/GetConf -> can't open Conf file: " + confFilePath)
        return nil, err
    }

    var anode map[string]map[string]string
    json.Unmarshal(jsonPath, &anode)

    for k,y := range loadData { 
        for y,_ := range y {
            if v, ok := anode[k][y]; ok {
                loadData[k][y] = v
            }else{
                loadData[k][y] = "None"
            }
        }
    }
    return loadData, nil
}

func UpdateBPFFile(path string, file string, bpf string) (err error) {
    //delete file content
    err = os.Truncate(path+file, 0)
    if err != nil {
        logs.Error("Error truncate BPF file: "+err.Error())
        return err
    }

    //write new bpf content
    bpfByteArray := []byte(bpf)
    err = WriteNewDataOnFile(path+file, bpfByteArray)
    if err != nil {logs.Error("Error writing new BPF data into file: "+err.Error()); return err}
    return nil
}

//create a BPF backup
func BackupFullPath(path string) (err error) { 
    t := time.Now()
    destFolder := path+"-"+strconv.FormatInt(t.Unix(), 10)
    cpCmd := exec.Command("cp", path, destFolder)
    err = cpCmd.Run()
    if err != nil{
        logs.Error("utils.BackupFullPath Error exec cmd command: "+err.Error())
        return err
    }

    return nil
}

func BackupFile(path string, fileName string) (err error) {  
    loadData := map[string]map[string]string{}
    loadData["node"] = map[string]string{}
    loadData["node"]["backupFolder"] = "" 
    loadData,err = GetConf(loadData)  
    backupFolder := loadData["node"]["backupFolder"]    
    if err != nil {logs.Error("utils.BackupFile Error getting backup path: "+err.Error()); return err}

    // check if folder exists
    if _, err := os.Stat(backupFolder); os.IsNotExist(err) {
        err = os.MkdirAll(backupFolder, 0755)
        if err != nil{logs.Error("utils.BackupFile Error creating main backup folder: "+err.Error()); return err}
    }

    //get older backup file
    listOfFiles,err := FilelistPathByFile(backupFolder, fileName)
    if err != nil{logs.Error("utils.BackupFile Error walking through backup folder: "+err.Error()); return err}
    count := 0
    previousBck := ""
    for x := range listOfFiles{
        count++
        if previousBck == "" {
            previousBck = x
            continue
        }else if previousBck > x{
            previousBck = x
        }
    }

    //delete older bck file if there are 5 bck files
    if count == 5 {
        err = os.Remove(backupFolder+previousBck)
        if err != nil{logs.Error("utils.BackupFile Error deleting older backup file: "+err.Error())}
    }

    //create backup
    t := time.Now()
    newFile := fileName+"-"+strconv.FormatInt(t.Unix(), 10)
    srcFolder := path+fileName
    destFolder := backupFolder+newFile

    //check if file exist
    if _, err := os.Stat(srcFolder); os.IsNotExist(err) {
        return errors.New("utils.BackupFile error: Source file doesn't exists")
    }else{
        cpCmd := exec.Command("cp", srcFolder, destFolder)
        err = cpCmd.Run()
        if err != nil{logs.Error("utils.BackupFile Error exec cmd command: "+err.Error()); return err}
    }
    return nil
}

//write data on a file
func WriteNewDataOnFile(path string, data []byte)(err error){
    err = ioutil.WriteFile(path, data, 0644)
    if err != nil {logs.Error("Error WriteNewData"); return err}
    
    return nil
}

//Read files 
func GetConfFiles()(loadDataReturn map[string]string, err error) { 
    confFilePath := "./conf/main.conf"
    JSONconf, err := ioutil.ReadFile(confFilePath)
    if err != nil {
        logs.Error("utils/GetConf -> can't open Conf file: " + confFilePath)
        return nil, err
    }
    var anode map[string]map[string]string
    json.Unmarshal(JSONconf, &anode)
    return anode["files"], nil
}

//Generate a 16 bytes unique id
func Generate()(uuid string)  {
    b := make([]byte, 16)
    _, err := rand.Read(b)
    if err != nil {
        logs.Error(err)
    }
    uuid = fmt.Sprintf("%x-%x-%x-%x-%x",b[0:4], b[4:6], b[6:8], b[8:10], b[10:])
    return uuid
}

func LoadDefaultServerData(fileName string)(json map[string]string, err error){
    //Get full path
    loadData := map[string]map[string]string{}
    loadData["files"] = map[string]string{}
    loadData["files"][fileName] = ""
    loadData,err = GetConf(loadData)
    if err != nil {
        logs.Error("LoadDefaultServerData Error getting data from main.conf: "+err.Error())
        return nil, err
    }
    fileContent := make(map[string]string)
    rawData, err := ioutil.ReadFile(loadData["files"][fileName])
    if err != nil {
        logs.Error("LoadDefaultServerData Error reading file: "+err.Error())
        return nil,err
    }
    fileContent["fileContent"] = string(rawData)
    return fileContent,nil
}

func CopyFile(dstfolder string, srcfolder string, file string, BUFFERSIZE int64) (err error) {
    if BUFFERSIZE == 0{
        BUFFERSIZE = 1000
    }
    sourceFileStat, err := os.Stat(srcfolder+file)
    if err != nil {
        logs.Error("Error checking file at CopyFile function" + err.Error())
        return err
    }
    if !sourceFileStat.Mode().IsRegular() {
        logs.Error("%s is not a regular file.", sourceFileStat)
        return errors.New(sourceFileStat.Name()+" is not a regular file.")
    } 
    source, err := os.Open(srcfolder+file)
    if err != nil {
        return err
    }
    defer source.Close()
    _, err = os.Stat(dstfolder+file)
    if err == nil {
        return errors.New("File "+dstfolder+file+" already exists.")
    }
    destination, err := os.Create(dstfolder+file)
    if err != nil {
        logs.Error("Error Create =-> "+err.Error())
        return err
    }
    defer destination.Close()
    logs.Info("copy file -> "+srcfolder+file)
    logs.Info("to file -> "+dstfolder+file)
    buf := make([]byte, BUFFERSIZE)
    for {
        n, err := source.Read(buf)
        if err != nil && err != io.EOF {
            logs.Error("Error no EOF=-> "+err.Error())
            return err
        }
        if n == 0 {
            break
        }
        if _, err := destination.Write(buf[:n]); err != nil {
            logs.Error("Error Writing File: "+err.Error())
            return err
        }
    }
    return err
}

func RemoveFile(path string, file string)(err error){
    err = os.Remove(path+file)
    if err != nil {
        logs.Error("Error deleting file "+path+file+": "+err.Error())
        return err
    }
    return nil
}

//read data from main.conf
func GetConfArray(loadData map[string]map[string][]string)(loadDataReturn map[string]map[string][]string, err error) { 
    confFilePath := "conf/main.conf"
    jsonPathBpf, err := ioutil.ReadFile(confFilePath)
    if err != nil {
        logs.Error("utils/GetConf -> can't open Conf file: " + confFilePath)
        return nil, err
    }

    var anode map[string]map[string][]string
    json.Unmarshal(jsonPathBpf, &anode)

    for k,y := range loadData { 
        for y,_ := range y {
            if v, ok := anode[k][y]; ok {
                loadData[k][y] = v
            }else{
                loadData[k][y] = nil
            }
        }
    }
    return loadData, nil
}

func RunCommand(cmdtxt string, params string)(err error){
    cmd := exec.Command(cmdtxt, params)
    logs.Notice("utils run command -> Running command "+cmdtxt+"with params " + params)
    err = cmd.Run()
    if err != nil {
        logs.Error("utils run command -> "+err.Error())
        return err
    }
    return err
}

func FilelistPathByFile(path string, fileToSearch string)(files map[string][]byte, err error){
    pathMap:= make(map[string][]byte)
    err = filepath.Walk(path,
        func(file string, info os.FileInfo, err error) error {
        if err != nil {return err}
        
        if !info.IsDir() {
            pathSplit := strings.Split(file, "/")
            if strings.Contains(pathSplit[len(pathSplit)-1], fileToSearch){
                content, err := ioutil.ReadFile(file)
                if err != nil {logs.Error("Error filepath walk: "+err.Error()); return err}
                pathMap[pathSplit[len(pathSplit)-1]] = content
            }
        }
        return nil
    })
    if err != nil {logs.Error("Error filepath walk finish: "+err.Error()); return nil, err}

    return pathMap, nil
}