## for what

- this project used to gin api server

# use lib

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
