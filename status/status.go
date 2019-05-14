package status

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/shirou/gopsutil/mem"
	"net/http"

	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/disk"
	"github.com/shirou/gopsutil/load"
)

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
)

// @Summary /status/health
// @Description HealthCheck shows OK as the ping-pong result.
// @Tags status
// @Success 200 "OK"
// @Router /status/health [get]
func HealthCheck(c *gin.Context) {
	message := "OK"
	c.String(http.StatusOK, "\n"+message)
}

type Disk struct {
	Info        string `json:"info,omitempty"`
	UsedMB      int    `json:"used_mb,omitempty"`
	UsedGB      int    `json:"used_gb,omitempty"`
	TotalMB     int    `json:"total_mb,omitempty"`
	TotalGB     int    `json:"total_gb,omitempty"`
	UsedPercent int    `json:"used_percent,omitempty"`
}

// @Summary /status/disk
// @Description HealthCheck DiskCheck checks the disk usage.
// @Tags status
// @Accept application/json
// @Produce application/json
// @Success 200 {object} status.Disk "value in monitor.status.Disk DISK OK less than 80% use, 80% to 90% is WARNING"
// @failure 429 {object} status.Disk "value in monitor.status.Disk DISK need check after 90%"
// @failure 500 "DISK CRITICAL must check!"
// @Router /status/disk [get]
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

// @Summary /status/ram
// @Description HealthCheck RAMCheck checks the disk usage.
// @Tags status
// @Accept application/json
// @Produce application/json
// @Success 200 {object} status.RAM "value in monitor.status.RAM OK less than 80% use, 80% to 95% is WARNING"
// @failure 429 {object} status.RAM "value in monitor.status.RAM need check after 95%"
// @failure 500 "RAM CRITICAL must check!"
// @Router /status/ram [get]
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

// @Summary /status/cpu
// @Description HealthCheck CPUCheck checks the cpu usage.
// @Tags status
// @Success 200 "CPU OK server run ok"
// @failure 429 "CPU WARNING need check"
// @failure 500 "CPU CRITICAL must check"
// @Router /status/cpu [get]
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
