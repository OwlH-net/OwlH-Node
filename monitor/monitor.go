package monitor

import (
    "github.com/astaxie/beego/logs"
    "strconv"
    "runtime"
    "fmt"
    "os"
    "time"
    "owlhnode/database"
    "owlhnode/utils"
    // "github.com/pbnjay/memory"
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
    Alloc          string     `json:"alloc"`
    TotalAlloc     string     `json:"totalalloc"`
    Sys            string     `json:"sys"`
    Gc             string     `json:"gc"`
    TotalMem       string     `json:"totalmem"`
    FreeMem        string     `json:"freemem"`
    UsedMem        string     `json:"usedmem"`
    Percentage     string    `json:"percentage"`
}

type Disk struct {
    UsedDisk       string     `json:"useddisk"`
    TotalDisk      string     `json:"totaldisk"`
    FreeDisk       string     `json:"freedisk"`
    Percentage     string    `json:"percentage"`
}

type Cpu struct {
    Id             string        `json:"id"`
    Percentage     string    `json:"percentage"`
}
// type Monitor struct {
//     Mem            Memory     `json:"mem"`
//     HD             Disk       `json:"disk"`
//     Cpus           []Cpu      `json:"cpus"`
// }

// type Memory struct {
//     Alloc          uint64     `json:"alloc"`
//     TotalAlloc     uint64     `json:"totalalloc"`
//     Sys            uint64     `json:"sys"`
//     Gc             uint32     `json:"gc"`
//     TotalMem       uint64     `json:"totalmem"`
//     FreeMem        uint64     `json:"freemem"`
//     UsedMem        uint64     `json:"usedmem"`
//     Percentage     float64    `json:"percentage"`
// }

// type Disk struct {
//     UsedDisk       uint64     `json:"useddisk"`
//     TotalDisk      uint64     `json:"totaldisk"`
//     FreeDisk       uint64     `json:"freedisk"`
//     Percentage     float64    `json:"percentage"`
// }

// type Cpu struct {
//     Id             int        `json:"id"`
//     Percentage     float64    `json:"percentage"`
// }

var GlobalMonitor Monitor

func doMonitor() {
    for {
        time.Sleep(time.Second * 20)
        PrintMemUsage()
        PrintDiskUsage()
        PrintCPUUsage()
        // testMonitor := GetLastMonitorInfo()
        // logs.Notice(" ===== Monitor ===== ")
        // logs.Notice(testMonitor)
        
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
    // logs.Info("Monitor -> Mem Stats")
    // logs.Notice("Alloc = %v MiB\tTotalAlloc = %v MiB\tSys = %v MiB\tNumGC = %v\tTotal system memory: %d MiB", bToMb(m.Alloc),bToMb(m.TotalAlloc),bToMb(m.Sys), m.NumGC,bToMb(memory.TotalMemory()))
    GlobalMonitor.Mem.Alloc = fmt.Sprintf("%v", bToMb(m.Alloc))
    GlobalMonitor.Mem.TotalAlloc = fmt.Sprintf("%v", bToMb(m.TotalAlloc))
    GlobalMonitor.Mem.Sys = fmt.Sprintf("%v", bToMb(m.Sys))
    GlobalMonitor.Mem.Gc = fmt.Sprintf("%v", m.NumGC)

    vmStat, _ := mem.VirtualMemory()
    // logs.Notice("Total Mem: %v MiB\t Free Mem: %v MiB\t Used Mem percentage: %s", bToMb(vmStat.Total), bToMb(vmStat.Free),  strconv.FormatFloat(vmStat.UsedPercent, 'f', 2, 64))
    GlobalMonitor.Mem.Percentage = fmt.Sprintf("%v",vmStat.UsedPercent) 
    GlobalMonitor.Mem.TotalMem = fmt.Sprintf("%v", bToMb(vmStat.Total))
    GlobalMonitor.Mem.FreeMem = fmt.Sprintf("%v", bToMb(vmStat.Free))
    GlobalMonitor.Mem.UsedMem = fmt.Sprintf("%v", bToMb(vmStat.Used))
}

func PrintDiskUsage() {
    diskStat, _ := disk.Usage("/")
    // logs.Notice("Total Disk: %v MiB\t Used Disk: %v MiB\t Free Disk: %v MiB\t Used Disk percentage: %s", bToMb(diskStat.Total), bToMb(diskStat.Used), bToMb(diskStat.Free),strconv.FormatFloat(diskStat.UsedPercent, 'f', 2, 64)) 
    GlobalMonitor.HD.Percentage = fmt.Sprintf("%v",diskStat.UsedPercent) 
    GlobalMonitor.HD.FreeDisk = fmt.Sprintf("%v", bToMb(diskStat.Free))
    GlobalMonitor.HD.TotalDisk = fmt.Sprintf("%v", bToMb(diskStat.Total))
    GlobalMonitor.HD.UsedDisk = fmt.Sprintf("%v", bToMb(diskStat.Used))
}

func PrintCPUUsage() {
    // cpuStat, _ := cpu.Info()
    percentage, _ := cpu.Percent(0, true)
    // logs.Notice("CPU cores: %v ", int64(cpuStat[0].Cores+1))
    GlobalMonitor.Cpus = nil
    for idx, cpupercent := range percentage {
        // logs.Info( "\t\tCurrent CPU utilization: [" + strconv.Itoa(idx) + "] " + strconv.FormatFloat(cpupercent, 'f', 2, 64) + "%")
        var acpu Cpu
        acpu.Id = fmt.Sprintf("%v",idx)
        acpu.Percentage = fmt.Sprintf("%v",cpupercent) 
        GlobalMonitor.Cpus = append(GlobalMonitor.Cpus, acpu)
        ///////////////////////////////////////////////////////
        // acpu.Id = fmt.Sprintf("%v",1)
        // acpu.Percentage = fmt.Sprintf("%v",11) 
        // GlobalMonitor.Cpus = append(GlobalMonitor.Cpus, acpu)
        // acpu.Id = fmt.Sprintf("%v",2)
        // acpu.Percentage = fmt.Sprintf("%v",3) 
        // GlobalMonitor.Cpus = append(GlobalMonitor.Cpus, acpu)
        ///////////////////////////////////////////////////////
    }
}

func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}

func AddMonitorFile(anode map[string]string)(err error) {
    uuid := utils.Generate()
    err = ndb.InsertMonitorValue(uuid, "path", anode["path"])
    if err != nil {logs.Error("AddMonitorFile error inserting new path into database: %s", err.Error());return err}
    
    return nil
}

func DeleteMonitorFile(anode map[string]string)(err error) {
    err = ndb.DeleteMonitorFile(anode["file"])
    if err != nil {logs.Error("DeleteMonitorFile error inserting new path into database: %s", err.Error());return err}
    
    return nil
}

func PingMonitorFiles()(data map[string]map[string]string, err error) {
    data,err = ndb.LoadMonitorFiles()
    if err != nil {logs.Error("PingMonitorFiles error getting monitor paths: %s", err.Error());return nil,err}

    for x := range data {
        fi, err := os.Stat(data[x]["path"]);
        if err != nil {
            data[x]["size"] = "-1"
        }else{
            data[x]["size"] = strconv.FormatInt(fi.Size(), 10)
        }
    }

    return data,err
}