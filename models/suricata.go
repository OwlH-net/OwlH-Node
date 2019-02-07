package models

import (
	"owlhnode/suricata"
)


func GetSuricata() (status bool) {
	return suricata.Installed()
}