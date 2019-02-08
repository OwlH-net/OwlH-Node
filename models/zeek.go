package models

import (
	"owlhnode/zeek"
)


func GetZeek() (status bool) {
	return zeek.Installed()
}