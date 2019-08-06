package models

import (
	"owlhnode/suricata"
)


func GetSuricata() (status map[string]bool, err error) {
	status,err = suricata.Installed()
	return status,err
}

func GetBPF() (bpf string, err error) {
	bpf,err = suricata.GetBPF()
    return bpf,err
}

func SetBPF(n map[string]string) (err error) {
    err = suricata.SetBPF(n)
    return err
}

func SyncRulesetFromMaster(n map[string][]byte) (err error) {
    err = suricata.SyncRulesetFromMaster(n)
    return err
}

func RunSuricata() (data string, err error) {
    data,err = suricata.RunSuricata()
    return data,err
}

func StopSuricata() (data string, err error) {
    data,err = suricata.StopSuricata()
    return data,err
}