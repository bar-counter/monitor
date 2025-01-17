package status

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/v4/cpu"
	"github.com/shirou/gopsutil/v4/disk"
	"github.com/shirou/gopsutil/v4/load"
	"github.com/shirou/gopsutil/v4/mem"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

// @Summary /status/health
// @Description HealthCheck shows OK as the ping-pong result. must of api not use [ BasePath ]
// @Tags status
// @Success 200 "OK"
// @Router /status/health [get]
// @BasePath /
func HealthCheck(c *gin.Context) {
	message := "OK"
	c.String(http.StatusOK, message)
}

type Disk struct {
	Info        string `json:"info,omitempty"`
	UsedMB      int    `json:"used_mb,omitempty"`
	UsedGB      int    `json:"used_gb,omitempty"`
	TotalMB     int    `json:"total_mb,omitempty"`
	TotalGB     int    `json:"total_gb,omitempty"`
	UsedPercent int    `json:"used_percent,omitempty"`
}

func DiskCheck(c *gin.Context) {
	u, _ := disk.Usage("/")

	usedMB := int(u.Used) / MB
	usedGB := int(u.Used) / GB
	totalMB := int(u.Total) / MB
	totalGB := int(u.Total) / GB
	usedPercent := int(u.UsedPercent)

	status := http.StatusOK
	text := "OK"

	if usedPercent >= 90 {
		status = http.StatusTooManyRequests
		text = "CRITICAL"
	} else if usedPercent >= 80 {
		status = http.StatusOK
		text = "WARNING"
	}
	diskInfo := &Disk{
		Info:        text,
		UsedMB:      usedMB,
		UsedGB:      usedGB,
		TotalMB:     totalMB,
		TotalGB:     totalGB,
		UsedPercent: usedPercent,
	}
	//message := fmt.Sprintf("%s - Free space: %dMB (%dGB) / %dMB (%dGB) | Used: %d%%", text, usedMB, usedGB, totalMB, totalGB, usedPercent)
	//c.String(status, "\n"+message)
	c.JSON(status, diskInfo)
}

type RAM struct {
	Info        string `json:"info,omitempty"`
	UsedMB      int    `json:"used_mb,omitempty"`
	UsedGB      int    `json:"used_gb,omitempty"`
	TotalMB     int    `json:"total_mb,omitempty"`
	TotalGB     int    `json:"total_gb,omitempty"`
	UsedPercent int    `json:"used_percent,omitempty"`
}

func RAMCheck(c *gin.Context) {
	u, _ := mem.VirtualMemory()

	usedMB := int(u.Used) / MB
	usedGB := int(u.Used) / GB
	totalMB := int(u.Total) / MB
	totalGB := int(u.Total) / GB
	usedPercent := int(u.UsedPercent)

	status := http.StatusOK
	text := "OK"

	if usedPercent >= 95 {
		status = http.StatusInternalServerError
		text = "CRITICAL"
	} else if usedPercent >= 80 {
		status = http.StatusTooManyRequests
		text = "WARNING"
	}

	//message := fmt.Sprintf("%s - Free space: %dMB (%dGB) / %dMB (%dGB) | Used: %d%%", text, usedMB, usedGB, totalMB, totalGB, usedPercent)
	//c.String(status, "\n"+message)
	ramInfo := &RAM{
		Info:        text,
		UsedMB:      usedMB,
		UsedGB:      usedGB,
		TotalMB:     totalMB,
		TotalGB:     totalGB,
		UsedPercent: usedPercent,
	}
	c.JSON(status, ramInfo)
}

func CPUCheck(c *gin.Context) {
	cores, _ := cpu.Counts(false)

	a, _ := load.Avg()
	l1 := a.Load1
	l5 := a.Load5
	l15 := a.Load15

	status := http.StatusOK
	text := "OK"

	if l5 >= float64(cores-1) {
		status = http.StatusInternalServerError
		text = "CRITICAL"
	} else if l5 >= float64(cores-2) {
		status = http.StatusTooManyRequests
		text = "WARNING"
	}

	message := fmt.Sprintf("%s - Load average: %.2f, %.2f, %.2f | Cores: %d", text, l1, l5, l15, cores)
	c.String(status, "\n"+message)
}

type CPU struct {
	L1  float64 `json:"l_1,omitempty"`
	L5  float64 `json:"l_5,omitempty"`
	L15 float64 `json:"l_15,omitempty"`

	CpuCnt int    `json:"c_cnt,omitempty"`
	Status string `json:"status"`
}

func CPUInfo(c *gin.Context) {
	cores, _ := cpu.Counts(false)

	a, _ := load.Avg()
	l1 := a.Load1
	l5 := a.Load5
	l15 := a.Load15

	status := http.StatusOK
	text := "OK"

	if l5 >= float64(cores-1) {
		status = http.StatusInternalServerError
		text = "CRITICAL"
	} else if l5 >= float64(cores-2) {
		status = http.StatusTooManyRequests
		text = "WARNING"
	}

	cpuInfo := &CPU{
		L1:     l1,
		L5:     l5,
		L15:    l15,
		CpuCnt: cores,
		Status: text,
	}
	c.JSON(status, cpuInfo)
}
