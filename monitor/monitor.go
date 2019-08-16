package monitor

import (
    "github.com/astaxie/beego/logs"
    "runtime"
    "time"
    "github.com/pbnjay/memory"
    "github.com/shirou/gopsutil/cpu"
    "github.com/shirou/gopsutil/disk"
    "github.com/shirou/gopsutil/host"
    "github.com/shirou/gopsutil/mem"
    "github.com/shirou/gopsutil/net"
)

func monitor() {
    for {
        time.Sleep(time.Second * 20)
        PrintMemUsage()
        PrintDiskUsage()
    }
}

func Init() {
    logs.Info("Monitor -> Starting Monitor Service")
    go monitor()
}


func PrintMemUsage() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    logs.Info("Monitor -> Mem Stats")
    logs.Notice("Alloc = %v MiB\tTotalAlloc = %v MiB\tSys = %v MiB\tNumGC = %v\tTotal system memory: %d MiB", bToMb(m.Alloc),bToMb(m.TotalAlloc),bToMb(m.Sys), m.NumGC,bToMb(memory.TotalMemory()))
}

func PrintDiskUsage() {
    diskStat, _ := disk.Usage("/")
    logs.Notice("Total Disk: %s\t Used Disk: %s\t Free Disk: %s", strconv.FormatUint(diskStat.Total, 10), strconv.FormatUint(diskStat.Used, 10), strconv.FormatUint(diskStat.Free, 10)) 
}

func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}