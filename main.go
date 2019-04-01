package main

import (

    // "github.com/astaxie/beego/logs"
    _ "owlhnode/routers"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/plugins/cors"
    "owlhnode/database"
    "owlhnode/stap"
    // "sync"
    // "owlhnode/suricata"
    // "owlhnode/wazuh"
    //"owlhnode/zeek"
)

func main() {

    ndb.SConn()
    ndb.Conn()

	// //put all the stap servers to true.
	// logs.Info("First launch, putting all STAP servers to true")
    // updateStap, err := ndb.Sdb.Prepare("update servers set server_value = ? where server_param = ?;")
    // _, err = updateStap.Exec("true", "status")  
    // defer updateStap.Close()
    // if (err != nil){
	// 	logs.Error("Error putting all the STAP servers to true for 1st time: "+err.Error())
	// }
	
    //Launch StapInit for 1st time for check status and go concurrency if status==true
    stap.StapInit()
    
    if beego.BConfig.RunMode == "dev" {
        beego.BConfig.WebConfig.DirectoryIndex = true
        beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
    }
    
    beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
        AllowOrigins:     []string{"*"},
        AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowHeaders:     []string{"Origin", "Authorization", "Access-Control-Allow-Origin"},
        ExposeHeaders:    []string{"Content-Length", "Access-Control-Allow-Origin"},
        AllowCredentials: true,
    }))

    beego.Run()
}
