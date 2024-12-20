//go:build !test

package main

import (
	_ "embed"
	"flag"
	"fmt"
	"github.com/bar-counter/monitor/v2"
	"github.com/gin-gonic/gin"
	"log"
	"os"
)

const cliVersion = "0.1.2"

var serverPort = flag.String("serverPort", "49002", "http service address")

var buildID string

func init() {
	if buildID == "" {
		buildID = "unknown"
	}
}

func main() {
	log.Printf("-> env:CI_DEBUG %s", os.Getenv("CI_DEBUG"))
	flag.Parse()
	log.Printf("-> run serverPort %v", *serverPort)
	log.Printf("=> now version %v, build id: %s", cliVersion, buildID)

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

	err = r.Run(fmt.Sprintf(":%s", *serverPort))
	if err != nil {
		fmt.Printf("run err %v\n", err)
		return
	}
}
