package monitor
 
import (
    "github.com/astaxie/beego/logs"
    "strconv"
    "runtime"
    "fmt"
    "os"
    "time"
    "syscall"
    "owlhnode/database"
    "owlhnode/utils"
    "errors"
    // "github.com/pbnjay/memory"
    "github.com/shirou/gopsutil/cpu"
    //"github.com/shirou/gopsutil/disk"
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
    Percentage     string     `json:"percentage"`
}
 
type Disk struct {
    UsedDisk       string     `json:"useddisk"`
    TotalDisk      string     `json:"totaldisk"`
    FreeDisk       string     `json:"freedisk"`
    Percentage     string     `json:"percentage"`
}
 
type Cpu struct {
    Id             string     `json:"id"`
    Percentage     string     `json:"percentage"`
}
 
const (
    B  = 1
    KB = 1024 * B
    MB = 1024 * KB
    GB = 1024 * MB
)
 
type DiskStatus struct {
    All  uint64 `json:"all"`
    Used uint64 `json:"used"`
    Free uint64 `json:"free"`
}
 
var GlobalMonitor Monitor
 
func DiskUsage(path string) (disk DiskStatus) {
    fs := syscall.Statfs_t{}
    err := syscall.Statfs(path, &fs)
    if err != nil {
        return
    }
    disk.All = fs.Blocks * uint64(fs.Bsize)
    disk.Free = fs.Bfree * uint64(fs.Bsize)
    disk.Used = disk.All - disk.Free
    return
}
 
func doMonitor() {
    for {
        t,err := utils.GetKeyValueString("loop", "monitor")
        if err != nil {logs.Error("Search Error: Cannot load node information.")}
        tDuration, err := strconv.Atoi(t)
        time.Sleep(time.Second * time.Duration(tDuration))
 
        PrintMemUsage()
        PrintDiskUsage()
        PrintCPUUsage()       
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
    GlobalMonitor.Mem.Alloc = fmt.Sprintf("%v", bToMb(m.Alloc))
    GlobalMonitor.Mem.TotalAlloc = fmt.Sprintf("%v", bToMb(m.TotalAlloc))
    GlobalMonitor.Mem.Sys = fmt.Sprintf("%v", bToMb(m.Sys))
    GlobalMonitor.Mem.Gc = fmt.Sprintf("%v", m.NumGC)
 
    vmStat, _ := mem.VirtualMemory()
    GlobalMonitor.Mem.Percentage = fmt.Sprintf("%v",vmStat.UsedPercent) 
    GlobalMonitor.Mem.TotalMem = fmt.Sprintf("%v", bToMb(vmStat.Total))
    GlobalMonitor.Mem.FreeMem = fmt.Sprintf("%v", bToMb(vmStat.Free))
    GlobalMonitor.Mem.UsedMem = fmt.Sprintf("%v", bToMb(vmStat.Used))
}
 
func PercentFloat(val, total uint64)(percent uint64) {
    return (val * 100) / total
}
 
func PrintDiskUsage() {
    // diskStat, _ := disk.Usage("/")
    // GlobalMonitor.HD.Percentage = fmt.Sprintf("%v",diskStat.UsedPercent) 
    // GlobalMonitor.HD.FreeDisk = fmt.Sprintf("%v", bToMb(diskStat.Free))
    // GlobalMonitor.HD.TotalDisk = fmt.Sprintf("%v", bToMb(diskStat.Total))
    // GlobalMonitor.HD.UsedDisk = fmt.Sprintf("%v", bToMb(diskStat.Used))
 
    disk := DiskUsage("/")
    GlobalMonitor.HD.Percentage = fmt.Sprintf("%.2f", PercentFloat(disk.Used, disk.All)) 
    GlobalMonitor.HD.TotalDisk = fmt.Sprintf("%.2f", float64(disk.All)/float64(GB))
    GlobalMonitor.HD.UsedDisk = fmt.Sprintf("%.2f", float64(disk.Used)/float64(GB))
    GlobalMonitor.HD.FreeDisk = fmt.Sprintf("%.2f", float64(disk.Free)/float64(GB))
}
 
func PrintCPUUsage() {
    percentage, _ := cpu.Percent(0, true)
    GlobalMonitor.Cpus = nil
    for idx, cpupercent := range percentage {
        var acpu Cpu
        acpu.Id = fmt.Sprintf("%v",idx)
        acpu.Percentage = fmt.Sprintf("%v",cpupercent) 
        GlobalMonitor.Cpus = append(GlobalMonitor.Cpus, acpu)
    }
}
 
func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}
 
func AddMonitorFile(anode map[string]string)(err error) {
    files,err := ndb.LoadMonitorFiles()
    if err != nil {logs.Error("AddMonitorFile error getting files from database: %s", err.Error());return err}
    for x := range files{
        if files[x]["path"] == anode["path"]{
            return errors.New("Cannot add duplicated files")
        }
    }
    uuid := utils.Generate()
    err = ndb.InsertMonitorValue(uuid, "path", anode["path"]); if err != nil {logs.Error("AddMonitorFile error inserting path into database: %s", err.Error());return err}
    err = ndb.InsertMonitorValue(uuid, "maxSize", anode["maxSize"]); if err != nil {logs.Error("AddMonitorFile error inserting maxSize into database: %s", err.Error());return err}
    err = ndb.InsertMonitorValue(uuid, "maxDays", anode["maxDays"]); if err != nil {logs.Error("AddMonitorFile error inserting maxDays into database: %s", err.Error());return err}
    err = ndb.InsertMonitorValue(uuid, "rotate", anode["rotate"]); if err != nil {logs.Error("AddMonitorFile error inserting rotate into database: %s", err.Error());return err}
    err = ndb.InsertMonitorValue(uuid, "maxLines", anode["maxLines"]); if err != nil {logs.Error("AddMonitorFile error inserting maxLines into database: %s", err.Error());return err}
    
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
 
func ChangeRotationStatus(anode map[string]string)(err error) {
    err = ndb.UpdateMonitorFileValue(anode["file"], "rotate", anode["status"])
    if err != nil {logs.Error("ChangeRotationStatus error updating monitor file value: %s", err.Error());return err}
 
    return nil
}