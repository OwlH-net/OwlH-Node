package models

import (
    "owlhnode/zeek"
    "github.com/astaxie/beego/logs"
)


func GetZeek() (status map[string]bool, err error) {
	logs.Info("Check whether Zeek is Installed!!")
    status,err = zeek.Installed()
	return status,err
}

func RunZeek() (data string, err error) {
    logs.Info("Run RunZeek system into node server")
    data,err = zeek.RunZeek()
    return data,err
}

func StopZeek() (data string, err error) {
    logs.Info("Stops StopZeek system into node server")
    data,err = zeek.StopZeek()
    return data,err
}