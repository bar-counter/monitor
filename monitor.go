package monitor

import (
	"github.com/bar-counter/monitor/status"
	"github.com/gin-gonic/gin"
	"net/http"
)

const (
	// DefaultPrefix url prefix of pprof
	DefaultStatusPrefix = "/status"
	DefaultDebugPrefix  = "/debug"
	DefaultVarsPrefix   = "/vars"
	DefaultPPROFPrefix  = "/pprof"
)

type Cfg struct {
	APIVersion   string
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
		APIVersion:   "/v1",
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
	mGroup := r.Group(cfg.APIVersion)
	{
		if cfg.Status {
			mGroup.GET(cfg.StatusPrefix+"/health", status.HealthCheck)
			mGroup.GET(cfg.StatusPrefix+"/disk", status.DiskCheck)
			mGroup.GET(cfg.StatusPrefix+"/ram", status.RAMCheck)
			mGroup.GET(cfg.StatusPrefix+"/cpu", status.CPUCheck)
		}
		if cfg.Debug {

		}
	}
	return nil
}

func checkCfg(cfg *Cfg) {
	defaultCfg := initDefaultCfg()
	if cfg.APIVersion != "" {
		cfg.APIVersion = defaultCfg.APIVersion
	}
	if cfg.StatusPrefix != "" {
		cfg.StatusPrefix = defaultCfg.StatusPrefix
	}
	if cfg.DebugPrefix != "" {
		cfg.DebugPrefix = defaultCfg.DebugPrefix
	}
	if cfg.VarsPrefix != "" {
		cfg.VarsPrefix = defaultCfg.VarsPrefix
	}
	if cfg.PProfPrefix != "" {
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
