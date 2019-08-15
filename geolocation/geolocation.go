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
    // If you are using strings that may be invalid, check that ip is not nil
    ip := net.ParseIP("137.101.151.36")
    record, err := db.City(ip)
    if err != nil {
        logs.Error(err)
    }
    logs.Info("Portuguese (BR) city name: %v\n", record.City.Names["en"])
    logs.Info("English subdivision name: %v\n", record.Subdivisions[0].Names["en"])
    logs.Info("Russian country name: %v\n", record.Country.Names["en"])
    logs.Info("ISO country code: %v\n", record.Country.IsoCode)
    logs.Info("Time zone: %v\n", record.Location.TimeZone)
    logs.Info("Coordinates: %v, %v\n", record.Location.Latitude, record.Location.Longitude)
}



func GetGeoInfo(ip string)(geoinfo map[string]string{}) {
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