package utils

import (
    "fmt"
    "github.com/astaxie/beego/logs"
    "net"
    "strings"
)

func LocalAddresses() (laddr []string) {
    ifaces, err := net.Interfaces()
    if err != nil {
        logs.Error(fmt.Errorf("localAddresses: %+v\n", err.Error()))
        return
    }
    for _, i := range ifaces {
        addrs, err := i.Addrs()
        if err != nil {
            logs.Error(fmt.Errorf("localAddresses: %+v\n", err.Error()))
            continue
        }
        for _, a := range addrs {
            switch v := a.(type) {
            case *net.IPNet:
                ip := v.IP
                laddr = append(laddr, ip.String())
            case *net.IPAddr:
                logs.Debug("%v : %s (%s)\n", i.Name, v, v.IP.DefaultMask())
            }
        }
    }
    return laddr
}

func IsLocalAddress(address string) bool {
    if strings.ToLower(address) == "localhost" {
        return true
    }
    localaddrs := LocalAddresses()
    for addr := range localaddrs {
        if address == localaddrs[addr] {
            return true
        }
    }
    return false
}
