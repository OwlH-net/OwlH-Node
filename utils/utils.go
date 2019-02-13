package utils

import (
    //"owlhnode/models"
    "encoding/json"
    //"strconv"
    //"github.com/astaxie/beego"
    "github.com/astaxie/beego/logs"
    "io/ioutil"
    "os"
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


