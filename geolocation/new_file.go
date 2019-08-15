package geolocation

import (
    "github.com/astaxie/beego/logs"
    "github.com/oschwald/geoip2-golang"
)

func GeoInfo() {
    db, err := Open("conf/GeoLite2-City.mmdb")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()
    // If you are using strings that may be invalid, check that ip is not nil
    ip := net.ParseIP("81.2.69.142")
    record, err := db.City(ip)
    if err != nil {
        logs.Error(err)
    }
    logs.Info("Portuguese (BR) city name: %v\n", record.City.Names["pt-BR"])
    logs.Info("English subdivision name: %v\n", record.Subdivisions[0].Names["en"])
    logs.Info("Russian country name: %v\n", record.Country.Names["ru"])
    logs.Info("ISO country code: %v\n", record.Country.IsoCode)
    logs.Info("Time zone: %v\n", record.Location.TimeZone)
    logs.Info("Coordinates: %v, %v\n", record.Location.Latitude, record.Location.Longitude)
}