package geolocation

import (
    "github.com/astaxie/beego/logs"
    "github.com/oschwald/geoip2-golang"
    "net"
)

var GeoDb *geoip2.Reader

func Init() {
    var err error
    GeoDb, err = geoip2.Open("conf/GeoLite2-City.mmdb")
    if err != nil {
        logs.Error("unable to load GEO database " + err.Error())
    }
    //defer GeoDb.Close()
}

func GetGeoInfo(nip string)(geoinfo map[string]string) {
    logs.Info("geolocation -> ip : "+nip)
    ip := net.ParseIP(nip)
    logs.Info("geolocation -> ip")
    logs.Notice(nip)
    record, err := GeoDb.City(ip)
    if err != nil {
        logs.Error(err)
    }
    logs.Notice(record)
    geoinfo = map[string]string{}
    logs.Info("geolocation -> city name")
    if record.City.Names != nil {
        logs.Info("geolocation -> city : "+record.City.Names["en"])
        geoinfo["city_name"] = record.City.Names["en"]
    } 
    logs.Info("geolocation -> province")
    if record.Subdivisions != nil {
        logs.Info("geolocation -> province : "+record.Subdivisions[0].Names["en"])
        geoinfo["state_province"] = record.Subdivisions[0].Names["en"]
    }
    logs.Info("geolocation -> country name")
    if record.Country.Names != nil {
        logs.Info("geolocation -> country_name : "+record.Country.Names["en"])
        geoinfo["country_name"] = record.Country.Names["en"]
    }
    logs.Info("geolocation -> country code")
    if record.Country.IsoCode != "" {
        logs.Info("geolocation -> country_code : "+record.Country.IsoCode)
        geoinfo["country_code"] = record.Country.IsoCode        
    }
    logs.Info("geolocation -> continent code")
    if record.Continent.Code != "" {
        logs.Info("geolocation -> continent : "+record.Continent.Code)
        geoinfo["continent_code"] = record.Continent.Code        
    }

    return geoinfo
}