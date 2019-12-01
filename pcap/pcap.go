package pcap


import (
    // "bytes"
    // "encoding/binary"
    // "errors"
    // "log"
    "net"
    // "sync"
    // "time"
    "github.com/astaxie/beego/logs"
    "github.com/google/gopacket"
    "github.com/google/gopacket/layers"
    "github.com/google/gopacket/pcap"
)


func readARP(iface string)(err error) {
    // var handle *pcap.Handle
    if handle1, err := pcap.OpenLive(iface, 65536, true, pcap.BlockForever); err != nil {
        logs.Error(err)
        return err
    } else {
        handle1.SetBPFFilter("arp")
        src := gopacket.NewPacketSource(handle1, layers.LayerTypeEthernet)
        in := src.Packets()
    
    // else if err := handle.SetBPFFilter("arp"); err != nil {  // optional
    //   logs.Error(err)
    //   return err
    // } 

    for {
        packet := <-in
        logs.Info ("new packet")
        arpLayer := packet.Layer(layers.LayerTypeARP)
        if arpLayer == nil {
            continue
        }
        arp := arpLayer.(*layers.ARP)
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
        var srchw net.HardwareAddr
        var dsthw net.HardwareAddr

        // logs.Info("arp.DstProtAddress    %v",arp.DstProtAddress   ) 
        if arp.Operation == 1 {
            srchw = arp.SourceHwAddress
            logs.Info("who has %v tells %v(%v)", net.IPv4(arp.DstProtAddress[0],arp.DstProtAddress[1],arp.DstProtAddress[2],arp.DstProtAddress[3]), net.IPv4(arp.SourceProtAddress[0],arp.SourceProtAddress[1],arp.SourceProtAddress[2],arp.SourceProtAddress[3]), srchw.String())
        }else if arp.Operation == 2{
            // dsthw := string(arp.DstHwAddress[0])+":"+string(arp.DstHwAddress[1])+":"+string(arp.DstHwAddress[2])+":"+string(arp.DstHwAddress[3])+":"+string(arp.DstHwAddress[4])+":"+string(arp.DstHwAddress[5])
            dsthw = arp.DstHwAddress
            logs.Info("%v is at %v", net.IPv4(arp.DstProtAddress[0],arp.DstProtAddress[1],arp.DstProtAddress[2],arp.DstProtAddress[3]), dsthw)
        }else {
            logs.Error("unkonwn operation %d", arp.Operation)
        }

        // if arp.Operation != layers.ARPReply || bytes.Equal([]byte(iface.HardwareAddr), arp.SourceHwAddress) {
        //     // This is a packet I sent.
        //     continue
        // }
        // Note:  we might get some packets here that aren't responses to ones we've sent,
        // if for example someone else sends US an ARP request.  Doesn't much matter, though...
        // all information is good information :)
        // log.Printf("IP %v is at %v", net.IP(arp.SourceProtAddress), net.HardwareAddr(arp.SourceHwAddress))

    }
    }
    return nil
}

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

func Init() {
    // Get a list of all interfaces.
    ifaces, err := net.Interfaces()
    if err != nil {
        logs.Error("error getting interfaces -> %s"+err.Error())
        return
    }
    logs.Info(ifaces)

    go readARP("eth0")

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
}