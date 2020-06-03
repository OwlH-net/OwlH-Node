package net 

import (
    "github.com/astaxie/beego/logs"
    "github.com/google/gopacket/pcap"
    "owlhnode/database"
    "owlhnode/zeek"
    "time"
    "owlhnode/utils"
    // "owlhnode/zeek"
    // "owlhnode/suricata"
    // "io/ioutil"
    // "regexp"
    // "strings"
)
func ListInterfaces(interfaces []pcap.Interface)(netValues map[string]string) {
    //historical log
    uuid := utils.Generate()
    currentTime := time.Now()
    timeFormated := currentTime.Format("2006-01-02T15:04:05")
    _ = ndb.InsertPluginCommand(uuid, "date", timeFormated)
    _ = ndb.InsertPluginCommand(uuid, "type", "Interfaces")
    _ = ndb.InsertPluginCommand(uuid, "action", "ListInterfaces")
    _ = ndb.InsertPluginCommand(uuid, "description", "Get system interfaces")

    data := make(map[string]string)
    for _, localInt := range interfaces {
        data[localInt.Name] = localInt.Name
        // logs.Info(localInt.Addresses)
        // logs.Info(localInt.Name)
    }
    if data == nil || len(data)<=0{
        _ = ndb.InsertPluginCommand(uuid, "status", "Error")
        _ = ndb.InsertPluginCommand(uuid, "output", "ListInterfaces is empty")
    }
    _ = ndb.InsertPluginCommand(uuid, "status", "Success")
    _ = ndb.InsertPluginCommand(uuid, "output", "ListInterfaces get interfaces successfully")
    return data
}
func ReadInterfaces()(devices []pcap.Interface, err error){
    //historical log
    uuid := utils.Generate()
    currentTime := time.Now()
    timeFormated := currentTime.Format("2006-01-02T15:04:05")
    _ = ndb.InsertPluginCommand(uuid, "date", timeFormated)
    _ = ndb.InsertPluginCommand(uuid, "type", "Interfaces")
    _ = ndb.InsertPluginCommand(uuid, "action", "ReadInterfaces")
    _ = ndb.InsertPluginCommand(uuid, "description", "Get node interfaces")

    devices, err = pcap.FindAllDevs()
    if err != nil {
        _ = ndb.InsertPluginCommand(uuid, "status", "Error")
        _ = ndb.InsertPluginCommand(uuid, "output", "ReadInterfaces error: "+err.Error())
        logs.Error("ReadInterfaces Error reading interfaces for Node: "+err.Error())
        return nil, err
    }

    _ = ndb.InsertPluginCommand(uuid, "status", "Success")
    _ = ndb.InsertPluginCommand(uuid, "output", "ReadInterfaces successfully")
    return devices, err
}

func GetNetworkData()(values map[string]string, err error) {
    //historical log
    uuid := utils.Generate()
    currentTime := time.Now()
    timeFormated := currentTime.Format("2006-01-02T15:04:05")
    _ = ndb.InsertPluginCommand(uuid, "date", timeFormated)
    _ = ndb.InsertPluginCommand(uuid, "type", "Interfaces")
    _ = ndb.InsertPluginCommand(uuid, "action", "GetNetworkData")
    _ = ndb.InsertPluginCommand(uuid, "description", "Get node interfaces")
    
    //get interfaces
    interfaces, err := ReadInterfaces()
    if err != nil {
        _ = ndb.InsertPluginCommand(uuid, "status", "Error")
        _ = ndb.InsertPluginCommand(uuid, "output", "GetNetworkData error on function ReadInterfaces: "+err.Error())
        return nil,err
    }
    data := ListInterfaces(interfaces)

    _ = ndb.InsertPluginCommand(uuid, "status", "Success")
    _ = ndb.InsertPluginCommand(uuid, "output", "StartSuricataMainConf successfully")

    return data, nil
}

func LoadNetworkValuesSelected()(values map[string]map[string]string, err error) {
    //get current values selected for network
    values,err = ndb.LoadNodeconfigValues()
    if err != nil {logs.Error("LoadNetworkValuesSelected Error reading nodeconfig values for Node: "+err.Error()); return nil, err}
    return values,err
}

func UpdateNetworkInterface(data map[string]string) (err error) {
    // //read suricata file
    // suricataPath := map[string]map[string]string{}
    // suricataPath["suriInit"] = map[string]string{}
    // suricataPath["suriInit"]["path"] = ""
    // suricataPath,err = utils.GetConf(suricataPath)
    // if err != nil {logs.Error("UpdateNetworkInterface Error readding GetConf: "+err.Error())}
    // suricataPathValue := suricataPath["suriInit"]["path"]
    
    // reg := regexp.MustCompile(`interface=`)
    // err = GetRegexpInterface(data,suricataPathValue, reg, "interface="+data["value"])
    // if err != nil {logs.Error("UpdateNetworkInterface Error calling function GetRegexpInterface: "+err.Error())}
    // reg = regexp.MustCompile(`INTERFACE="`)
    // err = GetRegexpInterface(data,suricataPathValue, reg, "INTERFACE=\""+data["value"]+"\"")
    // if err != nil {logs.Error("UpdateNetworkInterface Error calling function GetRegexpInterface: "+err.Error())}

    //read zeek file
    // zeekPath := map[string]map[string]string{}
    // zeekPath["loadDataZeekPath"] = map[string]string{}
    // zeekPath["loadDataZeekPath"]["nodeConfig"] = ""
    // zeekPath,err = utils.GetConf(zeekPath)
    // if err != nil {logs.Error("UpdateNetworkInterface Error readding GetConf: "+err.Error())}
    // zeekPathValue := zeekPath["loadDataZeekPath"]["nodeConfig"]

    // reg := regexp.MustCompile(`interface=`)
    // err = GetRegexpInterface(data,zeekPathValue, reg, "interface="+data["value"])
    // if err != nil {logs.Error("UpdateNetworkInterface Error calling function GetRegexpInterface: "+err.Error())}
    // reg = regexp.MustCompile(`INTERFACE="`)
    // err = GetRegexpInterface(data,zeekPathValue, reg, "INTERFACE=\""+data["value"]+"\"")
    // if err != nil {logs.Error("UpdateNetworkInterface Error calling function GetRegexpInterface: "+err.Error())}

    //update database with the new value
    // err = ndb.ChangeNodeconfigValues(data["uuid"],data["param"],data["value"])
    // if err != nil {logs.Error("UpdateNetworkInterface Error updating nodeconfig for Node: "+err.Error()); return err}
    
    // //restart suricata
    // _,err = suricata.StopSuricata()
    // if err != nil {logs.Error("UpdateNetworkInterface Error stopping Suricata"); return err}
    // _,err = suricata.RunSuricata()
    // if err != nil {logs.Error("UpdateNetworkInterface Error running Suricata"); return err}
    
    //update zeek db interface
    err = ndb.UpdatePluginValue(data["service"], "interface", data["value"]); if err != nil {logs.Error("UpdateNetworkInterface Zeek interface update Error: "+err.Error()); return err}
    err = zeek.SyncCluster(data, "standalone"); if err != nil {logs.Error("UpdateNetworkInterface Zeek update interface and node.cfg Error: "+err.Error()); return err}
    // //restart zeek
    // err = zeek.DeployZeek()
    // if err != nil {logs.Error("UpdateNetworkInterface Error restarting Zeek"); return err}
    
    return nil
}

// func GetRegexpInterface(data map[string]string, path string, regexpValue *regexp.Regexp, newLine string)(err error){
//     //get first occurrence suricata
//     input, err := ioutil.ReadFile(path)
//     if err != nil {logs.Error("GetRegexpInterface Error readding zeek interface file: "+err.Error()); return err}
//     lines := strings.Split(string(input), "\n")
    
//     found := false
//     for i := range lines {        
//         regexpresult := regexpValue.FindStringSubmatch(lines[i])
//         if regexpresult != nil && !found{
//             lines[i] = newLine
//             found = true
//         }
//     }
//     output := strings.Join(lines, "\n")
//     err = ioutil.WriteFile(path, []byte(output), 0644)
//     if err != nil {logs.Error("GetRegexpInterface Error writting new interface for Node: "+err.Error()); return err}

//     return nil
// }

// func LoadNetworkValuesSuricata()(values map[string]map[string]string, err error) {
//     GetNetworkData
//     return values,err
// }