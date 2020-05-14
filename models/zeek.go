package models

import (
    "owlhnode/zeek"
    "owlhnode/changeControl"
    "github.com/astaxie/beego/logs")


func GetZeek(username string) (status zeek.Zeek, err error) {
    logs.Info("Zeek Status!!")
    status,err = zeek.GetZeek()
    // changecontrol.ChangeControlInsertData(err, "GetZeek")    
    return status,err
}

func SetZeek(zeekdata zeek.Zeek, username string) (status zeek.Zeek, err error) {
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
    cc["username"] = username
    cc["actionDescription"] = "Set Zeek values and manage service status"
    
    changecontrol.InsertChangeControl(cc)
    
    // changecontrol.ChangeControlInsertData(err, "SetZeek")    
    return status,err

}

func RunZeek(username string) (data string, err error) {
    logs.Info("Run RunZeek system into node server")
    data,err = zeek.StartZeek("")
    // changecontrol.ChangeControlInsertData(err, "RunZeek")    
    return data,err
}

func StartZeek(action string, username string) (data string, err error) {
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
    cc["username"] = username
    changecontrol.InsertChangeControl(cc)
    
    // changecontrol.ChangeControlInsertData(err, "StartZeek")    
    return data,err
}

func StopZeek(username string) (data string, err error) {
    logs.Info("Stops StopZeek system into node server")
    data,err = zeek.StopZeek()
    // changecontrol.ChangeControlInsertData(err, "StopZeek")    
    return data,err
}

func DeployZeek(username string) ( err error) {
    logs.Info("DeployZeek system into node server")
    err = zeek.DeployZeek()
    // changecontrol.ChangeControlInsertData(err, "DeployZeek")    
    return err
}

func ChangeZeekMode(anode map[string]string, username string) (err error) {
    cc := anode
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

func AddClusterValue(anode map[string]string, username string) (err error) {
    cc := anode
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

func PingCluster(username string)(data map[string]map[string]string, err error) {
    logs.Info("Ping Zeek cluster values")
    data, err = zeek.PingCluster()
    // changecontrol.ChangeControlInsertData(err, "PingCluster")    
    return data, err
}

func EditClusterValue(anode map[string]string, username string) (err error) {
    cc := anode
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

func DeleteClusterValue(anode map[string]string, username string) (err error) {
    cc := anode
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

func SyncCluster(anode map[string]string, username string) (err error) {
    cc := anode
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

func SavePolicyFiles(files map[string]map[string][]byte, username string) (err error) {
    cc := files
    logs.Info("Zeek - SavePolicyFiles")
    for key :=range cc {
        logs.Info(key +" -> ")
    }

    err = zeek.SavePolicyFiles(files)
    // changecontrol.ChangeControlInsertData(err, "SavePolicyFiles")    
    return err
}

func SyncClusterFile(anode map[string][]byte, username string) (err error) {
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

func LaunchZeekMainConf(anode map[string]string, username string) (err error) {
    logs.Info("ZEEK - LaunchZeekMainConf")
    // for key :=range cc {
    //     logs.Info(key +" -> " + cc[key])
    // }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

    if anode["param"] == "start"{
        err = zeek.StartingZeek()
    }else if anode["param"] == "stop"{
        err = zeek.StoppingZeek()
        err = zeek.ChangeZeekPreviousStatus()
        // err = zeek.RemoveZeekData()
    }else if anode["param"] == "deploy"{
        err = zeek.DeployZeek()
    }
    // if anode["saveZeek"] == "true" && anode["param"] != "stop" {
    //     err = zeek.SaveZeekData()
    // }
    // err = zeek.LaunchZeekMainConf(anode)
    // changecontrol.ChangeControlInsertData(err, "SyncClusterFile")    
    return err
}

// func SaveZeekValues(anode map[string]string) (err error) {
//     cc := anode
//     logs.Info("============")
//     logs.Info("ZEEK - SaveZeekValues")
//     for key :=range cc {
//         logs.Info(key +" -> " + cc[key])
//     }
//     delete(anode,"action")
//     delete(anode,"controller")
//     delete(anode,"router")


//     err = zeek.SaveZeekValues(anode)
//     // changecontrol.ChangeControlInsertData(err, "SaveZeekValues")    
//     return err
// }

func SyncZeekValues(anode map[string]string, username string) (err error) {
    logs.Info("ZEEK - SyncZeekValues")
    // for key :=range cc {
    //     logs.Info(key +" -> " + cc[key])
    // }
    delete(anode,"action")
    delete(anode,"controller")
    delete(anode,"router")

    err = zeek.SyncZeekValues(anode)

    // err = zeek.SyncZeekValues(anode)
    // changecontrol.ChangeControlInsertData(err, "SyncZeekValues")    
    return err
}