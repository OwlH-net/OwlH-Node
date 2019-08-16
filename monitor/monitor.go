package monitor

import (
    "github.com/astaxie/beego/logs"
    "runtime"
    "time"
    "github.com/pbnjay/memory"
)

func monitor() {
    for {
        time.Sleep(time.Second * 20)
        PrintMemUsage()
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
        logs.Notice("Alloc = %v MiB\tTotalAlloc = %v MiB\tSys = %v MiB\tNumGC = %v\tTotal system memory: %d", bToMb(m.Alloc),bToMb(m.TotalAlloc),bToMb(m.Sys), m.NumGC,bToMb(memory.TotalMemory()))
}

func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}