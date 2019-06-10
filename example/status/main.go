package main

import (
	"fmt"
	"github.com/bar-counter/monitor"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	cfg := &monitor.Cfg{
		Status: true,
		//StatusPrefix: "/status",
	}
	err := monitor.Register(r, cfg)
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
