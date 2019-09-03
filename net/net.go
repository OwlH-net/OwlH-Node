package net 

import (
	"github.com/astaxie/beego/logs"
	"github.com/google/gopacket/pcap"
	"owlhnode/database"
	"owlhnode/utils"
	"io/ioutil"
	"regexp"
	"strings"
)
func ListInterfaces(interfaces []pcap.Interface)(netValues map[string]string) {
	data := make(map[string]string)
	for _, localInt := range interfaces {
		data[localInt.Name] = localInt.Name
        // logs.Info(localInt.Addresses)
        // logs.Info(localInt.Name)
	}
	return data
}
func ReadInterfaces()(devices []pcap.Interface, err error){
    devices, err = pcap.FindAllDevs()
    if err != nil {logs.Error("ReadInterfaces Error reading interfaces for Node: "+err.Error()); return nil, err}
    return devices, err
}

func GetNetworkData()(values map[string]string, err error) {
	//get interfaces
	interfaces, err := ReadInterfaces()
    if err != nil {
        return nil,err
    }
	data := ListInterfaces(interfaces)

	return data, nil
}

func LoadNetworkValuesSelected()(values map[string]map[string]string, err error) {
	//get current values selected for network
	values,err = ndb.LoadNodeconfigValues()
	if err != nil {logs.Error("LoadNetworkValuesSelected Error reading nodeconfig values for Node: "+err.Error()); return nil, err}
    return values,err
}

func UpdateNetworkInterface(data map[string]string) (err error) {
	//read suricata file
	suricataPath := map[string]map[string]string{}
    suricataPath["suriInit"] = map[string]string{}
    suricataPath["suriInit"]["path"] = ""
	suricataPath,err = utils.GetConf(suricataPath)
	if err != nil {logs.Error("UpdateNetworkInterface Error readding GetConf: "+err.Error())}
	suricataPathValue := suricataPath["suriInit"]["path"]
	
	reg := regexp.MustCompile(`interface=`)
	err = GetRegexpInterface(data,suricataPathValue, reg, "interface="+data["value"])
	if err != nil {logs.Error("UpdateNetworkInterface Error calling function GetRegexpInterface: "+err.Error())}
	reg = regexp.MustCompile(`INTERFACE="`)
	err = GetRegexpInterface(data,suricataPathValue, reg, "INTERFACE=\""+data["value"]+"\"")
	if err != nil {logs.Error("UpdateNetworkInterface Error calling function GetRegexpInterface: "+err.Error())}

	//read zeek file
	zeekPath := map[string]map[string]string{}
    zeekPath["loadDataZeekPath"] = map[string]string{}
    zeekPath["loadDataZeekPath"]["nodeConfig"] = ""
	zeekPath,err = utils.GetConf(zeekPath)
	if err != nil {logs.Error("UpdateNetworkInterface Error readding GetConf: "+err.Error())}
	zeekPathValue := zeekPath["loadDataZeekPath"]["nodeConfig"]

	reg = regexp.MustCompile(`interface=`)
	err = GetRegexpInterface(data,zeekPathValue, reg, "interface="+data["value"])
	if err != nil {logs.Error("UpdateNetworkInterface Error calling function GetRegexpInterface: "+err.Error())}
	reg = regexp.MustCompile(`INTERFACE="`)
	err = GetRegexpInterface(data,zeekPathValue, reg, "INTERFACE=\""+data["value"]+"\"")
	if err != nil {logs.Error("UpdateNetworkInterface Error calling function GetRegexpInterface: "+err.Error())}

	//update database with the new value
	err = ndb.ChangeNodeconfigValues(data["uuid"],data["param"],data["value"])
	if err != nil {logs.Error("UpdateNetworkInterface Error updating nodeconfig for Node: "+err.Error()); return err}
	
	//restart suricata
	err = utils.RestartSuricata()
	if err != nil {logs.Error("UpdateNetworkInterface Error restarting Suricata"); return err}
	
	//restart zeek
	err = utils.RestartZeek()
	if err != nil {logs.Error("UpdateNetworkInterface Error restarting Zeek"); return err}
	
    return nil
}

func GetRegexpInterface(data map[string]string, path string, regexpValue *regexp.Regexp, newLine string)(err error){
	//get first occurrence suricata
	input, err := ioutil.ReadFile(path)
	if err != nil {logs.Error("GetRegexpInterface Error readding zeek interface file: "+err.Error()); return err}
	lines := strings.Split(string(input), "\n")
	
	found := false
	for i := range lines {		
		regexpresult := regexpValue.FindStringSubmatch(lines[i])
		if regexpresult != nil && !found{
			lines[i] = newLine
			found = true
		}
	}
	output := strings.Join(lines, "\n")
	err = ioutil.WriteFile(path, []byte(output), 0644)
	if err != nil {logs.Error("GetRegexpInterface Error writting new interface for Node: "+err.Error()); return err}

	return nil
}

// func LoadNetworkValuesSuricata()(values map[string]map[string]string, err error) {
// 	GetNetworkData
//     return values,err
// }