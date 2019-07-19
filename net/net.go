package net 

import (
	"github.com/astaxie/beego/logs"
	"github.com/google/gopacket/pcap"
	"owlhnode/database"
)
func ListInterfaces(interfaces []pcap.Interface)(netValues map[string]string) {
	data := make(map[string]string)
	for _, localInt := range interfaces {
		data[localInt.Name] = localInt.Name
        // logs.Info(localInt.Addresses)
        // logs.Info(localInt.Name)
	}
	return data
}
func ReadInterfaces()(devices []pcap.Interface, err error){
    devices, err = pcap.FindAllDevs()
    if err != nil {logs.Error("ReadInterfaces Error reading interfaces for Node"); return nil, err}
    return devices, err
}

func GetNetworkData()(values map[string]string, err error) {
	//get interfaces
	interfaces, err := ReadInterfaces()
    if err != nil {
        return nil,err
    }
	data := ListInterfaces(interfaces)

	return data, nil
}

func UpdateNetworkInterface(data map[string]string) (err error) {
	err = ndb.ChangeNodeconfigValues(data["uuid"],data["param"],data["value"])
	if err != nil {logs.Error("UpdateNetworkInterface Error updating nodeconfig for Node"); return err}
    return nil
}

func LoadNetworkValuesSelected()(values map[string]map[string]string, err error) {
	//get current values selected for network
	values,err = ndb.LoadNodeconfigValues()
	if err != nil {logs.Error("LoadNetworkValuesSelected Error reading nodeconfig values for Node"); return nil, err}
    return values,err
}