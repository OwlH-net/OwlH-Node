package models

import (
    "owlhnode/knownports"
)

func ShowPorts() (data map[string]string, err error) {
	data,err = knownports.ShowPorts()
	return data,err
}

func PingPorts() (data map[string]map[string]string ,err error) {
	data, err = knownports.PingPorts()
	return data, err
}

func ChangeMode(anode map[string]string) (err error) {
	err = knownports.ChangeMode(anode)
	return err
}

func ChangeStatus(anode map[string]string) (err error) {
	err = knownports.ChangeStatus(anode)
	return err
}