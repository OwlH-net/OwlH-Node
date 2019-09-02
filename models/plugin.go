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

func DeleteService(anode map[string]string)(err error) {
	err = plugin.DeleteService(anode)
    return err
}

func AddPluginService(anode map[string]string) (err error) {
    err = plugin.AddPluginService(anode)
    return err
}

