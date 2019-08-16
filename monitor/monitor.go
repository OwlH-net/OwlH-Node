package monitor

import (
    "github.com/astaxie/beego/logs"
    "strconv"
    "runtime"
    "time"
    "github.com/pbnjay/memory"
    "github.com/shirou/gopsutil/cpu"
    "github.com/shirou/gopsutil/disk"
    "github.com/shirou/gopsutil/mem"
)

type Monitor struct {
    Mem            Memory     `json:"mem"`
    HD             Disk       `json:"disk"`
    Cpus           []Cpu      `json:"cpus"`
}

type Memory struct {
    Alloc          uint64     `json:"alloc"`
    TotalAlloc     uint64     `json:"totalalloc"`
    Sys            uint64     `json:"sys"`
    Gc             uint32     `json:"gc"`
    TotalMem       uint64     `json:"totalmem"`
    FreeMem        uint64     `json:"freemem"`
    UsedMem        uint64     `json:"usedmem"`
    Percentage     float64    `json:"percentage"`
}

type Disk struct {
    UsedDisk       uint64     `json:"useddisk"`
    TotalDisk      uint64     `json:"totaldisk"`
    FreeDisk       uint64     `json:"freedisk"`
    Percentage     float64    `json:"percentage"`
}

type Cpu struct {
    Id             int        `json:"id"`
    Percentage     float64    `json:"percentage"`
}

var GlobalMonitor Monitor

func doMonitor() {
    for {
        time.Sleep(time.Second * 20)
        PrintMemUsage()
        PrintDiskUsage()
        PrintCPUUsage()
        testMonitor := GetLastMonitorInfo()
        logs.Notice(" ===== Monitor ===== ")
        logs.Notice(testMonitor)
    }
}

func Init() {
    logs.Info("Monitor -> Starting Monitor Service")
    go doMonitor()
}

func GetLastMonitorInfo()(GMonitor Monitor) {
    return GlobalMonitor
}

func PrintMemUsage() {
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    logs.Info("Monitor -> Mem Stats")
    logs.Notice("Alloc = %v MiB\tTotalAlloc = %v MiB\tSys = %v MiB\tNumGC = %v\tTotal system memory: %d MiB", bToMb(m.Alloc),bToMb(m.TotalAlloc),bToMb(m.Sys), m.NumGC,bToMb(memory.TotalMemory()))
    GlobalMonitor.Mem.Alloc = m.Alloc
    GlobalMonitor.Mem.TotalAlloc = m.TotalAlloc
    GlobalMonitor.Mem.Sys = m.Sys
    GlobalMonitor.Mem.Gc = m.NumGC

    vmStat, _ := mem.VirtualMemory()
    logs.Notice("Total Mem: %v MiB\t Free Mem: %v MiB\t Used Mem percentage: %s", bToMb(vmStat.Total), bToMb(vmStat.Free),  strconv.FormatFloat(vmStat.UsedPercent, 'f', 2, 64))
    GlobalMonitor.Mem.Percentage = vmStat.UsedPercent
    GlobalMonitor.Mem.TotalMem = vmStat.Total
    GlobalMonitor.Mem.FreeMem = vmStat.Free
    GlobalMonitor.Mem.UsedMem = vmStat.Used
}

func PrintDiskUsage() {
    diskStat, _ := disk.Usage("/")
    logs.Notice("Total Disk: %v MiB\t Used Disk: %v MiB\t Free Disk: %v MiB\t Used Disk percentage: %s", bToMb(diskStat.Total), bToMb(diskStat.Used), bToMb(diskStat.Free),strconv.FormatFloat(diskStat.UsedPercent, 'f', 2, 64)) 
    GlobalMonitor.HD.FreeDisk = diskStat.Free
    GlobalMonitor.HD.Percentage = diskStat.UsedPercent
    GlobalMonitor.HD.TotalDisk = diskStat.Total
    GlobalMonitor.HD.UsedDisk = diskStat.Used
}

func PrintCPUUsage() {
    cpuStat, _ := cpu.Info()
    percentage, _ := cpu.Percent(0, true)
    logs.Notice("CPU cores: %v ", int64(cpuStat[0].Cores+1))
    GlobalMonitor.Cpus = nil
    for idx, cpupercent := range percentage {
        logs.Info( "\t\tCurrent CPU utilization: [" + strconv.Itoa(idx) + "] " + strconv.FormatFloat(cpupercent, 'f', 2, 64) + "%")
        var acpu Cpu
        acpu.Id = idx
        acpu.Percentage = cpupercent
        GlobalMonitor.Cpus = append(GlobalMonitor.Cpus, acpu)
    }
}

func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}