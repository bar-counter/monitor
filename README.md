<!-- TOC -->

- [for what](#for-what)
  - [dependInfo](#dependinfo)
- [demo](#demo)
- [use middleware](#use-middleware)
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
| gin | https://github.com/gin-gonic/gin | v1.4.0 |
| gopsutil | https://github.com/shirou/gopsutil | v2.19.10 |

# demo

```bash
make init
make checkDepends
# ensure right then
make dev
# and open url
# health http://127.0.0.1:38000/status/health
# pprof http://127.0.0.1:38000/debug/pprof/
```

# use middleware

```bash
# go get
go get -v github.com/bar-counter/monitor

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
curl 'http://127.0.0.1:38000/status/health' \                                                                                                                                                                          [3:34:08]
  -X GET

curl 'http://127.0.0.1:38000/status/hardware/disk' \                                                                                                                                                                          [3:34:08]
  -X GET

curl 'http://127.0.0.1:38000/status/hardware/ram' \                                                                                                                                                                          [3:34:08]
  -X GET

curl 'http://127.0.0.1:38000/status/hardware/cpu' \                                                                                                                                                                          [3:34:08]
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
curl 'http://127.0.0.1:38000/debug/vars' \                                                                                                                                                                             [4:02:09]
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