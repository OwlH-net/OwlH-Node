package models

import (
	"owlhnode/stap"
	"github.com/astaxie/beego/logs"
)

func AddServer(elem map[string]string) (err error) {
    logs.Info("AddServer stap into Node file")
	err = stap.AddServer(elem)
    return err
}