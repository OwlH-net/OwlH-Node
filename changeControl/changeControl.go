package changecontrol

import (
    "github.com/astaxie/beego/logs"
	"owlhnode/database"
	"owlhnode/utils"
	"time"
)

func GetChangeControlNode()(data map[string]map[string]string, err error) {
	data, err = ndb.GetChangeControlNode(); if err != nil{logs.Error("Error getting ChangeControl database values: "+err.Error()); return nil,err}
    return data, err
}

func InsertChangeControl(values map[string]string)(err error){
    uuid:= utils.Generate()
    currentTime := time.Now()
    timeFormated := currentTime.Format("2006-01-02T15:04:05")

    node,err := ndb.GetNodeData()
    for r := range node{
        err = ndb.InsertChangeControl(uuid, "deviceName", node[r]["name"]); if err != nil{logs.Error("Error inserting ChangeControl database value: "+err.Error()); return err}
        err = ndb.InsertChangeControl(uuid, "deviceIP", node[r]["ip"]); if err != nil{logs.Error("Error inserting ChangeControl database value: "+err.Error()); return err}
        err = ndb.InsertChangeControl(uuid, "devicePort", node[r]["port"]); if err != nil{logs.Error("Error inserting ChangeControl database value: "+err.Error()); return err}
    }
    err = ndb.InsertChangeControl(uuid, "nodeServerUuid", values["uuid"]); if err != nil{logs.Error("Error inserting ChangeControl database value: "+err.Error()); return err}
    err = ndb.InsertChangeControl(uuid, "time", timeFormated); if err != nil{logs.Error("Error inserting ChangeControl database value: "+err.Error()); return err}
    for x := range values {
        err = ndb.InsertChangeControl(uuid, x, values[x]); if err != nil{logs.Error("Error inserting ChangeControl database value: "+err.Error()); return err}
    }

    return nil
}