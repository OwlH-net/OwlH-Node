package models

import (
	"owlhnode/suricata"
	"github.com/astaxie/beego/logs"
)


func GetSuricata() (status bool) {
	return suricata.Installed()
}
/*
func GetBPF() (currentbpf string) {
    return suricata.GetBPF()
}
*/
func SetBPF(n map[string]string) (data string, err error) {
    logs.Info("Set Suricata BPF into Node file - %s",n)
    data,err = suricata.SetBPF(n)
    return data,err
}

func RetrieveFile(n map[string][]byte) (err error) {
    logs.Info("Set ruleset file into Node file")
    err = suricata.RetrieveFile(n)
    return err
}