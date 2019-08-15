package geolocation

import (
    "github.com/astaxie/beego/logs"
    "github.com/oschwald/geoip2-golang"
    "net"
)

var GeoDb geoip2.Reader

func Init() {
    GeoDb, err := geoip2.Open("conf/GeoLite2-City.mmdb")
    if err != nil {
        logs.Error(err)
    }
    defer GeoDb.Close()
}

func GetGeoInfo(nip string)(geoinfo map[string]string) {
    ip := net.ParseIP(nip)
    record, err := GeoDb.City(ip)
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