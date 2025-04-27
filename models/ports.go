package models

import (
	"github.com/OwlH-net/OwlH-Node/knownports"

	changecontrol "github.com/OwlH-net/OwlH-Node/changeControl"
	"github.com/astaxie/beego/logs"
)

func ShowPorts(username string) (data map[string]map[string]string, err error) {
	data, err = knownports.ShowPorts()
	//changecontrol.ChangeControlInsertData(err, "ShowPorts")
	return data, err
}

func ChangeMode(anode map[string]string, username string) (err error) {
	cc := anode
	logs.Info("============")
	logs.Info("PORTS - ChangeMode")
	for key := range cc {
		logs.Info(key + " -> " + cc[key])
	}
	delete(anode, "action")
	delete(anode, "controller")
	delete(anode, "router")

	err = knownports.ChangeMode(anode)

	if err != nil {
		cc["actionStatus"] = "error"
		cc["errorDescription"] = err.Error()
	} else {
		cc["actionStatus"] = "success"
	}
	cc["username"] = username
	cc["actionDescription"] = "ChangeMode knownports Mode"

	changecontrol.InsertChangeControl(cc)
	// changecontrol.ChangeControlInsertData(err, "ChangeMode")
	return err
}

func ChangeStatus(anode map[string]string, username string) (err error) {
	cc := anode
	logs.Info("============")
	logs.Info("PORTS - ChangeStatus")
	for key := range cc {
		logs.Info(key + " -> " + cc[key])
	}
	delete(anode, "action")
	delete(anode, "controller")
	delete(anode, "router")

	err = knownports.ChangeStatus(anode)
	// changecontrol.ChangeControlInsertData(err, "ChangeStatus")

	if err != nil {
		cc["actionStatus"] = "error"
		cc["errorDescription"] = err.Error()
	} else {
		cc["actionStatus"] = "success"
	}
	cc["username"] = username
	cc["actionDescription"] = "Start/stop knownports plugin"

	changecontrol.InsertChangeControl(cc)

	return err
}

func DeletePorts(anode map[string]string, username string) (err error) {
	cc := anode
	logs.Info("============")
	logs.Info("PORTS - DeletePorts")
	for key := range cc {
		logs.Info(key + " -> " + cc[key])
	}
	delete(anode, "action")
	delete(anode, "controller")
	delete(anode, "router")

	err = knownports.DeletePorts(anode)

	if err != nil {
		cc["actionStatus"] = "error"
		cc["errorDescription"] = err.Error()
	} else {
		cc["actionStatus"] = "success"
	}
	cc["username"] = username
	cc["actionDescription"] = "Clean Knownports data base"

	changecontrol.InsertChangeControl(cc)

	// changecontrol.ChangeControlInsertData(err, "DeletePorts")
	return err
}

func DeleteAllPorts(anode map[string]string, username string) (err error) {
	cc := anode
	logs.Info("============")
	logs.Info("PORTS - DeleteAllPorts")
	for key := range cc {
		logs.Info(key + " -> " + cc[key])
	}
	delete(anode, "action")
	delete(anode, "controller")
	delete(anode, "router")

	err = knownports.DeleteAllPorts()

	if err != nil {
		cc["actionStatus"] = "error"
		cc["errorDescription"] = err.Error()
	} else {
		cc["actionStatus"] = "success"
	}
	cc["username"] = username
	cc["actionDescription"] = "Clean all known ports"

	changecontrol.InsertChangeControl(cc)

	// changecontrol.ChangeControlInsertData(err, "DeleteAllPorts")
	return err
}

func PingPorts(username string) (data map[string]map[string]string, err error) {
	data, err = knownports.PingPorts()
	// changecontrol.ChangeControlInsertData(err, "PingPorts")
	return data, err
}
