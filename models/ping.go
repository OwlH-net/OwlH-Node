package models

import (
	"owlhnode/ping"
)

func PingService()(err error) {
	err = ping.PingService()
    return err
}

func DeployService()(err error) {
	err = ping.DeployService()
    return err
}