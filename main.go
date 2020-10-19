package main

import (
    "github.com/astaxie/beego/logs"
    // "github.com/astaxie/beego/context"
    "bufio"
    "crypto/tls"
    "github.com/astaxie/beego"
    "github.com/astaxie/beego/plugins/cors"
    "os"
    "os/signal"
    "owlhnode/analyzer"
    "owlhnode/configuration"
    "owlhnode/database"
    "owlhnode/geolocation"
    "owlhnode/knownports"
    "owlhnode/monitor"
    "owlhnode/pcap"
    "owlhnode/plugin"
    _ "owlhnode/routers"
    "owlhnode/stap"
    "owlhnode/utils"
    "owlhnode/zeek"
    "runtime"
    "strings"
    "syscall"
)

var version string

func main() {

    utils.Load()
    
    version = "0.17.0.20200908"
    logs.Info("OwlH Node : v%s", version)

    cancontinue := configuration.MainCheck()
    if !cancontinue {
        logs.Error("can't continue, see previous logs")
        // return
    }

    //operative system values
    data := OperativeSystemValues()
    for x := range data {
        if x == "ID" || x == "ID_LIKE" || x == "VERSION_ID" {
            logs.Info(x + " -- " + data[x])
        }
        // if (x == "ID" && data[x] == "debian") {
        //     logs.Info("debian")
        //     socatOutput, err := exec.Command("bash","-c","dpkg -l socat | grep socat").Output()
        //     if err != nil {logs.Error("Error checking socat for debian: "+err.Error())}
        //     logs.Info(socatOutput)

        //     tcpdumpOutput, err := exec.Command("bash","-c","dpkg -l tcpdump  | grep tcpdump").Output()
        //     if err != nil {logs.Error("Error checking tcpdump for debian: "+err.Error())}
        //     logs.Info(tcpdumpOutput)
        //     //check socat
        //         //install socat
        //     //check tcpdump
        //         //install tcpdump

        // }else if (data[x] != "debian") {
        //     logs.Info("not debian")
        //     // socatOutput, err := exec.Command("bash","-c","yum list socat | grep socat").Output()
        //     socatOutput, err := exec.Command("yum","list", "socat", "|", "grep", "socat").Output()
        //     if err != nil {logs.Error("Error checking socat for centos/redhat: "+err.Error())}
        //     logs.Info(socatOutput)

        //     // tcpdumpOutput, err := exec.Command("bash","-c","yum list tcpdump | grep tcpdump").Output()
        //     tcpdumpOutput, err := exec.Command("yum","list", "tcpdump", "|", "grep", "tcpdump").Output()
        //     if err != nil {logs.Error("Error checking tcpdump for centos/redhat: "+err.Error())}
        //     logs.Info(tcpdumpOutput)

        //     //check socat
        //         //install socat
        //     //check tcpdump
        //         //install tcpdump
        // }
    }

    logs.Info("Main Starting -> reading STAP DB")
    ndb.SConn() //stap database
    logs.Info("Main Starting -> reading PLUGINS DB")
    ndb.PConn() //plugins database
    logs.Info("Main Starting -> reading NODE DB")
    ndb.NConn() //node database
    logs.Info("Main Starting -> reading MONITOR DB")
    ndb.MConn() //monitor database
    logs.Info("Main Starting -> reading GROUP DB")
    ndb.GConn() //group database

    //launch logger
    monitor.Logger()
    go ManageSignals()
    go monitor.FileRotation()
    zeek.Init()
    plugin.CheckServicesStatus()
    stap.StapInit()
    knownports.Init()
    analyzer.Init()
    geolocation.Init()
    monitor.Init()
    pcap.Init()

    if beego.BConfig.RunMode == "dev" {
        beego.BConfig.WebConfig.DirectoryIndex = true
        beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
    }

    beego.BeeApp.Server.TLSConfig = &tls.Config{CipherSuites: []uint16{
        tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
        tls.TLS_RSA_WITH_AES_256_GCM_SHA384,
    },
        MinVersion:               tls.VersionTLS12,
        PreferServerCipherSuites: true,
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

func OperativeSystemValues() (values map[string]string) {
    if runtime.GOOS == "linux" {
        logs.Info("============" + runtime.GOOS + "============")
        var OSmap = make(map[string]string)
        file, err := os.Open("/etc/os-release")
        if err != nil {
            logs.Error("No os-release file")
        }
        defer file.Close()

        scanner := bufio.NewScanner(file)
        for scanner.Scan() {
            if scanner.Text() != "" {
                sidsSplit := strings.Split(scanner.Text(), "=")
                str := strings.Replace(sidsSplit[1], "\"", "", -1)
                OSmap[sidsSplit[0]] = str
            }
        }
        return OSmap
    } else {
        return nil
    }
}

func ManageSignals() {
    sigs := make(chan os.Signal, 1)
    signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR1, syscall.SIGUSR2)

    go func() {
        sig := <-sigs
        logs.Info("Signal received: " + sig.String())

        //kill plugins
        plugin.StopPluginsGracefully()

        //stop node
        logs.Critical("Stopping Node...")
        os.Exit(0)
    }()
}
