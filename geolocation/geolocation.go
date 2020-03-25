package geolocation

import (
    "github.com/astaxie/beego/logs"
    "github.com/oschwald/geoip2-golang"
    "net"
)

var GeoDb *geoip2.Reader
var running bool

func Init() {
    var err error
    running = false
    GeoDb, err = geoip2.Open("conf/GeoLite2-City.mmdb")
    if err != nil {
        logs.Error("unable to load GEO database " + err.Error())
        return
    }
    running = true
    //defer GeoDb.Close()
}

func GetGeoInfo(nip string)(geoinfo map[string]string) {
    if ! running {
        return nil
    }
    ip := net.ParseIP(nip)
    record, err := GeoDb.City(ip)
    if err != nil {
        logs.Error(err)
        return nil 
    }
    geoinfo = map[string]string{}
    if record.City.Names != nil {
        geoinfo["city_name"] = record.City.Names["en"]
    } 
    if record.Subdivisions != nil {
        geoinfo["state_province"] = record.Subdivisions[0].Names["en"]
    }
    if record.Country.Names != nil {
        geoinfo["country_name"] = record.Country.Names["en"]
    }
    if record.Country.IsoCode != "" {
        geoinfo["country_code"] = record.Country.IsoCode        
    }
    if record.Continent.Code != "" {
        geoinfo["continent_code"] = record.Continent.Code        
    }

    return geoinfo
}