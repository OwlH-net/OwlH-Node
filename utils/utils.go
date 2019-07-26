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
    "crypto/rand"
)


//read data from main.conf
func GetConf(loadData map[string]map[string]string)(loadDataReturn map[string]map[string]string, err error) { 
    confFilePath := "conf/main.conf"
    jsonPathBpf, err := ioutil.ReadFile(confFilePath)
    if err != nil {
        logs.Error("utils/GetConf -> can't open Conf file: " + confFilePath)
        return nil, err
	}

    var anode map[string]map[string]string
    json.Unmarshal(jsonPathBpf, &anode)

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
	if err != nil {
		logs.Error("Error writing new BPF data into file: "+err.Error())
		return err
	}
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
    t := time.Now()
    newFile := fileName+"-"+strconv.FormatInt(t.Unix(), 10)
    srcFolder := path+fileName
	destFolder := path+newFile
	//check if file exist
	if _, err := os.Stat(srcFolder); os.IsNotExist(err) {
		return nil
	}else{
		cpCmd := exec.Command("cp", srcFolder, destFolder)
		err = cpCmd.Run()
		if err != nil{
			logs.Error("utils.BackupFile Error exec cmd command: "+err.Error())
			return err
		}
	}
	return nil
}

//write data on a file
func WriteNewDataOnFile(path string, data []byte)(err error){
    err = ioutil.WriteFile(path, data, 0644)
	if err != nil {
        logs.Error("Error WriteNewData")
		return err
	}
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
	logs.Debug(anode["files"])
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

func RestartSuricata()(err error){
	//stop suricata
	suricataStop := map[string]map[string]string{}
	suricataStop["suriStop"] = map[string]string{}
	suricataStop["suriStop"]["stop"] = ""
	suricataStop["suriStop"]["param"] = ""
	suricataStop["suriStop"]["command"] = ""
	suricataStop,err = GetConf(suricataStop)
	if err != nil {logs.Error("RestartSuricata Error readding GetConf for stop: "+err.Error()); return err}
	
	_,err = exec.Command(suricataStop["suriStop"]["command"], suricataStop["suriStop"]["param"], suricataStop["suriStop"]["stop"]).Output()
	if err != nil{logs.Error("RestartSuricata Error exec cmd command: "+err.Error()); return err}

	//run suricata
	suricatastart := map[string]map[string]string{}
	suricatastart["suriStart"] = map[string]string{}
	suricatastart["suriStart"]["start"] = ""
	suricatastart["suriStart"]["param"] = ""
	suricatastart["suriStart"]["command"] = ""
	suricatastart,err = GetConf(suricatastart)
	if err != nil {logs.Error("RestartSuricata Error readding GetConf for start: "+err.Error()); return err}
	
	_,err = exec.Command(suricatastart["suriStart"]["command"], suricatastart["suriStart"]["param"], suricatastart["suriStart"]["start"]).Output()
	if err != nil{logs.Error("RestartSuricata Error exec cmd command: "+err.Error()); return err}

	return nil
}

func RunCommand(cmdtxt, params string)err error{
    cmd := exec.Command(cmdtxt, params)
    logs.Notice("utils run command -> Running command "+cmdtxt+"with params " + params)
    err = cmd.Run()
    if err != nil {
        logs.Error("utils run command -> "+err.Error())
        return nil, err
    }
    return nil
}
