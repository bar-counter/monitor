## for what

- this project used to gin api server

# use lib

- https://github.com/shirou/gopsutil for watch system status

```go
	r := gin.Default()
	cfg := &monitor.Cfg{
		Status:       true,
		//StatusPrefix: "/status",
	}
	err := monitor.Register(r, cfg)
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

curl 'http://127.0.0.1:38000/status/disk' \                                                                                                                                                                          [3:34:08]
  -X GET

curl 'http://127.0.0.1:38000/status/ram' \                                                                                                                                                                          [3:34:08]
  -X GET

curl 'http://127.0.0.1:38000/status/cpu' \                                                                                                                                                                          [3:34:08]
  -X GET
```

> StatusPrefix default is `/status` you can change by your self
