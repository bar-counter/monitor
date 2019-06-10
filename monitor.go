package monitor

import (
	"github.com/bar-counter/monitor/debug"
	"github.com/bar-counter/monitor/status"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	DefaultStatusPrefix = "/status"
	// DefaultVarsPrefix url prefix of /status/hardware
	DefaultStatusHardwarePrefix = "/hardware"
	DefaultDebugPrefix          = "/debug"
	// DefaultVarsPrefix url prefix of /debug/vars
	DefaultVarsPrefix = "/vars"
	// DefaultPrefix url prefix of /debug/pprof
	DefaultPPROFPrefix = "/pprof"
)

type Cfg struct {
	APIBase              string
	Status               bool
	StatusPrefix         string
	StatusHardware       bool
	StatusHardwarePrefix string
	Debug                bool
	DebugPrefix          string
	VarsPrefix           string
	PProf                bool
	PProfPrefix          string
}

var DefaultCfg *Cfg

func initDefaultCfg() *Cfg {
	cfg := Cfg{
		APIBase:              "",
		Status:               false,
		StatusPrefix:         DefaultStatusPrefix,
		StatusHardware:       false,
		StatusHardwarePrefix: DefaultStatusHardwarePrefix,
		Debug:                false,
		DebugPrefix:          DefaultDebugPrefix,
		VarsPrefix:           DefaultVarsPrefix,
		PProf:                false,
		PProfPrefix:          DefaultPPROFPrefix,
	}
	return &cfg
}

func Register(r *gin.Engine, cfg *Cfg) error {
	checkCfg(cfg)
	mGroup := r.Group(cfg.APIBase)
	{
		if cfg.Status {
			statusGroup := mGroup.Group(cfg.StatusPrefix)
			{
				statusGroup.GET("/health", status.HealthCheck)

				if cfg.StatusHardware {
					statusHardwareGroup := statusGroup.Group(cfg.StatusHardwarePrefix)
					{
						statusHardwareGroup.GET("/disk", status.DiskCheck)
						statusHardwareGroup.GET("/ram", status.RAMCheck)
						statusHardwareGroup.GET("/cpu", status.CPUCheck)
					}
				}
			}
		}
		if cfg.Debug {
			debugGroup := mGroup.Group(cfg.DebugPrefix)
			{
				debugGroup.GET(cfg.VarsPrefix, debug.GetMonitorRunningStats)
			}
		}
	}
	return nil
}

func checkCfg(cfg *Cfg) {
	defaultCfg := initDefaultCfg()
	if cfg.APIBase == "" {
		cfg.APIBase = defaultCfg.APIBase
	}
	if cfg.StatusPrefix == "" {
		cfg.StatusPrefix = defaultCfg.StatusPrefix
	}
	if cfg.StatusHardwarePrefix == "" {
		cfg.StatusHardwarePrefix = defaultCfg.StatusHardwarePrefix
	}
	if cfg.DebugPrefix == "" {
		cfg.DebugPrefix = defaultCfg.DebugPrefix
	}
	if cfg.VarsPrefix == "" {
		cfg.VarsPrefix = defaultCfg.VarsPrefix
	}
	if cfg.PProfPrefix == "" {
		cfg.PProfPrefix = defaultCfg.PProfPrefix
	}
	DefaultCfg = cfg
}

func handler(h http.HandlerFunc) gin.HandlerFunc {
	handler := http.HandlerFunc(h)
	return func(c *gin.Context) {
		handler.ServeHTTP(c.Writer, c.Request)
	}
}
