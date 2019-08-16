package monitor

import (
    "github.com/astaxie/beego/logs"
    "runtime"
    "time"
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
        logs.Notice("Alloc = %v MiB", bToMb(m.Alloc))
        logs.Notice("\tTotalAlloc = %v MiB", bToMb(m.TotalAlloc))
        logs.Notice("\tSys = %v MiB", bToMb(m.Sys))
        logs.Notice("\tNumGC = %v\n", m.NumGC)
}

func bToMb(b uint64) uint64 {
    return b / 1024 / 1024
}