package models

import (
    "owlhnode/zeek"
    "github.com/astaxie/beego/logs"
)


func GetZeek() (status map[string]bool, err error) {
	logs.Info("Check whether Zeek is Installed!!")
	status,err = zeek.Installed()
	logs.Warn(status)
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

func DeployZeek() ( err error) {
	logs.Info("DeployZeek system into node server")
	err = zeek.DeployZeek()
    return err
}

func ChangeZeekMode(anode map[string]string) (err error) {
	err = zeek.ChangeZeekMode(anode)
    return err
}

func AddClusterValue(anode map[string]string) (err error) {
	err = zeek.AddClusterValue(anode)
    return err
}

func PingCluster()(data map[string]map[string]string, err error) {
	data, err = zeek.PingCluster()
    return data, err
}

func EditClusterValue(anode map[string]string) (err error) {
	err = zeek.EditClusterValue(anode)
    return err
}

func DeleteClusterValue(anode map[string]string) (err error) {
	err = zeek.DeleteClusterValue(anode)
    return err
}

func SyncCluster(anode map[string]string) (err error) {
	err = zeek.SyncCluster(anode, "cluster")
    return err
}