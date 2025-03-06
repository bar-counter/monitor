[![golang-ci](https://github.com/bar-counter/monitor/actions/workflows/golang-ci.yml/badge.svg)](https://github.com/bar-counter/monitor/actions/workflows/golang-ci.yml)

[![go mod version](https://img.shields.io/github/go-mod/go-version/bar-counter/monitor?label=go.mod)](https://github.com/bar-counter/monitor)
[![GoDoc](https://godoc.org/github.com/bar-counter/monitor?status.png)](https://godoc.org/github.com/bar-counter/monitor)
[![goreportcard](https://goreportcard.com/badge/github.com/bar-counter/monitor)](https://goreportcard.com/report/github.com/bar-counter/monitor)

[![GitHub license](https://img.shields.io/github/license/bar-counter/monitor)](https://github.com/bar-counter/monitor)
[![GitHub latest SemVer tag)](https://img.shields.io/github/v/tag/bar-counter/monitor)](https://github.com/bar-counter/monitor/tags)
[![GitHub release)](https://img.shields.io/github/v/release/bar-counter/monitor)](https://github.com/bar-counter/monitor/releases)

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

## env

- minimum go version: go 1.21
- change `go 1.21`, `^1.21`, `1.23.4` to new go version

### libs

| lib      | url                                 | version |
|:---------|:------------------------------------|:--------|
| gin      | https://github.com/gin-gonic/gin    | v1.9.1  |
| gopsutil | https://github.com/shirou/gopsutil/ | v3.23.5 |

## Contributing

[![Contributor Covenant](https://img.shields.io/badge/contributor%20covenant-v1.4-ff69b4.svg)](.github/CONTRIBUTING_DOC/CODE_OF_CONDUCT.md)
[![GitHub contributors](https://img.shields.io/github/contributors/bar-counter/monitor)](https://github.com/bar-counter/monitor/graphs/contributors)

We welcome community contributions to this project.

Please read [Contributor Guide](.github/CONTRIBUTING_DOC/CONTRIBUTING.md) for more information on how to get started.

请阅读有关 [贡献者指南](.github/CONTRIBUTING_DOC/zh-CN/CONTRIBUTING.md) 以获取更多如何入门的信息

# demo

```bash
make init
make dep
# ensure right then
make exampleDebug
# and open url
# health http://127.0.0.1:38000/status/health
# pprof http://127.0.0.1:38000/debug/pprof/

make exampleStatus
# status http://127.0.0.1:38000/status/hardware/disk
# status http://127.0.0.1:38000/status/hardware/ram
# status http://127.0.0.1:38000/status/hardware/cpu
# status http://127.0.0.1:38000/status/hardware/cpu_info

make examplePprof
# pprof http://127.0.0.1:38000/debug/vars
# pprof http://127.0.0.1:38000/debug/pprof/
```

# use middleware lib

## import

```bash
# go get
go get -v github.com/bar-counter/monitor/v3

# go mod find out version
go list -mod readonly -m -versions github.com/bar-counter/monitor/v3
# all use awk to get script
echo "go mod edit -require=$(go list -m -versions github.com/bar-counter/monitor | awk '{print $1 "@" $NF}')"
# then use your want version like v3.0.0
go mod edit -require=github.com/bar-counter/monitor/v3@v3.0.0
go mod download -x
```

## gin server status

- use lib [https://github.com/shirou/gopsutil](https://github.com/shirou/gopsutil) for watch system status
- [see example statusdemo.go](example/status/statusdemo.go)

```go
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
```

and you can use to get status of server or run `make exampleStatus`

```bash
curl 'http://127.0.0.1:38000/status/health' \
  -X GET

curl 'http://127.0.0.1:38000/status/hardware/disk' \
  -X GET

curl 'http://127.0.0.1:38000/status/hardware/ram' \
  -X GET

curl 'http://127.0.0.1:38000/status/hardware/cpu' \
  -X GET

curl 'http://127.0.0.1:38000/status/hardware/cpu_info' \
  -X GET
```

> StatusPrefix default is `/status` you can change by your self
> StatusHardwarePrefix default is `/hardware`

## gin server debug

- [see example debugdemo.go](example/debug/debugdemo.go)

### vars

```go
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

```

or run `make exampleDebug` because use gin.BasicAuth must add `--user user:user`

```bash
curl 'http://127.0.0.1:38000/debug/vars' \
 --user user:user \
 -X GET
```

> DebugPrefix default is `/debug`

```json
{
  "cgo": 6,
  "cmdline": [
    "/var/folders/79/dw7nb8rx7kgcqty_9qq2nv640000gn/T/go-build2348802398/b001/exe/debugdemo"
  ],
  "gc_pause": 0,
  "go_version": "go1.18.2",
  "goroutine": 5,
  "memstats": {
  },
  "os": "darwin",
  "os_cores": 8,
  "run_time": "8.211807521s"
}
```

| item      | doc                         | desc                  |
|:----------|:----------------------------|:----------------------|
| cgo       | go doc runtime.NumCgoCall   |
| cmdline   |                             | server run cmd        |
| cores     | go doc runtime.NumCPU       |
| gc_pause  |                             | count last gc time    |
| goroutine | go doc runtime.NumGoroutine |                       |
| memstats  | go doc runtime.MemStats     |                       |
| os        | go doc runtime.GOOS         |                       |
| os_cores  | go doc runtime.NumCPU       |                       |
| run_time  |                             | count server run time |

> more info see `go doc expvar`

`DebugMiddleware` can use BasicAuth or other Middleware

### pprof

- pprof must Debug open
- [see example pprofdemo.go](example/pprof/pprofdemo.go)

```go
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
		//DebugMiddleware: gin.BasicAuth(gin.Accounts{
		//	"admin": "admin",
		//	"user":  "user",
		//}),
		PProf: true,
		//PProfPrefix: "/pprof",
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

```

or run `make examplePprof`

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

## License
[![FOSSA Status](https://app.fossa.com/api/projects/git%2Bgithub.com%2Fbar-counter%2Fmonitor.svg?type=large)](https://app.fossa.com/projects/git%2Bgithub.com%2Fbar-counter%2Fmonitor?ref=badge_large)