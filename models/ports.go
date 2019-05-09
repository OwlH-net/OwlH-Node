package models

import (
    "owlhnode/knownports"
)

func ShowPorts() (data map[string]string, err error) {
	data,err = knownports.ShowPorts()
	return data,err
}

func ChangeMode() (err error) {
	err = knownports.ChangeMode()
	return err
}

func ChangeStatus() (err error) {
	err = knownports.ChangeStatus()
	return err
}

func PingPorts() (data map[string]map[string]string ,err error) {
	data, err = knownports.PingPorts()
	return data, err
}