package models

import (
    "owlhnode/zeek"
    "owlhnode/changeControl"
    "github.com/astaxie/beego/logs")


func GetZeek() (status zeek.Zeek, err error) {
    logs.Info("Zeek Status!!")
    status = zeek.GetZeek()
    logs.Warn(status)
    // changecontrol.ChangeControlInsertData(err, "GetZeek")    
    return status,err
}

func SetZeek(zeekdata zeek.Zeek) (status zeek.Zeek, err error) {
    logs.Info("============")
    logs.Info("ZEEK - Set and manage Zeek")
    
    cc := map[string]string{}
    for key := range zeekdata.Extra {
        cc[key] = zeekdata.Extra[key]
    }
    status, err = zeek.SetZeek(zeekdata)
    logs.Warn(status)

    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }

    cc["actionDescription"] = "Set Zeek values and manage service status"
    cc["user"] = "admin"
    
    controlError := changecontrol.InsertChangeControl(cc)
    if controlError!=nil { logs.Error("Set Zeek params controlError: "+controlError.Error()) }

    // changecontrol.ChangeControlInsertData(err, "SetZeek")    
    return status,err

}

func RunZeek() (data string, err error) {
    logs.Info("Run RunZeek system into node server")
    data,err = zeek.StartZeek("")
    // changecontrol.ChangeControlInsertData(err, "RunZeek")    
    return data,err
}

func StartZeek(action string) (data string, err error) {
    logs.Info("============")
    logs.Info("ZEEK - Start Zeek with action " + action)
    
    cc := map[string]string{}
    cc["action"] = "PUT"
    cc["zeekAction"] = action
    logs.Info("Start Zeek using action " + action)
    data,err = zeek.StartZeek(action)
    if err!=nil { 
        cc["actionStatus"] = "error"
        cc["errorDescription"] = err.Error()
    }else{
        cc["actionStatus"] = "success"
    }
    cc["actionDescription"] = "Start Zeek with action " + action
    
    controlError := changecontrol.InsertChangeControl(cc)
    if controlError!=nil { logs.Error("Start Zeek controlError: "+controlError.Error()) }

    // changecontrol.ChangeControlInsertData(err, "StartZeek")    
    return data,err
}

func StopZeek() (data string, err error) {
    logs.Info("Stops StopZeek system into node server")
    data,err = zeek.StopZeek()
    // changecontrol.ChangeControlInsertData(err, "StopZeek")    
    return data,err
}

func DeployZeek() ( err error) {
    logs.Info("DeployZeek system into node server")
    err = zeek.DeployZeek()
    // changecontrol.ChangeControlInsertData(err, "DeployZeek")    
    return err
}

func ChangeZeekMode(anode map[string]string) (err error) {
    cc := anode
    logs.Info("============")
    logs.Info("ZEEK - ChangeZeekMode")
    for key :=range cc {
        logs.Info(key +" -> " + cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")
    err = zeek.ChangeZeekMode(anode)
    // changecontrol.ChangeControlInsertData(err, "ChangeZeekMode")    
    return err
}

func AddClusterValue(anode map[string]string) (err error) {
    cc := anode
    logs.Info("============")
    logs.Info("ZEEK - AddClusterValue")
    for key :=range cc {
        logs.Info(key +" -> ")
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

    err = zeek.AddClusterValue(anode)
    // changecontrol.ChangeControlInsertData(err, "AddClusterValue")    
    return err
}

func PingCluster()(data map[string]map[string]string, err error) {
    data, err = zeek.PingCluster()
    // changecontrol.ChangeControlInsertData(err, "PingCluster")    
    return data, err
}

func EditClusterValue(anode map[string]string) (err error) {
    cc := anode
    logs.Info("============")
    logs.Info("ZEEK - EditClusterValue")
    for key :=range cc {
        logs.Info(key +" -> " + cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

    err = zeek.EditClusterValue(anode)
    // changecontrol.ChangeControlInsertData(err, "EditClusterValue")    
    return err
}

func DeleteClusterValue(anode map[string]string) (err error) {
    cc := anode
    logs.Info("============")
    logs.Info("ZEEK - DeleteClusterValue")
    for key :=range cc {
        logs.Info(key +" -> " + cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

    err = zeek.DeleteClusterValue(anode)
    // changecontrol.ChangeControlInsertData(err, "DeleteClusterValue")    
    return err
}

func SyncCluster(anode map[string]string) (err error) {
    cc := anode
    logs.Info("============")
    logs.Info("ZEEK - SyncCluster")
    for key :=range cc {
        logs.Info(key +" -> " + cc[key])
    }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")


    err = zeek.SyncCluster(anode, "cluster")
    // changecontrol.ChangeControlInsertData(err, "SyncCluster")    
    return err
}

func SavePolicyFiles(files map[string]map[string][]byte) (err error) {
    cc := files
    logs.Info("============")
    logs.Info("Zeek - SaveConfigFile")
    for key :=range cc {
        logs.Info(key +" -> ")
    }

    err = zeek.SaveConfigFile(files)
    // changecontrol.ChangeControlInsertData(err, "SavePolicyFiles")    
    return err
}

func SyncClusterFile(anode map[string][]byte) (err error) {
    logs.Info("============")
    logs.Info("ZEEK - SyncClusterFile")
    // for key :=range cc {
    //     logs.Info(key +" -> " + cc[key])
    // }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

    err = zeek.SyncClusterFile(anode)
    // changecontrol.ChangeControlInsertData(err, "SyncClusterFile")    
    return err
}