package stap

import (
    "github.com/astaxie/beego/logs"
    "owlhnode/utils"
    "owlhnode/database"
    "errors"
    "encoding/json"
    "sync"
    "strconv"
)

var waitGroup sync.WaitGroup
var status bool
// var data chan string


func AddServer(elem map[string]string) (err error){
    nodeName:= elem["nodeName"]
    nodeIP:= elem["nodeIP"]
    uuidServer := utils.Generate()
    
    //database connection
    if ndb.Sdb == nil {
        logs.Error("AddServer stap -- Can't access to database")
        return errors.New("AddServer stap -- Can't access to database")
    } 
    
    //insert ip into server database
    logs.Info("stap/AddServer INSERT new server ip")
    err = ndb.AddServer(uuidServer, "ip", nodeIP)
    if err != nil {logs.Error("ERROR stap/AddServer INSERT new server ip: "+err.Error()); return err} 

    //insert name into server database
    logs.Info("stap/AddServer INSERT new server name")
    err = ndb.AddServer(uuidServer, "name", nodeName)
    if err != nil {logs.Error("ERROR stap/AddServer INSERT new server name: "+err.Error()); return err} 
    
    //Load default server data from main.conf
    jsonDefaultData,err := utils.LoadDefaultServerData("software TAP PULL mode conf")
    logs.Info("File readed !!!")
    jsonData := make(map[string]string)

    //parse raw data into byte array
    jsonDataArray := []byte(jsonDefaultData["fileContent"])
    err = json.Unmarshal(jsonDataArray, &jsonData)
    if err != nil {logs.Error("ERROR unmarshal json on AddServer func: "+err.Error()); return err}        

    //insert default data into server database
    for z,v := range jsonData{
        err = ndb.AddServer(uuidServer, z, v)
        if err != nil {logs.Error("ERROR stap/AddServer INSERT default server: "+err.Error()); return err} 
    }

    return nil
}

func GetAllServers()(data map[string]map[string]string, err error){
    allservers,err := ndb.GetAllServers()
    if err != nil {logs.Error("ERROR stap/GetAllServers: "+err.Error()); return nil, err} 
    return allservers, nil
} 

func GetServer(serveruuid string)(data map[string]map[string]string, err error){
    server, err := ndb.GetServer(serveruuid)
    if err != nil {logs.Error("ERROR stap/GetServer: "+err.Error()); return nil, err} 
    return server, nil
} 



func PingStap(uuid string) (isIt map[string]bool, err error){
    //database connection
    if ndb.Sdb == nil {
        logs.Error("PingStap stap -- Can't access to database")
        return nil, errors.New("PingStap stap -- Can't access to database")
    } 
    var res string
    stap := make(map[string]bool)
    stap["stapStatus"] = false
    
    sql := "select stap_value from stap where stap_param = \"status\";"
    rows, err := ndb.Sdb.Query(sql)
    if err != nil {
        logs.Info("PingStap Query Error immediately after retrieve data %s",err.Error())
        return stap,err
    }
    defer rows.Close()
    if rows.Next() {
        err := rows.Scan(&res)
        if err != nil {
            logs.Info("Stap query error %s",err.Error())
            return stap, err
        }
        if res=="true"{
            stap["stapStatus"]=true
        }else{
            stap["stapStatus"]=false
        }
    }else if uuid != "" {
        logs.Info("Put Stap status INSERT")
        insertStap, err := ndb.Sdb.Prepare("insert into stap (stap_uniqueid, stap_param, stap_value) values (?,?,?);")
        _, err = insertStap.Exec(&uuid, "status", "false")  
        defer insertStap.Close()
        if err != nil{
            logs.Info("Error Insert uuid !=")
            return stap,err
        }
    }
    logs.Info("checking stap value "+strconv.FormatBool(stap["stapStatus"]))
    return stap,nil
}


//Run stap main server on node
func RunStap(uuid string)(data string, err error){
    logs.Info("Starting Stap...  "+uuid)
    err = ndb.UpdateStapByparam("status", "true")
    if err != nil {logs.Error("ERROR stap/RunStap: "+err.Error()); return "", err} 

    //call Concurrency StapInit
    StapInit()

    return "Stap is Running!", err
}

//Stop Stap
func StopStap(uuid string)(data string, err error){
    logs.Info("Stopping Stap..."+uuid)
    err = ndb.UpdateStapByparam("status", "false")
    if err != nil {logs.Error("ERROR stap/StopStap: "+err.Error()); return "", err} 

    status = false
    return "Stap is stopped", err
}

//Run stap specific server
func RunStapServer(serveruuid string)(data string, err error){
    logs.Info("Running Stap server..."+serveruuid)
    err = ndb.UpdateStap(serveruuid, "status", "true")
    if err != nil {logs.Error("ERROR stap/RunStapServer: "+err.Error()); return "", err} 
    
    return "Stap specific server is Running!", err
}

//Stop Stap specific server
func StopStapServer(serveruuid string)(data string, err error){
    logs.Info("Stopping Stap server..."+serveruuid)
    err = ndb.UpdateStap(serveruuid, "status", "false")
    if err != nil {logs.Error("ERROR stap/StopStapServer: "+err.Error()); return "", err} 

    return "Stap specific server  is stoped", err
}
//Error Stap specific server
func ErrorStapServer(serveruuid string)(data string, err error){
    err = ndb.UpdateStap(serveruuid, "status", "error")
    if err != nil {logs.Error("ERROR stap/ErrorStapServer: "+err.Error()); return "", err} 

    return "Stap specific server status is error now", err
}

func PingServerStap(server string) (isIt map[string]string, err error){
    isIt,err = ndb.PingServerStap(server)
    if err != nil {logs.Error("ERROR stap/PingServerStap: "+err.Error()); return nil, err} 

    return isIt, nil
}

func GetStapUUID()(uuid string){
    uuid = ndb.GetStapUUID()
    return uuid
}

//delete specific stap server
func DeleteStapServer(serveruuid string)(data string, err error){
    err = ndb.DeleteServer(serveruuid)
    if err != nil {logs.Error("ERROR stap/DeleteStapServer: "+err.Error()); return "", err} 

    return "Stap specific server is deleted!", err
}

//Edit params for stap servers
func EditStapServer(data map[string]string) (err error){
    err = ndb.UpdateStap(data["server"], data["param"], data["value"])
    if err != nil {logs.Error("ERROR stap/EditStapServer: "+err.Error()); return err} 

    return nil
}