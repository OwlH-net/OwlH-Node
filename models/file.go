package models

import (
	"owlhnode/file"
	"github.com/astaxie/beego/logs"
)

func SendFile(anode string) (data map[string]string, err error) {
    logs.Info("SendFile into Node file")
	data,err = file.SendFile(anode)
	//data,err = SendFile(anode)
	logs.Info("anode vuelta")
    return data,err
}

func SaveFile(data map[string]string) (err error) {
    logs.Info("SaveFile into Node file")
	err = file.SaveFile(data)
	// err = SaveFile(data)
    return err
}