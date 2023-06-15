package debug

import (
	"encoding/json"
	"expvar"
	"fmt"
	"net/http"
	"runtime"
	"time"

	"github.com/gin-gonic/gin"
)

// PublishExpVarDebug
// let full vars publish at expvar
func PublishExpVarDebug() {
	expvar.Publish("cgo", expvar.Func(getNumCgoCall))
	expvar.Publish("gc_pause", expvar.Func(getLastGCPauseTime))
	expvar.Publish("go_version", expvar.Func(currentGoVersion))
	expvar.Publish("goroutine", expvar.Func(getNumGoroutine))
	expvar.Publish("os", expvar.Func(getGoOS))
	expvar.Publish("os_cores", expvar.Func(getNumCPU))
	expvar.Publish("run_time", expvar.Func(calculateUptime))
}

func GetMonitorRunningStats(c *gin.Context) {
	c.Writer.Header().Set("Content-Type", "application/json; charset=utf-8")

	first := true
	report := func(key string, value interface{}) {
		if !first {
			_, _ = fmt.Fprintf(c.Writer, ",\n")
		}
		first = false
		if str, ok := value.(string); ok {
			_, _ = fmt.Fprintf(c.Writer, "%q: %q", key, str)
		} else {
			_, _ = fmt.Fprintf(c.Writer, "%q: %v", key, value)
		}
	}

	_, _ = fmt.Fprintf(c.Writer, "{\n")
	expvar.Do(func(kv expvar.KeyValue) {
		report(kv.Key, kv.Value)
	})
	_, _ = fmt.Fprintf(c.Writer, "\n}\n")

	c.String(http.StatusOK, "")
}

// server start time
var start = time.Now()

// calculateUptime calculate uptime
func calculateUptime() interface{} {
	return time.Since(start).String()
}

// currentGoVersion now golang version
func currentGoVersion() interface{} {
	return runtime.Version()
}

// get Number of CPU core
func getNumCPU() interface{} {
	return runtime.NumCPU()
}

// get OS name
func getGoOS() interface{} {
	return runtime.GOOS
}

// get Number of goroutine
func getNumGoroutine() interface{} {
	return runtime.NumGoroutine()
}

// get Number of CgoCall
func getNumCgoCall() interface{} {
	return runtime.NumCgoCall()
}

// last pause time for record pause time
var lastPause uint32

// get last GC Pause Time
func getLastGCPauseTime() interface{} {
	var gcPause uint64
	ms := new(runtime.MemStats)

	statString := expvar.Get("memstats").String()
	if statString != "" {
		_ = json.Unmarshal([]byte(statString), ms)

		if lastPause == 0 || lastPause != ms.NumGC {
			gcPause = ms.PauseNs[(ms.NumGC+255)%256]
			lastPause = ms.NumGC
		}
	}

	return gcPause
}
