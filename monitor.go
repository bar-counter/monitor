package monitor

import (
	"github.com/bar-counter/monitor/debug"
	"github.com/bar-counter/monitor/status"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	DefaultStatusPrefix = "/status"
	DefaultDebugPrefix  = "/debug"
	// DefaultVarsPrefix url prefix of /debug/vars
	DefaultVarsPrefix = "/vars"
	// DefaultPrefix url prefix of /debug/pprof
	DefaultPPROFPrefix = "/pprof"
)

type Cfg struct {
	APIBase      string
	Status       bool
	StatusPrefix string
	Debug        bool
	DebugPrefix  string
	VarsPrefix   string
	PProf        bool
	PProfPrefix  string
}

var DefaultCfg *Cfg

func initDefaultCfg() *Cfg {
	cfg := Cfg{
		APIBase:      "",
		Status:       false,
		StatusPrefix: DefaultStatusPrefix,
		Debug:        false,
		DebugPrefix:  DefaultDebugPrefix,
		VarsPrefix:   DefaultVarsPrefix,
		PProf:        false,
		PProfPrefix:  DefaultPPROFPrefix,
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
				statusGroup.GET("/disk", status.DiskCheck)
				statusGroup.GET("/ram", status.RAMCheck)
				statusGroup.GET("/cpu", status.CPUCheck)
			}
		}
		if cfg.Debug {
			mGroup.GET(cfg.DebugPrefix, debug.GetMonitorRunningStats)
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
