package models

import (
	"owlhnode/net"
)

func GetNetworkData() (data map[string]string, err error) {
	data,err = net.Main()
    return data,err
}