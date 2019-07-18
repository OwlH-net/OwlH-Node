package models

import (
	"owlhnode/net"
)

func GetNetworkData() (data map[string]string, err error) {
	data,err = net.Main()
    return data,err
}

func UpdateNetworkInterface(data map[string]string) (err error) {
	err = net.UpdateNetworkInterface(data)
    return err
}