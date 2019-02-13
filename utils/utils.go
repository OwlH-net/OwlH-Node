package utils

import (
    //"owlhnode/models"
    "encoding/json"
    "strconv"
    //"github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
    "io/ioutil"
    "os"
    "time"
)


func GetConf(param string)(value string) {
    confFilePath := "/etc/owlh/conf/main.conf"
    confFile, err := os.Open(confFilePath)

    if err != nil {
        logs.Error ("utils/GetConf -> can't open Conf file -> " + confFilePath)
    }
    defer confFile.Close()

    byteValue, _ := ioutil.ReadAll(confFile)

    var config map[string]string
    json.Unmarshal([]byte(byteValue), &config)

    if value, exists := config[param]; exists {
        return value
    } else {
        return "ERROR"
    }
}

func backupFile(file string) (err error) {
    in, err := os.Open(file)
    if err != nil {
        return err
    }
    defer in.Close()
    t := time.Now()
    dst = in+"-"+strconv.FormatInt(t.Unix(), 10)
    out, err := os.Create(dst)
    if err != nil {
        return err
    }
    defer func() {
        cerr := out.Close()
        if err == nil {
            err = cerr
        }
    }()
    if _, err = io.Copy(out, in); err != nil {
        return err
    }
    err = out.Sync()
    return nil
}
