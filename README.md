[![golang-full](https://github.com/bar-counter/monitor/workflows/golang-full/badge.svg?branch=main)](https://github.com/bar-counter/monitor/actions)
[![TravisBuildStatus](https://api.travis-ci.org/bar-counter/monitor.svg?branch=master)](https://travis-ci.org/bar-counter/monitor)
[![GoDoc](https://godoc.org/github.com/bar-counter/monitor?status.png)](https://godoc.org/github.com/bar-counter/monitor/)
[![GoReportCard](https://goreportcard.com/badge/github.com/bar-counter/monitor)](https://goreportcard.com/report/github.com/bar-counter/monitor)
[![codecov](https://codecov.io/gh/bar-counter/monitor/branch/master/graph/badge.svg)](https://codecov.io/gh/bar-counter/monitor)

<!-- TOC -->

- [for what](#for-what)
  - [dependInfo](#dependinfo)
- [demo](#demo)
- [use middleware lib](#use-middleware-lib)
  - [import](#import)
  - [gin server status](#gin-server-status)
  - [gin server debug](#gin-server-debug)
    - [vars](#vars)
    - [pprof](#pprof)

<!-- /TOC -->

# for what

- this project used to gin api server status monitor

support check
- `health`
- `Hardware`
	- disk
	- cpu
	- ram
- `debug`
	- run vars
	- pprof

## dependInfo

| lib | url | version |
|:-----|:-----|:-----|
| gin | https://github.com/gin-gonic/gin | v1.6.3 |
| gopsutil | https://github.com/shirou/gopsutil | v2.20.9+incompatible |
| go-ole | https://github.com/go-ole/go-ole | v1.2.5 |

# demo

```bash
make init
make dep
# ensure right then
make dev
# and open url
# health http://127.0.0.1:38000/status/health
# pprof http://127.0.0.1:38000/debug/pprof/
```

# use middleware lib

## import

```bash
# go get
go get -v github.com/bar-counter/monitor

# go mod find out verison
go list -m -versions github.com/bar-counter/monitor
# all use awk to get script
echo "go mod edit -require=$(go list -m -versions github.com/bar-counter/monitor | awk '{print $1 "@" $NF}')"
# then use your want verison like v1.1.0
GO111MODULE=on go mod edit -require=github.com/bar-counter/monitor@v1.1.0
GO111MODULE=on go mod vendor

# dep go 1.7 -> 1.11
dep ensure --add github.com/bar-counter/monitor@1.0.1
dep ensure -v
```

## gin server status

- use lib [https://github.com/shirou/gopsutil](https://github.com/shirou/gopsutil) for watch system status
- [see example statusdemo.go](example/status/statusdemo.go)

```go
import (
	"github.com/bar-counter/monitor"
)

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
	r.Run(":38000")
```

and you can use to get status of server

```bash
curl 'http://127.0.0.1:38000/status/health' \
  -X GET

curl 'http://127.0.0.1:38000/status/hardware/disk' \
  -X GET

curl 'http://127.0.0.1:38000/status/hardware/ram' \
  -X GET

curl 'http://127.0.0.1:38000/status/hardware/cpu' \
  -X GET
```

> StatusPrefix default is `/status` you can change by your self
> StatusHardwarePrefix default is `/hardware`

## gin server debug

- [see example debugdemo.go](example/debug/debugdemo.go)

### vars

```go
import (
	"github.com/bar-counter/monitor"
)

	r := gin.Default()
	monitorCfg := &monitor.Cfg{
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
	r.Run(":38000")
```

```bash
curl 'http://127.0.0.1:38000/debug/vars' \
-X GET
```

> DebugPrefix default is `/debug`

```json
{
    "cgo": 6,
    "cmdline": [
        "/private/var/folders/3t/q6knzmrs09b0m5px2dljnpnr0000gn/T/___go_build_main_go__1_"
    ],
    "gc_pause": 0,
    "go_version": "go1.11.4",
    "goroutine": 3,
    "memstats": {
    },
    "os": "darwin",
    "os_cores": 8,
    "run_time": "8.211807521s"
}
```

| item | doc  | desc |
|:-----|:-----|:-----|
| cgo | go doc runtime.NumCgoCall |
| cmdline |  | server run cmd |
| cores | go doc runtime.NumCPU |
| gc_pause | | count last gc time |
| goroutine | go doc runtime.NumGoroutine| |
| memstats | go doc runtime.MemStats | |
| os | go doc runtime.GOOS | |
| os_cores | go doc runtime.NumCPU | |
| run_time | | count server run time |

> more info see `go doc expvar`

`DebugMiddleware` can use BasicAuth or other Middleware

### pprof

- pprof must Debug open
- [see example pprofdemo.go](example/pprof/pprofdemo.go)

```go
import (
	"github.com/bar-counter/monitor"
)

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
		PProf: true,
		//PProfPrefix: "/pprof",
	}
	err := monitor.Register(r, monitorCfg)
	if err != nil {
		fmt.Printf("monitor register err %v\n", err)
		return
	}
	r.Run(":38000")
```

then see at [http://127.0.0.1:38000/debug/pprof/](http://127.0.0.1:38000/debug/pprof/)
or use cli to

```bash
# cpu
go tool pprof http://localhost:38000/debug/pprof/profile
# mem
go tool pprof http://localhost:38000/debug/pprof/heap
# block
go tool pprof http://localhost:38000/debug/pprof/block
```