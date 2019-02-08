package main

import (

    "github.com/astaxie/beego/logs"
    _ "owlhnode/routers"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/plugins/cors"
    "owlhnode/database"
    "owlhnode/suricata"
    "owlhnode/wazuh"
    "owlhnode/zeek"
)

func main() {

    ndb.Conn()
    err := ndb.Get_master()

    suricata.Installed()
    wazuh.Installed()
    zeek.Installed()

    if err != nil {
        logs.Info("Main -> no puedo leer el master")
    }

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
