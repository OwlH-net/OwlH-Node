package pcap
 
import (
    // "bytes"
    // "encoding/binary"
    "encoding/json"
    "errors"
    // "log"
    "net"
    "io/ioutil"
    // "sync"
    "github.com/astaxie/beego/logs"
    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "github.com/google/gopacket/pcap"
    "time"
)
 
type MAC struct {
    Mac     net.HardwareAddr `json:"mac"`
    IPs     IPs              `json:"ips"`
    White   bool             `json:"white"`
    Onnewip bool             `json:"onnewip"`
    Alerted bool             `json:"alerted"`
    First   time.Time        `json:"first"`
    Last    time.Time        `json:"last"`
}
 
type IP struct {
    Ip      string    `json:"ip"`
    White   bool      `json:"white"`
    First   time.Time `json:"first"`
    Last    time.Time `json:"last"`
    Alerted bool      `json:"alerted"`
}
 
type Alert struct {
    Sid       string `json:"sid"`
    Type      string `json:"type"`
    Signature string `json:"signature"`
    Timestamp string `json:"timestamp"`
}
 
type IPs map[string]IP
 
type ARPConfig struct {
    Onnewmac    bool
    Onnewip     bool
    Enabled     bool
    Learning    bool
    Verbose     bool
    Interface   string
    KnownFile   string
    CurrentFile string
}
 
const (
    aMAC = iota
    anIP
)
 
var arpmain ARPConfig
var Knownmacs = map[string]MAC{}
var Currentmacs = map[string]MAC{}
 
func saveKnownMacs() {
 
    return
}
 
func isknownMac(arp *layers.ARP) (is bool) {
    var dsthw net.HardwareAddr
    dsthw = arp.DstHwAddress
 
    if _, ok := Knownmacs[dsthw.String()]; ok {
        logs.Info("mac %s exists", dsthw.String())
        return true
    } else {
        logs.Info("mac %s does not exist.", dsthw.String())
        return false
    }
    return false
}
 
func isKnownIPByMac(mac MAC, ip IP) (itis bool) {
 
    if _, ok := Knownmacs[mac.Mac.String()]; ok {
        if _, exists := Knownmacs[mac.Mac.String()].IPs[ip.Ip]; exists {
            return true
        }
    }
    return false
}
 
func isKnownIPAnyMac(ip IP) (itis bool) {
    for mac := range Knownmacs {
        macips := Knownmacs[mac].IPs
        if _, ok := macips[ip.Ip]; ok {
            logs.Info("ip %s does exist on mac %v", ip.Ip, Knownmacs[mac].Mac)
            return true
        }
    }
    logs.Info("ip %s does NOT exist on any mac", ip.Ip)
    return false
}
 
func addKnownIP(IPs IPs, newIP IP) (modips IPs) {
    if arpmain.Verbose {
        logs.Info("add ip %s", newIP.Ip)
    }
 
    // var modips = IPs
    // var iplist = map[string]IP{}
 
    // if ! modips.IPS[IPs]
 
    return modips
}
 
func addKnownMac(arp *layers.ARP) {
    var dsthw net.HardwareAddr
    dsthw = arp.DstHwAddress
 
    logs.Info("adding mac %s ", dsthw.String())
 
    // Manage MAC
    var cMac MAC
    cMac.Mac = arp.DstHwAddress
    cMac.First = time.Now()
    cMac.Last = time.Now()
    if arpmain.Verbose {
        logs.Info("added at %v", cMac.First)
    }
 
    // Manage new IP
    var newip IP
    newip.Ip = net.IP(arp.DstProtAddress).String()
    if arpmain.Verbose {
        logs.Info("lets add %s", newip.Ip)
    }
    newip.First = cMac.First
    newip.Last = cMac.First
 
    allips := make(IPs)
    allips[newip.Ip] = newip
    cMac.IPs = allips
 
    //    cMac.IPs = addKnownIP(cMac.IPs, newip)
    Knownmacs[dsthw.String()] = cMac
}
 
func updateLast(arp *layers.ARP, known bool) (err error) {
    var dsthw net.HardwareAddr
    dsthw = arp.DstHwAddress
    if known {
        cmac := Knownmacs[dsthw.String()]
        cmac.Last = time.Now()
        Knownmacs[dsthw.String()] = cmac
    } else {
        cmac := Currentmacs[dsthw.String()]
        cmac.Last = time.Now()
        Currentmacs[dsthw.String()] = cmac
    }
    return nil
}
 
func learnarp(arp *layers.ARP) (err error) {
    if !isknownMac(arp) {
        addKnownMac(arp)
    } else {
        updateLast(arp, true)
    }
    saveKnownMacs()
 
    return nil
 
}
 
func alertNewARP(arp *layers.ARP, alertabout int) {
    logs.Notice("alerting new arp")
    switch alertabout {
    case anIP:
        logs.Info("is IP - injecting ip alert")
    case aMAC:
        logs.Info("is MAC - injecting mac alert")
    default:
        logs.Error("have no idea what we try to alert about %v", alertabout)
    }
    return
}
 
func isknowMACIP(arp *layers.ARP) (isknown bool) {
    logs.Info("isknowMACIP")
    return true
}
 
func checkarp(arp *layers.ARP) {
 
    if !isknownMac(arp) {
        if arpmain.Onnewmac {
            alertNewARP(arp, aMAC)
        }
    } else if !isknowMACIP(arp) {
        if arpmain.Onnewip {
            alertNewARP(arp, anIP)
        }
    }
    return
}
 
func readARP(iface string) (err error) {
    logs.Info("starting read traffic from %s", arpmain.Interface)
    iface = arpmain.Interface
 
    if handle1, err := pcap.OpenLive(iface, 65536, true, pcap.BlockForever); err != nil {
        logs.Error(err)
        return err
    } else {
        handle1.SetBPFFilter("arp")
        src := gopacket.NewPacketSource(handle1, layers.LayerTypeEthernet)
        in := src.Packets()
 
        for {
            packet := <-in
            // logs.Info ("new packet")
            arpLayer := packet.Layer(layers.LayerTypeARP)
            if arpLayer == nil {
                continue
            }
            arp := arpLayer.(*layers.ARP)
            var dsthw net.HardwareAddr
 
            // logs.Info("arp.DstProtAddress    %v",arp.DstProtAddress   )
            if arp.Operation == 1 {
            } else if arp.Operation == 2 {
                if arpmain.Learning {
                    logs.Info("learning")
                    learnarp(arp)
                } else {
                    logs.Info("live")
                    checkarp(arp)
                }
 
                dsthw = arp.DstHwAddress
                logs.Info("%v is at %v", net.IPv4(arp.DstProtAddress[0], arp.DstProtAddress[1], arp.DstProtAddress[2], arp.DstProtAddress[3]), dsthw)
            } else {
                logs.Error("unkonwn operation %d", arp.Operation)
            }
        }
    }
    return nil
}
 
func Init() {
    ifaces, err := net.Interfaces()
    if err != nil {
        logs.Error("MAC Management - error getting interfaces -> %s" + err.Error())
        return
    }
    logs.Info(ifaces)
 
    arpmain.Enabled = true
    arpmain.Interface = "eth0"
    arpmain.Learning = true
    arpmain.Onnewip = true
    arpmain.Onnewmac = true
    arpmain.Verbose = true
    arpmain.KnownFile = "known.db"
    arpmain.CurrentFile = "current.db"
 
    go readARP("")
 
    // readARP("")
}
 
const (
    Known = iota
    Current
)

func ReadMacFileContent(alertabout int)(err error){    
    switch alertabout {
        case Known:
            byteValue, err := ioutil.ReadFile(arpmain.KnownFile)
            if err != nil {logs.Error("ReadFileContent - error getting KnownFile content -> %s" + err.Error()); return err}
            json.Unmarshal(byteValue, &Knownmacs)
        case Current:
            byteValue, err := ioutil.ReadFile(arpmain.CurrentFile)
            if err != nil {logs.Error("ReadFileContent - error getting CurrentFile content -> %s" + err.Error()); return err}
            json.Unmarshal(byteValue, &Currentmacs)
        default:
            return errors.New("ReadMacFileContent Invalid unmarshal variable")
    }

    return nil
}

func WriteMacFileContent(alertabout int)(err error){    
    switch alertabout {
        case Known:
            values, _ := json.Marshal(Knownmacs)
            err = ioutil.WriteFile(arpmain.KnownFile, values, 0644)
            if err!=nil{ logs.Error("ReadFileContent - error writing file content -> %s" + err.Error()); return err}
        case Current:
            values, _ := json.Marshal(Currentmacs)
            err = ioutil.WriteFile(arpmain.CurrentFile, values, 0644)
            if err!=nil{ logs.Error("ReadFileContent - error writing file content -> %s" + err.Error()); return err}
        default:
            return errors.New("ReadMacFileContent Invalid unmarshal variable")
    }

    return nil
}


// var wg sync.WaitGroup
// for _, iface := range ifaces {
//     // wg.Add(1)
//     // logs.Info("iface %s", iface.Name)
//     // if iface.Name == "eth0" {
//     //     go func(iface net.Interface) {
//     //         defer wg.Done()
//     //         if err := scan(&iface); err != nil {
//     //             logs.Error("interface %v: %v", iface.Name, err)
//     //         }
//     //     }(iface)
//     // }
// }
// wg.Wait()
 
// content := string(arp.BaseLayer.Payload)
// logs.Info("arp.BaseLayer         %v",content         )
// logs.Info("arp.AddrType          %v",arp.AddrType         )
// logs.Info("arp.Protocol          %v",arp.Protocol         )
// logs.Info("arp.HwAddressSize     %v",arp.HwAddressSize    )
// logs.Info("arp.ProtAddressSize   %v",arp.ProtAddressSize  )
// logs.Info("arp.Operation         %v",arp.Operation        )
// logs.Info("arp.SourceHwAddress   %v",arp.SourceHwAddress  )
// logs.Info("arp.SourceProtAddress %v",arp.SourceProtAddress)
// logs.Info("arp.DstHwAddress      %v",arp.DstHwAddress     )
// var srchw net.HardwareAddr
 
// func scan(iface string) error {
//     // We just look for IPv4 addresses, so try to find if the interface has one.
//     // var addr *net.IPNet
//     // if addrs, err := iface.Addrs(); err != nil {
//     //     return err
//     // } else {
//     //     for _, a := range addrs {
//     //         if ipnet, ok := a.(*net.IPNet); ok {
//     //             if ip4 := ipnet.IP.To4(); ip4 != nil {
//     //                 addr = &net.IPNet{
//     //                     IP:   ip4,
//     //                     Mask: ipnet.Mask[len(ipnet.Mask)-4:],
//     //                 }
//     //                 break
//     //             }
//     //         }
//     //     }
//     // }
//     // // Sanity-check that the interface has a good address.
//     // if addr == nil {
//     //     return errors.New("no good IP network found")
//     // } else if addr.IP[0] == 127 {
//     //     return errors.New("skipping localhost")
//     // } else if addr.Mask[0] != 0xff || addr.Mask[1] != 0xff {
//     //     return errors.New("mask means network is too large")
//     // }
//     // log.Printf("Using network range %v for interface %v", addr, iface.Name)
 
//     // // Open up a pcap handle for packet reads/writes.
//     // handle, err := pcap.OpenLive(iface.Name, 65536, true, pcap.BlockForever)
//     // if err != nil {
//     //     return err
//     // }
//     var handle *pcap.Handle
 
//     if handle, err := pcap.OpenLive(iface, 65536, true, pcap.BlockForever); err != nil {
//       logs.Error(err)
//     } else if err := handle.SetBPFFilter("arp"); err != nil {  // optional
//       logs.Error(err)
//     }
//     for {
//         logs.Info("let's scan ARP for iface: %s", iface)
//         readARP(handle)
//     }
//     return nil
// }
 
// if arp.Operation != layers.ARPReply || bytes.Equal([]byte(iface.HardwareAddr), arp.SourceHwAddress) {
//     // This is a packet I sent.
//     continue
// }
// Note:  we might get some packets here that aren't responses to ones we've sent,
// if for example someone else sends US an ARP request.  Doesn't much matter, though...
// all information is good information :)
// log.Printf("IP %v is at %v", net.IP(arp.SourceProtAddress), net.HardwareAddr(arp.SourceHwAddress))
 
// dsthw := string(arp.DstHwAddress[0])+":"+string(arp.DstHwAddress[1])+":"+string(arp.DstHwAddress[2])+":"+string(arp.DstHwAddress[3])+":"+string(arp.DstHwAddress[4])+":"+string(arp.DstHwAddress[5])
 
// srchw = arp.SourceHwAddress
// logs.Info("who has %v tells %v(%v)", net.IPv4(arp.DstProtAddress[0],arp.DstProtAddress[1],arp.DstProtAddress[2],arp.DstProtAddress[3]), net.IPv4(arp.SourceProtAddress[0],arp.SourceProtAddress[1],arp.SourceProtAddress[2],arp.SourceProtAddress[3]), srchw.String())