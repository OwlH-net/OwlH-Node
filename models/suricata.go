package models

import (
	"owlhnode/suricata"
	"github.com/astaxie/beego/logs"
)


func GetSuricata() (status map[string]bool, err error) {
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

func RunSuricata() (data string, err error) {
    logs.Info("Run suricata system into node server")
    data,err = suricata.RunSuricata()
    return data,err
}

func StopSuricata() (data string, err error) {
    logs.Info("Stops suricata system into node server")
    data,err = suricata.StopSuricata()
    return data,err
}