package main

import (
	"fmt"

	"github.com/bar-counter/monitor/v3"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	monitorCfg := &monitor.Cfg{
		Status: true,
		//StatusPrefix: "/status",
		StatusHardware: true,
		//StatusHardwarePrefix: "/hardware",
		Debug: true,
		//DebugPrefix: "/debug",
		DebugMiddleware: gin.BasicAuth(gin.Accounts{
			"admin": "admin",
			"user":  "user",
		}),
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
