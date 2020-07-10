package changecontrol

import (
    "github.com/astaxie/beego/logs"
    "owlhnode/database"
    "owlhnode/utils"
    "time"
)

func GetChangeControlNode() (data map[string]map[string]string, err error) {
    data, err = ndb.GetChangeControlNode()
    if err != nil {
        logs.Error("Error getting ChangeControl database values: " + err.Error())
        return nil, err
    }
    return data, err
}

func InsertChangeControl(cc map[string]string) {
    uuid := utils.Generate()
    currentTime := time.Now()
    timeFormated := currentTime.Format("2006-01-02T15:04:05")

    node, err := ndb.GetNodeData()
    if err != nil {
        return
    }
    if err == nil {
        for r := range node {
            ndb.InsertChangeControl(uuid, "deviceName", node[r]["name"])
            ndb.InsertChangeControl(uuid, "deviceIP", node[r]["ip"])
            ndb.InsertChangeControl(uuid, "devicePort", node[r]["port"])
        }
    }
    ndb.InsertChangeControl(uuid, "nodeServerUuid", cc["uuid"])
    ndb.InsertChangeControl(uuid, "time", timeFormated)
    // ndb.InsertChangeControl(uuid, "user", "admin")
    for x := range cc {
        ndb.InsertChangeControl(uuid, x, cc[x])
    }
}

func InsertChangeControlExtra(cc map[string]map[string]string) {
    // uuid:= utils.Generate()
    // currentTime := time.Now()
    // timeFormated := currentTime.Format("2006-01-02T15:04:05")

}

func InsertChangeControlByte(cc map[string][]byte) {
    // uuid:= utils.Generate()
    // currentTime := time.Now()
    // timeFormated := currentTime.Format("2006-01-02T15:04:05")

}

func ChangeControlInsertData(err error, desc string) {
    //check error
    return
    // n := make(map[string]string)
    // if err!=nil {
    //     n["actionStatus"] = "error"
    //     n["errorDescription"] = err.Error()
    // }else{
    //     n["actionStatus"] = "success"
    // }
    // n["action"] = "POST"
    // n["actionDescription"] = desc

    // //add incident
    // var controlError error
    // controlError = InsertChangeControl(n)
    // if controlError!=nil { logs.Error(desc+" controlError: "+controlError.Error()) }
}
