package net 

import (
	"github.com/astaxie/beego/logs"
    "github.com/google/gopacket/pcap"
)
func ListInterfaces(interfaces []pcap.Interface)(netValues map[string]string) {
	data := make(map[string]string)
	for _, localInt := range interfaces {
		data[localInt.Name] = localInt.Name
        // logs.Info(localInt.Addresses)
        // logs.Info(localInt.Name)
		// logs.Info("---")
	}
	return data
}
func ReadInterfaces()(devices []pcap.Interface, err error){
    devices, err = pcap.FindAllDevs()
    if err != nil {
        return nil, err
    }
    return devices, err
}

func Main()(data map[string]string, err error) {
    interfaces, err := ReadInterfaces()
    if err != nil {
        logs.Error(err)
        return nil,err
    }
	data = ListInterfaces(interfaces)
	return data, nil
}

func UpdateNetworkInterface(data map[string]string) (err error) {
	logs.Notice(data)
    return nil
}