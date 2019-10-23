package models

import (
    "owlhnode/knownports"
)

func ShowPorts() (data map[string]map[string]string, err error) {
	data,err = knownports.ShowPorts()
	return data,err
}

func ChangeMode(anode map[string]string) (err error) {
    var cc := anode
    logs.Info("============")
    logs.Info("PORTS - ChangeMode")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

	err = knownports.ChangeMode(anode)
	return err
}

func ChangeStatus(anode map[string]string) (err error) {
    var cc := anode
    logs.Info("============")
    logs.Info("PORTS - ChangeStatus")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

	err = knownports.ChangeStatus(anode)
	return err
}

func DeletePorts(anode map[string]string) (err error) {
    var cc := anode
    logs.Info("============")
    logs.Info("PORTS - DeletePorts")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

	err = knownports.DeletePorts(anode)
	return err
}

func DeleteAllPorts(anode map[string]string) (err error) {

    var cc := anode
    logs.Info("============")
    logs.Info("PORTS - DeleteAllPorts")
    for key :=range cc {
        logs.Info(key +" -> "+ cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

	err = knownports.DeleteAllPorts()
	return err
}

func PingPorts() (data map[string]map[string]string ,err error) {
	data, err = knownports.PingPorts()
	return data, err
}