package models

import (
	"owlhnode/suricata"
)


func GetSuricata() (status bool) {
	return suricata.Installed()
}

func GetBPF() (currentbpf string) {
    return suricata.GetBPF()
}

func SetBPF(newbpf string) (status bool) {
    return suricata.SetBPF(newbpf)
}