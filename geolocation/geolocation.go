package geolocation

import (
    "github.com/astaxie/beego/logs"
    "github.com/oschwald/geoip2-golang"
    "net"
)

func Init() {
    db, err := geoip2.Open("conf/GeoLite2-City.mmdb")
    if err != nil {
        logs.Error(err)
    }
    defer db.Close()
}

func GetGeoInfo(ip string)(geoinfo map[string]string) {
    ip := net.ParseIP(ip)
    record, err := db.City(ip)
    if err != nil {
        logs.Error(err)
    }
    //geoinfo = map[string]string{}
    geoinfo["city_name"] = record.City.Names["en"]
    geoinfo["county"] = record.Subdivisions[0].Names["en"]
    geoinfo["country"] = record.Country.Names["en"]
    geoinfo["country_code"] = record.Country.IsoCode
    geoinfo["continent_code"] = record.Continent.Code
    return geoinfo
}