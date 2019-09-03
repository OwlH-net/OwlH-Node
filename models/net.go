package models

import (
	"owlhnode/net"
)

func GetNetworkData()(values map[string]string, err error) {
	values,err = net.GetNetworkData()
    return values,err
}

func LoadNetworkValuesSelected()(values map[string]map[string]string, err error) {
	values,err = net.LoadNetworkValuesSelected()
    return values,err
}

func UpdateNetworkInterface(data map[string]string) (err error) {
	err = net.UpdateNetworkInterface(data)
    return err
}

// func LoadNetworkValuesSuricata()(values map[string]map[string]string, err error) {
// 	values,err = net.LoadNetworkValuesSuricata()
//     return values,err
// }