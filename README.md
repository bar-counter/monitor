## for what

- this project used to gin api server

# use lib

## gin server status

- https://github.com/shirou/gopsutil for watch system status

```go
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

## debug

### vars

```go
	r := gin.Default()
	monitorCfg := &monitor.Cfg{
		Debug: true,
		//DebugPrefix: "/debug",
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
