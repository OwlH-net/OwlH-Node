package dataflow

import (
    "github.com/astaxie/beego/logs"
    "owlhnode/database"
    "owlhnode/utils"
    "strings"
    "errors"
)

func ChangeDataflowValues(anode map[string]string)(err error) {
    err = ndb.ChangeDataflowValues(anode["FlowUUID"], anode["param"], anode["value"])
    if (err != nil){ logs.Error("ChangeDataflowValues UPDATE error for update dataflow values: "+err.Error()); return err}
    return nil
}

func LoadDataflowValues()(data map[string]map[string]string, err error) {
    data,err = ndb.LoadDataflowValues()
    if (err != nil){ logs.Error("LoadDataflowValues error loading dataflow values: "+err.Error()); return nil,err}
    return data,nil
}

func SaveSocketToNetwork(anode map[string]string) (err error) {
    data,err := ndb.LoadDataflowValues()
    anodeCapitalized := strings.Replace(strings.ToLower(anode["name"]), " ", "_", -1)
    exists := false

    for x := range data {
        dbCapitalized := strings.Replace(strings.ToLower(data[x]["name"]), " ", "_", -1)
        if (data[x]["type"] == "sockettonetwork" && anodeCapitalized == dbCapitalized){
            exists = true
            break
        }        
    }

    if (exists){
        return errors.New("Name in use. Use other name.")
    }else{
        uuid := utils.Generate()
        err = ndb.InsertDataflowValues(uuid, "name", anode["name"])
        if (err != nil){ logs.Error("SaveSocketToNetwork error inserting name dataflow values: "+err.Error()); return err}
        err = ndb.InsertDataflowValues(uuid, "cert", anode["cert"])
        if (err != nil){ logs.Error("SaveSocketToNetwork error inserting cert dataflow values: "+err.Error()); return err}
        err = ndb.InsertDataflowValues(uuid, "interface", anode["interface"])
        if (err != nil){ logs.Error("SaveSocketToNetwork error inserting interface dataflow values: "+err.Error()); return err}
        err = ndb.InsertDataflowValues(uuid, "port", anode["port"])
        if (err != nil){ logs.Error("SaveSocketToNetwork error inserting port dataflow values: "+err.Error()); return err}
        err = ndb.InsertDataflowValues(uuid, "type", "sockettonetwork")
        if (err != nil){ logs.Error("SaveSocketToNetwork error inserting sockettonetwork dataflow values: "+err.Error()); return err}
    }
    
    return nil
}

func SaveNewLocal(anode map[string]string)(err error) {
    data,err := ndb.LoadDataflowValues()
    anodeCapitalized := strings.Replace(strings.ToLower(anode["name"]), " ", "_", -1)
    exists := false

    for x := range data {
        dbCapitalized := strings.Replace(strings.ToLower(data[x]["name"]), " ", "_", -1)
        if (data[x]["type"] == "networknewlocal" && anodeCapitalized == dbCapitalized){
            exists = true
            break
        }        
    }

    if (exists){
        return errors.New("Name in use. Use other name.")
    }else{
        uuid := utils.Generate()
        err = ndb.InsertDataflowValues(uuid, "name", anode["name"])
        if (err != nil){ logs.Error("SaveSocketToNetwork error inserting name dataflow values: "+err.Error()); return err}
        err = ndb.InsertDataflowValues(uuid, "mtu", anode["mtu"])
        if (err != nil){ logs.Error("SaveSocketToNetwork error inserting cert dataflow values: "+err.Error()); return err}
        err = ndb.InsertDataflowValues(uuid, "type", "networknewlocal")
        if (err != nil){ logs.Error("SaveSocketToNetwork error inserting sockettonetwork dataflow values: "+err.Error()); return err}
    }

    return nil
}

func SaveVxLAN(anode map[string]string)(err error) {
    data,err := ndb.LoadDataflowValues()
    anodeCapitalized := strings.Replace(strings.ToLower(anode["interface"]), " ", "_", -1)
    exists := false

    for x := range data {
        dbCapitalized := strings.Replace(strings.ToLower(data[x]["name"]), " ", "_", -1)
        if (data[x]["type"] == "networkvxlan" && anodeCapitalized == dbCapitalized){
            exists = true
            break
        }        
    }

    if (exists){
        return errors.New("Name in use. Use other name.")
    }else{
        uuid := utils.Generate()
        err = ndb.InsertDataflowValues(uuid, "name", anode["interface"])
        if (err != nil){ logs.Error("SaveSocketToNetwork error inserting interface dataflow values: "+err.Error()); return err}
        err = ndb.InsertDataflowValues(uuid, "lanIp", anode["lanIp"])
        if (err != nil){ logs.Error("SaveSocketToNetwork error inserting interface dataflow values: "+err.Error()); return err}
        err = ndb.InsertDataflowValues(uuid, "localIp", anode["localIp"])
        if (err != nil){ logs.Error("SaveSocketToNetwork error inserting interface dataflow values: "+err.Error()); return err}
        err = ndb.InsertDataflowValues(uuid, "portIp", anode["portIp"])
        if (err != nil){ logs.Error("SaveSocketToNetwork error inserting interface dataflow values: "+err.Error()); return err}
        err = ndb.InsertDataflowValues(uuid, "type", anode["type"])
        if (err != nil){ logs.Error("SaveSocketToNetwork error inserting interface dataflow values: "+err.Error()); return err}
        err = ndb.InsertDataflowValues(uuid, "baseInterface", anode["baseInterface"])
        if (err != nil){ logs.Error("SaveSocketToNetwork error inserting interface dataflow values: "+err.Error()); return err}
    }
    
    return nil
}

func SaveSocketToNetworkSelected(anode map[string]string)(err error) {
    return nil
}

func DeleteDataFlowValueSelected(anode map[string]string)(err error) {
    err = ndb.DeleteDataFlowValueSelected(anode["uuidNode"])
    if (err != nil){ logs.Error("DeleteDataFlowValueSelected error deleting a socket->network: "+err.Error()); return err}
    return nil
}