package knownports

import (
	ndb "github.com/OwlH-net/OwlH-Node/database"
	"github.com/astaxie/beego/logs"
)

func PingPorts() (data map[string]map[string]string, err error) {
	plugins, err := ndb.GetPlugins()
	if err != nil {
		logs.Error("PingPorts knownports Error : %s", err.Error())
		return nil, err
	}
	return plugins, nil
}

func ShowPorts() (data map[string]map[string]string, err error) {
	ports, err := ndb.LoadPortsData()
	if err != nil {
		logs.Error("ShowPorts knownports Error : %s", err.Error())
		return nil, err
	}
	return ports, nil
}

func ChangeStatus(anode map[string]string) (err error) {
	if anode["value"] == "Enable" {
		Init()
	}
	err = ndb.UpdatePluginValue(anode["status"], "status", anode["plugin"])
	if err != nil {
		logs.Error("ChangeStatus knownports Error : %s", err.Error())
		return err
	}
	return nil
}

func ChangeMode(anode map[string]string) (err error) {
	err = ndb.UpdatePluginValue(anode["mode"], "mode", "knownports")
	if err != nil {
		logs.Error("ChangeMode knownports Error : %s", err.Error())
		return err
	}
	return nil
}

func DeletePorts(ports map[string]string) (err error) {
	err = ndb.DeletePort(ports)
	if err != nil {
		logs.Error("DeletePorts knownports Error : %s", err.Error())
		return err
	}

	anode := make(map[string]string)
	anode["plugin"] = "knownports"
	anode["status"] = "Reload"
	_ = ChangeStatus(anode)

	return nil
}

func DeleteAllPorts() (err error) {
	err = ndb.DeleteAllPorts()
	if err != nil {
		logs.Error("DeleteAllPorts knownports Error : %s", err.Error())
		return err
	}

	anode := make(map[string]string)
	anode["plugin"] = "knownports"
	anode["status"] = "Reload"
	_ = ChangeStatus(anode)

	return nil
}
