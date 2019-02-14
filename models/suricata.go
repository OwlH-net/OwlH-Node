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