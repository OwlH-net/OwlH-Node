package pcap

import (
    // "bytes"
    // "encoding/binary"
    "encoding/json"
    "errors"
    // "log"

    "bufio"
    "io/ioutil"
    "net"
    "os"
    "owlhnode/analyzer"
    "owlhnode/utils"
    "reflect"
    "regexp"
    // "sync"
    "github.com/astaxie/beego/logs"
    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "github.com/google/gopacket/pcap"
    "time"
)

type Oui struct {
    Oui                 string `json:"oui"`
    Isprivate           bool   `json:"isPrivate"`
    Companyname         string `json:"companyName"`
    Companyaddress      string `json:"companyAddress"`
    Countrycode         string `json:"countryCode"`
    Assignmentblocksize string `json:"assignmentBlockSize"`
    Datecreated         string `json:"dateCreated"`
    Dateupdated         string `json:"dateUpdated"`
}

var Ouis = map[string]Oui{}

type EventAlert struct {
    Timestamp   time.Time `json:"timestamp"`
    In_iface    string    `json:"in_iface"`
    Event_type  string    `json:"event_type"`
    Src_ip      string    `json:"src_ip"`
    Mac_address string    `json:"mac_address"`
    Proto       string    `json:"proto"`
    Alert       Alert     `json:"alert"`
}

type Alert struct {
    Action       string   `json:"action"`
    Signature_ID int      `json:"signature_id"`
    Rev          int      `json:"rev"`
    Signature    string   `json:"signature"`
    Category     string   `json:"category"`
    Severity     int      `json:"severity"`
    Times        int      `json:"times"`
    MultipleIP   []string `json:"multipleip"`
}

type Rule struct {
    Sid      int    `json:"sid"`
    Msg      string `json:"msg"`
    Severity int    `json:"severity"`
    Rev      int    `json:"int"`
    Category string `json:"category"`
}

type MAC struct {
    Mac       net.HardwareAddr `json:"mac"`
    IPs       map[string]IP    `json:"ips"`
    White     bool             `json:"white"`
    Onnewip   bool             `json:"onnewip"`
    Alerted   bool             `json:"alerted"`
    LastAlert time.Time        `json:"lastalert"`
    Times     int              `json:"times"`
    First     time.Time        `json:"first"`
    Last      time.Time        `json:"last"`
}

type IP struct {
    Ip        string    `json:"ip"`
    White     bool      `json:"white"`
    First     time.Time `json:"first"`
    Last      time.Time `json:"last"`
    Alerted   bool      `json:"alerted"`
    LastAlert time.Time `json:"lastalert"`
    Times     int       `json:"times"`
}

type TIPs map[string]IP

type ARPConfig struct {
    Onnewmac          bool
    Onnewip           bool
    Enabled           bool
    Learning          bool
    Verbose           bool
    Interface         string
    Timebetweenalerts int
    KnownFile         string
    CurrentFile       string
    Morethan          int
    OuiFile           string
    Ack               bool
}

const (
    aMAC = iota
    anIP
)

var ArpmainReturn ARPConfig
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

        return true
    } else {

        return false
    }
    return false
}

func isknownMac(arp *layers.ARP) (is bool) {
    var srcHw net.HardwareAddr
    srcHw = arp.SourceHwAddress

    if _, ok := Knownmacs[srcHw.String()]; ok {
        if arpmain.Verbose {
            logs.Info("KNOWN -> mac %s exists", srcHw.String())
        }
        return true
    } else {
        if arpmain.Verbose {
            logs.Info("KNOWN -> mac %s does not exist.", srcHw.String())
        }
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
            if arpmain.Verbose {
                logs.Info("ip %s does exist on mac %v", ip.Ip, Knownmacs[mac].Mac)
            }
            return true
        }
    }
    if arpmain.Verbose {
        logs.Info("ip %s does NOT exist on any mac", ip.Ip)
    }

    return false
}

func addKnownIP(arp *layers.ARP) {
    var srcHw net.HardwareAddr
    srcHw = arp.SourceHwAddress

    if mac, ok := Knownmacs[srcHw.String()]; ok {
        var newip IP
        newip.First = time.Now()
        newip.Last = time.Now()
        newip.Ip = net.IP(arp.SourceProtAddress).String()
        mac.IPs[newip.Ip] = newip
        Currentmacs[srcHw.String()] = mac
    }
}

func addCurrentMacIp(arp *layers.ARP) {
    var srcHw net.HardwareAddr
    srcHw = arp.SourceHwAddress

    // Manage MAC
    var newip IP
    var cMac MAC
    cMac = Currentmacs[srcHw.String()]

    newip.First = time.Now()
    newip.Last = time.Now()
    newip.Ip = net.IP(arp.SourceProtAddress).String()
    cMac.IPs[newip.Ip] = newip

    Currentmacs[srcHw.String()] = cMac
}

func addCurrentMac(arp *layers.ARP) {
    var srcHw net.HardwareAddr
    srcHw = arp.SourceHwAddress

    if arpmain.Verbose {
        logs.Info("CM adding mac %s ", srcHw.String())
    }

    // Manage MAC
    var cMac MAC
    isnow := time.Now()
    cMac.Mac = arp.SourceHwAddress
    cMac.First = isnow
    cMac.Last = isnow
    cMac.IPs = map[string]IP{}

    if arpmain.Verbose {
        logs.Info("CM mac added at %v", cMac.First)
    }

    // Manage new IP
    var newip IP

    newip.Ip = net.IP(arp.SourceProtAddress).String()
    newip.First = isnow
    newip.Last = isnow

    if arpmain.Verbose {
        logs.Info("CM lets add %s", newip.Ip)
        logs.Info("CM lets add  time %v", newip.First)
    }

    cMac.IPs[newip.Ip] = newip
    if arpmain.Verbose {
        logs.Info("CM - new mac details")
        logs.Warn("%+v", cMac)
    }
    Currentmacs[srcHw.String()] = cMac
    if arpmain.Verbose {
        logs.Warn("%+v", Currentmacs[srcHw.String()])
    }
}

func addKnownMac(arp *layers.ARP) {
    var srcHw net.HardwareAddr
    srcHw = arp.SourceHwAddress
    if arpmain.Verbose {
        logs.Info("adding mac %s ", srcHw.String())
    }
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

    allips := make(TIPs)
    allips[newip.Ip] = newip
    cMac.IPs = allips

    //    cMac.IPs = addKnownIP(cMac.IPs, newip)
    Knownmacs[srcHw.String()] = cMac
}

func updateLast(arp *layers.ARP, known bool) {
    var srcHw net.HardwareAddr
    srcHw = arp.SourceHwAddress

    var newIP IP
    newIP.Ip = net.IP(arp.SourceProtAddress).String()

    if known {
        cmac := Knownmacs[srcHw.String()]
        cmac.Last = time.Now()

        var oldIP IP
        if oldIP, ok := cmac.IPs[newIP.Ip]; ok {
            oldIP.Last = time.Now()
        } else {
            newIP.First = time.Now()
            newIP.Last = time.Now()
        }
        cmac.IPs[newIP.Ip] = oldIP
        Knownmacs[srcHw.String()] = cmac
    } else {
        cmac := Currentmacs[srcHw.String()]
        cmac.Last = time.Now()
        Currentmacs[srcHw.String()] = cmac
    }
}

func updateLastMac(arp *layers.ARP, known bool) {
    var srcHw net.HardwareAddr
    srcHw = arp.SourceHwAddress

    if known {
        cmac := Knownmacs[srcHw.String()]
        cmac.Last = time.Now()
        Knownmacs[srcHw.String()] = cmac
    } else {
        cmac := Currentmacs[srcHw.String()]
        cmac.Last = time.Now()
        Currentmacs[srcHw.String()] = cmac
    }
}

func updateLastIp(arp *layers.ARP, known bool) {
    var srcHw net.HardwareAddr
    srcHw = arp.SourceHwAddress
    var newIP IP
    newIP.Ip = net.IP(arp.SourceProtAddress).String()

    var oldIP IP
    var isnow = time.Now()

    if known {
        cmac := Knownmacs[srcHw.String()]
        if oldIP, ok := cmac.IPs[newIP.Ip]; ok {
            oldIP.Last = time.Now()
        } else {
            newIP.First = time.Now()
            newIP.Last = time.Now()
        }
        cmac.IPs[newIP.Ip] = oldIP
        Knownmacs[srcHw.String()] = cmac
    } else {
        cmac := Currentmacs[srcHw.String()]
        if arpmain.Verbose {
            logs.Debug("update last ip current -> cmac -> %v", cmac)
        }
        if oldIP, ok := cmac.IPs[newIP.Ip]; ok {
            oldIP.Last = isnow
            cmac.IPs[newIP.Ip] = oldIP
        } else {
            newIP.First = isnow
            newIP.Last = isnow
            cmac.IPs[newIP.Ip] = newIP
        }
        Currentmacs[srcHw.String()] = cmac
    }
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

func logVendor(arp *layers.ARP) (oui Oui, err error) {
    var loui Oui
    var srcHw net.HardwareAddr
    srcHw = arp.SourceHwAddress

    re := regexp.MustCompile(`^([^:]+:[^:]+:[^:]+:[^:]+:\w)`)
    hmac := re.FindStringSubmatch(srcHw.String())
    if len(hmac) > 0 {
        if foui, ok := Ouis[hmac[1]]; ok {
            if arpmain.Verbose {
                logs.Debug("Vendor -> %s", foui.Companyname)
                logs.Debug("Countrycode -> %s", foui.Countrycode)
            }
        } else {
            if arpmain.Verbose {
                logs.Debug("No mac vendor found")
            }
        }
    }
    return loui, nil
}

func alertMac(arp *layers.ARP) {
    var ev EventAlert
    var al Alert

    var srcHw net.HardwareAddr
    srcHw = arp.SourceHwAddress

    ev.Timestamp = time.Now()
    ev.In_iface = arpmain.Interface
    ev.Event_type = "alert"
    ev.Src_ip = net.IP(arp.SourceProtAddress).String()
    ev.Mac_address = srcHw.String()

    logVendor(arp)

    ev.Proto = "ARP"

    al.Action = "allowed"
    al.Signature_ID = 1150001
    al.Rev = 1
    al.Signature = "OwlH HWADD - new mac detected - " + ev.Mac_address
    al.Category = "Potentially Bad Traffic"
    al.Severity = 2
    al.Times = Currentmacs[srcHw.String()].Times

    ev.Alert = al

    values, _ := json.Marshal(ev)

    analyzer.ToDispatcher("start", string(values))
    if arpmain.Verbose {
        logs.Error(string(values))
    }
}

func alertIp(arp *layers.ARP) {
    var ev EventAlert
    var al Alert

    var srcHw net.HardwareAddr
    srcHw = arp.SourceHwAddress

    ev.Timestamp = time.Now()
    ev.In_iface = arpmain.Interface
    ev.Event_type = "alert"
    ev.Src_ip = net.IP(arp.SourceProtAddress).String()
    ev.Mac_address = srcHw.String()

    logVendor(arp)

    ev.Proto = "ARP"

    al.Action = "allowed"
    al.Signature_ID = 1150002
    al.Rev = 1
    al.Signature = "OwlH HWADD - new IP assigned to a MAC - " + ev.Mac_address
    al.Category = "Potentially Bad Traffic"
    al.Severity = 2
    al.Times = Currentmacs[srcHw.String()].Times

    if len(Currentmacs[srcHw.String()].IPs) > arpmain.Morethan {
        for ip := range Currentmacs[srcHw.String()].IPs {
            al.MultipleIP = append(al.MultipleIP, ip)
        }
    } else {
        return
    }

    ev.Alert = al

    values, _ := json.Marshal(ev)

    analyzer.ToDispatcher("start", string(values))
    if arpmain.Verbose {
        logs.Error(string(values))
    }
}

func timeToAlert(lastTime time.Time) bool {
    if arpmain.Verbose {
        logs.Info("Check time to alert")
        logs.Info(arpmain.Timebetweenalerts)
    }

    seconds := time.Second * time.Duration(arpmain.Timebetweenalerts)
    diff := time.Now().Sub(lastTime)

    if arpmain.Verbose {
        logs.Notice(seconds)
        logs.Warn(diff)
    }

    if diff >= seconds {
        return true
    }

    return false
}

func alertIfAlert(arp *layers.ARP) {
    var srcHw net.HardwareAddr
    srcHw = arp.SourceHwAddress

    if macAlert, ok := Currentmacs[srcHw.String()]; ok {
        if arpmain.Verbose {
            logs.Info("Alert new mac in current macs!")
        }
        if macAlert.Alerted {
            if timeToAlert(macAlert.LastAlert) {
                if arpmain.Verbose {
                    logs.Info("Time to alert again -- %d", macAlert.Times)
                }
                alertMac(arp)
                macAlert.Times = 0
                macAlert.LastAlert = time.Now()
            } else {
                macAlert.Times += 1
            }
        } else {
            if arpmain.Verbose {
                logs.Info("First time alert")
            }
            alertMac(arp)
            macAlert.Times = 0
            macAlert.LastAlert = time.Now()
            macAlert.Alerted = true
        }
        Currentmacs[srcHw.String()] = macAlert
    }
}

func alertIfAlertIp(arp *layers.ARP) {
    var srcHw net.HardwareAddr
    srcHw = arp.SourceHwAddress

    if macAlert, ok := Currentmacs[srcHw.String()]; ok {
        if arpmain.Verbose {
            logs.Info("Alert new mac in current macs!")
        }
        var newIP IP
        newIP.Ip = net.IP(arp.SourceProtAddress).String()
        if ipAlert, ok := macAlert.IPs[newIP.Ip]; ok {
            if ipAlert.Alerted {
                if timeToAlert(ipAlert.LastAlert) {
                    if arpmain.Verbose {
                        logs.Info("Time to alert again -- %d", macAlert.Times)
                    }
                    alertIp(arp)
                    ipAlert.Times = 0
                    ipAlert.LastAlert = time.Now()
                } else {
                    ipAlert.Times += 1
                }
            } else {
                if arpmain.Verbose {
                    logs.Info("First time alert")
                }
                alertIp(arp)
                ipAlert.Times = 0
                ipAlert.LastAlert = time.Now()
                ipAlert.Alerted = true
            }
            macAlert.IPs[newIP.Ip] = ipAlert
            Currentmacs[srcHw.String()] = macAlert
            WriteMacFileContent(Current)
        }
    }
}

func alertIfAlertMac(arp *layers.ARP) {
    var srcHw net.HardwareAddr
    srcHw = arp.SourceHwAddress

    if macAlert, ok := Currentmacs[srcHw.String()]; ok {
        if arpmain.Verbose {
            logs.Info("Alert new mac in current macs!")
        }
        if macAlert.Alerted {
            if timeToAlert(macAlert.LastAlert) {
                if arpmain.Verbose {
                    logs.Info("Time to alert again -- %d", macAlert.Times)
                }
                alertMac(arp)
                macAlert.Times = 0
                macAlert.LastAlert = time.Now()
            } else {
                macAlert.Times += 1
            }
        } else {
            if arpmain.Verbose {
                logs.Info("First time alert")
            }
            alertMac(arp)
            macAlert.Times = 0
            macAlert.LastAlert = time.Now()
            macAlert.Alerted = true
        }
        Currentmacs[srcHw.String()] = macAlert
        WriteMacFileContent(Current)
    }

}

func alertNewARP(arp *layers.ARP, alertabout int) {
    switch alertabout {
    case anIP:
        alertIfAlertIp(arp)
    case aMAC:
        alertIfAlertMac(arp)
    default:
        logs.Error("have no idea what we try to alert about %v", alertabout)
    }
    return
}

func isKnowMACIP(arp *layers.ARP) (isknown bool) {

    var srcHw net.HardwareAddr
    srcHw = arp.SourceHwAddress

    cMac := Knownmacs[srcHw.String()]
    if _, ok := cMac.IPs[net.IP(arp.SourceProtAddress).String()]; ok {
        return true
    }
    return false
}

func isCurrentMacIp(arp *layers.ARP) (isknown bool) {

    var srcHw net.HardwareAddr
    srcHw = arp.SourceHwAddress

    cMac := Currentmacs[srcHw.String()]
    if _, ok := cMac.IPs[net.IP(arp.SourceProtAddress).String()]; ok {
        return true
    }

    return false
}

func checkarp(arp *layers.ARP) {
    logVendor(arp)

    if !isknownMac(arp) {
        if arpmain.Verbose {
            logs.Debug("MAC is not known")
        }
        if !isCurrentMac(arp) {
            if arpmain.Verbose {
                logs.Debug("MAC is not known and is not current - add mac to current")
            }

            addCurrentMac(arp)
            if arpmain.Onnewmac {
                if arpmain.Verbose {
                    logs.Debug("MAC is not known and is not current - alert new mac ")
                }
                alertNewARP(arp, aMAC)
            }
        } else {
            if arpmain.Verbose {
                logs.Debug("MAC is not known but is current - update last mac known")
            }
            updateLastMac(arp, false)
            if isCurrentMacIp(arp) {
                if arpmain.Verbose {
                    logs.Debug("IP is current - update last ip current ")
                }
                updateLastIp(arp, false)
            } else {
                if arpmain.Verbose {
                    logs.Debug("IP is not current current - add ip current ")
                }
                addCurrentMacIp(arp)
            }
            if arpmain.Onnewip {
                if arpmain.Verbose {
                    logs.Debug("IP - alert new ip")
                }
                alertNewARP(arp, anIP)
            }
        }
    } else {
        if arpmain.Verbose {
            logs.Debug("MAC is known")
        }
        if !isCurrentMac(arp) {
            if arpmain.Verbose {
                logs.Debug("MAC is known but is not Current - add mac to current")
            }
            addCurrentMac(arp)
        } else {
            if arpmain.Verbose {
                logs.Debug("MAC is known and is Current - udpate last mac Current")
            }
            updateLastMac(arp, false)
        }
        if arpmain.Verbose {
            logs.Debug("MAC is known - udpate last mac known ")
        }
        updateLastMac(arp, true)
        if !isKnowMACIP(arp) {
            if arpmain.Verbose {
                logs.Debug("IP is known but is not Current - add ip to current ")
            }
            addCurrentMacIp(arp)
            if arpmain.Onnewip {
                if arpmain.Verbose {
                    logs.Debug("IP is known but is not Current - alert on new ip ")
                }
                alertNewARP(arp, anIP)
            }
        } else {
            if arpmain.Verbose {
                logs.Debug("IP is known and is Current - update last ip knwon and current ")
            }
            updateLastIp(arp, false)
            updateLastIp(arp, true)
        }
        if arpmain.Verbose {
            logs.Debug("wirte konwn")
        }
        WriteMacFileContent(Known)
    }
    if arpmain.Verbose {
        logs.Debug("wirte konwn")
    }
    WriteMacFileContent(Current)
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
            if !arpmain.Enabled {
                return nil
            }
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
                    if arpmain.Verbose {
                        logs.Info("learning")
                    }
                    learnarp(arp)
                } else {
                    if arpmain.Verbose {
                        logs.Info("live")
                    }
                    checkarp(arp)
                }

                srcHw = arp.SourceHwAddress
                if arpmain.Verbose {
                    logs.Info("%v is at %v", net.IPv4(arp.SourceProtAddress[0], arp.SourceProtAddress[1], arp.SourceProtAddress[2], arp.SourceProtAddress[3]), srcHw)
                }
            } else {
                if arpmain.Verbose {
                    logs.Error("unkonwn operation %d", arp.Operation)
                }
            }
        }
    }
    return nil
}

func addMac(mac string) {
    if _, ok := Knownmacs[mac]; !ok {
        macObj, _ := net.ParseMAC(mac)
        var newMAC MAC

        newMAC.Mac = macObj
        newMAC.First = time.Now()
        newMAC.Last = time.Now()

        Knownmacs[mac] = newMAC
        WriteMacFileContent(Known)
    }
}

func addIp(ip string, mac string) {
    if arpmain.Verbose {
        logs.Info("addIP function -> %s, %s", ip, mac)
    }
    if _, ok := Knownmacs[mac]; ok {
        if arpmain.Verbose {
            logs.Info("MAC EXISTS!! addIP function -> %s, %s", ip, mac)
        }
        if _, ok := Knownmacs[mac].IPs[ip]; !ok {
            if arpmain.Verbose {
                logs.Info("IP DOESN'T EXISTS!! addIP function -> %s, %s", ip, mac)
            }
            var newip IP
            newip.Ip = ip
            newip.First = time.Now()
            newip.Last = time.Now()

            mymac := Knownmacs[mac]
            cips := make(map[string]IP)

            for localip := range mymac.IPs {
                cips[localip] = mymac.IPs[localip]
            }

            cips[ip] = newip

            logs.Warn("cips2 -> %+v", cips)

            mymac.IPs = cips
            Knownmacs[mac] = mymac

            WriteMacFileContent(Known)
        }
    }
}

func AddMacIp(anode map[string]string) error {
    if _, ok := anode["mac"]; ok {
        addMac(anode["mac"])
        if _, okIP := anode["ip"]; okIP {
            addIp(anode["ip"], anode["mac"])
        }
    }
    return nil
}

func Db(data map[string]string) {
    for x := range data {
        if data[x] == "reload" {
            if x == "Known" {
                ReadMacFileContent(Known)
            } else if x == "Current" {
                ReadMacFileContent(Current)
            }
        } else if data[x] == "reset" {
            if x == "Known" {
                Knownmacs = map[string]MAC{}
                WriteMacFileContent(Known)
            } else if x == "Current" {
                Currentmacs = map[string]MAC{}
                WriteMacFileContent(Current)
            }
        }
    }
}

func Config(config map[string]interface{}) {
    v := reflect.ValueOf(&arpmain).Elem()

    for i := 0; i < v.NumField(); i++ {
        name := v.Type().Field(i).Name
        if _, ok := config[name]; ok {
            f1 := v.Field(i)
            f1.Set(reflect.ValueOf(config[name]))
        }
    }

    ArpmainReturn = arpmain
}

func LoadConfig() {
    isEnabled, err := utils.GetKeyValueBool("macmanagement", "enabled")
    if err != nil {
        logs.Error("AddPluginService Error getting data from main.conf: " + err.Error())
    } else {
        arpmain.Enabled = isEnabled
    }
    iface, err := utils.GetKeyValueString("macmanagement", "interface")
    if err != nil {
        logs.Error("AddPluginService Error getting data from main.conf: " + err.Error())
    } else {
        arpmain.Interface = iface
    }
    isLearning, err := utils.GetKeyValueBool("macmanagement", "learning")
    if err != nil {
        logs.Error("AddPluginService Error getting data from main.conf: " + err.Error())
    } else {
        arpmain.Learning = isLearning
    }
    isOnneip, err := utils.GetKeyValueBool("macmanagement", "onnewip")
    if err != nil {
        logs.Error("AddPluginService Error getting data from main.conf: " + err.Error())
    } else {
        arpmain.Enabled = isOnneip
    }
    isOnnewmac, err := utils.GetKeyValueBool("macmanagement", "onnewmac")
    if err != nil {
        logs.Error("AddPluginService Error getting data from main.conf: " + err.Error())
    } else {
        arpmain.Onnewmac = isOnnewmac
    }
    isVerbose, err := utils.GetKeyValueBool("macmanagement", "verbose")
    if err != nil {
        logs.Error("AddPluginService Error getting data from main.conf: " + err.Error())
    } else {
        arpmain.Verbose = isVerbose
    }
    timeBetweenAlerts, err := utils.GetKeyValueInt("macmanagement", "timebetweenalerts")
    if err != nil {
        logs.Error("AddPluginService Error getting data from main.conf: " + err.Error())
    } else {
        arpmain.Timebetweenalerts = timeBetweenAlerts
    }
    dataKnownFile, err := utils.GetKeyValueString("macmanagement", "knownFile")
    if err != nil {
        logs.Error("AddPluginService Error getting data from main.conf: " + err.Error())
    } else {
        arpmain.KnownFile = dataKnownFile
    }
    dataCurrentFile, err := utils.GetKeyValueString("macmanagement", "currentFile")
    if err != nil {
        logs.Error("AddPluginService Error getting data from main.conf: " + err.Error())
    } else {
        arpmain.CurrentFile = dataCurrentFile
    }
    moreThan, err := utils.GetKeyValueInt("macmanagement", "morethan")
    if err != nil {
        logs.Error("AddPluginService Error getting data from main.conf: " + err.Error())
    } else {
        arpmain.Morethan = moreThan
    }
    dataOui, err := utils.GetKeyValueString("macmanagement", "ouifile")
    if err != nil {
        logs.Error("AddPluginService Error getting data from main.conf: " + err.Error())
    } else {
        arpmain.OuiFile = dataOui
    }
    arpmain.Ack = true

    ArpmainReturn = arpmain
}

func Init() {

    arpmain.Enabled = true
    arpmain.Interface = "eth0"
    arpmain.Learning = true
    arpmain.Onnewip = true
    arpmain.Onnewmac = true
    arpmain.Verbose = true
    arpmain.Timebetweenalerts = 86400
    arpmain.KnownFile = "conf/known.db"
    arpmain.CurrentFile = "conf/current.db"
    arpmain.OuiFile = "conf/ouis.db"
    arpmain.Morethan = 1

    LoadConfig()

    // arpmain.Enabled = true
    // arpmain.Interface = "eth0"
    // arpmain.Learning = true
    // arpmain.Onnewip = true
    // arpmain.Onnewmac = true
    // arpmain.Verbose = true
    // arpmain.KnownFile = "conf/known.db"
    // arpmain.CurrentFile = "current.db"
    if !arpmain.Enabled {
        return
    }

    ReadOUI()
    ReadMacFileContent(Known)
    ReadMacFileContent(Current)
    go readARP("")

    // readARP("")
}

const (
    Known = iota
    Current
)

func ReadMacFileContent(alertabout int) (err error) {
    switch alertabout {
    case Known:
        byteValue, err := ioutil.ReadFile(arpmain.KnownFile)
        if err != nil {
            logs.Error("ReadFileContent - error getting KnownFile content -> %s", err.Error())
            return err
        }
        json.Unmarshal(byteValue, &Knownmacs)
    case Current:
        byteValue, err := ioutil.ReadFile(arpmain.CurrentFile)
        if err != nil {
            logs.Error("ReadFileContent - error getting CurrentFile content -> %s", err.Error())
            return err
        }
        json.Unmarshal(byteValue, &Currentmacs)
    default:
        return errors.New("ReadMacFileContent Invalid unmarshal variable")
    }

    return nil
}

func WriteMacFileContent(alertabout int) (err error) {
    switch alertabout {
    case Known:
        values, _ := json.Marshal(Knownmacs)
        err = ioutil.WriteFile(arpmain.KnownFile, values, 0644)
        if err != nil {
            logs.Error("ReadFileContent - error writing file content -> %s" + err.Error())
            return err
        }
    case Current:
        values, _ := json.Marshal(Currentmacs)
        if arpmain.Verbose {
            logs.Warn("current -> %+v", string(values))
        }
        err = ioutil.WriteFile(arpmain.CurrentFile, values, 0644)
        if err != nil {
            logs.Error("ReadFileContent - error writing file content -> %s" + err.Error())
            return err
        }
    default:
        return errors.New("ReadMacFileContent Invalid unmarshal variable")
    }

    return nil
}

func ReadOUI() (err error) {
    file, err := os.Open(arpmain.OuiFile)
    if err != nil {
        return err
    }
    defer file.Close()

    var newoui Oui
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        json.Unmarshal([]byte(scanner.Text()), &newoui)
        Ouis[newoui.Oui] = newoui
    }

    if err := scanner.Err(); err != nil {
        return err
    }

    logs.Warn("imported %d OUIs", len(Ouis))

    return nil
}
