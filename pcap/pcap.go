package pcap
 
import (
    // "bytes"
    // "encoding/binary"
    "encoding/json"
    "errors"
    // "log"
    "net"
    "owlhnode/utils"
    "owlhnode/analyzer"
    "io/ioutil"
    // "sync"
    "github.com/astaxie/beego/logs"
    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "github.com/google/gopacket/pcap"
    "time"
)



type EventAlert struct {
    Timestamp       time.Time               `json:"timestamp"`
    In_iface        string                  `json:"in_iface"`
    Event_type      string                  `json:"event_type"`
    Src_ip          string                  `json:"src_ip"`
    Mac_address     string                  `json:"mac_address"`
    Proto           string                  `json:"proto"`
    Alert           Alert                   `json:"alert"`
}

type Alert struct {
    Action          string      `json:"action"`
    Signature_ID    int         `json:"signature_id"`
    Rev             int         `json:"rev"`
    Signature       string      `json:"signature"`
    Category        string      `json:"category"`
    Severity        int         `json:"severity"`
    Times           int         `json:"times"`
}


type MAC struct {
    Mac         net.HardwareAddr `json:"mac"`
    IPs         IPs              `json:"ips"`
    White       bool             `json:"white"`
    Onnewip     bool             `json:"onnewip"`
    Alerted     bool             `json:"alerted"`
    LastAlert   time.Time        `json:"lastalert"`
    Times       int              `json:"times"`
    First       time.Time        `json:"first"`
    Last        time.Time        `json:"last"`
}   
 
type IP struct {
    Ip      string    `json:"ip"`
    White   bool      `json:"white"`
    First   time.Time `json:"first"`
    Last    time.Time `json:"last"`
    Alerted bool      `json:"alerted"`
}

type IPs map[string]IP
 
type ARPConfig struct {
    Onnewmac            bool
    Onnewip             bool
    Enabled             bool
    Learning            bool
    Verbose             bool        
    Interface           string         
    Timebetweenalerts   int         
    KnownFile           string
    CurrentFile         string
}
 
const (
    aMAC = iota
    anIP
)
 
var arpmain ARPConfig
var Knownmacs = map[string]MAC{}
var Currentmacs = map[string]MAC{}
 
func saveKnownMacs() {
    WriteMacFileContent(Known)
    return
}
 
func saveCurrentMacs() {
    WriteMacFileContent(Current)
    return
}
 
func isCurrentMac(arp *layers.ARP) (is bool) {
    var srcHw net.HardwareAddr
    srcHw = arp.SourceHwAddress
 
    if _, ok := Currentmacs[srcHw.String()]; ok {
        logs.Info("mac %s exists", srcHw.String())
        return true
    } else {
        logs.Info("mac %s does not exist.", srcHw.String())
        return false
    }
    return false
}
 
func isknownMac(arp *layers.ARP) (is bool) {
    var srcHw net.HardwareAddr
    srcHw = arp.SourceHwAddress
 
    if _, ok := Knownmacs[srcHw.String()]; ok {
        logs.Info("mac %s exists", srcHw.String())
        return true
    } else {
        logs.Info("mac %s does not exist.", srcHw.String())
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
 
func addCurrentMac(arp *layers.ARP) {
    var srcHw net.HardwareAddr
    srcHw = arp.SourceHwAddress
 
    logs.Info("CM adding mac %s ", srcHw.String())
 
    // Manage MAC
    var cMac MAC
    cMac.Mac = arp.SourceHwAddress
    cMac.First = time.Now()
    cMac.Last = time.Now()
    if arpmain.Verbose {
        logs.Info("CM added at %v", cMac.First)
    }
 
    // Manage new IP
    var newip IP
    newip.Ip = net.IP(arp.SourceProtAddress).String()
    if arpmain.Verbose {
        logs.Info("CM lets add %s", newip.Ip)
    }
    newip.First = cMac.First
    newip.Last = cMac.First
 
    allips := make(IPs)
    allips[newip.Ip] = newip
    cMac.IPs = allips
 
    Currentmacs[srcHw.String()] = cMac
}

func addKnownMac(arp *layers.ARP) {
    var srcHw net.HardwareAddr
    srcHw = arp.SourceHwAddress
 
    logs.Info("adding mac %s ", srcHw.String())
 
    // Manage MAC
    var cMac MAC
    cMac.Mac = arp.SourceHwAddress
    cMac.First = time.Now()
    cMac.Last = time.Now()
    if arpmain.Verbose {
        logs.Info("added at %v", cMac.First)
    }
 
    // Manage new IP
    var newip IP
    newip.Ip = net.IP(arp.SourceProtAddress).String()
    if arpmain.Verbose {
        logs.Info("lets add %s", newip.Ip)
    }
    newip.First = cMac.First
    newip.Last = cMac.First
 
    allips := make(IPs)
    allips[newip.Ip] = newip
    cMac.IPs = allips
 
    //    cMac.IPs = addKnownIP(cMac.IPs, newip)
    Knownmacs[srcHw.String()] = cMac
}
 
func updateLast(arp *layers.ARP, known bool) (err error) {
    var srcHw net.HardwareAddr
    srcHw = arp.SourceHwAddress
    var newIP IP
    newIP.Ip = net.IP(arp.SourceProtAddress).String()

    logs.Notice(newIP.Ip)

    if known {
        cmac := Knownmacs[srcHw.String()]
        cmac.Last = time.Now()

        if currentIP,ok := cmac.IPs[newIP.Ip]; ok {
            currentIP.Last = time.Now()
            cmac.IPs[newIP.Ip] = currentIP
        }
        Knownmacs[srcHw.String()] = cmac
    } else {
        cmac := Currentmacs[srcHw.String()]
        cmac.Last = time.Now()
        Currentmacs[srcHw.String()] = cmac
    }
    return nil
}
 
func learnarp(arp *layers.ARP) (err error) {
    if !isknownMac(arp) {
        addKnownMac(arp)
    } else {
        logs.Info("---Updating last---")
        updateLast(arp, true)
    }
    saveKnownMacs()
 
    return nil
 
}
 
func alertMac(arp *layers.ARP){
    var ev EventAlert
    var al Alert

    var srcHw net.HardwareAddr
    srcHw = arp.SourceHwAddress

    ev.Timestamp = time.Now()   
    ev.In_iface = arpmain.Interface
    ev.Event_type = "alert"
    ev.Src_ip = net.IP(arp.SourceProtAddress).String() 
    ev.Mac_address = srcHw.String()
    ev.Proto = "ARP"        
    
    al.Action = "allowed"          
    al.Signature_ID = 1150001    
    al.Rev = 1              
    al.Signature = "new mac detected - "+ev.Mac_address        
    al.Category = "Potentially Bad Traffic"          
    al.Severity = 2 
    al.Times = Currentmacs[srcHw.String()].Times

    ev.Alert = al    
    
    values, _ := json.Marshal(ev)
    
    analyzer.ToDispatcher("start",string(values))
    logs.Error(string(values))
}

func timeToAlert(lastTime time.Time) bool {
    logs.Info("Check time to alert")
    logs.Info(arpmain.Timebetweenalerts)

    seconds := time.Second * time.Duration(arpmain.Timebetweenalerts)
    diff := time.Now().Sub(lastTime)

    logs.Notice(seconds)
    logs.Warn(diff)

    if diff >= seconds {
        return true
    }

    return false
}

func alertIfAlert(arp *layers.ARP) {
    var srcHw net.HardwareAddr
    srcHw = arp.SourceHwAddress

    if macAlert,ok := Currentmacs[srcHw.String()]; ok  {
        logs.Info("Alert new mac in current macs!")  
        if macAlert.Alerted {
            if timeToAlert(macAlert.LastAlert) {
                logs.Info("Time to alert again -- %d",macAlert.Times)
                alertMac(arp)     
                macAlert.Times = 0
                macAlert.LastAlert = time.Now()
            }else{
                macAlert.Times += 1
            }
        }else{
            logs.Info("First time alert")
            alertMac(arp)     
            macAlert.Times = 0
            macAlert.LastAlert = time.Now()
            macAlert.Alerted = true
        }
        Currentmacs[srcHw.String()] = macAlert
    }
}

func alertNewARP(arp *layers.ARP, alertabout int) {
    logs.Notice("alerting new arp")
    switch alertabout {
    case anIP:
        logs.Info("is IP - injecting ip alert")
    case aMAC:
        logs.Info("is MAC - injecting mac alert")
        if !isCurrentMac(arp) {
            addCurrentMac(arp)            
        }
        alertIfAlert(arp)
        WriteMacFileContent(Current)
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
            var srcHw net.HardwareAddr
 
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
                
                srcHw = arp.SourceHwAddress
                logs.Info("%v is at %v", net.IPv4(arp.SourceProtAddress[0], arp.SourceProtAddress[1], arp.SourceProtAddress[2], arp.SourceProtAddress[3]), srcHw)
            } else {
                logs.Error("unkonwn operation %d", arp.Operation)
            }
        }
    }
    return nil
}
 
func Init() {
    // ifaces, err := net.Interfaces()
    // if err != nil {
    //     logs.Error("MAC Management - error getting interfaces -> %s" + err.Error())
    //     return
    // }
    // logs.Info(ifaces)
//  macmanagement

    isEnabled, err := utils.GetKeyValueBool("macmanagement", "enabled") 
    if err != nil {
        logs.Error("AddPluginService Error getting data from main.conf: "+err.Error())
        arpmain.Enabled = true
    }else{
        arpmain.Enabled = isEnabled
    }
    iface, err := utils.GetKeyValueString("macmanagement", "interface") 
    if err != nil {
        logs.Error("AddPluginService Error getting data from main.conf: "+err.Error())
        arpmain.Interface = "eth0"
    }else{
        arpmain.Interface = iface
    }
    isLearning, err := utils.GetKeyValueBool("macmanagement", "learning") 
    if err != nil {
        logs.Error("AddPluginService Error getting data from main.conf: "+err.Error())
        arpmain.Learning = true
    }else{
        arpmain.Learning = isLearning
    }
    isOnneip, err := utils.GetKeyValueBool("macmanagement", "onnewip") 
    if err != nil {
        logs.Error("AddPluginService Error getting data from main.conf: "+err.Error())
        arpmain.Onnewip = true
    }else{
        arpmain.Enabled = isOnneip
    }
    isOnnewmac, err := utils.GetKeyValueBool("macmanagement", "onnewmac") 
    if err != nil {
        logs.Error("AddPluginService Error getting data from main.conf: "+err.Error())
        arpmain.Onnewmac = true
    }else{
        arpmain.Onnewmac = isOnnewmac
    }
    isVerbose, err := utils.GetKeyValueBool("macmanagement", "verbose") 
    if err != nil {
        logs.Error("AddPluginService Error getting data from main.conf: "+err.Error())
        arpmain.Verbose = true
    }else{
        arpmain.Verbose = isVerbose
    }
    timeBetweenAlerts, err := utils.GetKeyValueInt("macmanagement", "timebetweenalerts") 
    if err != nil {
        logs.Error("AddPluginService Error getting data from main.conf: "+err.Error())
        arpmain.Timebetweenalerts = 86400
    }else{
        arpmain.Timebetweenalerts = timeBetweenAlerts
    }
    dataKnownFile, err := utils.GetKeyValueString("macmanagement", "knownFile") 
    if err != nil {
        logs.Error("AddPluginService Error getting data from main.conf: "+err.Error())
        arpmain.KnownFile = "conf/known.db"
    }else{
        arpmain.KnownFile = dataKnownFile
    }
    dataCurrentFile, err := utils.GetKeyValueString("macmanagement", "currentFile") 
    if err != nil {
        logs.Error("AddPluginService Error getting data from main.conf: "+err.Error())
        arpmain.CurrentFile = "conf/current.db"
    }else{
        arpmain.CurrentFile = dataCurrentFile
    }
    
    // arpmain.Enabled = true
    // arpmain.Interface = "eth0"
    // arpmain.Learning = true
    // arpmain.Onnewip = true
    // arpmain.Onnewmac = true
    // arpmain.Verbose = true
    // arpmain.KnownFile = "conf/known.db"
    // arpmain.CurrentFile = "current.db"
 
    ReadMacFileContent(Known)
    ReadMacFileContent(Current)
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