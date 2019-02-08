package models

import (
	"owlhnode/wazuh"
)


func GetWazuh() (status bool) {
	return wazuh.Installed()
}