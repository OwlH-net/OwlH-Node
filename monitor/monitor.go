package monitor

import (

)

func monitor() {
    for {
        
    }
}

func Init() {
    logs.Info("Monitor -> Starting Monitor Service")
    go monitor()
}
