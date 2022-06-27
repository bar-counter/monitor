package main

import (
	"fmt"

	"github.com/bar-counter/monitor/v2"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	monitorCfg := &monitor.Cfg{
		Status: true,
		//StatusPrefix: "/status",
		StatusHardware: true,
		//StatusHardwarePrefix: "/hardware",
	}
	err := monitor.Register(r, monitorCfg)
	if err != nil {
		fmt.Printf("monitor register err %v\n", err)
		return
	}
	err = r.Run(":38000")
	if err != nil {
		fmt.Printf("run err %v\n", err)
		return
	}
}
