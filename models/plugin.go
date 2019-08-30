package models

import (
	"owlhnode/plugin"
)

func ChangeServiceStatus(anode map[string]string)(err error) {
	err = plugin.ChangeServiceStatus(anode)
    return err
}

func ChangeMainServiceStatus(anode map[string]string)(err error) {
	err = plugin.ChangeMainServiceStatus(anode)
    return err
}